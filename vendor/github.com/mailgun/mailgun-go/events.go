package mailgun

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// Events are open-ended, loosely-defined JSON documents.
// They will always have an event and a timestamp field, however.
type Event map[string]interface{}

// Parse the timestamp field for this event into a Time object
func (e Event) ParseTimeStamp() (time.Time, error) {
	obj, ok := e["timestamp"]
	if !ok {
		return time.Time{}, errors.New("'timestamp' field not found in event")
	}
	timestamp, ok := obj.(float64)
	if !ok {
		return time.Time{}, errors.New("'timestamp' field not a float64")
	}
	microseconds := int64(timestamp * 1000000)
	return time.Unix(0, microseconds*int64(time.Microsecond/time.Nanosecond)).UTC(), nil
}

func (e Event) ParseMessageId() (string, error) {
	message, err := toMapInterface("message", e)
	if err != nil {
		return "", err
	}
	headers, err := toMapInterface("headers", message)
	if err != nil {
		return "", err
	}
	return headers["message-id"].(string), nil
}

// noTime always equals an uninitialized Time structure.
// It's used to detect when a time parameter is provided.
var noTime time.Time

// GetEventsOptions lets the caller of GetEvents() specify how the results are to be returned.
// Begin and End time-box the results returned.
// ForceAscending and ForceDescending are used to force Mailgun to use a given traversal order of the events.
// If both ForceAscending and ForceDescending are true, an error will result.
// If none, the default will be inferred from the Begin and End parameters.
// Limit caps the number of results returned.  If left unspecified, Mailgun assumes 100.
// Compact, if true, compacts the returned JSON to minimize transmission bandwidth.
// Otherwise, the JSON is spaced appropriately for human consumption.
// Filter allows the caller to provide more specialized filters on the query.
// Consult the Mailgun documentation for more details.
type EventsOptions struct {
	Begin, End                               time.Time
	ForceAscending, ForceDescending, Compact bool
	Limit                                    int
	Filter                                   map[string]string
	ThresholdAge                             time.Duration
	PollInterval                             time.Duration
}

// Depreciated See `ListEvents()`
type GetEventsOptions struct {
	Begin, End                               time.Time
	ForceAscending, ForceDescending, Compact bool
	Limit                                    int
	Filter                                   map[string]string
}

// EventIterator maintains the state necessary for paging though small parcels of a larger set of events.
type EventIterator struct {
	events                              []Event
	NextURL, PrevURL, FirstURL, LastURL string
	mg                                  Mailgun
	err                                 error
}

// NewEventIterator creates a new iterator for events.
// Use GetFirstPage to retrieve the first batch of events.
// Use GetNext and GetPrevious thereafter as appropriate to iterate through sets of data.
//
// *This call is Deprecated, use ListEvents() instead*
func (mg *MailgunImpl) NewEventIterator() *EventIterator {
	return &EventIterator{mg: mg}
}

// Create an new iterator to fetch a page of events from the events api
//	it := mg.ListEvents(EventsOptions{})
//	var events []Event
//	for it.Next(&events) {
//	    	for _, event := range events {
//		        // Do things with events
//		}
//	}
//	if it.Err() != nil {
//		log.Fatal(it.Err())
//	}
func (mg *MailgunImpl) ListEvents(opts *EventsOptions) *EventIterator {
	req := newHTTPRequest(generateApiUrl(mg, eventsEndpoint))
	if opts != nil {
		if opts.Limit != 0 {
			req.addParameter("limit", fmt.Sprintf("%d", opts.Limit))
		}
		if opts.Compact {
			req.addParameter("pretty", "no")
		}
		if opts.ForceAscending {
			req.addParameter("ascending", "yes")
		}
		if opts.ForceDescending {
			req.addParameter("ascending", "no")
		}
		if opts.Begin != noTime {
			req.addParameter("begin", formatMailgunTime(&opts.Begin))
		}
		if opts.End != noTime {
			req.addParameter("end", formatMailgunTime(&opts.End))
		}
		if opts.Filter != nil {
			for k, v := range opts.Filter {
				req.addParameter(k, v)
			}
		}
	}
	url, err := req.generateUrlWithParameters()
	return &EventIterator{
		mg:       mg,
		NextURL:  url,
		FirstURL: url,
		PrevURL:  "",
		err:      err,
	}
}

// If an error occurred during iteration `Err()` will return non nil
func (ei *EventIterator) Err() error {
	return ei.err
}

// Events returns the most recently retrieved batch of events.
// The length is guaranteed to fall between 0 and the limit set in the GetEventsOptions structure passed to GetFirstPage.
func (ei *EventIterator) Events() []Event {
	return ei.events
}

// GetFirstPage retrieves the first batch of events, according to your criteria.
// See the GetEventsOptions structure for more details on how the fields affect the data returned.
func (ei *EventIterator) GetFirstPage(opts GetEventsOptions) error {
	if opts.ForceAscending && opts.ForceDescending {
		return fmt.Errorf("collation cannot at once be both ascending and descending")
	}

	payload := newUrlEncodedPayload()
	if opts.Limit != 0 {
		payload.addValue("limit", fmt.Sprintf("%d", opts.Limit))
	}
	if opts.Compact {
		payload.addValue("pretty", "no")
	}
	if opts.ForceAscending {
		payload.addValue("ascending", "yes")
	}
	if opts.ForceDescending {
		payload.addValue("ascending", "no")
	}
	if opts.Begin != noTime {
		payload.addValue("begin", formatMailgunTime(&opts.Begin))
	}
	if opts.End != noTime {
		payload.addValue("end", formatMailgunTime(&opts.End))
	}
	if opts.Filter != nil {
		for k, v := range opts.Filter {
			payload.addValue(k, v)
		}
	}

	url, err := generateParameterizedUrl(ei.mg, eventsEndpoint, payload)
	if err != nil {
		return err
	}
	return ei.fetch(url)
}

// Retrieves the chronologically previous batch of events, if any exist.
// You know you're at the end of the list when len(Events())==0.
func (ei *EventIterator) GetPrevious() error {
	return ei.fetch(ei.PrevURL)
}

// Retrieves the chronologically next batch of events, if any exist.
// You know you're at the end of the list when len(Events())==0.
func (ei *EventIterator) GetNext() error {
	return ei.fetch(ei.NextURL)
}

// Retrieves the next page of events from the api. Returns false when there
// no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error
func (ei *EventIterator) Next(events *[]Event) bool {
	if ei.err != nil {
		return false
	}
	ei.err = ei.fetch(ei.NextURL)
	if ei.err != nil {
		return false
	}
	*events = ei.events
	if len(ei.events) == 0 {
		return false
	}
	return true
}

// Retrieves the first page of events from the api. Returns false if there
// was an error. It also sets the iterator object to the first page.
// Use `.Err()` to retrieve the error.
func (ei *EventIterator) First(events *[]Event) bool {
	if ei.err != nil {
		return false
	}
	ei.err = ei.fetch(ei.FirstURL)
	if ei.err != nil {
		return false
	}
	*events = ei.events
	return true
}

// Retrieves the last page of events from the api.
// Calling Last() is invalid unless you first call First() or Next()
// Returns false if there was an error. It also sets the iterator object
// to the last page. Use `.Err()` to retrieve the error.
func (ei *EventIterator) Last(events *[]Event) bool {
	if ei.err != nil {
		return false
	}
	ei.err = ei.fetch(ei.LastURL)
	if ei.err != nil {
		return false
	}
	*events = ei.events
	return true
}

// Retrieves the previous page of events from the api. Returns false when there
// no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error if any
func (ei *EventIterator) Previous(events *[]Event) bool {
	if ei.err != nil {
		return false
	}
	if ei.PrevURL == "" {
		return false
	}
	ei.err = ei.fetch(ei.PrevURL)
	if ei.err != nil {
		return false
	}
	*events = ei.events
	if len(ei.events) == 0 {
		return false
	}
	return true
}

// EventPoller maintains the state necessary for polling events
type EventPoller struct {
	it            *EventIterator
	opts          EventsOptions
	thresholdTime time.Time
	sleepUntil    time.Time
	mg            Mailgun
	err           error
}

// Poll the events api and return new events as they occur
// 	it = mg.PollEvents(&EventsOptions{
//			// Poll() returns after this threshold is met, or events older than this threshold appear
// 			ThresholdAge: time.Second * 10,
//			// Only events with a timestamp after this date/time will be returned
//			Begin:        time.Now().Add(time.Second * -3),
//			// How often we poll the api for new events
//			PollInterval: time.Second * 4})
//	var events []Event
//	// Blocks until new events appear
//	for it.Poll(&events) {
//		for _, event := range(events) {
//			fmt.Printf("Event %+v\n", event)
//		}
//	}
//	if it.Err() != nil {
//		log.Fatal(it.Err())
//	}
func (mg *MailgunImpl) PollEvents(opts *EventsOptions) *EventPoller {
	now := time.Now()
	// ForceAscending must be set
	opts.ForceAscending = true

	// Default begin time is 30 minutes ago
	if opts.Begin == noTime {
		opts.Begin = now.Add(time.Minute * -30)
	}

	// Default threshold age is 30 minutes
	if opts.ThresholdAge.Nanoseconds() == 0 {
		opts.ThresholdAge = time.Duration(time.Minute * 30)
	}

	// Set a 15 second poll interval if none set
	if opts.PollInterval.Nanoseconds() == 0 {
		opts.PollInterval = time.Duration(time.Second * 15)
	}

	return &EventPoller{
		it:   mg.ListEvents(opts),
		opts: *opts,
		mg:   mg,
	}
}

// If an error occurred during polling `Err()` will return non nil
func (ep *EventPoller) Err() error {
	return ep.err
}

func (ep *EventPoller) Poll(events *[]Event) bool {
	var currentPage string
	ep.thresholdTime = time.Now().UTC().Add(ep.opts.ThresholdAge)
	for {
		if ep.sleepUntil != noTime {
			// Sleep the rest of our duration
			time.Sleep(ep.sleepUntil.Sub(time.Now()))
		}

		// Remember our current page url
		currentPage = ep.it.NextURL

		// Attempt to get a page of events
		var page []Event
		if ep.it.Next(&page) == false {
			if ep.it.Err() == nil && len(page) == 0 {
				// No events, sleep for our poll interval
				ep.sleepUntil = time.Now().Add(ep.opts.PollInterval)
				continue
			}
			ep.err = ep.it.Err()
			return false
		}

		// Last event on the page
		lastEvent := page[len(page)-1]

		timeStamp, err := lastEvent.ParseTimeStamp()
		if err != nil {
			ep.err = errors.Wrap(err, "event timestamp error")
			return false
		}
		// Record the next time we should query for new events
		ep.sleepUntil = time.Now().Add(ep.opts.PollInterval)

		// If the last event on the page is older than our threshold time
		// or we have been polling for longer than our threshold time
		if timeStamp.After(ep.thresholdTime) || time.Now().UTC().After(ep.thresholdTime) {
			ep.thresholdTime = time.Now().UTC().Add(ep.opts.ThresholdAge)
			// Return the page of events to the user
			*events = page
			return true
		}
		// Since we didn't find an event older than our
		// threshold, fetch this same page again
		ep.it.NextURL = currentPage
	}
}

// GetFirstPage, GetPrevious, and GetNext all have a common body of code.
// fetch completes the API fetch common to all three of these functions.
func (ei *EventIterator) fetch(url string) error {
	r := newHTTPRequest(url)
	r.setClient(ei.mg.Client())
	r.setBasicAuth(basicAuthUser, ei.mg.ApiKey())
	var response map[string]interface{}
	err := getResponseFromJSON(r, &response)
	if err != nil {
		return err
	}

	items := response["items"].([]interface{})
	ei.events = make([]Event, len(items))
	for i, item := range items {
		ei.events[i] = item.(map[string]interface{})
	}

	pagings := response["paging"].(map[string]interface{})
	links := make(map[string]string, len(pagings))
	for key, page := range pagings {
		links[key] = page.(string)
	}
	ei.NextURL = links["next"]
	ei.PrevURL = links["previous"]
	ei.FirstURL = links["first"]
	ei.LastURL = links["last"]
	return err
}

func toMapInterface(field string, thingy map[string]interface{}) (map[string]interface{}, error) {
	var empty map[string]interface{}
	obj, ok := thingy[field]
	if !ok {
		return empty, errors.Errorf("'%s' field not found in event", field)
	}
	result, ok := obj.(map[string]interface{})
	if !ok {
		return empty, errors.Errorf("'%s' field not a map[string]interface{}", field)
	}
	return result, nil
}
