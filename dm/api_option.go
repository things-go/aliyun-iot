package dm

import (
	"time"

	"github.com/thinkgos/go-core-package/lib/logger"
)

type Option func(*Client)

// WithConfigCache 设备消息缓存超时时间
func WithCache(expiration, cleanupInterval time.Duration) Option {
	return func(c *Client) {
		c.cacheExpiration = expiration
		c.cacheCleanupInterval = cleanupInterval
	}
}

// WithConfigWork 工作在...之上,WorkOnCOAP,WorkOnHTTP,WorkOnMQTT(默认)
func WithWork(on int) Option {
	return func(c *Client) {
		switch on {
		case WorkOnCOAP:
			c.workOnWho = WorkOnCOAP
			c.uriOffset = 1
		case WorkOnHTTP:
			c.workOnWho = WorkOnHTTP
			c.uriOffset = 1
		default:
			c.workOnWho = WorkOnMQTT
			c.uriOffset = 0
		}
	}
}

// WithConfigEnableNTP 使能NTP
func WithEnableNTP() Option {
	return func(c *Client) {
		c.hasNTP = true
	}
}

// WithConfigEnableModelRaw 使能透传
func WithEnableModelRaw() Option {
	return func(c *Client) {
		c.hasRawModel = true
	}
}

// WithConfigEnableDesired 使能期望属性
func WithEnableDesired() Option {
	return func(c *Client) {
		c.hasDesired = true
	}
}

// WithConfigEnableExtRRPC 使能扩展RRPC功能
func WithEnableExtRRPC() Option {
	return func(c *Client) {
		c.hasExtRRPC = true
	}
}

// WithConfigEnableGateway 使能网关功能
func WithEnableGateway() Option {
	return func(c *Client) {
		c.isGateway = true
	}
}

// WithConfigEnableOTA 使能ota功能
func WithEnableOTA() Option {
	return func(c *Client) {
		c.hasOTA = true
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

func WithLogger(l logger.Logger) Option {
	return func(c *Client) {
		c.log = l
	}
}
