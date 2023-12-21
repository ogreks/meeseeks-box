package aide

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

// ProcessMessage
func ProcessMessage(msg interface{}) (string, error) {
	msg = strings.TrimSpace(msg.(string))
	msgB, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	msgStr := string(msgB)
	if len(msgStr) > 2 {
		msgStr = msgStr[1 : len(msgStr)-1]
	}

	return msgStr, nil
}

// ProcessNewLine replace \n
func ProcessNewLine(msg string) string {
	return strings.Replace(msg, "\\n", `
`, -1)
}

// ProcessQuote replace \\\"
func ProcessQuote(msg string) string {
	return strings.Replace(msg, "\\\"", "\"", -1)
}

// ProcessUnicode
func ProcessUnicode(msg string) string {
	regex := regexp.MustCompile(`\\u[0-9a-fA-F]{4}`)
	return regex.ReplaceAllStringFunc(msg, func(s string) string {
		r, _ := regexp.Compile(`\\u`)
		s = r.ReplaceAllString(s, "")
		i, _ := strconv.ParseInt(s, 16, 32)
		return string(rune(i))
	})
}

// CleanTextBlock msg
func CleanTextBlock(msg string) string {
	msg = ProcessNewLine(msg)
	msg = ProcessUnicode(msg)
	msg = ProcessQuote(msg)

	return msg
}
