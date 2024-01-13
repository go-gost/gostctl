package list

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/go-gost/gui/config"
	"github.com/go-gost/gui/ui/icons"
	"github.com/go-gost/gui/ui/page"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

type serverState struct {
	btn widget.Clickable
}

type serverList struct {
	router *page.Router
	list   layout.List
	states []serverState
}

func Server(r *page.Router) List {
	return &serverList{
		router: r,
		list: layout.List{
			Axis: layout.Vertical,
		},
		states: make([]serverState, 16),
	}
}

func (l *serverList) Layout(gtx C, th *material.Theme) D {
	cfg := config.Global()
	servers := cfg.Servers
	if len(servers) > len(l.states) {
		states := l.states
		l.states = make([]serverState, len(servers))
		copy(l.states, states)
	}

	return l.list.Layout(gtx, len(servers), func(gtx C, index int) D {
		if l.states[index].btn.Clicked(gtx) {
			l.router.Goto(page.Route{
				Path: page.PageServerEdit,
				ID:   servers[index].Name,
			})
		}

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
								layout.Rigid(material.Body2(th, servers[index].URL).Layout),
								layout.Rigid(layout.Spacer{Height: 5}.Layout),
								layout.Rigid(func(gtx C) D {
									return layout.Flex{
										Spacing:   layout.SpaceBetween,
										Alignment: layout.Middle,
									}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											gtx.Constraints.Min.X = gtx.Dp(16)
											return icons.IconActionUpdate.Layout(gtx, color.NRGBA(colornames.Grey800))
										}),
										layout.Rigid(layout.Spacer{Width: 5}.Layout),
										layout.Flexed(1, material.Body2(th, servers[index].Interval.String()).Layout),
									)
								}),
								layout.Rigid(layout.Spacer{Height: 5}.Layout),
								layout.Rigid(func(gtx C) D {
									return layout.Flex{
										Spacing:   layout.SpaceBetween,
										Alignment: layout.Middle,
									}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											gtx.Constraints.Min.X = gtx.Dp(16)
											return icons.IconActionHourGlassEmpty.Layout(gtx, color.NRGBA(colornames.Grey800))
										}),
										layout.Rigid(layout.Spacer{Width: 5}.Layout),
										layout.Flexed(1, material.Body2(th, servers[index].Timeout.String()).Layout),
									)
								}),
							)
						}),
						layout.Rigid(layout.Spacer{Width: 10}.Layout),
						layout.Rigid(func(gtx C) D {
							if index == cfg.CurrentServer {
								gtx.Constraints.Min.X = gtx.Dp(15)
								return icons.IconCircle.Layout(gtx, color.NRGBA(colornames.Green500))
							}
							return D{}
						}),
					)
				})
			})
		})
	})
}
