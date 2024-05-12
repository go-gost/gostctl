package widget

import (
	"github.com/go-gost/gostctl/ui/page"
)

type Widget interface {
	Layout(gtx page.C, th *page.T) page.D
}
