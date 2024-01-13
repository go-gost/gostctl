package list

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
)

type C = layout.Context
type D = layout.Dimensions

type List interface {
	Layout(gtx C, th *material.Theme) D
}
