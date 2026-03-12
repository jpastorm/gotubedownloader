//go:build ignore
// +build ignore

// Generates menu bar template icons from new-menu-app-icon.png using only stdlib.

package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
)

func main() {
	f, _ := os.Open("new-menu-app-icon.png")
	src, _ := png.Decode(f)
	f.Close()

	save(src, 18, "menu_icon_18x18.png")
	save(src, 36, "menu_icon_36x36@2x.png")
}

// Nearest-neighbor + area-average hybrid downscale (stdlib only)
func downscale(src image.Image, size int) *image.NRGBA {
	sb := src.Bounds()
	sw, sh := sb.Dx(), sb.Dy()
	out := image.NewNRGBA(image.Rect(0, 0, size, size))

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			// Map output pixel to source area
			sx0 := float64(x) * float64(sw) / float64(size)
			sy0 := float64(y) * float64(sh) / float64(size)
			sx1 := float64(x+1) * float64(sw) / float64(size)
			sy1 := float64(y+1) * float64(sh) / float64(size)

			var rr, gg, bb, aa float64
			var count float64

			for sy := int(sy0); sy < int(math.Ceil(sy1)); sy++ {
				for sx := int(sx0); sx < int(math.Ceil(sx1)); sx++ {
					if sx >= sw || sy >= sh {
						continue
					}
					r, g, b, a := src.At(sx+sb.Min.X, sy+sb.Min.Y).RGBA()
					rr += float64(r)
					gg += float64(g)
					bb += float64(b)
					aa += float64(a)
					count++
				}
			}

			if count > 0 {
				out.SetNRGBA(x, y, color.NRGBA{
					R: uint8(rr / count / 257),
					G: uint8(gg / count / 257),
					B: uint8(bb / count / 257),
					A: uint8(aa / count / 257),
				})
			}
		}
	}
	return out
}

func save(src image.Image, size int, filename string) {
	scaled := downscale(src, size)

	// Convert to template: black + alpha derived from darkness * original alpha
	out := image.NewNRGBA(image.Rect(0, 0, size, size))
	draw.Draw(out, out.Bounds(), image.Transparent, image.Point{}, draw.Src)

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			c := scaled.NRGBAAt(x, y)
			if c.A == 0 {
				continue
			}
			// Perceived luminance
			lum := 0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B)
			darkness := 1.0 - (lum / 255.0)
			// Combine with original alpha
			alpha := darkness * (float64(c.A) / 255.0)
			// Boost contrast
			alpha = math.Pow(alpha, 0.6)
			a := uint8(math.Min(255, alpha*255))
			out.SetNRGBA(x, y, color.NRGBA{0, 0, 0, a})
		}
	}

	w, _ := os.Create(filename)
	defer w.Close()
	png.Encode(w, out)
}
