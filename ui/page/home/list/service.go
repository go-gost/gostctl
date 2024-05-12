package list

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"
	"time"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/go-gost/gostctl/api"
	"github.com/go-gost/gostctl/ui/i18n"
	"github.com/go-gost/gostctl/ui/icons"
	"github.com/go-gost/gostctl/ui/page"
	"github.com/go-gost/gostctl/ui/theme"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

type serviceList struct {
	router *page.Router
	list   layout.List
	states []state
}

func Service(r *page.Router) List {
	return &serviceList{
		router: r,
		list: layout.List{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		},
		states: make([]state, 16),
	}
}

func (l *serviceList) Layout(gtx page.C, th *page.T) page.D {
	cfg := api.GetConfig()
	services := cfg.Services

	if len(services) > len(l.states) {
		states := l.states
		l.states = make([]state, len(services))
		copy(l.states, states)
	}

	return l.list.Layout(gtx, len(services), func(gtx page.C, index int) page.D {
		if l.states[index].clk.Clicked(gtx) {
			l.router.Goto(page.Route{
				Path: page.PageService,
				ID:   services[index].Name,
				Perm: page.PermReadWriteDelete,
			})
		}

		service := services[index]
		handler := service.Handler
		if handler == nil {
			handler = &api.HandlerConfig{
				Type: "auto",
			}
		}
		listener := service.Listener
		if listener == nil {
			listener = &api.ListenerConfig{
				Type: "tcp",
			}
		}
		status := service.Status

		return layout.Inset{
			Top:    8,
			Bottom: 8,
			Left:   8,
			Right:  8,
		}.Layout(gtx, func(gtx page.C) page.D {
			return material.ButtonLayoutStyle{
				Background:   theme.Current().ListBg,
				CornerRadius: 12,
				Button:       &l.states[index].clk,
			}.Layout(gtx, func(gtx page.C) page.D {
				return layout.UniformInset(16).Layout(gtx, func(gtx page.C) page.D {
					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(gtx,
						// title, state
						layout.Rigid(func(gtx page.C) page.D {
							var state string
							if status != nil {
								state = status.State
							}

							c := colornames.Grey500
							key := i18n.Unknown

							switch state {
							case "running":
								c = colornames.DeepOrange500
								key = i18n.Running
							case "ready":
								c = colornames.Green500
								key = i18n.Ready
							case "failed":
								c = colornames.Red500
								key = i18n.Failed
							case "closed":
								c = colornames.Grey500
								key = i18n.Closed
							default:
								key = i18n.Unknown
							}

							return layout.Flex{
								Alignment: layout.Middle,
								Spacing:   layout.SpaceBetween,
							}.Layout(gtx,
								layout.Flexed(1, func(gtx page.C) page.D {
									label := material.Body1(th, service.Name)
									label.Font.Weight = font.SemiBold
									return label.Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Width: 4}.Layout),
								layout.Rigid(func(gtx page.C) page.D {
									gtx.Constraints.Min.X = gtx.Dp(10)
									return icons.IconCircle.Layout(gtx, color.NRGBA(c))
								}),
								layout.Rigid(layout.Spacer{Width: 4}.Layout),
								layout.Rigid(material.Body2(th, key.Value()).Layout),
							)
						}),
						layout.Rigid(layout.Spacer{Height: 4}.Layout),
						layout.Rigid(func(gtx page.C) page.D {
							return layout.Flex{
								Alignment: layout.Middle,
								Spacing:   layout.SpaceBetween,
							}.Layout(gtx,
								layout.Flexed(1, func(gtx page.C) page.D {
									return material.Body2(th, service.Addr).Layout(gtx)
								}),
								layout.Rigid(func(gtx page.C) page.D {
									if service.Status != nil && service.Status.CreateTime > 0 {
										createdAt := time.Unix(service.Status.CreateTime, 0)
										v, unit := formatDuration(time.Since(createdAt))
										return material.Body2(th, fmt.Sprintf("%d%s", v, unit)).Layout(gtx)
									}
									return material.Body2(th, "N/A").Layout(gtx)
								}),
							)
						}),
						// layout.Rigid(material.Body2(th, fmt.Sprintf("Type: %s, %s", handler.Type, listener.Type)).Layout),
						layout.Rigid(layout.Spacer{Height: 4}.Layout),
						layout.Rigid(func(gtx page.C) page.D {
							return layout.Flex{
								Alignment: layout.Middle,
								Spacing:   layout.SpaceBetween,
							}.Layout(gtx,
								layout.Rigid(func(gtx page.C) page.D {
									return icons.IconActionCode.Layout(gtx, th.Fg)
								}),
								layout.Rigid(layout.Spacer{Width: 4}.Layout),
								layout.Flexed(1, func(gtx page.C) page.D {
									if status != nil && status.Stats != nil {
										current, unitCurrent := format(int64(status.Stats.CurrentConns), 1000)
										current = float64(int64(current*10)) / 10

										total, unitTotal := format(int64(status.Stats.TotalConns), 1000)
										total = float64(int64(total*10)) / 10
										return material.Body2(th, fmt.Sprintf("%s%s / %s%s",
											strconv.FormatFloat(current, 'f', -1, 64), strings.ToLower(unitCurrent),
											strconv.FormatFloat(total, 'f', -1, 64), strings.ToLower(unitTotal))).Layout(gtx)
									}
									return material.Body2(th, "N/A").Layout(gtx)
								}),
								layout.Rigid(func(gtx page.C) page.D {
									if status != nil && status.Stats != nil {
										rate := status.Stats.RequestRate
										rate = float64(int64(rate*100)) / 100
										return material.Body2(th, fmt.Sprintf("%s R/s", strconv.FormatFloat(rate, 'f', -1, 64))).Layout(gtx)
									}
									return material.Body2(th, "N/A").Layout(gtx)
								}),
							)
						}),
						layout.Rigid(layout.Spacer{Height: 4}.Layout),
						layout.Rigid(func(gtx page.C) page.D {
							return layout.Flex{
								Alignment: layout.Middle,
								Spacing:   layout.SpaceBetween,
							}.Layout(gtx,
								layout.Rigid(func(gtx page.C) page.D {
									return icons.IconNavExpandLess.Layout(gtx, th.Fg)
								}),
								layout.Rigid(layout.Spacer{Width: 4}.Layout),
								layout.Flexed(1, func(gtx page.C) page.D {
									if status != nil && status.Stats != nil {
										v, unit := format(int64(status.Stats.OutputBytes), 1024)
										v = float64(int64(v*100)) / 100
										return material.Body2(th, fmt.Sprintf("%s %sB", strconv.FormatFloat(v, 'f', -1, 64), unit)).Layout(gtx)
									}
									return material.Body2(th, "N/A").Layout(gtx)
								}),
								layout.Rigid(func(gtx page.C) page.D {
									if status != nil && status.Stats != nil {
										v, unit := format(int64(status.Stats.OutputRateBytes), 1024)
										v = float64(int64(v*100)) / 100
										return material.Body2(th, fmt.Sprintf("%s %sB/s", strconv.FormatFloat(v, 'f', -1, 64), unit)).Layout(gtx)
									}
									return material.Body2(th, "N/A").Layout(gtx)
								}),
							)
						}),
						layout.Rigid(layout.Spacer{Height: 4}.Layout),
						layout.Rigid(func(gtx page.C) page.D {
							return layout.Flex{
								Alignment: layout.Middle,
								Spacing:   layout.SpaceBetween,
							}.Layout(gtx,
								layout.Rigid(func(gtx page.C) page.D {
									return icons.IconNavExpandMore.Layout(gtx, th.Fg)
								}),
								layout.Rigid(layout.Spacer{Width: 4}.Layout),
								layout.Flexed(1, func(gtx page.C) page.D {
									if status != nil && status.Stats != nil {
										v, unit := format(int64(status.Stats.InputBytes), 1024)
										v = float64(int64(v*100)) / 100
										return material.Body2(th, fmt.Sprintf("%s %sB", strconv.FormatFloat(v, 'f', -1, 64), unit)).Layout(gtx)
									}
									return material.Body2(th, "N/A").Layout(gtx)
								}),
								layout.Rigid(func(gtx page.C) page.D {
									if status != nil && status.Stats != nil {
										v, unit := format(int64(status.Stats.InputRateBytes), 1024)
										v = float64(int64(v*100)) / 100
										return material.Body2(th, fmt.Sprintf("%s %sB/s", strconv.FormatFloat(v, 'f', -1, 64), unit)).Layout(gtx)
									}
									return material.Body2(th, "N/A").Layout(gtx)
								}),
							)
						}),
					)
				})
			})
		})
	})
}

var (
	units = []string{"", "K", "M", "G", "T", "P", "E"}
)

func format(n int64, scale int64) (v float64, unit string) {
	var remain float64
	for i := range units {
		unit = units[i]

		r := n % scale
		if n = n / scale; n == 0 {
			v = float64(r) + remain/math.Pow(float64(scale), float64(i))
			return
		}
		remain += float64(r) * math.Pow(float64(scale), float64(i))
	}
	return
}

var (
	dunits = []string{"s", "m", "h"}
)

func formatDuration(d time.Duration) (v int64, unit string) {
	if d.Hours() >= 24 {
		v = int64(d.Hours() / 24)
		unit = "d"
		return
	}

	var scale int64 = 60
	n := int64(d.Seconds())
	for i := range dunits {
		v = n % scale
		unit = dunits[i]

		if n = n / scale; n == 0 {
			return
		}
	}
	return
}
