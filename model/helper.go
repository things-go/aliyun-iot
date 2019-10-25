package model

// URIService 获得本设备URI
func (sf *Manager) URIService(prefix, name string) string {
	return URIService(prefix, name, sf.ProductKey, sf.DeviceName)
}
