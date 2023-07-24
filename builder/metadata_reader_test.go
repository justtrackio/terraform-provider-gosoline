package builder_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/justtrackio/terraform-provider-gosoline/builder"
	"github.com/stretchr/testify/assert"
)

func TestReadApplicationMetadataSuccess(t *testing.T) {
	actual := &builder.MetadataApplication{
		Cloud: builder.MetadataCloud{
			Aws: builder.MetadataCloudAws{
				Dynamodb: builder.MetadataCloudAwsDynamodb{
					Tables: []builder.MetadataCloudAwsDynamodbTable{
						{
							AwsClientName: "default",
							TableName:     "foo",
						},
						{
							AwsClientName: "default",
							TableName:     "bar",
						},
					},
				},
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		bytes, err := json.Marshal(actual)
		assert.NoError(t, err)

		writer.Write(bytes)
	}))

	reader := builder.NewMetadataReaderWithHostBuilder(func(_ builder.AppId) string {
		return ts.URL
	})

	data, err := reader.ReadMetadata(builder.AppId{})
	assert.NoError(t, err)
	assert.Equal(t, actual, data)
}

func TestReadApplicationMetadataRetry(t *testing.T) {
	attempts := 0

	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		attempts++

		if attempts < 3 {
			writer.WriteHeader(http.StatusBadGateway)
			return
		}

		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(`{}`))
	}))

	reader := builder.NewMetadataReaderWithHostBuilder(func(_ builder.AppId) string {
		return ts.URL
	}, func(bo *backoff.ExponentialBackOff) {
		bo.InitialInterval = time.Millisecond * 100
		bo.MaxInterval = time.Millisecond * 100
	})

	data, err := reader.ReadMetadata(builder.AppId{})
	assert.NoError(t, err)
	assert.Equal(t, &builder.MetadataApplication{}, data)
	assert.Equal(t, 3, attempts)
}

func TestReadApplicationMetadataRetryTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		time.Sleep(time.Millisecond * 100)
		writer.WriteHeader(http.StatusBadGateway)
	}))

	reader := builder.NewMetadataReaderWithHostBuilder(func(_ builder.AppId) string {
		return ts.URL
	}, func(bo *backoff.ExponentialBackOff) {
		bo.InitialInterval = time.Millisecond
		bo.MaxInterval = time.Millisecond * 10
		bo.MaxElapsedTime = time.Millisecond * 100
	})

	data, err := reader.ReadMetadata(builder.AppId{})
	assert.EqualError(t, err, "can not read application metadata from "+ts.URL+": got response code 502: metadata not yet available")
	assert.Nil(t, data)
}
