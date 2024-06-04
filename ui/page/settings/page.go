package settings

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"gioui.org/x/pref/locale"
	gio_theme "gioui.org/x/pref/theme"
	"github.com/go-gost/gostctl/config"
	"github.com/go-gost/gostctl/ui/i18n"
	"github.com/go-gost/gostctl/ui/icons"
	"github.com/go-gost/gostctl/ui/page"
	"github.com/go-gost/gostctl/ui/theme"
	ui_widget "github.com/go-gost/gostctl/ui/widget"
	"github.com/go-gost/gostctl/version"
	"golang.org/x/text/language"
)

type settingsPage struct {
	router *page.Router
	menu   ui_widget.Menu
	list   widget.List

	btnBack widget.Clickable

	lang  ui_widget.Selector
	theme ui_widget.Selector
}

func NewPage(r *page.Router) page.Page {
	return &settingsPage{
		router: r,
		list: widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		lang:  ui_widget.Selector{Title: i18n.Language},
		theme: ui_widget.Selector{Title: i18n.Theme},
	}
}

func (p *settingsPage) Init(opts ...page.PageOption) {
	settings := config.Get().Settings

	p.lang.Clear()
	for _, lang := range i18n.Langs() {
		if lang.Value == settings.Lang {
			p.lang.Select(ui_widget.SelectorItem{
				Name:  lang.Name,
				Value: lang.Value,
			})
			break
		}
	}

	p.theme.Clear()
	switch settings.Theme {
	case theme.Light:
		p.theme.Select(ui_widget.SelectorItem{
			Key:   i18n.Light,
			Value: settings.Theme,
		})
	case theme.Dark:
		p.theme.Select(ui_widget.SelectorItem{
			Key:   i18n.Dark,
			Value: settings.Theme,
		})
	}
}

func (p *settingsPage) Layout(gtx layout.Context) layout.Dimensions {
	if p.btnBack.Clicked(gtx) {
		p.router.Back()
	}

	th := p.router.Theme

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		// header
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    8,
				Bottom: 8,
				Left:   8,
				Right:  8,
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.IconButton(th, &p.btnBack, icons.IconBack, "Back")
						btn.Color = th.Fg
						btn.Background = th.Bg
						return btn.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: 8}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						title := material.H6(th, i18n.Settings.Value())
						return title.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: 8}.Layout),
				)
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return p.list.Layout(gtx, 1, func(gtx layout.Context, _ int) layout.Dimensions {
				return layout.Inset{
					Top:    8,
					Bottom: 8,
					Left:   8,
					Right:  8,
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return p.layout(gtx, th)
				})
			})
		}),
	)
}

func (p *settingsPage) layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Vertical,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max.X = gtx.Dp(80)
						return icons.IconApp.Layout(gtx)
					})
				}),
				layout.Rigid(layout.Spacer{Height: 16}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					label := material.Body1(th, "GOST")
					label.Font.Weight = font.SemiBold
					return layout.Center.Layout(gtx, label.Layout)
				}),
				layout.Rigid(layout.Spacer{Height: 8}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, material.Body1(th, version.Version).Layout)
				}),
			)
		}),
		layout.Rigid(layout.Spacer{Height: 32}.Layout),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return component.SurfaceStyle{
				Theme: th,
				ShadowStyle: component.ShadowStyle{
					CornerRadius: 12,
				},
				Fill: theme.Current().ContentSurfaceBg,
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(16).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if p.lang.Clicked(gtx) {
								p.showLangMenu(gtx)
							}
							return p.lang.Layout(gtx, th)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if p.theme.Clicked(gtx) {
								p.showThemeMenu(gtx)
							}
							return p.theme.Layout(gtx, th)
						}),
					)
				})
			})
		}),
	)
}

func (p *settingsPage) showLangMenu(gtx layout.Context) {
	options := []ui_widget.MenuOption{}
	for _, lang := range i18n.Langs() {
		options = append(options, ui_widget.MenuOption{
			Name:  lang.Name,
			Value: lang.Value,
		})
	}

	for i := range options {
		if found := p.lang.AnyValue(options[i].Value); found {
			options[i].Selected = found
			break
		}
	}

	p.menu.Title = i18n.Language
	p.menu.Options = options
	p.menu.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		p.lang.Clear()

		for index := range p.menu.Options {
			if p.menu.Options[index].Selected {
				p.lang.Select(ui_widget.SelectorItem{
					Key:   p.menu.Options[index].Key,
					Name:  p.menu.Options[index].Name,
					Value: p.menu.Options[index].Value,
				})
				break
			}
		}

		cfg := config.Get()
		cfg.Settings.Lang = p.lang.Item().Value

		config.Set(cfg)
		cfg.Write()

		i18n.Set(cfg.Settings.Lang)
		if cfg.Settings.Lang == "" {
			lang, err := locale.Language()
			if err != nil {
				lang = language.English
			}
			if ls := lang.String(); ls != i18n.Current().Value {
				i18n.Set(ls)
			}
		}
	}

	p.router.ShowModal(gtx, func(gtx page.C, th *material.Theme) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *settingsPage) showThemeMenu(gtx layout.Context) {
	options := []ui_widget.MenuOption{
		{Key: i18n.Light, Value: theme.Light},
		{Key: i18n.Dark, Value: theme.Dark},
	}

	for i := range options {
		if found := p.theme.AnyValue(options[i].Value); found {
			options[i].Selected = found
			break
		}
	}

	p.menu.Title = i18n.Theme
	p.menu.Options = options
	p.menu.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		p.theme.Clear()

		for index := range p.menu.Options {
			if p.menu.Options[index].Selected {
				p.theme.Select(ui_widget.SelectorItem{
					Key:   p.menu.Options[index].Key,
					Value: p.menu.Options[index].Value,
				})
				break
			}
		}

		cfg := config.Get()
		cfg.Settings.Theme = p.theme.Item().Value

		config.Set(cfg)
		cfg.Write()

		theme.Use(cfg.Settings.Theme)
		if cfg.Settings.Theme == "" {
			if dark, _ := gio_theme.IsDarkMode(); dark {
				if theme.Current().Name == theme.Light {
					theme.Use(theme.Dark)
				}
			} else {
				if theme.Current().Name == theme.Dark {
					theme.Use(theme.Light)
				}
			}
		}

		p.router.Emit(page.Event{ID: page.EventThemeChanged})
	}

	p.router.ShowModal(gtx, func(gtx page.C, th *material.Theme) page.D {
		return p.menu.Layout(gtx, th)
	})
}
