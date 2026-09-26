package main

import (
	"Sapper/helper"
	"time"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
)

const Title = "Sapper"

func run() {
	var icons []pixel.Picture
	for _, f := range []string{"app/app_icon_16.png", "app/app_icon_32.png", "app/app_icon_256.png"} {
		if pic, err := helper.Load_system(f); err == nil {
			icons = append(icons, pic)
		}
	}

	cfg := opengl.WindowConfig{
		Title:  Title,
		Bounds: pixel.R(0, 0, 100, 100),
		Icon:   icons,
		VSync:  true,
	}
	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	app, err := NewApp(win)
	if err != nil {
		panic(err)
	}

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		app.Update(dt)
		app.Draw()

		win.Update()
	}
}

func main() {
	opengl.Run(run)
}
