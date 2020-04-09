package dynamic

import (
	"testing"

	"github.com/thinkgos/aiot/infra"
)

func Test_calcSign(t *testing.T) {
	t.Run("calcSign", func(t *testing.T) {
		met := MetaSign{
			ProductKey:    "a1iJcssSlPC",
			ProductSecret: "lw3QzKHNfh7XvOxO",
			DeviceName:    "dynamic",
			Random:        "8Ygb7ULYh53B6OA",
			SignMethod:    infra.SignMethodHMACSHA256,
		}

		s, err := calcSign(&met)
		if err != nil {
			t.Errorf("calcSign() = %+v", err)
		}
		t.Logf("sign: %s", s)
	})
}
