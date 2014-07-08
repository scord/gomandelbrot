package main

import (
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"runtime"
)

type pixel struct {
	x, y int
}

func mandelbrot(w, h, i int, z float32, seed int64) *image.RGBA {

	work := make(chan pixel)

	colors := make([]color.RGBA, i)

	rand.Seed(seed)

	for index := range colors {
		colors[index] = randomColor()
	}

	zoom := 1 / z

	m := image.NewRGBA(image.Rect(0, 0, w, h))

	for t := 0; t < 256; t++ {
		go func() {
			for p := range work {
				setColor(m, colors, p.x, p.y, i, zoom)
			}
		}()
	}

	b := m.Bounds()

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			work <- pixel{x, y}
		}
	}

	close(work)

	return m

}

func randomColor() color.RGBA {

	r := uint8(rand.Float32() * 255)
	g := uint8(rand.Float32() * 255)
	b := uint8(rand.Float32() * 255)

	return color.RGBA{r, g, b, 255}

}

func setColor(m *image.RGBA, colors []color.RGBA, px, py, maxi int, zoom float32) {

	x0 := zoom * (3.5*float32(px)/float32(m.Bounds().Size().X) - 2.5)
	y0 := zoom * (2*float32(py)/float32(m.Bounds().Size().Y) - 1.0)
	x := float32(0)
	y := float32(0)

	i := 0

	for x*x+y*y < 2*2 && i < maxi {

		xtemp := x*x - y*y + x0

		y = 2*x*y + y0
		x = xtemp

		i++
	}

	m.Set(px, py, colors[i-1])

	return
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	m := mandelbrot(1000, 1000, 100, 1.0, 1)

	w, _ := os.Create("mandelbrot.png")
	defer w.Close()
	png.Encode(w, m)
}
