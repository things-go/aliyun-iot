package model

func (sf *Manager) URIService(prefix, name string) string {
	return URIService(prefix, name, sf.ProductKey, sf.DeviceName)
}
