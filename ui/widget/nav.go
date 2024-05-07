package widget

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/go-gost/gostctl/ui/i18n"
	"github.com/go-gost/gostctl/ui/theme"
)

type Nav struct {
	list    layout.List
	btns    []*NavButton
	current int
}

func NewNav(btns ...*NavButton) *Nav {
	return &Nav{
		btns: btns,
	}
}

func (p *Nav) Current() int {
	return p.current
}

func (p *Nav) SetCurrent(n int) {
	p.current = n
}

func (p *Nav) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	for i, btn := range p.btns {
		if btn.btn.Clicked(gtx) {
			p.current = i
			break
		}
	}

	return p.list.Layout(gtx, len(p.btns), func(gtx layout.Context, index int) layout.Dimensions {
		btn := p.btns[index]

		if p.current == index {
			btn.background = theme.Current().NavButtonContrastBg
			btn.borderWidth = 0
		} else {
			btn.background = theme.Current().Material.Bg
			btn.borderWidth = 1
		}

		return layout.UniformInset(8).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return btn.Layout(gtx, th)
		})
	})
}

type NavButton struct {
	btn          widget.Clickable
	cornerRadius unit.Dp
	borderWidth  unit.Dp
	background   color.NRGBA
	text         i18n.Key
}

func NewNavButton(text i18n.Key) *NavButton {
	return &NavButton{
		cornerRadius: 18,
		borderWidth:  1,
		text:         text,
	}
}

func (btn *NavButton) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return material.ButtonLayoutStyle{
		Background:   btn.background,
		CornerRadius: btn.cornerRadius,
		Button:       &btn.btn,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return widget.Border{
			Color:        theme.Current().NavButtonContrastBg,
			Width:        btn.borderWidth,
			CornerRadius: btn.cornerRadius,
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    9,
				Bottom: 9,
				Left:   16,
				Right:  16,
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				label := material.Body1(th, btn.text.Value())
				return label.Layout(gtx)
			})
		})
	})
}
