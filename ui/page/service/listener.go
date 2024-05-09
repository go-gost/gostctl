package service

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

type listener struct {
	service *servicePage
	menu    ui_widget.Menu

	typ   ui_widget.Selector
	chain ui_widget.Selector

	authType widget.Enum
	username component.TextField
	password component.TextField
	auther   ui_widget.Selector

	enableTLS   ui_widget.Switcher
	tlsCertFile component.TextField
	tlsKeyFile  component.TextField
	tlsCAFile   component.TextField

	metadata         []page.Metadata
	metadataSelector ui_widget.Selector
	metadataDialog   ui_widget.MetadataDialog
}

func newListener(service *servicePage) *listener {
	return &listener{
		service:   service,
		typ:       ui_widget.Selector{Title: i18n.Type},
		chain:     ui_widget.Selector{Title: i18n.Chain},
		auther:    ui_widget.Selector{Title: i18n.Auther},
		enableTLS: ui_widget.Switcher{Title: i18n.TLS},
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

func (l *listener) init(cfg *api.ListenerConfig) {
	if cfg == nil {
		cfg = &api.ListenerConfig{}
	}

	l.typ.Clear()
	for i := range listenerTypeAdvancedOptions {
		if listenerTypeAdvancedOptions[i].Value == cfg.Type {
			l.typ.Select(ui_widget.SelectorItem{Name: listenerTypeAdvancedOptions[i].Name, Key: listenerTypeAdvancedOptions[i].Key, Value: listenerTypeAdvancedOptions[i].Value})
			break
		}
	}

	l.chain.Clear()
	l.chain.Select(ui_widget.SelectorItem{Value: cfg.Chain})

	{
		l.username.Clear()
		l.password.Clear()
		l.authType.Value = ""

		if cfg.Auth != nil {
			l.username.SetText(cfg.Auth.Username)
			l.password.SetText(cfg.Auth.Password)
			l.authType.Value = string(page.AuthSimple)
		}

		l.auther.Clear()
		var items []ui_widget.SelectorItem
		if cfg.Auther != "" {
			items = append(items, ui_widget.SelectorItem{Value: cfg.Auther})
		}
		for _, v := range cfg.Authers {
			items = append(items, ui_widget.SelectorItem{Value: v})
		}
		l.auther.Select(items...)
		if len(cfg.Authers) > 0 || cfg.Auther != "" {
			l.authType.Value = string(page.AuthAuther)
		}
	}

	{
		l.enableTLS.SetValue(false)
		l.tlsCertFile.Clear()
		l.tlsKeyFile.Clear()
		l.tlsCAFile.Clear()

		if tls := cfg.TLS; tls != nil {
			l.enableTLS.SetValue(true)
			l.tlsCertFile.SetText(tls.CertFile)
			l.tlsKeyFile.SetText(tls.KeyFile)
			l.tlsCAFile.SetText(tls.CAFile)
		}
	}

	l.metadata = nil
	meta := api.NewMetadata(cfg.Metadata)
	for k := range cfg.Metadata {
		md := page.Metadata{
			K: k,
			V: meta.GetString(k),
		}
		l.metadata = append(l.metadata, md)
	}
	l.metadataSelector.Clear()
	l.metadataSelector.Select(ui_widget.SelectorItem{Value: strconv.Itoa(len(l.metadata))})
}

func (l *listener) Layout(gtx page.C, th *page.T) page.D {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx page.C) page.D {
			if l.typ.Clicked(gtx) {
				l.showTypeMenu(gtx)
			}
			return l.typ.Layout(gtx, th)
		}),

		layout.Rigid(func(gtx page.C) page.D {
			if !l.canChain() {
				return page.D{}
			}

			if l.chain.Clicked(gtx) {
				l.showChainMenu(gtx)
			}
			return l.chain.Layout(gtx, th)
		}),

		// Auth Config
		layout.Rigid(func(gtx page.C) page.D {
			if !l.canAuth() {
				return page.D{}
			}

			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx page.C) page.D {
					return layout.Inset{
						Top:    4,
						Bottom: 4,
					}.Layout(gtx, func(gtx page.C) page.D {
						return layout.Flex{
							Alignment: layout.Middle,
							Spacing:   layout.SpaceBetween,
						}.Layout(gtx,
							layout.Flexed(1, func(gtx page.C) page.D {
								return material.Body1(th, i18n.Auth.Value()).Layout(gtx)
							}),
							layout.Rigid(func(gtx page.C) page.D {
								return material.RadioButton(th, &l.authType, string(page.AuthSimple), i18n.AuthSimple.Value()).Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Width: 4}.Layout),
							layout.Rigid(func(gtx page.C) page.D {
								return material.RadioButton(th, &l.authType, string(page.AuthAuther), i18n.AuthAuther.Value()).Layout(gtx)
							}),
						)
					})
				}),

				layout.Rigid(func(gtx page.C) page.D {
					if l.authType.Value != string(page.AuthSimple) {
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
								return l.username.Layout(gtx, th, "")
							}),
							layout.Rigid(layout.Spacer{Height: 8}.Layout),

							layout.Rigid(func(gtx page.C) page.D {
								return material.Body1(th, i18n.Password.Value()).Layout(gtx)
							}),
							layout.Rigid(func(gtx page.C) page.D {
								return l.password.Layout(gtx, th, "")
							}),
							layout.Rigid(layout.Spacer{Height: 8}.Layout),
						)
					})
				}),

				layout.Rigid(func(gtx page.C) page.D {
					if l.authType.Value != string(page.AuthAuther) {
						return page.D{}
					}

					if l.auther.Clicked(gtx) {
						l.showAutherMenu(gtx)
					}
					return layout.UniformInset(8).Layout(gtx, func(gtx page.C) page.D {
						return l.auther.Layout(gtx, th)
					})
				}),
			)
		}),

		// TLS config
		layout.Rigid(func(gtx page.C) page.D {
			if !l.canTLS() {
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
						return l.enableTLS.Layout(gtx, th)
					}),

					layout.Rigid(func(gtx page.C) page.D {
						if !l.enableTLS.Value() {
							return page.D{}
						}
						return layout.UniformInset(8).Layout(gtx, func(gtx page.C) page.D {
							return layout.Flex{
								Axis: layout.Vertical,
							}.Layout(gtx,
								layout.Rigid(func(gtx page.C) page.D {
									return material.Body1(th, i18n.CertFile.Value()).Layout(gtx)
								}),
								layout.Rigid(func(gtx page.C) page.D {
									return l.tlsCertFile.Layout(gtx, th, "")
								}),
								layout.Rigid(layout.Spacer{Height: 8}.Layout),

								layout.Rigid(func(gtx page.C) page.D {
									return material.Body1(th, i18n.KeyFile.Value()).Layout(gtx)
								}),
								layout.Rigid(func(gtx page.C) page.D {
									return l.tlsKeyFile.Layout(gtx, th, "")
								}),
								layout.Rigid(layout.Spacer{Height: 8}.Layout),

								layout.Rigid(func(gtx page.C) page.D {
									return material.Body1(th, i18n.CAFile.Value()).Layout(gtx)
								}),
								layout.Rigid(func(gtx page.C) page.D {
									return l.tlsCAFile.Layout(gtx, th, "")
								}),
							)
						})
					}),
				)
			})
		}),

		layout.Rigid(func(gtx page.C) page.D {
			if l.metadataSelector.Clicked(gtx) {
				l.showMetadataDialog(gtx)
			}
			return l.metadataSelector.Layout(gtx, th)
		}),
	)
}

var (
	listenerTypeOptions = []ui_widget.MenuOption{
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
	}
	listenerTypeAdvancedOptions = []ui_widget.MenuOption{
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

func (l *listener) canChain() bool {
	return l.typ.AnyValue("rtcp", "rudp")
}

func (l *listener) canAuth() bool {
	return l.typ.AnyValue("ssh", "sshd")
}

func (l *listener) canTLS() bool {
	return l.typ.AnyValue("tls", "mtls", "wss", "mwss", "http2", "h2", "grpc", "quic", "http3", "h3", "wt", "dtls")
}

func (l *listener) showTypeMenu(gtx page.C) {
	options := listenerTypeOptions
	if l.service.mode.Value == string(page.AdvancedMode) {
		options = listenerTypeAdvancedOptions
	}

	for i := range options {
		options[i].Selected = l.typ.AnyValue(options[i].Value)
	}

	l.menu.Title = i18n.Listener
	l.menu.Options = options
	l.menu.OnClick = func(ok bool) {
		l.service.router.HideModal(gtx)
		if !ok {
			return
		}

		l.typ.Clear()
		for i := range l.menu.Options {
			if l.menu.Options[i].Selected {
				l.typ.Select(ui_widget.SelectorItem{Name: l.menu.Options[i].Name, Key: l.menu.Options[i].Key, Value: l.menu.Options[i].Value})
			}
		}
	}
	l.menu.ShowAdd = false
	l.menu.Multiple = false

	l.service.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return l.menu.Layout(gtx, th)
	})
}

func (l *listener) showChainMenu(gtx page.C) {
	options := []ui_widget.MenuOption{}
	for _, v := range api.GetConfig().Chains {
		options = append(options, ui_widget.MenuOption{
			Value: v.Name,
		})
	}
	for i := range options {
		options[i].Selected = l.chain.AnyValue(options[i].Value)
	}

	l.menu.Title = i18n.Chain
	l.menu.Options = options
	l.menu.OnClick = func(ok bool) {
		l.service.router.HideModal(gtx)
		if !ok {
			return
		}

		l.chain.Clear()
		for i := range l.menu.Options {
			if l.menu.Options[i].Selected {
				l.chain.Select(ui_widget.SelectorItem{Value: l.menu.Options[i].Value})
			}
		}
	}
	l.menu.ShowAdd = true
	l.menu.Multiple = false

	l.service.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return l.menu.Layout(gtx, th)
	})
}

func (l *listener) showAutherMenu(gtx page.C) {
	options := []ui_widget.MenuOption{}
	for _, v := range api.GetConfig().Authers {
		options = append(options, ui_widget.MenuOption{
			Value: v.Name,
		})
	}
	for i := range options {
		options[i].Selected = l.auther.AnyValue(options[i].Value)
	}

	l.menu.Title = i18n.Auther
	l.menu.Options = options
	l.menu.OnClick = func(ok bool) {
		l.service.router.HideModal(gtx)
		if !ok {
			return
		}

		l.auther.Clear()
		for i := range l.menu.Options {
			if l.menu.Options[i].Selected {
				l.auther.Select(ui_widget.SelectorItem{Value: l.menu.Options[i].Value})
			}
		}
	}
	l.menu.ShowAdd = true
	l.menu.Multiple = true

	l.service.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return l.menu.Layout(gtx, th)
	})
}

func (l *listener) showMetadataDialog(gtx page.C) {
	l.metadataDialog.Clear()
	for _, md := range l.metadata {
		l.metadataDialog.Add(md.K, md.V)
	}
	l.metadataDialog.OnClick = func(ok bool) {
		l.service.router.HideModal(gtx)
		if !ok {
			return
		}
		l.metadata = nil
		for _, kv := range l.metadataDialog.Metadata() {
			k, v := kv.Get()
			l.metadata = append(l.metadata, page.Metadata{
				K: k,
				V: v,
			})
		}
		l.metadataSelector.Clear()
		l.metadataSelector.Select(ui_widget.SelectorItem{Value: strconv.Itoa(len(l.metadata))})
	}

	l.service.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return l.metadataDialog.Layout(gtx, th)
	})
}
