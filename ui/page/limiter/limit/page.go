package limit

import (
	"strings"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/go-gost/gostctl/ui/i18n"
	"github.com/go-gost/gostctl/ui/icons"
	"github.com/go-gost/gostctl/ui/page"
	ui_widget "github.com/go-gost/gostctl/ui/widget"
)

type limit struct {
	rule   string
	clk    widget.Clickable
	delete widget.Clickable
}

type limitPage struct {
	router *page.Router
	list   layout.List

	btnBack   widget.Clickable
	btnDelete widget.Clickable
	btnEdit   widget.Clickable
	btnSave   widget.Clickable
	btnAdd    widget.Clickable

	limits []limit

	id       string
	perm     page.Perm
	callback page.Callback

	edit bool

	ruleDialog ruleDialog
	delDialog  ui_widget.Dialog
}

func NewPage(r *page.Router) page.Page {
	p := &limitPage{
		router: r,

		list: layout.List{
			Axis: layout.Vertical,
		},
		ruleDialog: ruleDialog{
			rule: rule{
				value: component.TextField{
					Editor: widget.Editor{
						MaxLen:     255,
						SingleLine: true,
					},
				},
			},
		},
		delDialog: ui_widget.Dialog{Title: i18n.DeleteLimits},
	}

	return p
}

func (p *limitPage) Init(opts ...page.PageOption) {
	var options page.PageOptions
	for _, opt := range opts {
		opt(&options)
	}

	p.id = options.ID
	p.callback = options.Callback
	p.perm = options.Perm
	p.edit = p.perm&page.PermWrite > 0

	p.limits = nil
	limits, _ := options.Value.([]string)
	for i := range limits {
		if limits[i] == "" {
			continue
		}
		p.limits = append(p.limits, limit{
			rule: limits[i],
		})
	}
}

func (p *limitPage) Layout(gtx page.C) page.D {
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

func (p *limitPage) layout(gtx page.C, th *page.T) page.D {
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
						title := material.H6(th, i18n.Limits.Value())
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

		/*
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				div := component.Divider(th)
				div.Inset = layout.Inset{}
				return div.Layout(gtx)
			}),
		*/

		layout.Flexed(1, func(gtx page.C) page.D {
			for i := range p.limits {
				if p.limits[i].delete.Clicked(gtx) {
					p.limits = append(p.limits[:i], p.limits[i+1:]...)
					break
				}

				if p.limits[i].clk.Clicked(gtx) {
					if !p.edit {
						break
					}

					p.showDialog(gtx, &p.limits[i])
					break
				}
			}

			return p.list.Layout(gtx, len(p.limits), func(gtx page.C, index int) page.D {
				limit := &p.limits[index]

				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Alignment: layout.Middle,
						}.Layout(gtx,
							layout.Flexed(1, func(gtx page.C) page.D {
								return material.Clickable(gtx, &limit.clk, func(gtx page.C) page.D {
									return layout.UniformInset(16).Layout(gtx, func(gtx page.C) page.D {
										return layout.Flex{
											Axis: layout.Vertical,
										}.Layout(gtx,
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												label := material.Body1(th, limit.rule)
												return label.Layout(gtx)
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
									btn := material.IconButton(th, &limit.delete, icons.IconDelete, "Remove")
									btn.Color = th.Fg
									btn.Background = th.Bg
									return btn.Layout(gtx)
								})
							}),
						)
					}),
					/*
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							div := component.Divider(th)
							div.Inset = layout.Inset{}
							return div.Layout(gtx)
						}),
					*/
				)
			})
		}),
	)
}

func (p *limitPage) showDialog(gtx page.C, m *limit) {
	if m != nil {
		p.ruleDialog.rule.Set(m.rule)
	} else {
		p.ruleDialog.rule.Set("")
	}

	p.ruleDialog.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		v := p.ruleDialog.rule.Get()
		if m != nil {
			m.rule = v
			return
		}

		p.limits = append(p.limits, limit{
			rule: v,
		})
	}

	p.router.ShowModal(gtx, p.ruleDialog.Layout)
}

func (p *limitPage) generateConfig() []string {
	rules := []string{}
	for i := range p.limits {
		rule := strings.TrimSpace(p.limits[i].rule)
		if rule == "" {
			continue
		}
		rules = append(rules, rule)
	}
	return rules
}

func (p *limitPage) save() bool {
	if p.callback != nil {
		p.callback(page.ActionUpdate, p.id, p.generateConfig())
	}

	return true
}

func (p *limitPage) delete() {
	if p.callback != nil {
		p.callback(page.ActionDelete, p.id, nil)
	}
}
