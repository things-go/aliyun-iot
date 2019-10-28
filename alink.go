package aliIOT

import (
	"github.com/thinkgos/aliIOT/feature"
	"github.com/thinkgos/aliIOT/model"
)

type Client struct {
	*model.Manager
	features *feature.Options
}
