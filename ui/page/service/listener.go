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
	router *page.Router
	menu   ui_widget.Menu
	mode   *widget.Enum

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

					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(gtx,
						layout.Rigid(func(gtx page.C) page.D {
							return l.username.Layout(gtx, th, i18n.Username.Value())
						}),
						layout.Rigid(func(gtx page.C) page.D {
							return l.password.Layout(gtx, th, i18n.Password.Value())
						}),
						layout.Rigid(layout.Spacer{Height: 4}.Layout),
					)
				}),

				layout.Rigid(func(gtx page.C) page.D {
					if l.authType.Value != string(page.AuthAuther) {
						return page.D{}
					}

					if l.auther.Clicked(gtx) {
						l.showAutherMenu(gtx)
					}
					return l.auther.Layout(gtx, th)
				}),
			)
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
					layout.Rigid(layout.Spacer{Height: 4}.Layout),
					layout.Rigid(func(gtx page.C) page.D {
						if !l.enableTLS.Value() {
							return page.D{}
						}
						return l.tlsCertFile.Layout(gtx, th, i18n.CertFile.Value())
					}),
					layout.Rigid(func(gtx page.C) page.D {
						if !l.enableTLS.Value() {
							return page.D{}
						}
						return l.tlsKeyFile.Layout(gtx, th, i18n.KeyFile.Value())
					}),
					layout.Rigid(func(gtx page.C) page.D {
						if !l.enableTLS.Value() {
							return page.D{}
						}
						return l.tlsCAFile.Layout(gtx, th, i18n.CAFile.Value())
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
	options := []ui_widget.MenuOption{
		{Key: "TCP", Value: "tcp"},
		{Key: "UDP", Value: "udp"},
		{Key: "RTCP", Value: "rtcp"},
		{Key: "RUDP", Value: "rudp"},
		{Key: "TLS", Value: "tls"},
		{Key: "WS", Value: "ws"},
		{Key: "WSS", Value: "wss"},
		{Key: "HTTP2", Value: "http2"},
		{Key: "H2", Value: "h2"},
		{Key: "H2C", Value: "h2c"},
		{Key: "KCP", Value: "kcp"},
	}
	if l.mode.Value == string(page.AdvancedMode) {
		options = append(options, []ui_widget.MenuOption{
			{Key: "MTLS", Value: "mtls"},
			{Key: "MWS", Value: "mws"},
			{Key: "MWSS", Value: "mwss"},
			{Key: "gRPC", Value: "grpc"},
			{Key: "QUIC", Value: "quic"},
			{Key: "HTTP3", Value: "http3"},

			{Key: "SSH", Value: "ssh"},
			{Key: "SSHD", Value: "sshd"},
			{Key: "DNS", Value: "dns"},

			{Key: "TCP Redirect", Value: "red"},
			{Key: "UDP Redirect", Value: "redu"},
			{Key: "TUN", Value: "tun"},
			{Key: "TAP", Value: "tap"},

			{Key: "PHT", Value: "pht"},
			{Key: "Obfs-HTTP", Value: "ohttp"},
			{Key: "Obfs-TLS", Value: "otls"},
			{Key: "WebTransport", Value: "wt"},
			{Key: "DTLS", Value: "dtls"},

			{Key: "MTCP", Value: "mtcp"},
			{Key: "Fake TCP", Value: "ftcp"},
			{Key: "ICMP", Value: "icmp"},

			{Key: "Serial Port Redirector", Value: "serial"},
			{Key: "Unix Domain Socket", Value: "unix"},
		}...)
	}

	for i := range options {
		options[i].Selected = l.typ.AnyValue(options[i].Value)
	}

	l.menu.Title = i18n.Listener
	l.menu.Options = options
	l.menu.OnClick = func(ok bool) {
		l.router.HideModal(gtx)
		if !ok {
			return
		}

		l.typ.Clear()
		for i := range l.menu.Options {
			if l.menu.Options[i].Selected {
				l.typ.Select(ui_widget.SelectorItem{Value: l.menu.Options[i].Value})
			}
		}
	}
	l.menu.ShowAdd = false
	l.menu.Multiple = false

	l.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
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
		l.router.HideModal(gtx)
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

	l.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
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
		l.router.HideModal(gtx)
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

	l.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return l.menu.Layout(gtx, th)
	})
}

func (l *listener) showMetadataDialog(gtx page.C) {
	l.metadataDialog.Clear()
	for _, md := range l.metadata {
		l.metadataDialog.Add(md.K, md.V)
	}
	l.metadataDialog.OnClick = func(ok bool) {
		l.router.HideModal(gtx)
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

	l.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return l.metadataDialog.Layout(gtx, th)
	})
}
