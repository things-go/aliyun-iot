package dynamic

import (
	"testing"

	"github.com/thinkgos/aliyun-iot/infra"
)

func Test_requestBody(t *testing.T) {
	t.Run("hmacsha256", func(t *testing.T) {
		meta := infra.MetaTetrad{
			ProductKey:    "ProductKey",
			ProductSecret: "ProductSecret",
			DeviceName:    "DeviceName",
		}

		s := requestBody(&meta, "hmacsha256")
		t.Logf("sign: %s", s)
	})
	t.Run("default", func(t *testing.T) {
		meta := infra.MetaTetrad{
			ProductKey:    "ProductKey",
			ProductSecret: "ProductSecret",
			DeviceName:    "DeviceName",
		}

		s := requestBody(&meta, "default")
		t.Logf("sign: %s", s)
	})
}
