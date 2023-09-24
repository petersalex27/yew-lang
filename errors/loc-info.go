package errors

import "strconv"

func locInfo(path string, line, char int) string {
	line_s, char_s := strconv.Itoa(line), strconv.Itoa(char)
	return "[" + path + ":" + line_s + ":" + char_s + "]"
}
