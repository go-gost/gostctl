package nameserver

import (
	"fmt"
	"strings"

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
	"github.com/google/uuid"
)

type nameserverPage struct {
	router *page.Router

	menu ui_widget.Menu
	mode widget.Enum
	list layout.List

	btnBack   widget.Clickable
	btnDelete widget.Clickable
	btnEdit   widget.Clickable
	btnSave   widget.Clickable

	// name component.TextField
	addr    component.TextField
	ttl     component.TextField
	timeout component.TextField
	async   ui_widget.Switcher

	chain    ui_widget.Selector
	prefer   ui_widget.Selector
	only     ui_widget.Selector
	clientIP component.TextField
	hostname component.TextField

	id       string
	perm     page.Perm
	callback page.Callback

	edit   bool
	create bool

	delDialog ui_widget.Dialog
}

func NewPage(r *page.Router) page.Page {
	p := &nameserverPage{
		router: r,

		list: layout.List{
			// NOTE: the list must be vertical
			Axis: layout.Vertical,
		},

		addr: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     255,
			},
		},
		ttl: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     16,
				Filter:     "1234567890",
			},
			Suffix: material.Body1(r.Theme, i18n.TimeSecond.Value()).Layout,
		},
		timeout: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     16,
				Filter:     "1234567890",
			},
			Suffix: material.Body1(r.Theme, i18n.TimeSecond.Value()).Layout,
		},
		async:  ui_widget.Switcher{Title: i18n.Async},
		chain:  ui_widget.Selector{Title: i18n.Chain},
		prefer: ui_widget.Selector{Title: i18n.Prefer},
		only:   ui_widget.Selector{Title: i18n.Only},
		clientIP: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     64,
			},
		},
		hostname: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     64,
			},
		},

		delDialog: ui_widget.Dialog{Title: i18n.DeleteNameserver},
	}

	return p
}

func (p *nameserverPage) Init(opts ...page.PageOption) {
	var options page.PageOptions
	for _, opt := range opts {
		opt(&options)
	}

	p.id = options.ID
	nameserver, _ := options.Value.(*api.NameserverConfig)
	if nameserver == nil {
		nameserver = &api.NameserverConfig{}
	}
	p.callback = options.Callback

	if p.id != "" {
		p.edit = false
		p.create = false
	} else {
		p.edit = true
		p.create = true
	}

	p.perm = options.Perm
	p.mode.Value = string(page.BasicMode)

	p.addr.SetText(nameserver.Addr)

	p.ttl.Clear()
	p.ttl.SetText(fmt.Sprintf("%d", int(nameserver.TTL.Seconds())))

	p.timeout.Clear()
	p.timeout.SetText(fmt.Sprintf("%d", int(nameserver.Timeout.Seconds())))

	p.async.SetValue(nameserver.Async)

	p.chain.Clear()
	p.chain.Select(ui_widget.SelectorItem{Value: nameserver.Chain})

	p.prefer.Clear()
	for _, v := range page.IPTypeOptions {
		if v.Value == nameserver.Prefer {
			p.prefer.Select(ui_widget.SelectorItem{Name: v.Name, Value: v.Value})
		}
	}
	p.only.Clear()
	for _, v := range page.IPTypeOptions {
		if v.Value == nameserver.Prefer {
			p.only.Select(ui_widget.SelectorItem{Name: v.Name, Value: v.Value})
		}
	}
	p.clientIP.SetText(nameserver.ClientIP)
	p.hostname.SetText(nameserver.Hostname)

}

func (p *nameserverPage) Layout(gtx page.C) page.D {
	if p.btnBack.Clicked(gtx) {
		p.router.Back()
	}
	if p.btnEdit.Clicked(gtx) {
		p.edit = true
	}
	if p.btnSave.Clicked(gtx) {
		if p.save() {
			p.router.Back()
		}
	}

	if p.btnDelete.Clicked(gtx) {
		p.delDialog.OnClick = func(ok bool) {
			if ok {
				p.delete()
				p.router.Back()
			}
			p.router.HideModal(gtx)
		}
		p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
			return p.delDialog.Layout(gtx, th)
		})
	}

	th := p.router.Theme

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		// header
		layout.Rigid(func(gtx page.C) page.D {
			return layout.Inset{
				Top:    8,
				Bottom: 8,
				Left:   8,
				Right:  8,
			}.Layout(gtx, func(gtx page.C) page.D {
				return layout.Flex{
					Spacing:   layout.SpaceBetween,
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(func(gtx page.C) page.D {
						btn := material.IconButton(th, &p.btnBack, icons.IconBack, "Back")
						btn.Color = th.Fg
						btn.Background = th.Bg
						return btn.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: 8}.Layout),
					layout.Flexed(1, func(gtx page.C) page.D {
						title := material.H6(th, i18n.Nameserver.Value())
						return title.Layout(gtx)
					}),
					layout.Rigid(func(gtx page.C) page.D {
						if p.perm&page.PermDelete == 0 || p.create {
							return page.D{}
						}
						btn := material.IconButton(th, &p.btnDelete, icons.IconDelete, "Delete")
						btn.Color = th.Fg
						btn.Background = th.Bg
						return btn.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: 8}.Layout),
					layout.Rigid(func(gtx page.C) page.D {
						if p.perm&page.PermWrite == 0 {
							return page.D{}
						}

						if p.edit {
							btn := material.IconButton(th, &p.btnSave, icons.IconDone, "Done")
							btn.Color = th.Fg
							btn.Background = th.Bg
							return btn.Layout(gtx)
						} else {
							btn := material.IconButton(th, &p.btnEdit, icons.IconEdit, "Edit")
							btn.Color = th.Fg
							btn.Background = th.Bg
							return btn.Layout(gtx)
						}
					}),
				)
			})
		}),
		layout.Flexed(1, func(gtx page.C) page.D {
			return p.list.Layout(gtx, 1, func(gtx page.C, index int) page.D {
				return layout.Inset{
					Top:    8,
					Bottom: 8,
					Left:   8,
					Right:  8,
				}.Layout(gtx, func(gtx page.C) page.D {
					return p.layout(gtx, th)
				})
			})
		}),
	)
}

func (p *nameserverPage) layout(gtx page.C, th *page.T) page.D {
	src := gtx.Source

	if !p.edit {
		gtx = gtx.Disabled()
	}

	return component.SurfaceStyle{
		Theme: th,
		ShadowStyle: component.ShadowStyle{
			CornerRadius: 12,
		},
		Fill: theme.Current().ContentSurfaceBg,
	}.Layout(gtx, func(gtx page.C) page.D {
		return layout.UniformInset(16).Layout(gtx, func(gtx page.C) page.D {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Flexed(1, layout.Spacer{Width: 8}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							gtx.Source = src
							return material.RadioButton(th, &p.mode, string(page.BasicMode), i18n.Basic.Value()).Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: 8}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							gtx.Source = src
							return material.RadioButton(th, &p.mode, string(page.AdvancedMode), i18n.Advanced.Value()).Layout(gtx)
						}),
					)
				}),

				layout.Rigid(material.Body1(th, i18n.Address.Value()).Layout),
				layout.Rigid(func(gtx page.C) page.D {
					return p.addr.Layout(gtx, th, "")
				}),
				layout.Rigid(layout.Spacer{Height: 4}.Layout),

				layout.Rigid(material.Body1(th, "TTL").Layout),
				layout.Rigid(func(gtx page.C) page.D {
					return p.ttl.Layout(gtx, th, "")
				}),
				layout.Rigid(layout.Spacer{Height: 4}.Layout),

				layout.Rigid(material.Body1(th, i18n.Timeout.Value()).Layout),
				layout.Rigid(func(gtx page.C) page.D {
					return p.timeout.Layout(gtx, th, "")
				}),

				layout.Rigid(layout.Spacer{Height: 4}.Layout),
				layout.Rigid(func(gtx page.C) page.D {
					return p.async.Layout(gtx, th)
				}),

				layout.Rigid(func(gtx page.C) page.D {
					if p.prefer.Clicked(gtx) {
						p.showPreferMenu(gtx)
					}
					return p.prefer.Layout(gtx, th)
				}),

				// advanced mode
				layout.Rigid(func(gtx page.C) page.D {
					if p.mode.Value == string(page.BasicMode) {
						return page.D{}
					}

					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(gtx,
						layout.Rigid(func(gtx page.C) page.D {
							if p.only.Clicked(gtx) {
								p.showOnlyMenu(gtx)
							}
							return p.only.Layout(gtx, th)
						}),

						layout.Rigid(func(gtx page.C) page.D {
							if p.chain.Clicked(gtx) {
								p.showChainMenu(gtx)
							}
							return p.chain.Layout(gtx, th)
						}),

						layout.Rigid(layout.Spacer{Height: 8}.Layout),
						layout.Rigid(material.Body1(th, i18n.ClientIP.Value()).Layout),
						layout.Rigid(func(gtx page.C) page.D {
							return p.clientIP.Layout(gtx, th, "")
						}),

						layout.Rigid(layout.Spacer{Height: 4}.Layout),
						layout.Rigid(material.Body1(th, i18n.Hostname.Value()).Layout),
						layout.Rigid(func(gtx page.C) page.D {
							return p.hostname.Layout(gtx, th, "")
						}),
					)
				}),
			)
		})
	})
}

func (p *nameserverPage) showChainMenu(gtx page.C) {
	options := []ui_widget.MenuOption{}
	for _, v := range api.GetConfig().Chains {
		options = append(options, ui_widget.MenuOption{
			Value: v.Name,
		})
	}
	for i := range options {
		options[i].Selected = p.chain.AnyValue(options[i].Value)
	}

	p.menu.Title = i18n.Chain
	p.menu.Options = options
	p.menu.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		p.chain.Clear()
		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				p.chain.Select(ui_widget.SelectorItem{Value: p.menu.Options[i].Value})
			}
		}
	}
	p.menu.OnAdd = func() {
		p.router.Goto(page.Route{
			Path: page.PageChain,
			Perm: page.PermReadWrite,
		})
		p.router.HideModal(gtx)
	}
	p.menu.Multiple = false

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *nameserverPage) showPreferMenu(gtx page.C) {
	for i := range page.IPTypeOptions {
		page.IPTypeOptions[i].Selected = p.prefer.AnyValue(page.IPTypeOptions[i].Value)
	}

	p.menu.Title = i18n.Prefer
	p.menu.Options = page.IPTypeOptions
	p.menu.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		p.prefer.Clear()
		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				p.prefer.Select(ui_widget.SelectorItem{Name: p.menu.Options[i].Name, Key: p.menu.Options[i].Key, Value: p.menu.Options[i].Value})
			}
		}
	}
	p.menu.OnAdd = nil
	p.menu.Multiple = false

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *nameserverPage) showOnlyMenu(gtx page.C) {
	for i := range page.IPTypeOptions {
		page.IPTypeOptions[i].Selected = p.only.AnyValue(page.IPTypeOptions[i].Value)
	}

	p.menu.Title = i18n.Only
	p.menu.Options = page.IPTypeOptions
	p.menu.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		p.only.Clear()
		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				p.only.Select(ui_widget.SelectorItem{Name: p.menu.Options[i].Name, Key: p.menu.Options[i].Key, Value: p.menu.Options[i].Value})
			}
		}
	}
	p.menu.OnAdd = nil
	p.menu.Multiple = false

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *nameserverPage) generateConfig() *api.NameserverConfig {
	nameserver := &api.NameserverConfig{
		Addr: strings.TrimSpace(p.addr.Text()),
	}

	return nameserver
}

func (p *nameserverPage) save() bool {
	nameserver := p.generateConfig()

	if p.id == "" {
		if p.callback != nil {
			p.callback(page.ActionCreate, uuid.New().String(), nameserver)
		}

	} else {
		if p.callback != nil {
			p.callback(page.ActionUpdate, p.id, nameserver)
		}
	}

	return true
}

func (p *nameserverPage) delete() {
	if p.callback != nil {
		p.callback(page.ActionDelete, p.id, nil)
	}
}
