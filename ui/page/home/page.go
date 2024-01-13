package home

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/go-gost/gui/ui/icons"
	"github.com/go-gost/gui/ui/page"
	"github.com/go-gost/gui/ui/page/home/list"
)

type C = layout.Context
type D = layout.Dimensions

type homePage struct {
	router      *page.Router
	nav         nav
	btnCreate   widget.Clickable
	btnSettings widget.Clickable
}

func NewPage(r *page.Router) page.Page {
	return &homePage{
		router: r,
		nav: nav{
			btns: []*navButton{
				NavButton("Server", list.Server(r), page.PageServerEdit),
				NavButton("Service", list.Service(r), page.PageServiceEdit),
				NavButton("Chain", nil, ""),
				NavButton("Hop", nil, ""),
				NavButton("Auther", nil, ""),
				NavButton("Admission", nil, ""),
				NavButton("Bypass", nil, ""),
				NavButton("Resolver", nil, ""),
				NavButton("Hosts", nil, ""),
				NavButton("Limiter", nil, ""),
				NavButton("Ingress", nil, ""),
				NavButton("Logger", nil, ""),
			},
			current: 1,
		},
	}
}

func (p *homePage) Init(opts ...page.PageOption) {

}

func (p *homePage) Layout(gtx C, th *material.Theme) D {
	if p.btnCreate.Clicked(gtx) {
		p.router.Goto(page.Route{Path: p.nav.btns[p.nav.current].Create})
	}

	return layout.Stack{
		Alignment: layout.SE,
	}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				// header
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Top:    5,
						Bottom: 5,
						Left:   10,
						Right:  10,
					}.Layout(gtx, func(gtx C) D {
						return layout.Flex{
							Spacing:   layout.SpaceBetween,
							Alignment: layout.Middle,
						}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return icons.IconApp.Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Width: 5}.Layout),
							layout.Rigid(func(gtx C) D {
								label := material.H6(th, "GOST")
								label.Font.Weight = font.Bold
								return label.Layout(gtx)
							}),
							layout.Flexed(1, layout.Spacer{}.Layout),
							layout.Rigid(func(gtx C) D {
								btn := material.IconButton(th, &p.btnSettings, icons.IconSettings, "Settings")
								btn.Color = th.Fg
								btn.Background = th.Bg
								return btn.Layout(gtx)
							}),
						)
					})
				}),
				// nav
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Top:    5,
						Bottom: 5,
					}.Layout(gtx, func(gtx C) D {
						return p.nav.Layout(gtx, th)
					})
				}),
				// list
				layout.Flexed(1, func(gtx C) D {
					current := p.nav.current
					if current >= len(p.nav.btns) {
						current = 0
					}
					list := p.nav.btns[current].List
					if list == nil {
						return D{
							Size: gtx.Constraints.Max,
						}
					}

					inset := layout.Inset{
						Top:    5,
						Bottom: 5,
					}
					width := unit.Dp(800)
					if x := gtx.Metric.PxToDp(gtx.Constraints.Max.X); x > width {
						inset.Left = (x - width) / 2
						inset.Right = inset.Left
					}
					return inset.Layout(gtx, func(gtx C) D {
						return list.Layout(gtx, th)
					})
				}),
			)
		}),
		layout.Stacked(func(gtx C) D {
			return layout.Inset{
				Top:    10,
				Bottom: 10,
				Left:   10,
				Right:  10,
			}.Layout(gtx, func(gtx C) D {
				btn := material.IconButton(th, &p.btnCreate, icons.IconAdd, "Add")
				btn.Inset = layout.UniformInset(16)

				return btn.Layout(gtx)
			})
		}),
	)
}
