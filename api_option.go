// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package aiot

import (
	"time"

	"github.com/thinkgos/x/lib/logger"
)

// Option 配置选项
type Option func(*Client)

// WithCache 设备消息缓存超时时间
func WithCache(expiration, cleanupInterval time.Duration) Option {
	return func(c *Client) {
		c.cacheExpiration = expiration
		c.cacheCleanupInterval = cleanupInterval
	}
}

// WithMode 设置工作模式 支持 ModeCOAP ,ModeHTTP, ModeMQTT(默认)
func WithMode(m Mode) Option {
	return func(c *Client) {
		c.mode = m
	}
}

// WithVersion 设置平台版本,默认为 DefaultVersion
func WithVersion(ver string) Option {
	return func(c *Client) {
		c.version = ver
	}
}

// WithEnableNTP 使能NTP
func WithEnableNTP() Option {
	return func(c *Client) {
		c.hasNTP = true
	}
}

// WithEnableModelRaw 使能透传
func WithEnableModelRaw() Option {
	return func(c *Client) {
		c.hasRawModel = true
	}
}

// WithEnableDesired 使能期望属性
func WithEnableDesired() Option {
	return func(c *Client) {
		c.hasDesired = true
	}
}

// WithEnableExtRRPC 使能扩展RRPC功能
func WithEnableExtRRPC() Option {
	return func(c *Client) {
		c.hasExtRRPC = true
	}
}

// WithEnableGateway 使能网关功能
func WithEnableGateway() Option {
	return func(c *Client) {
		c.isGateway = true
	}
}

// WithEnableOTA 使能ota功能
func WithEnableOTA() Option {
	return func(c *Client) {
		c.hasOTA = true
	}
}

// WithEnableDiag 使能diag功能
func WithEnableDiag() Option {
	return func(c *Client) {
		c.hasDiag = true
	}
}

// WithCallback 设置事件处理接口
func WithCallback(cb Callback) Option {
	return func(c *Client) {
		c.cb = cb
	}
}

// WithGwCallback 设备网关事件接口
func WithGwCallback(cb GwCallback) Option {
	return func(c *Client) {
		c.gwCb = cb
	}
}

// WithLogger 设置日志
func WithLogger(l logger.Logger) Option {
	return func(c *Client) {
		c.Log = l
	}
}
