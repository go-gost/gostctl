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

type handler struct {
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

	limiter     string
	btnLimiter  widget.Clickable
	observer    string
	btnObserver widget.Clickable

	metadata    []metadata
	addMetadata widget.Clickable
}

func (h *handler) Layout(gtx C, th *material.Theme) D {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if h.btnTyp.Clicked(gtx) {
				h.showTypeMenu(gtx)
			}

			return h.btnTyp.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top:    8,
					Bottom: 8,
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Flexed(1, material.Body1(th, "Type").Layout),
						layout.Rigid(material.Body2(th, h.typ).Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return icons.IconNavRight.Layout(gtx, th.Fg)
						}),
					)
				})
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if h.btnChain.Clicked(gtx) {
				h.showChainMenu(gtx)
			}

			return h.btnChain.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top:    8,
					Bottom: 8,
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Flexed(1, material.Body1(th, "Chain").Layout),
						layout.Rigid(material.Body2(th, h.chain).Layout),
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
					return material.RadioButton(th, &h.authType, "simple", "Simple").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.RadioButton(th, &h.authType, "auther", "Auther").Layout(gtx)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if h.authType.Value != "simple" {
				return D{}
			}
			return h.username.Layout(gtx, th, "Username")
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if h.authType.Value != "simple" {
				return D{}
			}
			return h.password.Layout(gtx, th, "Password")
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if h.authType.Value != "auther" {
				return D{}
			}

			if h.btnAuther.Clicked(gtx) {
				h.showAutherMenu(gtx)
			}

			return h.btnAuther.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
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
							return material.Body2(th, strings.Join(h.authers, ", ")).Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return icons.IconNavRight.Layout(gtx, th.Fg)
						}),
					)
				})
			})
		}),

		layout.Rigid(layout.Spacer{Height: 5}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if h.btnLimiter.Clicked(gtx) {
				h.showLimiterMenu(gtx)
			}

			return h.btnLimiter.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top:    8,
					Bottom: 8,
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Flexed(1, material.Body1(th, "Limiter").Layout),
						layout.Rigid(material.Body2(th, h.limiter).Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return icons.IconNavRight.Layout(gtx, th.Fg)
						}),
					)
				})
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if h.btnObserver.Clicked(gtx) {
				h.showObserverMenu(gtx)
			}

			return h.btnObserver.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top:    8,
					Bottom: 8,
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Flexed(1, material.Body1(th, "Observer").Layout),
						layout.Rigid(material.Body2(th, h.observer).Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return icons.IconNavRight.Layout(gtx, th.Fg)
						}),
					)
				})
			})
		}),

		layout.Rigid(layout.Spacer{Height: 5}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if h.addMetadata.Clicked(gtx) {
				h.metadata = append(h.metadata, metadata{})
			}

			return layout.Flex{
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Flexed(1, material.Body1(th, "Metadata").Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btn := material.IconButton(th, &h.addMetadata, icons.IconAdd, "Add")
					btn.Background = th.Bg
					btn.Color = th.Fg
					btn.Inset = layout.UniformInset(0)
					return btn.Layout(gtx)
				}),
			)
		}),
		layout.Rigid(layout.Spacer{Height: 5}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return h.layoutMetadata(gtx, th)
		}),
	)
}

func (h *handler) layoutMetadata(gtx C, th *material.Theme) D {
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

func (h *handler) showTypeMenu(gtx C) {
	items := []page.MenuItem{
		{Key: "Auto", Value: "auto"},
		{Key: "HTTP", Value: "http"},
		{Key: "SOCKS4", Value: "socks4"},
		{Key: "SOCKS5", Value: "socks5"},
		{Key: "Relay", Value: "relay"},
		{Key: "Shadowsocks", Value: "ss"},
		{Key: "HTTP2", Value: "http2"},
		{Key: "HTTP3", Value: "http3"},

		{Key: "TCP", Value: "tcp"},
		{Key: "UDP", Value: "udp"},
		{Key: "RTCP", Value: "rtcp"},
		{Key: "RUDP", Value: "rudp"},

		{Key: "SNI", Value: "sni"},
		{Key: "DNS", Value: "dns"},
		{Key: "SSHD", Value: "sshd"},

		{Key: "TCP Redirect", Value: "red"},
		{Key: "UDP Redirect", Value: "redu"},
		{Key: "TUN", Value: "tun"},
		{Key: "TAP", Value: "tap"},
		{Key: "Tunnel", Value: "tunnel"},

		{Key: "File Server", Value: "file"},
		{Key: "Serial Port Redirector", Value: "serial"},
		{Key: "Unix Domain Socket", Value: "unix"},
	}
	for i := range items {
		if items[i].Value == h.typ {
			items[i].Selected = true
		}
	}

	h.menu.Title = "Handler Type"
	h.menu.Items = items
	h.menu.Selected = func(index int) {
		h.typ = h.menu.Items[index].Value
		h.modal.Disappear(gtx.Now)
	}
	h.menu.ShowAdd = false

	h.modal.Widget = func(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
		return h.menu.Layout(gtx, th)
	}
	h.modal.Appear(gtx.Now)
}

func (h *handler) showChainMenu(gtx C) {
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
		if items[i].Value == h.chain {
			items[i].Selected = true
		}
	}

	h.menu.Title = "Chain"
	h.menu.Items = items
	h.menu.Selected = func(index int) {
		h.chain = h.menu.Items[index].Value
		h.modal.Disappear(gtx.Now)
	}
	h.menu.ShowAdd = true

	h.modal.Widget = func(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
		return h.menu.Layout(gtx, th)
	}
	h.modal.Appear(gtx.Now)
}

func (h *handler) showAutherMenu(gtx C) {
	items := []page.MenuItem{}
	for _, v := range api.GetConfig().Authers {
		items = append(items, page.MenuItem{
			Key:   v.Name,
			Value: v.Name,
		})
	}
	for i := range items {
		for _, v := range h.authers {
			if items[i].Value == v {
				items[i].Selected = true
			}
		}
	}

	h.menu.Title = "Auther"
	h.menu.Items = items
	h.menu.Selected = func(index int) {
		h.authers = nil
		for i := range h.menu.Items {
			if h.menu.Items[i].Selected {
				h.authers = append(h.authers, h.menu.Items[i].Value)
			}
		}
	}
	h.menu.ShowAdd = true
	h.menu.Multiple = true

	h.modal.Widget = func(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
		return h.menu.Layout(gtx, th)
	}
	h.modal.Appear(gtx.Now)
}

func (h *handler) showLimiterMenu(gtx C) {
	items := []page.MenuItem{
		{Key: "N/A"},
	}
	for _, v := range api.GetConfig().Limiters {
		items = append(items, page.MenuItem{
			Key:   v.Name,
			Value: v.Name,
		})
	}
	for i := range items {
		if items[i].Value == h.limiter {
			items[i].Selected = true
		}
	}

	h.menu.Title = "Limiter"
	h.menu.Items = items
	h.menu.Selected = func(index int) {
		h.limiter = h.menu.Items[index].Value
		h.modal.Disappear(gtx.Now)
	}
	h.menu.ShowAdd = true

	h.modal.Widget = func(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
		return h.menu.Layout(gtx, th)
	}
	h.modal.Appear(gtx.Now)
}

func (h *handler) showObserverMenu(gtx C) {
	items := []page.MenuItem{
		{Key: "N/A"},
	}
	for _, v := range api.GetConfig().Observers {
		items = append(items, page.MenuItem{
			Key:   v.Name,
			Value: v.Name,
		})
	}
	for i := range items {
		if items[i].Value == h.observer {
			items[i].Selected = true
		}
	}

	h.menu.Title = "Observer"
	h.menu.Items = items
	h.menu.Selected = func(index int) {
		h.observer = h.menu.Items[index].Value
		h.modal.Disappear(gtx.Now)
	}
	h.menu.ShowAdd = true

	h.modal.Widget = func(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
		return h.menu.Layout(gtx, th)
	}
	h.modal.Appear(gtx.Now)
}
