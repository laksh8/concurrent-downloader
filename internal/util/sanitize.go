package util

import (
	"net/url"
	"path"
	"strings"
	"unicode"
)

func SafeFileName(rawurl string, fallback string) string {

	u, err := url.Parse(rawurl)
	if err != nil {
		return fallback
	}

	str := strings.TrimSpace(path.Base(u.Path))
	res := strings.Builder{}

	for _, c := range str {
		if unicode.IsLetter(c) || c == '_' || unicode.IsDigit(c) || c == '-' || c == '.' {
			res.WriteRune(c)
		} else {
			res.WriteRune('_')
		}
	}

	result := strings.Trim(res.String(), "_")

	if result == "" || result == "." {
		return fallback
	}
	if result[0] == '.' { // Hidden files under unix systems
		result = "file" + result
	}

	switch strings.ToUpper(result) { // Windows reserved keywords
	case "CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9":
		return fallback
	}

	return result
}
