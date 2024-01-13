package server

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"time"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/go-gost/gui/api/util"
	"github.com/go-gost/gui/config"
	"github.com/go-gost/gui/ui/icons"
	"github.com/go-gost/gui/ui/page"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

type C = layout.Context
type D = layout.Dimensions

type serverPage struct {
	router *page.Router

	btnBack   widget.Clickable
	btnDelete widget.Clickable
	btnEdit   widget.Clickable
	btnSave   widget.Clickable

	list layout.List

	name component.TextField
	url  component.TextField

	basicAuth widget.Bool
	username  component.TextField
	password  component.TextField
	interval  component.TextField
	timeout   component.TextField

	btnPasswordVisible widget.Clickable
	passwordVisible    bool

	edit   bool
	create bool
	id     string
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

	var server config.Server
	for _, srv := range config.Global().Servers {
		if srv.Name == p.id {
			server = srv
			break
		}
	}

	p.name.Clear()
	p.name.SetText(server.Name)

	p.url.Clear()
	p.url.SetText(server.URL)

	if server.Username != "" {
		p.basicAuth.Value = true
	} else {
		p.basicAuth.Value = false
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

func (p *serverPage) Layout(gtx C, th *material.Theme) D {
	if p.btnBack.Clicked(gtx) {
		p.router.Back()
	}
	if p.btnEdit.Clicked(gtx) {
		p.edit = true
	}
	if p.btnSave.Clicked(gtx) {
		if p.save() {
			util.RestartGetConfigTask()
			p.router.Back()
		}
	}
	if p.btnDelete.Clicked(gtx) {
		p.delete()
		p.router.Back()
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		// header
		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Top:    5,
				Bottom: 5,
				Left:   10,
				Right:  10,
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
					layout.Rigid(layout.Spacer{Width: 10}.Layout),
					layout.Flexed(1, func(gtx C) D {
						title := material.H6(th, "Server")
						title.Font.Weight = font.Bold
						return title.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						if p.create {
							return D{}
						}
						btn := material.IconButton(th, &p.btnDelete, icons.IconDelete, "Delete")
						btn.Color = th.Fg
						btn.Background = th.Bg
						return btn.Layout(gtx)
					}),
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
		layout.Rigid(layout.Spacer{Height: 10}.Layout),
		layout.Flexed(1, func(gtx C) D {
			inset := layout.Inset{
				Top:    5,
				Bottom: 5,
			}
			width := unit.Dp(800)
			if x := gtx.Metric.PxToDp(gtx.Constraints.Max.X); x > width {
				inset.Left = (x - width) / 2
				inset.Right = inset.Left
			}
			return inset.Layout(gtx, func(gtx C) D {
				return p.list.Layout(gtx, 1, func(gtx C, index int) D {
					return layout.Inset{
						Top:    10,
						Bottom: 10,
						Left:   10,
						Right:  10,
					}.Layout(gtx, func(gtx C) D {
						return p.layout(gtx, th)
					})
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
			CornerRadius: 20,
		},
		Fill: color.NRGBA(colornames.Grey50),
	}.Layout(gtx, func(gtx C) D {
		return layout.Inset{
			Top:    10,
			Bottom: 10,
			Left:   10,
			Right:  10,
		}.Layout(gtx, func(gtx C) D {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return material.Body1(th, "Name").Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return p.name.Layout(gtx, th, "")
				}),
				layout.Rigid(layout.Spacer{Height: 10}.Layout),
				layout.Rigid(func(gtx C) D {
					return material.Body1(th, "URL (e.g. http://localhost:8000)").Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return p.url.Layout(gtx, th, "")
				}),
				layout.Rigid(layout.Spacer{Height: 10}.Layout),

				layout.Rigid(func(gtx C) D {
					return material.Body1(th, "Interval, the period for obtaining configuration").Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					p.interval.Suffix = material.Body1(th, "s").Layout
					return p.interval.Layout(gtx, th, "Seconds")
				}),
				layout.Rigid(layout.Spacer{Height: 10}.Layout),

				layout.Rigid(func(gtx C) D {
					return material.Body1(th, "Timeout, request timeout when obtaining configuration").Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					p.timeout.Suffix = material.Body1(th, "s").Layout
					return p.timeout.Layout(gtx, th, "Seconds")
				}),

				layout.Rigid(layout.Spacer{Height: 10}.Layout),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: 10, Bottom: 10}.Layout(gtx, func(gtx C) D {
						return layout.Flex{
							Spacing:   layout.SpaceBetween,
							Alignment: layout.Middle,
						}.Layout(gtx,
							layout.Flexed(1, material.Body1(th, "Use basic auth").Layout),
							layout.Rigid(material.Switch(th, &p.basicAuth, "use basic auth").Layout),
						)
					})
				}),
				layout.Rigid(func(gtx C) D {
					if !p.basicAuth.Value {
						p.username.SetText("")
						return D{}
					}
					return p.username.Layout(gtx, th, "Username")
				}),
				layout.Rigid(func(gtx C) D {
					if !p.basicAuth.Value {
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

					return p.password.Layout(gtx, th, "Password")
				}),
			)
		})
	})
}

func (p *serverPage) save() bool {
	server := config.Server{
		Name: strings.TrimSpace(p.name.Text()),
		URL:  strings.TrimSpace(p.url.Text()),
	}

	if p.basicAuth.Value {
		server.Username = strings.TrimSpace(p.username.Text())
		server.Password = strings.TrimSpace(p.password.Text())
	}
	if interval, _ := strconv.Atoi(strings.TrimSpace(p.interval.Text())); interval > 0 {
		server.Interval = time.Duration(interval) * time.Second
	}
	if timeout, _ := strconv.Atoi(strings.TrimSpace(p.timeout.Text())); timeout > 0 {
		server.Timeout = time.Duration(timeout) * time.Second
	}

	cfg := config.Global()

	ok := func() bool {
		ok := true

		if server.Name == "" {
			p.name.SetError("Name is required")
			ok = false
		}

		if p.create {
			for _, svc := range cfg.Servers {
				if svc.Name == server.Name {
					p.name.SetError("Name already exists")
					ok = false
					break
				}
			}
		}

		if server.URL == "" {
			p.url.SetError("URL is required")
		}

		return ok
	}()
	if !ok {
		return false
	}

	if p.create {
		cfg.Servers = append(cfg.Servers, server)
	} else {
		for i := range cfg.Servers {
			if cfg.Servers[i].Name == server.Name {
				cfg.Servers[i] = server
			}
		}
	}
	config.Set(cfg)
	cfg.Write()

	return true
}

func (p *serverPage) delete() {
	var servers []config.Server

	cfg := config.Global()
	for _, server := range cfg.Servers {
		if server.Name == p.id {
			continue
		}
		servers = append(servers, server)
	}
	cfg.Servers = servers
	config.Set(cfg)
	cfg.Write()
}
