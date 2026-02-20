package util

import "testing"

func TestSafeFilename(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{"simple", "https://example.com/img/pic.jpg", "pic.jpg"},
		{"query", "https://site/x?y=1", "x"},
		{"no-path", "https://site", "fallback.bin"},
		{"unsafe", "https://a/%00%20bad:name.png", "bad_name.png"},
		{"hidden-file", "comeone.com/_.png", "file.png"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SafeFileName(tt.url, "fallback.bin")
			if got != tt.want {
				t.Fatalf("got %q want %q", got, tt.want)
			}
		})
	}
}
