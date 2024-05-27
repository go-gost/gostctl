package node

import (
	"strings"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/go-gost/gostctl/api"
	"github.com/go-gost/gostctl/ui/i18n"
	"github.com/go-gost/gostctl/ui/icons"
	"github.com/go-gost/gostctl/ui/page"
	"github.com/go-gost/gostctl/ui/theme"
	ui_widget "github.com/go-gost/gostctl/ui/widget"
	"github.com/google/uuid"
)

type nodePage struct {
	router *page.Router

	menu ui_widget.Menu
	mode widget.Enum
	list layout.List

	btnBack   widget.Clickable
	btnDelete widget.Clickable
	btnEdit   widget.Clickable
	btnSave   widget.Clickable

	name component.TextField
	addr component.TextField

	bypass ui_widget.Selector

	enableFilter   ui_widget.Switcher
	protocolFilter ui_widget.Selector
	hostFilter     component.TextField
	pathFilter     component.TextField

	enableHTTP   ui_widget.Switcher
	httpHost     component.TextField
	httpUsername component.TextField
	httpPassword component.TextField

	enableTLS     ui_widget.Switcher
	tlsSecure     ui_widget.Switcher
	tlsServerName component.TextField

	id       string
	perm     page.Perm
	callback page.Callback

	edit   bool
	create bool

	delDialog ui_widget.Dialog
}

func NewPage(r *page.Router) page.Page {
	p := &nodePage{
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
		addr: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     128,
			},
		},

		bypass:         ui_widget.Selector{Title: i18n.Bypass},
		enableFilter:   ui_widget.Switcher{Title: i18n.Filter},
		protocolFilter: ui_widget.Selector{Title: i18n.Protocol},
		hostFilter: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     255,
			},
		},
		pathFilter: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     255,
			},
		},
		enableHTTP: ui_widget.Switcher{Title: i18n.HTTP},
		httpHost: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     255,
			},
		},
		httpUsername: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     255,
			},
		},
		httpPassword: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     255,
			},
		},

		enableTLS: ui_widget.Switcher{Title: i18n.TLS},
		tlsSecure: ui_widget.Switcher{Title: i18n.VerifyServerCert},
		tlsServerName: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     255,
			},
		},

		delDialog: ui_widget.Dialog{Title: i18n.DeleteService},
	}

	return p
}

func (p *nodePage) Init(opts ...page.PageOption) {
	var options page.PageOptions
	for _, opt := range opts {
		opt(&options)
	}

	p.id = options.ID
	node, _ := options.Value.(*api.ForwardNodeConfig)
	if node == nil {
		node = &api.ForwardNodeConfig{}
	}
	p.callback = options.Callback

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

	p.mode.Value = string(page.BasicMode)

	p.name.SetText(node.Name)
	p.addr.SetText(node.Addr)

	{
		p.bypass.Clear()
		var items []ui_widget.SelectorItem
		if node.Bypass != "" {
			items = append(items, ui_widget.SelectorItem{Value: node.Bypass})
		}
		for _, v := range node.Bypasses {
			items = append(items, ui_widget.SelectorItem{
				Value: v,
			})
		}
		p.bypass.Select(items...)
	}

	p.enableFilter.SetValue(false)
	filter := node.Filter
	if filter == nil {
		if node.Protocol != "" || node.Host != "" || node.Path != "" {
			filter = &api.NodeFilterConfig{
				Protocol: node.Protocol,
				Host:     node.Host,
				Path:     node.Path,
			}
		}
	}
	if filter != nil {
		p.enableFilter.SetValue(true)
		p.protocolFilter.Clear()
		for i := range protocolOptions {
			if protocolOptions[i].Value == filter.Protocol {
				p.protocolFilter.Select(ui_widget.SelectorItem{Name: protocolOptions[i].Name, Key: protocolOptions[i].Key, Value: protocolOptions[i].Value})
				break
			}
		}
		p.hostFilter.SetText(filter.Host)
		p.pathFilter.SetText(filter.Path)
	}

	p.enableHTTP.SetValue(false)
	if node.HTTP != nil {
		p.enableHTTP.SetValue(true)
		p.httpHost.SetText(node.HTTP.Host)
		if node.HTTP.Auth != nil {
			p.httpUsername.SetText(node.HTTP.Auth.Username)
			p.httpPassword.SetText(node.HTTP.Auth.Password)
		}
	}

	p.enableTLS.SetValue(false)
	if node.TLS != nil {
		p.enableTLS.SetValue(true)
		p.tlsSecure.SetValue(node.TLS.Secure)
		p.tlsServerName.SetText(node.TLS.ServerName)
	}
}

func (p *nodePage) Layout(gtx page.C) page.D {
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
						title := material.H6(th, i18n.Node.Value())
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

func (p *nodePage) layout(gtx page.C, th *page.T) page.D {
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
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Flexed(1, layout.Spacer{Width: 8}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							gtx.Source = src
							return material.RadioButton(th, &p.mode, string(page.BasicMode), i18n.Basic.Value()).Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: 8}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							gtx.Source = src
							return material.RadioButton(th, &p.mode, string(page.AdvancedMode), i18n.Advanced.Value()).Layout(gtx)
						}),
					)
				}),

				layout.Rigid(layout.Spacer{Height: 16}.Layout),

				layout.Rigid(material.Body1(th, i18n.Name.Value()).Layout),
				layout.Rigid(func(gtx page.C) page.D {
					return p.name.Layout(gtx, th, "")
				}),
				layout.Rigid(layout.Spacer{Height: 4}.Layout),
				layout.Rigid(material.Body1(th, i18n.Address.Value()).Layout),
				layout.Rigid(func(gtx page.C) page.D {
					return p.addr.Layout(gtx, th, "")
				}),
				layout.Rigid(layout.Spacer{Height: 4}.Layout),

				// advanced mode
				layout.Rigid(func(gtx page.C) page.D {
					if p.mode.Value == string(page.BasicMode) {
						return page.D{}
					}

					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if p.bypass.Clicked(gtx) {
								p.showBypassMenu(gtx)
							}
							return p.bypass.Layout(gtx, th)
						}),
						layout.Rigid(func(gtx page.C) page.D {
							return p.enableFilter.Layout(gtx, th)
						}),

						layout.Rigid(func(gtx page.C) page.D {
							if !p.enableFilter.Value() {
								return page.D{}
							}

							return layout.UniformInset(8).Layout(gtx, func(gtx page.C) page.D {
								return layout.Flex{
									Axis: layout.Vertical,
								}.Layout(gtx,
									layout.Rigid(func(gtx page.C) page.D {
										if p.protocolFilter.Clicked(gtx) {
											p.showProtocolMenu(gtx)
										}
										return p.protocolFilter.Layout(gtx, th)
									}),

									layout.Rigid(func(gtx page.C) page.D {
										return material.Body1(th, i18n.Host.Value()).Layout(gtx)
									}),
									layout.Rigid(func(gtx page.C) page.D {
										return p.hostFilter.Layout(gtx, th, "")
									}),
									layout.Rigid(layout.Spacer{Height: 8}.Layout),

									layout.Rigid(func(gtx page.C) page.D {
										return material.Body1(th, i18n.Path.Value()).Layout(gtx)
									}),
									layout.Rigid(func(gtx page.C) page.D {
										return p.pathFilter.Layout(gtx, th, "")
									}),
								)
							})
						}),

						layout.Rigid(func(gtx page.C) page.D {
							return p.enableHTTP.Layout(gtx, th)
						}),

						layout.Rigid(func(gtx page.C) page.D {
							if !p.enableHTTP.Value() {
								return page.D{}
							}

							return layout.UniformInset(8).Layout(gtx, func(gtx page.C) page.D {
								return layout.Flex{
									Axis: layout.Vertical,
								}.Layout(gtx,
									layout.Rigid(func(gtx page.C) page.D {
										return p.httpHost.Layout(gtx, th, i18n.RewriteHostHeader.Value())
									}),
									layout.Rigid(func(gtx page.C) page.D {
										return p.httpUsername.Layout(gtx, th, i18n.Username.Value())
									}),
									layout.Rigid(func(gtx page.C) page.D {
										return p.httpPassword.Layout(gtx, th, i18n.Password.Value())
									}),
									layout.Rigid(layout.Spacer{Height: 4}.Layout),
								)
							})
						}),

						layout.Rigid(func(gtx page.C) page.D {
							return p.enableTLS.Layout(gtx, th)
						}),

						layout.Rigid(func(gtx page.C) page.D {
							if !p.enableTLS.Value() {
								return page.D{}
							}

							return layout.UniformInset(8).Layout(gtx, func(gtx page.C) page.D {
								return layout.Flex{
									Axis: layout.Vertical,
								}.Layout(gtx,
									layout.Rigid(func(gtx page.C) page.D {
										return p.tlsSecure.Layout(gtx, th)
									}),

									layout.Rigid(func(gtx page.C) page.D {
										return material.Body1(th, i18n.ServerName.Value()).Layout(gtx)
									}),

									layout.Rigid(func(gtx page.C) page.D {
										return p.tlsServerName.Layout(gtx, th, "")
									}),
								)
							})
						}),
					)
				}),
			)
		})
	})
}

func (p *nodePage) showBypassMenu(gtx page.C) {
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

var (
	protocolOptions = []ui_widget.MenuOption{
		{Name: "HTTP", Value: "http"},
		{Name: "TLS", Value: "tls"},
		{Name: "SSH", Value: "ssh"},
	}
)

func (p *nodePage) showProtocolMenu(gtx page.C) {
	for i := range protocolOptions {
		protocolOptions[i].Selected = p.protocolFilter.AnyValue(protocolOptions[i].Value)
	}

	p.menu.Title = i18n.Protocol
	p.menu.Options = protocolOptions
	p.menu.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		p.protocolFilter.Clear()
		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				p.protocolFilter.Select(ui_widget.SelectorItem{Name: p.menu.Options[i].Name, Key: p.menu.Options[i].Key, Value: p.menu.Options[i].Value})
			}
		}
	}
	p.menu.OnAdd = nil
	p.menu.Multiple = false

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *nodePage) generateConfig() *api.ForwardNodeConfig {
	node := &api.ForwardNodeConfig{
		Name:     strings.TrimSpace(p.name.Text()),
		Addr:     strings.TrimSpace(p.addr.Text()),
		Bypasses: p.bypass.Values(),
	}
	if p.enableFilter.Value() {
		node.Filter = &api.NodeFilterConfig{
			Protocol: p.protocolFilter.Value(),
			Host:     strings.TrimSpace(p.hostFilter.Text()),
			Path:     strings.TrimSpace(p.pathFilter.Text()),
		}
	}
	if p.enableHTTP.Value() {
		node.HTTP = &api.HTTPNodeConfig{
			Host: strings.TrimSpace(p.httpHost.Text()),
		}
		username := strings.TrimSpace(p.httpUsername.Text())
		password := strings.TrimSpace(p.httpPassword.Text())
		if username != "" {
			node.HTTP.Auth = &api.AuthConfig{
				Username: username,
				Password: password,
			}
		}
	}
	if p.enableTLS.Value() {
		node.TLS = &api.TLSNodeConfig{
			ServerName: strings.TrimSpace(p.tlsServerName.Text()),
			Secure:     p.tlsSecure.Value(),
		}
	}

	return node
}

func (p *nodePage) save() bool {
	node := p.generateConfig()

	if p.id == "" {
		if p.callback != nil {
			p.callback(page.ActionCreate, uuid.New().String(), node)
		}

	} else {
		if p.callback != nil {
			p.callback(page.ActionUpdate, p.id, node)
		}
	}

	return true
}

func (p *nodePage) delete() {
	if p.callback != nil {
		p.callback(page.ActionDelete, p.id, nil)
	}
}
