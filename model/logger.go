package model

// 对调试的wrapper

func (sf *Manager) critical(format string, v ...interface{}) {
	sf.LogProvider().Critical(format, v...)
}

func (sf *Manager) error(format string, v ...interface{}) {
	sf.LogProvider().Error(format, v...)
}

func (sf *Manager) warn(format string, v ...interface{}) {
	sf.LogProvider().Warn(format, v...)
}

func (sf *Manager) debug(format string, v ...interface{}) {
	sf.LogProvider().Debug(format, v...)
}
