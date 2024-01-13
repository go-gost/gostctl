package page

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
)

type PagePath string

const (
	PageHome PagePath = "/"
)

type Page interface {
	Layout(gtx layout.Context, th *material.Theme) layout.Dimensions
}
