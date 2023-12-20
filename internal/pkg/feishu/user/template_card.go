package feishu

import (
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
)

func templateI18nCardElements(zhCn, enUs string) *larkcard.MessageCardI18nElements {
	zhCnElements := make([]larkcard.MessageCardElement, 1)
	zhCnElements[0] = larkcard.NewMessageCardMarkdown().
		Content(zhCn).
		Build()

	enUsElements := make([]larkcard.MessageCardElement, 1)
	enUsElements[0] = larkcard.NewMessageCardMarkdown().
		Content(enUs).
		Build()

	return larkcard.NewMessageCardI18nElements().
		ZhCN(zhCnElements).
		EnUS(enUsElements).
		Build()
}

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

	element := templateI18nCardElements(content_zh, content_en)

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

//func VersionTemplate(version string) string {
//	content_zh := `当前版本：{{ .Version }}
//[点我查看](https://no0overtime0group.feishu.cn/docx/TQSkdZizGoeFbmxe0apcQncdnMe)`
//	content_en := `Current: {{ .Version }}
//[Click here to view](https://no0overtime0group.feishu.cn/docx/TQSkdZizGoeFbmxe0apcQncdnMe)
//	`
//}
