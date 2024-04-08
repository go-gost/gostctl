package service

import (
	"strconv"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/go-gost/gostctl/api"
	"github.com/go-gost/gostctl/ui/i18n"
	ui_widget "github.com/go-gost/gostctl/ui/widget"
)

type handler struct {
	modal *component.ModalLayer
	menu  ui_widget.Menu
	mode  *widget.Enum

	typ   ui_widget.Selector
	chain ui_widget.Selector

	authType widget.Enum
	username component.TextField
	password component.TextField
	auther   ui_widget.Selector

	limiter  ui_widget.Selector
	observer ui_widget.Selector

	metadata         []metadata
	metadataSelector ui_widget.Selector
	metadataDialog   ui_widget.MetadataDialog
}

func (h *handler) Layout(gtx C, th *material.Theme) D {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if h.typ.Clicked(gtx) {
				h.showTypeMenu(gtx)
			}

			return h.typ.Layout(gtx, th)
		}),

		// auth for handler
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if !h.canAuth() {
				return D{}
			}

			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Top:    4,
						Bottom: 4,
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Alignment: layout.Middle,
							Spacing:   layout.SpaceBetween,
						}.Layout(gtx,
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return material.Body1(th, i18n.Auth.Value()).Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return material.RadioButton(th, &h.authType, AuthTypeSimple, i18n.AuthSimple.Value()).Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Width: 8}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return material.RadioButton(th, &h.authType, AuthTypeAuther, i18n.AuthAuther.Value()).Layout(gtx)
							}),
						)
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if h.authType.Value != AuthTypeSimple {
						return D{}
					}

					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return h.username.Layout(gtx, th, i18n.Username.Value())
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return h.password.Layout(gtx, th, i18n.Password.Value())
						}),
						layout.Rigid(layout.Spacer{Height: 4}.Layout),
					)
				}),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if h.authType.Value != AuthTypeAuther {
						return D{}
					}

					if h.auther.Clicked(gtx) {
						h.showAutherMenu(gtx)
					}

					return h.auther.Layout(gtx, th)
				}),
			)
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if h.chain.Clicked(gtx) {
				h.showChainMenu(gtx)
			}

			return h.chain.Layout(gtx, th)
		}),

		// advanced mode
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if h.mode.Value == BasicMode {
				return D{}
			}

			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if h.limiter.Clicked(gtx) {
						h.showLimiterMenu(gtx)
					}
					return h.limiter.Layout(gtx, th)
				}),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if h.observer.Clicked(gtx) {
						h.showObserverMenu(gtx)
					}
					return h.observer.Layout(gtx, th)
				}),
			)
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if h.metadataSelector.Clicked(gtx) {
				h.showMetadataDialog(gtx)
			}
			return h.metadataSelector.Layout(gtx, th)
		}),
	)
}

func (h *handler) canAuth() bool {
	return h.typ.AnyValue("auto", "http", "http2", "socks4", "socks5", "socks", "relay", "ss", "file")
}

func (h *handler) canForward() bool {
	return h.typ.AnyValue("tcp", "udp", "rtcp", "rudp", "forward", "rforward", "dns", "serial", "unix")
}

func (h *handler) showTypeMenu(gtx C) {
	options := []ui_widget.MenuOption{
		{Key: "Auto", Value: "auto"},
		{Key: "HTTP", Value: "http"},
		{Key: "SOCKS4", Value: "socks4"},
		{Key: "SOCKS5", Value: "socks5"},
		{Key: "Relay", Value: "relay"},
		{Key: "Shadowsocks", Value: "ss"},
		{Key: "HTTP2", Value: "http2"},

		{Key: "TCP", Value: "tcp"},
		{Key: "UDP", Value: "udp"},
		{Key: "RTCP", Value: "rtcp"},
		{Key: "RUDP", Value: "rudp"},
	}

	if h.mode.Value == AdvancedMode {
		options = append(options, []ui_widget.MenuOption{
			{Key: "SNI", Value: "sni"},
			{Key: "DNS", Value: "dns"},
			{Key: "SSHD", Value: "sshd"},
			{Key: "HTTP3", Value: "http3"},

			{Key: "TCP Redirect", Value: "red"},
			{Key: "UDP Redirect", Value: "redu"},
			{Key: "TUN", Value: "tun"},
			{Key: "TAP", Value: "tap"},
			{Key: "Tunnel", Value: "tunnel"},

			{Key: "File Server", Value: "file"},
			{Key: "Serial Port Redirector", Value: "serial"},
			{Key: "Unix Domain Socket", Value: "unix"},
		}...)
	}

	for i := range options {
		options[i].Selected = h.typ.AnyValue(options[i].Value)
	}

	h.menu.Title = i18n.Handler
	h.menu.Options = options
	h.menu.OnClick = func(ok bool) {
		h.modal.Disappear(gtx.Now)
		if !ok {
			return
		}

		h.typ.Clear()
		for i := range h.menu.Options {
			if h.menu.Options[i].Selected {
				h.typ.Select(ui_widget.SelectorItem{Value: h.menu.Options[i].Value})
			}
		}
		h.modal.Disappear(gtx.Now)
	}
	h.menu.ShowAdd = false

	h.modal.Widget = func(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
		return h.menu.Layout(gtx, th)
	}
	h.modal.Appear(gtx.Now)
}

func (h *handler) showChainMenu(gtx C) {
	options := []ui_widget.MenuOption{}
	for _, v := range api.GetConfig().Chains {
		options = append(options, ui_widget.MenuOption{
			Value: v.Name,
		})
	}
	for i := range options {
		options[i].Selected = h.chain.AnyValue(options[i].Value)
	}

	h.menu.Title = i18n.Chain
	h.menu.Options = options
	h.menu.OnClick = func(ok bool) {
		h.modal.Disappear(gtx.Now)
		if !ok {
			return
		}

		h.chain.Clear()
		for i := range h.menu.Options {
			if h.menu.Options[i].Selected {
				h.chain.Select(ui_widget.SelectorItem{Value: h.menu.Options[i].Value})
			}
		}
		h.modal.Disappear(gtx.Now)
	}
	h.menu.ShowAdd = true

	h.modal.Widget = func(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
		return h.menu.Layout(gtx, th)
	}
	h.modal.Appear(gtx.Now)
}

func (h *handler) showAutherMenu(gtx C) {
	options := []ui_widget.MenuOption{}
	for _, v := range api.GetConfig().Authers {
		options = append(options, ui_widget.MenuOption{
			Value: v.Name,
		})
	}
	for i := range options {
		options[i].Selected = h.auther.AnyValue(options[i].Value)
	}

	h.menu.Title = i18n.Auther
	h.menu.Options = options
	h.menu.OnClick = func(ok bool) {
		h.modal.Disappear(gtx.Now)
		if !ok {
			return
		}

		h.auther.Clear()
		for i := range h.menu.Options {
			if h.menu.Options[i].Selected {
				h.auther.Select(ui_widget.SelectorItem{Value: h.menu.Options[i].Value})
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
	options := []ui_widget.MenuOption{}
	for _, v := range api.GetConfig().Limiters {
		options = append(options, ui_widget.MenuOption{
			Value: v.Name,
		})
	}
	for i := range options {
		options[i].Selected = h.limiter.AnyValue(options[i].Value)
	}

	h.menu.Title = i18n.Limiter
	h.menu.Options = options
	h.menu.OnClick = func(ok bool) {
		h.modal.Disappear(gtx.Now)
		if !ok {
			return
		}

		h.limiter.Clear()
		for i := range h.menu.Options {
			if h.menu.Options[i].Selected {
				h.limiter.Select(ui_widget.SelectorItem{Value: h.menu.Options[i].Value})
			}
		}
		h.modal.Disappear(gtx.Now)
	}
	h.menu.ShowAdd = true

	h.modal.Widget = func(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
		return h.menu.Layout(gtx, th)
	}
	h.modal.Appear(gtx.Now)
}

func (h *handler) showObserverMenu(gtx C) {
	options := []ui_widget.MenuOption{}
	for _, v := range api.GetConfig().Observers {
		options = append(options, ui_widget.MenuOption{
			Value: v.Name,
		})
	}
	for i := range options {
		options[i].Selected = h.observer.AnyValue(options[i].Value)
	}

	h.menu.Title = i18n.Observer
	h.menu.Options = options
	h.menu.OnClick = func(ok bool) {
		h.modal.Disappear(gtx.Now)
		if !ok {
			return
		}

		h.observer.Clear()
		for i := range h.menu.Options {
			if h.menu.Options[i].Selected {
				h.observer.Select(ui_widget.SelectorItem{Value: h.menu.Options[i].Value})
			}
		}
		h.modal.Disappear(gtx.Now)
	}
	h.menu.ShowAdd = true

	h.modal.Widget = func(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
		return h.menu.Layout(gtx, th)
	}
	h.modal.Appear(gtx.Now)
}

func (h *handler) showMetadataDialog(gtx layout.Context) {
	h.metadataDialog.Clear()
	for _, md := range h.metadata {
		h.metadataDialog.Add(md.k, md.v)
	}
	h.metadataDialog.OnClick = func(ok bool) {
		h.modal.Disappear(gtx.Now)
		if !ok {
			return
		}
		h.metadata = nil
		for _, kv := range h.metadataDialog.Metadata() {
			k, v := kv.Get()
			h.metadata = append(h.metadata, metadata{
				k: k,
				v: v,
			})
		}
		h.metadataSelector.Clear()
		h.metadataSelector.Select(ui_widget.SelectorItem{Value: strconv.Itoa(len(h.metadata))})
	}

	h.modal.Widget = func(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
		return h.metadataDialog.Layout(gtx, th)
	}
	h.modal.Appear(gtx.Now)
}
