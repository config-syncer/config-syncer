package es

import (
	"testing"
	api "github.com/appscode/api/kubernetes/v1beta1"
	"github.com/stretchr/testify/assert"
	elastic "gopkg.in/olivere/elastic.v3"
)

func TestEsJanitor(t *testing.T) {
	// ElasticSearch client
	esClient, err := elastic.NewClient(
		//elastic.SetSniff(false),
		elastic.SetURL(""),
	)
	assert.Nil(t, err)

	settings := &api.ClusterSettings{
		LogIndexPrefix:     "",
		LogStorageLifetime: 0,
	}
	err = DeleteIndices(esClient, settings)
	assert.Nil(t, err)
}
