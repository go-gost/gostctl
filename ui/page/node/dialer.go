package node

import (
	"strconv"
	"strings"

	"gioui.org/font"
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

	metadata          []metadata
	mdSelector        ui_widget.Selector
	mdFolded          bool
	mdAdd             widget.Clickable
	mdDialog          ui_widget.MetadataDialog
	delMetadataDialog ui_widget.Dialog
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
		mdSelector: ui_widget.Selector{Title: i18n.Metadata},
		mdDialog: ui_widget.MetadataDialog{
			K: component.TextField{
				Editor: widget.Editor{
					SingleLine: true,
					MaxLen:     255,
				},
			},
			V: component.TextField{
				Editor: widget.Editor{
					SingleLine: true,
					MaxLen:     255,
				},
			},
		},
		delMetadataDialog: ui_widget.Dialog{Title: i18n.DeleteMetadata},
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
	md := api.NewMetadata(cfg.Metadata)
	for k := range md {
		if k == "" {
			continue
		}
		p.metadata = append(p.metadata, metadata{
			k: k,
			v: md.GetString(k),
		})
	}
	p.mdSelector.Clear()
	p.mdSelector.Select(ui_widget.SelectorItem{Value: strconv.Itoa(len(p.metadata))})
	p.mdFolded = true
}

func (p *dialer) Layout(gtx page.C, th *page.T) page.D {
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
			if p.mdAdd.Clicked(gtx) {
				p.showMetadataDialog(gtx, -1)
			}

			return layout.Flex{
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Flexed(1, func(gtx page.C) page.D {
					gtx.Source = src

					if p.mdSelector.Clicked(gtx) {
						p.mdFolded = !p.mdFolded
					}
					return p.mdSelector.Layout(gtx, th)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if !p.node.edit {
						return page.D{}
					}
					return layout.Spacer{Width: 8}.Layout(gtx)
				}),
				layout.Rigid(func(gtx page.C) page.D {
					if !p.node.edit {
						return page.D{}
					}
					btn := material.IconButton(th, &p.mdAdd, icons.IconAdd, "Add")
					btn.Background = theme.Current().ContentSurfaceBg
					btn.Color = th.Fg
					// btn.Inset = layout.UniformInset(8)
					return btn.Layout(gtx)
				}),
			)
		}),
		layout.Rigid(func(gtx page.C) page.D {
			if p.mdFolded {
				return page.D{}
			}

			gtx.Source = src
			return p.layoutMetadata(gtx, th)
		}),
	)
}

func (p *dialer) layoutMetadata(gtx page.C, th *page.T) page.D {
	for i := range p.metadata {
		if p.metadata[i].clk.Clicked(gtx) {
			if p.node.edit {
				p.showMetadataDialog(gtx, i)
			}
			break
		}

		if p.metadata[i].delete.Clicked(gtx) {
			p.delMetadataDialog.OnClick = func(ok bool) {
				p.node.router.HideModal(gtx)
				if !ok {
					return
				}
				p.metadata = append(p.metadata[:i], p.metadata[i+1:]...)

				p.mdSelector.Clear()
				p.mdSelector.Select(ui_widget.SelectorItem{Value: strconv.Itoa(len(p.metadata))})
			}
			p.node.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
				return p.delMetadataDialog.Layout(gtx, th)
			})
			break
		}
	}

	var children []layout.FlexChild
	for i := range p.metadata {
		md := &p.metadata[i]

		children = append(children,
			layout.Rigid(func(gtx page.C) page.D {
				return layout.Flex{
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Flexed(1, func(gtx page.C) page.D {
						return material.Clickable(gtx, &md.clk, func(gtx page.C) page.D {
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
										label := material.Body2(th, md.k)
										label.Font.Weight = font.SemiBold
										return label.Layout(gtx)
									}),
									layout.Rigid(layout.Spacer{Height: 4}.Layout),
									layout.Rigid(func(gtx page.C) page.D {
										return material.Body2(th, md.v).Layout(gtx)
									}),
								)
							})
						})
					}),
					layout.Rigid(func(gtx page.C) page.D {
						if !p.node.edit {
							return page.D{}
						}
						return layout.Spacer{Width: 8}.Layout(gtx)
					}),
					layout.Rigid(func(gtx page.C) page.D {
						if !p.node.edit {
							return page.D{}
						}
						btn := material.IconButton(th, &md.delete, icons.IconDelete, "delete")
						btn.Background = theme.Current().ContentSurfaceBg
						btn.Color = th.Fg
						// btn.Inset = layout.UniformInset(8)
						return btn.Layout(gtx)
					}),
				)
			}),
		)
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx, children...)
}

func (p *dialer) showMetadataDialog(gtx page.C, i int) {
	p.mdDialog.K.Clear()
	p.mdDialog.V.Clear()

	if i >= 0 && i < len(p.metadata) {
		p.mdDialog.K.SetText(p.metadata[i].k)
		p.mdDialog.V.SetText(p.metadata[i].v)
	}

	p.mdDialog.OnClick = func(ok bool) {
		p.node.router.HideModal(gtx)
		if !ok {
			return
		}

		k, v := strings.TrimSpace(p.mdDialog.K.Text()), strings.TrimSpace(p.mdDialog.V.Text())
		if k == "" {
			return
		}

		if i >= 0 && i < len(p.metadata) {
			p.metadata[i].k = k
			p.metadata[i].v = v
		} else {
			p.metadata = append(p.metadata, metadata{
				k: k,
				v: v,
			})
		}

		p.mdSelector.Clear()
		p.mdSelector.Select(ui_widget.SelectorItem{Value: strconv.Itoa(len(p.metadata))})
	}

	p.node.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.mdDialog.Layout(gtx, th)
	})
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
	p.menu.OnAdd = nil
	p.menu.Multiple = false

	p.node.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}
