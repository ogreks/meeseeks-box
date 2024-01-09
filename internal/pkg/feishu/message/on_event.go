package message

import (
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
		m.ReplyJudgeMessage(
			*a.Ctx,
			SendHelperCard(m.(*MessageHandle).actions...),
			larkim.MsgTypeInteractive,
			a.Info,
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
		m.ReplyJudgeMessage(*a.Ctx, larkim.NewMessageTextBuilder().
			Text("🤖️：你想知道什么呢~").
			Build(),
			larkim.MsgTypeText,
			a.Info,
		)
		return false
	}

	return true
}

type VersionAction struct {
	BaseAction
}

func (*VersionAction) Helper() []string {
	return []string{"🔘 **查看版本**", "文本回复 *当前版本* 或者 *版本*、*/version*、*/v*"}
}

func (*VersionAction) Execute(a *ActionInfo, m MessageHandleInterface) bool {
	_, foundHelp := aide.EitherTrimEqual(
		a.Info.QParsed,
		"/version", "/v", "当前版本", "版本", "目前版本",
	)
	if foundHelp {
		m.ReplyJudgeMessage(
			*a.Ctx,
			larkim.NewTextMsgBuilder().
				Text("🔘当前版本：dev-0.0.1-pre 版本号详解：【https://no0overtime0group.feishu.cn/docx/TQSkdZizGoeFbmxe0apcQncdnMe】").
				Build(),
			larkim.MsgTypeText,
			a.Info,
		)
		return false
	}

	return true
}

// StartAction /start
type StartAction struct {
	BaseAction
}

func (*StartAction) Execute(a *ActionInfo, m MessageHandleInterface) bool {
	if a.Info.HandlerType != UserHandler {
		return true
	}

	_, e := aide.EitherTrimEqual(
		a.Info.QParsed,
		"/start",
	)

	if e {
		m.SendMessage(
			*a.Ctx,
			SendStartCard(),
			larkim.MsgTypeInteractive,
			a.Info,
		)
		return false
	}

	return true
}
