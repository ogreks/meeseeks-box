package card

import larkcard "github.com/larksuite/oapi-sdk-go/v3/card"

// SetTag is a helper function to add a tag to the card.
func SetTag(elements []larkcard.MessageCardElement) []larkcard.MessageCardElement {
	tag := larkcard.NewMessageCardNote().
		Elements([]larkcard.MessageCardNoteElement{
			larkcard.NewMessageCardPlainText().
				Content("Mr.Meeeeks!").
				Build(),
		}).
		Build()

	elements = append(elements, tag)

	return elements
}

// SetHr is a helper function to add a hr to the card.
func SetHr(elements []larkcard.MessageCardElement) []larkcard.MessageCardElement {
	hr := larkcard.NewMessageCardHr().
		Build()

	elements = append(elements, hr)

	return elements
}

// SetHrAndTag is a helper function to add a hr and a tag to the card.
func SetHrAndTag(elements []larkcard.MessageCardElement) []larkcard.MessageCardElement {
	elements = SetHr(elements)
	elements = SetTag(elements)

	return elements
}
