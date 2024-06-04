package i18n

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

var (
	cat catalog.Catalog
)

func init() {
	builder := catalog.NewBuilder()
	for k, v := range en_US {
		builder.SetString(language.English, string(k), v)
	}
	for k, v := range zh_CN {
		builder.SetString(language.Chinese, string(k), v)
	}
	cat = builder
}
