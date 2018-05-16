//go:generate mockgen -destination ../mock_stride/mock_stride.go bitbucket.org/atlassian/go-stride/pkg/stride Client
package stride

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

const (
	apiBaseURL     = "https://api.atlassian.com"
	apiAudience    = "api.atlassian.com"
	authAPIBaseURL = "https://auth.atlassian.com"
)

type Action string
type BadgeState string
type BadgeTheme string
type LozengeState string

const (
	// Actions
	ActionAdd       = Action("add")
	ActionArchive   = Action("archive")
	ActionCreate    = Action("create")
	ActionDelete    = Action("delete")
	ActionRemove    = Action("remove")
	ActionUnarchive = Action("unarchive")
	ActionUpdate    = Action("update")
	// Badges
	BadgeAdded        = BadgeState("added")
	BadgeDefault      = BadgeState("default")
	BadgeImportant    = BadgeState("important")
	BadgePrimary      = BadgeState("primary")
	BadgeRemoved      = BadgeState("removed")
	BadgeThemeDefault = BadgeTheme("default")
	BadgeThemeDark    = BadgeTheme("dark")
	// Conditions
	ConditionRoomPublic    = "room_is_public"
	ConditionUserAdmin     = "user_is_admin"
	ConditionUserGuest     = "user_is_guest"
	ConditionUserRoomOwner = "user_is_room_owner"
	// Events
	EventConversationUpdates = "conversation:updates"
	EventRosterUpdates       = "roster:updates"
	// Lozenges
	LozengeDefault    = LozengeState("default")
	LozengeInProgress = LozengeState("inprogress")
	LozengeMoved      = LozengeState("moved")
	LozengeNew        = LozengeState("new")
	LozengeRemoved    = LozengeState("removed")
	LozengeSuccess    = LozengeState("success")
	// Scopes
	ScopeParticipateConversation = "participate:conversation"
	ScopeManageConversation      = "manage:conversation"
)

type clientImpl struct {
	clientID     string
	clientSecret string
	// Protect access to token.
	mutex        *sync.Mutex
	token        string
	tokenExpires time.Time
	// Override URLs for testing.
	apiBaseURL     string
	authAPIBaseURL string
	httpClient     HttpClient
}

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
	Post(string, string, io.Reader) (*http.Response, error)
}

func New(clientID, clientSecret string) Client {
	return NewClient(clientID, clientSecret, http.DefaultClient)
}

func NewClient(clientID, clientSecret string, httpClient HttpClient) Client {
	return &clientImpl{
		clientID:       clientID,
		clientSecret:   clientSecret,
		mutex:          &sync.Mutex{},
		apiBaseURL:     apiBaseURL,
		authAPIBaseURL: authAPIBaseURL,
		httpClient:     httpClient,
	}
}

// NewRoomClient returns a Client configured with a room specific token.
func NewRoomClient(roomToken string, httpClient HttpClient) Client {
	return &clientImpl{
		token:          roomToken,
		tokenExpires:   time.Now().Add(time.Hour * 24 * 365 * 10),
		mutex:          &sync.Mutex{},
		apiBaseURL:     apiBaseURL,
		authAPIBaseURL: authAPIBaseURL,
		httpClient:     httpClient,
	}
}

type Client interface {
	SendMessage(cloudID, conversationID string, payload *Payload) error
	SendMarkdown(cloudID, conversationID string, markdown string) error
	SendUserMessage(cloudID, userID string, payload *Payload) error
	GetConversation(cloudID, conversationID string) (*ConversationCommon, error)
	GetConversations(cloudID string) ([]*ConversationCommon, error)
	GetConversationByName(cloudID, conversationName string) (*ConversationCommon, error)
	GetUser(cloudID, userID string) (user *User, err error)
	ConvertDocToText(doc *Document) (plain string, err error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
}

// Deprecated: This function will be removed in a future version. Use ParagraphDocument instead.
func PlainText(text string) *Payload {
	return ParagraphDocument(text)
}

func ParagraphDocument(text string) *Payload {
	return &Payload{
		Body: &Document{
			Version: 1,
			Content: []interface{}{
				&Paragraph{
					Content: []InlineGroupNode{
						&Text{
							Text: text,
						},
					},
				},
			},
		},
	}
}

func ApplicationCardDocument(text, title, description, previewURL string) *Payload {
	c := &ApplicationCard{
		&ApplicationCardAttrs{
			Text: text,
			Title: CardAttrTitle{
				Text: title,
			},
		},
	}
	if description != "" {
		c.Attrs.Description = &CardAttrText{
			Text: description,
		}
	}
	if previewURL != "" {
		c.Attrs.Preview = &URL{
			URL: previewURL,
		}
	}

	return &Payload{
		Body: &Document{
			Version: 1,
			Content: []interface{}{
				c,
			},
		},
	}
}

func SendText(client Client, cloudID, conversationID, text string) error {
	return client.SendMessage(cloudID, conversationID, ParagraphDocument(text))
}

func SendCard(client Client, cloudID, conversationID, text, title, description, previewURL string) error {
	return client.SendMessage(cloudID, conversationID, ApplicationCardDocument(text, title, description, previewURL))
}

func (c *clientImpl) postToAPI(url string, obj interface{}) error {
	accessToken, err := c.GetAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get access token: %v", err)
	}

	b, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %v", err)
	}

	body := bytes.NewReader(b)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Add("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 204 {
		return fmt.Errorf("unable to send message: bad response code (%d) from Stride", resp.StatusCode)
	}
	return nil
}

func (c *clientImpl) postMarkdownToAPI(url string, markdown string) error {
	accessToken, err := c.GetAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get access token: %v", err)
	}
	body := strings.NewReader(markdown)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Add("Content-Type", "text/markdown")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 204 {
		return fmt.Errorf("unable to send message: bad response code (%d) from Stride", resp.StatusCode)
	}
	return nil
}

func (c *clientImpl) SendMessage(cloudID, conversationID string, payload *Payload) error {
	url := fmt.Sprintf("%s/site/%s/conversation/%s/message", c.apiBaseURL, cloudID, conversationID)
	err := c.postToAPI(url, payload)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Trying to contact cloud [ %s ] conversation [ %s ]", cloudID, conversationID))
	}
	return err
}

func (c *clientImpl) SendMarkdown(cloudID, conversationID, markdown string) error {
	url := fmt.Sprintf("%s/site/%s/conversation/%s/message", c.apiBaseURL, cloudID, conversationID)
	err := c.postMarkdownToAPI(url, markdown)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Trying to contact cloud [ %s ] conversation [ %s ]", cloudID, conversationID))
	}
	return err
}

func (c *clientImpl) SendUserMessage(cloudID, userID string, payload *Payload) error {
	url := fmt.Sprintf("%s/site/%s/conversation/user/%s/message", c.apiBaseURL, cloudID, userID)
	return c.postToAPI(url, payload)
}

func (c *clientImpl) GetConversation(cloudID, conversationID string) (*ConversationCommon, error) {
	url := fmt.Sprintf("%s/site/%s/conversation/%s", c.apiBaseURL, cloudID, conversationID)
	accessToken, err := c.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %v", err)
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	conv := &ConversationCommon{}
	json.Unmarshal(body, conv)
	return conv, nil
}

// GetConversations returns all conversations this app has access to.
func (c *clientImpl) GetConversations(cloudID string) ([]*ConversationCommon, error) {
	url := fmt.Sprintf("%s/site/%s/conversation", c.apiBaseURL, cloudID)
	accessToken, err := c.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %v", err)
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	tmp := struct {
		Values []*ConversationCommon
	}{}
	json.Unmarshal(body, &tmp)
	return tmp.Values, nil
}

// GetConversationByName finds and returns a conversation by its name.
// It is possible that the conversation name may be a substring of more rooms than
// the query limit in which case it may not be found.
func (c *clientImpl) GetConversationByName(cloudID, conversationName string) (*ConversationCommon, error) {
	url := fmt.Sprintf("%s/site/%s/conversation", c.apiBaseURL, cloudID)
	accessToken, err := c.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %v", err)
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	// The query parameter limits conversations to those with this substring in the name.
	q := req.URL.Query()
	q.Add("query", conversationName)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	conv := &Conversations{}
	err = json.Unmarshal(body, conv)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal conversations: %v", err)
	}
	for _, c := range conv.Values {
		if strings.EqualFold(c.Name, conversationName) {
			return c, nil
		}
	}
	return nil, nil
}

//   /*
//    * https://developer.atlassian.com/cloud/stride/apis/rest/#api-site-cloudId-conversation-post
//    */
//   function createConversation(token, cloudId, name, privacy, topic, callback) {
//     var body = {
//       name: name,
//       privacy: privacy,
//       topic: topic
//     }
//     var options = {
//       uri: API_BASE_URL + '/site/' + cloudId + '/conversation',
//       method: 'POST',
//       headers: {
//         authorization: "Bearer " + token,
//         "cache-control": "no-cache"
//       },
//       json: body
//     }
//     request(options, function (err, response, body) {
//       callback(err, body);
//     });
//   }
//
//   /**
//    * https://developer.atlassian.com/cloud/stride/apis/rest/#api-site-cloudId-conversation-conversationId-archive-put
//    */
//   function archiveConversation(token, cloudId, conversationId, callback) {
//
//     var options = {
//       uri: API_BASE_URL + '/site/' + cloudId + '/conversation/' + conversationId + '/archive',
//       method: 'PUT',
//       headers: {
//         authorization: "Bearer " + token,
//         "cache-control": "no-cache"
//       }
//     }
//     request(options, function (err, response, body) {
//       callback(err, body);
//     });
//   }
//
//   /**
//    * https://developer.atlassian.com/cloud/stride/apis/rest/#api-site-cloudId-conversation-conversationId-message-get
//    */
//   function getConversationHistory(token, cloudId, conversationId, callback) {
//     var options = {
//       uri: API_BASE_URL + '/site/' + cloudId + '/conversation/' + conversationId + "/message?limit=5",
//       method: 'GET',
//       headers: {
//         authorization: "Bearer " + token,
//         "cache-control": "no-cache"
//       }
//     }
//     request(options, function (err, response, body) {
//       callback(err, JSON.parse(body));
//     });
//   }
//
//   /**
//    * https://developer.atlassian.com/cloud/stride/apis/rest/#api-site-cloudId-conversation-conversationId-roster-get
//    */
//   function getConversationRoster(token, cloudId, conversationId, callback) {
//     var options = {
//       uri: API_BASE_URL + '/site/' + cloudId + '/conversation/' + conversationId + "/roster",
//       method: 'GET',
//       headers: {
//         authorization: "Bearer " + token,
//         "cache-control": "no-cache"
//       }
//     }
//     request(options, function (err, response, body) {
//       callback(err, JSON.parse(body));
//     });
//   }
//
//   /**
//    * Send a file to a conversation. you can then include this file when sending a message
//    */
//   function sendMedia(cloudId, conversationId, name, stream, callback) {
//     getAccessToken(function (err, accessToken) {
//       if (err) {
//         callback(err);
//       } else {
//         var options = {
//           uri: API_BASE_URL + '/site/' + cloudId + '/conversation/' + conversationId + '/media?name=' + name,
//           method: 'POST',
//           headers: {
//             authorization: "Bearer " + accessToken,
//             'content-type': 'application/octet-stream'
//           },
//           body: stream
//         }
//         request(options, function (err, response, body) {
//           console.log("upload file: " + response.statusCode)
//           callback(err, body);
//         });
//       }
//     });
//   }
//

//
//   function updateConfigurationState(cloudId, conversationId, configKey, state, callback) {
//     getAccessToken(function (err, accessToken) {
//       if (err) {
//         callback(err);
//       } else {
//         var uri = API_BASE_URL + '/app/module/chat/conversation/chat:configuration/' + configKey + '/state';
//         console.log(uri);
//         var options = {
//           uri: uri,
//           method: 'POST',
//           headers: {
//             authorization: "Bearer " + accessToken,
//             "cache-control": "no-cache"
//           },
//           json: {
//             "context": {
//               "cloudId": cloudId,
//               "conversationId": conversationId
//             },
//             "configured": state
//           }
//         }
//
//         request(options, function (err, response, body) {
//           console.log(response.statusCode);
//           console.log(response.body);
//           callback(err, body);
//         });
//       }
//     });
//   }

/*
   Atlassian Users API
*/

func (c *clientImpl) GetUser(cloudID, userID string) (user *User, err error) {
	url := fmt.Sprintf("%s/scim/site/%s/Users/%s", c.apiBaseURL, cloudID, userID)
	accessToken, err := c.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %v", err)
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	if resp.StatusCode != 200 {
		// TODO: Look at other status codes and return a better response.
		return nil, fmt.Errorf("could not get user")
	}
	u := &User{}
	if err = json.Unmarshal(b, u); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %v", err)
	}
	return u, nil
}

/*
   Utility functions
*/

func SendReply(c Client, m *Message, reply *Payload) error {
	return c.SendMessage(m.CloudID, m.Conversation.ID, reply)
}

func SendTextReply(c Client, m *Message, reply string) error {
	return SendText(c, m.CloudID, m.Conversation.ID, reply)
}

/*
   Convert an Atlassian document to plain text
*/

func (c *clientImpl) ConvertDocToText(doc *Document) (plain string, err error) {
	url := fmt.Sprintf("%s/pf-editor-service/render", c.apiBaseURL)
	accessToken, err := c.GetAccessToken()
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %v", err)
	}
	body, err := json.Marshal(doc)
	if err != nil {
		return "", fmt.Errorf("failed to marshal document: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Accept", "text/plain")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get rendered text: %v", err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("error converting document: %s", string(b))
	}
	return string(b), nil
}

func RefreshDescriptor(url string) error {
	time.Sleep(2 * time.Second)
	resp, err := http.DefaultClient.Post(
		url,
		"text/plain",
		&bytes.Buffer{})
	if err != nil {
		return errors.Wrap(err, "Posting request: ")
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "Reading response: ")
	}

	if resp.StatusCode != 200 {
		return errors.New("Status Code not OK " + string(respBody))
	}
	return nil
}

var NotImplemented = errors.New("HTTP 501 Not Implemented")

func ReadMessage(req *http.Request) (msg *Message, err error) {
	if req.Method != http.MethodPost {
		err = NotImplemented
		return
	}

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}

	msg = &Message{}
	err = json.Unmarshal(b, msg)
	str := string(b)
	if err != nil {
		fmt.Println(str)
		return
	}
	return
}

func ReadPayload(bs []byte) (p *Payload, err error) {
	p = new(Payload)
	p.Body = new(Document)
	err = json.Unmarshal(bs, p)
	return
}
