package matcher

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

type matcher struct {
	rule   string
	clk    widget.Clickable
	delete widget.Clickable
}

type matcherPage struct {
	router *page.Router
	list   layout.List

	btnBack   widget.Clickable
	btnDelete widget.Clickable
	btnEdit   widget.Clickable
	btnSave   widget.Clickable
	btnAdd    widget.Clickable

	matchers []matcher

	id       string
	perm     page.Perm
	callback page.Callback

	edit bool

	ruleDialog ruleDialog
	delDialog  ui_widget.Dialog
}

func NewPage(r *page.Router) page.Page {
	p := &matcherPage{
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
		delDialog: ui_widget.Dialog{Title: i18n.DeleteRules},
	}

	return p
}

func (p *matcherPage) Init(opts ...page.PageOption) {
	var options page.PageOptions
	for _, opt := range opts {
		opt(&options)
	}

	p.id = options.ID
	p.callback = options.Callback
	p.perm = options.Perm
	p.edit = p.perm&page.PermWrite > 0

	p.matchers = nil
	matchers, _ := options.Value.([]string)
	for i := range matchers {
		if matchers[i] == "" {
			continue
		}
		p.matchers = append(p.matchers, matcher{
			rule: matchers[i],
		})
	}
}

func (p *matcherPage) Layout(gtx page.C) page.D {
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

func (p *matcherPage) layout(gtx page.C, th *page.T) page.D {
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
						title := material.H6(th, i18n.Rules.Value())
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
			for i := range p.matchers {
				if p.matchers[i].delete.Clicked(gtx) {
					p.matchers = append(p.matchers[:i], p.matchers[i+1:]...)
					break
				}

				if p.matchers[i].clk.Clicked(gtx) {
					if !p.edit {
						break
					}

					p.showDialog(gtx, &p.matchers[i])
					break
				}
			}

			return p.list.Layout(gtx, len(p.matchers), func(gtx page.C, index int) page.D {
				matcher := &p.matchers[index]

				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Alignment: layout.Middle,
						}.Layout(gtx,
							layout.Flexed(1, func(gtx page.C) page.D {
								return material.Clickable(gtx, &matcher.clk, func(gtx page.C) page.D {
									return layout.UniformInset(16).Layout(gtx, func(gtx page.C) page.D {
										return layout.Flex{
											Axis: layout.Vertical,
										}.Layout(gtx,
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												label := material.Body1(th, matcher.rule)
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
									btn := material.IconButton(th, &matcher.delete, icons.IconDelete, "Remove")
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

func (p *matcherPage) showDialog(gtx page.C, m *matcher) {
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

		p.matchers = append(p.matchers, matcher{
			rule: v,
		})
	}

	p.router.ShowModal(gtx, p.ruleDialog.Layout)
}

func (p *matcherPage) generateConfig() []string {
	rules := []string{}
	for i := range p.matchers {
		rule := strings.TrimSpace(p.matchers[i].rule)
		if rule == "" {
			continue
		}
		rules = append(rules, rule)
	}
	return rules
}

func (p *matcherPage) save() bool {
	if p.callback != nil {
		p.callback(page.ActionUpdate, p.id, p.generateConfig())
	}

	return true
}

func (p *matcherPage) delete() {
	if p.callback != nil {
		p.callback(page.ActionDelete, p.id, nil)
	}
}
