package service

import (
	"errors"
	"testing"
	"time"
)

func TestResolveHoldDuration(t *testing.T) {
	tests := []struct {
		name     string
		minutes  int
		fallback time.Duration
		want     time.Duration
		wantErr  bool
	}{
		{name: "uses configured default", minutes: 0, fallback: 10 * time.Minute, want: 10 * time.Minute},
		{name: "accepts one minute", minutes: 1, fallback: 10 * time.Minute, want: time.Minute},
		{name: "accepts max duration", minutes: 60, fallback: 10 * time.Minute, want: 60 * time.Minute},
		{name: "rejects negative duration", minutes: -1, fallback: 10 * time.Minute, wantErr: true},
		{name: "rejects duration over one hour", minutes: 61, fallback: 10 * time.Minute, wantErr: true},
		{name: "rejects missing configured default", minutes: 0, fallback: 0, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveHoldDuration(tt.minutes, tt.fallback)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidHoldDuration) {
					t.Fatalf("expected ErrInvalidHoldDuration, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %s, got %s", tt.want, got)
			}
		})
	}
}
