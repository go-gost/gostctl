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

type chainList struct {
	router *page.Router
	list   layout.List
	states []state
}

func Chain(r *page.Router) List {
	return &chainList{
		router: r,
		list: layout.List{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		},
		states: make([]state, 16),
	}
}

func (p *chainList) Layout(gtx page.C, th *page.T) page.D {
	cfg := api.GetConfig()
	chains := cfg.Chains

	if len(chains) > len(p.states) {
		states := p.states
		p.states = make([]state, len(chains))
		copy(p.states, states)
	}

	return p.list.Layout(gtx, len(chains), func(gtx page.C, index int) page.D {
		if p.states[index].clk.Clicked(gtx) {
			p.router.Goto(page.Route{
				Path: page.PageChain,
				ID:   chains[index].Name,
				Perm: page.PermReadWriteDelete,
			})
		}

		chain := chains[index]

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
						// title, state
						layout.Rigid(func(gtx page.C) page.D {
							return layout.Flex{
								Alignment: layout.Middle,
								Spacing:   layout.SpaceBetween,
							}.Layout(gtx,
								layout.Flexed(1, func(gtx page.C) page.D {
									label := material.Body1(th, chain.Name)
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
									}.Layout(gtx, len(chain.Hops), func(gtx page.C, i int) page.D {
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
												}.Layout(gtx, material.Body2(th, chain.Hops[i].Name).Layout)
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
