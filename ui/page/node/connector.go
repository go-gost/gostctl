package node

import (
	"strconv"

	"gioui.org/layout"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/go-gost/gostctl/api"
	"github.com/go-gost/gostctl/ui/i18n"
	"github.com/go-gost/gostctl/ui/page"
	ui_widget "github.com/go-gost/gostctl/ui/widget"
	"github.com/google/uuid"
)

type connector struct {
	node *nodePage
	menu ui_widget.Menu

	typ ui_widget.Selector

	enableAuth ui_widget.Switcher
	username   component.TextField
	password   component.TextField

	metadata   api.Metadata
	mdSelector ui_widget.Selector
}

func newConnector(node *nodePage) *connector {
	return &connector{
		node:       node,
		typ:        ui_widget.Selector{Title: i18n.Type},
		enableAuth: ui_widget.Switcher{Title: i18n.Auth},
		mdSelector: ui_widget.Selector{Title: i18n.Metadata},
	}

}

func (p *connector) init(cfg *api.ConnectorConfig) {
	if cfg == nil {
		cfg = &api.ConnectorConfig{}
	}

	p.typ.Clear()
	for i := range connectorTypeAdvancedOptions {
		if connectorTypeAdvancedOptions[i].Value == cfg.Type {
			p.typ.Select(ui_widget.SelectorItem{Name: connectorTypeAdvancedOptions[i].Name, Key: connectorTypeAdvancedOptions[i].Key, Value: connectorTypeAdvancedOptions[i].Value})
			break
		}
	}

	{
		p.enableAuth.SetValue(false)
		p.username.Clear()
		p.password.Clear()

		if cfg.Auth != nil {
			p.enableAuth.SetValue(true)
			p.username.SetText(cfg.Auth.Username)
			p.password.SetText(cfg.Auth.Password)
		}
	}

	p.metadata = api.NewMetadata(cfg.Metadata)
	p.mdSelector.Clear()
	p.mdSelector.Select(ui_widget.SelectorItem{Value: strconv.Itoa(len(p.metadata))})
}

func (p *connector) Layout(gtx page.C, th *page.T) page.D {
	src := gtx.Source

	if !p.node.edit {
		gtx = gtx.Disabled()
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx page.C) page.D {
			if p.typ.Clicked(gtx) {
				p.showTypeMenu(gtx)
			}

			return p.typ.Layout(gtx, th)
		}),

		layout.Rigid(func(gtx page.C) page.D {
			if !p.canAuth() {
				return page.D{}
			}

			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx page.C) page.D {
					return p.enableAuth.Layout(gtx, th)
				}),
				layout.Rigid(func(gtx page.C) page.D {
					if !p.enableAuth.Value() {
						return page.D{}
					}
					return layout.UniformInset(8).Layout(gtx, func(gtx page.C) page.D {
						return layout.Flex{
							Axis: layout.Vertical,
						}.Layout(gtx,
							layout.Rigid(func(gtx page.C) page.D {
								return material.Body1(th, i18n.Username.Value()).Layout(gtx)
							}),
							layout.Rigid(func(gtx page.C) page.D {
								return p.username.Layout(gtx, th, "")
							}),
							layout.Rigid(layout.Spacer{Height: 8}.Layout),

							layout.Rigid(func(gtx page.C) page.D {
								return material.Body1(th, i18n.Password.Value()).Layout(gtx)
							}),
							layout.Rigid(func(gtx page.C) page.D {
								return p.password.Layout(gtx, th, "")
							}),
							layout.Rigid(layout.Spacer{Height: 8}.Layout),
						)
					})
				}),
			)
		}),

		layout.Rigid(func(gtx page.C) page.D {
			gtx.Source = src

			if p.mdSelector.Clicked(gtx) {
				perm := page.PermRead
				if p.node.edit {
					perm = page.PermWrite | page.PermDelete
				}
				p.node.router.Goto(page.Route{
					Path:     page.PageMetadata,
					ID:       uuid.New().String(),
					Value:    p.metadata,
					Callback: p.mdCallback,
					Perm:     perm,
				})
			}

			return p.mdSelector.Layout(gtx, th)
		}),
	)
}

var (
	connectorTypeOptions = []ui_widget.MenuOption{
		{Name: "HTTP", Value: "http"},
		{Name: "SOCKS4", Value: "socks4"},
		{Name: "SOCKS5", Value: "socks5"},
		{Name: "Relay", Value: "relay"},
		{Name: "Shadowsocks", Value: "ss"},
		{Name: "HTTP/2", Value: "http2"},
	}

	connectorTypeAdvancedOptions = []ui_widget.MenuOption{
		{Name: "HTTP", Value: "http"},
		{Name: "SOCKS4", Value: "socks4"},
		{Name: "SOCKS5", Value: "socks5"},
		{Name: "Relay", Value: "relay"},
		{Name: "Shadowsocks", Value: "ss"},
		{Name: "HTTP/2", Value: "http2"},

		{Name: "SNI", Value: "sni"},
		{Name: "SSHD", Value: "sshd"},
		{Key: i18n.ReverseProxyTunnel, Value: "tunnel"},
		{Key: i18n.SerialPortRedirector, Value: "serial"},
		{Key: i18n.UnixDomainSocket, Value: "unix"},
	}
)

func (p *connector) canAuth() bool {
	return p.typ.AnyValue("http", "http2", "socks4", "socks5", "socks", "relay", "ss")
}

func (p *connector) showTypeMenu(gtx page.C) {
	options := connectorTypeOptions
	if p.node.mode.Value == string(page.AdvancedMode) {
		options = connectorTypeAdvancedOptions
	}

	for i := range options {
		options[i].Selected = p.typ.AnyValue(options[i].Value)
	}

	p.menu.Title = i18n.Connector
	p.menu.Options = options
	p.menu.OnClick = func(ok bool) {
		p.node.router.HideModal(gtx)
		if !ok {
			return
		}

		p.typ.Clear()
		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				p.typ.Select(ui_widget.SelectorItem{Name: options[i].Name, Key: options[i].Key, Value: options[i].Value})
			}
		}
	}
	p.menu.ShowAdd = false
	p.menu.Multiple = false

	p.node.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *connector) mdCallback(action page.Action, id string, value any) {
	if id == "" {
		return
	}

	switch action {
	case page.ActionUpdate:
		p.metadata, _ = value.(api.Metadata)

	case page.ActionDelete:
		p.metadata = nil
	}

	p.mdSelector.Clear()
	p.mdSelector.Select(ui_widget.SelectorItem{Value: strconv.Itoa(len(p.metadata))})
}
