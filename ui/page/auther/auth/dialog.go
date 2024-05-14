package auth

import (
	"strings"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/go-gost/gostctl/ui/i18n"
)

type authDialog struct {
	kv        kv
	OnClick   func(ok bool)
	btnCancel widget.Clickable
	btnOK     widget.Clickable
}

func (p *authDialog) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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
									layout.Flexed(1, material.H6(th, i18n.Auth.Value()).Layout),
								)
							})
						}),

						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{
								Top:    8,
								Bottom: 8,
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return p.kv.Layout(gtx, th)
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
											CornerRadius: 18,
											Button:       &p.btnCancel,
										}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											return layout.Inset{
												Top:    8,
												Bottom: 8,
												Left:   20,
												Right:  20,
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
											CornerRadius: 18,
											Button:       &p.btnOK,
										}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											return layout.Inset{
												Top:    8,
												Bottom: 8,
												Left:   20,
												Right:  20,
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
}

type kv struct {
	k component.TextField
	v component.TextField
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
			Axis: layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(material.Body1(th, i18n.Username.Value()).Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return p.k.Layout(gtx, th, "")
			}),
			layout.Rigid(layout.Spacer{Height: 8}.Layout),

			layout.Rigid(material.Body1(th, i18n.Password.Value()).Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return p.v.Layout(gtx, th, "")
			}),
		)
	})
}
