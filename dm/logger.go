package dm

// 对调试的wrapper

func (sf *Client) criticalf(format string, v ...interface{}) {
	sf.LogProvider().Criticalf(format, v...)
}

func (sf *Client) errorf(format string, v ...interface{}) {
	sf.LogProvider().Errorf(format, v...)
}

func (sf *Client) warnf(format string, v ...interface{}) {
	sf.LogProvider().Warnf(format, v...)
}

func (sf *Client) debugf(format string, v ...interface{}) {
	sf.LogProvider().Debugf(format, v...)
}
