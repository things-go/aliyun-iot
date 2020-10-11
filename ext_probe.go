package aiot

// ProcExtNetworkProbeRequest 处理平台测试延迟请求
// request:  /ext/network/probe/${messageId}
// subscribe: /ext/network/probe/+
func ProcExtNetworkProbeRequest(c *Client, rawURI string, _ []byte) error {
	c.Log.Debugf("ext.network.probe --> uri: %s", rawURI)
	return nil
}
