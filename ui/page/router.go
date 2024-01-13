package page

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
)

type C = layout.Context
type D = layout.Dimensions

type Router struct {
	pages map[PagePath]Page
}

func NewRouter() *Router {
	r := &Router{
		pages: make(map[PagePath]Page),
	}

	return r
}

func (r *Router) Register(path PagePath, page Page) {
	r.pages[path] = page
}

func (r *Router) Layout(gtx C, th *material.Theme) D {
	return layout.Background{}.Layout(gtx,
		func(gtx C) D {
			// paint.Fill(gtx.Ops, color.NRGBA(colornames.Grey600))
			return D{
				Size: gtx.Constraints.Max,
			}
		},
		func(gtx C) D {
			return r.pages[PageHome].Layout(gtx, th)
		},
	)
}
