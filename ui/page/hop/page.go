package hop

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

type node struct {
	id     string
	cfg    *api.NodeConfig
	clk    widget.Clickable
	delete widget.Clickable
}

type hopPage struct {
	router *page.Router

	menu ui_widget.Menu
	mode widget.Enum
	list layout.List

	btnBack   widget.Clickable
	btnDelete widget.Clickable
	btnEdit   widget.Clickable
	btnSave   widget.Clickable

	name component.TextField

	enableSelector      ui_widget.Switcher
	selectorStrategy    ui_widget.Selector
	selectorMaxFails    component.TextField
	selectorFailTimeout component.TextField

	bypass     ui_widget.Selector
	resolver   ui_widget.Selector
	hostMapper ui_widget.Selector

	addNode widget.Clickable
	nodes   []node

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
	p := &hopPage{
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

		enableSelector:   ui_widget.Switcher{Title: i18n.Selector},
		selectorStrategy: ui_widget.Selector{Title: i18n.SelectorStrategy},
		selectorMaxFails: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     8,
				Filter:     "1234567890",
			},
		},
		selectorFailTimeout: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     16,
				Filter:     "1234567890",
			},
			Suffix: func(gtx page.C) page.D {
				return material.Body1(r.Theme, i18n.TimeSecond.Value()).Layout(gtx)
			},
		},

		bypass:     ui_widget.Selector{Title: i18n.Bypass},
		resolver:   ui_widget.Selector{Title: i18n.Resolver},
		hostMapper: ui_widget.Selector{Title: i18n.Hosts},

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

		delDialog: ui_widget.Dialog{Title: i18n.DeleteHop},
	}

	return p
}

func (p *hopPage) Init(opts ...page.PageOption) {
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

	hop, _ := options.Value.(*api.HopConfig)

	if hop == nil {
		cfg := api.GetConfig()
		for _, v := range cfg.Hops {
			if v.Name == p.id {
				hop = v
				break
			}
		}
		if hop == nil {
			hop = &api.HopConfig{}
		}
	}

	p.mode.Value = string(page.BasicMode)
	if hop.File != nil || hop.HTTP != nil || hop.Redis != nil {
		p.mode.Value = string(page.AdvancedMode)
	}

	p.name.SetText(hop.Name)

	{
		p.enableSelector.SetValue(false)
		p.selectorStrategy.Clear()
		p.selectorMaxFails.Clear()
		p.selectorFailTimeout.Clear()

		if hop.Selector != nil {
			p.enableSelector.SetValue(true)
			for i := range page.SelectorStrategyOptions {
				if page.SelectorStrategyOptions[i].Value == hop.Selector.Strategy {
					p.selectorStrategy.Select(ui_widget.SelectorItem{Name: page.SelectorStrategyOptions[i].Name, Key: page.SelectorStrategyOptions[i].Key, Value: page.SelectorStrategyOptions[i].Value})
					break
				}
			}
			p.selectorMaxFails.SetText(strconv.Itoa(hop.Selector.MaxFails))
			p.selectorFailTimeout.SetText(hop.Selector.FailTimeout.String())
		}
	}

	{
		p.bypass.Clear()
		var items []ui_widget.SelectorItem
		if hop.Bypass != "" {
			items = append(items, ui_widget.SelectorItem{Value: hop.Bypass})
		}
		for _, v := range hop.Bypasses {
			items = append(items, ui_widget.SelectorItem{
				Value: v,
			})
		}
		p.bypass.Select(items...)
	}

	p.resolver.Clear()
	if hop.Resolver != "" {
		p.resolver.Select(ui_widget.SelectorItem{Value: hop.Resolver})
	}

	p.hostMapper.Clear()
	if hop.Hosts != "" {
		p.hostMapper.Select(ui_widget.SelectorItem{Value: hop.Hosts})
	}

	{
		p.nodes = nil
		for _, v := range hop.Nodes {
			p.nodes = append(p.nodes, node{
				id:  uuid.New().String(),
				cfg: v,
			})
		}
	}

	{
		p.reload.SetText(strconv.Itoa(int(hop.Reload.Seconds())))

		p.enableFileDataSource.SetValue(false)
		if hop.File != nil {
			p.enableFileDataSource.SetValue(true)
			p.filePath.SetText(hop.File.Path)
		}

		p.enableRedisDataSource.SetValue(false)
		if hop.Redis != nil {
			p.enableRedisDataSource.SetValue(true)
			p.redisAddr.SetText(hop.Redis.Addr)
			p.redisDB.SetText(strconv.Itoa(hop.Redis.DB))
			p.redisPassword.SetText(hop.Redis.Password)
			p.redisKey.SetText(hop.Redis.Key)
		}

		p.enableHTTPDataSource.SetValue(false)
		if hop.HTTP != nil {
			p.enableHTTPDataSource.SetValue(true)
			p.httpURL.SetText(hop.HTTP.URL)
			p.httpTimeout.SetText(strconv.Itoa(int(hop.HTTP.Timeout.Seconds())))
		}
	}

	{
		p.pluginType.Clear()
		p.pluginAddr.Clear()
		p.pluginEnableTLS.SetValue(false)
		p.pluginTLSSecure.SetValue(false)
		p.pluginTLSServerName.Clear()

		if hop.Plugin != nil {
			p.mode.Value = string(page.PluginMode)
			for i := range page.PluginTypeOptions {
				if page.PluginTypeOptions[i].Value == hop.Plugin.Type {
					p.pluginType.Select(ui_widget.SelectorItem{Key: page.PluginTypeOptions[i].Key, Value: page.PluginTypeOptions[i].Value})
					break
				}
			}
			p.pluginAddr.SetText(hop.Plugin.Addr)

			if hop.Plugin.TLS != nil {
				p.pluginEnableTLS.SetValue(true)
				p.pluginTLSSecure.SetValue(hop.Plugin.TLS.Secure)
				p.pluginTLSServerName.SetText(hop.Plugin.TLS.ServerName)
			}
		}
	}
}

func (p *hopPage) Layout(gtx page.C) page.D {
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
						title := material.H6(th, i18n.Hop.Value())
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
				return layout.Inset{
					Top:    8,
					Bottom: 8,
					Left:   8,
					Right:  8,
				}.Layout(gtx, func(gtx page.C) page.D {
					return p.layout(gtx, th)
				})
			})
		}),
	)
}

func (p *hopPage) layout(gtx page.C, th *page.T) page.D {
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
						layout.Rigid(layout.Spacer{Width: 8}.Layout),
						layout.Rigid(func(gtx page.C) page.D {
							return material.RadioButton(th, &p.mode, string(page.AdvancedMode), i18n.Advanced.Value()).Layout(gtx)
						}),
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

				// advanced mode
				layout.Rigid(func(gtx page.C) page.D {
					if p.mode.Value != string(page.AdvancedMode) {
						return page.D{}
					}

					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(gtx,
						// Selector
						layout.Rigid(func(gtx page.C) page.D {
							return p.enableSelector.Layout(gtx, th)
						}),

						layout.Rigid(func(gtx page.C) page.D {
							if !p.enableSelector.Value() {
								return page.D{}
							}

							return layout.UniformInset(8).Layout(gtx, func(gtx page.C) page.D {
								return layout.Flex{
									Axis: layout.Vertical,
								}.Layout(gtx,
									layout.Rigid(func(gtx page.C) page.D {
										if p.selectorStrategy.Clicked(gtx) {
											p.showSelectorStrategyMenu(gtx)
										}
										return p.selectorStrategy.Layout(gtx, th)
									}),
									layout.Rigid(func(gtx page.C) page.D {
										return p.selectorMaxFails.Layout(gtx, th, "Max fails")
									}),
									layout.Rigid(func(gtx page.C) page.D {
										return p.selectorFailTimeout.Layout(gtx, th, "Fail timeout in seconds")
									}),
								)
							})
						}),

						layout.Rigid(func(gtx page.C) page.D {
							if p.bypass.Clicked(gtx) {
								p.showBypassMenu(gtx)
							}
							return p.bypass.Layout(gtx, th)
						}),

						layout.Rigid(func(gtx page.C) page.D {
							if p.resolver.Clicked(gtx) {
								p.showResolverMenu(gtx)
							}
							return p.resolver.Layout(gtx, th)
						}),

						layout.Rigid(func(gtx page.C) page.D {
							if p.hostMapper.Clicked(gtx) {
								p.showHostMapperMenu(gtx)
							}
							return p.hostMapper.Layout(gtx, th)
						}),

						layout.Rigid(layout.Spacer{Height: 8}.Layout),
					)
				}),

				layout.Rigid(func(gtx page.C) page.D {
					if p.mode.Value == string(page.PluginMode) {
						return page.D{}
					}

					if p.addNode.Clicked(gtx) {
						p.router.Goto(page.Route{
							Path:     page.PageNode,
							Callback: p.nodeCallback,
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
							}.Layout(gtx, material.Body1(th, i18n.Nodes.Value()).Layout)
						}),
						layout.Rigid(func(gtx page.C) page.D {
							if !p.edit {
								return page.D{}
							}
							btn := material.IconButton(th, &p.addNode, icons.IconAdd, "Add")
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
					return p.layoutNodes(gtx, th)
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

func (p *hopPage) layoutNodes(gtx page.C, th *page.T) page.D {
	for i := range p.nodes {
		if p.nodes[i].clk.Clicked(gtx) {
			perm := page.PermRead
			if p.edit {
				perm = page.PermReadWrite
			}
			p.router.Goto(page.Route{
				Path:     page.PageNode,
				ID:       p.nodes[i].id,
				Value:    p.nodes[i].cfg,
				Callback: p.nodeCallback,
				Perm:     perm,
			})
			break
		}

		if p.nodes[i].delete.Clicked(gtx) {
			p.delDialog.Title = i18n.DeleteNode
			p.delDialog.OnClick = func(ok bool) {
				p.router.HideModal(gtx)
				if !ok {
					return
				}
				p.nodes = append(p.nodes[:i], p.nodes[i+1:]...)
			}
			p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
				return p.delDialog.Layout(gtx, th)
			})
			break
		}
	}

	var children []layout.FlexChild
	for i := range p.nodes {
		node := &p.nodes[i]

		children = append(children, layout.Rigid(func(gtx page.C) page.D {
			return layout.Flex{
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Flexed(1, func(gtx page.C) page.D {
					return material.Clickable(gtx, &node.clk, func(gtx page.C) page.D {
						return layout.Inset{
							Top:    8,
							Bottom: 8,
							Left:   8,
							Right:  8,
						}.Layout(gtx, func(gtx page.C) page.D {
							return layout.Flex{
								Axis: layout.Vertical,
							}.Layout(gtx,
								layout.Rigid(func(gtx page.C) page.D {
									return material.Body1(th, node.cfg.Name).Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Height: 4}.Layout),
								layout.Rigid(func(gtx page.C) page.D {
									return material.Body2(th, node.cfg.Addr).Layout(gtx)
								}),
							)
						})
					})
				}),
				layout.Rigid(func(gtx page.C) page.D {
					if !p.edit {
						return page.D{}
					}
					btn := material.IconButton(th, &node.delete, icons.IconDelete, "delete")
					btn.Background = theme.Current().ContentSurfaceBg
					btn.Color = th.Fg
					// btn.Inset = layout.UniformInset(8)
					return btn.Layout(gtx)
				}),
			)
		}))
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx, children...)
}

func (p *hopPage) layoutDataSource(gtx page.C, th *page.T) page.D {
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

var ()

func (p *hopPage) showSelectorStrategyMenu(gtx page.C) {
	for i := range page.SelectorStrategyOptions {
		page.SelectorStrategyOptions[i].Selected = p.selectorStrategy.AnyValue(page.SelectorStrategyOptions[i].Value)
	}

	p.menu.Title = i18n.SelectorStrategy
	p.menu.Options = page.SelectorStrategyOptions
	p.menu.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		p.selectorStrategy.Clear()
		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				p.selectorStrategy.Select(ui_widget.SelectorItem{Key: p.menu.Options[i].Key, Value: p.menu.Options[i].Value})
			}
		}
	}
	p.menu.OnAdd = nil
	p.menu.Multiple = false

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *hopPage) showBypassMenu(gtx page.C) {
	options := []ui_widget.MenuOption{}
	for _, v := range api.GetConfig().Bypasses {
		options = append(options, ui_widget.MenuOption{
			Value: v.Name,
		})
	}
	for i := range options {
		options[i].Selected = p.bypass.AnyValue(options[i].Value)
	}

	p.menu.Title = i18n.Bypass
	p.menu.Options = options
	p.menu.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		p.bypass.Clear()
		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				p.bypass.Select(ui_widget.SelectorItem{Value: p.menu.Options[i].Value})
			}
		}
	}
	p.menu.OnAdd = func() {
		p.router.Goto(page.Route{
			Path: page.PageBypass,
			Perm: page.PermReadWrite,
		})
		p.router.HideModal(gtx)
	}
	p.menu.Multiple = true

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *hopPage) showResolverMenu(gtx page.C) {
	options := []ui_widget.MenuOption{}
	for _, v := range api.GetConfig().Resolvers {
		options = append(options, ui_widget.MenuOption{
			Value: v.Name,
		})
	}
	for i := range options {
		options[i].Selected = p.resolver.AnyValue(options[i].Value)
	}

	p.menu.Title = i18n.Resolver
	p.menu.Options = options
	p.menu.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		p.resolver.Clear()
		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				p.resolver.Select(ui_widget.SelectorItem{Value: p.menu.Options[i].Value})
			}
		}
	}
	p.menu.OnAdd = func() {

	}
	p.menu.Multiple = false

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *hopPage) showHostMapperMenu(gtx page.C) {
	options := []ui_widget.MenuOption{}
	for _, v := range api.GetConfig().Hosts {
		options = append(options, ui_widget.MenuOption{
			Value: v.Name,
		})
	}
	for i := range options {
		options[i].Selected = p.hostMapper.AnyValue(options[i].Value)
	}

	p.menu.Title = i18n.Hosts
	p.menu.Options = options
	p.menu.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		p.hostMapper.Clear()
		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				p.hostMapper.Select(ui_widget.SelectorItem{Value: p.menu.Options[i].Value})
			}
		}
	}
	p.menu.OnAdd = func() {

	}
	p.menu.Multiple = false

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *hopPage) showPluginTypeMenu(gtx page.C) {
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

func (p *hopPage) nodeCallback(action page.Action, id string, value any) {
	if id == "" {
		return
	}

	switch action {
	case page.ActionCreate:
		cfg, _ := value.(*api.NodeConfig)
		if cfg == nil {
			return
		}
		p.nodes = append(p.nodes, node{
			id:  id,
			cfg: cfg,
		})
	case page.ActionUpdate:
		cfg, _ := value.(*api.NodeConfig)
		if cfg == nil {
			return
		}
		for i := range p.nodes {
			if p.nodes[i].id == id {
				p.nodes[i].cfg = cfg
				break
			}
		}
	case page.ActionDelete:
		for i := range p.nodes {
			if p.nodes[i].id == id {
				p.nodes = append(p.nodes[:i], p.nodes[i+1:]...)
				break
			}
		}
	}
}

func (p *hopPage) save() bool {
	cfg := p.generateConfig()

	var err error
	if p.id == "" {
		err = runner.Exec(context.Background(),
			task.CreateHop(cfg),
			runner.WithCancel(true),
		)
	} else {
		err = runner.Exec(context.Background(),
			task.UpdateHop(cfg),
			runner.WithCancel(true),
		)
	}
	util.RestartGetConfigTask()

	return err == nil
}

func (p *hopPage) generateConfig() *api.HopConfig {
	cfg := &api.HopConfig{
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

	if p.selectorStrategy.Value() != "" {
		maxFails, _ := strconv.Atoi(p.selectorMaxFails.Text())
		failTimeout, _ := strconv.Atoi(p.selectorFailTimeout.Text())
		cfg.Selector = &api.SelectorConfig{
			Strategy:    p.selectorStrategy.Value(),
			MaxFails:    maxFails,
			FailTimeout: time.Duration(failTimeout) * time.Second,
		}
	}
	cfg.Bypasses = p.bypass.Values()
	cfg.Resolver = p.resolver.Value()
	cfg.Hosts = p.hostMapper.Value()

	for i := range p.nodes {
		cfg.Nodes = append(cfg.Nodes, p.nodes[i].cfg)
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

func (p *hopPage) delete() {
	runner.Exec(context.Background(),
		task.DeleteHop(p.id),
		runner.WithCancel(true),
	)
	util.RestartGetConfigTask()
}
