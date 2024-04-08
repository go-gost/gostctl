package widget

import (
	"sync"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/go-gost/gostctl/ui/i18n"
	"github.com/go-gost/gostctl/ui/icons"
)

type Menu struct {
	Title     i18n.Key
	Options   []MenuOption
	OnClick   func(ok bool)
	list      material.ListStyle
	btnAdd    widget.Clickable
	ShowAdd   bool
	Multiple  bool
	btnCancel widget.Clickable
	btnOK     widget.Clickable
	once      sync.Once
}

func (p *Menu) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	p.once.Do(func() {
		p.list = material.List(th, &widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		})
	})

	var cl widget.Clickable
	return cl.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		if gtx.Constraints.Max.X > gtx.Dp(800) {
			gtx.Constraints.Max.X = gtx.Dp(800)
		}
		gtx.Constraints.Max.X = gtx.Constraints.Max.X * 2 / 3
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    16,
				Bottom: 16,
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return component.SurfaceStyle{
					Theme: th,
					ShadowStyle: component.ShadowStyle{
						CornerRadius: 28,
					},
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Top:    16,
						Bottom: 16,
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis: layout.Vertical,
						}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{
									Top:    8,
									Bottom: 8,
									Left:   24,
									Right:  24,
								}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return layout.Flex{
										Alignment: layout.Middle,
									}.Layout(gtx,
										layout.Flexed(1, material.H6(th, p.Title.Value()).Layout),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											if !p.ShowAdd {
												return layout.Dimensions{}
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
								gtx.Constraints.Max.Y -= gtx.Dp(80)

								return layout.Inset{
									Top:    8,
									Bottom: 8,
								}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return p.list.Layout(gtx, len(p.Options), func(gtx layout.Context, index int) layout.Dimensions {
										if p.Options[index].state.Clicked(gtx) {
											if p.Multiple {
												p.Options[index].Selected = !p.Options[index].Selected
											} else {
												for i := range p.Options {
													p.Options[i].Selected = false
												}
												p.Options[index].Selected = true
											}
										}

										return p.Options[index].Layout(gtx, th)
									})
								})
							}),

							layout.Rigid(layout.Spacer{Height: 8}.Layout),

							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{
									Left:  24,
									Right: 24,
								}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return layout.Flex{
										Spacing:   layout.SpaceBetween,
										Alignment: layout.Middle,
									}.Layout(gtx,
										layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
											return layout.Spacer{Width: 8}.Layout(gtx)
										}),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											if p.btnCancel.Clicked(gtx) && p.OnClick != nil {
												p.OnClick(false)
											}

											return material.ButtonLayoutStyle{
												Background:   th.Bg,
												CornerRadius: 20,
												Button:       &p.btnCancel,
											}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
												return layout.Inset{
													Top:    8,
													Bottom: 8,
													Left:   24,
													Right:  24,
												}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
													label := material.Body1(th, i18n.Cancel.Value())
													label.Color = th.Fg
													return label.Layout(gtx)
												})

											})
										}),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											return layout.Spacer{Width: 8}.Layout(gtx)
										}),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											if p.btnOK.Clicked(gtx) && p.OnClick != nil {
												p.OnClick(true)
											}

											return material.ButtonLayoutStyle{
												Background:   th.Bg,
												CornerRadius: 20,
												Button:       &p.btnOK,
											}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
												return layout.Inset{
													Top:    8,
													Bottom: 8,
													Left:   24,
													Right:  24,
												}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
													label := material.Body1(th, i18n.OK.Value())
													label.Color = th.Fg
													return label.Layout(gtx)
												})

											})
										}),
									)
								})
							}),
						)
					})
				})
			})
		})
	})
}

type MenuOption struct {
	state    widget.Clickable
	Key      i18n.Key
	Value    string
	Selected bool
}

func (p *MenuOption) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return material.ButtonLayoutStyle{
		Background: th.Bg,
		Button:     &p.state,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top:    8,
			Bottom: 8,
			Left:   24,
			Right:  24,
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Spacing:   layout.SpaceBetween,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					key := p.Key.Value()
					if key == "" {
						key = p.Value
					}
					return material.Body2(th, key).Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Width: 8}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if p.Selected {
						gtx.Constraints.Min.X = gtx.Dp(16)
						return icons.IconDone.Layout(gtx, th.Fg)
					}
					return layout.Dimensions{}
				}),
			)
		})
	})
}