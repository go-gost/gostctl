package list

import (
	"fmt"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/go-gost/gui/config"
	"github.com/go-gost/gui/ui/icons"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

type serverState struct {
	btn widget.Clickable
}

type serverList struct {
	list   layout.List
	states []serverState
}

func Server() List {
	return &serverList{
		list: layout.List{
			Axis: layout.Vertical,
		},
		states: make([]serverState, 16),
	}
}

func (l *serverList) Layout(gtx C, th *material.Theme) D {
	servers := config.Global().Servers
	if len(servers) > len(l.states) {
		states := l.states
		l.states = make([]serverState, len(servers))
		copy(l.states, states)
	}

	return l.list.Layout(gtx, len(servers), func(gtx C, index int) D {
		return layout.Inset{
			Top:    5,
			Bottom: 5,
			Left:   10,
			Right:  10,
		}.Layout(gtx, func(gtx C) D {
			return material.ButtonLayoutStyle{
				Background:   color.NRGBA(colornames.BlueGrey50),
				CornerRadius: 10,
				Button:       &l.states[index].btn,
			}.Layout(gtx, func(gtx C) D {
				return layout.Inset{
					Top:    8,
					Bottom: 8,
					Left:   10,
					Right:  10,
				}.Layout(gtx, func(gtx C) D {
					return layout.Flex{
						Alignment: layout.Middle,
						Spacing:   layout.SpaceBetween,
					}.Layout(gtx,
						layout.Flexed(1, func(gtx C) D {
							return layout.Flex{
								Axis: layout.Vertical,
							}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									label := material.Body1(th, servers[index].Name)
									label.Font.Weight = font.SemiBold
									return label.Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Height: 10}.Layout),
								layout.Rigid(material.Body2(th, fmt.Sprintf("URL: %s", servers[index].URL)).Layout),
								layout.Rigid(layout.Spacer{Height: 5}.Layout),
								layout.Rigid(material.Body2(th, fmt.Sprintf("Interval: %s", servers[index].Interval)).Layout),
								layout.Rigid(layout.Spacer{Height: 5}.Layout),
								layout.Rigid(material.Body2(th, fmt.Sprintf("Timeout: %s", servers[index].Timeout)).Layout),
							)
						}),
						layout.Rigid(layout.Spacer{Width: 10}.Layout),
						layout.Rigid(func(gtx C) D {
							return icons.IconDone.Layout(gtx, color.NRGBA(colornames.Green500))
						}),
					)
				})
			})
		})
	})
}
