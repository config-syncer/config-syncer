package es

import (
	"fmt"
	"time"

	"github.com/appscode/go/log"
	"github.com/appscode/kubed/pkg/config"
	elastic "gopkg.in/olivere/elastic.v3"
)

type Janitor struct {
	Spec config.ElasticSearchSpec
	TTL  time.Duration
}

func (j *Janitor) Cleanup() error {
	client, err := elastic.NewClient(
		// elastic.SetSniff(false),
		elastic.SetURL(j.Spec.Endpoint),
	)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	oldDate := now.Add(-j.TTL)

	// how many index should we check to delete? I set it to 7
	for i := 1; i <= 7; i++ {
		date := oldDate.AddDate(0, 0, -i)
		prefix := fmt.Sprintf("%s%s", j.Spec.LogIndexPrefix, date.Format("2006.01.02"))

		if _, err := client.Search(prefix).Do(); err == nil {
			if _, err := client.DeleteIndex(prefix).Do(); err != nil {
				log.Errorln(err)
				return err
			}
			log.Debugf("Index [%s] deleted", prefix)
		}
	}
	log.Debugf("ElasticSearch cleanup process complete")
	return nil
}
