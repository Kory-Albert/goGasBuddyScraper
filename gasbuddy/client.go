package gasbuddy

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Client struct {
	http      *http.Client
	baseURL   string
	userAgent string
}

func NewClient() *Client {
	jar, _ := cookiejar.New(nil)

	return &Client{
		baseURL:   "https://www.gasbuddy.com",
		userAgent: "Mozilla/5.0 (X11; Linux x86_64; rv:148.0) Gecko/20100101 Firefox/148.0",
		http: &http.Client{
			Timeout: 20 * time.Second,
			Jar:     jar,
		},
	}
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", c.baseURL)
	req.Header.Set("Referer", c.baseURL+"/station/")
	req.Header.Set("apollo-require-preflight", "true")

	return c.http.Do(req)
}
