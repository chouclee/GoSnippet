package main

import (
	"bufio"
	"code.google.com/p/draw2d/draw2d"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

type Canvas struct {
	gc     *draw2d.ImageGraphicContext
	img    image.Image
	width  int
	height int
}

func (c *Canvas) MoveTo(x, y float64) {
	c.gc.MoveTo(x, y)
}

func (c *Canvas) LineTo(x, y float64) {
	c.gc.LineTo(x, y)
}

func (c *Canvas) SetStrokeColor(col color.Color) {
	c.gc.SetStrokeColor(col)
}

func (c *Canvas) SetFillColor(col color.Color) {
	c.gc.SetFillColor(col)
}

func (c *Canvas) SetLineWidth(w float64) {
	c.gc.SetLineWidth(w)
}

func (c *Canvas) Stroke() {
	c.gc.Stroke()
}

func (c *Canvas) FillStroke() {
	c.gc.FillStroke()
}

func (c *Canvas) Fill() {
	c.gc.Fill()
}

func (c *Canvas) Clear() {
	c.gc.Clear()
}

func (c *Canvas) ClearRect(x1, y1, x2, y2 int) {
	c.gc.ClearRect(x1, y1, x2, y2)
}

func (c *Canvas) SaveToPNG(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	b := bufio.NewWriter(f)
	err = png.Encode(b, c.img)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Wrote %s OK.\n", filename)
}

func (c *Canvas) Width() int {
	return c.width
}

func (c *Canvas) Height() int {
	return c.height
}

func CreateNewCanvas(w, h int) Canvas {
	i := image.NewRGBA(image.Rect(0, 0, w, h))
	gc := draw2d.NewGraphicContext(i)

	gc.SetStrokeColor(image.Black)
	gc.SetFillColor(image.White)
	// fill the background
	gc.Clear()
	gc.SetFillColor(image.Black)

	return Canvas{gc, i, w, h}
}

func MakeColor(r, g, b uint8) color.Color {
	return &color.RGBA{r, g, b, 255}
}
