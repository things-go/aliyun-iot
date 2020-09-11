package ahttp

import (
	"net/http"
	"strings"

	"github.com/thinkgos/aliyun-iot/clog"
)

// Option client option
type Option func(c *Client)

// WithHTTPClient with custom http.Client
func WithHTTPClient(c *http.Client) Option {
	return func(client *Client) {
		client.httpc = c
	}
}

// WithHost 设置远程主机
func WithHost(h string) Option {
	return func(c *Client) {
		if !strings.Contains(h, "://") {
			h = "http://" + h
		}
		if h != "" {
			c.host = h
		}
	}
}

// WithSignMethod 设置签名方法,目前支持hmacsha1,hmacmd5(默认)
func WithSignMethod(method string) Option {
	return func(c *Client) {
		if method == hmacsha1 {
			c.signMethod = hmacsha1
		} else {
			c.signMethod = hmacmd5
		}
	}
}

// WithLogger 设置日志
func WithLogger(l clog.LogProvider) Option {
	return func(c *Client) {
		c.log = l
	}
}
