package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"sync"
)

func renderPixel(img *image.RGBA, ximg, yimg int, xplane, yplane float64) {
	c := complex(xplane, yplane)
	z := complex(0, 0)
	i, max := 0, 255

	// Save time by ensuring that z doesn't exceed a radius
	// of 2 from the origin, using the Pythagorean theorem
	for real(z) * real(z) + imag(z) * imag(z) <= 4 && i < max {
		z = z * z + c
		i++
	}

	var rg, b uint8

	if i < max {
		// Circular function for an extra highlight
		x := float64(i) / float64(max)
		f := math.Sqrt(2 * x - x * x)

		rg = uint8(255 * f)
		b = 111
	}

	img.SetRGBA(ximg, yimg, color.RGBA {
		rg, rg, b, 255,
	})
}

func render(img *image.RGBA, minX, maxX, minY, maxY float64) {
	pixelW := (maxX - minX) / float64(img.Bounds().Dx())
	pixelH := (maxY - minY) / float64(img.Bounds().Dy())

	ximg := 0
	xplane := minX

	var wg sync.WaitGroup
	wg.Add(img.Bounds().Dx())

	for ximg < img.Bounds().Dx() {
		go func(ximg int, xplane float64) {
			yimg := 0
			yplane := minY

			for yimg < img.Bounds().Dy() {
				renderPixel(img, ximg, yimg, xplane, yplane)

				yimg++
				yplane += pixelH
			}

			wg.Done()
		}(ximg, xplane)

		ximg++
		xplane += pixelW
	}

	wg.Wait()
}

func encode(img *image.RGBA, name string) {
	f, err := os.Create(name)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
}

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 1680, 1050))

	render(img, -2.5, 1.5, -1.25, 1.25)
	encode(img, "image.png")
}
