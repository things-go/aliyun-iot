package model

type dm_client_uri_map_t struct {
	devType   int
	uriPrefix string
	uriName   string
	proc      ProcDownStreamFunc
}
