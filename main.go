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
	"github.com/go-gost/gostctl/ui/page"
	"github.com/go-gost/gostctl/ui/theme"
	"github.com/go-gost/gostctl/ui/widget"
	_ "github.com/go-gost/gostctl/winres"
)

func main() {
	Init()

	go func() {
		err := run()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run() error {
	ui := ui.NewUI()

	go handleEvent(ui)

	w := ui.Window()
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

func handleEvent(ui *ui.UI) {
	for {
		select {
		case e := <-ui.Router().Event():
			switch e.ID {
			case page.EventThemeChanged:
				slog.Debug("theme changed", "event", e.ID)
				ui.Window().Option(app.StatusColor(theme.Current().Material.Bg))
			}

		case e := <-runner.Event():
			switch e.TaskID {
			case runner.TaskGetConfig:
				server := config.CurrentServer()
				if server == nil {
					break
				}

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
}

func Init() {
	config.Init()

	util.RestartGetConfigTask()
}
