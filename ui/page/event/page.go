package event

import (
	"slices"
	"time"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/go-gost/gostctl/ui/i18n"
	"github.com/go-gost/gostctl/ui/icons"
	"github.com/go-gost/gostctl/ui/page"
	"github.com/go-gost/gostctl/ui/theme"
)

type eventPage struct {
	router *page.Router
	list   widget.List

	btnBack widget.Clickable

	events []page.ServerEvent
}

func NewPage(r *page.Router) page.Page {
	p := &eventPage{
		router: r,

		list: widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
	}

	return p
}

func (p *eventPage) Init(opts ...page.PageOption) {
	var options page.PageOptions
	for _, opt := range opts {
		opt(&options)
	}

	p.events, _ = options.Value.([]page.ServerEvent)

	slices.Reverse(p.events)
}

func (p *eventPage) Layout(gtx page.C) page.D {
	th := p.router.Theme

	if p.btnBack.Clicked(gtx) {
		p.router.Back()
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		// header
		layout.Rigid(func(gtx page.C) page.D {
			return layout.UniformInset(8).Layout(gtx, func(gtx page.C) page.D {
				return layout.Flex{
					Spacing:   layout.SpaceBetween,
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(func(gtx page.C) page.D {
						btn := material.IconButton(th, &p.btnBack, icons.IconBack, "Back")
						btn.Color = th.Fg
						btn.Background = th.Bg
						return btn.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: 8}.Layout),
					layout.Rigid(func(gtx page.C) page.D {
						title := material.H6(th, i18n.Event.Value())
						return title.Layout(gtx)
					}),
					layout.Flexed(1, layout.Spacer{Width: 8}.Layout),
				)
			})
		}),
		layout.Flexed(1, func(gtx page.C) page.D {
			return layout.UniformInset(8).Layout(gtx, func(gtx page.C) page.D {
				return p.layout(gtx, th)
			})
		}),
	)
}

func (p *eventPage) layout(gtx page.C, th *page.T) page.D {
	return component.SurfaceStyle{
		Theme: th,
		ShadowStyle: component.ShadowStyle{
			CornerRadius: 12,
		},
		Fill: theme.Current().ContentSurfaceBg,
	}.Layout(gtx, func(gtx page.C) page.D {
		return layout.UniformInset(16).Layout(gtx, func(gtx page.C) page.D {
			return material.List(th, &p.list).Layout(gtx, len(p.events), func(gtx page.C, index int) page.D {
				return layout.Inset{
					Bottom: 8,
				}.Layout(gtx, func(gtx page.C) page.D {
					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(gtx,
						layout.Rigid(func(gtx page.C) page.D {
							label := material.Body1(th, p.events[index].Msg)
							label.Font.Weight = font.SemiBold
							return label.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Height: 4}.Layout),
						layout.Rigid(func(gtx page.C) page.D {
							return material.Body2(th, p.events[index].Time.Local().Format(time.RFC3339)).Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Height: 8}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							div := component.Divider(th)
							return div.Layout(gtx)
						}),
					)
				})
			})
		})
	})
}
