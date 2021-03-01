// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ahttp

import (
	"net/http"
	"strings"

	"github.com/thinkgos/aliyun-iot/logger"
)

// Option client option
type Option func(c *Client)

// WithHTTPClient with custom http.Client
func WithHTTPClient(c *http.Client) Option {
	return func(client *Client) {
		client.httpc = c
	}
}

// WithEndpoint 设置Endpoint地址,也就是Host
func WithEndpoint(h string) Option {
	return func(c *Client) {
		if !strings.Contains(h, "://") {
			h = "http://" + h
		}
		if h != "" {
			c.endpoint = h
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
func WithLogger(l logger.Logger) Option {
	return func(c *Client) {
		c.log = l
	}
}
