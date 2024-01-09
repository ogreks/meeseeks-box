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

// Helper help message
func (b BaseAction) Helper() []string { return nil }

// Execute run server
func (b BaseAction) Execute(a *ActionInfo, m MessageHandleInterface) bool { return true }

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

// SendStartCard
func SendStartCard() string {
	newCard, _ := aide.NewSendCard(
		aide.WithCardHeader(larkcard.TemplateBlue, "📒Hello，World!"),
		aide.WithCardMdContent(
			`**🤠你好呀~ 我是使命必达盒，一个出自瑞克莫蒂动漫的虚拟助手！**

欢迎使用 Meeseeks Box's，我支持很多奇奇怪怪的功能，当然这些都需要你去发掘

同时我也支持基础任务的编辑和操作 **目前仍在开发中**

[点我前往查看](https://no0overtime0group.feishu.cn/docx/TQSkdZizGoeFbmxe0apcQncdnMe)

🤖 如果你想知道我的快捷指令请发送 **/help**`,
		),
	)

	return newCard
}
