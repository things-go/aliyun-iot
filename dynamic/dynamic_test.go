package dynamic

import (
	"testing"

	"github.com/thinkgos/aliyun-iot/infra"
)

func Test_calcSign(t *testing.T) {
	t.Run("calcSign", func(t *testing.T) {
		meta := infra.MetaInfo{
			ProductKey:    "a1iJcssSlPC",
			ProductSecret: "lw3QzKHNfh7XvOxO",
			DeviceName:    "dynamic",
		}

		s, err := calcSign("hmacsha256", "8Ygb7ULYh53B6OA", &meta)
		if err != nil {
			t.Errorf("calcSign() = %+v", err)
		}
		t.Logf("sign: %s", s)
	})
}
