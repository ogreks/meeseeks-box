package message

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/ogreks/meeseeks-box/internal/pkg/feishu/aide"
)

type (
	CardKind     string
	CardChatType string
)

var (
	GroupChatType = CardChatType("group")
	UserChatType  = CardChatType("personal")
)

const (
	GroupHandler = "group"
	UserHandler  = "personal"
)

type Message struct {
	HandlerType string
	MsgType     string
	MsgId       *string
	ChatId      *string
	QParsed     string
	FileKey     string
	ImageKey    string
	ImageKeys   []string // post 消息图片组
	SessionId   *string
	Mention     []*larkim.MentionEvent
}

type MessageHandleInterface interface {
	RegisterRoute(path string, g gin.IRouter)

	RegisterAction(action Action) MessageHandleInterface
	RegisterActions(actions ...Action) MessageHandleInterface

	Reply(ctx context.Context, messageId string, message string, msgType string) error
}

func (m *MessageHandle) RegisterAction(action Action) MessageHandleInterface {
	m.actions = append(m.actions, action)
	return m
}

func (m *MessageHandle) RegisterActions(actions ...Action) MessageHandleInterface {
	m.actions = append(m.actions, actions...)
	return m
}

func (m *MessageHandle) Reply(ctx context.Context, messageId string, message string, msgType string) error {
	replyBody := larkim.NewReplyMessageReqBodyBuilder().
		MsgType(msgType).
		Content(message).
		Build()

	replyReq := larkim.NewReplyMessageReqBuilder().
		MessageId(messageId).
		Body(replyBody).
		Build()

	resp, err := m.cli.Im.Message.Reply(ctx, replyReq)
	if err != nil {
		return err
	}

	fmt.Println(resp)
	return nil
}

// chain of responsibility
func (m *MessageHandle) chain(data *ActionInfo, actions ...Action) bool {
	for _, v := range actions {
		if !v.Execute(data, m) {
			return false
		}
	}
	return true
}

// judgeMsgType
func judgeMsgType(event *larkim.P2MessageReceiveV1) (string, error) {
	msgType := event.Event.Message.MessageType

	switch *msgType {
	case "text", "image", "audio", "post":
		return *msgType, nil
	default:
		return "", fmt.Errorf("unknow message type: %v", *msgType)
	}
}

func judgeCardType(cardAction *larkcard.CardAction) string {
	actionValue := cardAction.Action.Value
	chatType := actionValue["chatType"]
	if chatType == "group" {
		return GroupHandler
	}

	if chatType == "personal" {
		return UserHandler
	}

	return "otherChat"
}

func judgeChatType(event *larkim.P2MessageReceiveV1) string {
	chatType := event.Event.Message.ChatType
	if *chatType == "group" {
		return GroupHandler
	}

	if *chatType == "p2p" {
		return UserHandler
	}

	return "otherChat"
}

func (m *MessageHandle) msgReceivedHandler(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	handleType := judgeChatType(event)
	// log
	if handleType == "otherChat" {
		fmt.Println("unknow chat type")
		return nil
	}

	fmt.Printf("收到消息: %s\v", larkcore.Prettify(event.Event.Message))

	msgType, err := judgeMsgType(event)
	if err != nil {
		fmt.Printf("error getting message type: %v\n", err)
		return nil
	}

	content := event.Event.Message.Content
	msgId := event.Event.Message.MessageId
	rootId := event.Event.Message.RootId
	chatId := event.Event.Message.ChatId
	mention := event.Event.Message.Mentions

	sessionId := rootId
	if sessionId == nil || *sessionId == "" {
		sessionId = msgId
	}

	msg := Message{
		HandlerType: handleType,
		MsgType:     msgType,
		MsgId:       msgId,
		ChatId:      chatId,
		QParsed:     strings.Trim(aide.ParseContent(*content, msgType), " "),
		FileKey:     aide.ParseFileKey(*content),
		ImageKey:    aide.ParseImageKey(*content),
		ImageKeys:   aide.ParsePostImageKeys(*content),
		SessionId:   sessionId,
		Mention:     mention,
	}

	data := &ActionInfo{
		Ctx:     &ctx,
		Handler: m,
		Info:    &msg,
	}

	m.chain(data, m.actions...)

	return nil
}
