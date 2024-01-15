package service

import (
	"strings"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/go-gost/gui/api"
	"github.com/go-gost/gui/ui/icons"
	"github.com/go-gost/gui/ui/page"
)

type listener struct {
	modal *component.ModalLayer
	menu  *page.Menu

	typ      string
	btnTyp   widget.Clickable
	chain    string
	btnChain widget.Clickable

	authType  widget.Enum
	username  component.TextField
	password  component.TextField
	authers   []string
	btnAuther widget.Clickable

	enableTLS   widget.Bool
	tlsCertFile component.TextField
	tlsKeyFile  component.TextField
	tlsCAFile   component.TextField

	metadata    []metadata
	addMetadata widget.Clickable
}

func (l *listener) Layout(gtx C, th *material.Theme) D {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if l.btnTyp.Clicked(gtx) {
				l.showTypeMenu(gtx)
			}

			return l.btnTyp.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top:    8,
					Bottom: 8,
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Flexed(1, material.Body1(th, "Type").Layout),
						layout.Rigid(material.Body2(th, l.typ).Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return icons.IconNavRight.Layout(gtx, th.Fg)
						}),
					)
				})
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if l.btnChain.Clicked(gtx) {
				l.showChainMenu(gtx)
			}

			return l.btnChain.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top:    8,
					Bottom: 8,
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Flexed(1, material.Body1(th, "Chain").Layout),
						layout.Rigid(material.Body2(th, l.chain).Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return icons.IconNavRight.Layout(gtx, th.Fg)
						}),
					)
				})
			})
		}),
		layout.Rigid(layout.Spacer{Height: 5}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Alignment: layout.Middle,
				Spacing:   layout.SpaceBetween,
			}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return material.Body1(th, "Auth").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.RadioButton(th, &l.authType, "simple", "Simple").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.RadioButton(th, &l.authType, "auther", "Auther").Layout(gtx)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if l.authType.Value != "simple" {
				return D{}
			}
			return l.username.Layout(gtx, th, "Username")
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if l.authType.Value != "simple" {
				return D{}
			}
			return l.password.Layout(gtx, th, "Password")
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if l.authType.Value != "auther" {
				return D{}
			}

			if l.btnAuther.Clicked(gtx) {
				l.showAutherMenu(gtx)
			}

			return l.btnAuther.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top:    8,
					Bottom: 8,
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{
								Right: 5,
							}.Layout(gtx, material.Body1(th, "Auther").Layout)
						}),
						layout.Flexed(1, layout.Spacer{Width: 5}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return material.Body2(th, strings.Join(l.authers, ", ")).Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return icons.IconNavRight.Layout(gtx, th.Fg)
						}),
					)
				})
			})
		}),
		layout.Rigid(layout.Spacer{Height: 10}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Flexed(1, material.Body1(th, "TLS").Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Switch(th, &l.enableTLS, "TLS").Layout(gtx)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if !l.enableTLS.Value {
				return D{}
			}
			return l.tlsCertFile.Layout(gtx, th, "Cert File")
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if !l.enableTLS.Value {
				return D{}
			}
			return l.tlsKeyFile.Layout(gtx, th, "Key File")
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if !l.enableTLS.Value {
				return D{}
			}
			return l.tlsCAFile.Layout(gtx, th, "CA File")
		}),

		layout.Rigid(layout.Spacer{Height: 15}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if l.addMetadata.Clicked(gtx) {
				l.metadata = append(l.metadata, metadata{})
			}

			return layout.Flex{
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Flexed(1, material.Body1(th, "Metadata").Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btn := material.IconButton(th, &l.addMetadata, icons.IconAdd, "Add")
					btn.Background = th.Bg
					btn.Color = th.Fg
					btn.Inset = layout.UniformInset(0)
					return btn.Layout(gtx)
				}),
			)
		}),
		layout.Rigid(layout.Spacer{Height: 5}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return l.layoutMetadata(gtx, th)
		}),
	)
}

func (h *listener) layoutMetadata(gtx C, th *material.Theme) D {
	for i := range h.metadata {
		if h.metadata[i].delete.Clicked(gtx) {
			h.metadata = append(h.metadata[:i], h.metadata[i+1:]...)
			break
		}
	}

	var children []layout.FlexChild
	for i := range h.metadata {
		i := i
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return h.metadata[i].k.Layout(gtx, th, "Key")
					}),
					layout.Rigid(layout.Spacer{Width: 5}.Layout),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return h.metadata[i].v.Layout(gtx, th, "Value")
					}),
					layout.Rigid(layout.Spacer{Width: 5}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.IconButton(th, &h.metadata[i].delete, icons.IconDelete, "delete")
						btn.Background = th.Bg
						btn.Color = th.Fg
						btn.Inset = layout.UniformInset(5)
						return btn.Layout(gtx)
					}),
				)
			})
		}))
	}
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx, children...)
}

func (l *listener) showTypeMenu(gtx C) {
	items := []page.MenuItem{
		{Key: "TCP", Value: "tcp"},
		{Key: "UDP", Value: "udp"},
		{Key: "RTCP", Value: "rtcp"},
		{Key: "RUDP", Value: "rudp"},
		{Key: "TLS", Value: "tls"},
		{Key: "Multiplexed TLS", Value: "mtls"},
		{Key: "Websocket", Value: "ws"},
		{Key: "Multiplexed Websocket", Value: "mws"},
		{Key: "Websocket Secure", Value: "wss"},
		{Key: "Multiplexed Websocket Secure", Value: "mwss"},
		{Key: "HTTP2", Value: "http2"},
		{Key: "HTTP3", Value: "http3"},
		{Key: "gRPC", Value: "grpc"},
		{Key: "QUIC", Value: "quic"},
		{Key: "WebTransport", Value: "wt"},
		{Key: "KCP", Value: "kcp"},
		{Key: "DTLS", Value: "dtls"},
		{Key: "DTLS", Value: "dtls"},

		{Key: "Plain HTTP Tunnel", Value: "pht"},
		{Key: "Obfs-HTTP", Value: "ohttp"},
		{Key: "Obfs-TLS", Value: "otls"},

		{Key: "SSH", Value: "ssh"},
		{Key: "SSHD", Value: "sshd"},
		{Key: "DNS", Value: "dns"},

		{Key: "Multiplexed TCP", Value: "mtcp"},
		{Key: "Fake TCP", Value: "ftcp"},
		{Key: "ICMP", Value: "icmp"},

		{Key: "TCP Redirect", Value: "red"},
		{Key: "UDP Redirect", Value: "redu"},
		{Key: "TUN", Value: "tun"},
		{Key: "TAP", Value: "tap"},

		{Key: "Serial Port Redirector", Value: "serial"},
		{Key: "Unix Domain Socket", Value: "unix"},
	}
	for i := range items {
		if items[i].Value == l.typ {
			items[i].Selected = true
		}
	}

	l.menu.Title = "Listener Type"
	l.menu.Items = items
	l.menu.Selected = func(index int) {
		l.typ = l.menu.Items[index].Value
		l.modal.Disappear(gtx.Now)
	}
	l.menu.ShowAdd = false

	l.modal.Widget = func(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
		return l.menu.Layout(gtx, th)
	}
	l.modal.Appear(gtx.Now)
}

func (l *listener) showChainMenu(gtx C) {
	items := []page.MenuItem{
		{Key: "N/A"},
	}
	for _, v := range api.GetConfig().Chains {
		items = append(items, page.MenuItem{
			Key:   v.Name,
			Value: v.Name,
		})
	}
	for i := range items {
		if items[i].Value == l.chain {
			items[i].Selected = true
		}
	}

	l.menu.Title = "Chain"
	l.menu.Items = items
	l.menu.Selected = func(index int) {
		l.chain = l.menu.Items[index].Value
		l.modal.Disappear(gtx.Now)
	}
	l.menu.ShowAdd = true

	l.modal.Widget = func(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
		return l.menu.Layout(gtx, th)
	}
	l.modal.Appear(gtx.Now)
}

func (l *listener) showAutherMenu(gtx C) {
	items := []page.MenuItem{}
	for _, v := range api.GetConfig().Authers {
		items = append(items, page.MenuItem{
			Key:   v.Name,
			Value: v.Name,
		})
	}
	for i := range items {
		for _, v := range l.authers {
			if items[i].Value == v {
				items[i].Selected = true
			}
		}
	}

	l.menu.Title = "Auther"
	l.menu.Items = items
	l.menu.Selected = func(index int) {
		l.authers = nil
		for i := range l.menu.Items {
			if l.menu.Items[i].Selected {
				l.authers = append(l.authers, l.menu.Items[i].Value)
			}
		}
	}
	l.menu.ShowAdd = true
	l.menu.Multiple = true

	l.modal.Widget = func(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
		return l.menu.Layout(gtx, th)
	}
	l.modal.Appear(gtx.Now)
}
