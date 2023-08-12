package stringutil

import "strings"

func StringToLines(text string) []string {
	// zp := regexp.MustCompile(`[\t\n\f\r]`)
	// lines := zp.Split(text, -1)
	// return lines

	s := strings.ReplaceAll(text, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	lines := strings.Split(s, "\n")
	return lines
}

func GetSubStr(prefix, suffix, rawStr string) string {
	i1 := strings.Index(rawStr, prefix)
	if i1 < 1 {
		return ""
	}
	i2 := strings.Index(rawStr[i1+len(prefix):], suffix)
	if i2 < 0 {
		return rawStr[i1+len(prefix):]
	}
	return rawStr[i1+len(prefix):][:i2]
}
