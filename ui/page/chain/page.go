package chain

import (
	"context"
	"strconv"
	"strings"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/go-gost/gostctl/api"
	"github.com/go-gost/gostctl/api/runner"
	"github.com/go-gost/gostctl/api/runner/task"
	"github.com/go-gost/gostctl/api/util"
	"github.com/go-gost/gostctl/ui/i18n"
	"github.com/go-gost/gostctl/ui/icons"
	"github.com/go-gost/gostctl/ui/page"
	"github.com/go-gost/gostctl/ui/theme"
	ui_widget "github.com/go-gost/gostctl/ui/widget"
)

type metadata struct {
	k string
	v string
}

type chainHop struct {
	cfg    *api.HopConfig
	clk    widget.Clickable
	delete widget.Clickable
}

type chainPage struct {
	router *page.Router

	menu ui_widget.Menu
	list layout.List

	btnBack   widget.Clickable
	btnDelete widget.Clickable
	btnEdit   widget.Clickable
	btnSave   widget.Clickable

	name component.TextField

	hops   []chainHop
	addHop widget.Clickable

	id   string
	perm page.Perm

	edit   bool
	create bool

	metadata         []metadata
	metadataSelector ui_widget.Selector
	// metadataDialog   ui_widget.MetadataDialog
	delDialog ui_widget.Dialog
}

func NewPage(r *page.Router) page.Page {
	return &chainPage{
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
		delDialog:        ui_widget.Dialog{Title: i18n.DeleteChain},
		metadataSelector: ui_widget.Selector{Title: i18n.Metadata},
		// metadataDialog:   ui_widget.MetadataDialog{},
	}
}

func (p *chainPage) Init(opts ...page.PageOption) {
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
	chain := &api.ChainConfig{}
	for _, ch := range cfg.Chains {
		if ch.Name == p.id {
			chain = ch
			break
		}
	}

	p.name.SetText(chain.Name)

	{
		p.hops = nil
		for _, hop := range chain.Hops {
			p.hops = append(p.hops, chainHop{
				cfg: hop,
			})
		}
	}

	{
		p.metadata = nil
		meta := api.NewMetadata(chain.Metadata)
		for k := range chain.Metadata {
			md := metadata{
				k: k,
				v: meta.GetString(k),
			}
			p.metadata = append(p.metadata, md)
		}
		p.metadataSelector.Clear()
		p.metadataSelector.Select(ui_widget.SelectorItem{Value: strconv.Itoa(len(p.metadata))})
	}
}

func (p *chainPage) Layout(gtx page.C) page.D {
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
		p.delDialog.Title = i18n.DeleteChain
		p.delDialog.OnClick = func(ok bool) {
			if ok {
				p.delete()
				p.router.Back()
			}
			p.router.HideModal(gtx)
		}

		p.router.ShowModal(gtx, func(gtx page.C, th *material.Theme) page.D {
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
						title := material.H6(th, i18n.Chain.Value())
						return title.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: 4}.Layout),
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
			return p.list.Layout(gtx, 1, func(gtx page.C, _ int) page.D {
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

func (p *chainPage) layout(gtx page.C, th *page.T) page.D {
	src := gtx.Source

	if !p.edit {
		gtx = gtx.Disabled()
	}

	if p.addHop.Clicked(gtx) {
		p.showHopMenu(gtx)
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
					return material.Body1(th, i18n.Name.Value()).Layout(gtx)
				}),
				layout.Rigid(func(gtx page.C) page.D {
					return p.name.Layout(gtx, th, "")
				}),
				layout.Rigid(layout.Spacer{Height: 8}.Layout),
				layout.Rigid(func(gtx page.C) page.D {
					return layout.Inset{
						Top:    8,
						Bottom: 8,
					}.Layout(gtx, func(gtx page.C) page.D {
						return layout.Flex{
							Alignment: layout.Middle,
						}.Layout(gtx,
							layout.Flexed(1, func(gtx page.C) page.D {
								return layout.Inset{
									Top:    10,
									Bottom: 10,
								}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return material.Body1(th, i18n.Hop.Value()).Layout(gtx)
								})
							}),
							layout.Rigid(func(gtx page.C) page.D {
								if !p.edit {
									return page.D{}
								}
								btn := material.IconButton(th, &p.addHop, icons.IconAdd, "Add")
								btn.Background = theme.Current().ContentSurfaceBg
								btn.Color = th.Fg
								// btn.Inset = layout.UniformInset(8)
								return btn.Layout(gtx)
							}),
						)
					})
				}),
				layout.Rigid(func(gtx page.C) page.D {
					gtx.Source = src
					return p.layoutHops(gtx, th)
				}),
			)
		})
	})
}

func (p *chainPage) layoutHops(gtx page.C, th *page.T) page.D {
	for i := range p.hops {
		if p.hops[i].clk.Clicked(gtx) {
			route := page.Route{
				Path: page.PageHop,
				ID:   p.hops[i].cfg.Name,
				Perm: page.PermRead,
			}
			if len(p.hops[i].cfg.Nodes) > 0 {
				route.Value = p.hops[i].cfg
			}
			p.router.Goto(route)
			break
		}

		if p.hops[i].delete.Clicked(gtx) {
			p.delDialog.Title = i18n.DeleteHop
			p.delDialog.OnClick = func(ok bool) {
				p.router.HideModal(gtx)
				if !ok {
					return
				}
				p.hops = append(p.hops[:i], p.hops[i+1:]...)
			}
			p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
				return p.delDialog.Layout(gtx, th)
			})
			break
		}
	}

	var children []layout.FlexChild
	for i := range p.hops {
		hop := &p.hops[i]

		children = append(children, layout.Rigid(func(gtx page.C) page.D {
			return layout.Flex{
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Flexed(1, func(gtx page.C) page.D {
					return material.Clickable(gtx, &hop.clk, func(gtx page.C) page.D {
						return layout.Inset{
							Top:    16,
							Bottom: 16,
							Left:   8,
							Right:  8,
						}.Layout(gtx, func(gtx page.C) page.D {
							return material.Body2(th, hop.cfg.Name).Layout(gtx)
						})
					})
				}),
				layout.Rigid(func(gtx page.C) page.D {
					if !p.edit {
						return page.D{}
					}
					btn := material.IconButton(th, &hop.delete, icons.IconDelete, "delete")
					btn.Background = theme.Current().ContentSurfaceBg
					btn.Color = th.Fg
					// btn.Inset = layout.UniformInset(8)
					return btn.Layout(gtx)
				}),
			)
		}))
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx, children...)
}

func (p *chainPage) showHopMenu(gtx page.C) {
	options := []ui_widget.MenuOption{}
	for _, v := range api.GetConfig().Hops {
		options = append(options, ui_widget.MenuOption{
			Value: v.Name,
		})
	}

	p.menu.Title = i18n.Hop
	p.menu.Options = options
	p.menu.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				p.hops = append(p.hops, chainHop{cfg: &api.HopConfig{Name: p.menu.Options[i].Value}})
			}
		}
	}

	p.menu.OnAdd = func() {
		p.router.Goto(page.Route{
			Path: page.PageHop,
			Perm: page.PermReadWrite,
		})
		p.router.HideModal(gtx)
	}
	p.menu.Multiple = false

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *chainPage) save() bool {
	cfg := p.generateConfig()

	var err error
	if p.id == "" {
		err = runner.Exec(context.Background(),
			task.CreateChain(cfg),
			runner.WithCancel(true),
		)
	} else {
		err = runner.Exec(context.Background(),
			task.UpdateChain(cfg),
			runner.WithCancel(true),
		)
	}
	util.RestartGetConfigTask()

	return err == nil
}

func (p *chainPage) generateConfig() *api.ChainConfig {
	chainCfg := &api.ChainConfig{
		Name: strings.TrimSpace(p.name.Text()),
	}

	for i := range p.hops {
		chainCfg.Hops = append(chainCfg.Hops, p.hops[i].cfg)
	}
	return chainCfg
}

func (p *chainPage) delete() {
	runner.Exec(context.Background(),
		task.DeleteChain(p.id),
		runner.WithCancel(true),
	)
	util.RestartGetConfigTask()
}
