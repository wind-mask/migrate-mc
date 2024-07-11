package api

import (
	"io"
	"net/http"

	"github.com/wind-mask/migrate-mc/api/modrinth"
)

const (
	ModrinthBaseUrl = "https://api.modrinth.com/v2/"
)

type Client struct {
	BaseUrl string
}

func (c *Client) get(api string) (resp *http.Response, err error) {
	return http.Get(c.BaseUrl + api)
}
func (c *Client) post(api string, body io.Reader) (resp *http.Response, err error) {
	return http.Post(c.BaseUrl+api, "application/json", body)
}
func NewClient(baseUrl string) *Client {
	return &Client{BaseUrl: baseUrl}
}
func NewModrinthClient() *Client {
	return NewClient(ModrinthBaseUrl)
}
func (c *Client) GetModFileFromVersionFile(versionFile modrinth.VersionFile) (ModFile io.ReadCloser, err error) {
	resq, err := http.Get(versionFile.Url)
	if err != nil {
		return nil, err
	}
	return resq.Body, nil
}
