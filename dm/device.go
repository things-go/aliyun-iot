// Package dm 实现阿里云物模型
package dm

// @see https://help.aliyun.com/document_detail/89301.html?spm=a2c4g.11186623.6.706.570f3f69J3fW5z

// ProcDownStream 处理下行数据
type ProcDownStream func(c *Client, rawURI string, payload []byte) error
