package config

import (
	"os"
	"path"
	"path/filepath"
	"testing"
)

func Test_getConfig(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name: "test1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("getConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Http == nil || got.Grpc == nil {
				t.Errorf("config.Http is nil. %#+v", got)
			}

			if got.Http.DB == nil {
				t.Errorf("config.Http is invalid. %#+v", got.Http)
			}

			if got.Grpc.DB == nil || got.Grpc.Host == "" || got.Grpc.Port == 0 {
				t.Errorf("config.Grcp is invalid. %#+v", got.Grpc)
			}
		})
	}
}

func Test_getPath(t *testing.T) {
	cp, _ := os.Getwd()
	fp := "./config.yaml"
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{filePath: fp},
			want: path.Clean(filepath.Join(cp, fp)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPath(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getPath() = %v, want %v", got, tt.want)
			}
			t.Log(got)
		})
	}
}
