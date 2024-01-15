package service

import (
	"image/color"
	"strings"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/go-gost/gui/api"
	"github.com/go-gost/gui/ui/icons"
	"github.com/go-gost/gui/ui/page"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

type metadata struct {
	k      component.TextField
	v      component.TextField
	delete widget.Clickable
}

type C = layout.Context
type D = layout.Dimensions

type servicePage struct {
	router *page.Router

	modal *component.ModalLayer
	menu  page.Menu
	list  layout.List

	btnBack   widget.Clickable
	btnDelete widget.Clickable
	btnEdit   widget.Clickable
	btnSave   widget.Clickable

	name component.TextField
	addr component.TextField

	admissions   []string
	btnAdmission widget.Clickable
	bypasses     []string
	btnBypass    widget.Clickable
	resolver     string
	btnResolver  widget.Clickable
	hosts        string
	btnHosts     widget.Clickable
	limiter      string
	btnLimiter   widget.Clickable
	loggers      []string
	btnLogger    widget.Clickable
	observer     string
	btnObserver  widget.Clickable

	enableStats widget.Bool

	handler  handler
	listener listener

	id            string
	edit          bool
	create        bool
	deleteConfirm bool
}

func NewPage(r *page.Router) page.Page {
	p := &servicePage{
		router: r,

		modal: component.NewModal(),
		menu: page.Menu{
			Surface: component.SurfaceStyle{
				Theme: r.Theme,
				ShadowStyle: component.ShadowStyle{
					CornerRadius: 12,
				},
				Fill: r.Theme.Bg,
			},
			List: layout.List{
				Axis: layout.Vertical,
			},
		},

		list: layout.List{
			// NOTE: the list must be vertical
			Axis: layout.Vertical,
		},
		name: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     32,
			},
		},
		addr: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     255,
			},
		},
	}
	p.handler.modal = p.modal
	p.handler.menu = &p.menu
	p.listener.modal = p.modal
	p.listener.menu = &p.menu

	return p
}

func (p *servicePage) Init(opts ...page.PageOption) {
	var options page.PageOptions
	for _, opt := range opts {
		opt(&options)
	}
	p.id = options.ID

	if p.id != "" {
		p.edit = false
		p.create = false
		p.name.ReadOnly = true
	} else {
		p.edit = true
		p.create = true
		p.name.ReadOnly = false
	}
	p.deleteConfirm = false

	cfg := api.GetConfig()
	var service *api.ServiceConfig
	for _, svc := range cfg.Services {
		if svc.Name == p.id {
			service = svc
			break
		}
	}
	if service == nil {
		service = &api.ServiceConfig{}
	}

	p.name.Clear()
	p.name.SetText(service.Name)
}

func (p *servicePage) Layout(gtx C) D {
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
		if p.deleteConfirm {
			p.delete()
			p.router.Back()
		} else {
			p.deleteConfirm = true
		}
	}

	th := p.router.Theme

	defer p.modal.Layout(gtx, th)

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		// header
		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Top:    5,
				Bottom: 5,
				Left:   10,
				Right:  10,
			}.Layout(gtx, func(gtx C) D {
				return layout.Flex{
					Spacing:   layout.SpaceBetween,
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						btn := material.IconButton(th, &p.btnBack, icons.IconBack, "Back")
						btn.Color = th.Fg
						btn.Background = th.Bg
						return btn.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: 10}.Layout),
					layout.Flexed(1, func(gtx C) D {
						title := material.H6(th, "Service")
						title.Font.Weight = font.Bold
						return title.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						if p.create {
							return D{}
						}
						btn := material.IconButton(th, &p.btnDelete, icons.IconDelete, "Delete")
						if p.deleteConfirm {
							btn = material.IconButton(th, &p.btnDelete, icons.IconDeleteForever, "Delete")
						}
						btn.Color = th.Fg
						btn.Background = th.Bg
						return btn.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
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
		layout.Flexed(1, func(gtx C) D {
			inset := layout.Inset{
				Top:    5,
				Bottom: 5,
			}
			width := unit.Dp(800)
			if x := gtx.Metric.PxToDp(gtx.Constraints.Max.X); x > width {
				inset.Left = (x - width) / 2
				inset.Right = inset.Left
			}
			return inset.Layout(gtx, func(gtx C) D {
				return p.list.Layout(gtx, 1, func(gtx C, index int) D {
					return layout.Inset{
						Top:    5,
						Bottom: 5,
						Left:   10,
						Right:  10,
					}.Layout(gtx, func(gtx C) D {
						return p.layout(gtx, th)
					})
				})
			})
		}),
	)
}

func (p *servicePage) layout(gtx C, th *material.Theme) D {
	if !p.edit {
		gtx = gtx.Disabled()
	}

	return component.SurfaceStyle{
		Theme: th,
		ShadowStyle: component.ShadowStyle{
			CornerRadius: 20,
		},
		Fill: color.NRGBA(colornames.Grey50),
	}.Layout(gtx, func(gtx C) D {
		return layout.Inset{
			Top:    10,
			Bottom: 10,
			Left:   10,
			Right:  10,
		}.Layout(gtx, func(gtx C) D {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(material.Body1(th, "Name").Layout),
				layout.Rigid(func(gtx C) D {
					return p.name.Layout(gtx, th, "")
				}),

				layout.Rigid(layout.Spacer{Height: 10}.Layout),
				layout.Rigid(material.Body1(th, "Addr").Layout),
				layout.Rigid(func(gtx C) D {
					return p.addr.Layout(gtx, th, "")
				}),

				layout.Rigid(layout.Spacer{Height: 10}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if p.btnAdmission.Clicked(gtx) {
						p.showAdmissionMenu(gtx)
					}

					return p.btnAdmission.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
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
									}.Layout(gtx, material.Body1(th, "Admission").Layout)
								}),
								layout.Flexed(1, layout.Spacer{Width: 5}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return material.Body2(th, strings.Join(p.admissions, ", ")).Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return icons.IconNavRight.Layout(gtx, th.Fg)
								}),
							)
						})
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if p.btnBypass.Clicked(gtx) {
						p.showBypassMenu(gtx)
					}

					return p.btnBypass.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
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
									}.Layout(gtx, material.Body1(th, "Bypass").Layout)
								}),
								layout.Flexed(1, layout.Spacer{Width: 5}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return material.Body2(th, strings.Join(p.bypasses, ", ")).Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return icons.IconNavRight.Layout(gtx, th.Fg)
								}),
							)
						})
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if p.btnResolver.Clicked(gtx) {
						p.showResolverMenu(gtx)
					}

					return p.btnResolver.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{
							Top:    8,
							Bottom: 8,
						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{
								Alignment: layout.Middle,
							}.Layout(gtx,
								layout.Flexed(1, material.Body1(th, "Resolver").Layout),
								layout.Rigid(material.Body2(th, p.resolver).Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return icons.IconNavRight.Layout(gtx, th.Fg)
								}),
							)
						})
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if p.btnHosts.Clicked(gtx) {
						p.showHostMapperMenu(gtx)
					}

					return p.btnHosts.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{
							Top:    8,
							Bottom: 8,
						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{
								Alignment: layout.Middle,
							}.Layout(gtx,
								layout.Flexed(1, material.Body1(th, "Host Mapper").Layout),
								layout.Rigid(material.Body2(th, p.hosts).Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return icons.IconNavRight.Layout(gtx, th.Fg)
								}),
							)
						})
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if p.btnLimiter.Clicked(gtx) {
						p.showLimiterMenu(gtx)
					}

					return p.btnLimiter.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{
							Top:    8,
							Bottom: 8,
						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{
								Alignment: layout.Middle,
							}.Layout(gtx,
								layout.Flexed(1, material.Body1(th, "Limiter").Layout),
								layout.Rigid(material.Body2(th, p.limiter).Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return icons.IconNavRight.Layout(gtx, th.Fg)
								}),
							)
						})
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if p.btnObserver.Clicked(gtx) {
						p.showObserverMenu(gtx)
					}

					return p.btnObserver.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{
							Top:    8,
							Bottom: 8,
						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{
								Alignment: layout.Middle,
							}.Layout(gtx,
								layout.Flexed(1, material.Body1(th, "Observer").Layout),
								layout.Rigid(material.Body2(th, p.observer).Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return icons.IconNavRight.Layout(gtx, th.Fg)
								}),
							)
						})
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if p.btnLogger.Clicked(gtx) {
						p.showLoggerMenu(gtx)
					}

					return p.btnLogger.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
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
									}.Layout(gtx, material.Body1(th, "Logger").Layout)
								}),
								layout.Flexed(1, layout.Spacer{Width: 5}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return material.Body2(th, strings.Join(p.loggers, ", ")).Layout(gtx)
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
						layout.Flexed(1, material.Body1(th, "Enable Stats").Layout),
						layout.Rigid(material.Switch(th, &p.enableStats, "").Layout),
					)
				}),

				layout.Rigid(layout.Spacer{Height: 10}.Layout),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Top:    20,
						Bottom: 20,
					}.Layout(gtx, material.H6(th, "Handler").Layout)
				}),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return p.handler.Layout(gtx, th)
				}),

				layout.Rigid(layout.Spacer{Height: 10}.Layout),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Top:    20,
						Bottom: 20,
					}.Layout(gtx, material.H6(th, "Listener").Layout)
				}),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return p.listener.Layout(gtx, th)
				}),
			)
		})
	})
}

func (p *servicePage) showAdmissionMenu(gtx C) {
	items := []page.MenuItem{}
	for _, v := range api.GetConfig().Admissions {
		items = append(items, page.MenuItem{
			Key:   v.Name,
			Value: v.Name,
		})
	}
	for i := range items {
		for _, v := range p.admissions {
			if items[i].Value == v {
				items[i].Selected = true
			}
		}
	}

	p.menu.Title = "Admission"
	p.menu.Items = items
	p.menu.Selected = func(index int) {
		p.admissions = nil
		for i := range p.menu.Items {
			if p.menu.Items[i].Selected {
				p.admissions = append(p.admissions, p.menu.Items[i].Value)
			}
		}
	}
	p.menu.ShowAdd = true
	p.menu.Multiple = true

	p.modal.Widget = func(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
		return p.menu.Layout(gtx, th)
	}
	p.modal.Appear(gtx.Now)
}

func (p *servicePage) showBypassMenu(gtx C) {
	items := []page.MenuItem{}
	for _, v := range api.GetConfig().Bypasses {
		items = append(items, page.MenuItem{
			Key:   v.Name,
			Value: v.Name,
		})
	}
	for i := range items {
		for _, v := range p.bypasses {
			if items[i].Value == v {
				items[i].Selected = true
			}
		}
	}

	p.menu.Title = "Bypass"
	p.menu.Items = items
	p.menu.Selected = func(index int) {
		p.bypasses = nil
		for i := range p.menu.Items {
			if p.menu.Items[i].Selected {
				p.bypasses = append(p.bypasses, p.menu.Items[i].Value)
			}
		}
	}
	p.menu.ShowAdd = true
	p.menu.Multiple = true

	p.modal.Widget = func(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
		return p.menu.Layout(gtx, th)
	}
	p.modal.Appear(gtx.Now)
}

func (p *servicePage) showResolverMenu(gtx C) {
	items := []page.MenuItem{
		{Key: "N/A"},
	}
	for _, v := range api.GetConfig().Resolvers {
		items = append(items, page.MenuItem{
			Key:   v.Name,
			Value: v.Name,
		})
	}
	for i := range items {
		if items[i].Value == p.resolver {
			items[i].Selected = true
		}
	}

	p.menu.Title = "Resolver"
	p.menu.Items = items
	p.menu.Selected = func(index int) {
		p.resolver = p.menu.Items[index].Value
		p.modal.Disappear(gtx.Now)
	}
	p.menu.ShowAdd = true

	p.modal.Widget = func(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
		return p.menu.Layout(gtx, th)
	}
	p.modal.Appear(gtx.Now)
}

func (p *servicePage) showHostMapperMenu(gtx C) {
	items := []page.MenuItem{
		{Key: "N/A"},
	}
	for _, v := range api.GetConfig().Hosts {
		items = append(items, page.MenuItem{
			Key:   v.Name,
			Value: v.Name,
		})
	}
	for i := range items {
		if items[i].Value == p.hosts {
			items[i].Selected = true
		}
	}

	p.menu.Title = "Host Mapper"
	p.menu.Items = items
	p.menu.Selected = func(index int) {
		p.hosts = p.menu.Items[index].Value
		p.modal.Disappear(gtx.Now)
	}
	p.menu.ShowAdd = true

	p.modal.Widget = func(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
		return p.menu.Layout(gtx, th)
	}
	p.modal.Appear(gtx.Now)
}

func (p *servicePage) showLimiterMenu(gtx C) {
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
		if items[i].Value == p.limiter {
			items[i].Selected = true
		}
	}

	p.menu.Title = "Limiter"
	p.menu.Items = items
	p.menu.Selected = func(index int) {
		p.limiter = p.menu.Items[index].Value
		p.modal.Disappear(gtx.Now)
	}
	p.menu.ShowAdd = true

	p.modal.Widget = func(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
		return p.menu.Layout(gtx, th)
	}
	p.modal.Appear(gtx.Now)
}

func (p *servicePage) showObserverMenu(gtx C) {
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
		if items[i].Value == p.observer {
			items[i].Selected = true
		}
	}

	p.menu.Title = "Observer"
	p.menu.Items = items
	p.menu.Selected = func(index int) {
		p.observer = p.menu.Items[index].Value
		p.modal.Disappear(gtx.Now)
	}
	p.menu.ShowAdd = true

	p.modal.Widget = func(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
		return p.menu.Layout(gtx, th)
	}
	p.modal.Appear(gtx.Now)
}

func (p *servicePage) showLoggerMenu(gtx C) {
	items := []page.MenuItem{}
	for _, v := range api.GetConfig().Loggers {
		items = append(items, page.MenuItem{
			Key:   v.Name,
			Value: v.Name,
		})
	}
	for i := range items {
		for _, v := range p.loggers {
			if items[i].Value == v {
				items[i].Selected = true
			}
		}
	}

	p.menu.Title = "Logger"
	p.menu.Items = items
	p.menu.Selected = func(index int) {
		p.loggers = nil
		for i := range p.menu.Items {
			if p.menu.Items[i].Selected {
				p.loggers = append(p.loggers, p.menu.Items[i].Value)
			}
		}
	}
	p.menu.ShowAdd = true
	p.menu.Multiple = true

	p.modal.Widget = func(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
		return p.menu.Layout(gtx, th)
	}
	p.modal.Appear(gtx.Now)
}

func (p *servicePage) save() bool {
	return false
}

func (p *servicePage) delete() {
}
