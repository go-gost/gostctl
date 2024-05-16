package mapping

import (
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
	ui_widget "github.com/go-gost/gostctl/ui/widget"
)

type mapping struct {
	hostname string
	ip       string
	alias    []string
	clk      widget.Clickable
	delete   widget.Clickable
}

type mappingPage struct {
	router *page.Router
	list   layout.List

	btnBack   widget.Clickable
	btnDelete widget.Clickable
	btnEdit   widget.Clickable
	btnSave   widget.Clickable
	btnAdd    widget.Clickable

	mappings []mapping

	id       string
	perm     page.Perm
	callback page.Callback

	edit bool

	mappingDialog mappingDialog
	delDialog     ui_widget.Dialog
}

func NewPage(r *page.Router) page.Page {
	p := &mappingPage{
		router: r,

		list: layout.List{
			Axis: layout.Vertical,
		},
		mappingDialog: mappingDialog{
			hostname: component.TextField{
				Editor: widget.Editor{
					MaxLen:     255,
					SingleLine: true,
				},
			},
			ip: component.TextField{
				Editor: widget.Editor{
					MaxLen:     128,
					SingleLine: true,
				},
			},
			alias: component.TextField{
				Editor: widget.Editor{
					MaxLen:     255,
					SingleLine: true,
				},
			},
		},
		delDialog: ui_widget.Dialog{Title: i18n.DeleteAuth},
	}

	return p
}

func (p *mappingPage) Init(opts ...page.PageOption) {
	var options page.PageOptions
	for _, opt := range opts {
		opt(&options)
	}

	p.id = options.ID
	p.callback = options.Callback
	p.perm = options.Perm
	p.edit = p.perm&page.PermWrite > 0

	p.mappings = nil
	mappings, _ := options.Value.([]*api.HostMappingConfig)
	for i := range mappings {
		if mappings[i] == nil || mappings[i].Hostname == "" {
			continue
		}
		p.mappings = append(p.mappings, mapping{
			hostname: mappings[i].Hostname,
			ip:       mappings[i].IP,
			alias:    mappings[i].Aliases,
		})
	}
}

func (p *mappingPage) Layout(gtx page.C) page.D {
	if p.btnAdd.Clicked(gtx) {
		p.showDialog(gtx, nil)
	}

	th := p.router.Theme

	return layout.Stack{
		Alignment: layout.SE,
	}.Layout(gtx,
		layout.Expanded(func(gtx page.C) page.D {
			// gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return p.layout(gtx, th)
		}),
		layout.Stacked(func(gtx page.C) page.D {
			if !p.edit {
				return page.D{}
			}

			return layout.UniformInset(16).Layout(gtx, func(gtx page.C) page.D {
				btn := material.IconButton(th, &p.btnAdd, icons.IconAdd, "Add")
				btn.Inset = layout.UniformInset(16)

				return btn.Layout(gtx)
			})
		}),
	)
}

func (p *mappingPage) layout(gtx page.C, th *page.T) page.D {
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

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		// header
		layout.Rigid(func(gtx page.C) page.D {
			return layout.UniformInset(8).Layout(gtx, func(gtx page.C) page.D {
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
						title := material.H6(th, i18n.HostMappings.Value())
						return title.Layout(gtx)
					}),
					layout.Rigid(func(gtx page.C) page.D {
						if p.perm&page.PermDelete == 0 {
							return page.D{}
						}
						btn := material.IconButton(th, &p.btnDelete, icons.IconDelete, "Delete")
						btn.Color = th.Fg
						btn.Background = th.Bg
						return btn.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: 8}.Layout),
					layout.Rigid(func(gtx page.C) page.D {
						if !p.edit {
							return page.D{}
						}

						btn := material.IconButton(th, &p.btnSave, icons.IconDone, "Done")
						btn.Color = th.Fg
						btn.Background = th.Bg
						return btn.Layout(gtx)
					}),
				)
			})
		}),
		layout.Rigid(layout.Spacer{Height: 16}.Layout),

		layout.Flexed(1, func(gtx page.C) page.D {
			for i := range p.mappings {
				if p.mappings[i].delete.Clicked(gtx) {
					p.mappings = append(p.mappings[:i], p.mappings[i+1:]...)
					break
				}

				if p.mappings[i].clk.Clicked(gtx) {
					if !p.edit {
						break
					}

					p.showDialog(gtx, &p.mappings[i])
					break
				}
			}

			return p.list.Layout(gtx, len(p.mappings), func(gtx page.C, index int) page.D {
				mapping := &p.mappings[index]

				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Alignment: layout.Middle,
						}.Layout(gtx,
							layout.Flexed(1, func(gtx page.C) page.D {
								return material.Clickable(gtx, &mapping.clk, func(gtx page.C) page.D {
									return layout.UniformInset(16).Layout(gtx, func(gtx page.C) page.D {
										return layout.Flex{
											Axis: layout.Vertical,
										}.Layout(gtx,
											layout.Rigid(func(gtx page.C) page.D {
												label := material.Body1(th, mapping.hostname)
												label.Font.Weight = font.SemiBold
												return label.Layout(gtx)
											}),
											layout.Rigid(layout.Spacer{Height: 8}.Layout),
											layout.Rigid(material.Body2(th, mapping.ip).Layout),

											layout.Rigid(func(gtx page.C) page.D {
												if len(mapping.alias) == 0 {
													return page.D{}
												}
												return layout.Spacer{Height: 8}.Layout(gtx)
											}),
											layout.Rigid(func(gtx page.C) page.D {
												if len(mapping.alias) == 0 {
													return page.D{}
												}
												return material.Body2(th, strings.Join(mapping.alias, ", ")).Layout(gtx)
											}),
										)
									})
								})
							}),
							layout.Rigid(func(gtx page.C) page.D {
								if !p.edit {
									return page.D{}
								}

								return layout.Inset{
									Top:    8,
									Bottom: 8,
									Left:   8,
									Right:  16,
								}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									btn := material.IconButton(th, &mapping.delete, icons.IconDelete, "Remove")
									btn.Color = th.Fg
									btn.Background = th.Bg
									return btn.Layout(gtx)
								})
							}),
						)
					}),
				)
			})
		}),
	)
}

func (p *mappingPage) showDialog(gtx page.C, m *mapping) {
	p.mappingDialog.hostname.Clear()
	p.mappingDialog.ip.Clear()
	p.mappingDialog.alias.Clear()
	if m != nil {
		p.mappingDialog.hostname.SetText(m.hostname)
		p.mappingDialog.ip.SetText(m.ip)
		p.mappingDialog.alias.SetText(strings.Join(m.alias, ", "))
	}

	p.mappingDialog.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		hostname := strings.TrimSpace(p.mappingDialog.hostname.Text())
		ip := strings.TrimSpace(p.mappingDialog.ip.Text())
		var alias []string
		for _, s := range strings.Split(p.mappingDialog.alias.Text(), ",") {
			if s = strings.TrimSpace(s); s != "" {
				alias = append(alias, s)
			}
		}

		if m != nil {
			m.hostname = hostname
			m.ip = ip
			m.alias = alias
			return
		}

		p.mappings = append(p.mappings, mapping{
			hostname: hostname,
			ip:       ip,
			alias:    alias,
		})
	}

	p.router.ShowModal(gtx, p.mappingDialog.Layout)
}

func (p *mappingPage) generateConfig() []*api.HostMappingConfig {
	mappings := []*api.HostMappingConfig{}
	for i := range p.mappings {
		hostname := strings.TrimSpace(p.mappings[i].hostname)
		if hostname == "" {
			continue
		}
		mappings = append(mappings, &api.HostMappingConfig{
			Hostname: hostname,
			IP:       strings.TrimSpace(p.mappings[i].ip),
			Aliases:  p.mappings[i].alias,
		})
	}
	return mappings
}

func (p *mappingPage) save() bool {
	if p.callback != nil {
		p.callback(page.ActionUpdate, p.id, p.generateConfig())
	}

	return true
}

func (p *mappingPage) delete() {
	if p.callback != nil {
		p.callback(page.ActionDelete, p.id, nil)
	}
}
