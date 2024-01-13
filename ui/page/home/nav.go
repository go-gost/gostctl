package home

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/go-gost/gui/ui/page/home/list"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

type nav struct {
	list    layout.List
	btns    []*navButton
	current int
}

func (p *nav) Layout(gtx C, th *material.Theme) D {
	for i, btn := range p.btns {
		if btn.btn.Clicked(gtx) {
			p.current = i
			break
		}
	}

	return p.list.Layout(gtx, len(p.btns), func(gtx C, index int) D {
		btn := p.btns[index]
		if p.current == index {
			btn.Color = color.NRGBA(colornames.Black)
			btn.Background = color.NRGBA(colornames.BlueGrey100)
			btn.BorderWidth = 0
		} else {
			btn.Color = color.NRGBA(colornames.Grey700)
			btn.Background = color.NRGBA(colornames.White)
			btn.BorderWidth = 1
		}

		return layout.Inset{
			Top:    5,
			Bottom: 5,
			Left:   10,
			Right:  10,
		}.Layout(gtx, func(gtx C) D {
			return btn.Layout(gtx, th)
		})
	})
}

type navButton struct {
	btn          widget.Clickable
	CornerRadius unit.Dp
	BorderWidth  unit.Dp
	BorderColor  color.NRGBA
	Color        color.NRGBA
	Background   color.NRGBA
	Text         string
	List         list.List
}

func NavButton(text string, list list.List) *navButton {
	return &navButton{
		CornerRadius: 18,
		BorderWidth:  1,
		BorderColor:  color.NRGBA(colornames.Grey200),
		Text:         text,
		List:         list,
	}
}

func (btn *navButton) Layout(gtx C, th *material.Theme) D {
	return material.ButtonLayoutStyle{
		Background:   btn.Background,
		CornerRadius: btn.CornerRadius,
		Button:       &btn.btn,
	}.Layout(gtx, func(gtx C) D {
		return widget.Border{
			Color:        btn.BorderColor,
			Width:        btn.BorderWidth,
			CornerRadius: btn.CornerRadius,
		}.Layout(gtx, func(gtx C) D {
			return layout.Inset{
				Top:    8,
				Bottom: 8,
				Left:   20,
				Right:  20,
			}.Layout(gtx, func(gtx C) D {
				label := material.Body1(th, btn.Text)
				label.Color = btn.Color
				return label.Layout(gtx)
			})
		})
	})
}
