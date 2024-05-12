package metadata

import (
	"sort"
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

type metadata struct {
	md     page.Metadata
	clk    widget.Clickable
	delete widget.Clickable
}

type mdPage struct {
	router *page.Router
	list   layout.List

	btnBack   widget.Clickable
	btnDelete widget.Clickable
	btnEdit   widget.Clickable
	btnSave   widget.Clickable
	btnAdd    widget.Clickable

	metadata []metadata

	id       string
	perm     page.Perm
	callback page.Callback

	edit bool

	mdDialog  metadataDialog
	delDialog ui_widget.Dialog
}

func NewPage(r *page.Router) page.Page {
	p := &mdPage{
		router: r,

		list: layout.List{
			Axis: layout.Vertical,
		},
		mdDialog: metadataDialog{
			kv: kv{
				k: component.TextField{
					Editor: widget.Editor{
						MaxLen:     128,
						SingleLine: true,
					},
				},
				v: component.TextField{
					Editor: widget.Editor{
						MaxLen:     128,
						SingleLine: true,
					},
				},
			},
		},
		delDialog: ui_widget.Dialog{Title: i18n.DeleteMetadata},
	}

	return p
}

func (p *mdPage) Init(opts ...page.PageOption) {
	var options page.PageOptions
	for _, opt := range opts {
		opt(&options)
	}

	p.id = options.ID
	p.callback = options.Callback
	p.perm = options.Perm
	p.edit = p.perm&page.PermWrite > 0

	p.metadata = nil
	md, _ := options.Value.(api.Metadata)
	for k := range md {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		p.metadata = append(p.metadata, metadata{
			md: page.Metadata{
				K: k,
				V: md.GetString(k),
			},
		})
	}
	sort.SliceStable(p.metadata, func(i, j int) bool {
		return strings.Compare(p.metadata[i].md.K, p.metadata[j].md.K) < 0
	})
}

func (p *mdPage) Layout(gtx page.C) page.D {
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

func (p *mdPage) layout(gtx page.C, th *page.T) page.D {
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
						title := material.H6(th, i18n.Metadata.Value())
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
			for i := range p.metadata {
				if p.metadata[i].delete.Clicked(gtx) {
					p.metadata = append(p.metadata[:i], p.metadata[i+1:]...)
					break
				}

				if p.metadata[i].clk.Clicked(gtx) {
					if !p.edit {
						break
					}

					p.showDialog(gtx, &p.metadata[i].md)
					break
				}
			}

			return p.list.Layout(gtx, len(p.metadata), func(gtx page.C, index int) page.D {
				md := &p.metadata[index]

				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Alignment: layout.Middle,
						}.Layout(gtx,
							layout.Flexed(1, func(gtx page.C) page.D {
								return material.Clickable(gtx, &md.clk, func(gtx page.C) page.D {
									return layout.UniformInset(16).Layout(gtx, func(gtx page.C) page.D {
										return layout.Flex{
											Axis: layout.Vertical,
										}.Layout(gtx,
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												label := material.Body1(th, md.md.K)
												label.Font.Weight = font.SemiBold
												return label.Layout(gtx)
											}),
											layout.Rigid(layout.Spacer{Height: 8}.Layout),
											layout.Rigid(material.Body2(th, md.md.V).Layout),
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
									btn := material.IconButton(th, &md.delete, icons.IconDelete, "Remove")
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

func (p *mdPage) showDialog(gtx page.C, md *page.Metadata) {
	if md != nil {
		p.mdDialog.kv.Set(md.K, md.V)
	} else {
		p.mdDialog.kv.Set("", "")
	}

	p.mdDialog.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		k, v := p.mdDialog.kv.Get()
		if md != nil {
			md.K = k
			md.V = v
			return
		}

		p.metadata = append(p.metadata, metadata{
			md: page.Metadata{
				K: k,
				V: v,
			},
		})
	}

	p.router.ShowModal(gtx, p.mdDialog.Layout)
}

func (p *mdPage) generateConfig() api.Metadata {
	md := api.Metadata{}
	for i := range p.metadata {
		k := strings.TrimSpace(p.metadata[i].md.K)
		if k == "" {
			continue
		}
		md[k] = strings.TrimSpace(p.metadata[i].md.V)
	}
	return md
}

func (p *mdPage) save() bool {
	if p.callback != nil {
		p.callback(page.ActionUpdate, p.id, p.generateConfig())
	}

	return true
}

func (p *mdPage) delete() {
	if p.callback != nil {
		p.callback(page.ActionDelete, p.id, nil)
	}
}
