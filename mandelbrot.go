package main

import (
	"fmt"
	"image"
	"image/color"
	"runtime"
	"time"
)

func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", what, time.Since(start))
	}
}

var palette [][3]uint8

func fillPalette() {
	for i := 0; i < 25; i++ {
		palette = append(palette, [3]uint8{uint8(i) * 10, uint8(i) * 8, 50 + uint8(i)*8})
	}
	for i := 25; i > 5; i-- {
		palette = append(palette, [3]uint8{50 + uint8(i)*8, 150 + uint8(i)*2, uint8(i) * 10})
	}
	for i := 10; i > 2; i-- {
		palette = append(palette, [3]uint8{0, uint8(i) * 15, 48})
	}
}

func computePixelColor(real float64, imag float64, iterNumber int) int {
	var zReal float64 = 0.0
	var zImag float64 = 0.0
	for i := 0; i < iterNumber; i++ {
		zReal, zImag = zReal*zReal-zImag*zImag+real,
			2*zReal*zImag+imag
		if (zReal*zReal + zImag*zImag) >= 4 {
			return i
		}
	}
	return -1
}

func computeRowColors(y int, size int, minX float64, minY float64, resolution float64, iterNumber int, result *image.RGBA) {

	//hsvColor := [3]float64{0.0, 0.8, 0.5}
	rgbaColor := color.RGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 0xff,
	}

	zImag := minY + float64(y)*resolution
	for x := 0; x < size; x++ {
		zReal := minX + float64(x)*resolution
		colorIndex := computePixelColor(zReal, zImag, iterNumber)

		if colorIndex >= 0 {
			//hsvColor[0] = 1 - float64(colorIndex % 100) / 100.0
			//rgbColor := hsv_to_rgb(hsvColor[0], hsvColor[1], hsvColor[2])
			rgbColor := palette[colorIndex%len(palette)]
			rgbaColor.R = rgbColor[0]
			rgbaColor.G = rgbColor[1]
			rgbaColor.B = rgbColor[2]
		} else {
			rgbaColor.R = 0
			rgbaColor.G = 0
			rgbaColor.B = 0
		}
		result.Set(x, y, &rgbaColor)
	}
}

func mandelbrot(width int, height int, min_x float64, min_y float64, resolution float64, iter_number int, image_buffer *image.RGBA) {
	for y := 0; y < height; y++ {
		computeRowColors(y, width, min_x, min_y, resolution, iter_number, image_buffer)
	}
}

func computePixelColorAsync(yFrom int, yTo int, size int, minX float64, minY float64, pixelSize float64, iterNumber int, result *image.RGBA, c chan int) {

	for y := yFrom; y < yTo; y++ {
		computeRowColors(y, size, minX, minY, pixelSize, iterNumber, result)
	}

	c <- 1
}

func mandelbrotAsync(width int, height int, min_x float64, min_y float64, resolution float64, iter_number int, image_buffer *image.RGBA) {
	numCpu := runtime.NumCPU() - 1
	c := make(chan int, numCpu)
	step := height / numCpu
	for i := 0; i < numCpu; i++ {
		yFrom, yTo := i*step, (i+1)*step
		go computePixelColorAsync(yFrom, yTo, width, min_x, min_y, resolution, iter_number, image_buffer, c)
	}

	for i := 0; i < numCpu; i++ {
		<-c
	}
}

// func mandelbrotImage(width uint32, height uint32, min_x float64, max_x float64, min_y float64, max_y float64, iter_number uint32) *bytes.Buffer {
// 	defer elapsed("mandelbrotImage")()

// 	var imageData *image.RGBA
// 	if isAsync {
// 		imageData = mandelbrotAsync(size, minX, minY, maxX, iterNumber)
// 	} else {
// 		imageData = mandelbrot(size, minX, minY, maxX, iterNumber)
// 	}

// 	imageBuffer := new(bytes.Buffer)
// 	err := png.Encode(imageBuffer, imageData)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	return imageBuffer
// }
