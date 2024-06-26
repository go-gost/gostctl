package service

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

	metadata          []metadata
	mdSelector        ui_widget.Selector
	mdFolded          bool
	mdAdd             widget.Clickable
	mdDialog          ui_widget.MetadataDialog
	delMetadataDialog ui_widget.Dialog
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
	md := api.NewMetadata(cfg.Metadata)
	for k := range md {
		if k == "" {
			continue
		}
		l.metadata = append(l.metadata, metadata{
			k: k,
			v: md.GetString(k),
		})
	}
	l.mdSelector.Clear()
	l.mdSelector.Select(ui_widget.SelectorItem{Value: strconv.Itoa(len(l.metadata))})
	l.mdFolded = true
}

func (l *listener) Layout(gtx page.C, th *page.T) page.D {
	src := gtx.Source

	if !l.service.edit {
		gtx = gtx.Disabled()
	}

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
							layout.Rigid(layout.Spacer{Width: 8}.Layout),
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
			if l.mdAdd.Clicked(gtx) {
				l.showMetadataDialog(gtx, -1)
			}

			return layout.Flex{
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Flexed(1, func(gtx page.C) page.D {
					gtx.Source = src

					if l.mdSelector.Clicked(gtx) {
						l.mdFolded = !l.mdFolded
					}
					return l.mdSelector.Layout(gtx, th)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if !l.service.edit {
						return page.D{}
					}
					return layout.Spacer{Width: 8}.Layout(gtx)
				}),
				layout.Rigid(func(gtx page.C) page.D {
					if !l.service.edit {
						return page.D{}
					}
					btn := material.IconButton(th, &l.mdAdd, icons.IconAdd, "Add")
					btn.Background = theme.Current().ContentSurfaceBg
					btn.Color = th.Fg
					// btn.Inset = layout.UniformInset(8)
					return btn.Layout(gtx)
				}),
			)
		}),
		layout.Rigid(func(gtx page.C) page.D {
			if l.mdFolded {
				return page.D{}
			}

			gtx.Source = src
			return l.layoutMetadata(gtx, th)
		}),
	)
}

func (l *listener) layoutMetadata(gtx page.C, th *page.T) page.D {
	for i := range l.metadata {
		if l.metadata[i].clk.Clicked(gtx) {
			if l.service.edit {
				l.showMetadataDialog(gtx, i)
			}
			break
		}

		if l.metadata[i].delete.Clicked(gtx) {
			l.delMetadataDialog.OnClick = func(ok bool) {
				l.service.router.HideModal(gtx)
				if !ok {
					return
				}
				l.metadata = append(l.metadata[:i], l.metadata[i+1:]...)

				l.mdSelector.Clear()
				l.mdSelector.Select(ui_widget.SelectorItem{Value: strconv.Itoa(len(l.metadata))})
			}
			l.service.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
				return l.delMetadataDialog.Layout(gtx, th)
			})
			break
		}
	}

	var children []layout.FlexChild
	for i := range l.metadata {
		md := &l.metadata[i]

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
						if !l.service.edit {
							return page.D{}
						}
						return layout.Spacer{Width: 8}.Layout(gtx)
					}),
					layout.Rigid(func(gtx page.C) page.D {
						if !l.service.edit {
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

func (l *listener) showMetadataDialog(gtx page.C, i int) {
	l.mdDialog.K.Clear()
	l.mdDialog.V.Clear()

	if i >= 0 && i < len(l.metadata) {
		l.mdDialog.K.SetText(l.metadata[i].k)
		l.mdDialog.V.SetText(l.metadata[i].v)
	}

	l.mdDialog.OnClick = func(ok bool) {
		l.service.router.HideModal(gtx)
		if !ok {
			return
		}

		k, v := strings.TrimSpace(l.mdDialog.K.Text()), strings.TrimSpace(l.mdDialog.V.Text())
		if k == "" {
			return
		}

		if i >= 0 && i < len(l.metadata) {
			l.metadata[i].k = k
			l.metadata[i].v = v
		} else {
			l.metadata = append(l.metadata, metadata{
				k: k,
				v: v,
			})
		}

		l.mdSelector.Clear()
		l.mdSelector.Select(ui_widget.SelectorItem{Value: strconv.Itoa(len(l.metadata))})
	}

	l.service.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return l.mdDialog.Layout(gtx, th)
	})
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

		{Name: "RED", Value: "red"},
		{Name: "REDU", Value: "redu"},
		{Name: "TUN", Value: "tun"},
		{Name: "TAP", Value: "tap"},

		{Name: "PHT", Value: "pht"},
		{Name: "OHTTP", Value: "ohttp"},
		{Name: "OTLS", Value: "otls"},

		{Name: "HTTP3", Value: "http3"},
		{Name: "FTCP", Value: "ftcp"},
		{Name: "ICMP", Value: "icmp"},

		{Key: "Serial", Value: "serial"},
		{Key: "Unix", Value: "unix"},
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
	l.menu.OnAdd = nil
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
	l.menu.OnAdd = func() {
		l.service.router.Goto(page.Route{
			Path: page.PageChain,
			Perm: page.PermReadWrite,
		})
		l.service.router.HideModal(gtx)
	}
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
	l.menu.OnAdd = func() {
		l.service.router.Goto(page.Route{
			Path: page.PageAuther,
			Perm: page.PermReadWrite,
		})
		l.service.router.HideModal(gtx)
	}
	l.menu.Multiple = true

	l.service.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return l.menu.Layout(gtx, th)
	})
}
