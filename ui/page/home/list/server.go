package list

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/go-gost/gostctl/config"
	"github.com/go-gost/gostctl/ui/icons"
	"github.com/go-gost/gostctl/ui/page"
	"github.com/go-gost/gostctl/ui/theme"
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
	cfg := config.Get()
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
			Top:    8,
			Bottom: 8,
			Left:   8,
			Right:  8,
		}.Layout(gtx, func(gtx C) D {
			return material.ButtonLayoutStyle{
				Background:   theme.Current().ListBg,
				CornerRadius: 12,
				Button:       &l.states[index].btn,
			}.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(16).Layout(gtx, func(gtx C) D {
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
								layout.Rigid(layout.Spacer{Height: 4}.Layout),
								layout.Rigid(material.Body2(th, servers[index].URL).Layout),
								layout.Rigid(layout.Spacer{Height: 4}.Layout),
								layout.Rigid(func(gtx C) D {
									return layout.Flex{
										Spacing:   layout.SpaceBetween,
										Alignment: layout.Middle,
									}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											gtx.Constraints.Min.X = gtx.Dp(16)
											return icons.IconActionUpdate.Layout(gtx, th.Fg)
										}),
										layout.Rigid(layout.Spacer{Width: 4}.Layout),
										layout.Flexed(1, material.Body2(th, servers[index].Interval.String()).Layout),
									)
								}),
								layout.Rigid(layout.Spacer{Height: 4}.Layout),
								layout.Rigid(func(gtx C) D {
									return layout.Flex{
										Spacing:   layout.SpaceBetween,
										Alignment: layout.Middle,
									}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											gtx.Constraints.Min.X = gtx.Dp(16)
											return icons.IconActionHourGlassEmpty.Layout(gtx, th.Fg)
										}),
										layout.Rigid(layout.Spacer{Width: 4}.Layout),
										layout.Flexed(1, material.Body2(th, servers[index].Timeout.String()).Layout),
									)
								}),
							)
						}),
						layout.Rigid(layout.Spacer{Width: 4}.Layout),
						layout.Rigid(func(gtx C) D {
							if index == cfg.CurrentServer {
								gtx.Constraints.Min.X = gtx.Dp(12)
								if state := servers[index].State(); state == config.ServerError {
									return icons.IconCircle.Layout(gtx, color.NRGBA(colornames.Red500))
								}
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
