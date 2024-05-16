package resolver

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
	"github.com/google/uuid"
)

type nameserver struct {
	id     string
	cfg    *api.NameserverConfig
	clk    widget.Clickable
	delete widget.Clickable
}

type resolverPage struct {
	router *page.Router

	menu ui_widget.Menu
	mode widget.Enum
	list layout.List

	btnBack   widget.Clickable
	btnDelete widget.Clickable
	btnEdit   widget.Clickable
	btnSave   widget.Clickable

	name component.TextField

	addNameserver widget.Clickable
	nameservers   []nameserver

	reload component.TextField

	enableFileDataSource ui_widget.Switcher
	filePath             component.TextField

	enableRedisDataSource ui_widget.Switcher
	redisAddr             component.TextField
	redisDB               component.TextField
	redisPassword         component.TextField
	redisKey              component.TextField

	enableHTTPDataSource ui_widget.Switcher
	httpURL              component.TextField
	httpTimeout          component.TextField

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
	p := &resolverPage{
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

		reload: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     16,
				Filter:     "1234567890",
			},
			Suffix: func(gtx page.C) page.D {
				return material.Body1(r.Theme, i18n.TimeSecond.Value()).Layout(gtx)
			},
		},
		enableFileDataSource: ui_widget.Switcher{Title: i18n.FileDataSource},
		filePath: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     128,
			},
		},

		enableRedisDataSource: ui_widget.Switcher{Title: i18n.RedisDataSource},
		redisAddr: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     128,
			},
		},
		redisDB: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     16,
				Filter:     "1234567890",
			},
		},
		redisPassword: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     255,
			},
		},
		redisKey: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     255,
			},
		},

		enableHTTPDataSource: ui_widget.Switcher{Title: i18n.HTTPDataSource},
		httpURL: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     255,
			},
		},
		httpTimeout: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     16,
				Filter:     "1234567890",
			},
			Suffix: func(gtx page.C) page.D {
				return material.Body1(r.Theme, i18n.TimeSecond.Value()).Layout(gtx)
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

		delDialog: ui_widget.Dialog{Title: i18n.DeleteResolver},
	}

	return p
}

func (p *resolverPage) Init(opts ...page.PageOption) {
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

	resolver, _ := options.Value.(*api.ResolverConfig)

	if resolver == nil {
		cfg := api.GetConfig()
		for _, v := range cfg.Resolvers {
			if v.Name == p.id {
				resolver = v
				break
			}
		}
		if resolver == nil {
			resolver = &api.ResolverConfig{}
		}
	}

	p.mode.Value = string(page.BasicMode)
	p.name.SetText(resolver.Name)

	{
		p.nameservers = nil
		for _, v := range resolver.Nameservers {
			p.nameservers = append(p.nameservers, nameserver{
				id:  uuid.New().String(),
				cfg: v,
			})
		}
	}

	/*
		{
			p.reload.SetText(strconv.Itoa(int(resolver.Reload.Seconds())))

			p.enableFileDataSource.SetValue(false)
			if resolver.File != nil {
				p.enableFileDataSource.SetValue(true)
				p.filePath.SetText(resolver.File.Path)
			}

			p.enableRedisDataSource.SetValue(false)
			if resolver.Redis != nil {
				p.enableRedisDataSource.SetValue(true)
				p.redisAddr.SetText(resolver.Redis.Addr)
				p.redisDB.SetText(strconv.Itoa(resolver.Redis.DB))
				p.redisPassword.SetText(resolver.Redis.Password)
				p.redisKey.SetText(resolver.Redis.Key)
			}

			p.enableHTTPDataSource.SetValue(false)
			if resolver.HTTP != nil {
				p.enableHTTPDataSource.SetValue(true)
				p.httpURL.SetText(resolver.HTTP.URL)
				p.httpTimeout.SetText(strconv.Itoa(int(resolver.HTTP.Timeout.Seconds())))
			}
		}
	*/

	{
		p.pluginType.Clear()
		p.pluginAddr.Clear()
		p.pluginEnableTLS.SetValue(false)
		p.pluginTLSSecure.SetValue(false)
		p.pluginTLSServerName.Clear()

		if resolver.Plugin != nil {
			p.mode.Value = string(page.PluginMode)
			for i := range page.PluginTypeOptions {
				if page.PluginTypeOptions[i].Value == resolver.Plugin.Type {
					p.pluginType.Select(ui_widget.SelectorItem{Key: page.PluginTypeOptions[i].Key, Value: page.PluginTypeOptions[i].Value})
					break
				}
			}
			p.pluginAddr.SetText(resolver.Plugin.Addr)

			if resolver.Plugin.TLS != nil {
				p.pluginEnableTLS.SetValue(true)
				p.pluginTLSSecure.SetValue(resolver.Plugin.TLS.Secure)
				p.pluginTLSServerName.SetText(resolver.Plugin.TLS.ServerName)
			}
		}
	}
}

func (p *resolverPage) Layout(gtx page.C) page.D {
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
			return layout.Inset{
				Top:    8,
				Bottom: 8,
				Left:   8,
				Right:  8,
			}.Layout(gtx, func(gtx page.C) page.D {
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
						title := material.H6(th, i18n.Resolver.Value())
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

func (p *resolverPage) layout(gtx page.C, th *page.T) page.D {
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
						layout.Flexed(1, layout.Spacer{Width: 4}.Layout),
						layout.Rigid(func(gtx page.C) page.D {
							return material.RadioButton(th, &p.mode, string(page.BasicMode), i18n.Basic.Value()).Layout(gtx)
						}),
						/*
							layout.Rigid(layout.Spacer{Width: 8}.Layout),
							layout.Rigid(func(gtx page.C) page.D {
								return material.RadioButton(th, &p.mode, string(page.AdvancedMode), i18n.Advanced.Value()).Layout(gtx)
							}),
						*/
						layout.Rigid(layout.Spacer{Width: 4}.Layout),
						layout.Rigid(func(gtx page.C) page.D {
							return material.RadioButton(th, &p.mode, string(page.PluginMode), i18n.Plugin.Value()).Layout(gtx)
						}),
					)
				}),

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

				layout.Rigid(func(gtx page.C) page.D {
					if p.mode.Value == string(page.PluginMode) {
						return page.D{}
					}

					if p.addNameserver.Clicked(gtx) {
						p.router.Goto(page.Route{
							Path:     page.PageNameServer,
							Callback: p.callback,
							Perm:     page.PermReadWrite,
						})
					}

					return layout.Flex{
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Flexed(1, func(gtx page.C) page.D {
							return layout.Inset{
								Top:    16,
								Bottom: 16,
							}.Layout(gtx, material.Body1(th, i18n.Nameserver.Value()).Layout)
						}),
						layout.Rigid(func(gtx page.C) page.D {
							if !p.edit {
								return page.D{}
							}
							btn := material.IconButton(th, &p.addNameserver, icons.IconAdd, "Add")
							btn.Background = theme.Current().ContentSurfaceBg
							btn.Color = th.Fg
							return btn.Layout(gtx)
						}),
					)
				}),
				layout.Rigid(func(gtx page.C) page.D {
					if p.mode.Value == string(page.PluginMode) {
						return page.D{}
					}
					gtx.Source = src
					return p.layoutNameservers(gtx, th)
				}),
				layout.Rigid(func(gtx page.C) page.D {
					if p.mode.Value != string(page.AdvancedMode) {
						return page.D{}
					}
					return p.layoutDataSource(gtx, th)
				}),
			)
		})
	})
}

func (p *resolverPage) layoutNameservers(gtx page.C, th *page.T) page.D {
	for i := range p.nameservers {
		if p.nameservers[i].clk.Clicked(gtx) {
			perm := page.PermRead
			if p.edit {
				perm = page.PermReadWrite
			}
			p.router.Goto(page.Route{
				Path:     page.PageNameServer,
				ID:       p.nameservers[i].id,
				Value:    p.nameservers[i].cfg,
				Callback: p.callback,
				Perm:     perm,
			})
			break
		}

		if p.nameservers[i].delete.Clicked(gtx) {
			p.delDialog.Title = i18n.DeleteNode
			p.delDialog.OnClick = func(ok bool) {
				p.router.HideModal(gtx)
				if !ok {
					return
				}
				p.nameservers = append(p.nameservers[:i], p.nameservers[i+1:]...)
			}
			p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
				return p.delDialog.Layout(gtx, th)
			})
			break
		}
	}

	var children []layout.FlexChild
	for i := range p.nameservers {
		nameserver := &p.nameservers[i]

		children = append(children, layout.Rigid(func(gtx page.C) page.D {
			return layout.Flex{
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Flexed(1, func(gtx page.C) page.D {
					return material.Clickable(gtx, &nameserver.clk, func(gtx page.C) page.D {
						return layout.UniformInset(8).Layout(gtx, func(gtx page.C) page.D {
							return layout.Flex{
								Axis: layout.Vertical,
							}.Layout(gtx,
								layout.Rigid(func(gtx page.C) page.D {
									return material.Body1(th, nameserver.cfg.Addr).Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Height: 4}.Layout),
								layout.Rigid(func(gtx page.C) page.D {
									return material.Body2(th, nameserver.cfg.TTL.String()).Layout(gtx)
								}),
							)
						})
					})
				}),
				layout.Rigid(func(gtx page.C) page.D {
					if !p.edit {
						return page.D{}
					}
					btn := material.IconButton(th, &nameserver.delete, icons.IconDelete, "delete")
					btn.Background = theme.Current().ContentSurfaceBg
					btn.Color = th.Fg
					return btn.Layout(gtx)
				}),
			)
		}))
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx, children...)
}

func (p *resolverPage) layoutDataSource(gtx page.C, th *page.T) page.D {
	if p.mode.Value != string(page.AdvancedMode) {
		return page.D{}
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx page.C) page.D {
			return layout.Inset{
				Top:    16,
				Bottom: 16,
			}.Layout(gtx, material.H6(th, i18n.DataSource.Value()).Layout)
		}),

		layout.Rigid(material.Body1(th, i18n.DataSourceReload.Value()).Layout),
		layout.Rigid(func(gtx page.C) page.D {
			return p.reload.Layout(gtx, th, "")
		}),
		layout.Rigid(layout.Spacer{Height: 8}.Layout),

		layout.Rigid(func(gtx page.C) page.D {
			return p.enableFileDataSource.Layout(gtx, th)
		}),
		layout.Rigid(func(gtx page.C) page.D {
			if !p.enableFileDataSource.Value() {
				return page.D{}
			}
			return layout.UniformInset(8).Layout(gtx, func(gtx page.C) page.D {
				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(material.Body1(th, i18n.FilePath.Value()).Layout),
					layout.Rigid(func(gtx page.C) page.D {
						return p.filePath.Layout(gtx, th, "")
					}),
				)
			})
		}),

		layout.Rigid(func(gtx page.C) page.D {
			return p.enableRedisDataSource.Layout(gtx, th)
		}),
		layout.Rigid(func(gtx page.C) page.D {
			if !p.enableRedisDataSource.Value() {
				return page.D{}
			}
			return layout.UniformInset(8).Layout(gtx, func(gtx page.C) page.D {
				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(material.Body1(th, i18n.RedisAddr.Value()).Layout),
					layout.Rigid(func(gtx page.C) page.D {
						return p.redisAddr.Layout(gtx, th, "")
					}),
					layout.Rigid(layout.Spacer{Height: 8}.Layout),

					layout.Rigid(material.Body1(th, i18n.RedisDB.Value()).Layout),
					layout.Rigid(func(gtx page.C) page.D {
						return p.redisDB.Layout(gtx, th, "")
					}),
					layout.Rigid(layout.Spacer{Height: 8}.Layout),

					layout.Rigid(material.Body1(th, i18n.RedisPassword.Value()).Layout),
					layout.Rigid(func(gtx page.C) page.D {
						return p.redisPassword.Layout(gtx, th, "")
					}),
					layout.Rigid(layout.Spacer{Height: 8}.Layout),

					layout.Rigid(material.Body1(th, i18n.RedisKey.Value()).Layout),
					layout.Rigid(func(gtx page.C) page.D {
						return p.redisKey.Layout(gtx, th, "")
					}),
				)
			})
		}),

		layout.Rigid(func(gtx page.C) page.D {
			return p.enableHTTPDataSource.Layout(gtx, th)
		}),
		layout.Rigid(func(gtx page.C) page.D {
			if !p.enableHTTPDataSource.Value() {
				return page.D{}
			}
			return layout.UniformInset(8).Layout(gtx, func(gtx page.C) page.D {
				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(material.Body1(th, i18n.HTTPURL.Value()).Layout),
					layout.Rigid(func(gtx page.C) page.D {
						return p.httpURL.Layout(gtx, th, "")
					}),
					layout.Rigid(layout.Spacer{Height: 8}.Layout),

					layout.Rigid(material.Body1(th, i18n.HTTPTimeout.Value()).Layout),
					layout.Rigid(func(gtx page.C) page.D {
						return p.httpTimeout.Layout(gtx, th, "")
					}),
				)
			})
		}),
	)
}

func (p *resolverPage) showPluginTypeMenu(gtx page.C) {
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

func (p *resolverPage) callback(action page.Action, id string, value any) {
	if id == "" {
		return
	}

	switch action {
	case page.ActionCreate:
		cfg, _ := value.(*api.NameserverConfig)
		if cfg == nil {
			return
		}
		p.nameservers = append(p.nameservers, nameserver{
			id:  id,
			cfg: cfg,
		})
	case page.ActionUpdate:
		cfg, _ := value.(*api.NameserverConfig)
		if cfg == nil {
			return
		}
		for i := range p.nameservers {
			if p.nameservers[i].id == id {
				p.nameservers[i].cfg = cfg
				break
			}
		}
	case page.ActionDelete:
		for i := range p.nameservers {
			if p.nameservers[i].id == id {
				p.nameservers = append(p.nameservers[:i], p.nameservers[i+1:]...)
				break
			}
		}
	}
}

func (p *resolverPage) save() bool {
	cfg := p.generateConfig()

	var err error
	if p.id == "" {
		err = runner.Exec(context.Background(),
			task.CreateResolver(cfg),
			runner.WithCancel(true),
		)
	} else {
		err = runner.Exec(context.Background(),
			task.UpdateResolver(cfg),
			runner.WithCancel(true),
		)
	}
	util.RestartGetConfigTask()

	return err == nil
}

func (p *resolverPage) generateConfig() *api.ResolverConfig {
	cfg := &api.ResolverConfig{
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

	for i := range p.nameservers {
		cfg.Nameservers = append(cfg.Nameservers, p.nameservers[i].cfg)
	}

	/*
		reload, _ := strconv.Atoi(p.reload.Text())
		cfg.Reload = time.Duration(reload) * time.Second

		if p.enableFileDataSource.Value() {
			cfg.File = &api.FileLoader{
				Path: strings.TrimSpace(p.filePath.Text()),
			}
		}

		if p.enableRedisDataSource.Value() {
			db, _ := strconv.Atoi(p.redisDB.Text())
			cfg.Redis = &api.RedisLoader{
				Addr:     strings.TrimSpace(p.redisAddr.Text()),
				DB:       db,
				Password: strings.TrimSpace(p.redisPassword.Text()),
				Key:      strings.TrimSpace(p.redisKey.Text()),
			}
		}

		if p.enableHTTPDataSource.Value() {
			timeout, _ := strconv.Atoi(p.httpTimeout.Text())
			cfg.HTTP = &api.HTTPLoader{
				URL:     strings.TrimSpace(p.httpURL.Text()),
				Timeout: time.Duration(timeout) * time.Second,
			}
		}
	*/

	return cfg
}

func (p *resolverPage) delete() {
	runner.Exec(context.Background(),
		task.DeleteResolver(p.id),
		runner.WithCancel(true),
	)
	util.RestartGetConfigTask()
}
