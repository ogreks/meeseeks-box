package aide

import (
	"encoding/json"
	"fmt"
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

func MsgFilter(msg string) string {
	regex := regexp.MustCompile(`@[^ ]*`)
	return regex.ReplaceAllString(msg, "")
}

func ParsePostContent(content string) string {
	var contentMap map[string]interface{}
	err := json.Unmarshal([]byte(content), &contentMap)

	if err != nil {
		fmt.Println(err)
	}

	if contentMap["content"] == nil {
		return ""
	}

	var text string
	if contentMap["title"] != nil && contentMap["title"] != "" {
		text += contentMap["title"].(string) + "\n"
	}

	contentList := contentMap["content"].([]interface{})
	for _, v := range contentList {
		for _, v1 := range v.([]interface{}) {
			if v1.(map[string]interface{})["tag"] == "text" {
				text += v1.(map[string]interface{})["text"].(string)
			}
		}
		// add new line
		text += "\n"
	}

	return MsgFilter(text)
}

func ParseContent(content, msgType string) string {
	//"{\"text\":\"@_user_1  hahaha\"}",
	//only get text content hahaha
	if msgType == "post" {
		return ParsePostContent(content)
	}

	var contentMap map[string]interface{}
	err := json.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		fmt.Println(err)
	}

	if contentMap["text"] == nil {
		return ""
	}

	text := contentMap["text"].(string)
	return MsgFilter(text)
}

func ParseFileKey(content string) string {
	var contentMap map[string]interface{}
	err := json.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if contentMap["file_key"] == nil {
		return ""
	}

	fileKey := contentMap["file_key"].(string)
	return fileKey
}

func ParseImageKey(content string) string {
	var contentMap map[string]interface{}
	err := json.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if contentMap["image_key"] == nil {
		return ""
	}

	imageKey := contentMap["image_key"].(string)
	return imageKey
}

func ParsePostImageKeys(content string) []string {
	var contentMap map[string]interface{}
	err := json.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var imageKeys []string
	if contentMap["content"] == nil {
		return imageKeys
	}

	contentList := contentMap["content"].([]interface{})
	for _, v := range contentList {
		for _, v1 := range v.([]interface{}) {
			if v1.(map[string]interface{})["tag"] == "img" {
				imageKeys = append(imageKeys, v1.(map[string]interface{})["image_key"].(string))
			}
		}
	}

	return imageKeys
}
