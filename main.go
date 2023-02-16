package main

import (
	"bytes"
	"image"
	"image/png"

	"github.com/gofiber/fiber/v2"
)

type Params struct {
	MinX       float64 `query:"min_x"`
	MinY       float64 `query:"min_y"`
	MaxX       float64 `query:"max_x"`
	IterNumber int     `query:"iter_number"`
	Width      int     `query:"width"`
	Height     int     `query:"height"`
	IsAsync    bool    `query:"is_async"`
}

type ParamsCenter struct {
	X          float64 `query:"x"`
	Y          float64 `query:"y"`
	Resolution float64 `query:"resolution"`
	IterNumber int     `query:"iter_number"`
	Width      int     `query:"width"`
	Height     int     `query:"height"`
	IsAsync    bool    `query:"is_async"`
}

func main() {
	app := fiber.New()

	app.Get("/", Mandelbrot)
	app.Get("/fom-center", MandelbrotFromCenter)

	app.Listen(":3000")
}

func Mandelbrot(c *fiber.Ctx) error {
	println("Mandelbrot started")
	defer elapsed("Mandelbrot")()

	p := new(Params)

	if err := c.QueryParser(p); err != nil {
		println(err)
	}

	fillPalette()

	imageBuffer := image.NewRGBA(image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{int(p.Width), int(p.Height)},
	})
	if p.IsAsync {
		mandelbrotAsync(p.Width, p.Height, p.MinX, p.MinY, (p.MaxX-p.MinX)/float64(p.Width), p.IterNumber, imageBuffer)
	} else {
		mandelbrot(p.Width, p.Height, p.MinX, p.MinY, (p.MaxX-p.MinX)/float64(p.Width), p.IterNumber, imageBuffer)
	}

	buffer := new(bytes.Buffer)
	png.Encode(buffer, imageBuffer)

	c.Response().Header.SetContentType("image/png")
	c.SendStream(buffer, buffer.Len())
	// c.Response.SetBodyStream(imageBuffer, imageBuffer.Len())
	// c.Response.Header.SetContentType("image/png")

	return nil
}

func MandelbrotFromCenter(c *fiber.Ctx) error {
	println("MandelbrotFromCenter started")
	defer elapsed("MandelbrotFromCenter")()

	p := new(ParamsCenter)

	if err := c.QueryParser(p); err != nil {
		println(err)
	}

	minX := p.X - float64(p.Width/2)*p.Resolution
	if p.Width%2 != 0 {
		minX += p.Resolution / 2.0
	}
	minY := p.Y - float64(p.Height/2)*p.Resolution
	if p.Height%2 != 0 {
		minY += p.Resolution / 2.0
	}

	fillPalette()

	imageBuffer := image.NewRGBA(image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{int(p.Width), int(p.Height)},
	})
	if p.IsAsync {
		mandelbrotAsync(p.Width, p.Height, minX, minY, p.Resolution, p.IterNumber, imageBuffer)
	} else {
		mandelbrot(p.Width, p.Height, minX, minY, p.Resolution, p.IterNumber, imageBuffer)
	}

	buffer := new(bytes.Buffer)
	png.Encode(buffer, imageBuffer)

	c.Response().Header.SetContentType("image/png")
	return c.SendStream(buffer, buffer.Len())
}
