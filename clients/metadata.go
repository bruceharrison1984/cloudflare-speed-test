package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bruceharrison1984/cloudflare-speed-test/providers"
	"github.com/bruceharrison1984/cloudflare-speed-test/types"
)

type IMetadataClient interface {
	FetchMetadata() (*types.CloudflareMetadata, error)
}

type metadataClient struct {
	Http        *http.Client
	urlProvider providers.IUrlProvider
}

func NewMetadataClient(http *http.Client, urlProvider providers.IUrlProvider) IMetadataClient {
	return &metadataClient{http, urlProvider}
}

func (mdc metadataClient) FetchMetadata() (*types.CloudflareMetadata, error) {
	client := mdc.Http

	resp, err := client.Get(mdc.urlProvider.GetMetadataUrl())
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
