package server

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"time"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/go-gost/gostctl/api"
	"github.com/go-gost/gostctl/api/util"
	"github.com/go-gost/gostctl/config"
	"github.com/go-gost/gostctl/ui/i18n"
	"github.com/go-gost/gostctl/ui/icons"
	"github.com/go-gost/gostctl/ui/page"
	"github.com/go-gost/gostctl/ui/theme"
	ui_widget "github.com/go-gost/gostctl/ui/widget"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

type serverPage struct {
	router *page.Router

	btnBack   widget.Clickable
	btnActive widget.Clickable
	btnDelete widget.Clickable
	btnEdit   widget.Clickable
	btnSave   widget.Clickable

	btnConfig widget.Clickable

	list layout.List

	name component.TextField
	url  component.TextField

	basicAuth          ui_widget.Switcher
	username           component.TextField
	password           component.TextField
	btnPasswordVisible widget.Clickable
	passwordVisible    bool

	interval component.TextField
	timeout  component.TextField

	autoSave ui_widget.Switcher
	saveFile component.TextField
	readonly ui_widget.Switcher

	id   string
	perm page.Perm

	edit   bool
	create bool
	active bool

	delDialog ui_widget.Dialog
}

func NewPage(r *page.Router) page.Page {
	return &serverPage{
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
		basicAuth: ui_widget.Switcher{Title: i18n.BasicAuth},
		username: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     64,
			},
		},
		password: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     64,
			},
		},
		interval: component.TextField{
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

		autoSave: ui_widget.Switcher{Title: i18n.AutoSave},
		saveFile: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     255,
			},
		},
		readonly: ui_widget.Switcher{Title: i18n.Readonly},
		delDialog: ui_widget.Dialog{
			Title: i18n.DeleteServer,
		},
	}
}

func (p *serverPage) Init(opts ...page.PageOption) {
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

	p.active = false
	cfg := config.Get()
	server := &config.Server{}
	for i, srv := range cfg.Servers {
		if srv.Name == p.id {
			server = srv
			if cfg.CurrentServer == i {
				p.active = true
			}
			break
		}
	}

	p.name.Clear()
	p.name.SetText(server.Name)

	p.url.Clear()
	p.url.SetText(server.URL)

	if server.Username != "" {
		p.basicAuth.SetValue(true)
	} else {
		p.basicAuth.SetValue(false)
	}

	p.username.SetText(server.Username)
	p.password.SetText(server.Password)
	p.passwordVisible = false

	p.interval.Clear()
	p.interval.SetText(fmt.Sprintf("%d", int(server.Interval.Seconds())))

	p.timeout.Clear()
	p.timeout.SetText(fmt.Sprintf("%d", int(server.Timeout.Seconds())))

	p.autoSave.SetValue(false)
	p.saveFile.Clear()
	if server.AutoSave != "" {
		p.autoSave.SetValue(true)
		p.saveFile.SetText(server.AutoSave)
	}

	p.readonly.SetValue(server.Readonly)
}

func (p *serverPage) Layout(gtx page.C) page.D {
	if p.btnBack.Clicked(gtx) {
		if page := p.router.Back(); page != nil {
			page.Init()
		}
	}
	if p.btnActive.Clicked(gtx) && !p.active {
		p.activate()
		p.active = true
	}
	if p.btnEdit.Clicked(gtx) {
		p.edit = true
	}
	if p.btnSave.Clicked(gtx) {
		if p.save() {
			if page := p.router.Back(); page != nil {
				page.Init()
			}
		}
	}

	if p.btnDelete.Clicked(gtx) {
		p.delDialog.OnClick = func(ok bool) {
			if ok {
				p.delete()
				if page := p.router.Back(); page != nil {
					page.Init()
				}
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
						title := material.H6(th, i18n.Server.Value())
						return title.Layout(gtx)
					}),
					layout.Rigid(func(gtx page.C) page.D {
						if p.create {
							return page.D{}
						}
						btn := material.IconButton(th, &p.btnActive, icons.IconCircle, "Active")
						btn.Background = th.Bg
						if p.active {
							btn.Color = color.NRGBA(colornames.Green500)
						} else {
							btn.Color = th.Fg
						}
						return btn.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: 8}.Layout),
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
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return p.layout(gtx, th)
				})
			})
		}),
	)
}

func (p *serverPage) layout(gtx page.C, th *page.T) page.D {
	if p.btnConfig.Clicked(gtx) {
		p.router.Goto(page.Route{
			Path:  page.PageConfig,
			Value: api.GetConfig(),
		})
	}

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
					if !p.active {
						return page.D{}
					}

					gtx.Source = src
					return layout.Flex{
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Flexed(1, layout.Spacer{Width: 4}.Layout),
						layout.Rigid(func(gtx page.C) page.D {
							btn := material.IconButton(th, &p.btnConfig, icons.IconCode, "Config")
							btn.Color = th.Fg
							btn.Background = theme.Current().ContentSurfaceBg
							return btn.Layout(gtx)
						}),
					)
				}),

				layout.Rigid(layout.Spacer{Height: 16}.Layout),

				layout.Rigid(func(gtx page.C) page.D {
					return material.Body1(th, i18n.Name.Value()).Layout(gtx)
				}),
				layout.Rigid(func(gtx page.C) page.D {
					return p.name.Layout(gtx, th, "")
				}),
				layout.Rigid(layout.Spacer{Height: 16}.Layout),
				layout.Rigid(func(gtx page.C) page.D {
					return layout.Flex{
						Alignment: layout.Baseline,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return material.Body1(th, i18n.URL.Value()).Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: 4}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return material.Body2(th, "("+i18n.URLHint.Value()+")").Layout(gtx)
						}),
					)
				}),
				layout.Rigid(func(gtx page.C) page.D {
					return p.url.Layout(gtx, th, "")
				}),
				layout.Rigid(layout.Spacer{Height: 16}.Layout),

				layout.Rigid(func(gtx page.C) page.D {
					return layout.Flex{
						Alignment: layout.Baseline,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return material.Body1(th, i18n.Interval.Value()).Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: 4}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return material.Body2(th, "("+i18n.IntervalHint.Value()+")").Layout(gtx)
						}),
					)
				}),
				layout.Rigid(func(gtx page.C) page.D {
					return p.interval.Layout(gtx, th, "")
				}),
				layout.Rigid(layout.Spacer{Height: 16}.Layout),

				layout.Rigid(func(gtx page.C) page.D {
					return layout.Flex{
						Alignment: layout.Baseline,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return material.Body1(th, i18n.Timeout.Value()).Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: 4}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return material.Body2(th, "("+i18n.TimeoutHint.Value()+")").Layout(gtx)
						}),
					)
				}),
				layout.Rigid(func(gtx page.C) page.D {
					return p.timeout.Layout(gtx, th, "")
				}),

				layout.Rigid(layout.Spacer{Height: 16}.Layout),
				layout.Rigid(func(gtx page.C) page.D {
					return p.basicAuth.Layout(gtx, th)
				}),
				layout.Rigid(func(gtx page.C) page.D {
					if !p.basicAuth.Value() {
						return page.D{}
					}

					return layout.UniformInset(8).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis: layout.Vertical,
						}.Layout(gtx,
							layout.Rigid(material.Body1(th, i18n.Username.Value()).Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return p.username.Layout(gtx, th, "")
							}),
							layout.Rigid(layout.Spacer{Height: 8}.Layout),
							layout.Rigid(material.Body1(th, i18n.Password.Value()).Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								if p.btnPasswordVisible.Clicked(gtx) {
									p.passwordVisible = !p.passwordVisible
								}

								if p.passwordVisible {
									p.password.Suffix = func(gtx page.C) page.D {
										return p.btnPasswordVisible.Layout(gtx, func(gtx page.C) page.D {
											return icons.IconVisibility.Layout(gtx, color.NRGBA(colornames.Grey500))
										})
									}
									p.password.Mask = 0
								} else {
									p.password.Suffix = func(gtx page.C) page.D {
										return p.btnPasswordVisible.Layout(gtx, func(gtx page.C) page.D {
											return icons.IconVisibilityOff.Layout(gtx, color.NRGBA(colornames.Grey500))
										})
									}
									p.password.Mask = '*'
								}
								return p.password.Layout(gtx, th, "")
							}),
						)
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return p.readonly.Layout(gtx, th)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return p.autoSave.Layout(gtx, th)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if !p.autoSave.Value() {
						return page.D{}
					}

					return layout.UniformInset(8).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis: layout.Vertical,
						}.Layout(gtx,
							layout.Rigid(material.Body1(th, i18n.FilePath.Value()).Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return p.saveFile.Layout(gtx, th, "")
							}),
						)
					})
				}),
			)
		})
	})
}

func (p *serverPage) activate() {
	cfg := config.Get()
	for i, server := range cfg.Servers {
		if server.Name == p.id {
			cfg.CurrentServer = i
			break
		}
	}
	config.Set(cfg)
	cfg.Write()
	util.RestartGetConfigTask()
}

func (p *serverPage) save() bool {
	server := &config.Server{
		Name: strings.TrimSpace(p.name.Text()),
		URL:  strings.TrimSpace(p.url.Text()),
	}

	if p.basicAuth.Value() {
		server.Username = strings.TrimSpace(p.username.Text())
		server.Password = strings.TrimSpace(p.password.Text())
	}
	if interval, _ := strconv.Atoi(strings.TrimSpace(p.interval.Text())); interval > 0 {
		server.Interval = time.Duration(interval) * time.Second
	}
	if timeout, _ := strconv.Atoi(strings.TrimSpace(p.timeout.Text())); timeout > 0 {
		server.Timeout = time.Duration(timeout) * time.Second
	}
	server.Readonly = p.readonly.Value()
	if p.autoSave.Value() {
		server.AutoSave = strings.TrimSpace(p.saveFile.Text())
	}

	cfg := config.Get()

	ok := func() bool {
		ok := true

		if server.Name == "" {
			p.name.SetError(i18n.ErrNameRequired.Value())
			ok = false
		}

		if p.create {
			for _, srv := range cfg.Servers {
				if srv.Name == server.Name {
					p.name.SetError(i18n.ErrNameExists.Value())
					ok = false
					break
				}
			}
		}

		if server.URL == "" {
			p.url.SetError(i18n.ErrURLRequired.Value())
		}

		return ok
	}()
	if !ok {
		return false
	}

	servers := make([]*config.Server, len(cfg.Servers))
	copy(servers, cfg.Servers)

	if p.create {
		servers = append(servers, server)
	} else {
		for i := range servers {
			if servers[i].Name == server.Name {
				servers[i] = server
			}
		}
	}
	cfg.Servers = servers

	config.Set(cfg)
	cfg.Write()

	util.RestartGetConfigTask()

	return true
}

func (p *serverPage) delete() {
	var servers []*config.Server

	cfg := config.Get()
	for _, server := range cfg.Servers {
		if server.Name == p.id {
			if p.active {
				cfg.CurrentServer = 0
			}
			continue
		}
		servers = append(servers, server)
	}
	cfg.Servers = servers
	config.Set(cfg)
	cfg.Write()

	if p.active {
		util.RestartGetConfigTask()
	}
}
