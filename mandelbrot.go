package main

import (
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"runtime"
	"sync"
)

func mandelbrot(w, h, i int, z float32, seed int64) *image.RGBA {

	var wg sync.WaitGroup

	colors := make([]color.RGBA, i)

	rand.Seed(seed)

	for index := range colors {
		colors[index] = randomColor()
	}

	zoom := 1 / z

	m := image.NewRGBA(image.Rect(0, 0, w, h))

	b := m.Bounds()
	wg.Add(b.Size().Y)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		go func(y int) {
			for x := b.Min.X; x < b.Max.X; x++ {
				setColor(m, colors, x, y, i, zoom)
			}
			wg.Done()
		}(y)
	}

	wg.Wait()

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
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	m := mandelbrot(3000, 3000, 50, 1.0, 2)

	w, _ := os.Create("mandelbrot.png")
	defer w.Close()
	png.Encode(w, m)
}
