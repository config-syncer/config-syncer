package es

import (
	"fmt"
	"time"

	"github.com/appscode/log"
	elastic "gopkg.in/olivere/elastic.v3"
)

func DeleteIndices(esClient *elastic.Client, logIndexPrefix string, logStorageLifetime int64) error {
	now := time.Now().UTC()
	oldDate := now.Add(time.Duration(-(logStorageLifetime)) * time.Second)

	// how many index should we check to delete? I set it to 7
	for i := 1; i <= 7; i++ {
		date := oldDate.AddDate(0, 0, -i)
		prefix := fmt.Sprintf("%s%s", logIndexPrefix, date.Format("2006.01.02"))

		if _, err := esClient.Search(prefix).Do(); err == nil {
			if _, err := esClient.DeleteIndex(prefix).Do(); err != nil {
				log.Errorln(err)
				return err
			}
			log.Debugf("Index [%s] deleted", prefix)
		}
	}
	log.Debugf("ElasticSearch cleanup process complete")
	return nil
}
