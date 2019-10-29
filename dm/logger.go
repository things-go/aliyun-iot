package dm

// 对调试的wrapper

func (sf *Client) critical(format string, v ...interface{}) {
	sf.LogProvider().Critical(format, v...)
}

func (sf *Client) error(format string, v ...interface{}) {
	sf.LogProvider().Error(format, v...)
}

func (sf *Client) warn(format string, v ...interface{}) {
	sf.LogProvider().Warn(format, v...)
}

func (sf *Client) debug(format string, v ...interface{}) {
	sf.LogProvider().Debug(format, v...)
}
