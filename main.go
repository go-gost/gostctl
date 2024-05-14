package main

import (
	"fmt"
	"log"
	"log/slog"
	_ "net"
	"os"
	"time"

	"gioui.org/app"
	_ "gioui.org/app/permission/storage"
	"gioui.org/op"
	"github.com/go-gost/gostctl/api/runner"
	"github.com/go-gost/gostctl/api/util"
	"github.com/go-gost/gostctl/config"
	"github.com/go-gost/gostctl/ui"
	"github.com/go-gost/gostctl/ui/widget"
	_ "github.com/go-gost/gostctl/winres"
)

func main() {
	Init()

	go func() {
		var w app.Window
		w.Option(app.Title("GOST"))
		// w.Option(app.MinSize(800, 600))
		err := run(&w)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	ui := ui.NewUI(w)

	go handleTaskEvent(ui)

	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			ui.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func handleTaskEvent(ui *ui.UI) {
	for e := range runner.Event() {
		switch e.TaskID {
		case runner.TaskGetConfig:
			cfg := config.Get()
			if cfg.CurrentServer >= 0 && cfg.CurrentServer < len(cfg.Servers) {
				server := cfg.Servers[cfg.CurrentServer]
				if e.Err != nil {
					slog.Error(fmt.Sprintf("task: %s", e.Err), "task", e.TaskID)
					server.SetState(config.ServerError)
					server.AddEvent(config.ServerEvent{
						Time: time.Now(),
						Msg:  e.Err.Error(),
					})
				} else {
					server.SetState(config.ServerReady)
				}
				ui.Window().Invalidate()
			}

		default:
			if e.Err != nil {
				slog.Error(fmt.Sprintf("task: %s", e.Err), "task", e.TaskID)
				ui.Router().Notify(widget.Message{
					Type:    widget.Error,
					Content: e.Err.Error(),
				})
			}
		}
	}
}

func Init() {
	config.Init()

	util.RestartGetConfigTask()
}
