package page

import (
	"gioui.org/layout"
)

type PagePath string

const (
	PageHome     PagePath = "/"
	PageServer   PagePath = "/server"
	PageService  PagePath = "/service"
	PageChain    PagePath = "/chain"
	PageHop      PagePath = "/hop"
	PageNode     PagePath = "/node"
	PageSettings PagePath = "/settings"
)

type Perm uint8

const (
	PermRead   Perm = 1
	PermWrite  Perm = 2
	PermDelete Perm = 4

	PermReadWrite       Perm = PermRead | PermWrite
	PermReadWriteDelete Perm = PermReadWrite | PermDelete
)

type Action string

const (
	ActionCreate Action = "create"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
)

type Callback func(action Action, id string, value any)

type PageOptions struct {
	ID       string
	Value    any
	Callback Callback
	Perm     Perm
}

type PageOption func(opts *PageOptions)

func WithPageID(id string) PageOption {
	return func(opts *PageOptions) {
		opts.ID = id
	}
}

func WithPageValue(v any) PageOption {
	return func(opts *PageOptions) {
		opts.Value = v
	}
}

func WithPageCallback(cb Callback) PageOption {
	return func(opts *PageOptions) {
		opts.Callback = cb
	}
}

func WithPagePerm(perm Perm) PageOption {
	return func(opts *PageOptions) {
		opts.Perm = perm
	}
}

type Page interface {
	Init(opts ...PageOption)
	Layout(gtx layout.Context) layout.Dimensions
}

type PageMode string

const (
	BasicMode    PageMode = "basic"
	AdvancedMode PageMode = "advanced"
	PluginMode   PageMode = "plugin"
)

type Metadata struct {
	K string
	V string
}

type AuthType string

const (
	AuthSimple AuthType = "simple"
	AuthAuther AuthType = "auther"
)
