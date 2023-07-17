package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bruceharrison1984/cloudflare-speed-test/providers"
	"github.com/bruceharrison1984/cloudflare-speed-test/types"
)

type MetadataClient struct {
	Http        *http.Client
	urlProvider providers.UrlProvider
}

func NewMetadataClient(http *http.Client, urlProvider providers.UrlProvider) *MetadataClient {
	return &MetadataClient{http, urlProvider}
}

func (metadataClient MetadataClient) FetchMetadata() (*types.CloudflareMetadata, error) {
	client := metadataClient.Http

	resp, err := client.Get(metadataClient.urlProvider.GetMetadataUrl())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	metadataBody, _ := io.ReadAll(resp.Body)

	var metadata types.CloudflareMetadata
	if err := json.Unmarshal(metadataBody, &metadata); err != nil {
		fmt.Println("failed to unmarshal:", err)
	}
	return &metadata, nil
}
