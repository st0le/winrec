package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kbinani/screenshot"
	"github.com/lxn/win"
)

var (
	output string
	hwnd   int64
)

func init() {
	flag.StringVar(&output, "o", "screen.gif", "name of output file")
	flag.Int64Var(&hwnd, "hwnd", 0, "hWnd of the window to record")
	flag.Parse()

}

func main() {

	var r win.RECT

	palette := append(palette.WebSafe, color.Transparent)
	outGif := &gif.GIF{}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\r- Recording Stopped. Writing to %s.\n", output)

		file, _ := os.Create(output)
		defer file.Close()
		gif.EncodeAll(file, outGif)

		os.Exit(0)
	}()

	for {
		win.GetWindowRect(win.HWND(hwnd), &r)
		bounds := image.Rect(int(r.Left), int(r.Top), int(r.Right), int(r.Bottom))
		imgbounds := image.Rect(0, 0, bounds.Dx(), bounds.Dy())
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			panic(err)
		}

		palettedImage := image.NewPaletted(imgbounds, palette)
		draw.Draw(palettedImage, imgbounds, img, image.ZP, draw.Src)
		outGif.Image = append(outGif.Image, palettedImage)
		outGif.Delay = append(outGif.Delay, 40)
		time.Sleep(time.Duration(40) * time.Millisecond)
	}

}
