package page

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/go-gost/gui/ui/icons"
)

type Menu struct {
	Theme    *material.Theme
	Surface  component.SurfaceStyle
	List     layout.List
	Items    []MenuItem
	Title    string
	Selected func(index int)
	btnAdd   widget.Clickable
	ShowAdd  bool
	Multiple bool
}

func (p *Menu) Layout(gtx C, th *material.Theme) D {
	if gtx.Constraints.Max.X > gtx.Dp(800) {
		gtx.Constraints.Max.X = gtx.Dp(800)
	}
	gtx.Constraints.Max.X = gtx.Constraints.Max.X * 2 / 3
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return p.Surface.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    10,
				Bottom: 10,
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{
							Left:   10,
							Right:  10,
							Bottom: 20,
						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{
								Alignment: layout.Middle,
							}.Layout(gtx,
								layout.Flexed(1, material.Body1(th, p.Title).Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if !p.ShowAdd {
										return D{}
									}
									btn := material.IconButton(th, &p.btnAdd, icons.IconAdd, "Add")
									btn.Background = th.Bg
									btn.Color = th.Fg
									btn.Inset = layout.UniformInset(0)
									return btn.Layout(gtx)
								}),
							)
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return p.List.Layout(gtx, len(p.Items), func(gtx layout.Context, index int) layout.Dimensions {
							if p.Items[index].state.Clicked(gtx) {
								if p.Multiple {
									p.Items[index].Selected = !p.Items[index].Selected
								} else {
									for i := range p.Items {
										p.Items[i].Selected = false
									}
									p.Items[index].Selected = true
								}
								if p.Selected != nil {
									p.Selected(index)
								}
							}
							return p.Items[index].Layout(gtx, th)
						})
					}),
				)
			})
		})
	})
}

type MenuItem struct {
	state    widget.Clickable
	Key      string
	Value    string
	Selected bool
}

func (p *MenuItem) Layout(gtx C, th *material.Theme) D {
	return material.ButtonLayoutStyle{
		Background: th.Bg,
		Button:     &p.state,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top:    10,
			Bottom: 10,
			Left:   10,
			Right:  10,
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Spacing:   layout.SpaceBetween,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Flexed(1, material.Body2(th, p.Key).Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if p.Selected {
						gtx.Constraints.Min.X = gtx.Dp(16)
						return icons.IconDone.Layout(gtx, th.Fg)
					}
					return D{}
				}),
			)
		})
	})
}
