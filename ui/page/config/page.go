package config

import (
	"bytes"
	"encoding/json"
	"io"

	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/go-gost/gostctl/ui/i18n"
	"github.com/go-gost/gostctl/ui/icons"
	"github.com/go-gost/gostctl/ui/page"
	"github.com/go-gost/gostctl/ui/theme"
	"gopkg.in/yaml.v3"
)

const (
	FormatYAML = "yaml"
	FormatJSON = "json"
)

type configPage struct {
	router *page.Router
	list   widget.List
	editor widget.Editor

	btnBack widget.Clickable
	btnCopy widget.Clickable

	format widget.Enum

	cfg any
}

func NewPage(r *page.Router) page.Page {
	p := &configPage{
		router: r,

		list: widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		editor: widget.Editor{
			ReadOnly:   true,
			WrapPolicy: text.WrapGraphemes,
		},
	}

	return p
}

func (p *configPage) Init(opts ...page.PageOption) {
	var options page.PageOptions
	for _, opt := range opts {
		opt(&options)
	}

	p.cfg = options.Value
	p.format.Value = FormatYAML

	value, _ := yaml.Marshal(p.cfg)

	p.editor.SetText(string(value))
}

func (p *configPage) Layout(gtx page.C) page.D {
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
						title := material.H6(th, i18n.Config.Value())
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

func (p *configPage) layout(gtx page.C, th *page.T) page.D {
	if p.btnCopy.Clicked(gtx) {
		gtx.Execute(clipboard.WriteCmd{
			Data: io.NopCloser(bytes.NewBufferString(p.editor.Text())),
		})
	}
	if p.format.Update(gtx) {
		switch p.format.Value {
		case FormatJSON:
			value, _ := json.Marshal(p.cfg)
			var out bytes.Buffer
			json.Indent(&out, value, "", "  ")
			p.editor.SetText(out.String())
		default:
			value, _ := yaml.Marshal(p.cfg)
			p.editor.SetText(string(value))
		}
	}

	return component.SurfaceStyle{
		Theme: th,
		ShadowStyle: component.ShadowStyle{
			CornerRadius: 12,
		},
		Fill: theme.Current().ContentSurfaceBg,
	}.Layout(gtx, func(gtx page.C) page.D {
		return layout.UniformInset(16).Layout(gtx, func(gtx page.C) page.D {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx page.C) page.D {
					return layout.Flex{
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Rigid(func(gtx page.C) page.D {
							return material.RadioButton(th, &p.format, FormatYAML, FormatYAML).Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: 4}.Layout),
						layout.Rigid(func(gtx page.C) page.D {
							return material.RadioButton(th, &p.format, FormatJSON, FormatJSON).Layout(gtx)
						}),
						layout.Flexed(1, layout.Spacer{Width: 4}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							btn := material.IconButton(th, &p.btnCopy, icons.IconCopy, "Copy")
							btn.Color = th.Fg
							btn.Background = theme.Current().ContentSurfaceBg
							return btn.Layout(gtx)
						}),
					)
				}),

				layout.Rigid(layout.Spacer{Height: 16}.Layout),

				layout.Flexed(1, func(gtx page.C) page.D {
					return widget.Border{
						CornerRadius: 5,
						Color:        th.ContrastBg,
						Width:        1,
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{
							Top:    8,
							Bottom: 8,
							Left:   8,
						}.Layout(gtx, func(gtx page.C) page.D {
							return material.List(th, &p.list).Layout(gtx, 1, func(gtx page.C, _ int) page.D {
								return layout.Flex{}.Layout(gtx,
									layout.Flexed(1, func(gtx page.C) page.D {
										return layout.UniformInset(4).Layout(gtx, func(gtx page.C) page.D {
											/*
												editor := material.Editor(th, &p.editor, "")
												editor.TextSize = th.TextSize
												// editor.Font.Weight = font.Medium
												editor.Font.Typeface = "monospace"
												return editor.Layout(gtx)
											*/
											label := material.Label(th, th.TextSize, p.editor.Text())
											return label.Layout(gtx)
										})
									}),
								)
							})
						})
					})
				}),
			)
		})
	})
}
