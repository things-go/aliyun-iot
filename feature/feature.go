package feature

type Options struct {
	NTP      bool
	RawModel bool
	Desired  bool
	ExtRRPC  bool
	Gateway  bool
}

func New() *Options {
	return &Options{}
}

func (sf *Options) EnableNTP() *Options {
	sf.NTP = true
	return sf
}

// 使能透传 EnableModelRaw
func (sf *Options) EnableModelRaw() *Options {
	sf.RawModel = true
	return sf
}

// EnableDesired 使能期望属性
func (sf *Options) EnableDesired() *Options {
	sf.Desired = true
	return sf
}

// EnableExtRRPC 使能扩展RRPC功能
func (sf *Options) EnableExtRRPC() *Options {
	sf.ExtRRPC = true
	return sf
}

func (sf *Options) EnableGateway() *Options {
	sf.Gateway = true
	return sf
}
