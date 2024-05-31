package fonts

import (
	_ "embed"
	"fmt"

	"gioui.org/font"
	"gioui.org/font/opentype"
)

var (

	//go:embed NotoSans-Regular.ttf
	fontNotoSansRegular []byte
	//go:embed NotoSans-SemiBold.ttf
	fontNotoSansSemiBold []byte
	//go:embed NotoSans-Bold.ttf
	fontNotoSansBold []byte

	//go:embed NotoSansSC-Regular.ttf
	fontNotoSansSCRegular []byte
	//go:embed NotoSansSC-SemiBold.ttf
	fontNotoSansSCSemiBold []byte
	//go:embed NotoSansSC-Bold.ttf
	fontNotoSansSCBold []byte
)

var (
	collection []font.FontFace
)

func init() {
	register(fontNotoSansRegular)
	register(fontNotoSansSemiBold)
	register(fontNotoSansBold)
	register(fontNotoSansSCRegular)
	register(fontNotoSansSCSemiBold)
	register(fontNotoSansSCBold)
}

func Collection() []font.FontFace {
	return collection
}

func register(ttf []byte) {
	faces, err := opentype.ParseCollection(ttf)
	if err != nil {
		panic(fmt.Errorf("failed to parse font: %v", err))
	}
	collection = append(collection, faces...)
}
