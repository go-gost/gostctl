package mapping

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/go-gost/gostctl/ui/i18n"
)

type mappingDialog struct {
	hostname component.TextField
	ip       component.TextField
	alias    component.TextField

	OnClick   func(ok bool)
	btnCancel widget.Clickable
	btnOK     widget.Clickable
}

func (p *mappingDialog) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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
									layout.Flexed(1, material.H6(th, i18n.Mapping.Value()).Layout),
								)
							})
						}),

						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{
								Top:    16,
								Bottom: 16,
								Left:   24,
								Right:  24,
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{
									Axis: layout.Vertical,
								}.Layout(gtx,
									layout.Rigid(material.Body1(th, i18n.Hostname.Value()).Layout),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return p.hostname.Layout(gtx, th, "")
									}),
									layout.Rigid(layout.Spacer{Height: 8}.Layout),

									layout.Rigid(material.Body1(th, "IP").Layout),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return p.ip.Layout(gtx, th, "")
									}),
									layout.Rigid(layout.Spacer{Height: 8}.Layout),

									/*
										layout.Rigid(material.Body1(th, i18n.HostAlias.Value()).Layout),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											return p.alias.Layout(gtx, th, "")
										}),
									*/
								)
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
