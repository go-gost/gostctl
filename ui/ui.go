package ui

import (
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/go-gost/gostctl/config"
	"github.com/go-gost/gostctl/ui/i18n"
	"github.com/go-gost/gostctl/ui/page"
	"github.com/go-gost/gostctl/ui/page/chain"
	"github.com/go-gost/gostctl/ui/page/home"
	"github.com/go-gost/gostctl/ui/page/hop"
	"github.com/go-gost/gostctl/ui/page/node"
	"github.com/go-gost/gostctl/ui/page/server"
	"github.com/go-gost/gostctl/ui/page/service"
	"github.com/go-gost/gostctl/ui/page/settings"
	"github.com/go-gost/gostctl/ui/theme"
)

type C = layout.Context
type D = layout.Dimensions

type UI struct {
	router *page.Router
}

func NewUI() *UI {
	if settings := config.Get().Settings; settings != nil {
		switch settings.Theme {
		case theme.Dark:
			theme.UseDark()
		default:
			theme.UseLight()
		}
		i18n.Set(settings.Lang)
	}

	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	// th.Shaper = text.NewShaper(text.WithCollection(fonts.Collection()))
	th.Palette = theme.Current().Material

	router := page.NewRouter(th)
	router.Register(page.PageHome, home.NewPage(router))
	router.Register(page.PageServer, server.NewPage(router))
	router.Register(page.PageService, service.NewPage(router))
	router.Register(page.PageChain, chain.NewPage(router))
	router.Register(page.PageHop, hop.NewPage(router))
	router.Register(page.PageNode, node.NewPage(router))
	router.Register(page.PageSettings, settings.NewPage(router))

	router.Goto(page.Route{
		Path: page.PageHome,
	})

	return &UI{
		router: router,
	}
}

func (ui *UI) Layout(gtx C) D {
	return ui.router.Layout(gtx)
}
