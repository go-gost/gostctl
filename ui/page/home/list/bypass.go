package list

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/go-gost/gostctl/api"
	"github.com/go-gost/gostctl/ui/page"
	"github.com/go-gost/gostctl/ui/theme"
)

type bypassList struct {
	router *page.Router
	list   layout.List
	states []state
}

func Bypass(r *page.Router) List {
	return &bypassList{
		router: r,
		list: layout.List{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		},
		states: make([]state, 16),
	}
}

func (p *bypassList) Layout(gtx page.C, th *page.T) page.D {
	cfg := api.GetConfig()
	bypasses := cfg.Bypasses

	if len(bypasses) > len(p.states) {
		states := p.states
		p.states = make([]state, len(bypasses))
		copy(p.states, states)
	}

	return p.list.Layout(gtx, len(bypasses), func(gtx page.C, index int) page.D {
		if p.states[index].clk.Clicked(gtx) {
			p.router.Goto(page.Route{
				Path: page.PageBypass,
				ID:   bypasses[index].Name,
				Perm: page.PermReadWriteDelete,
			})
		}

		bypass := bypasses[index]

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
									label := material.Body1(th, bypass.Name)
									label.Font.Weight = font.SemiBold
									return label.Layout(gtx)
								}),
							)
						}),
					)
				})
			})
		})
	})
}
