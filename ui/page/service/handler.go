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

type handler struct {
	service *servicePage
	menu    ui_widget.Menu

	typ   ui_widget.Selector
	chain ui_widget.Selector

	authType widget.Enum
	username component.TextField
	password component.TextField
	auther   ui_widget.Selector

	limiter  ui_widget.Selector
	observer ui_widget.Selector

	metadata          []metadata
	mdSelector        ui_widget.Selector
	mdFolded          bool
	mdAdd             widget.Clickable
	mdDialog          ui_widget.MetadataDialog
	delMetadataDialog ui_widget.Dialog
}

func newHandler(service *servicePage) *handler {
	return &handler{
		service:    service,
		typ:        ui_widget.Selector{Title: i18n.Type},
		chain:      ui_widget.Selector{Title: i18n.Chain},
		auther:     ui_widget.Selector{Title: i18n.Auther},
		limiter:    ui_widget.Selector{Title: i18n.Limiter},
		observer:   ui_widget.Selector{Title: i18n.Observer},
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

func (h *handler) init(cfg *api.HandlerConfig) {
	if cfg == nil {
		cfg = &api.HandlerConfig{}
	}

	h.typ.Clear()
	for i := range handlerTypeAdvancedOptions {
		if handlerTypeAdvancedOptions[i].Value == cfg.Type {
			h.typ.Select(ui_widget.SelectorItem{Name: handlerTypeAdvancedOptions[i].Name, Key: handlerTypeAdvancedOptions[i].Key, Value: handlerTypeAdvancedOptions[i].Value})
			break
		}
	}

	h.chain.Clear()
	h.chain.Select(ui_widget.SelectorItem{Value: cfg.Chain})

	{
		h.username.Clear()
		h.password.Clear()
		h.authType.Value = ""

		if cfg.Auth != nil {
			h.username.SetText(cfg.Auth.Username)
			h.password.SetText(cfg.Auth.Password)
			h.authType.Value = string(page.AuthSimple)
		}

		h.auther.Clear()
		var items []ui_widget.SelectorItem
		if cfg.Auther != "" {
			items = append(items, ui_widget.SelectorItem{Value: cfg.Auther})
		}
		for _, v := range cfg.Authers {
			items = append(items, ui_widget.SelectorItem{Value: v})
		}
		h.auther.Select(items...)

		if len(cfg.Authers) > 0 || cfg.Auther != "" {
			h.authType.Value = string(page.AuthAuther)
		}
	}

	h.metadata = nil
	md := api.NewMetadata(cfg.Metadata)
	for k := range md {
		if k == "" {
			continue
		}
		h.metadata = append(h.metadata, metadata{
			k: k,
			v: md.GetString(k),
		})
	}
	h.mdSelector.Clear()
	h.mdSelector.Select(ui_widget.SelectorItem{Value: strconv.Itoa(len(h.metadata))})
	h.mdFolded = true
}

func (h *handler) Layout(gtx page.C, th *page.T) page.D {
	src := gtx.Source

	if !h.service.edit {
		gtx = gtx.Disabled()
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx page.C) page.D {
			if h.typ.Clicked(gtx) {
				h.showTypeMenu(gtx)
			}

			return h.typ.Layout(gtx, th)
		}),

		layout.Rigid(func(gtx page.C) page.D {
			if h.chain.Clicked(gtx) {
				h.showChainMenu(gtx)
			}

			return h.chain.Layout(gtx, th)
		}),

		// auth for handler
		layout.Rigid(func(gtx page.C) page.D {
			if !h.canAuth() {
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
								return material.RadioButton(th, &h.authType, string(page.AuthSimple), i18n.AuthSimple.Value()).Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Width: 8}.Layout),
							layout.Rigid(func(gtx page.C) page.D {
								return material.RadioButton(th, &h.authType, string(page.AuthAuther), i18n.AuthAuther.Value()).Layout(gtx)
							}),
						)
					})
				}),
				layout.Rigid(func(gtx page.C) page.D {
					if h.authType.Value != string(page.AuthSimple) {
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
								return h.username.Layout(gtx, th, "")
							}),
							layout.Rigid(layout.Spacer{Height: 8}.Layout),

							layout.Rigid(func(gtx page.C) page.D {
								return material.Body1(th, i18n.Password.Value()).Layout(gtx)
							}),
							layout.Rigid(func(gtx page.C) page.D {
								return h.password.Layout(gtx, th, "")
							}),
							layout.Rigid(layout.Spacer{Height: 8}.Layout),
						)
					})
				}),

				layout.Rigid(func(gtx page.C) page.D {
					if h.authType.Value != string(page.AuthAuther) {
						return page.D{}
					}

					if h.auther.Clicked(gtx) {
						h.showAutherMenu(gtx)
					}

					return layout.UniformInset(8).Layout(gtx, func(gtx page.C) page.D {
						return h.auther.Layout(gtx, th)
					})
				}),
			)
		}),

		// advanced mode
		layout.Rigid(func(gtx page.C) page.D {
			if h.service.mode.Value == string(page.BasicMode) {
				return page.D{}
			}

			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx page.C) page.D {
					if h.limiter.Clicked(gtx) {
						h.showLimiterMenu(gtx)
					}
					return h.limiter.Layout(gtx, th)
				}),

				layout.Rigid(func(gtx page.C) page.D {
					if h.observer.Clicked(gtx) {
						h.showObserverMenu(gtx)
					}
					return h.observer.Layout(gtx, th)
				}),
			)
		}),

		layout.Rigid(func(gtx page.C) page.D {
			if h.mdAdd.Clicked(gtx) {
				h.showMetadataDialog(gtx, -1)
			}

			return layout.Flex{
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Flexed(1, func(gtx page.C) page.D {
					gtx.Source = src

					if h.mdSelector.Clicked(gtx) {
						h.mdFolded = !h.mdFolded
					}
					return h.mdSelector.Layout(gtx, th)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if !h.service.edit {
						return page.D{}
					}
					return layout.Spacer{Width: 8}.Layout(gtx)
				}),
				layout.Rigid(func(gtx page.C) page.D {
					if !h.service.edit {
						return page.D{}
					}
					btn := material.IconButton(th, &h.mdAdd, icons.IconAdd, "Add")
					btn.Background = theme.Current().ContentSurfaceBg
					btn.Color = th.Fg
					// btn.Inset = layout.UniformInset(8)
					return btn.Layout(gtx)
				}),
			)
		}),
		layout.Rigid(func(gtx page.C) page.D {
			if h.mdFolded {
				return page.D{}
			}

			gtx.Source = src
			return h.layoutMetadata(gtx, th)
		}),
	)
}

func (h *handler) layoutMetadata(gtx page.C, th *page.T) page.D {
	for i := range h.metadata {
		if h.metadata[i].clk.Clicked(gtx) {
			if h.service.edit {
				h.showMetadataDialog(gtx, i)
			}
			break
		}

		if h.metadata[i].delete.Clicked(gtx) {
			h.delMetadataDialog.OnClick = func(ok bool) {
				h.service.router.HideModal(gtx)
				if !ok {
					return
				}
				h.metadata = append(h.metadata[:i], h.metadata[i+1:]...)

				h.mdSelector.Clear()
				h.mdSelector.Select(ui_widget.SelectorItem{Value: strconv.Itoa(len(h.metadata))})
			}
			h.service.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
				return h.delMetadataDialog.Layout(gtx, th)
			})
			break
		}
	}

	var children []layout.FlexChild
	for i := range h.metadata {
		md := &h.metadata[i]

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
						if !h.service.edit {
							return page.D{}
						}
						return layout.Spacer{Width: 8}.Layout(gtx)
					}),
					layout.Rigid(func(gtx page.C) page.D {
						if !h.service.edit {
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

func (h *handler) showMetadataDialog(gtx page.C, i int) {
	h.mdDialog.K.Clear()
	h.mdDialog.V.Clear()

	if i >= 0 && i < len(h.metadata) {
		h.mdDialog.K.SetText(h.metadata[i].k)
		h.mdDialog.V.SetText(h.metadata[i].v)
	}

	h.mdDialog.OnClick = func(ok bool) {
		h.service.router.HideModal(gtx)
		if !ok {
			return
		}

		k, v := strings.TrimSpace(h.mdDialog.K.Text()), strings.TrimSpace(h.mdDialog.V.Text())
		if k == "" {
			return
		}

		if i >= 0 && i < len(h.metadata) {
			h.metadata[i].k = k
			h.metadata[i].v = v
		} else {
			h.metadata = append(h.metadata, metadata{
				k: k,
				v: v,
			})
		}

		h.mdSelector.Clear()
		h.mdSelector.Select(ui_widget.SelectorItem{Value: strconv.Itoa(len(h.metadata))})
	}

	h.service.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return h.mdDialog.Layout(gtx, th)
	})
}

var (
	handlerTypeOptions = []ui_widget.MenuOption{
		{Name: "Auto", Value: "auto"},
		{Name: "HTTP", Value: "http"},
		{Name: "SOCKS4", Value: "socks4"},
		{Name: "SOCKS5", Value: "socks5"},
		{Name: "Relay", Value: "relay"},
		{Name: "Shadowsocks", Value: "ss"},
		{Name: "HTTP/2", Value: "http2"},

		{Name: "TCP", Value: "tcp"},
		{Name: "UDP", Value: "udp"},
		{Name: "RTCP", Value: "rtcp"},
		{Name: "RUDP", Value: "rudp"},
	}

	handlerTypeAdvancedOptions = []ui_widget.MenuOption{
		{Name: "Auto", Value: "auto"},
		{Name: "HTTP", Value: "http"},
		{Name: "SOCKS4", Value: "socks4"},
		{Name: "SOCKS5", Value: "socks5"},
		{Name: "Relay", Value: "relay"},
		{Name: "Shadowsocks", Value: "ss"},
		{Name: "HTTP/2", Value: "http2"},

		{Name: "TCP", Value: "tcp"},
		{Name: "UDP", Value: "udp"},
		{Name: "RTCP", Value: "rtcp"},
		{Name: "RUDP", Value: "rudp"},

		{Name: "SNI", Value: "sni"},
		{Name: "DNS", Value: "dns"},
		{Name: "SSHD", Value: "sshd"},
		{Name: "HTTP/3", Value: "http3"},

		{Name: "RED", Value: "red"},
		{Name: "REDU", Value: "redu"},
		{Name: "TUN", Value: "tun"},
		{Name: "TAP", Value: "tap"},
		{Key: "Tunnel", Value: "tunnel"},

		{Key: "File", Value: "file"},
		{Key: "Serial", Value: "serial"},
		{Key: "Unix", Value: "unix"},
	}
)

func (h *handler) canAuth() bool {
	return h.typ.AnyValue("auto", "http", "http2", "socks4", "socks5", "socks", "relay", "ss", "file")
}

func (h *handler) canForward() bool {
	return h.typ.AnyValue("tcp", "udp", "rtcp", "rudp", "forward", "rforward", "dns", "serial", "unix")
}

func (h *handler) showTypeMenu(gtx page.C) {
	options := handlerTypeOptions
	if h.service.mode.Value == string(page.AdvancedMode) {
		options = handlerTypeAdvancedOptions
	}

	for i := range options {
		options[i].Selected = h.typ.AnyValue(options[i].Value)
	}

	h.menu.Title = i18n.Handler
	h.menu.Options = options
	h.menu.OnClick = func(ok bool) {
		h.service.router.HideModal(gtx)
		if !ok {
			return
		}

		h.typ.Clear()
		for i := range h.menu.Options {
			if h.menu.Options[i].Selected {
				h.typ.Select(ui_widget.SelectorItem{Name: h.menu.Options[i].Name, Key: h.menu.Options[i].Key, Value: h.menu.Options[i].Value})
			}
		}
	}
	h.menu.OnAdd = nil
	h.menu.Multiple = false

	h.service.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return h.menu.Layout(gtx, th)
	})
}

func (h *handler) showChainMenu(gtx page.C) {
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
		h.service.router.HideModal(gtx)
		if !ok {
			return
		}

		h.chain.Clear()
		for i := range h.menu.Options {
			if h.menu.Options[i].Selected {
				h.chain.Select(ui_widget.SelectorItem{Value: h.menu.Options[i].Value})
			}
		}
	}
	h.menu.OnAdd = func() {
		h.service.router.Goto(page.Route{
			Path: page.PageChain,
			Perm: page.PermReadWrite,
		})
		h.service.router.HideModal(gtx)
	}
	h.menu.Multiple = false

	h.service.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return h.menu.Layout(gtx, th)
	})
}

func (h *handler) showAutherMenu(gtx page.C) {
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
		h.service.router.HideModal(gtx)
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
	h.menu.OnAdd = func() {
		h.service.router.Goto(page.Route{
			Path: page.PageAuther,
			Perm: page.PermReadWrite,
		})
		h.service.router.HideModal(gtx)
	}
	h.menu.Multiple = true

	h.service.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return h.menu.Layout(gtx, th)
	})
}

func (h *handler) showLimiterMenu(gtx page.C) {
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
		h.service.router.HideModal(gtx)
		if !ok {
			return
		}

		h.limiter.Clear()
		for i := range h.menu.Options {
			if h.menu.Options[i].Selected {
				h.limiter.Select(ui_widget.SelectorItem{Value: h.menu.Options[i].Value})
			}
		}
	}
	h.menu.OnAdd = func() {}
	h.menu.Multiple = false

	h.service.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return h.menu.Layout(gtx, th)
	})
}

func (h *handler) showObserverMenu(gtx page.C) {
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
		h.service.router.HideModal(gtx)
		if !ok {
			return
		}

		h.observer.Clear()
		for i := range h.menu.Options {
			if h.menu.Options[i].Selected {
				h.observer.Select(ui_widget.SelectorItem{Value: h.menu.Options[i].Value})
			}
		}
	}
	h.menu.OnAdd = func() {
		h.service.router.Goto(page.Route{
			Path: page.PageObserver,
			Perm: page.PermReadWrite,
		})
		h.service.router.HideModal(gtx)
	}
	h.menu.Multiple = false

	h.service.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return h.menu.Layout(gtx, th)
	})
}
