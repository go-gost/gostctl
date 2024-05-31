package service

import (
	"context"
	"strconv"
	"strings"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/go-gost/gostctl/api"
	"github.com/go-gost/gostctl/api/runner"
	"github.com/go-gost/gostctl/api/runner/task"
	"github.com/go-gost/gostctl/api/util"
	"github.com/go-gost/gostctl/config"
	"github.com/go-gost/gostctl/ui/i18n"
	"github.com/go-gost/gostctl/ui/icons"
	"github.com/go-gost/gostctl/ui/page"
	"github.com/go-gost/gostctl/ui/theme"
	ui_widget "github.com/go-gost/gostctl/ui/widget"
	"github.com/google/uuid"
)

type record struct {
	id     string
	cfg    *api.RecorderObject
	clk    widget.Clickable
	delete widget.Clickable
}

type metadata struct {
	k      string
	v      string
	clk    widget.Clickable
	delete widget.Clickable
}

type servicePage struct {
	readonly bool

	router *page.Router

	menu ui_widget.Menu
	mode widget.Enum
	list layout.List

	btnBack   widget.Clickable
	btnDelete widget.Clickable
	btnEdit   widget.Clickable
	btnSave   widget.Clickable

	name component.TextField
	addr component.TextField

	// stats ui_widget.Switcher

	// customInterface ui_widget.Selector
	// interfaceDialog ui_widget.InputDialog

	admission  ui_widget.Selector
	bypass     ui_widget.Selector
	resolver   ui_widget.Selector
	hostMapper ui_widget.Selector
	limiter    ui_widget.Selector
	observer   ui_widget.Selector

	records          []record
	recorderSelector ui_widget.Selector
	recordFolded     bool
	recordAdd        widget.Clickable

	id   string
	perm page.Perm

	edit   bool
	create bool

	metadata   []metadata
	mdSelector ui_widget.Selector
	mdFolded   bool
	mdAdd      widget.Clickable
	mdDialog   ui_widget.MetadataDialog

	delDialog         ui_widget.Dialog
	delMetadataDialog ui_widget.Dialog
	delRecordDialog   ui_widget.Dialog

	handler   *handler
	listener  *listener
	forwarder *forwarder
}

func NewPage(r *page.Router) page.Page {
	p := &servicePage{
		router: r,

		list: layout.List{
			// NOTE: the list must be vertical
			Axis: layout.Vertical,
		},
		name: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     128,
			},
		},
		addr: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     255,
			},
		},
		delDialog:         ui_widget.Dialog{Title: i18n.DeleteService},
		delMetadataDialog: ui_widget.Dialog{Title: i18n.DeleteMetadata},
		delRecordDialog:   ui_widget.Dialog{Title: i18n.DeleteMetadata},

		// stats:           ui_widget.Switcher{Title: "Stats"},
		/*
			customInterface: ui_widget.Selector{Title: i18n.Interface},
			interfaceDialog: ui_widget.InputDialog{
				Title: i18n.Interface,
				Hint:  i18n.InterfaceHint,
				Input: component.TextField{
					Editor: widget.Editor{
						SingleLine: true,
						MaxLen:     255,
					},
				},
			},
		*/

		admission:        ui_widget.Selector{Title: i18n.Admission},
		bypass:           ui_widget.Selector{Title: i18n.Bypass},
		resolver:         ui_widget.Selector{Title: i18n.Resolver},
		hostMapper:       ui_widget.Selector{Title: i18n.Hosts},
		limiter:          ui_widget.Selector{Title: i18n.Limiter},
		observer:         ui_widget.Selector{Title: i18n.Observer},
		recorderSelector: ui_widget.Selector{Title: i18n.Recorder},

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
	}

	p.handler = newHandler(p)
	p.listener = newListener(p)
	p.forwarder = newForwarder(p)

	return p
}

func (p *servicePage) Init(opts ...page.PageOption) {
	if server := config.CurrentServer(); server != nil {
		p.readonly = server.Readonly
	}

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

	p.perm = options.Perm

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

	p.mode.Value = string(page.BasicMode)

	p.name.Clear()
	p.name.SetText(service.Name)

	p.addr.Clear()
	p.addr.SetText(service.Addr)

	// md := api.NewMetadata(service.Metadata)
	// p.stats.SetValue(md.GetBool("enableStats"))

	// p.customInterface.Select(ui_widget.SelectorItem{Value: service.Interface})

	{
		p.admission.Clear()
		var items []ui_widget.SelectorItem
		if service.Admission != "" {
			items = append(items, ui_widget.SelectorItem{Value: service.Admission})
		}
		for _, v := range service.Admissions {
			items = append(items, ui_widget.SelectorItem{
				Value: v,
			})
		}
		p.admission.Select(items...)
	}

	{
		p.bypass.Clear()
		var items []ui_widget.SelectorItem
		if service.Bypass != "" {
			items = append(items, ui_widget.SelectorItem{Value: service.Bypass})
		}
		for _, v := range service.Bypasses {
			items = append(items, ui_widget.SelectorItem{
				Value: v,
			})
		}
		p.bypass.Select(items...)
	}

	p.resolver.Clear()
	if service.Resolver != "" {
		p.resolver.Select(ui_widget.SelectorItem{Value: service.Resolver})
	}

	p.hostMapper.Clear()
	if service.Hosts != "" {
		p.hostMapper.Select(ui_widget.SelectorItem{Value: service.Hosts})
	}

	p.limiter.Clear()
	if service.Limiter != "" {
		p.limiter.Select(ui_widget.SelectorItem{Value: service.Limiter})
	}

	p.observer.Clear()
	if service.Observer != "" {
		p.observer.Select(ui_widget.SelectorItem{Value: service.Observer})
	}

	p.records = nil
	for _, v := range service.Recorders {
		if v == nil {
			continue
		}

		record := record{
			id:  uuid.New().String(),
			cfg: v,
		}

		p.records = append(p.records, record)
	}
	p.recorderSelector.Clear()
	p.recorderSelector.Select(ui_widget.SelectorItem{Value: strconv.Itoa(len(p.records))})
	p.recordFolded = true

	p.metadata = nil
	md := api.NewMetadata(service.Metadata)
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

	p.handler.init(service.Handler)
	p.listener.init(service.Listener)
	p.forwarder.init(service.Forwarder)
}

func (p *servicePage) Layout(gtx page.C) page.D {
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
						title := material.H6(th, i18n.Service.Value())
						return title.Layout(gtx)
					}),
					layout.Rigid(func(gtx page.C) page.D {
						if p.readonly || p.perm&page.PermDelete == 0 || p.create {
							return page.D{}
						}
						btn := material.IconButton(th, &p.btnDelete, icons.IconDelete, "Delete")
						btn.Color = th.Fg
						btn.Background = th.Bg
						return btn.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: 8}.Layout),
					layout.Rigid(func(gtx page.C) page.D {
						if p.readonly || p.perm&page.PermWrite == 0 {
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
				return layout.UniformInset(8).Layout(gtx, func(gtx page.C) page.D {
					return p.layout(gtx, th)
				})
			})
		}),
	)
}

func (p *servicePage) layout(gtx page.C, th *page.T) page.D {
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
				layout.Rigid(func(gtx page.C) page.D {
					gtx.Source = src
					return layout.Flex{
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Rigid(func(gtx page.C) page.D {
							return material.RadioButton(th, &p.mode, string(page.BasicMode), i18n.Basic.Value()).Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: 4}.Layout),
						layout.Rigid(func(gtx page.C) page.D {
							return material.RadioButton(th, &p.mode, string(page.AdvancedMode), i18n.Advanced.Value()).Layout(gtx)
						}),
					)
				}),

				layout.Rigid(layout.Spacer{Height: 16}.Layout),

				layout.Rigid(material.Body1(th, i18n.Name.Value()).Layout),
				layout.Rigid(func(gtx page.C) page.D {
					return p.name.Layout(gtx, th, "")
				}),
				layout.Rigid(layout.Spacer{Height: 4}.Layout),

				layout.Rigid(material.Body1(th, i18n.Address.Value()).Layout),
				layout.Rigid(func(gtx page.C) page.D {
					return p.addr.Layout(gtx, th, "")
				}),

				layout.Rigid(layout.Spacer{Height: 4}.Layout),

				/*
					layout.Rigid(func(gtx page.C) page.D {
						return p.stats.Layout(gtx, th)
					}),
				*/

				// advanced mode
				layout.Rigid(func(gtx page.C) page.D {
					if p.mode.Value == string(page.BasicMode) {
						return page.D{}
					}

					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(gtx,
						/*
							layout.Rigid(func(gtx page.C) page.D {
								if p.customInterface.Clicked(gtx) {
									p.showInterfaceDialog(gtx)
								}
								return p.customInterface.Layout(gtx, th)
							}),
						*/

						layout.Rigid(func(gtx page.C) page.D {
							if p.admission.Clicked(gtx) {
								p.showAdmissionMenu(gtx)
							}
							return p.admission.Layout(gtx, th)
						}),

						layout.Rigid(func(gtx page.C) page.D {
							if p.bypass.Clicked(gtx) {
								p.showBypassMenu(gtx)
							}
							return p.bypass.Layout(gtx, th)
						}),

						layout.Rigid(func(gtx page.C) page.D {
							if p.resolver.Clicked(gtx) {
								p.showResolverMenu(gtx)
							}
							return p.resolver.Layout(gtx, th)
						}),

						layout.Rigid(func(gtx page.C) page.D {
							if p.hostMapper.Clicked(gtx) {
								p.showHostMapperMenu(gtx)
							}
							return p.hostMapper.Layout(gtx, th)
						}),

						layout.Rigid(func(gtx page.C) page.D {
							if p.limiter.Clicked(gtx) {
								p.showLimiterMenu(gtx)
							}
							return p.limiter.Layout(gtx, th)
						}),

						layout.Rigid(func(gtx page.C) page.D {
							if p.observer.Clicked(gtx) {
								p.showObserverMenu(gtx)
							}
							return p.observer.Layout(gtx, th)
						}),

						layout.Rigid(func(gtx page.C) page.D {
							if p.recordAdd.Clicked(gtx) {
								p.router.Goto(page.Route{
									Path:     page.PageServiceRecord,
									Callback: p.recordCallback,
									Perm:     page.PermReadWrite,
								})
							}

							return layout.Flex{
								Alignment: layout.Middle,
							}.Layout(gtx,
								layout.Flexed(1, func(gtx page.C) page.D {
									gtx.Source = src

									if p.recorderSelector.Clicked(gtx) {
										p.recordFolded = !p.recordFolded
									}
									return p.recorderSelector.Layout(gtx, th)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if !p.edit {
										return page.D{}
									}
									return layout.Spacer{Width: 8}.Layout(gtx)
								}),
								layout.Rigid(func(gtx page.C) page.D {
									if !p.edit {
										return page.D{}
									}
									btn := material.IconButton(th, &p.recordAdd, icons.IconAdd, "Add")
									btn.Background = theme.Current().ContentSurfaceBg
									btn.Color = th.Fg
									// btn.Inset = layout.UniformInset(8)
									return btn.Layout(gtx)
								}),
							)
						}),
						layout.Rigid(func(gtx page.C) page.D {
							if p.recordFolded {
								return page.D{}
							}

							gtx.Source = src
							return p.layoutRecorder(gtx, th)
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
									if !p.edit {
										return page.D{}
									}
									return layout.Spacer{Width: 8}.Layout(gtx)
								}),
								layout.Rigid(func(gtx page.C) page.D {
									if !p.edit {
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
						layout.Rigid(layout.Spacer{Height: 8}.Layout),
					)

				}),

				layout.Rigid(func(gtx page.C) page.D {
					return layout.Inset{
						Top:    16,
						Bottom: 16,
					}.Layout(gtx, material.H6(th, i18n.Handler.Value()).Layout)
				}),

				layout.Rigid(func(gtx page.C) page.D {
					gtx.Source = src
					return p.handler.Layout(gtx, th)
				}),

				layout.Rigid(func(gtx page.C) page.D {
					return layout.Inset{
						Top:    16,
						Bottom: 16,
					}.Layout(gtx, material.H6(th, i18n.Listener.Value()).Layout)
				}),

				layout.Rigid(func(gtx page.C) page.D {
					gtx.Source = src
					return p.listener.Layout(gtx, th)
				}),

				layout.Rigid(func(gtx page.C) page.D {
					if !p.handler.canForward() {
						return page.D{}
					}

					return layout.Inset{
						Top:    16,
						Bottom: 16,
					}.Layout(gtx, material.H6(th, i18n.Forwarder.Value()).Layout)
				}),

				layout.Rigid(func(gtx page.C) page.D {
					if !p.handler.canForward() {
						return page.D{}
					}
					gtx.Source = src
					return p.forwarder.Layout(gtx, th)
				}),
			)
		})
	})
}

func (p *servicePage) layoutRecorder(gtx page.C, th *page.T) page.D {
	for i := range p.records {
		if p.records[i].clk.Clicked(gtx) {
			perm := page.PermRead
			if p.edit {
				perm = page.PermReadWrite
			}
			p.router.Goto(page.Route{
				Path:     page.PageServiceRecord,
				ID:       p.records[i].id,
				Value:    p.records[i].cfg,
				Callback: p.recordCallback,
				Perm:     perm,
			})
			break
		}

		if p.records[i].delete.Clicked(gtx) {
			p.delRecordDialog.OnClick = func(ok bool) {
				p.router.HideModal(gtx)
				if !ok {
					return
				}
				p.records = append(p.records[:i], p.records[i+1:]...)

				p.recorderSelector.Clear()
				p.recorderSelector.Select(ui_widget.SelectorItem{Value: strconv.Itoa(len(p.records))})
			}
			p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
				return p.delRecordDialog.Layout(gtx, th)
			})
			break
		}
	}

	var children []layout.FlexChild
	for i := range p.records {
		ro := &p.records[i]

		children = append(children,
			layout.Rigid(func(gtx page.C) page.D {
				return layout.Flex{
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Flexed(1, func(gtx page.C) page.D {
						return material.Clickable(gtx, &ro.clk, func(gtx page.C) page.D {
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
										label := material.Body2(th, ro.cfg.Name)
										label.Font.Weight = font.SemiBold
										return label.Layout(gtx)
									}),
									layout.Rigid(layout.Spacer{Height: 4}.Layout),
									layout.Rigid(func(gtx page.C) page.D {
										return material.Body2(th, ro.cfg.Record).Layout(gtx)
									}),
								)
							})
						})
					}),
					layout.Rigid(func(gtx page.C) page.D {
						if !p.edit {
							return page.D{}
						}
						return layout.Spacer{Width: 8}.Layout(gtx)
					}),
					layout.Rigid(func(gtx page.C) page.D {
						if !p.edit {
							return page.D{}
						}
						btn := material.IconButton(th, &ro.delete, icons.IconDelete, "delete")
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

func (p *servicePage) layoutMetadata(gtx page.C, th *page.T) page.D {
	for i := range p.metadata {
		if p.metadata[i].clk.Clicked(gtx) {
			if p.edit {
				p.showMetadataDialog(gtx, i)
			}
			break
		}

		if p.metadata[i].delete.Clicked(gtx) {
			p.delMetadataDialog.OnClick = func(ok bool) {
				p.router.HideModal(gtx)
				if !ok {
					return
				}
				p.metadata = append(p.metadata[:i], p.metadata[i+1:]...)

				p.mdSelector.Clear()
				p.mdSelector.Select(ui_widget.SelectorItem{Value: strconv.Itoa(len(p.metadata))})
			}
			p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
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
						if !p.edit {
							return page.D{}
						}
						return layout.Spacer{Width: 8}.Layout(gtx)
					}),
					layout.Rigid(func(gtx page.C) page.D {
						if !p.edit {
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

func (p *servicePage) showMetadataDialog(gtx page.C, i int) {
	p.mdDialog.K.Clear()
	p.mdDialog.V.Clear()

	if i >= 0 && i < len(p.metadata) {
		p.mdDialog.K.SetText(p.metadata[i].k)
		p.mdDialog.V.SetText(p.metadata[i].v)
	}

	p.mdDialog.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
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

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.mdDialog.Layout(gtx, th)
	})
}

func (p *servicePage) showAdmissionMenu(gtx page.C) {
	options := []ui_widget.MenuOption{}
	for _, v := range api.GetConfig().Admissions {
		options = append(options, ui_widget.MenuOption{
			Value: v.Name,
		})
	}
	for i := range options {
		options[i].Selected = p.admission.AnyValue(options[i].Value)
	}

	p.menu.Title = i18n.Admission
	p.menu.Options = options
	p.menu.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		p.admission.Clear()
		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				p.admission.Select(ui_widget.SelectorItem{Value: p.menu.Options[i].Value})
			}
		}
	}
	p.menu.OnAdd = func() {
		p.router.Goto(page.Route{
			Path: page.PageAdmission,
			Perm: page.PermReadWrite,
		})
		p.router.HideModal(gtx)
	}
	p.menu.Multiple = true

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *servicePage) showBypassMenu(gtx page.C) {
	options := []ui_widget.MenuOption{}
	for _, v := range api.GetConfig().Bypasses {
		options = append(options, ui_widget.MenuOption{
			Value: v.Name,
		})
	}
	for i := range options {
		options[i].Selected = p.bypass.AnyValue(options[i].Value)
	}

	p.menu.Title = i18n.Bypass
	p.menu.Options = options
	p.menu.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		p.bypass.Clear()
		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				p.bypass.Select(ui_widget.SelectorItem{Value: p.menu.Options[i].Value})
			}
		}
	}
	p.menu.OnAdd = func() {
		p.router.Goto(page.Route{
			Path: page.PageBypass,
			Perm: page.PermReadWrite,
		})
		p.router.HideModal(gtx)
	}
	p.menu.Multiple = true

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *servicePage) showResolverMenu(gtx page.C) {
	options := []ui_widget.MenuOption{}
	for _, v := range api.GetConfig().Resolvers {
		options = append(options, ui_widget.MenuOption{
			Value: v.Name,
		})
	}
	for i := range options {
		options[i].Selected = p.resolver.AnyValue(options[i].Value)
	}

	p.menu.Title = i18n.Resolver
	p.menu.Options = options
	p.menu.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		p.resolver.Clear()
		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				p.resolver.Select(ui_widget.SelectorItem{Value: p.menu.Options[i].Value})
			}
		}
	}
	p.menu.OnAdd = func() {
		p.router.Goto(page.Route{
			Path: page.PageResolver,
			Perm: page.PermReadWrite,
		})
		p.router.HideModal(gtx)
	}
	p.menu.Multiple = false

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *servicePage) showHostMapperMenu(gtx page.C) {
	options := []ui_widget.MenuOption{}
	for _, v := range api.GetConfig().Hosts {
		options = append(options, ui_widget.MenuOption{
			Value: v.Name,
		})
	}
	for i := range options {
		options[i].Selected = p.hostMapper.AnyValue(options[i].Value)
	}

	p.menu.Title = i18n.Hosts
	p.menu.Options = options
	p.menu.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		p.hostMapper.Clear()
		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				p.hostMapper.Select(ui_widget.SelectorItem{Value: p.menu.Options[i].Value})
			}
		}
	}
	p.menu.OnAdd = func() {
		p.router.Goto(page.Route{
			Path: page.PageHosts,
			Perm: page.PermReadWrite,
		})
		p.router.HideModal(gtx)
	}
	p.menu.Multiple = false

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *servicePage) showLimiterMenu(gtx page.C) {
	options := []ui_widget.MenuOption{}
	for _, v := range api.GetConfig().Limiters {
		options = append(options, ui_widget.MenuOption{
			Value: v.Name,
		})
	}
	for i := range options {
		options[i].Selected = p.limiter.AnyValue(options[i].Value)
	}

	p.menu.Title = i18n.Limiter
	p.menu.Options = options
	p.menu.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		p.limiter.Clear()
		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				p.limiter.Select(ui_widget.SelectorItem{Value: p.menu.Options[i].Value})
			}
		}
	}
	p.menu.OnAdd = func() {
		p.router.Goto(page.Route{
			Path: page.PageLimiter,
			Perm: page.PermReadWrite,
		})
		p.router.HideModal(gtx)
	}
	p.menu.Multiple = false

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *servicePage) showObserverMenu(gtx page.C) {
	options := []ui_widget.MenuOption{}
	for _, v := range api.GetConfig().Observers {
		options = append(options, ui_widget.MenuOption{
			Value: v.Name,
		})
	}
	for i := range options {
		options[i].Selected = p.observer.AnyValue(options[i].Value)
	}

	p.menu.Title = i18n.Observer
	p.menu.Options = options
	p.menu.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		p.observer.Clear()
		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				p.observer.Select(ui_widget.SelectorItem{Value: p.menu.Options[i].Value})
			}
		}
	}
	p.menu.OnAdd = func() {
		p.router.Goto(page.Route{
			Path: page.PageObserver,
			Perm: page.PermReadWrite,
		})
		p.router.HideModal(gtx)
	}
	p.menu.Multiple = false

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

/*
func (p *servicePage) showLoggerMenu(gtx page.C) {
	options := []ui_widget.MenuOption{}
	for _, v := range api.GetConfig().Loggers {
		options = append(options, ui_widget.MenuOption{
			Value: v.Name,
		})
	}
	for i := range options {
		options[i].Selected = p.logger.AnyValue(options[i].Value)
	}

	p.menu.Title = i18n.Logger
	p.menu.Options = options
	p.menu.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		p.logger.Clear()
		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				p.logger.Select(ui_widget.SelectorItem{Value: p.menu.Options[i].Value})
			}
		}
	}
	p.menu.OnAdd = func() {}
	p.menu.Multiple = true

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}
*/

func (p *servicePage) recordCallback(action page.Action, id string, value any) {
	if id == "" {
		return
	}

	switch action {
	case page.ActionCreate:
		cfg, _ := value.(*api.RecorderObject)
		if cfg == nil {
			return
		}
		p.records = append(p.records, record{
			id:  id,
			cfg: cfg,
		})
	case page.ActionUpdate:
		cfg, _ := value.(*api.RecorderObject)
		if cfg == nil {
			return
		}
		for i := range p.records {
			if p.records[i].id == id {
				p.records[i].cfg = cfg
				break
			}
		}

	case page.ActionDelete:
		for i := range p.records {
			if p.records[i].id == id {
				p.records = append(p.records[:i], p.records[i+1:]...)
				break
			}
		}
	}

	p.recorderSelector.Clear()
	p.recorderSelector.Select(ui_widget.SelectorItem{Value: strconv.Itoa(len(p.records))})
}

func (p *servicePage) save() bool {
	cfg := p.generateConfig()

	var err error
	if p.id == "" {
		err = runner.Exec(context.Background(),
			task.CreateService(cfg),
			runner.WithCancel(true),
		)
	} else {
		err = runner.Exec(context.Background(),
			task.UpdateService(cfg),
			runner.WithCancel(true),
		)
	}
	util.RestartGetConfigTask()

	return err == nil
}

func (p *servicePage) generateConfig() *api.ServiceConfig {
	var svcCfg *api.ServiceConfig

	if p.id != "" {
		for _, svc := range api.GetConfig().Services {
			if svc == nil {
				continue
			}
			if svc.Name == p.id {
				svcCfg = svc.Copy()
				break
			}
		}
	}

	if svcCfg == nil {
		svcCfg = &api.ServiceConfig{}
	}

	svcCfg.Name = p.name.Text()
	svcCfg.Addr = p.addr.Text()
	// svcCfg.Interface = p.customInterface.Value()

	svcCfg.Admission = ""
	svcCfg.Admissions = nil
	if admissions := p.admission.Values(); len(admissions) > 1 {
		svcCfg.Admissions = admissions
	} else {
		if len(admissions) > 0 {
			svcCfg.Admission = admissions[0]
		}
	}

	svcCfg.Bypass = ""
	svcCfg.Bypasses = nil
	if bypasses := p.bypass.Values(); len(bypasses) > 1 {
		svcCfg.Bypasses = bypasses
	} else {
		if len(bypasses) > 0 {
			svcCfg.Bypass = bypasses[0]
		}
	}

	svcCfg.Resolver = p.resolver.Value()
	svcCfg.Hosts = p.hostMapper.Value()
	svcCfg.Limiter = p.limiter.Value()
	svcCfg.Observer = p.observer.Value()

	/*
		svcCfg.Logger = ""
		svcCfg.Loggers = nil
		if loggers := p.logger.Values(); len(loggers) > 1 {
			svcCfg.Loggers = loggers
		} else {
			if len(loggers) > 0 {
				svcCfg.Logger = loggers[0]
			}
		}
	*/

	svcCfg.Recorders = nil
	for i := range p.records {
		svcCfg.Recorders = append(svcCfg.Recorders, p.records[i].cfg)
	}

	svcCfg.Metadata = make(map[string]any)
	svcCfg.Metadata["enablestats"] = true
	for i := range p.metadata {
		svcCfg.Metadata[p.metadata[i].k] = p.metadata[i].v
	}

	if svcCfg.Handler == nil {
		svcCfg.Handler = &api.HandlerConfig{}
	}

	svcCfg.Handler.Type = p.handler.typ.Value()
	svcCfg.Handler.Chain = p.handler.chain.Value()

	svcCfg.Handler.Auther = ""
	svcCfg.Handler.Authers = nil
	svcCfg.Handler.Auth = nil
	if p.handler.authType.Value == string(page.AuthAuther) {
		if authers := p.handler.auther.Values(); len(authers) > 1 {
			svcCfg.Handler.Authers = authers
		} else {
			if len(authers) > 0 {
				svcCfg.Handler.Auther = authers[0]
			}
		}
	}
	if p.handler.authType.Value == string(page.AuthSimple) {
		username := strings.TrimSpace(p.handler.username.Text())
		password := strings.TrimSpace(p.handler.password.Text())
		if username != "" {
			svcCfg.Handler.Auth = &api.AuthConfig{
				Username: username,
				Password: password,
			}
		}
	}

	svcCfg.Handler.Limiter = p.handler.limiter.Value()
	svcCfg.Handler.Observer = p.handler.observer.Value()

	svcCfg.Handler.Metadata = make(map[string]any)
	for i := range p.handler.metadata {
		svcCfg.Handler.Metadata[p.handler.metadata[i].k] = p.handler.metadata[i].v
	}

	if svcCfg.Listener == nil {
		svcCfg.Listener = &api.ListenerConfig{}
	}

	svcCfg.Listener.Type = p.listener.typ.Value()
	svcCfg.Listener.Chain = p.listener.chain.Value()

	svcCfg.Listener.Auther = ""
	svcCfg.Listener.Authers = nil
	svcCfg.Listener.Auth = nil
	if p.listener.authType.Value == string(page.AuthAuther) {
		if authers := p.listener.auther.Values(); len(authers) > 1 {
			svcCfg.Listener.Authers = authers
		} else {
			if len(authers) > 0 {
				svcCfg.Listener.Auther = authers[0]
			}
		}
	}
	if p.listener.authType.Value == string(page.AuthSimple) {
		username := strings.TrimSpace(p.listener.username.Text())
		password := strings.TrimSpace(p.listener.password.Text())
		if username != "" {
			svcCfg.Listener.Auth = &api.AuthConfig{
				Username: username,
				Password: password,
			}
		}
	}

	svcCfg.Listener.TLS = nil
	if p.listener.enableTLS.Value() {
		svcCfg.Listener.TLS = &api.TLSConfig{
			CertFile: strings.TrimSpace(p.listener.tlsCertFile.Text()),
			KeyFile:  strings.TrimSpace(p.listener.tlsKeyFile.Text()),
			CAFile:   strings.TrimSpace(p.listener.tlsCAFile.Text()),
		}
	}

	svcCfg.Listener.Metadata = make(map[string]any)
	for i := range p.listener.metadata {
		svcCfg.Listener.Metadata[p.listener.metadata[i].k] = p.listener.metadata[i].v
	}

	svcCfg.Forwarder = nil
	if len(p.forwarder.nodes) > 0 {
		svcCfg.Forwarder = &api.ForwarderConfig{}

		for _, node := range p.forwarder.nodes {
			svcCfg.Forwarder.Nodes = append(svcCfg.Forwarder.Nodes, node.cfg)
		}
	}

	return svcCfg
}

func (p *servicePage) delete() {
	runner.Exec(context.Background(),
		task.DeleteService(p.id),
		runner.WithCancel(true),
	)
	util.RestartGetConfigTask()
}
