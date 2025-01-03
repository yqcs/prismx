package netutil

import (
	"testing"
)

func TestTryJoinHostPort(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		port    string
		want    string
		wantErr bool
	}{
		{
			name:    "both host and port provided",
			host:    "localhost",
			port:    "8080",
			want:    "localhost:8080",
			wantErr: false,
		},
		{
			name:    "empty host",
			host:    "",
			port:    "8080",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty port",
			host:    "localhost",
			port:    "",
			want:    "localhost",
			wantErr: true,
		},
		{
			name:    "both host and port empty",
			host:    "",
			port:    "",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TryJoinHostPort(tt.host, tt.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("TryJoinHostPort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TryJoinHostPort() = %v, want %v", got, tt.want)
			}
		})
	}
}
