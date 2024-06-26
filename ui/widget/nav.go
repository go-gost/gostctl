package widget

import (
	"image/color"

	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/clip"
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
	defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
	event.Op(gtx.Ops, &p.list)
	for {
		ev, ok := gtx.Event(
			key.Filter{Name: key.NameLeftArrow},
			key.Filter{Name: key.NameRightArrow},
		)
		if !ok {
			break
		}
		switch ev := ev.(type) {
		case key.Event:
			if ev.State != key.Press {
				break
			}
			switch ev.Name {
			case key.NameLeftArrow:
				current := p.Current() - 1
				if current < 0 {
					current = 0
				}
				p.SetCurrent(current)
				if current < p.list.Position.First {
					p.list.ScrollTo(current)
				}

			case key.NameRightArrow:
				current := p.Current() + 1
				if current >= len(p.btns) {
					current = len(p.btns) - 1
				}
				p.SetCurrent(current)
				if current+1 >= p.list.Position.First+p.list.Position.Count {
					p.list.ScrollTo(p.list.Position.First + 1)
				}
			}
		}
	}

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
		} else {
			btn.background = theme.Current().NavButtonBg
		}

		return layout.UniformInset(8).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return btn.Layout(gtx, th)
		})
	})
}

type NavButton struct {
	btn          widget.Clickable
	cornerRadius unit.Dp
	background   color.NRGBA
	text         i18n.Key
}

func NewNavButton(text i18n.Key) *NavButton {
	return &NavButton{
		cornerRadius: 20,
		text:         text,
	}
}

func (btn *NavButton) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return material.ButtonLayoutStyle{
		Background:   btn.background,
		CornerRadius: btn.cornerRadius,
		Button:       &btn.btn,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top:    8,
			Bottom: 8,
			Left:   16,
			Right:  16,
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			label := material.Body1(th, btn.text.Value())
			return label.Layout(gtx)
		})
	})
}
