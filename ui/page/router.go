package page

import (
	"fmt"
	"log/slog"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/go-gost/gostctl/ui/theme"
	ui_widget "github.com/go-gost/gostctl/ui/widget"
)

const (
	MaxWidth = 800
)

type Route struct {
	Path     PagePath
	ID       string
	Value    any
	Callback Callback
	Perm     Perm
}

type Router struct {
	w       *app.Window
	pages   map[PagePath]Page
	stack   routeStack
	current Route
	*material.Theme
	modal        *component.ModalLayer
	notification *ui_widget.Notification
}

func NewRouter(w *app.Window, th *T) *Router {
	r := &Router{
		w:     w,
		pages: make(map[PagePath]Page),
		Theme: th,
		modal: component.NewModal(),
		notification: ui_widget.NewNotification(3*time.Second, func() {
			w.Invalidate()
		}),
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

	page.Init(
		WithPageID(route.ID),
		WithPageValue(route.Value),
		WithPageCallback(route.Callback),
		WithPagePerm(route.Perm),
	)
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

	// page.Init(WithPageID(route.ID))
	slog.Debug(fmt.Sprintf("back to %s", route.Path), "kind", "router", "route.id", route.ID)
}

func (r *Router) Layout(gtx C) D {
	r.Theme.Palette = theme.Current().Material

	defer r.modal.Layout(gtx, r.Theme)

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
			width := unit.Dp(MaxWidth)
			if x := gtx.Metric.PxToDp(gtx.Constraints.Max.X); x > width {
				inset.Left = (x - width) / 2
				inset.Right = inset.Left
			}
			return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Stack{
					Alignment: layout.N,
				}.Layout(gtx,
					layout.Expanded(page.Layout),
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{
							Top:    16,
							Bottom: 16,
							Right:  gtx.Metric.PxToDp(gtx.Constraints.Max.X) / 5,
							Left:   gtx.Metric.PxToDp(gtx.Constraints.Max.X) / 5,
						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return r.notification.Layout(gtx, r.Theme)
						})
					}),
				)
			})
		},
	)
}

func (r *Router) ShowModal(gtx C, w func(gtx C, th *T) D) {
	r.modal.Widget = func(gtx C, th *T, anim *component.VisibilityAnimation) D {
		if gtx.Constraints.Max.X > gtx.Dp(MaxWidth) {
			gtx.Constraints.Max.X = gtx.Dp(MaxWidth)
		}
		gtx.Constraints.Max.X = gtx.Constraints.Max.X * 3 / 4

		var clk widget.Clickable
		return clk.Layout(gtx, func(gtx C) D {
			return w(gtx, th)
		})
	}
	r.modal.Appear(gtx.Now)
}

func (r *Router) HideModal(gtx C) {
	r.modal.Disappear(gtx.Now)
}

func (r *Router) Notify(message ui_widget.Message) {
	r.notification.Show(message)
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
