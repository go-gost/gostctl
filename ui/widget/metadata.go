package widget

import (
	"strings"
	"sync"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/go-gost/gostctl/ui/i18n"
	"github.com/go-gost/gostctl/ui/icons"
)

type MetadataDialog struct {
	OnClick   func(ok bool)
	list      material.ListStyle
	btnAdd    widget.Clickable
	btnCancel widget.Clickable
	btnOK     widget.Clickable
	metadata  []*kv
	once      sync.Once
}

func (p *MetadataDialog) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	p.once.Do(func() {
		p.list = material.List(th, &widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		})
	})

	if p.btnAdd.Clicked(gtx) {
		p.metadata = append(p.metadata, &kv{
			k: component.TextField{
				Editor: widget.Editor{
					SingleLine: true,
					MaxLen:     255,
				},
			},
			v: component.TextField{
				Editor: widget.Editor{
					SingleLine: true,
					MaxLen:     255,
				},
			},
		})
	}

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
										layout.Flexed(1, material.H6(th, "Metadata").Layout),
										layout.Rigid(layout.Spacer{Width: 8}.Layout),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											btn := material.IconButton(th, &p.btnAdd, icons.IconAdd, "Add")
											btn.Background = th.Bg
											btn.Color = th.Fg
											return btn.Layout(gtx)
										}),
									)
								})
							}),

							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								for i := range p.metadata {
									if p.metadata[i].remove.Clicked(gtx) {
										p.metadata = append(p.metadata[:i], p.metadata[i+1:]...)
										break
									}
								}

								gtx.Constraints.Max.Y -= gtx.Dp(80)
								return layout.Inset{
									Top:    8,
									Bottom: 8,
								}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return p.list.Layout(gtx, len(p.metadata), func(gtx layout.Context, index int) layout.Dimensions {
										return p.metadata[index].Layout(gtx, th)
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

func (p *MetadataDialog) Add(k, v string) {
	kv := &kv{}
	kv.k.SetText(k)
	kv.v.SetText(v)
	p.metadata = append(p.metadata, kv)
}

func (p *MetadataDialog) Clear() {
	p.metadata = nil
}

func (p *MetadataDialog) Metadata() []*kv {
	return p.metadata
}

type kv struct {
	k      component.TextField
	v      component.TextField
	remove widget.Clickable
}

func (p *kv) Get() (string, string) {
	return strings.TrimSpace(p.k.Text()), strings.TrimSpace(p.v.Text())
}

func (p *kv) Set(k, v string) {
	p.k.SetText(k)
	p.v.SetText(v)
}

func (p *kv) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Inset{
		Top:    8,
		Bottom: 8,
		Left:   24,
		Right:  24,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Alignment: layout.End,
		}.Layout(gtx,
			layout.Flexed(0.4, func(gtx layout.Context) layout.Dimensions {
				return p.k.Layout(gtx, th, "Key")
			}),
			layout.Rigid(layout.Spacer{Width: 8}.Layout),
			layout.Flexed(0.4, func(gtx layout.Context) layout.Dimensions {
				return p.v.Layout(gtx, th, "Value")
			}),
			layout.Rigid(layout.Spacer{Width: 8}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.IconButton(th, &p.remove, icons.IconDelete, "Remove")
				btn.Color = th.Fg
				btn.Background = th.Bg
				return btn.Layout(gtx)
			}),
		)
	})
}
