package es

import (
	"fmt"
	"time"

	api "github.com/appscode/api/kubernetes/v1beta1"
	"github.com/appscode/log"
	elastic "gopkg.in/olivere/elastic.v3"
)

func DeleteIndices(esClient *elastic.Client, j *api.ClusterSettings) error {
	now := time.Now().UTC()
	oldDate := now.Add(time.Duration(-(j.LogStorageLifetime)) * time.Second)

	// how many index should we check to delete? I set it to 7
	for i := 1; i <= 7; i++ {
		date := oldDate.AddDate(0, 0, -i)
		prefix := fmt.Sprintf("%s%s", j.LogIndexPrefix, date.Format("2006.01.02"))

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
