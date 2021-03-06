package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sapcc/hermes/pkg/util"
	"github.com/spf13/viper"
	elastic "gopkg.in/olivere/elastic.v5"
)

//ElasticSearch struct holds esclient
type ElasticSearch struct {
	esClient *elastic.Client
}

// Mapping for attributes based on return values to API
var esFieldMapping = map[string]string{
	"time": "eventTime",

	"action":         "action",
	"outcome":        "outcome",
	"observer_id":    "observer.id",
	"observer_type":  "observer.typeURI",
	"target_id":      "target.id",
	"target_type":    "target.typeURI",
	"initiator_id":   "initiator.id",
	"initiator_type": "initiator.typeURI",
}

func (es *ElasticSearch) client() *elastic.Client {
	// Lazy initialisation - don't connect to ElasticSearch until we need to
	if es.esClient == nil {
		es.init()
	}
	return es.esClient
}

func (es *ElasticSearch) init() {
	util.LogDebug("Initializing ElasticSearch()")

	// Create a client
	var err error
	var url = viper.GetString("elasticsearch.url")
	util.LogDebug("Using ElasticSearch URL: %s", url)

	// Kubernetes LB with Elasticsearch causes challenges with IP being held on connections.
	// We can create our own custom http client, but then connections take awhile to be marked dead.
	// Syntax below...

	// Create our own http client with Transport set to deal with kubernetes lb and cached IP
	// httpClient := &http.Client{
	//	Transport: &http.Transport{
	//		DisableKeepAlives: true, // change to "false" for cached IP
	//	},
	// }
	// Connect to Elasticsearch, no sniffing due to load balancer. Custom client so no caching of IP.
	// es.esClient, err = elastic.NewClient(elastic.SetURL(url), elastic.SetHttpClient(httpClient), elastic.SetSniff(false))

	// However, that is slow to recover from a connection. We can be faster with simple client where we
	// create the connection each time we want it. I expect this will end up slow at scale, and we'll
	// have to revert to the above implementation.

	es.esClient, err = elastic.NewSimpleClient(elastic.SetURL(url))

	if err != nil {
		// TODO - Add instrumentation here for failed elasticsearch connection
		// If issues - https://github.com/olivere/elastic/wiki/Connection-Problems#how-to-figure-out-connection-problems
		panic(err)
	}
}

// GetEvents grabs events for a given tenantID with filtering.
func (es ElasticSearch) GetEvents(filter *EventFilter, tenantID string) ([]*EventDetail, int, error) {
	index := indexName(tenantID)
	util.LogDebug("Looking for events in index %s", index)

	query := elastic.NewBoolQuery()
	if filter.ObserverType != "" {
		util.LogDebug("Filtering on ObserverType %s", filter.ObserverType)
		query = query.Filter(elastic.NewMatchPhrasePrefixQuery("observer.typeURI", filter.ObserverType))
	}
	if filter.TargetType != "" {
		query = query.Filter(elastic.NewMatchPhrasePrefixQuery("target.typeURI", filter.TargetType))
	}
	// Adding .raw because it's tokenizing the ID in searches, and won't match. .raw is not analyzed, and not tokenized.
	if filter.TargetID != "" {
		query = query.Filter(elastic.NewTermQuery("target.id.raw", filter.TargetID))
	}
	if filter.InitiatorType != "" {
		query = query.Filter(elastic.NewMatchPhrasePrefixQuery("initiator.typeURI", filter.InitiatorType))
	}
	if filter.InitiatorID != "" {
		query = query.Filter(elastic.NewTermQuery("initiator.id.raw", filter.InitiatorID))
	}
	if filter.Action != "" {
		query = query.Filter(elastic.NewMatchPhrasePrefixQuery("action", filter.Action))
	}
	if filter.Outcome != "" {
		query = query.Filter(elastic.NewMatchPhrasePrefixQuery("outcome", filter.Outcome))
	}
	if filter.Time != nil && len(filter.Time) > 0 {
		for key, value := range filter.Time {
			timeField := "eventTime"
			switch key {
			case "lt":
				query = query.Filter(elastic.NewRangeQuery(timeField).Lt(value))
			case "lte":
				query = query.Filter(elastic.NewRangeQuery(timeField).Lte(value))
			case "gt":
				query = query.Filter(elastic.NewRangeQuery(timeField).Gt(value))
			case "gte":
				query = query.Filter(elastic.NewRangeQuery(timeField).Gte(value))
			}
		}
	}

	esSearch := es.client().Search().
		Index(index).
		Query(query)

	if filter.Sort != nil {
		for _, fieldOrder := range filter.Sort {
			switch fieldOrder.Order {
			case "asc":
				esSearch = esSearch.Sort(esFieldMapping[fieldOrder.Fieldname], true)
			case "desc":
				esSearch = esSearch.Sort(esFieldMapping[fieldOrder.Fieldname], false)
			}
		}
	}

	esSearch = esSearch.
		Sort("eventTime", false).
		From(int(filter.Offset)).Size(int(filter.Limit))

	searchResult, err := esSearch.Do(context.Background()) // execute
	if err != nil {
		return nil, 0, err
	}

	util.LogDebug("Got %d hits", searchResult.TotalHits())

	//Construct EventDetail array from search results
	var events []*EventDetail
	for _, hit := range searchResult.Hits.Hits {
		var de EventDetail
		err := json.Unmarshal(*hit.Source, &de)
		if err != nil {
			return nil, 0, err
		}
		events = append(events, &de)
	}
	total := searchResult.TotalHits()

	return events, int(total), nil
}

// GetEvent Returns EventDetail for a single event.
func (es ElasticSearch) GetEvent(eventID string, tenantID string) (*EventDetail, error) {
	index := indexName(tenantID)
	util.LogDebug("Looking for event %s in index %s", eventID, index)

	// we use .raw on ID fields because Elasticsearch tokenizes fields with - in them. .raw is not tokenized.
	query := elastic.NewTermQuery("id.raw", eventID)
	esSearch := es.client().Search().
		Index(index).
		Query(query)

	searchResult, err := esSearch.Do(context.Background())
	if err != nil {
		util.LogDebug("Query failed: %s", err.Error())
		return nil, err
	}
	total := searchResult.TotalHits()
	util.LogDebug("Results: %d", total)

	if total > 0 {
		hit := searchResult.Hits.Hits[0]
		var de EventDetail
		err := json.Unmarshal(*hit.Source, &de)
		return &de, err
	}
	return nil, nil
}

// GetAttributes Return all unique attributes available for filtering
// Possible queries, event_type, dns, identity, etc..
func (es ElasticSearch) GetAttributes(filter *AttributeFilter, tenantID string) ([]string, error) {
	index := indexName(tenantID)

	util.LogDebug("Looking for unique attributes for %s in index %s", filter.QueryName, index)

	// ObserverType in this case is not the cadf source, but instead the first part of event_type
	var esName string
	// Append .raw onto queryName, in Elasticsearch. Aggregations turned on for .raw
	if val, ok := esFieldMapping[filter.QueryName]; ok {
		esName = val + ".raw"
	} else {
		esName = filter.QueryName + ".raw"
	}
	util.LogDebug("Mapped Queryname: %s --> %s", filter.QueryName, esName)

	queryAgg := elastic.NewTermsAggregation().Size(int(filter.Limit)).Field(esName)

	esSearch := es.client().Search().Index(index).Size(int(filter.Limit)).Aggregation("attributes", queryAgg)
	searchResult, err := esSearch.Do(context.Background())
	if err != nil {
		return nil, err
	}

	if searchResult.Hits == nil {
		util.LogDebug("expected Hits != nil; got: nil")
	}

	agg := searchResult.Aggregations
	if agg == nil {
		util.LogDebug("expected Aggregations, got nil")
	}

	termsAggRes, found := agg.Terms("attributes")
	if !found {
		util.LogDebug("Term %s not found in Aggregation", esName)
	}
	if termsAggRes == nil {
		util.LogDebug("termsAggRes is nil")
	}
	util.LogDebug("Number of Buckets: %d", len(termsAggRes.Buckets))

	var unique []string
	for _, bucket := range termsAggRes.Buckets {
		util.LogDebug("key: %s count: %d", bucket.Key, bucket.DocCount)
		attribute := bucket.Key.(string)

		// Hierarchical Depth Handling
		var att string
		if filter.MaxDepth != 0 && strings.Contains(attribute, "/") {
			s := strings.Split(attribute, "/")
			l := len(s)
			for i := 0; i < int(filter.MaxDepth) && i < l; i++ {
				if i != 0 {
					att += "/"
				}
				att += s[i]
			}
			attribute = att
		}

		unique = append(unique, attribute)

	}

	unique = SliceUniqMap(unique)
	return unique, nil
}

// SliceUniqMap Removes duplicates from slice
func SliceUniqMap(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}

// MaxLimit grabs the configured maxlimit for results
func (es ElasticSearch) MaxLimit() uint {
	return uint(viper.GetInt("elasticsearch.max_result_window"))
}

// indexName Generates the index name for a given TenantId. If no tenantID defaults to audit-*
// Records for audit-* will not be accessible from a given Tenant.
func indexName(tenantID string) string {
	index := "audit-*"
	if tenantID != "" {
		index = fmt.Sprintf("audit-%s-*", tenantID)
	}
	return index
}
