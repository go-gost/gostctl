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

type metadata struct {
	k      string
	v      string
	clk    widget.Clickable
	delete widget.Clickable
}

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

	bypass     ui_widget.Selector
	resolver   ui_widget.Selector
	hostMapper ui_widget.Selector

	connector *connector
	dialer    *dialer

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

		bypass:     ui_widget.Selector{Title: i18n.Bypass},
		resolver:   ui_widget.Selector{Title: i18n.Resolver},
		hostMapper: ui_widget.Selector{Title: i18n.Hosts},

		delDialog: ui_widget.Dialog{Title: i18n.DeleteNode},
	}
	p.connector = newConnector(p)
	p.dialer = newDialer(p)

	return p
}

func (p *nodePage) Init(opts ...page.PageOption) {
	var options page.PageOptions
	for _, opt := range opts {
		opt(&options)
	}

	p.id = options.ID
	node, _ := options.Value.(*api.NodeConfig)
	if node == nil {
		node = &api.NodeConfig{}
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

	p.resolver.Clear()
	if node.Resolver != "" {
		p.resolver.Select(ui_widget.SelectorItem{Value: node.Resolver})
	}

	p.hostMapper.Clear()
	if node.Hosts != "" {
		p.hostMapper.Select(ui_widget.SelectorItem{Value: node.Hosts})
	}

	p.connector.init(node.Connector)
	p.dialer.init(node.Dialer)
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
				return layout.UniformInset(8).Layout(gtx, func(gtx page.C) page.D {
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
					gtx.Source = src
					return layout.Flex{
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return material.RadioButton(th, &p.mode, string(page.BasicMode), i18n.Basic.Value()).Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: 8}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return material.RadioButton(th, &p.mode, string(page.AdvancedMode), i18n.Advanced.Value()).Layout(gtx)
						}),
					)
				}),

				layout.Rigid(layout.Spacer{Height: 16}.Layout),

				layout.Rigid(material.Body1(th, i18n.Name.Value()).Layout),
				layout.Rigid(func(gtx page.C) page.D {
					return p.name.Layout(gtx, th, "")
				}),
				layout.Rigid(layout.Spacer{Height: 8}.Layout),
				layout.Rigid(material.Body1(th, i18n.Address.Value()).Layout),
				layout.Rigid(func(gtx page.C) page.D {
					return p.addr.Layout(gtx, th, "")
				}),
				layout.Rigid(layout.Spacer{Height: 8}.Layout),

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

						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if p.resolver.Clicked(gtx) {
								p.showResolverMenu(gtx)
							}
							return p.resolver.Layout(gtx, th)
						}),

						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if p.hostMapper.Clicked(gtx) {
								p.showHostMapperMenu(gtx)
							}
							return p.hostMapper.Layout(gtx, th)
						}),

						layout.Rigid(layout.Spacer{Height: 8}.Layout),
					)
				}),
				layout.Rigid(func(gtx page.C) page.D {
					return layout.Inset{
						Top:    16,
						Bottom: 16,
					}.Layout(gtx, material.H6(th, i18n.Connector.Value()).Layout)
				}),
				layout.Rigid(func(gtx page.C) page.D {
					gtx.Source = src
					return p.connector.Layout(gtx, th)
				}),

				layout.Rigid(func(gtx page.C) page.D {
					return layout.Inset{
						Top:    16,
						Bottom: 16,
					}.Layout(gtx, material.H6(th, i18n.Dialer.Value()).Layout)
				}),
				layout.Rigid(func(gtx page.C) page.D {
					gtx.Source = src
					return p.dialer.Layout(gtx, th)
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

func (p *nodePage) showResolverMenu(gtx page.C) {
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
		p.router.Goto(page.Route{
			Path: page.PageResolver,
			Perm: page.PermReadWrite,
		})
		p.router.HideModal(gtx)
	}
	p.menu.Multiple = false

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *nodePage) showHostMapperMenu(gtx page.C) {
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
		p.router.Goto(page.Route{
			Path: page.PageHosts,
			Perm: page.PermReadWrite,
		})
		p.router.HideModal(gtx)
	}
	p.menu.Multiple = false

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *nodePage) generateConfig() *api.NodeConfig {
	node := &api.NodeConfig{
		Name:     strings.TrimSpace(p.name.Text()),
		Addr:     strings.TrimSpace(p.addr.Text()),
		Bypasses: p.bypass.Values(),
		Resolver: p.resolver.Value(),
		Hosts:    p.hostMapper.Value(),
	}

	connector := &api.ConnectorConfig{
		Type:     p.connector.typ.Value(),
		Metadata: make(map[string]any),
	}
	for i := range p.dialer.metadata {
		connector.Metadata[p.connector.metadata[i].k] = p.connector.metadata[i].v
	}

	if p.connector.enableAuth.Value() {
		connector.Auth = &api.AuthConfig{
			Username: strings.TrimSpace(p.connector.username.Text()),
			Password: strings.TrimSpace(p.connector.password.Text()),
		}
	}
	node.Connector = connector

	dialer := &api.DialerConfig{
		Type:     p.dialer.typ.Value(),
		Metadata: make(map[string]any),
	}
	for i := range p.dialer.metadata {
		dialer.Metadata[p.dialer.metadata[i].k] = p.dialer.metadata[i].v
	}

	if p.dialer.enableAuth.Value() {
		dialer.Auth = &api.AuthConfig{
			Username: strings.TrimSpace(p.dialer.username.Text()),
			Password: strings.TrimSpace(p.dialer.password.Text()),
		}
	}
	if p.dialer.enableTLS.Value() {
		dialer.TLS = &api.TLSConfig{
			Secure:     p.dialer.tlsSecure.Value(),
			ServerName: strings.TrimSpace(p.dialer.tlsServerName.Text()),
			CertFile:   strings.TrimSpace(p.dialer.tlsCertFile.Text()),
			KeyFile:    strings.TrimSpace(p.dialer.tlsKeyFile.Text()),
			CAFile:     strings.TrimSpace(p.dialer.tlsCAFile.Text()),
		}
	}
	node.Dialer = dialer

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
