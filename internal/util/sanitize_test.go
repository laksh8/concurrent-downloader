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

func TestSafeFileName_WindowsReservedNames(t *testing.T) {
	fallback := "default.txt"

	cases := []string{
		"CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9",
	}

	for _, name := range cases {
		url := "https://example.com/" + name
		got := SafeFileName(url, fallback)

		if got != fallback {
			t.Errorf("expected fallback for %s, got %s", name, got)
		}
	}
}
