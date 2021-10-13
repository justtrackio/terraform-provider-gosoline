package builder

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type ApiServer struct {
	Routes []struct {
		Method string `json:"method"`
		Path   string `json:"path"`
	} `json:"routes"`
}

type ApplicationMetadata struct {
	ApiServer struct {
		Routes []struct {
			Method string `json:"method"`
			Path   string `json:"path"`
		} `json:"routes"`
	} `json:"apiserver"`
	Cloud struct {
		Aws struct {
			Dynamodb struct {
				Tables []string `json:"tables"`
			} `json:"dynamodb"`
			Sqs struct {
				Queues []string `json:"queues"`
			} `json:"sqs"`
		} `json:"aws"`
	} `json:"cloud"`
}

type MetadataReader struct {
	client         *resty.Client
	metadataDomain string
}

func NewMetadataReader(metadataDomain string) *MetadataReader {
	return &MetadataReader{
		client:         resty.New(),
		metadataDomain: metadataDomain,
	}
}

func (r *MetadataReader) ReadMetadata(appId AppId) (*ApplicationMetadata, error) {
	metadata := &ApplicationMetadata{}
	path := fmt.Sprintf("http://%s.%s.%s.%s:8070", appId.Application, appId.Family, appId.Environment, r.metadataDomain)

	resp, err := r.client.R().
		SetResult(metadata).
		ForceContentType("application/json").
		Get(path)

	if err != nil {
		return nil, fmt.Errorf("can not read application metadata: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("can not read application metadata: unexpected response code %d", resp.StatusCode())
	}

	return metadata, nil
}
