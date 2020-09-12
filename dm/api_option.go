package dm

import (
	"time"

	"github.com/thinkgos/go-core-package/lib/logger"
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

// WithWork 工作在...之上,WorkOnCOAP,WorkOnHTTP,WorkOnMQTT(默认)
func WithWork(on int) Option {
	return func(c *Client) {
		switch on {
		case WorkOnCOAP:
			c.workOnWho = WorkOnCOAP
		case WorkOnHTTP:
			c.workOnWho = WorkOnHTTP
		default:
			c.workOnWho = WorkOnMQTT
		}
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
		c.log = l
	}
}
