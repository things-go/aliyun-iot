package infra

import (
	"testing"
)

func TestHmac(t *testing.T) {
	type args struct {
		method string
		key    string
		val    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "invalid method",
			args: args{"", "thinkgos", "thinkgos"},
			want: "thinkgos",
		},
		{
			name: "md5",
			args: args{"hmacmd5", "thinkgos", "thinkgos"},
			want: "ed9dc7baf84a9740a3ceb6f1f26bfb4f",
		},
		{
			name: "sha1",
			args: args{"hmacsha1", "thinkgos", "thinkgos"},
			want: "f9cebe2044ea375cff1a46f4dc05eb15ff9870ee",
		},
		{
			name: "sha224",
			args: args{"hmacsha224", "thinkgos", "thinkgos"},
			want: "4c95a98768fabdff9756e2e92eda74ca062e00532c8c42eb67481701",
		},
		{
			name: "sha256",
			args: args{"hmacsha256", "thinkgos", "thinkgos"},
			want: "e9403e3a615fad72d1dd1fe90c225cbec4ba81a03e5474d91a72844d2218954f",
		},
		{
			name: "sha384",
			args: args{"hmacsha384", "thinkgos", "thinkgos"},
			want: "ba0c2ef006f64e7db43fb085abac27e960ed1c43e4604838c10ead6ffaa31dfd139f66343e1db84027271a267428ebbd",
		},
		{
			name: "sha512",
			args: args{"hmacsha512", "", "thinkgos"},
			want: "c2181a3b42befba66ac95cf9fc6e11971c8bac0ec25bbf6805342b8166dd450e18ca5872e011ef1dd56bb960d96e7f93e1d2593d84f78e8e9a88892630393ce2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Hmac(tt.args.method, tt.args.key, tt.args.val); got != tt.want {
				t.Errorf("Hmac(%s) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
