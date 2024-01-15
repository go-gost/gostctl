package page

import (
	"gioui.org/layout"
)

type PagePath string

const (
	PageHome        PagePath = "/"
	PageServerEdit  PagePath = "/server/edit"
	PageServiceEdit PagePath = "/service/edit"
)

type PageOptions struct {
	ID string
}

type PageOption func(opts *PageOptions)

func WithPageID(id string) PageOption {
	return func(opts *PageOptions) {
		opts.ID = id
	}
}

type Page interface {
	Init(opts ...PageOption)
	Layout(gtx layout.Context) layout.Dimensions
}
