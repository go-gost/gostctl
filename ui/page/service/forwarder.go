package service

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
	"github.com/go-gost/gostctl/ui/theme"
	ui_widget "github.com/go-gost/gostctl/ui/widget"
)

type node struct {
	fold    bool
	btnFold widget.Clickable
	delete  widget.Clickable

	name component.TextField
	addr component.TextField

	bypass ui_widget.Selector

	enableFilter ui_widget.Switcher
	protocol     ui_widget.Selector
	host         component.TextField
	path         component.TextField

	enableHTTP   ui_widget.Switcher
	httpHost     component.TextField
	httpUsername component.TextField
	httpPassword component.TextField

	enableTLS     ui_widget.Switcher
	tlsSecure     ui_widget.Switcher
	tlsServerName component.TextField
}

type forwarder struct {
	modal *component.ModalLayer
	menu  ui_widget.Menu
	mode  *widget.Enum

	addNode widget.Clickable
	nodes   []node
}

func (p *forwarder) Layout(gtx C, th *material.Theme) D {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(layout.Spacer{Height: 8}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if p.addNode.Clicked(gtx) {
				p.nodes = append(p.nodes, node{
					bypass:       ui_widget.Selector{Title: i18n.Bypass},
					enableFilter: ui_widget.Switcher{Title: i18n.Filter},
					protocol:     ui_widget.Selector{Title: i18n.Protocol},
					enableHTTP:   ui_widget.Switcher{Title: "HTTP"},
					enableTLS:    ui_widget.Switcher{Title: "TLS"},
					tlsSecure:    ui_widget.Switcher{Title: i18n.VerifyServerCert},
				})
			}

			return layout.Flex{
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Flexed(1, material.Body1(th, i18n.Nodes.Value()).Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btn := material.IconButton(th, &p.addNode, icons.IconAdd, "Add")
					btn.Background = theme.Current().ContentSurfaceBg
					btn.Color = th.Fg
					btn.Inset = layout.UniformInset(4)
					return btn.Layout(gtx)
				}),
			)
		}),
		layout.Rigid(layout.Spacer{Height: 8}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return p.layoutNodes(gtx, th)
		}),
	)
}

func (p *forwarder) layoutNodes(gtx C, th *material.Theme) D {
	for i := range p.nodes {
		if p.nodes[i].delete.Clicked(gtx) {
			p.nodes = append(p.nodes[:i], p.nodes[i+1:]...)
			break
		}
	}

	var children []layout.FlexChild
	for i := range p.nodes {
		node := &p.nodes[i]

		children = append(children,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				name := strings.TrimSpace(node.name.Text())
				if name == "" {
					name = fmt.Sprintf("Node-%d", i+1)
				}
				if node.btnFold.Clicked(gtx) {
					node.fold = !node.fold
				}

				return layout.Inset{
					Left:  8,
					Right: 8,
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(gtx,
						// node header
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return node.btnFold.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{
									Top:    8,
									Bottom: 8,
								}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return layout.Flex{
										Alignment: layout.Middle,
									}.Layout(gtx,
										layout.Flexed(1, material.Body1(th, name).Layout),
										layout.Rigid(layout.Spacer{Width: 8}.Layout),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											btn := material.IconButton(th, &node.delete, icons.IconDelete, "delete")
											btn.Background = theme.Current().ContentSurfaceBg
											btn.Color = th.Fg
											btn.Inset = layout.UniformInset(4)
											return btn.Layout(gtx)
										}),
										layout.Rigid(layout.Spacer{Width: 8}.Layout),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											if node.fold {
												return icons.IconNavRight.Layout(gtx, th.Fg)
											}
											return icons.IconNavExpandMore.Layout(gtx, th.Fg)
										}),
									)
								})
							})
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if node.fold {
								return D{}
							}

							return layout.Flex{
								Axis: layout.Vertical,
							}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return node.name.Layout(gtx, th, i18n.Name.Value())
								}),
								layout.Rigid(layout.Spacer{Height: 4}.Layout),

								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return node.addr.Layout(gtx, th, i18n.Address.Value())
								}),
								layout.Rigid(layout.Spacer{Height: 4}.Layout),

								// reverse proxy options
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if p.mode.Value == BasicMode {
										return D{}
									}

									return layout.Flex{
										Axis: layout.Vertical,
									}.Layout(gtx,
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											if node.bypass.Clicked(gtx) {
												p.showBypassMenu(gtx, node)
											}
											return node.bypass.Layout(gtx, th)
										}),

										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											return node.enableFilter.Layout(gtx, th)
										}),

										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											if !node.enableFilter.Value() {
												return layout.Dimensions{}
											}

											return layout.Flex{
												Axis: layout.Vertical,
											}.Layout(gtx,
												layout.Rigid(func(gtx layout.Context) layout.Dimensions {
													if node.protocol.Clicked(gtx) {
														p.showProtocolMenu(gtx, node)
													}
													return node.protocol.Layout(gtx, th)
												}),
												layout.Rigid(func(gtx layout.Context) layout.Dimensions {
													return node.host.Layout(gtx, th, i18n.Host.Value())
												}),
												layout.Rigid(func(gtx layout.Context) layout.Dimensions {
													return node.path.Layout(gtx, th, i18n.Path.Value())
												}),
												layout.Rigid(layout.Spacer{Height: 4}.Layout),
											)
										}),

										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											return node.enableHTTP.Layout(gtx, th)
										}),

										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											if !node.enableHTTP.Value() {
												return layout.Dimensions{}
											}

											return layout.Flex{
												Axis: layout.Vertical,
											}.Layout(gtx,
												layout.Rigid(func(gtx layout.Context) layout.Dimensions {
													return node.httpHost.Layout(gtx, th, i18n.RewriteHostHeader.Value())
												}),
												layout.Rigid(func(gtx layout.Context) layout.Dimensions {
													return node.httpUsername.Layout(gtx, th, i18n.Username.Value())
												}),
												layout.Rigid(func(gtx layout.Context) layout.Dimensions {
													return node.httpPassword.Layout(gtx, th, i18n.Password.Value())
												}),
												layout.Rigid(layout.Spacer{Height: 4}.Layout),
											)
										}),

										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											return node.enableTLS.Layout(gtx, th)
										}),

										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											if !node.enableTLS.Value() {
												return layout.Dimensions{}
											}

											return layout.Flex{
												Axis: layout.Vertical,
											}.Layout(gtx,
												layout.Rigid(func(gtx layout.Context) layout.Dimensions {
													return node.tlsSecure.Layout(gtx, th)
												}),

												layout.Rigid(func(gtx layout.Context) layout.Dimensions {
													return node.tlsServerName.Layout(gtx, th, i18n.ServerName.Value())
												}),
												layout.Rigid(layout.Spacer{Height: 4}.Layout),
											)
										}),
									)
								}),
							)
						}),
					)
				})
			}),
		)
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx, children...)
}

func (p *forwarder) showBypassMenu(gtx C, node *node) {
	options := []ui_widget.MenuOption{}
	for _, v := range api.GetConfig().Bypasses {
		options = append(options, ui_widget.MenuOption{
			Value: v.Name,
		})
	}
	for i := range options {
		options[i].Selected = node.bypass.AnyValue(options[i].Value)
	}

	p.menu.Title = i18n.Bypass
	p.menu.Options = options
	p.menu.OnClick = func(ok bool) {
		p.modal.Disappear(gtx.Now)
		if !ok {
			return
		}

		node.bypass.Clear()
		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				node.bypass.Select(ui_widget.SelectorItem{Value: p.menu.Options[i].Value})
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

func (p *forwarder) showProtocolMenu(gtx C, node *node) {
	options := []ui_widget.MenuOption{
		{Key: "HTTP", Value: "http"},
		{Key: "TLS", Value: "tls"},
		{Key: "SSH", Value: "ssh"},
	}

	for i := range options {
		options[i].Selected = node.protocol.AnyValue(options[i].Value)
	}

	p.menu.Title = i18n.Protocol
	p.menu.Options = options
	p.menu.OnClick = func(ok bool) {
		p.modal.Disappear(gtx.Now)
		if !ok {
			return
		}

		node.protocol.Clear()
		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				node.protocol.Select(ui_widget.SelectorItem{Value: p.menu.Options[i].Value})
			}
		}
		p.modal.Disappear(gtx.Now)
	}
	p.menu.ShowAdd = false

	p.modal.Widget = func(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
		return p.menu.Layout(gtx, th)
	}
	p.modal.Appear(gtx.Now)
}
