package page

import (
	"fmt"
	"log/slog"

	"gioui.org/layout"
	"gioui.org/widget/material"
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
	slog.Debug(fmt.Sprintf("go to %s %s", route.Path, route.ID), "kind", "router")
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
	slog.Debug(fmt.Sprintf("back to %s %s", route.Path, route.ID), "kind", "router")
}

func (r *Router) Layout(gtx C) D {
	return layout.Background{}.Layout(gtx,
		func(gtx C) D {
			return D{
				Size: gtx.Constraints.Max,
			}
		},
		func(gtx C) D {
			page := r.pages[r.current.Path]
			if page == nil {
				page = r.pages[PageHome]
			}
			return page.Layout(gtx)
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
