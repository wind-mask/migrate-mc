package api

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/wind-mask/migrate-mc/api/modrinth"
)
var Logger = log.Default()
const (
	ModrinthBaseUrl = "https://api.modrinth.com/v2/"
)

type Client struct {
	BaseUrl             string
	Ratelimit_limit     int
	Ratelimit_remaining int
	Ratelimit_reset     int
}

func (c *Client) updateRatelimit(header *http.Header) error {
	var err error
	c.Ratelimit_limit, err = strconv.Atoi(header.Get("x-ratelimit-limit"))
	if err != nil {
		return err
	}

	c.Ratelimit_remaining, err = strconv.Atoi(header.Get("x-ratelimit-remaining"))
	if err != nil {
		return err
	}

	rr, err := strconv.Atoi(header.Get("x-ratelimit-reset"))
	if err != nil {
		return err
	}
	c.Ratelimit_reset = rr+int(time.Now().Unix())
	return nil
}
func (c *Client)waitIfLimit() {
	if c.Ratelimit_remaining == 0 {
		time.Sleep(time.Duration(c.Ratelimit_reset-int(time.Now().Unix())) * time.Second)
	}
}
func (c *Client) get(api string) (resp *http.Response, err error) {
	c.waitIfLimit()
	r, err := http.Get(c.BaseUrl + api)
	if err != nil {
		return nil, err
	}
	err = c.updateRatelimit(&r.Header)
	if err != nil {
		return nil, err
	}
	return r, nil

}
func (c *Client) post(api string, body io.Reader) (resp *http.Response, err error) {
	c.waitIfLimit()
	r,err:=http.Post(c.BaseUrl+api, "application/json", body)
	if err != nil {
		return nil, err
	}
	err = c.updateRatelimit(&r.Header)
	if err != nil {
		return nil, err
	}
	return r, nil
}
func NewClient(baseUrl string) *Client {
	return &Client{BaseUrl: baseUrl}
}
func NewModrinthClient() *Client {
	return NewClient(ModrinthBaseUrl)
}
func (c *Client) GetModFileFromVersionFile(versionFile modrinth.VersionFile) (ModFile io.ReadCloser, err error) {
	Logger.Println("GetModFileFromVersionFile:", versionFile.Url)
	resq, err := http.Get(versionFile.Url)
	if err != nil {
		return nil, err
	}
	return resq.Body, nil
}
