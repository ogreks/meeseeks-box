package aide

import (
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
)

// WithCardHr set card hr
func WithCardHr() larkcard.MessageCardElement {
	return larkcard.NewMessageCardHr().Build()
}

// WithCardHeader set card header
func WithCardHeader(color string, title ...string) *larkcard.MessageCardHeader {
	if len(title) == 0 {
		title[0] = "🤖️机器人提醒"
	}

	return larkcard.NewMessageCardHeader().
		Template(color).
		Title(larkcard.NewMessageCardPlainText().
			Content(title[0]).
			Build(),
		).
		Build()
}

// WithCardNote set note
func WithCardNote(note string) larkcard.MessageCardElement {
	return larkcard.NewMessageCardNote().
		Elements([]larkcard.MessageCardNoteElement{
			larkcard.NewMessageCardPlainText().Content(note).Build(),
		}).
		Build()
}

// WithCardMdContent generate markdown message
func WithCardMdContent(msg string) larkcard.MessageCardElement {
	msg, err := ProcessMessage(msg)
	if err != nil {
		return nil
	}
	msg = ProcessNewLine(msg)

	return larkcard.NewMessageCardDiv().
		Fields([]*larkcard.MessageCardField{
			larkcard.NewMessageCardField().
				Text(larkcard.NewMessageCardLarkMd().
					Content(msg).
					Build(),
				).
				IsShort(true).
				Build(),
		}).
		Build()
}

// WithCardPlainText
func WithCardPlainText(msg string) larkcard.MessageCardElement {
	msg, err := ProcessMessage(msg)
	if err != nil {
		return nil
	}
	msg = ProcessNewLine(msg)

	return larkcard.NewMessageCardDiv().
		Fields([]*larkcard.MessageCardField{
			larkcard.NewMessageCardField().
				Text(larkcard.NewMessageCardPlainText().
					Content(msg).
					Build()).
				IsShort(true).
				Build(),
		}).
		Build()
}

func NewSendCard(header *larkcard.MessageCardHeader, elements ...larkcard.MessageCardElement) (string, error) {
	config := larkcard.NewMessageCardConfig().
		WideScreenMode(false).
		EnableForward(true).
		UpdateMulti(false).
		Build()

	var elementPool []larkcard.MessageCardElement
	elementPool = append(elementPool, elements...)
	elementPool = append(elementPool, WithCardHr(), WithCardNote("Mr.Meeseeks!"))

	return larkcard.NewMessageCard().
		Config(config).
		Header(header).
		Elements(elementPool).
		String()
}
