package message

import (
	"context"
	"encoding/json"
	"fmt"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	"github.com/ogreks/meeseeks-box/internal/pkg/feishu/aide"
	expand_card "github.com/ogreks/meeseeks-box/internal/pkg/feishu/expand/card"
)

type Action interface {
	Helper() []string
	Execute(a *ActionInfo, m MessageHandleInterface) bool
}

type BaseAction struct{}

func (b BaseAction) Helper() []string { return nil }

type ActionInfo struct {
	Handler MessageHandleInterface
	Ctx     *context.Context
	Info    *Message
}

func helperLayout(title, descript string) []byte {
	var columnMap = map[string]interface{}{
		"tag":                "column_set",
		"flex_mode":          "bisect",
		"background_style":   "default",
		"horizontal_spacing": "default",
		"columns": []map[string]interface{}{
			{
				"tag":            "column",
				"width":          "weighted",
				"weight":         1,
				"vertical_align": "top",
				"elements": []map[string]interface{}{
					{
						"tag":     "markdown",
						"content": title,
					},
				},
			},
			{
				"tag":            "column",
				"width":          "weighted",
				"weight":         2,
				"vertical_align": "top",
				"elements": []map[string]interface{}{
					{
						"tag":     "markdown",
						"content": descript,
					},
				},
			},
		},
	}

	data, _ := json.Marshal(columnMap)
	return data
}

func SendHelperCard(actions ...Action) string {
	var elements []larkcard.MessageCardElement = []larkcard.MessageCardElement{
		aide.WithCardMdContent("**🤠你好呀~ 我是使命必达盒，一个出自瑞克莫蒂动漫的虚拟助手！**\n\n - 下面是支持的命令"),
		aide.WithCardHr(),
	}

	for i := 0; i < len(actions); i++ {
		action := actions[i]

		fmt.Printf("helper: %+v, helper: %s\n", action, action.Helper())
		if len(action.Helper()) < 2 {
			continue
		}

		help := action.Helper()
		elements = append(
			elements,
			expand_card.NewCardColumnSet(helperLayout(help[0], help[1])),
		)
	}

	newCard, _ := aide.NewSendCard(
		aide.WithCardHeader(larkcard.TemplateBlue, "🎒需要帮助吗？"),
		elements...,
	)

	return newCard
}
