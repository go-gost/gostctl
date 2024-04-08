package home

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/go-gost/gostctl/ui/i18n"
	"github.com/go-gost/gostctl/ui/icons"
	"github.com/go-gost/gostctl/ui/page"
	"github.com/go-gost/gostctl/ui/page/home/list"
	ui_widget "github.com/go-gost/gostctl/ui/widget"
)

type C = layout.Context
type D = layout.Dimensions

type navPage struct {
	list list.List
	path page.PagePath
}

type homePage struct {
	router      *page.Router
	nav         *ui_widget.Nav
	pages       []navPage
	btnAdd      widget.Clickable
	btnSettings widget.Clickable
}

func NewPage(r *page.Router) page.Page {
	return &homePage{
		router: r,
		nav: ui_widget.NewNav(
			ui_widget.NewNavButton(i18n.Server),
			ui_widget.NewNavButton(i18n.Service),
			ui_widget.NewNavButton(i18n.Chain),
			ui_widget.NewNavButton(i18n.Hop),
			ui_widget.NewNavButton(i18n.Auther),
			ui_widget.NewNavButton(i18n.Admission),
			ui_widget.NewNavButton(i18n.Bypass),
			ui_widget.NewNavButton(i18n.Resolver),
			ui_widget.NewNavButton(i18n.Hosts),
			ui_widget.NewNavButton(i18n.Limiter),
			ui_widget.NewNavButton(i18n.Ingress),
			ui_widget.NewNavButton(i18n.Logger),
		),
		pages: []navPage{
			{
				list: list.Server(r),
				path: page.PageServerEdit,
			},
			{
				list: list.Service(r),
				path: page.PageServiceEdit,
			},
		},
	}
}

func (p *homePage) Init(opts ...page.PageOption) {
}

func (p *homePage) Layout(gtx C) D {
	if p.btnAdd.Clicked(gtx) {
		if current := p.nav.Current(); current < len(p.pages) {
			p.router.Goto(page.Route{Path: p.pages[current].path})
		}
	}

	th := p.router.Theme

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
						Top:    8,
						Bottom: 8,
						Left:   8,
						Right:  8,
					}.Layout(gtx, func(gtx C) D {
						return layout.Flex{
							Spacing:   layout.SpaceBetween,
							Alignment: layout.Middle,
						}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								gtx.Constraints.Max.X = gtx.Dp(50)
								return icons.IconApp.Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Width: 8}.Layout),
							layout.Rigid(func(gtx C) D {
								label := material.H6(th, "GOST")
								label.Font.Weight = font.Bold
								return label.Layout(gtx)
							}),
							layout.Flexed(1, layout.Spacer{Width: 8}.Layout),
							layout.Rigid(func(gtx C) D {
								if p.btnSettings.Clicked(gtx) {
									p.router.Goto(page.Route{
										Path: page.PageSettings,
									})
								}

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
						Top:    4,
						Bottom: 4,
					}.Layout(gtx, func(gtx C) D {
						return p.nav.Layout(gtx, th)
					})
				}),
				// list
				layout.Flexed(1, func(gtx C) D {
					current := p.nav.Current()
					if current >= len(p.pages) {
						current = 0
					}
					pg := p.pages[current]
					if pg.list == nil {
						return D{
							Size: gtx.Constraints.Max,
						}
					}

					return pg.list.Layout(gtx, th)
				}),
			)
		}),
		layout.Stacked(func(gtx C) D {
			return layout.Inset{
				Top:    16,
				Bottom: 16,
				Left:   16,
				Right:  16,
			}.Layout(gtx, func(gtx C) D {
				btn := material.IconButton(th, &p.btnAdd, icons.IconAdd, "Add")
				btn.Inset = layout.UniformInset(16)

				return btn.Layout(gtx)
			})
		}),
	)
}
