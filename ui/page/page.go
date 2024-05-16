package page

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/go-gost/gostctl/ui/i18n"
	ui_widget "github.com/go-gost/gostctl/ui/widget"
)

type PagePath string

const (
	PageHome          PagePath = "/"
	PageServer        PagePath = "/server"
	PageService       PagePath = "/service"
	PageChain         PagePath = "/chain"
	PageHop           PagePath = "/hop"
	PageNode          PagePath = "/node"
	PageForwarderNode PagePath = "/forwarder/node"
	PageMetadata      PagePath = "/metadata"
	PageAuther        PagePath = "/auther"
	PageAutherAuths   PagePath = "/auther/auths"
	PageMatcher       PagePath = "/matcher"
	PageAdmission     PagePath = "/admission"
	PageBypass        PagePath = "/bypass"
	PageHosts         PagePath = "/hosts"
	PageHostMapping   PagePath = "/hosts/mapping"
	PageSettings      PagePath = "/settings"
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

type C = layout.Context
type D = layout.Dimensions
type T = material.Theme

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
	Layout(gtx C) D
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

var (
	SelectorStrategyOptions = []ui_widget.MenuOption{
		{Key: i18n.SelectorRound, Value: "round"},
		{Key: i18n.SelectorRandom, Value: "rand"},
		{Key: i18n.SelectorFIFO, Value: "fifo"},
	}

	PluginTypeOptions = []ui_widget.MenuOption{
		{Key: i18n.PluginGRPC, Value: "grpc"},
		{Key: i18n.PluginHTTP, Value: "http"},
	}
)
