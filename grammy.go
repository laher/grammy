package main

import (
	"image"
	"image/color"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
)

type item struct {
	w, h float64
}

type columns []item

func main() {
	space := float64(10)
	items := []item{
		{100, 50},
		{100, 50},
		{100, 50},
	}
	cols := columns{}
	rows := []item{}
	var h, w float64
	h += space * 2
	w += space * 2
	cols = append(cols, item{space, 0})
	rows = append(rows, item{space, 0})
	maxH := h
	for _, i := range items {
		cols = append(cols, item{i.w, 0})
		cols = append(cols, item{space, 0})
		w += i.w + space
		if maxH < i.h {
			maxH = i.h
			h += i.h + space
		}
	}
	rows = append(rows, item{0, maxH})

	// Initialize the graphic context on an RGBA image
	dest := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	gc := draw2dimg.NewGraphicContext(dest)

	x := float64(space)
	y := float64(space)
	for _, i := range items {
		// Set some properties
		//#DFDCD3
		gc.SetFillColor(color.RGBA{0xff, 0xff, 0x99, 0xff})
		gc.SetStrokeColor(color.RGBA{0x44, 0x44, 0x44, 0xff})
		gc.SetLineWidth(1)

		gc.MoveTo(x, y) // should always be called first for a new path
		// Draw a closed shape
		//gc.LineTo(i.w+x, i.h+y)
		//gc.QuadCurveTo(i.w+x, y, x, y)
		//gc.Close()
		draw2dkit.Rectangle(gc, x, y, i.w+x, i.h+y)
		gc.FillStroke()

		//y += i.h + 10
		x += i.w + space
	}
	// Save to file
	draw2dimg.SaveToPngFile("hello.png", dest)
}
