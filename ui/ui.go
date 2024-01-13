package ui

import (
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/go-gost/gui/ui/page"
	"github.com/go-gost/gui/ui/page/home"
)

type C = layout.Context
type D = layout.Dimensions

type UI struct {
	th     *material.Theme
	router *page.Router
}

func NewUI() *UI {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	router := page.NewRouter()
	router.Register(page.PageHome, home.NewPage(router))

	return &UI{
		th:     th,
		router: router,
	}
}

func (ui *UI) Layout(gtx C) D {
	return ui.router.Layout(gtx, ui.th)
}
