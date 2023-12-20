package feishu

import larkcard "github.com/larksuite/oapi-sdk-go/v3/card"

func TemplateUnopenedAbilityCard() string {
	content_zh := `这个功能极为重要 ***@开发者***
这个功能可以暂缓 **[!-_-!]**
这个功能什么时候能有~ 
[点我前往查看](https://no0overtime0group.feishu.cn/docx/TQSkdZizGoeFbmxe0apcQncdnMe)`
	content_en := `This function is extremely important ***@Developer***
This function can be postponed **[!-_-!]**
When will this function be available~
[Click here to view](https://no0overtime0group.feishu.cn/docx/TQSkdZizGoeFbmxe0apcQncdnMe)
	`
	element := larkcard.NewMessageCardI18nElements().
		ZhCN([]larkcard.MessageCardElement{
			larkcard.NewMessageCardMarkdown().
				Content(content_zh).
				Build(),
			larkcard.NewMessageCardHr().
				Build(),
			larkcard.NewMessageCardNote().
				Elements([]larkcard.MessageCardNoteElement{
					larkcard.NewMessageCardPlainText().
						Content("Mr.Meeeeks!").
						Build(),
				}).
				Build(),
		}).
		EnUS([]larkcard.MessageCardElement{
			larkcard.NewMessageCardMarkdown().
				Content(content_en).
				Build(),
			larkcard.NewMessageCardHr().
				Build(),
			larkcard.NewMessageCardNote().
				Elements([]larkcard.MessageCardNoteElement{
					larkcard.NewMessageCardPlainText().
						Content("Mr.Meeeeks!").
						Build(),
				}).
				Build(),
		}).
		Build()

	header := larkcard.NewMessageCardHeader().
		Template(larkcard.TemplateBlue).
		Title(larkcard.NewMessageCardPlainText().
			I18n(
				larkcard.NewMessageCardPlainTextI18n().
					ZhCN("🤡 这个功能还未开发!!!").
					EnUS("🤡 This feature has not been developed yet!!!").
					Build(),
			).
			Build(),
		).
		Build()

	config := larkcard.NewMessageCardConfig().
		WideScreenMode(true).
		UpdateMulti(false).
		Build()

	data, _ := larkcard.NewMessageCard().
		Config(config).
		Header(header).
		I18nElements(element).
		String()

	return data
}
