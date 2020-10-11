package uri

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestURI(t *testing.T) {
	type args struct {
		prefix     string
		name       string
		productKey string
		deviceName string
		inName     []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"all",
			args{
				SysPrefix,
				ThingEventPropertyPost,
				"productKey",
				"deviceName",
				nil,
			},
			fmt.Sprintf(SysPrefix+ThingEventPropertyPost, "productKey", "deviceName"),
		},
		{
			"in name",
			args{
				prefix:     SysPrefix,
				name:       ThingEventPost,
				productKey: "productKey",
				deviceName: "deviceName",
				inName:     []string{"event"},
			},
			fmt.Sprintf(SysPrefix+ThingEventPost, "productKey", "deviceName", "event"),
		},
		{
			"空prefix",
			args{
				"",
				ThingEventPropertyPost,
				"productKey",
				"deviceName",
				nil,
			},
			ThingEventPropertyPost,
		},
		{
			"空name",
			args{
				SysPrefix,
				"",
				"productKey",
				"deviceName",
				nil,
			},
			fmt.Sprintf(SysPrefix, "productKey", "deviceName"),
		},
		{
			"空prefix和name",
			args{
				"",
				"",
				"productKey",
				"deviceName",
				nil,
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := URI(tt.args.prefix, tt.args.name, tt.args.productKey, tt.args.deviceName, tt.args.inName...); got != tt.want {
				t.Errorf("uriService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestURIReplyWithRequestURI(t *testing.T) {
	tests := []struct {
		name string
		uri  string
		want string
	}{
		{
			"/topic",
			"/topic",
			"/topic_reply",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplyWithRequestURI(tt.uri); got != tt.want {
				t.Errorf("ReplyWithRequestURI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestURISpilt(t *testing.T) {
	tests := []struct {
		name string
		uri  string
		want []string
	}{
		{
			"/a/b/c",
			"/a/b/c",
			[]string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Spilt(tt.uri); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Spilt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtRRPC(t *testing.T) {
	require.Equal(t, "/ext/rrpc/+//a/b/c", ExtRRPC("+", "/a/b/c"))
	require.Equal(t, "/ext/rrpc/+//a/b/c", ExtRRPCWildcardOne("/a/b/c"))
}
