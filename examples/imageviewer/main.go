package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"os/user"

	"github.com/faiface/mainthread"
	"github.com/juanefec/gui"
	"github.com/juanefec/gui/win"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/gofont/goregular"
)

func run() {
	face, err := TTFToFace(goregular.TTF, 18)
	if err != nil {
		panic(err)
	}

	theme := &Theme{
		Face:       face,
		Background: colornames.White,
		Empty:      colornames.Darkgrey,
		Text:       colornames.Black,
		Highlight:  colornames.Blueviolet,
		ButtonUp:   colornames.Lightgrey,
		ButtonDown: colornames.Grey,
	}

	icon := "iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAAAXNSR0IArs4c6QAAAPxJREFUWIXtlzkSwjAMRc1ymCQlPSehZOg4DW0GrsNQMTTAbaAAJx5H1pe8NfC7LOg9yRkSG/PrmVXmvXz2sgJoyLb/nD7uxr5zCExgFiSJViAJphFgR5gzrsBQPRfo/niarm3EAlk79OEhmbm0mPaaBC4W4MbYtQ0riH4vEkBB68wFClDdoY4190IBqjvJ2C0cTSd6CajCvpRkaYICmjHb+2OeBfKfMFSMgyD44rSSC8RCONhtsx4PzldeQBMIA1EJpMJUAiVgUMCFloBBgVpQN1neBSkp9VFK5rI/1BGgQN9MtgHJAgGYeL+hEkiFiQU0I0yNW3CybfqnRt6Qs2CdXlIOlAAAAABJRU5ErkJggg=="
	bicons, err := base64.StdEncoding.DecodeString(icon)
	if err != nil {
		panic(err)
	}
	ic, err := png.Decode(bytes.NewReader(bicons))
	if err != nil {
		panic(err)
	}

	w, err := win.New(win.Title("Image Viewer"), win.Size(900, 600), win.Resizable(), win.Icon([]image.Image{ic}))
	if err != nil {
		panic(err)
	}

	mux, env := gui.NewMux(w)

	cd := make(chan string)
	view := make(chan string)

	go Browser(FixedBottom(FixedLeft(mux.MakeEnv(), 300), 30), theme, ".", cd, view)
	go Viewer(FixedRight(mux.MakeEnv(), 300), theme, view)

	go Button(EvenHorizontal(FixedTop(FixedLeft(mux.MakeEnv(), 300), 30), 0, 1, 3), theme, "Dir Up", func() {
		cd <- ".."
	})
	go Button(EvenHorizontal(FixedTop(FixedLeft(mux.MakeEnv(), 300), 30), 1, 2, 3), theme, "Refresh", func() {
		cd <- "."
	})
	go Button(EvenHorizontal(FixedTop(FixedLeft(mux.MakeEnv(), 300), 30), 2, 3, 3), theme, "Home", func() {
		user, err := user.Current()
		if err != nil {
			return
		}
		cd <- user.HomeDir
	})

	for e := range env.Events() {
		switch e.(type) {
		case win.WiClose:
			close(env.Draw())
		}
	}
}

func main() {
	mainthread.Run(run)
}
