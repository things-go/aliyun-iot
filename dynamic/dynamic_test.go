package dynamic

import (
	"testing"

	"github.com/thinkgos/aliyun-iot/infra"
)

func Test_requestBody(t *testing.T) {
	t.Run("calcSign", func(t *testing.T) {
		meta := infra.MetaTetrad{
			ProductKey:    "a1iJcssSlPC",
			ProductSecret: "lw3QzKHNfh7XvOxO",
			DeviceName:    "dynamic",
		}

		s := requestBody(&meta, "hmacsha256")
		t.Logf("sign: %s", s)
	})
}
