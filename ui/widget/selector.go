package widget

import (
	"strings"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/go-gost/gostctl/ui/i18n"
	"github.com/go-gost/gostctl/ui/icons"
)

type SelectorItem struct {
	Name  i18n.Key
	Value string
}

type Selector struct {
	Title     i18n.Key
	items     []SelectorItem
	clickable widget.Clickable
}

func (p *Selector) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return material.Clickable(gtx, &p.clickable, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top:    8,
			Bottom: 8,
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(material.Body1(th, p.Title.Value()).Layout),
				layout.Rigid(layout.Spacer{Width: 8}.Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						var values []string
						for _, item := range p.items {
							if item.Value == "" {
								continue
							}

							value := item.Name.Value()
							if value == "" {
								value = item.Value
							}
							values = append(values, value)
						}
						return material.Body2(th, strings.Join(values, ",")).Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return icons.IconNavRight.Layout(gtx, th.Fg)
				}),
			)
		})
	})
}

func (p *Selector) Clicked(gtx layout.Context) bool {
	return p.clickable.Clicked(gtx)
}

func (p *Selector) Select(items ...SelectorItem) {
	for _, item := range items {
		if item.Name == "" && item.Value == "" {
			continue
		}
		p.items = append(p.items, item)
	}
}

func (p *Selector) Any(items ...SelectorItem) bool {
	for _, item := range items {
		if p.contains(item) {
			return true
		}
	}
	return false
}

func (p *Selector) contains(item SelectorItem) bool {
	for i := range p.items {
		if p.items[i] == item {
			return true
		}
	}
	return false
}

func (p *Selector) AnyValue(values ...string) bool {
	for _, value := range values {
		for i := range p.items {
			if p.items[i].Value == value {
				return true
			}
		}
	}
	return false
}

func (p *Selector) Item() SelectorItem {
	if len(p.items) == 0 {
		return SelectorItem{}
	}
	return p.items[0]
}

func (p *Selector) Items() []SelectorItem {
	return p.items
}

func (p *Selector) Value() string {
	if len(p.items) == 0 {
		return ""
	}
	return p.items[0].Value
}

func (p *Selector) Values() (values []string) {
	for i := range p.items {
		values = append(values, p.items[i].Value)
	}
	return
}

func (p *Selector) Clear() {
	p.items = nil
}
