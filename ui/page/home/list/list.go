package list

import (
	"gioui.org/widget"
	"github.com/go-gost/gostctl/ui/page"
)

type List interface {
	Layout(gtx page.C, th *page.T) page.D
}

type state struct {
	clk widget.Clickable
}
