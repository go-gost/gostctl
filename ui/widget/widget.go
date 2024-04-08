package widget

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
)

type Widget interface {
	Layout(gtx layout.Context, th *material.Theme) layout.Dimensions
}
