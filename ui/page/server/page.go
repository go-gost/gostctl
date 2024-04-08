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
	"github.com/go-gost/gostctl/api/util"
	"github.com/go-gost/gostctl/config"
	"github.com/go-gost/gostctl/ui/i18n"
	"github.com/go-gost/gostctl/ui/icons"
	"github.com/go-gost/gostctl/ui/page"
	"github.com/go-gost/gostctl/ui/theme"
	ui_widget "github.com/go-gost/gostctl/ui/widget"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

type C = layout.Context
type D = layout.Dimensions

type serverPage struct {
	router *page.Router
	modal  *component.ModalLayer

	btnBack   widget.Clickable
	btnActive widget.Clickable
	btnDelete widget.Clickable
	btnEdit   widget.Clickable
	btnSave   widget.Clickable

	list layout.List

	name component.TextField
	url  component.TextField

	basicAuth ui_widget.Switcher
	username  component.TextField
	password  component.TextField

	interval component.TextField
	timeout  component.TextField

	btnPasswordVisible widget.Clickable
	passwordVisible    bool

	id     string
	edit   bool
	create bool
	active bool

	delDialog ui_widget.Dialog
}

func NewPage(r *page.Router) page.Page {
	return &serverPage{
		router: r,

		modal: component.NewModal(),

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
				MaxLen:     10,
			},
		},
		timeout: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				MaxLen:     10,
			},
		},
		delDialog: ui_widget.Dialog{
			Title: i18n.DeleteServer,
		},

		basicAuth: ui_widget.Switcher{Title: i18n.BasicAuth},
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

	p.username.Clear()
	p.username.SetText(server.Username)

	p.password.Clear()
	p.password.SetText(server.Password)
	p.passwordVisible = false

	p.interval.Clear()
	p.interval.SetText(fmt.Sprintf("%d", int(server.Interval.Seconds())))

	p.timeout.Clear()
	p.timeout.SetText(fmt.Sprintf("%d", int(server.Timeout.Seconds())))

}

func (p *serverPage) Layout(gtx C) D {
	if p.btnBack.Clicked(gtx) {
		p.router.Back()
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
			p.router.Back()
		}
	}

	if p.btnDelete.Clicked(gtx) {
		p.delDialog.OnClick = func(ok bool) {
			if ok {
				p.delete()
				p.router.Back()
			}
			p.modal.Disappear(gtx.Now)
		}
		p.modal.Widget = func(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
			return p.delDialog.Layout(gtx, th)
		}
		p.modal.Appear(gtx.Now)
	}

	th := p.router.Theme

	defer p.modal.Layout(gtx, th)

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		// header
		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Top:    8,
				Bottom: 8,
				Left:   8,
				Right:  8,
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
					layout.Rigid(layout.Spacer{Width: 8}.Layout),
					layout.Flexed(1, func(gtx C) D {
						title := material.H6(th, "Server")
						// title.Font.Weight = font.Bold
						return title.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						if p.create {
							return D{}
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
					layout.Rigid(func(gtx C) D {
						if p.create {
							return D{}
						}
						btn := material.IconButton(th, &p.btnDelete, icons.IconDelete, "Delete")

						btn.Color = th.Fg
						btn.Background = th.Bg
						return btn.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: 8}.Layout),
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
			return p.list.Layout(gtx, 1, func(gtx C, _ int) D {
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

func (p *serverPage) layout(gtx C, th *material.Theme) D {
	if !p.edit {
		gtx = gtx.Disabled()
	}

	return component.SurfaceStyle{
		Theme: th,
		ShadowStyle: component.ShadowStyle{
			CornerRadius: 12,
		},
		Fill: theme.Current().ContentSurfaceBg,
	}.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(16).Layout(gtx, func(gtx C) D {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return material.Body1(th, i18n.Name.Value()).Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return p.name.Layout(gtx, th, "")
				}),
				layout.Rigid(layout.Spacer{Height: 16}.Layout),
				layout.Rigid(func(gtx C) D {
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
				layout.Rigid(func(gtx C) D {
					return p.url.Layout(gtx, th, "")
				}),
				layout.Rigid(layout.Spacer{Height: 16}.Layout),

				layout.Rigid(func(gtx C) D {
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
				layout.Rigid(func(gtx C) D {
					p.interval.Suffix = material.Body1(th, "s").Layout
					return p.interval.Layout(gtx, th, i18n.Seconds.Value())
				}),
				layout.Rigid(layout.Spacer{Height: 16}.Layout),

				layout.Rigid(func(gtx C) D {
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
				layout.Rigid(func(gtx C) D {
					p.timeout.Suffix = material.Body1(th, "s").Layout
					return p.timeout.Layout(gtx, th, i18n.Seconds.Value())
				}),

				layout.Rigid(layout.Spacer{Height: 16}.Layout),
				layout.Rigid(func(gtx C) D {
					return p.basicAuth.Layout(gtx, th)
				}),
				layout.Rigid(func(gtx C) D {
					if !p.basicAuth.Value() {
						p.username.SetText("")
						return D{}
					}
					return p.username.Layout(gtx, th, i18n.Username.Value())
				}),
				layout.Rigid(func(gtx C) D {
					if !p.basicAuth.Value() {
						p.password.SetText("")
						return D{}
					}

					if p.btnPasswordVisible.Clicked(gtx) {
						p.passwordVisible = !p.passwordVisible
					}

					if p.passwordVisible {
						p.password.Suffix = func(gtx C) D {
							return p.btnPasswordVisible.Layout(gtx, func(gtx C) D {
								return icons.IconVisibility.Layout(gtx, color.NRGBA(colornames.Grey500))
							})
						}
						p.password.Mask = 0
					} else {
						p.password.Suffix = func(gtx C) D {
							return p.btnPasswordVisible.Layout(gtx, func(gtx C) D {
								return icons.IconVisibilityOff.Layout(gtx, color.NRGBA(colornames.Grey500))
							})
						}
						p.password.Mask = '*'
					}

					return p.password.Layout(gtx, th, i18n.Password.Value())
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
