package builder_test

import (
	"sort"
	"testing"

	"github.com/justtrackio/terraform-provider-gosoline/builder"
	"github.com/stretchr/testify/assert"
)

// TestCollectionSorting verifies that collections can be sorted consistently
// This tests the sorting approach used in the dashboard datasource
func TestCollectionSorting(t *testing.T) {
	// Test stream consumers sorting
	t.Run("StreamConsumers", func(t *testing.T) {
		consumers := []builder.MetadataStreamConsumer{
			{Name: "consumer-z"},
			{Name: "consumer-a"},
			{Name: "consumer-m"},
		}

		// Sort by name
		sort.Slice(consumers, func(i, j int) bool {
			return consumers[i].Name < consumers[j].Name
		})

		assert.Equal(t, "consumer-a", consumers[0].Name)
		assert.Equal(t, "consumer-m", consumers[1].Name)
		assert.Equal(t, "consumer-z", consumers[2].Name)
	})

	// Test stream producers sorting
	t.Run("StreamProducers", func(t *testing.T) {
		producers := []builder.MetadataStreamProducer{
			{Name: "producer-y"},
			{Name: "producer-b"},
			{Name: "producer-k"},
		}

		// Sort by name
		sort.Slice(producers, func(i, j int) bool {
			return producers[i].Name < producers[j].Name
		})

		assert.Equal(t, "producer-b", producers[0].Name)
		assert.Equal(t, "producer-k", producers[1].Name)
		assert.Equal(t, "producer-y", producers[2].Name)
	})

	// Test Kinesis kinsumers sorting
	t.Run("KinesisKinsumers", func(t *testing.T) {
		kinsumers := []builder.MetadataCloudAwsKinesisKinsumer{
			{Name: "kinsumer-x"},
			{Name: "kinsumer-c"},
			{Name: "kinsumer-p"},
		}

		// Sort by name
		sort.Slice(kinsumers, func(i, j int) bool {
			return kinsumers[i].Name < kinsumers[j].Name
		})

		assert.Equal(t, "kinsumer-c", kinsumers[0].Name)
		assert.Equal(t, "kinsumer-p", kinsumers[1].Name)
		assert.Equal(t, "kinsumer-x", kinsumers[2].Name)
	})

	// Test Kinesis record writers sorting
	t.Run("KinesisRecordWriters", func(t *testing.T) {
		writers := []builder.MetadataCloudAwsKinesisRecordWriter{
			{StreamName: "writer-stream-z"},
			{StreamName: "writer-stream-a"},
			{StreamName: "writer-stream-m"},
		}

		// Sort by stream name
		sort.Slice(writers, func(i, j int) bool {
			return writers[i].StreamName < writers[j].StreamName
		})

		assert.Equal(t, "writer-stream-a", writers[0].StreamName)
		assert.Equal(t, "writer-stream-m", writers[1].StreamName)
		assert.Equal(t, "writer-stream-z", writers[2].StreamName)
	})

	// Test SQS queues sorting
	t.Run("SqsQueues", func(t *testing.T) {
		queues := []builder.MetadataCloudAwsSqsQueue{
			{QueueNameFull: "full-queue-z"},
			{QueueNameFull: "full-queue-a"},
			{QueueNameFull: "full-queue-m"},
		}

		// Sort by queue name full
		sort.Slice(queues, func(i, j int) bool {
			return queues[i].QueueNameFull < queues[j].QueueNameFull
		})

		assert.Equal(t, "full-queue-a", queues[0].QueueNameFull)
		assert.Equal(t, "full-queue-m", queues[1].QueueNameFull)
		assert.Equal(t, "full-queue-z", queues[2].QueueNameFull)
	})

	// Test DynamoDB tables sorting
	t.Run("DynamoDbTables", func(t *testing.T) {
		tables := []builder.MetadataCloudAwsDynamodbTable{
			{TableName: "table-z"},
			{TableName: "table-a"},
			{TableName: "table-m"},
		}

		// Sort by table name
		sort.Slice(tables, func(i, j int) bool {
			return tables[i].TableName < tables[j].TableName
		})

		assert.Equal(t, "table-a", tables[0].TableName)
		assert.Equal(t, "table-m", tables[1].TableName)
		assert.Equal(t, "table-z", tables[2].TableName)
	})

	// Test HTTP servers sorting
	t.Run("HttpServers", func(t *testing.T) {
		servers := []builder.MetadataHttpServer{
			{Name: "server-z"},
			{Name: "server-a"},
			{Name: "server-m"},
		}

		// Sort by name
		sort.Slice(servers, func(i, j int) bool {
			return servers[i].Name < servers[j].Name
		})

		assert.Equal(t, "server-a", servers[0].Name)
		assert.Equal(t, "server-m", servers[1].Name)
		assert.Equal(t, "server-z", servers[2].Name)
	})

	// Test HTTP server handlers sorting
	t.Run("HttpServerHandlers", func(t *testing.T) {
		handlers := []builder.MetadataHttpServerHandler{
			{Method: "POST", Path: "/api/z"},
			{Method: "GET", Path: "/api/z"},
			{Method: "GET", Path: "/api/a"},
			{Method: "PUT", Path: "/api/a"},
		}

		// Sort by method first, then path
		sort.Slice(handlers, func(i, j int) bool {
			if handlers[i].Method != handlers[j].Method {
				return handlers[i].Method < handlers[j].Method
			}
			return handlers[i].Path < handlers[j].Path
		})

		// Expected order: GET /api/a, GET /api/z, POST /api/z, PUT /api/a
		assert.Equal(t, "GET", handlers[0].Method)
		assert.Equal(t, "/api/a", handlers[0].Path)

		assert.Equal(t, "GET", handlers[1].Method)
		assert.Equal(t, "/api/z", handlers[1].Path)

		assert.Equal(t, "POST", handlers[2].Method)
		assert.Equal(t, "/api/z", handlers[2].Path)

		assert.Equal(t, "PUT", handlers[3].Method)
		assert.Equal(t, "/api/a", handlers[3].Path)
	})
}

// TestSortingIdempotency verifies that sorting the same collection multiple times
// produces the same result
func TestSortingIdempotency(t *testing.T) {
	// Create a collection that will be sorted multiple times
	consumers := []builder.MetadataStreamConsumer{
		{Name: "consumer-z"},
		{Name: "consumer-a"},
		{Name: "consumer-m"},
		{Name: "consumer-b"},
		{Name: "consumer-x"},
	}

	// Create copies to sort independently
	copy1 := make([]builder.MetadataStreamConsumer, len(consumers))
	copy2 := make([]builder.MetadataStreamConsumer, len(consumers))
	copy3 := make([]builder.MetadataStreamConsumer, len(consumers))

	copy(copy1, consumers)
	copy(copy2, consumers)
	copy(copy3, consumers)

	// Sort each copy
	sort.Slice(copy1, func(i, j int) bool {
		return copy1[i].Name < copy1[j].Name
	})

	sort.Slice(copy2, func(i, j int) bool {
		return copy2[i].Name < copy2[j].Name
	})

	sort.Slice(copy3, func(i, j int) bool {
		return copy3[i].Name < copy3[j].Name
	})

	// All should be identical
	assert.Equal(t, copy1, copy2)
	assert.Equal(t, copy1, copy3)
	assert.Equal(t, copy2, copy3)

	// Verify expected order
	expectedNames := []string{"consumer-a", "consumer-b", "consumer-m", "consumer-x", "consumer-z"}
	for i, expected := range expectedNames {
		assert.Equal(t, expected, copy1[i].Name)
	}
}
