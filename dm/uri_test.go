package dm

import (
	"fmt"
	"reflect"
	"testing"
)

func TestURIService(t *testing.T) {
	type args struct {
		prefix     string
		name       string
		productKey string
		deviceName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"all",
			args{
				prefix:     URISysPrefix,
				name:       URIThingEventPropertyPost,
				productKey: "productKey",
				deviceName: "deviceName",
			},
			fmt.Sprintf(URISysPrefix+URIThingEventPropertyPost, "productKey", "deviceName"),
		},
		{
			"空prefix",
			args{
				prefix:     "",
				name:       URIThingEventPropertyPost,
				productKey: "productKey",
				deviceName: "deviceName",
			},
			URIThingEventPropertyPost,
		},
		{
			"空name",
			args{
				prefix:     URISysPrefix,
				name:       "",
				productKey: "productKey",
				deviceName: "deviceName",
			},
			fmt.Sprintf(URISysPrefix, "productKey", "deviceName"),
		},
		{
			"空prefix和name",
			args{
				prefix:     "",
				name:       "",
				productKey: "productKey",
				deviceName: "deviceName",
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := URIService(tt.args.prefix, tt.args.name, tt.args.productKey, tt.args.deviceName); got != tt.want {
				t.Errorf("URIService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestURIServiceReplyWithRequestURI(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"/topic",
			args{uri: "/topic"},
			"/topic_reply",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := URIServiceReplyWithRequestURI(tt.args.uri); got != tt.want {
				t.Errorf("URIServiceReplyWithRequestURI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestURIServiceSpilt(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := URIServiceSpilt(tt.args.uri); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("URIServiceSpilt() = %v, want %v", got, tt.want)
			}
		})
	}
}
