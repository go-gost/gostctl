package service

import (
	"context"
	"strconv"
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

type servicePage struct {
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

	// stats ui_widget.Switcher

	// customInterface ui_widget.Selector
	// interfaceDialog ui_widget.InputDialog

	admission  ui_widget.Selector
	bypass     ui_widget.Selector
	resolver   ui_widget.Selector
	hostMapper ui_widget.Selector
	limiter    ui_widget.Selector
	logger     ui_widget.Selector
	observer   ui_widget.Selector

	id   string
	perm page.Perm

	edit   bool
	create bool

	metadata         []page.Metadata
	metadataSelector ui_widget.Selector
	metadataDialog   ui_widget.MetadataDialog
	delDialog        ui_widget.Dialog

	handler   handler
	listener  listener
	forwarder forwarder
}

func NewPage(r *page.Router) page.Page {
	p := &servicePage{
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
				MaxLen:     255,
			},
		},
		delDialog: ui_widget.Dialog{Title: i18n.DeleteService},

		// stats:           ui_widget.Switcher{Title: "Stats"},
		/*
			customInterface: ui_widget.Selector{Title: i18n.Interface},
			interfaceDialog: ui_widget.InputDialog{
				Title: i18n.Interface,
				Hint:  i18n.InterfaceHint,
				Input: component.TextField{
					Editor: widget.Editor{
						SingleLine: true,
						MaxLen:     255,
					},
				},
			},
		*/

		admission:  ui_widget.Selector{Title: i18n.Admission},
		bypass:     ui_widget.Selector{Title: i18n.Bypass},
		resolver:   ui_widget.Selector{Title: i18n.Resolver},
		hostMapper: ui_widget.Selector{Title: i18n.Hosts},
		limiter:    ui_widget.Selector{Title: i18n.Limiter},
		logger:     ui_widget.Selector{Title: i18n.Logger},
		observer:   ui_widget.Selector{Title: i18n.Observer},

		metadataSelector: ui_widget.Selector{Title: i18n.Metadata},
		metadataDialog:   ui_widget.MetadataDialog{},
	}

	p.handler = handler{
		router:           r,
		mode:             &p.mode,
		typ:              ui_widget.Selector{Title: i18n.Type},
		chain:            ui_widget.Selector{Title: i18n.Chain},
		auther:           ui_widget.Selector{Title: i18n.Auther},
		limiter:          ui_widget.Selector{Title: i18n.Limiter},
		observer:         ui_widget.Selector{Title: i18n.Observer},
		metadataSelector: ui_widget.Selector{Title: i18n.Metadata},
		metadataDialog:   ui_widget.MetadataDialog{},
	}
	p.listener = listener{
		router:           r,
		mode:             &p.mode,
		typ:              ui_widget.Selector{Title: i18n.Type},
		chain:            ui_widget.Selector{Title: i18n.Chain},
		auther:           ui_widget.Selector{Title: i18n.Auther},
		enableTLS:        ui_widget.Switcher{Title: i18n.TLS},
		metadataSelector: ui_widget.Selector{Title: i18n.Metadata},
		metadataDialog:   ui_widget.MetadataDialog{},
	}
	p.forwarder = forwarder{
		router: r,
		mode:   &p.mode,
		delDialog: ui_widget.Dialog{
			Title: i18n.DeleteNode,
		},
	}

	return p
}

func (p *servicePage) Init(opts ...page.PageOption) {
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

	cfg := api.GetConfig()
	var service *api.ServiceConfig
	for _, svc := range cfg.Services {
		if svc.Name == p.id {
			service = svc
			break
		}
	}
	if service == nil {
		service = &api.ServiceConfig{}
	}

	p.mode.Value = string(page.BasicMode)

	p.name.Clear()
	p.name.SetText(service.Name)

	p.addr.Clear()
	p.addr.SetText(service.Addr)

	// md := api.NewMetadata(service.Metadata)
	// p.stats.SetValue(md.GetBool("enableStats"))

	// p.customInterface.Select(ui_widget.SelectorItem{Value: service.Interface})

	{
		p.admission.Clear()
		var items []ui_widget.SelectorItem
		if service.Admission != "" {
			items = append(items, ui_widget.SelectorItem{Value: service.Admission})
		}
		for _, v := range service.Admissions {
			items = append(items, ui_widget.SelectorItem{
				Value: v,
			})
		}
		p.admission.Select(items...)
	}

	{
		p.bypass.Clear()
		var items []ui_widget.SelectorItem
		if service.Bypass != "" {
			items = append(items, ui_widget.SelectorItem{Value: service.Bypass})
		}
		for _, v := range service.Bypasses {
			items = append(items, ui_widget.SelectorItem{
				Value: v,
			})
		}
		p.bypass.Select(items...)
	}

	p.resolver.Clear()
	if service.Resolver != "" {
		p.resolver.Select(ui_widget.SelectorItem{Value: service.Resolver})
	}

	p.hostMapper.Clear()
	if service.Hosts != "" {
		p.hostMapper.Select(ui_widget.SelectorItem{Value: service.Hosts})
	}

	p.limiter.Clear()
	if service.Limiter != "" {
		p.limiter.Select(ui_widget.SelectorItem{Value: service.Limiter})
	}

	p.observer.Clear()
	if service.Observer != "" {
		p.observer.Select(ui_widget.SelectorItem{Value: service.Observer})
	}

	{
		p.logger.Clear()
		var items []ui_widget.SelectorItem
		if service.Logger != "" {
			items = append(items, ui_widget.SelectorItem{Value: service.Logger})
		}

		for _, v := range service.Loggers {
			items = append(items, ui_widget.SelectorItem{
				Value: v,
			})
		}
		p.logger.Select(items...)
	}

	{
		p.metadata = nil
		meta := api.NewMetadata(service.Metadata)
		for k := range service.Metadata {
			md := page.Metadata{
				K: k,
				V: meta.GetString(k),
			}
			p.metadata = append(p.metadata, md)
		}
		p.metadataSelector.Clear()
		p.metadataSelector.Select(ui_widget.SelectorItem{Value: strconv.Itoa(len(p.metadata))})
	}

	{
		h := service.Handler
		if h == nil {
			h = &api.HandlerConfig{}
		}

		p.handler.typ.Clear()
		for i := range handlerTypeAdvancedOptions {
			if handlerTypeAdvancedOptions[i].Value == h.Type {
				p.handler.typ.Select(ui_widget.SelectorItem{Name: handlerTypeAdvancedOptions[i].Name, Key: handlerTypeAdvancedOptions[i].Key, Value: handlerTypeAdvancedOptions[i].Value})
				break
			}
		}

		p.handler.chain.Clear()
		p.handler.chain.Select(ui_widget.SelectorItem{Value: h.Chain})

		{
			p.handler.username.Clear()
			p.handler.password.Clear()
			p.handler.authType.Value = ""

			if h.Auth != nil {
				p.handler.username.SetText(h.Auth.Username)
				p.handler.password.SetText(h.Auth.Password)
				p.handler.authType.Value = string(page.AuthSimple)
			}

			p.handler.auther.Clear()
			var items []ui_widget.SelectorItem
			if h.Auther != "" {
				items = append(items, ui_widget.SelectorItem{Value: h.Auther})
			}
			for _, v := range h.Authers {
				items = append(items, ui_widget.SelectorItem{Value: v})
			}
			p.handler.auther.Select(items...)

			if len(h.Authers) > 0 || h.Auther != "" {
				p.handler.authType.Value = string(page.AuthAuther)
			}
		}

		p.handler.metadata = nil
		meta := api.NewMetadata(h.Metadata)
		for k := range h.Metadata {
			md := page.Metadata{
				K: k,
				V: meta.GetString(k),
			}
			p.handler.metadata = append(p.handler.metadata, md)
		}
		p.handler.metadataSelector.Clear()
		p.handler.metadataSelector.Select(ui_widget.SelectorItem{Value: strconv.Itoa(len(p.handler.metadata))})
	}

	{
		ln := service.Listener
		if ln == nil {
			ln = &api.ListenerConfig{}
		}

		p.listener.typ.Clear()
		p.listener.typ.Select(ui_widget.SelectorItem{Value: ln.Type})

		p.listener.chain.Clear()
		p.listener.chain.Select(ui_widget.SelectorItem{Value: ln.Chain})

		{
			p.listener.username.Clear()
			p.listener.password.Clear()
			p.listener.authType.Value = ""

			if ln.Auth != nil {
				p.listener.username.SetText(ln.Auth.Username)
				p.listener.password.SetText(ln.Auth.Password)
				p.listener.authType.Value = string(page.AuthSimple)
			}

			p.listener.auther.Clear()
			var items []ui_widget.SelectorItem
			if ln.Auther != "" {
				items = append(items, ui_widget.SelectorItem{Value: ln.Auther})
			}
			for _, v := range ln.Authers {
				items = append(items, ui_widget.SelectorItem{Value: v})
			}
			p.listener.auther.Select(items...)
			if len(ln.Authers) > 0 || ln.Auther != "" {
				p.listener.authType.Value = string(page.AuthAuther)
			}
		}

		{
			p.listener.enableTLS.SetValue(false)
			p.listener.tlsCertFile.Clear()
			p.listener.tlsKeyFile.Clear()
			p.listener.tlsCAFile.Clear()

			if tls := ln.TLS; tls != nil {
				p.listener.enableTLS.SetValue(true)
				p.listener.tlsCertFile.SetText(tls.CertFile)
				p.listener.tlsKeyFile.SetText(tls.KeyFile)
				p.listener.tlsCAFile.SetText(tls.CAFile)
			}
		}

		p.listener.metadata = nil
		meta := api.NewMetadata(ln.Metadata)
		for k := range ln.Metadata {
			md := page.Metadata{
				K: k,
				V: meta.GetString(k),
			}
			p.listener.metadata = append(p.listener.metadata, md)
		}
		p.listener.metadataSelector.Clear()
		p.listener.metadataSelector.Select(ui_widget.SelectorItem{Value: strconv.Itoa(len(p.listener.metadata))})
	}

	{
		fwd := service.Forwarder
		if fwd == nil {
			fwd = &api.ForwarderConfig{}
		}

		p.forwarder.nodes = nil

		for _, v := range fwd.Nodes {
			if v == nil {
				continue
			}

			nd := node{
				bypass:       ui_widget.Selector{Title: i18n.Bypass},
				enableFilter: ui_widget.Switcher{Title: i18n.Filter},
				protocol:     ui_widget.Selector{Title: i18n.Protocol},
				enableHTTP:   ui_widget.Switcher{Title: i18n.HTTP},
				enableTLS:    ui_widget.Switcher{Title: i18n.TLS},
				tlsSecure:    ui_widget.Switcher{Title: i18n.VerifyServerCert},
				fold:         true,
			}
			nd.name.SetText(v.Name)
			nd.addr.SetText(v.Addr)

			var items []ui_widget.SelectorItem
			if v.Bypass != "" {
				items = append(items, ui_widget.SelectorItem{Value: v.Bypass})
			}
			for _, v := range v.Bypasses {
				items = append(items, ui_widget.SelectorItem{Value: v})
			}
			nd.bypass.Select(items...)

			nd.protocol.Select(ui_widget.SelectorItem{Value: v.Protocol})
			nd.host.SetText(v.Host)
			nd.path.SetText(v.Path)
			if v.Protocol != "" || v.Host != "" || v.Path != "" {
				nd.enableFilter.SetValue(true)
			}

			if v.HTTP != nil {
				nd.enableHTTP.SetValue(true)
				nd.httpHost.SetText(v.HTTP.Host)
				if v.HTTP.Auth != nil {
					nd.httpUsername.SetText(v.HTTP.Auth.Username)
					nd.httpPassword.SetText(v.HTTP.Auth.Password)
				}
			}

			if v.TLS != nil {
				nd.enableTLS.SetValue(true)
				nd.tlsSecure.SetValue(v.TLS.Secure)
				nd.tlsServerName.SetText(v.TLS.ServerName)
			}

			p.forwarder.nodes = append(p.forwarder.nodes, nd)
		}
	}
}

func (p *servicePage) Layout(gtx page.C) page.D {
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
						title := material.H6(th, i18n.Service.Value())
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

func (p *servicePage) layout(gtx page.C, th *page.T) page.D {
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
					return layout.Flex{
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Flexed(1, layout.Spacer{Width: 8}.Layout),
						layout.Rigid(func(gtx page.C) page.D {
							gtx.Source = src
							return material.RadioButton(th, &p.mode, string(page.BasicMode), i18n.Basic.Value()).Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: 8}.Layout),
						layout.Rigid(func(gtx page.C) page.D {
							gtx.Source = src
							return material.RadioButton(th, &p.mode, string(page.AdvancedMode), i18n.Advanced.Value()).Layout(gtx)
						}),
					)
				}),

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

				/*
					layout.Rigid(func(gtx page.C) page.D {
						return p.stats.Layout(gtx, th)
					}),
				*/

				// advanced mode
				layout.Rigid(func(gtx page.C) page.D {
					if p.mode.Value == string(page.BasicMode) {
						return page.D{}
					}

					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(gtx,
						/*
							layout.Rigid(func(gtx page.C) page.D {
								if p.customInterface.Clicked(gtx) {
									p.showInterfaceDialog(gtx)
								}
								return p.customInterface.Layout(gtx, th)
							}),
						*/

						layout.Rigid(func(gtx page.C) page.D {
							if p.admission.Clicked(gtx) {
								p.showAdmissionMenu(gtx)
							}
							return p.admission.Layout(gtx, th)
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

						layout.Rigid(func(gtx page.C) page.D {
							if p.limiter.Clicked(gtx) {
								p.showLimiterMenu(gtx)
							}
							return p.limiter.Layout(gtx, th)
						}),

						layout.Rigid(func(gtx page.C) page.D {
							if p.observer.Clicked(gtx) {
								p.showObserverMenu(gtx)
							}
							return p.observer.Layout(gtx, th)
						}),

						layout.Rigid(func(gtx page.C) page.D {
							if p.logger.Clicked(gtx) {
								p.showLoggerMenu(gtx)
							}
							return p.logger.Layout(gtx, th)
						}),

						layout.Rigid(func(gtx page.C) page.D {
							if p.metadataSelector.Clicked(gtx) {
								p.showMetadataDialog(gtx)
							}
							return p.metadataSelector.Layout(gtx, th)
						}),
						layout.Rigid(layout.Spacer{Height: 8}.Layout),
					)

				}),

				layout.Rigid(func(gtx page.C) page.D {
					return layout.Inset{
						Top:    16,
						Bottom: 16,
					}.Layout(gtx, material.H6(th, i18n.Handler.Value()).Layout)
				}),

				layout.Rigid(func(gtx page.C) page.D {
					return p.handler.Layout(gtx, th)
				}),

				layout.Rigid(func(gtx page.C) page.D {
					return layout.Inset{
						Top:    16,
						Bottom: 16,
					}.Layout(gtx, material.H6(th, i18n.Listener.Value()).Layout)
				}),

				layout.Rigid(func(gtx page.C) page.D {
					return p.listener.Layout(gtx, th)
				}),

				layout.Rigid(func(gtx page.C) page.D {
					if !p.handler.canForward() {
						return page.D{}
					}

					return layout.Inset{
						Top:    16,
						Bottom: 16,
					}.Layout(gtx, material.H6(th, i18n.Forwarder.Value()).Layout)
				}),

				layout.Rigid(func(gtx page.C) page.D {
					if !p.handler.canForward() {
						return page.D{}
					}

					return p.forwarder.Layout(gtx, th)
				}),
			)
		})
	})
}

/*
func (p *servicePage) showInterfaceDialog(gtx  page.C)  {
	p.interfaceDialog.Input.Clear()
	p.interfaceDialog.Input.SetText(p.customInterface.Value())
	p.interfaceDialog.OnClick = func(ok bool) {
		p.modal.Disappear(gtx.Now)

		if ok {
			p.customInterface.Clear()
			p.customInterface.Select(ui_widget.SelectorItem{Value: strings.TrimSpace(p.interfaceDialog.Input.Text())})
		}
	}
	p.modal.Widget = func(gtx page.C, th *material.Theme, anim *component.VisibilityAnimation) page.D {
		return p.interfaceDialog.Layout(gtx, th)
	}
	p.modal.Appear(gtx.Now)
}
*/

func (p *servicePage) showAdmissionMenu(gtx page.C) {
	options := []ui_widget.MenuOption{}
	for _, v := range api.GetConfig().Admissions {
		options = append(options, ui_widget.MenuOption{
			Value: v.Name,
		})
	}
	for i := range options {
		options[i].Selected = p.admission.AnyValue(options[i].Value)
	}

	p.menu.Title = i18n.Admission
	p.menu.Options = options
	p.menu.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		p.admission.Clear()
		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				p.admission.Select(ui_widget.SelectorItem{Value: p.menu.Options[i].Value})
			}
		}
	}
	p.menu.ShowAdd = true
	p.menu.Multiple = true

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *servicePage) showBypassMenu(gtx page.C) {
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
	p.menu.ShowAdd = true
	p.menu.Multiple = true

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *servicePage) showResolverMenu(gtx page.C) {
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
	p.menu.ShowAdd = true
	p.menu.Multiple = false

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *servicePage) showHostMapperMenu(gtx page.C) {
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
	p.menu.ShowAdd = true
	p.menu.Multiple = false

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *servicePage) showLimiterMenu(gtx page.C) {
	options := []ui_widget.MenuOption{}
	for _, v := range api.GetConfig().Limiters {
		options = append(options, ui_widget.MenuOption{
			Value: v.Name,
		})
	}
	for i := range options {
		options[i].Selected = p.limiter.AnyValue(options[i].Value)
	}

	p.menu.Title = i18n.Limiter
	p.menu.Options = options
	p.menu.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		p.limiter.Clear()
		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				p.limiter.Select(ui_widget.SelectorItem{Value: p.menu.Options[i].Value})
			}
		}
	}
	p.menu.ShowAdd = true
	p.menu.Multiple = false

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *servicePage) showObserverMenu(gtx page.C) {
	options := []ui_widget.MenuOption{}
	for _, v := range api.GetConfig().Observers {
		options = append(options, ui_widget.MenuOption{
			Value: v.Name,
		})
	}
	for i := range options {
		options[i].Selected = p.observer.AnyValue(options[i].Value)
	}

	p.menu.Title = i18n.Observer
	p.menu.Options = options
	p.menu.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		p.observer.Clear()
		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				p.observer.Select(ui_widget.SelectorItem{Value: p.menu.Options[i].Value})
			}
		}
	}
	p.menu.ShowAdd = true
	p.menu.Multiple = false

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *servicePage) showLoggerMenu(gtx page.C) {
	options := []ui_widget.MenuOption{}
	for _, v := range api.GetConfig().Loggers {
		options = append(options, ui_widget.MenuOption{
			Value: v.Name,
		})
	}
	for i := range options {
		options[i].Selected = p.logger.AnyValue(options[i].Value)
	}

	p.menu.Title = i18n.Logger
	p.menu.Options = options
	p.menu.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		p.logger.Clear()
		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				p.logger.Select(ui_widget.SelectorItem{Value: p.menu.Options[i].Value})
			}
		}
	}
	p.menu.ShowAdd = true
	p.menu.Multiple = true

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *servicePage) showMetadataDialog(gtx page.C) {
	p.metadataDialog.Clear()
	for _, md := range p.metadata {
		p.metadataDialog.Add(md.K, md.V)
	}
	p.metadataDialog.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}
		p.metadata = nil
		for _, kv := range p.metadataDialog.Metadata() {
			k, v := kv.Get()
			p.metadata = append(p.metadata, page.Metadata{
				K: k,
				V: v,
			})
		}
		p.metadataSelector.Clear()
		p.metadataSelector.Select(ui_widget.SelectorItem{Value: strconv.Itoa(len(p.metadata))})
	}

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.metadataDialog.Layout(gtx, th)
	})
}

func (p *servicePage) save() bool {
	cfg := p.generateConfig()

	var err error
	if p.id == "" {
		err = runner.Exec(context.Background(),
			task.CreateService(cfg),
			runner.WithCancel(true),
		)
	} else {
		err = runner.Exec(context.Background(),
			task.UpdateService(cfg),
			runner.WithCancel(true),
		)
	}
	util.RestartGetConfigTask()

	return err == nil
}

func (p *servicePage) generateConfig() *api.ServiceConfig {
	var svcCfg *api.ServiceConfig

	if p.id != "" {
		for _, svc := range api.GetConfig().Services {
			if svc == nil {
				continue
			}
			if svc.Name == p.id {
				svcCfg = svc.Copy()
				break
			}
		}
	}

	if svcCfg == nil {
		svcCfg = &api.ServiceConfig{}
	}

	svcCfg.Name = p.name.Text()
	svcCfg.Addr = p.addr.Text()
	// svcCfg.Interface = p.customInterface.Value()

	svcCfg.Admission = ""
	svcCfg.Admissions = p.admission.Values()

	svcCfg.Bypass = ""
	svcCfg.Bypasses = p.bypass.Values()

	svcCfg.Resolver = p.resolver.Value()
	svcCfg.Hosts = p.hostMapper.Value()
	svcCfg.Limiter = p.limiter.Value()

	svcCfg.Logger = ""
	svcCfg.Loggers = p.logger.Values()
	svcCfg.Observer = p.observer.Value()

	if svcCfg.Metadata == nil {
		svcCfg.Metadata = make(map[string]any)
	}
	svcCfg.Metadata["enablestats"] = true

	if svcCfg.Handler == nil {
		svcCfg.Handler = &api.HandlerConfig{}
	}

	svcCfg.Handler.Type = p.handler.typ.Value()
	svcCfg.Handler.Chain = p.handler.chain.Value()

	svcCfg.Handler.Auther = ""
	svcCfg.Handler.Authers = nil
	svcCfg.Handler.Auth = nil
	if p.handler.authType.Value == string(page.AuthAuther) {
		svcCfg.Handler.Authers = p.handler.auther.Values()
	}
	if p.handler.authType.Value == string(page.AuthSimple) {
		username := strings.TrimSpace(p.handler.username.Text())
		password := strings.TrimSpace(p.handler.password.Text())
		if username != "" {
			svcCfg.Handler.Auth = &api.AuthConfig{
				Username: username,
				Password: password,
			}
		}
	}

	svcCfg.Limiter = p.handler.limiter.Value()
	svcCfg.Observer = p.handler.observer.Value()

	svcCfg.Handler.Metadata = nil
	if len(p.handler.metadata) > 0 {
		svcCfg.Handler.Metadata = make(map[string]any)
	}
	for _, md := range p.handler.metadata {
		svcCfg.Handler.Metadata[md.K] = md.V
	}

	if svcCfg.Listener == nil {
		svcCfg.Listener = &api.ListenerConfig{}
	}

	svcCfg.Listener.Type = p.listener.typ.Value()
	svcCfg.Listener.Chain = p.listener.chain.Value()

	svcCfg.Listener.Auther = ""
	svcCfg.Listener.Authers = nil
	svcCfg.Listener.Auth = nil
	if p.listener.authType.Value == string(page.AuthAuther) {
		svcCfg.Listener.Authers = p.listener.auther.Values()
	}
	if p.listener.authType.Value == string(page.AuthSimple) {
		username := strings.TrimSpace(p.listener.username.Text())
		password := strings.TrimSpace(p.listener.password.Text())
		if username != "" {
			svcCfg.Listener.Auth = &api.AuthConfig{
				Username: username,
				Password: password,
			}
		}
	}

	svcCfg.Listener.TLS = nil
	if p.listener.enableTLS.Value() {
		svcCfg.Listener.TLS = &api.TLSConfig{
			CertFile: strings.TrimSpace(p.listener.tlsCertFile.Text()),
			KeyFile:  strings.TrimSpace(p.listener.tlsKeyFile.Text()),
			CAFile:   strings.TrimSpace(p.listener.tlsCAFile.Text()),
		}
	}

	svcCfg.Listener.Metadata = nil
	if len(p.listener.metadata) > 0 {
		svcCfg.Listener.Metadata = make(map[string]any)
	}
	for _, md := range p.listener.metadata {
		svcCfg.Listener.Metadata[md.K] = md.V
	}

	svcCfg.Forwarder = nil
	if len(p.forwarder.nodes) > 0 {
		svcCfg.Forwarder = &api.ForwarderConfig{}

		for _, node := range p.forwarder.nodes {
			nodeCfg := &api.ForwardNodeConfig{
				Name:     node.name.Text(),
				Addr:     node.addr.Text(),
				Bypasses: node.bypass.Values(),
			}
			if node.enableFilter.Value() {
				nodeCfg.Host = node.host.Text()
				nodeCfg.Protocol = node.protocol.Value()
				nodeCfg.Path = node.path.Text()
			}
			if node.enableHTTP.Value() {
				nodeCfg.HTTP = &api.HTTPNodeConfig{
					Host: node.httpHost.Text(),
				}
				username := strings.TrimSpace(node.httpUsername.Text())
				password := strings.TrimSpace(node.httpPassword.Text())
				if username != "" {
					nodeCfg.HTTP.Auth = &api.AuthConfig{
						Username: username,
						Password: password,
					}
				}
			}
			if node.enableTLS.Value() {
				nodeCfg.TLS = &api.TLSNodeConfig{
					Secure:     node.tlsSecure.Value(),
					ServerName: node.tlsServerName.Text(),
				}
			}
			svcCfg.Forwarder.Nodes = append(svcCfg.Forwarder.Nodes, nodeCfg)
		}
	}

	return svcCfg
}

func (p *servicePage) delete() {
	runner.Exec(context.Background(),
		task.DeleteService(p.id),
		runner.WithCancel(true),
	)
	util.RestartGetConfigTask()
}
