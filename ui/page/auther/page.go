package auther

import (
	"context"
	"strconv"
	"strings"
	"time"

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

type autherPage struct {
	router *page.Router

	menu ui_widget.Menu
	mode widget.Enum
	list layout.List

	btnBack   widget.Clickable
	btnDelete widget.Clickable
	btnEdit   widget.Clickable
	btnSave   widget.Clickable

	name component.TextField

	auths        []*api.AuthConfig
	authSelector ui_widget.Selector

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
	p := &autherPage{
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

		authSelector: ui_widget.Selector{Title: i18n.Auths},

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

		delDialog: ui_widget.Dialog{Title: i18n.DeleteAuther},
	}

	return p
}

func (p *autherPage) Init(opts ...page.PageOption) {
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

	auther, _ := options.Value.(*api.AutherConfig)

	if auther == nil {
		cfg := api.GetConfig()
		for _, v := range cfg.Authers {
			if v.Name == p.id {
				auther = v
				break
			}
		}
		if auther == nil {
			auther = &api.AutherConfig{}
		}
	}

	p.mode.Value = string(page.BasicMode)
	if auther.File != nil || auther.HTTP != nil || auther.Redis != nil {
		p.mode.Value = string(page.AdvancedMode)
	}

	p.name.SetText(auther.Name)
	p.auths = auther.Auths
	p.authSelector.Clear()
	p.authSelector.Select(ui_widget.SelectorItem{Value: strconv.Itoa(len(p.auths))})

	{
		p.reload.SetText(strconv.Itoa(int(auther.Reload.Seconds())))

		p.enableFileDataSource.SetValue(false)
		if auther.File != nil {
			p.enableFileDataSource.SetValue(true)
			p.filePath.SetText(auther.File.Path)
		}

		p.enableRedisDataSource.SetValue(false)
		if auther.Redis != nil {
			p.enableRedisDataSource.SetValue(true)
			p.redisAddr.SetText(auther.Redis.Addr)
			p.redisDB.SetText(strconv.Itoa(auther.Redis.DB))
			p.redisPassword.SetText(auther.Redis.Password)
			p.redisKey.SetText(auther.Redis.Key)
		}

		p.enableHTTPDataSource.SetValue(false)
		if auther.HTTP != nil {
			p.enableHTTPDataSource.SetValue(true)
			p.httpURL.SetText(auther.HTTP.URL)
			p.httpTimeout.SetText(strconv.Itoa(int(auther.HTTP.Timeout.Seconds())))
		}
	}

	{
		p.pluginType.Clear()
		p.pluginAddr.Clear()
		p.pluginEnableTLS.SetValue(false)
		p.pluginTLSSecure.SetValue(false)
		p.pluginTLSServerName.Clear()

		if auther.Plugin != nil {
			p.mode.Value = string(page.PluginMode)
			for i := range page.PluginTypeOptions {
				if page.PluginTypeOptions[i].Value == auther.Plugin.Type {
					p.pluginType.Select(ui_widget.SelectorItem{Key: page.PluginTypeOptions[i].Key, Value: page.PluginTypeOptions[i].Value})
					break
				}
			}
			p.pluginAddr.SetText(auther.Plugin.Addr)

			if auther.Plugin.TLS != nil {
				p.pluginEnableTLS.SetValue(true)
				p.pluginTLSSecure.SetValue(auther.Plugin.TLS.Secure)
				p.pluginTLSServerName.SetText(auther.Plugin.TLS.ServerName)
			}
		}
	}
}

func (p *autherPage) Layout(gtx page.C) page.D {
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
						title := material.H6(th, i18n.Auther.Value())
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

func (p *autherPage) layout(gtx page.C, th *page.T) page.D {
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
							return material.RadioButton(th, &p.mode, string(page.BasicMode), i18n.Basic.Value()).Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: 8}.Layout),
						layout.Rigid(func(gtx page.C) page.D {
							return material.RadioButton(th, &p.mode, string(page.AdvancedMode), i18n.Advanced.Value()).Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: 8}.Layout),
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
				layout.Rigid(func(gtx page.C) page.D {
					if p.mode.Value == string(page.PluginMode) {
						return page.D{}
					}
					gtx.Source = src

					if p.authSelector.Clicked(gtx) {
						perm := page.PermRead
						if p.edit {
							perm = page.PermWrite | page.PermDelete
						}
						p.router.Goto(page.Route{
							Path:     page.PageAutherAuths,
							ID:       uuid.New().String(),
							Value:    p.auths,
							Callback: p.authCallback,
							Perm:     perm,
						})
					}
					return p.authSelector.Layout(gtx, th)
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

func (p *autherPage) layoutDataSource(gtx page.C, th *page.T) page.D {
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

func (p *autherPage) showPluginTypeMenu(gtx page.C) {
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

func (p *autherPage) authCallback(action page.Action, id string, value any) {
	if id == "" {
		return
	}

	switch action {
	case page.ActionUpdate:
		p.auths, _ = value.([]*api.AuthConfig)

	case page.ActionDelete:
		p.auths = nil
	}

	p.authSelector.Clear()
	p.authSelector.Select(ui_widget.SelectorItem{Value: strconv.Itoa(len(p.auths))})
}

func (p *autherPage) save() bool {
	cfg := p.generateConfig()

	var err error
	if p.id == "" {
		err = runner.Exec(context.Background(),
			task.CreateAuther(cfg),
			runner.WithCancel(true),
		)
	} else {
		err = runner.Exec(context.Background(),
			task.UpdateAuther(cfg),
			runner.WithCancel(true),
		)
	}
	util.RestartGetConfigTask()

	return err == nil
}

func (p *autherPage) generateConfig() *api.AutherConfig {
	cfg := &api.AutherConfig{
		Name:  strings.TrimSpace(p.name.Text()),
		Auths: p.auths,
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

	return cfg
}

func (p *autherPage) delete() {
	runner.Exec(context.Background(),
		task.DeleteAuther(p.id),
		runner.WithCancel(true),
	)
	util.RestartGetConfigTask()
}
