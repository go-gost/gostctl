package node

import (
	"strconv"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/go-gost/gostctl/api"
	"github.com/go-gost/gostctl/ui/i18n"
	"github.com/go-gost/gostctl/ui/page"
	ui_widget "github.com/go-gost/gostctl/ui/widget"
)

type dialer struct {
	node *nodePage
	menu ui_widget.Menu

	typ ui_widget.Selector

	enableAuth ui_widget.Switcher
	username   component.TextField
	password   component.TextField

	enableTLS     ui_widget.Switcher
	tlsSecure     ui_widget.Switcher
	tlsServerName component.TextField
	tlsCertFile   component.TextField
	tlsKeyFile    component.TextField
	tlsCAFile     component.TextField

	metadata         []page.Metadata
	metadataSelector ui_widget.Selector
	metadataDialog   ui_widget.MetadataDialog
}

func newDialer(node *nodePage) *dialer {
	return &dialer{
		node:       node,
		typ:        ui_widget.Selector{Title: i18n.Type},
		enableAuth: ui_widget.Switcher{Title: i18n.Auth},
		enableTLS:  ui_widget.Switcher{Title: i18n.TLS},
		tlsSecure:  ui_widget.Switcher{Title: i18n.VerifyServerCert},
		tlsServerName: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     255,
			},
		},
		tlsCertFile: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     255,
			},
		},
		tlsKeyFile: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     255,
			},
		},
		tlsCAFile: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     255,
			},
		},
		metadataSelector: ui_widget.Selector{Title: i18n.Metadata},
		metadataDialog:   ui_widget.MetadataDialog{},
	}

}

func (p *dialer) init(cfg *api.DialerConfig) {
	if cfg == nil {
		cfg = &api.DialerConfig{}
	}

	p.typ.Clear()
	for i := range dialerTypeAdvancedOptions {
		if dialerTypeAdvancedOptions[i].Value == cfg.Type {
			p.typ.Select(ui_widget.SelectorItem{Name: dialerTypeAdvancedOptions[i].Name, Key: dialerTypeAdvancedOptions[i].Key, Value: dialerTypeAdvancedOptions[i].Value})
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

	{
		p.enableTLS.SetValue(false)
		p.tlsCertFile.Clear()
		p.tlsKeyFile.Clear()
		p.tlsCAFile.Clear()
		p.tlsSecure.SetValue(false)
		p.tlsServerName.Clear()

		if tls := cfg.TLS; tls != nil {
			p.enableTLS.SetValue(true)
			p.tlsCertFile.SetText(tls.CertFile)
			p.tlsKeyFile.SetText(tls.KeyFile)
			p.tlsCAFile.SetText(tls.CAFile)

			p.tlsSecure.SetValue(tls.Secure)
			p.tlsServerName.SetText(tls.ServerName)
		}
	}

	p.metadata = nil
	meta := api.NewMetadata(cfg.Metadata)
	for k := range cfg.Metadata {
		md := page.Metadata{
			K: k,
			V: meta.GetString(k),
		}
		p.metadata = append(p.metadata, md)
	}
	p.metadataSelector.Clear()
	p.metadataSelector.Select(ui_widget.SelectorItem{Value: strconv.Itoa(len(p.metadata))})
}

func (p *dialer) Layout(gtx page.C, th *page.T) page.D {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx page.C) page.D {
			if p.typ.Clicked(gtx) {
				p.showTypeMenu(gtx)
			}

			return p.typ.Layout(gtx, th)
		}),

		// Auth Config
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

		// TLS config
		layout.Rigid(func(gtx page.C) page.D {
			if !p.canTLS() {
				return page.D{}
			}

			return layout.Inset{
				Top:    4,
				Bottom: 4,
			}.Layout(gtx, func(gtx page.C) page.D {
				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
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
								layout.Rigid(layout.Spacer{Height: 8}.Layout),

								layout.Rigid(func(gtx page.C) page.D {
									return material.Body1(th, i18n.CertFile.Value()).Layout(gtx)
								}),
								layout.Rigid(func(gtx page.C) page.D {
									return p.tlsCertFile.Layout(gtx, th, "")
								}),
								layout.Rigid(layout.Spacer{Height: 8}.Layout),

								layout.Rigid(func(gtx page.C) page.D {
									return material.Body1(th, i18n.KeyFile.Value()).Layout(gtx)
								}),
								layout.Rigid(func(gtx page.C) page.D {
									return p.tlsKeyFile.Layout(gtx, th, "")
								}),
								layout.Rigid(layout.Spacer{Height: 8}.Layout),

								layout.Rigid(func(gtx page.C) page.D {
									return material.Body1(th, i18n.CAFile.Value()).Layout(gtx)
								}),
								layout.Rigid(func(gtx page.C) page.D {
									return p.tlsCAFile.Layout(gtx, th, "")
								}),
							)
						})
					}),
				)
			})
		}),
		layout.Rigid(func(gtx page.C) page.D {
			if p.metadataSelector.Clicked(gtx) {
				p.showMetadataDialog(gtx)
			}
			return p.metadataSelector.Layout(gtx, th)
		}),
	)
}

var (
	dialerTypeOptions = []ui_widget.MenuOption{
		{Name: "TCP", Value: "tcp"},
		{Name: "UDP", Value: "udp"},
		{Name: "TLS", Value: "tls"},
		{Name: "MTLS", Value: "mtls"},
		{Name: "WS", Value: "ws"},
		{Name: "MWS", Value: "mws"},
		{Name: "WSS", Value: "wss"},
		{Name: "MWSS", Value: "mwss"},
		{Name: "HTTP/2", Value: "http2"},
		{Name: "gRPC", Value: "grpc"},
		{Name: "QUIC", Value: "quic"},
		{Name: "KCP", Value: "kcp"},
	}
	dialerTypeAdvancedOptions = []ui_widget.MenuOption{
		{Name: "TCP", Value: "tcp"},
		{Name: "UDP", Value: "udp"},
		{Name: "RTCP", Value: "rtcp"},
		{Name: "RUDP", Value: "rudp"},
		{Name: "TLS", Value: "tls"},
		{Name: "MTLS", Value: "mtls"},
		{Name: "WS", Value: "ws"},
		{Name: "MWS", Value: "mws"},
		{Name: "WSS", Value: "wss"},
		{Name: "MWSS", Value: "mwss"},
		{Name: "HTTP/2", Value: "http2"},
		{Name: "gRPC", Value: "grpc"},
		{Name: "QUIC", Value: "quic"},
		{Name: "KCP", Value: "kcp"},
		{Name: "H2", Value: "h2"},
		{Name: "H2C", Value: "h2c"},

		{Name: "WebTransport", Value: "wt"},
		{Name: "DTLS", Value: "dtls"},
		{Name: "MTCP", Value: "mtcp"},

		{Name: "SSH", Value: "ssh"},
		{Name: "SSHD", Value: "sshd"},
		{Name: "DNS", Value: "dns"},

		{Name: "TCP Redirector", Value: "red"},
		{Name: "UDP Redirector", Value: "redu"},
		{Name: "TUN", Value: "tun"},
		{Name: "TAP", Value: "tap"},

		{Name: "PHT", Value: "pht"},
		{Name: "Obfs-HTTP", Value: "ohttp"},
		{Name: "Obfs-TLS", Value: "otls"},

		{Name: "HTTP3", Value: "http3"},
		{Name: "Fake TCP", Value: "ftcp"},
		{Name: "ICMP", Value: "icmp"},

		{Key: i18n.SerialPortRedirector, Value: "serial"},
		{Key: i18n.UnixDomainSocket, Value: "unix"},
	}
)

func (p *dialer) canAuth() bool {
	return p.typ.AnyValue("ssh", "sshd")
}

func (p *dialer) canTLS() bool {
	return p.typ.AnyValue("tls", "mtls", "wss", "mwss", "http2", "h2", "grpc", "quic", "http3", "h3", "wt", "dtls")
}

func (p *dialer) showTypeMenu(gtx page.C) {
	options := dialerTypeOptions
	if p.node.mode.Value == string(page.AdvancedMode) {
		options = dialerTypeAdvancedOptions
	}

	for i := range options {
		options[i].Selected = p.typ.AnyValue(options[i].Value)
	}

	p.menu.Title = i18n.Dialer
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

func (p *dialer) showMetadataDialog(gtx page.C) {
	p.metadataDialog.Clear()
	for _, md := range p.metadata {
		p.metadataDialog.Add(md.K, md.V)
	}
	p.metadataDialog.OnClick = func(ok bool) {
		p.node.router.HideModal(gtx)
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

	p.node.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.metadataDialog.Layout(gtx, th)
	})
}
