package service

import (
	"fmt"
	"strconv"
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

type node struct {
	id     string
	cfg    *api.ForwardNodeConfig
	clk    widget.Clickable
	delete widget.Clickable
}

type forwarder struct {
	service *servicePage
	menu    ui_widget.Menu

	enableSelector      ui_widget.Switcher
	selectorStrategy    ui_widget.Selector
	selectorMaxFails    component.TextField
	selectorFailTimeout component.TextField

	addNode   widget.Clickable
	nodes     []node
	delDialog ui_widget.Dialog
}

func newForwarder(service *servicePage) *forwarder {
	forwarder := &forwarder{
		service: service,

		enableSelector:   ui_widget.Switcher{Title: i18n.Selector},
		selectorStrategy: ui_widget.Selector{Title: i18n.SelectorStrategy},
		selectorMaxFails: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     8,
				Filter:     "1234567890",
			},
		},
		selectorFailTimeout: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     16,
				Filter:     "1234567890",
			},
			Suffix: func(gtx page.C) page.D {
				return material.Body1(service.router.Theme, i18n.TimeSecond.Value()).Layout(gtx)
			},
		},

		delDialog: ui_widget.Dialog{
			Title: i18n.DeleteNode,
		},
	}

	return forwarder
}

func (p *forwarder) init(cfg *api.ForwarderConfig) {
	if cfg == nil {
		cfg = &api.ForwarderConfig{}
	}

	{
		p.enableSelector.SetValue(false)
		p.selectorStrategy.Clear()
		p.selectorMaxFails.Clear()
		p.selectorFailTimeout.Clear()

		if selector := cfg.Selector; selector != nil {
			p.enableSelector.SetValue(true)
			for i := range page.SelectorStrategyOptions {
				if page.SelectorStrategyOptions[i].Value == selector.Strategy {
					p.selectorStrategy.Select(ui_widget.SelectorItem{Name: page.SelectorStrategyOptions[i].Name, Key: page.SelectorStrategyOptions[i].Key, Value: page.SelectorStrategyOptions[i].Value})
					break
				}
			}
			p.selectorMaxFails.SetText(strconv.Itoa(selector.MaxFails))
			p.selectorFailTimeout.SetText(selector.FailTimeout.String())
		}
	}

	p.nodes = nil
	for _, v := range cfg.Nodes {
		if v == nil {
			continue
		}

		nd := node{
			id:  uuid.New().String(),
			cfg: v,
		}

		p.nodes = append(p.nodes, nd)
	}
}

func (p *forwarder) Layout(gtx page.C, th *page.T) page.D {
	src := gtx.Source
	if !p.service.edit {
		gtx = gtx.Disabled()
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		// advanced mode
		layout.Rigid(func(gtx page.C) page.D {
			if p.service.mode.Value != string(page.AdvancedMode) {
				return page.D{}
			}

			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				// Selector
				layout.Rigid(func(gtx page.C) page.D {
					return p.enableSelector.Layout(gtx, th)
				}),

				layout.Rigid(func(gtx page.C) page.D {
					if !p.enableSelector.Value() {
						return page.D{}
					}

					return layout.UniformInset(8).Layout(gtx, func(gtx page.C) page.D {
						return layout.Flex{
							Axis: layout.Vertical,
						}.Layout(gtx,
							layout.Rigid(func(gtx page.C) page.D {
								if p.selectorStrategy.Clicked(gtx) {
									p.showSelectorStrategyMenu(gtx)
								}
								return p.selectorStrategy.Layout(gtx, th)
							}),
							layout.Rigid(func(gtx page.C) page.D {
								return p.selectorMaxFails.Layout(gtx, th, "Max fails")
							}),
							layout.Rigid(func(gtx page.C) page.D {
								return p.selectorFailTimeout.Layout(gtx, th, "Fail timeout in seconds")
							}),
						)
					})
				}),
			)
		}),
		layout.Rigid(func(gtx page.C) page.D {
			if p.addNode.Clicked(gtx) {
				p.service.router.Goto(page.Route{
					Path:     page.PageForwarderNode,
					Callback: p.nodeCallback,
					Perm:     page.PermReadWrite,
				})
			}

			return layout.Flex{
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Flexed(1, func(gtx page.C) page.D {
					return layout.Inset{
						Top:    16,
						Bottom: 16,
					}.Layout(gtx, material.Body1(th, i18n.Nodes.Value()).Layout)
				}),
				layout.Rigid(func(gtx page.C) page.D {
					if !p.service.edit {
						return page.D{}
					}
					btn := material.IconButton(th, &p.addNode, icons.IconAdd, "Add")
					btn.Background = theme.Current().ContentSurfaceBg
					btn.Color = th.Fg
					// btn.Inset = layout.UniformInset(8)
					return btn.Layout(gtx)
				}),
			)
		}),
		layout.Rigid(func(gtx page.C) page.D {
			gtx.Source = src
			return p.layoutNodes(gtx, th)
		}),
	)
}

func (p *forwarder) layoutNodes(gtx page.C, th *page.T) page.D {
	for i := range p.nodes {
		if p.nodes[i].clk.Clicked(gtx) {
			perm := page.PermRead
			if p.service.edit {
				perm = page.PermReadWrite
			}
			p.service.router.Goto(page.Route{
				Path:     page.PageForwarderNode,
				ID:       p.nodes[i].id,
				Value:    p.nodes[i].cfg,
				Callback: p.nodeCallback,
				Perm:     perm,
			})
			break
		}

		if p.nodes[i].delete.Clicked(gtx) {
			p.delDialog.OnClick = func(ok bool) {
				p.service.router.HideModal(gtx)
				if !ok {
					return
				}
				p.nodes = append(p.nodes[:i], p.nodes[i+1:]...)
			}
			p.service.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
				return p.delDialog.Layout(gtx, th)
			})
			break
		}
	}

	var children []layout.FlexChild
	for i := range p.nodes {
		node := &p.nodes[i]

		children = append(children,
			layout.Rigid(func(gtx page.C) page.D {
				name := strings.TrimSpace(node.cfg.Name)
				if name == "" {
					name = fmt.Sprintf("node-%d", i+1)
				}

				return layout.Flex{
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Flexed(1, func(gtx page.C) page.D {
						return material.Clickable(gtx, &node.clk, func(gtx page.C) page.D {
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
										return material.Body1(th, name).Layout(gtx)
									}),
									layout.Rigid(layout.Spacer{Height: 4}.Layout),
									layout.Rigid(func(gtx page.C) page.D {
										return material.Body2(th, node.cfg.Addr).Layout(gtx)
									}),
								)
							})
						})
					}),
					layout.Rigid(func(gtx page.C) page.D {
						if !p.service.edit {
							return page.D{}
						}
						btn := material.IconButton(th, &node.delete, icons.IconDelete, "delete")
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

func (p *forwarder) showSelectorStrategyMenu(gtx page.C) {
	for i := range page.SelectorStrategyOptions {
		page.SelectorStrategyOptions[i].Selected = p.selectorStrategy.AnyValue(page.SelectorStrategyOptions[i].Value)
	}

	p.menu.Title = i18n.SelectorStrategy
	p.menu.Options = page.SelectorStrategyOptions
	p.menu.OnClick = func(ok bool) {
		p.service.router.HideModal(gtx)
		if !ok {
			return
		}

		p.selectorStrategy.Clear()
		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				p.selectorStrategy.Select(ui_widget.SelectorItem{Name: p.menu.Options[i].Name, Key: p.menu.Options[i].Key, Value: p.menu.Options[i].Value})
			}
		}
	}
	p.menu.OnAdd = nil
	p.menu.Multiple = false

	p.service.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *forwarder) nodeCallback(action page.Action, id string, value any) {
	if id == "" {
		return
	}

	switch action {
	case page.ActionCreate:
		cfg, _ := value.(*api.ForwardNodeConfig)
		if cfg == nil {
			return
		}
		p.nodes = append(p.nodes, node{
			id:  id,
			cfg: cfg,
		})
	case page.ActionUpdate:
		cfg, _ := value.(*api.ForwardNodeConfig)
		if cfg == nil {
			return
		}
		for i := range p.nodes {
			if p.nodes[i].id == id {
				p.nodes[i].cfg = cfg
				break
			}
		}
	case page.ActionDelete:
		for i := range p.nodes {
			if p.nodes[i].id == id {
				p.nodes = append(p.nodes[:i], p.nodes[i+1:]...)
				break
			}
		}
	}
}
