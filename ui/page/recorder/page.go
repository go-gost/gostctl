package recorder

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
	"github.com/go-gost/gostctl/config"
	"github.com/go-gost/gostctl/ui/i18n"
	"github.com/go-gost/gostctl/ui/icons"
	"github.com/go-gost/gostctl/ui/page"
	"github.com/go-gost/gostctl/ui/theme"
	ui_widget "github.com/go-gost/gostctl/ui/widget"
)

type recorderPage struct {
	readonly bool
	router   *page.Router

	menu ui_widget.Menu
	mode widget.Enum
	list layout.List

	btnBack   widget.Clickable
	btnDelete widget.Clickable
	btnEdit   widget.Clickable
	btnSave   widget.Clickable

	name component.TextField

	enableFileDataSource ui_widget.Switcher
	filePath             component.TextField
	fileSep              component.TextField

	enableRedisDataSource ui_widget.Switcher
	redisAddr             component.TextField
	redisDB               component.TextField
	redisPassword         component.TextField
	redisKey              component.TextField

	enableHTTPDataSource ui_widget.Switcher
	httpURL              component.TextField
	httpTimeout          component.TextField

	enableTCPDataSource ui_widget.Switcher
	tcpAddr             component.TextField
	tcpTimeout          component.TextField

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
	p := &recorderPage{
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

		enableFileDataSource: ui_widget.Switcher{Title: i18n.FileDataSource},
		filePath: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     255,
			},
		},
		fileSep: component.TextField{
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

		enableTCPDataSource: ui_widget.Switcher{Title: i18n.TCPDataSource},
		tcpAddr: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     255,
			},
		},
		tcpTimeout: component.TextField{
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

		delDialog: ui_widget.Dialog{Title: i18n.DeleteRecorder},
	}

	return p
}

func (p *recorderPage) Init(opts ...page.PageOption) {
	if server := config.CurrentServer(); server != nil {
		p.readonly = server.Readonly
	}

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

	recorder, _ := options.Value.(*api.RecorderConfig)

	if recorder == nil {
		cfg := api.GetConfig()
		for _, v := range cfg.Recorders {
			if v.Name == p.id {
				recorder = v
				break
			}
		}
		if recorder == nil {
			recorder = &api.RecorderConfig{}
		}
	}

	p.mode.Value = string(page.BasicMode)
	p.name.SetText(recorder.Name)

	{
		p.enableFileDataSource.SetValue(false)
		if recorder.File != nil {
			p.enableFileDataSource.SetValue(true)
			p.filePath.SetText(recorder.File.Path)
			p.fileSep.SetText(recorder.File.Sep)
		}

		p.enableRedisDataSource.SetValue(false)
		if recorder.Redis != nil {
			p.enableRedisDataSource.SetValue(true)
			p.redisAddr.SetText(recorder.Redis.Addr)
			p.redisDB.SetText(strconv.Itoa(recorder.Redis.DB))
			p.redisPassword.SetText(recorder.Redis.Password)
			p.redisKey.SetText(recorder.Redis.Key)
		}

		p.enableHTTPDataSource.SetValue(false)
		if recorder.HTTP != nil {
			p.enableHTTPDataSource.SetValue(true)
			p.httpURL.SetText(recorder.HTTP.URL)
			p.httpTimeout.SetText(strconv.Itoa(int(recorder.HTTP.Timeout.Seconds())))
		}

		p.enableTCPDataSource.SetValue(false)
		if recorder.TCP != nil {
			p.enableTCPDataSource.SetValue(true)
			p.tcpAddr.SetText(recorder.TCP.Addr)
			p.tcpTimeout.SetText(strconv.Itoa(int(recorder.TCP.Timeout.Seconds())))
		}
	}

	{
		p.pluginType.Clear()
		p.pluginAddr.Clear()
		p.pluginEnableTLS.SetValue(false)
		p.pluginTLSSecure.SetValue(false)
		p.pluginTLSServerName.Clear()

		if recorder.Plugin != nil {
			p.mode.Value = string(page.PluginMode)
			for i := range page.PluginTypeOptions {
				if page.PluginTypeOptions[i].Value == recorder.Plugin.Type {
					p.pluginType.Select(ui_widget.SelectorItem{Key: page.PluginTypeOptions[i].Key, Value: page.PluginTypeOptions[i].Value})
					break
				}
			}
			p.pluginAddr.SetText(recorder.Plugin.Addr)

			if recorder.Plugin.TLS != nil {
				p.pluginEnableTLS.SetValue(true)
				p.pluginTLSSecure.SetValue(recorder.Plugin.TLS.Secure)
				p.pluginTLSServerName.SetText(recorder.Plugin.TLS.ServerName)
			}
		}
	}
}

func (p *recorderPage) Layout(gtx page.C) page.D {
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
						title := material.H6(th, i18n.Recorder.Value())
						return title.Layout(gtx)
					}),
					layout.Rigid(func(gtx page.C) page.D {
						if p.readonly || p.perm&page.PermDelete == 0 || p.create {
							return page.D{}
						}

						btn := material.IconButton(th, &p.btnDelete, icons.IconDelete, "Delete")
						btn.Color = th.Fg
						btn.Background = th.Bg
						return btn.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: 8}.Layout),
					layout.Rigid(func(gtx page.C) page.D {
						if p.readonly || p.perm&page.PermWrite == 0 {
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

func (p *recorderPage) layout(gtx page.C, th *page.T) page.D {
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
					return p.layoutDataSource(gtx, th)
				}),
			)
		})
	})
}

func (p *recorderPage) layoutDataSource(gtx page.C, th *page.T) page.D {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
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

					layout.Rigid(layout.Spacer{Height: 8}.Layout),
					layout.Rigid(material.Body1(th, i18n.FileSep.Value()).Layout),
					layout.Rigid(func(gtx page.C) page.D {
						return p.fileSep.Layout(gtx, th, "")
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

		layout.Rigid(func(gtx page.C) page.D {
			return p.enableTCPDataSource.Layout(gtx, th)
		}),
		layout.Rigid(func(gtx page.C) page.D {
			if !p.enableTCPDataSource.Value() {
				return page.D{}
			}
			return layout.UniformInset(8).Layout(gtx, func(gtx page.C) page.D {
				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(material.Body1(th, i18n.TCPAddr.Value()).Layout),
					layout.Rigid(func(gtx page.C) page.D {
						return p.tcpAddr.Layout(gtx, th, "")
					}),
					layout.Rigid(layout.Spacer{Height: 8}.Layout),

					layout.Rigid(material.Body1(th, i18n.TCPTimeout.Value()).Layout),
					layout.Rigid(func(gtx page.C) page.D {
						return p.tcpTimeout.Layout(gtx, th, "")
					}),
				)
			})
		}),
	)
}

func (p *recorderPage) showPluginTypeMenu(gtx page.C) {
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

func (p *recorderPage) save() bool {
	cfg := p.generateConfig()

	var err error
	if p.id == "" {
		err = runner.Exec(context.Background(),
			task.CreateRecorder(cfg),
			runner.WithCancel(true),
		)
	} else {
		err = runner.Exec(context.Background(),
			task.UpdateRecorder(cfg),
			runner.WithCancel(true),
		)
	}
	util.RestartGetConfigTask()

	return err == nil
}

func (p *recorderPage) generateConfig() *api.RecorderConfig {
	cfg := &api.RecorderConfig{
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

	if p.enableFileDataSource.Value() {
		cfg.File = &api.FileRecorder{
			Path: strings.TrimSpace(p.filePath.Text()),
		}
	}

	if p.enableRedisDataSource.Value() {
		db, _ := strconv.Atoi(p.redisDB.Text())
		cfg.Redis = &api.RedisRecorder{
			Addr:     strings.TrimSpace(p.redisAddr.Text()),
			DB:       db,
			Password: strings.TrimSpace(p.redisPassword.Text()),
			Key:      strings.TrimSpace(p.redisKey.Text()),
		}
	}

	if p.enableHTTPDataSource.Value() {
		timeout, _ := strconv.Atoi(p.httpTimeout.Text())
		cfg.HTTP = &api.HTTPRecorder{
			URL:     strings.TrimSpace(p.httpURL.Text()),
			Timeout: time.Duration(timeout) * time.Second,
		}
	}

	if p.enableTCPDataSource.Value() {
		timeout, _ := strconv.Atoi(p.tcpTimeout.Text())
		cfg.TCP = &api.TCPRecorder{
			Addr:    strings.TrimSpace(p.tcpAddr.Text()),
			Timeout: time.Duration(timeout) * time.Second,
		}
	}

	return cfg
}

func (p *recorderPage) delete() {
	runner.Exec(context.Background(),
		task.DeleteRecorder(p.id),
		runner.WithCancel(true),
	)
	util.RestartGetConfigTask()
}
