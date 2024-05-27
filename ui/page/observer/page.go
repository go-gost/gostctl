package observer

import (
	"context"
	"strings"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/go-gost/gostctl/api"
	"github.com/go-gost/gostctl/api/runner"
	"github.com/go-gost/gostctl/api/runner/task"
	"github.com/go-gost/gostctl/api/util"
	"github.com/go-gost/gostctl/ui/i18n"
	"github.com/go-gost/gostctl/ui/icons"
	"github.com/go-gost/gostctl/ui/page"
	"github.com/go-gost/gostctl/ui/theme"
	ui_widget "github.com/go-gost/gostctl/ui/widget"
)

type observerPage struct {
	router *page.Router

	menu ui_widget.Menu
	mode widget.Enum
	list layout.List

	btnBack   widget.Clickable
	btnDelete widget.Clickable
	btnEdit   widget.Clickable
	btnSave   widget.Clickable

	name component.TextField

	pluginType          ui_widget.Selector
	pluginAddr          component.TextField
	pluginEnableTLS     ui_widget.Switcher
	pluginTLSSecure     ui_widget.Switcher
	pluginTLSServerName component.TextField

	id   string
	perm page.Perm

	edit   bool
	create bool

	delDialog ui_widget.Dialog
}

func NewPage(r *page.Router) page.Page {
	p := &observerPage{
		router: r,

		list: layout.List{
			// NOTE: the list must be vertical
			Axis: layout.Vertical,
		},
		name: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     128,
			},
		},

		pluginType: ui_widget.Selector{Title: i18n.Type},
		pluginAddr: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     128,
			},
		},
		pluginEnableTLS: ui_widget.Switcher{Title: i18n.TLS},
		pluginTLSSecure: ui_widget.Switcher{Title: i18n.VerifyServerCert},
		pluginTLSServerName: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     255,
			},
		},

		delDialog: ui_widget.Dialog{Title: i18n.DeleteObserver},
	}

	return p
}

func (p *observerPage) Init(opts ...page.PageOption) {
	var options page.PageOptions
	for _, opt := range opts {
		opt(&options)
	}
	p.id = options.ID

	if p.id != "" {
		p.edit = false
		p.create = false
		p.name.ReadOnly = true
	} else {
		p.edit = true
		p.create = true
		p.name.ReadOnly = false
	}

	p.perm = options.Perm

	observer, _ := options.Value.(*api.ObserverConfig)

	if observer == nil {
		cfg := api.GetConfig()
		for _, v := range cfg.Observers {
			if v.Name == p.id {
				observer = v
				break
			}
		}
		if observer == nil {
			observer = &api.ObserverConfig{}
		}
	}

	p.mode.Value = string(page.PluginMode)
	p.name.SetText(observer.Name)

	{
		p.pluginType.Clear()
		p.pluginAddr.Clear()
		p.pluginEnableTLS.SetValue(false)
		p.pluginTLSSecure.SetValue(false)
		p.pluginTLSServerName.Clear()

		if observer.Plugin != nil {
			for i := range page.PluginTypeOptions {
				if page.PluginTypeOptions[i].Value == observer.Plugin.Type {
					p.pluginType.Select(ui_widget.SelectorItem{Key: page.PluginTypeOptions[i].Key, Value: page.PluginTypeOptions[i].Value})
					break
				}
			}
			p.pluginAddr.SetText(observer.Plugin.Addr)

			if observer.Plugin.TLS != nil {
				p.pluginEnableTLS.SetValue(true)
				p.pluginTLSSecure.SetValue(observer.Plugin.TLS.Secure)
				p.pluginTLSServerName.SetText(observer.Plugin.TLS.ServerName)
			}
		}
	}
}

func (p *observerPage) Layout(gtx page.C) page.D {
	if p.btnBack.Clicked(gtx) {
		p.router.Back()
	}
	if p.btnEdit.Clicked(gtx) {
		p.edit = true
	}
	if p.btnSave.Clicked(gtx) {
		if p.save() {
			p.router.Back()
		}
	}

	if p.btnDelete.Clicked(gtx) {
		p.delDialog.OnClick = func(ok bool) {
			if ok {
				p.delete()
				p.router.Back()
			}
			p.router.HideModal(gtx)
		}
		p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
			return p.delDialog.Layout(gtx, th)
		})
	}

	th := p.router.Theme

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
					layout.Flexed(1, func(gtx page.C) page.D {
						title := material.H6(th, i18n.Observer.Value())
						return title.Layout(gtx)
					}),
					layout.Rigid(func(gtx page.C) page.D {
						if p.perm&page.PermDelete == 0 || p.create {
							return page.D{}
						}

						btn := material.IconButton(th, &p.btnDelete, icons.IconDelete, "Delete")
						btn.Color = th.Fg
						btn.Background = th.Bg
						return btn.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: 8}.Layout),
					layout.Rigid(func(gtx page.C) page.D {
						if p.perm&page.PermWrite == 0 {
							return page.D{}
						}

						if p.edit {
							btn := material.IconButton(th, &p.btnSave, icons.IconDone, "Done")
							btn.Color = th.Fg
							btn.Background = th.Bg
							return btn.Layout(gtx)
						} else {
							btn := material.IconButton(th, &p.btnEdit, icons.IconEdit, "Edit")
							btn.Color = th.Fg
							btn.Background = th.Bg
							return btn.Layout(gtx)
						}
					}),
				)
			})
		}),
		layout.Flexed(1, func(gtx page.C) page.D {
			return p.list.Layout(gtx, 1, func(gtx page.C, index int) page.D {
				return layout.UniformInset(8).Layout(gtx, func(gtx page.C) page.D {
					return p.layout(gtx, th)
				})
			})
		}),
	)
}

func (p *observerPage) layout(gtx page.C, th *page.T) page.D {
	src := gtx.Source

	if !p.edit {
		gtx = gtx.Disabled()
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
					gtx.Source = src

					return layout.Flex{
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Rigid(func(gtx page.C) page.D {
							return material.RadioButton(th, &p.mode, string(page.PluginMode), i18n.Plugin.Value()).Layout(gtx)
						}),
					)
				}),

				layout.Rigid(layout.Spacer{Height: 16}.Layout),

				layout.Rigid(material.Body1(th, i18n.Name.Value()).Layout),
				layout.Rigid(func(gtx page.C) page.D {
					return p.name.Layout(gtx, th, "")
				}),
				layout.Rigid(layout.Spacer{Height: 8}.Layout),

				// plugin
				layout.Rigid(func(gtx page.C) page.D {
					if p.mode.Value != string(page.PluginMode) {
						return page.D{}
					}

					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(gtx,
						layout.Rigid(material.Body1(th, i18n.Address.Value()).Layout),
						layout.Rigid(func(gtx page.C) page.D {
							return p.pluginAddr.Layout(gtx, th, "")
						}),
						layout.Rigid(layout.Spacer{Height: 4}.Layout),
						layout.Rigid(func(gtx page.C) page.D {
							if p.pluginType.Clicked(gtx) {
								p.showPluginTypeMenu(gtx)
							}
							return p.pluginType.Layout(gtx, th)
						}),
						layout.Rigid(func(gtx page.C) page.D {
							return p.pluginEnableTLS.Layout(gtx, th)
						}),
						layout.Rigid(func(gtx page.C) page.D {
							if !p.pluginEnableTLS.Value() {
								return page.D{}
							}

							return layout.Flex{
								Axis: layout.Vertical,
							}.Layout(gtx,
								layout.Rigid(func(gtx page.C) page.D {
									return layout.UniformInset(8).Layout(gtx, func(gtx page.C) page.D {
										return layout.Flex{
											Axis: layout.Vertical,
										}.Layout(gtx,
											layout.Rigid(func(gtx page.C) page.D {
												return p.pluginTLSSecure.Layout(gtx, th)
											}),

											layout.Rigid(func(gtx page.C) page.D {
												return material.Body1(th, i18n.ServerName.Value()).Layout(gtx)
											}),
											layout.Rigid(func(gtx page.C) page.D {
												return p.pluginTLSServerName.Layout(gtx, th, "")
											}),
										)
									})
								}),
							)
						}),
					)
				}),
			)
		})
	})
}

func (p *observerPage) showPluginTypeMenu(gtx page.C) {
	for i := range page.PluginTypeOptions {
		page.PluginTypeOptions[i].Selected = p.pluginType.AnyValue(page.PluginTypeOptions[i].Value)
	}

	p.menu.Title = i18n.Type
	p.menu.Options = page.PluginTypeOptions
	p.menu.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		p.pluginType.Clear()
		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				p.pluginType.Select(ui_widget.SelectorItem{Key: p.menu.Options[i].Key, Value: p.menu.Options[i].Value})
			}
		}
	}
	p.menu.OnAdd = nil
	p.menu.Multiple = false

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *observerPage) save() bool {
	cfg := p.generateConfig()

	var err error
	if p.id == "" {
		err = runner.Exec(context.Background(),
			task.CreateObserver(cfg),
			runner.WithCancel(true),
		)
	} else {
		err = runner.Exec(context.Background(),
			task.UpdateObserver(cfg),
			runner.WithCancel(true),
		)
	}
	util.RestartGetConfigTask()

	return err == nil
}

func (p *observerPage) generateConfig() *api.ObserverConfig {
	cfg := &api.ObserverConfig{
		Name: strings.TrimSpace(p.name.Text()),
	}
	if p.mode.Value == string(page.PluginMode) && p.pluginType.Value() != "" {
		cfg.Plugin = &api.PluginConfig{
			Type: p.pluginType.Value(),
			Addr: p.pluginAddr.Text(),
		}
		if p.pluginEnableTLS.Value() {
			cfg.Plugin.TLS = &api.TLSConfig{
				Secure:     p.pluginTLSSecure.Value(),
				ServerName: strings.TrimSpace(p.pluginTLSServerName.Text()),
			}
		}
		return cfg
	}

	return cfg
}

func (p *observerPage) delete() {
	runner.Exec(context.Background(),
		task.DeleteObserver(p.id),
		runner.WithCancel(true),
	)
	util.RestartGetConfigTask()
}
