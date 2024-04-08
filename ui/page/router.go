package page

import (
	"fmt"
	"log/slog"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/go-gost/gostctl/ui/theme"
)

type C = layout.Context
type D = layout.Dimensions

type Route struct {
	Path PagePath
	ID   string
}

type Router struct {
	pages   map[PagePath]Page
	stack   routeStack
	current Route
	*material.Theme
}

func NewRouter(th *material.Theme) *Router {
	r := &Router{
		pages: make(map[PagePath]Page),
		Theme: th,
	}

	return r
}

func (r *Router) Register(path PagePath, page Page) {
	r.pages[path] = page
}

func (r *Router) Goto(route Route) {
	page := r.pages[route.Path]
	if page == nil {
		return
	}

	r.current = route
	r.stack.Push(route)

	page.Init(WithPageID(route.ID))
	slog.Debug(fmt.Sprintf("go to %s", route.Path), "kind", "router", "route.id", route.ID)
}

func (r *Router) Back() {
	r.stack.Pop()
	route := r.stack.Peek()

	page := r.pages[route.Path]
	if page == nil {
		return
	}
	r.current = route

	page.Init(WithPageID(route.ID))
	slog.Debug(fmt.Sprintf("back to %s", route.Path), "kind", "router", "route.id", route.ID)
}

func (r *Router) Layout(gtx C) D {
	r.Theme.Palette = theme.Current().Material

	return layout.Background{}.Layout(gtx,
		func(gtx C) D {
			defer clip.Rect{
				Max: gtx.Constraints.Max,
			}.Op().Push(gtx.Ops).Pop()

			paint.ColorOp{
				Color: r.Theme.Bg,
			}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)

			return layout.Dimensions{
				Size: gtx.Constraints.Max,
			}
		},
		func(gtx C) D {
			page := r.pages[r.current.Path]
			if page == nil {
				page = r.pages[PageHome]
			}

			inset := layout.Inset{}
			width := unit.Dp(800)
			if x := gtx.Metric.PxToDp(gtx.Constraints.Max.X); x > width {
				inset.Left = (x - width) / 2
				inset.Right = inset.Left
			}
			return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return page.Layout(gtx)
			})
		},
	)
}

type routeStack struct {
	routes []Route
}

func (p *routeStack) Push(route Route) {
	p.routes = append(p.routes, route)
}

func (p *routeStack) Pop() (route Route) {
	if len(p.routes) == 0 {
		return
	}

	n := len(p.routes) - 1
	route = p.routes[n]
	p.routes = p.routes[:n]

	return
}

func (p *routeStack) Peek() (route Route) {
	if len(p.routes) == 0 {
		return
	}

	return p.routes[len(p.routes)-1]
}
