package message

import (
	"fmt"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/ogreks/meeseeks-box/internal/pkg/feishu/aide"
)

type HelpAction struct {
	BaseAction
}

func (*HelpAction) Helper() []string {
	return []string{"🙊 **需要帮助**", "文本回复 *帮助* 或者 */help*、*/h*"}
}

func (*HelpAction) Execute(a *ActionInfo, m MessageHandleInterface) bool {
	if _, foundHelp := aide.EitherTrimEqual(a.Info.QParsed, "/help", "/h", "帮助", "命令辅助", "辅助"); foundHelp {
		fmt.Printf("%v\n", m.(*MessageHandle).actions)
		m.Reply(
			*a.Ctx,
			*a.Info.SessionId,
			SendHelperCard(m.(*MessageHandle).actions...),
			larkim.MsgTypeInteractive,
		)

		return false
	}

	return true
}

type EmptyAction struct {
	BaseAction
}

func (EmptyAction) Execute(a *ActionInfo, m MessageHandleInterface) bool {
	if len(a.Info.QParsed) == 0 {
		m.Reply(
			*a.Ctx,
			*a.Info.SessionId,
			"🤖️：你想知道什么呢~",
			larkim.MsgTypeText,
		)
		return false
	}

	return true
}
