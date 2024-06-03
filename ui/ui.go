package ui

import (
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/go-gost/gostctl/config"
	"github.com/go-gost/gostctl/ui/fonts"
	"github.com/go-gost/gostctl/ui/i18n"
	"github.com/go-gost/gostctl/ui/page"
	"github.com/go-gost/gostctl/ui/page/admission"
	"github.com/go-gost/gostctl/ui/page/auther"
	"github.com/go-gost/gostctl/ui/page/auther/auth"
	"github.com/go-gost/gostctl/ui/page/bypass"
	"github.com/go-gost/gostctl/ui/page/chain"
	page_config "github.com/go-gost/gostctl/ui/page/config"
	page_event "github.com/go-gost/gostctl/ui/page/event"
	"github.com/go-gost/gostctl/ui/page/home"
	"github.com/go-gost/gostctl/ui/page/hop"
	"github.com/go-gost/gostctl/ui/page/hosts"
	"github.com/go-gost/gostctl/ui/page/hosts/mapping"
	"github.com/go-gost/gostctl/ui/page/limiter"
	"github.com/go-gost/gostctl/ui/page/limiter/limit"
	"github.com/go-gost/gostctl/ui/page/matcher"
	"github.com/go-gost/gostctl/ui/page/node"
	"github.com/go-gost/gostctl/ui/page/observer"
	"github.com/go-gost/gostctl/ui/page/recorder"
	"github.com/go-gost/gostctl/ui/page/resolver"
	"github.com/go-gost/gostctl/ui/page/resolver/nameserver"
	"github.com/go-gost/gostctl/ui/page/server"
	"github.com/go-gost/gostctl/ui/page/service"
	forwarder_node "github.com/go-gost/gostctl/ui/page/service/node"
	"github.com/go-gost/gostctl/ui/page/service/record"
	"github.com/go-gost/gostctl/ui/page/settings"
	"github.com/go-gost/gostctl/ui/theme"
)

type C = layout.Context
type D = layout.Dimensions

type UI struct {
	w      *app.Window
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
	// th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	th.Shaper = text.NewShaper(text.WithCollection(fonts.Collection()))
	th.Palette = theme.Current().Material

	w := &app.Window{}
	w.Option(
		app.Title("GOST"),
		app.MinSize(800, 600),
		app.StatusColor(th.Bg),
	)

	router := page.NewRouter(w, th)
	router.Register(page.PageHome, home.NewPage(router))
	router.Register(page.PageServer, server.NewPage(router))
	router.Register(page.PageService, service.NewPage(router))
	router.Register(page.PageServiceRecord, record.NewPage(router))
	router.Register(page.PageChain, chain.NewPage(router))
	router.Register(page.PageHop, hop.NewPage(router))
	router.Register(page.PageNode, node.NewPage(router))
	router.Register(page.PageForwarderNode, forwarder_node.NewPage(router))
	router.Register(page.PageAuther, auther.NewPage(router))
	router.Register(page.PageAutherAuths, auth.NewPage(router))
	router.Register(page.PageMatcher, matcher.NewPage(router))
	router.Register(page.PageAdmission, admission.NewPage(router))
	router.Register(page.PageBypass, bypass.NewPage(router))
	router.Register(page.PageResolver, resolver.NewPage(router))
	router.Register(page.PageNameServer, nameserver.NewPage(router))
	router.Register(page.PageHosts, hosts.NewPage(router))
	router.Register(page.PageHostMapping, mapping.NewPage(router))
	router.Register(page.PageLimiter, limiter.NewPage(router))
	router.Register(page.PageLimit, limit.NewPage(router))
	router.Register(page.PageObserver, observer.NewPage(router))
	router.Register(page.PageRecorder, recorder.NewPage(router))
	router.Register(page.PageEvent, page_event.NewPage(router))
	router.Register(page.PageConfig, page_config.NewPage(router))
	router.Register(page.PageSettings, settings.NewPage(router))

	router.Goto(page.Route{
		Path: page.PageHome,
	})

	return &UI{
		w:      w,
		router: router,
	}
}

func (ui *UI) Layout(gtx C) D {
	return ui.router.Layout(gtx)
}

func (ui *UI) Window() *app.Window {
	return ui.w
}

func (ui *UI) Router() *page.Router {
	return ui.router
}
