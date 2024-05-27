package record

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/go-gost/gostctl/api"
	"github.com/go-gost/gostctl/ui/i18n"
	"github.com/go-gost/gostctl/ui/icons"
	"github.com/go-gost/gostctl/ui/page"
	"github.com/go-gost/gostctl/ui/theme"
	ui_widget "github.com/go-gost/gostctl/ui/widget"
	"github.com/google/uuid"
)

type recordPage struct {
	router *page.Router

	menu ui_widget.Menu
	list layout.List

	btnBack   widget.Clickable
	btnDelete widget.Clickable
	btnEdit   widget.Clickable
	btnSave   widget.Clickable

	name   ui_widget.Selector
	record ui_widget.Selector

	id       string
	perm     page.Perm
	callback page.Callback

	edit   bool
	create bool

	delDialog ui_widget.Dialog
}

func NewPage(r *page.Router) page.Page {
	p := &recordPage{
		router: r,

		list: layout.List{
			// NOTE: the list must be vertical
			Axis: layout.Vertical,
		},
		name:      ui_widget.Selector{Title: i18n.Recorder},
		record:    ui_widget.Selector{Title: i18n.Record},
		delDialog: ui_widget.Dialog{Title: i18n.DeleteRecorder},
	}

	return p
}

func (p *recordPage) Init(opts ...page.PageOption) {
	var options page.PageOptions
	for _, opt := range opts {
		opt(&options)
	}

	p.id = options.ID
	record, _ := options.Value.(*api.RecorderObject)
	if record == nil {
		record = &api.RecorderObject{}
	}
	p.callback = options.Callback

	if p.id != "" {
		p.edit = false
		p.create = false
	} else {
		p.edit = true
		p.create = true
	}

	p.perm = options.Perm

	p.name.Clear()
	if record.Name != "" {
		p.name.Select(ui_widget.SelectorItem{Value: record.Name})
	}

	p.record.Clear()
	if record.Record != "" {
		p.record.Select(ui_widget.SelectorItem{Value: record.Record})
	}
}

func (p *recordPage) Layout(gtx page.C) page.D {
	if p.btnBack.Clicked(gtx) {
		p.router.Back()
	}
	if p.btnEdit.Clicked(gtx) {
		p.edit = true
	}
	if p.btnSave.Clicked(gtx) {
		if p.save() {
			p.router.Back()
		}
	}

	if p.btnDelete.Clicked(gtx) {
		p.delDialog.OnClick = func(ok bool) {
			if ok {
				p.delete()
				p.router.Back()
			}
			p.router.HideModal(gtx)
		}
		p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
			return p.delDialog.Layout(gtx, th)
		})
	}

	th := p.router.Theme

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		// header
		layout.Rigid(func(gtx page.C) page.D {
			return layout.Inset{
				Top:    8,
				Bottom: 8,
				Left:   8,
				Right:  8,
			}.Layout(gtx, func(gtx page.C) page.D {
				return layout.Flex{
					Spacing:   layout.SpaceBetween,
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(func(gtx page.C) page.D {
						btn := material.IconButton(th, &p.btnBack, icons.IconBack, "Back")
						btn.Color = th.Fg
						btn.Background = th.Bg
						return btn.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: 8}.Layout),
					layout.Flexed(1, func(gtx page.C) page.D {
						title := material.H6(th, i18n.Record.Value())
						return title.Layout(gtx)
					}),
					layout.Rigid(func(gtx page.C) page.D {
						if p.perm&page.PermDelete == 0 || p.create {
							return page.D{}
						}
						btn := material.IconButton(th, &p.btnDelete, icons.IconDelete, "Delete")
						btn.Color = th.Fg
						btn.Background = th.Bg
						return btn.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: 8}.Layout),
					layout.Rigid(func(gtx page.C) page.D {
						if p.perm&page.PermWrite == 0 {
							return page.D{}
						}

						if p.edit {
							btn := material.IconButton(th, &p.btnSave, icons.IconDone, "Done")
							btn.Color = th.Fg
							btn.Background = th.Bg
							return btn.Layout(gtx)
						} else {
							btn := material.IconButton(th, &p.btnEdit, icons.IconEdit, "Edit")
							btn.Color = th.Fg
							btn.Background = th.Bg
							return btn.Layout(gtx)
						}
					}),
				)
			})
		}),
		layout.Flexed(1, func(gtx page.C) page.D {
			return p.list.Layout(gtx, 1, func(gtx page.C, index int) page.D {
				return layout.Inset{
					Top:    8,
					Bottom: 8,
					Left:   8,
					Right:  8,
				}.Layout(gtx, func(gtx page.C) page.D {
					return p.layout(gtx, th)
				})
			})
		}),
	)
}

func (p *recordPage) layout(gtx page.C, th *page.T) page.D {
	if !p.edit {
		gtx = gtx.Disabled()
	}

	return component.SurfaceStyle{
		Theme: th,
		ShadowStyle: component.ShadowStyle{
			CornerRadius: 12,
		},
		Fill: theme.Current().ContentSurfaceBg,
	}.Layout(gtx, func(gtx page.C) page.D {
		return layout.UniformInset(16).Layout(gtx, func(gtx page.C) page.D {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx page.C) page.D {
					if p.name.Clicked(gtx) {
						p.showRecorderMenu(gtx)
					}
					return p.name.Layout(gtx, th)
				}),
				layout.Rigid(func(gtx page.C) page.D {
					if p.record.Clicked(gtx) {
						p.showRecordMenu(gtx)
					}
					return p.record.Layout(gtx, th)
				}),
			)
		})
	})
}

func (p *recordPage) showRecorderMenu(gtx page.C) {
	options := []ui_widget.MenuOption{}
	for _, v := range api.GetConfig().Recorders {
		options = append(options, ui_widget.MenuOption{
			Value: v.Name,
		})
	}
	for i := range options {
		options[i].Selected = p.name.AnyValue(options[i].Value)
	}

	p.menu.Title = i18n.Recorder
	p.menu.Options = options
	p.menu.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		p.name.Clear()
		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				p.name.Select(ui_widget.SelectorItem{Value: p.menu.Options[i].Value})
			}
		}
	}
	p.menu.OnAdd = func() {
		p.router.Goto(page.Route{
			Path: page.PageRecorder,
			Perm: page.PermReadWrite,
		})
		p.router.HideModal(gtx)
	}
	p.menu.Multiple = false

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

var (
	recordOptions = []ui_widget.MenuOption{
		{Value: "recorder.service.client.address"},
		{Value: "recorder.service.router.dial.address"},
		{Value: "recorder.service.router.dial.address.error"},
		{Value: "recorder.service.handler.serial"},
	}
)

func (p *recordPage) showRecordMenu(gtx page.C) {
	for i := range recordOptions {
		recordOptions[i].Selected = p.record.AnyValue(recordOptions[i].Value)
	}

	p.menu.Title = i18n.Record
	p.menu.Options = recordOptions
	p.menu.OnClick = func(ok bool) {
		p.router.HideModal(gtx)
		if !ok {
			return
		}

		p.record.Clear()
		for i := range p.menu.Options {
			if p.menu.Options[i].Selected {
				p.record.Select(ui_widget.SelectorItem{Name: p.menu.Options[i].Name, Key: p.menu.Options[i].Key, Value: p.menu.Options[i].Value})
			}
		}
	}
	p.menu.OnAdd = nil
	p.menu.Multiple = false

	p.router.ShowModal(gtx, func(gtx page.C, th *page.T) page.D {
		return p.menu.Layout(gtx, th)
	})
}

func (p *recordPage) generateConfig() *api.RecorderObject {
	ro := &api.RecorderObject{
		Name:   p.name.Value(),
		Record: p.record.Value(),
	}
	return ro
}

func (p *recordPage) save() bool {
	ro := p.generateConfig()

	if p.id == "" {
		if p.callback != nil {
			p.callback(page.ActionCreate, uuid.New().String(), ro)
		}

	} else {
		if p.callback != nil {
			p.callback(page.ActionUpdate, p.id, ro)
		}
	}

	return true
}

func (p *recordPage) delete() {
	if p.callback != nil {
		p.callback(page.ActionDelete, p.id, nil)
	}
}
