package build

import (
	"testing"
	"time"
)

func TestFormatDisplayDuration(t *testing.T) {
	tests := []struct {
		name string
		d    time.Duration
		want string
	}{
		{
			name: "sub-second nanoseconds to integer ms",
			d:    949802583 * time.Nanosecond,
			want: "949ms",
		},
		{
			name: "just over one second to two decimal places",
			d:    1366963417 * time.Nanosecond,
			want: "1.37s",
		},
		{
			name: "microseconds unchanged integer",
			d:    500 * time.Microsecond,
			want: "500µs",
		},
		{
			name: "integer milliseconds",
			d:    42 * time.Millisecond,
			want: "42ms",
		},
		{
			name: "exact one second",
			d:    time.Second,
			want: "1s",
		},
		{
			name: "sub-millisecond nanoseconds to integer ms",
			d:    1 * time.Millisecond,
			want: "1ms",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDisplayDuration(tt.d)
			if got != tt.want {
				t.Fatalf("formatDisplayDuration(%v) = %q, want %q", tt.d, got, tt.want)
			}
		})
	}
}