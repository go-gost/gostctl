package widget

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/go-gost/gostctl/ui/i18n"
)

type Switcher struct {
	Title i18n.Key
	b     widget.Bool
	ck    widget.Clickable
}

func (p *Switcher) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if p.ck.Clicked(gtx) && !p.b.Update(gtx) {
		p.b.Value = !p.b.Value
	}

	return material.Clickable(gtx, &p.ck, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top:    10,
			Bottom: 10,
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Flexed(1, material.Body1(th, p.Title.Value()).Layout),
				layout.Rigid(material.Switch(th, &p.b, "").Layout),
			)
		})
	})
}

func (p *Switcher) Value() bool {
	return p.b.Value
}

func (p *Switcher) SetValue(b bool) {
	p.b.Value = b
}
