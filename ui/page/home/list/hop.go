package list

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"gioui.org/x/outlay"
	"github.com/go-gost/gostctl/api"
	"github.com/go-gost/gostctl/ui/page"
	"github.com/go-gost/gostctl/ui/theme"
)

type hopList struct {
	router *page.Router
	list   layout.List
	states []state
}

func Hop(r *page.Router) List {
	return &hopList{
		router: r,
		list: layout.List{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		},
		states: make([]state, 16),
	}
}

func (p *hopList) Layout(gtx page.C, th *page.T) page.D {
	cfg := api.GetConfig()
	hops := cfg.Hops

	if len(hops) > len(p.states) {
		states := p.states
		p.states = make([]state, len(hops))
		copy(p.states, states)
	}

	return p.list.Layout(gtx, len(hops), func(gtx page.C, index int) page.D {
		if p.states[index].clk.Clicked(gtx) {
			p.router.Goto(page.Route{
				Path: page.PageHop,
				ID:   hops[index].Name,
				Perm: page.PermReadWriteDelete,
			})
		}

		hop := hops[index]

		return layout.Inset{
			Top:    8,
			Bottom: 8,
			Left:   8,
			Right:  8,
		}.Layout(gtx, func(gtx page.C) page.D {
			return material.ButtonLayoutStyle{
				Background:   theme.Current().ListBg,
				CornerRadius: 12,
				Button:       &p.states[index].clk,
			}.Layout(gtx, func(gtx page.C) page.D {
				return layout.UniformInset(16).Layout(gtx, func(gtx page.C) page.D {
					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(gtx,
						// title
						layout.Rigid(func(gtx page.C) page.D {
							return layout.Flex{
								Alignment: layout.Middle,
								Spacing:   layout.SpaceBetween,
							}.Layout(gtx,
								layout.Flexed(1, func(gtx page.C) page.D {
									label := material.Body1(th, hop.Name)
									label.Font.Weight = font.SemiBold
									return label.Layout(gtx)
								}),
							)
						}),
						layout.Rigid(layout.Spacer{Height: 4}.Layout),
						layout.Rigid(func(gtx page.C) page.D {
							return layout.Flex{
								Alignment: layout.Middle,
							}.Layout(gtx,
								layout.Rigid(func(gtx page.C) page.D {
									return outlay.FlowWrap{
										Alignment: layout.Middle,
									}.Layout(gtx, len(hop.Nodes), func(gtx page.C, i int) page.D {
										return layout.Inset{
											Top:    4,
											Bottom: 4,
											Right:  8,
										}.Layout(gtx, func(gtx page.C) page.D {
											return component.SurfaceStyle{
												Theme:       th,
												ShadowStyle: component.ShadowStyle{CornerRadius: 14},
												Fill:        theme.Current().ItemBg,
											}.Layout(gtx, func(gtx page.C) page.D {
												return layout.Inset{
													Top:    4,
													Bottom: 4,
													Left:   10,
													Right:  10,
												}.Layout(gtx, material.Body2(th, hop.Nodes[i].Name).Layout)
											})
										})
									})
								}),
							)
						}),
					)
				})
			})
		})
	})
}
