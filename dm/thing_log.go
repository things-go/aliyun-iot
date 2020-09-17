package dm

import (
	"encoding/json"

	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/uri"
)

// ConfigLogParam 设备获取日志配置参数域
type ConfigLogParam struct {
	// 配置范围，目前日志只有设备维度配置，默认为device
	ConfigScope string `json:"configScope"`
	// 获取内容类型,默认为content.因日志配置内容较少,默认直接返回内容
	GetType string `json:"getType"`
}

// ThingConfigLogGet 获取日志配置
// request： /sys/${productKey}/${deviceName}/thing/config/log/get
// response：/sys/${productKey}/${deviceName}/thing/config/log/get_reply
func (sf *Client) ThingConfigLogGet(pk, dn string, _ ConfigLogParam) (*Token, error) {
	_uri := uri.URI(uri.SysPrefix, uri.ThingConfigLogGet, pk, dn)
	return sf.SendRequest(_uri, infra.MethodConfigLogGet, ConfigLogParam{
		"device",
		"content",
	})
}

// LogParam 日志内容参数域
type LogParam struct {
	UtcTime      string `json:"utcTime"`
	LogLevel     string `json:"logLevel"`
	Module       string `json:"module"`
	Code         string `json:"code"`
	TraceContext string `json:"traceContext"`
	LogContent   string `json:"logContent"`
}

// ThingConfigLogGet 设备上报日志内容
// request： /sys/${productKey}/${deviceName}/thing/config/log/post
// response：/sys/${productKey}/${deviceName}/thing/config/log/post_reply
func (sf *Client) ThingLogPost(pk, dn string, lp []LogParam) (*Token, error) {
	_uri := uri.URI(uri.SysPrefix, uri.ThingConfigLogGet, pk, dn)
	return sf.SendRequest(_uri, infra.MethodLogPost, lp)
}

// ConfigLogMode 日志配置的日志上报模式
type ConfigLogMode struct {
	Mode int `json:"mode"` // 0 表示不上报, 1 表示上报
}

// ConfigLogData 日志配置的参数域或配置域
type ConfigLogParamData struct {
	GetType string        `json:"getType"`
	Content ConfigLogMode `json:"content"`
}

// ConfigLogResponse 日志配置回复
type ConfigLogResponse struct {
	ID      uint               `json:"id,string"`
	Code    int                `json:"code"`
	Data    ConfigLogParamData `json:"data"`
	Message string             `json:"message,omitempty"`
}

// ProcThingConfigLogGetReply 处理获取日志配置应答
// request：  /sys/${productKey}/${deviceName}/thing/config/log/get
// response： /sys/${productKey}/${deviceName}/thing/config/log/get_reply
// subscribe: /sys/${productKey}/${deviceName}/thing/config/log/get_reply
func ProcThingConfigLogGetReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 7 {
		return ErrInvalidURI
	}
	rsp := &ConfigLogResponse{}
	err := json.Unmarshal(payload, rsp)
	if err != nil {
		return err
	}

	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}
	c.signalPending(Message{rsp.ID, rsp.Data, err})
	return nil
}

// ProcThingLogPostReply 处理日志上报应答
// request： /sys/${productKey}/${deviceName}/thing/log/post
// response：/sys/${productKey}/${deviceName}/thing/log/post_reply
func ProcThingLogPostReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}
	rsp := &Response{}
	err := json.Unmarshal(payload, rsp)
	if err != nil {
		return err
	}
	c.log.Debugf("thing.log.post @%d", rsp.ID)
	c.signalPending(Message{rsp.ID, nil, err})
	return nil
}

// ConfigLogPush 日志配置推送
type ConfigLogPush struct {
	ID      uint               `json:"id,string"`
	Version string             `json:"version"`
	Params  ConfigLogParamData `json:"params"`
}

// ProcThingConfigLogPush 处理日志配置推送
// subscribe: /sys/${productKey}/${deviceName}/thing/config/log/push
func ProcThingConfigLogPush(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 7 {
		return ErrInvalidURI
	}

	rsp := &ConfigLogPush{}
	err := json.Unmarshal(payload, rsp)
	if err != nil {
		return err
	}
	// TODO: 回调
	return nil
}
