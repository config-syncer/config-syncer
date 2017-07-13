package es

import (
	"testing"

	"github.com/stretchr/testify/assert"
	elastic "gopkg.in/olivere/elastic.v3"
)

func TestEsJanitor(t *testing.T) {
	// ElasticSearch client
	esClient, err := elastic.NewClient(
		elastic.SetURL(""),
	)
	assert.Nil(t, err)

	err = DeleteIndices(esClient, "", 0)
	assert.Nil(t, err)
}
