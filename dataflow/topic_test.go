package dataflow

import (
	"reflect"
	"testing"
)

func TestParseTopicStatus(t *testing.T) {
	type args struct {
		topic string
	}
	tests := []struct {
		name    string
		args    args
		wantTi  TopicInfo
		wantErr bool
	}{
		{
			"status",
			args{"/as/mqtt/status/productKey/deviceName"},
			TopicInfo{"productKey", "deviceName", ""},
			false,
		},
		{
			"status",
			args{""},
			TopicInfo{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTi, err := ParseTopicStatus(tt.args.topic)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTopicStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTi, tt.wantTi) {
				t.Errorf("ParseTopicStatus() gotTi = %v, want %v", gotTi, tt.wantTi)
			}
		})
	}
}

func TestParseTopicEvent(t *testing.T) {
	type args struct {
		topic string
	}
	tests := []struct {
		name    string
		args    args
		wantTi  TopicInfo
		wantErr bool
	}{
		{
			"event property",
			args{"/productKey/deviceName/thing/event/property/post"},
			TopicInfo{"productKey", "deviceName", "property"},
			false,
		},
		{
			"event tsl.event.identifier",
			args{"/productKey/deviceName/thing/event/tsl.event.identifier/post"},
			TopicInfo{"productKey", "deviceName", "tsl.event.identifier"},
			false,
		},
		{
			"status",
			args{""},
			TopicInfo{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTi, err := ParseTopicEvent(tt.args.topic)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTopicEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTi, tt.wantTi) {
				t.Errorf("ParseTopicEvent() gotTi = %v, want %v", gotTi, tt.wantTi)
			}
		})
	}
}

func TestParseTopicLifecycle(t *testing.T) {
	type args struct {
		topic string
	}
	tests := []struct {
		name    string
		args    args
		wantTi  TopicInfo
		wantErr bool
	}{
		{
			"status",
			args{"/productKey/deviceName/thing/lifecycle"},
			TopicInfo{"productKey", "deviceName", ""},
			false,
		},
		{
			"status",
			args{""},
			TopicInfo{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTi, err := ParseTopicLifecycle(tt.args.topic)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTopicLifecycle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTi, tt.wantTi) {
				t.Errorf("ParseTopicLifecycle() gotTi = %v, want %v", gotTi, tt.wantTi)
			}
		})
	}
}

func TestParseTopicTopoLifecycle(t *testing.T) {
	type args struct {
		topic string
	}
	tests := []struct {
		name    string
		args    args
		wantTi  TopicInfo
		wantErr bool
	}{
		{
			"status",
			args{"/productKey/deviceName/thing/topo/lifecycle"},
			TopicInfo{"productKey", "deviceName", ""},
			false,
		},
		{
			"status",
			args{""},
			TopicInfo{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTi, err := ParseTopicTopoLifecycle(tt.args.topic)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTopicTopoLifecycle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTi, tt.wantTi) {
				t.Errorf("ParseTopicTopoLifecycle() gotTi = %v, want %v", gotTi, tt.wantTi)
			}
		})
	}
}

func TestParseTopicListFound(t *testing.T) {
	type args struct {
		topic string
	}
	tests := []struct {
		name    string
		args    args
		wantTi  TopicInfo
		wantErr bool
	}{
		{
			"status",
			args{"/productKey/deviceName/thing/list/found"},
			TopicInfo{"productKey", "deviceName", ""},
			false,
		},
		{
			"status",
			args{""},
			TopicInfo{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTi, err := ParseTopicListFound(tt.args.topic)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTopicListFound() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTi, tt.wantTi) {
				t.Errorf("ParseTopicListFound() gotTi = %v, want %v", gotTi, tt.wantTi)
			}
		})
	}
}

func TestParseTopicDownLinkReply(t *testing.T) {
	type args struct {
		topic string
	}
	tests := []struct {
		name    string
		args    args
		wantTi  TopicInfo
		wantErr bool
	}{
		{
			"status",
			args{"/productKey/deviceName/thing/downlink/reply/message"},
			TopicInfo{"productKey", "deviceName", ""},
			false,
		},
		{
			"status",
			args{""},
			TopicInfo{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTi, err := ParseTopicDownLinkReply(tt.args.topic)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTopicDownLinkReply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTi, tt.wantTi) {
				t.Errorf("ParseTopicDownLinkReply() gotTi = %v, want %v", gotTi, tt.wantTi)
			}
		})
	}
}

func TestParseHistoryEvent(t *testing.T) {
	type args struct {
		topic string
	}
	tests := []struct {
		name    string
		args    args
		wantTi  TopicInfo
		wantErr bool
	}{
		{
			"event property",
			args{"/sys/productKey/deviceName/thing/event/property/history/post"},
			TopicInfo{"productKey", "deviceName", "property"},
			false,
		},
		{
			"event tsl.event.identifier",
			args{"/sys/productKey/deviceName/thing/event/tsl.event.identifier/history/post"},
			TopicInfo{"productKey", "deviceName", "tsl.event.identifier"},
			false,
		},
		{
			"status",
			args{""},
			TopicInfo{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTi, err := ParseTopicHistoryEvent(tt.args.topic)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTopicHistoryEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTi, tt.wantTi) {
				t.Errorf("ParseTopicHistoryEvent() gotTi = %v, want %v", gotTi, tt.wantTi)
			}
		})
	}
}

func TestParseOtaUpgrade(t *testing.T) {
	type args struct {
		topic string
	}
	tests := []struct {
		name    string
		args    args
		wantTi  TopicInfo
		wantErr bool
	}{
		{
			"status",
			args{"/sys/productKey/deviceName/ota/upgrade"},
			TopicInfo{"productKey", "deviceName", ""},
			false,
		},
		{
			"status",
			args{""},
			TopicInfo{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTi, err := ParseTopicOtaUpgrade(tt.args.topic)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTopicOtaUpgrade() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTi, tt.wantTi) {
				t.Errorf("ParseTopicOtaUpgrade() gotTi = %v, want %v", gotTi, tt.wantTi)
			}
		})
	}
}
