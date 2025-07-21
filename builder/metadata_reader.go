package builder

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/go-resty/resty/v2"
)

type (
	MetadataHostBuilder func(appId AppId) string
	MetadataReaderOpt   func(bo *backoff.ExponentialBackOff)
	BackoffFactory      func() *backoff.ExponentialBackOff
)

type MetadataReader struct {
	client         *resty.Client
	hostBuilder    MetadataHostBuilder
	backoffFactory BackoffFactory
}

func NewMetadataReader(metadataHostnameNamePattern string, additionalReplacements map[string]string, opts ...MetadataReaderOpt) *MetadataReader {
	return NewMetadataReaderWithHostBuilder(func(appId AppId) string {
		return Augment(metadataHostnameNamePattern, appId, additionalReplacements)
	}, opts...)
}

func NewMetadataReaderWithHostBuilder(hostBuilder MetadataHostBuilder, opts ...MetadataReaderOpt) *MetadataReader {
	bof := func() *backoff.ExponentialBackOff {
		bo := backoff.NewExponentialBackOff()
		bo.InitialInterval = time.Second * 10
		bo.MaxInterval = time.Minute
		bo.MaxElapsedTime = time.Minute * 5

		for _, opt := range opts {
			opt(bo)
		}

		return bo
	}

	return &MetadataReader{
		client:         resty.New(),
		hostBuilder:    hostBuilder,
		backoffFactory: bof,
	}
}

func (r *MetadataReader) ReadMetadata(appId AppId) (*MetadataApplication, error) {
	metadata := &MetadataApplication{}
	path := r.hostBuilder(appId)

	bo := r.backoffFactory()
	err := backoff.Retry(func() error {
		resp, err := r.client.R().
			SetResult(metadata).
			ForceContentType("application/json").
			Get(path)
		if err != nil {
			return backoff.Permanent(err)
		}

		if resp.StatusCode() == http.StatusOK {
			return nil
		}

		if resp.StatusCode() == http.StatusBadGateway {
			return fmt.Errorf("got response code %d: metadata not yet available", resp.StatusCode())
		}

		return backoff.Permanent(fmt.Errorf("unexpected response code %d", resp.StatusCode()))
	}, bo)
	if err != nil {
		return nil, fmt.Errorf("can not read application metadata from %s: %w", path, err)
	}

	return metadata, nil
}
