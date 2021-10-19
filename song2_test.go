package song2_test

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/anthonynsimon/bild/blur"
	"github.com/matsuyoshi30/song2"
)

var (
	img image.Image
	r   float64
)

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	file, err := os.Open(filepath.Join(pwd, "assets", "sample.png"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		file.Close()
	}()

	img, _, err = image.Decode(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	r = 5.0
}

func BenchmarkGaussianBlur1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		SimpleGaussianBlur(img, r)
	}
}

func BenchmarkGaussianBlur2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GaussianBlurUsingBox(img, r)
	}
}

func BenchmarkGaussianBlur3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GaussianBlurHT(img, r)
	}
}

func BenchmarkGaussianBlur(b *testing.B) {
	for n := 0; n < b.N; n++ {
		song2WithoutGoroutine(img, r)
	}
}

func BenchmarkStackblur(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Stackblur(img, uint32(r))
	}
}

func BenchmarkBildBlur(b *testing.B) {
	for n := 0; n < b.N; n++ {
		blur.Gaussian(img, r)
	}
}

func BenchmarkGaussianBlurUsingGoroutine(b *testing.B) {
	for n := 0; n < b.N; n++ {
		song2.GaussianBlur(img, r)
	}
}

// SimpleGaussianBlur implements super naive Gaussian Blur
func SimpleGaussianBlur(src image.Image, r float64) *image.RGBA {
	clone := song2.CloneToRGBA(src)

	dst := image.NewRGBA(src.Bounds())

	rs := int(math.Ceil(r * 2.57))
	for i := clone.Bounds().Min.Y; i < clone.Bounds().Max.Y; i++ {
		for j := clone.Bounds().Min.X; j < clone.Bounds().Max.X; j++ {
			val := color.RGBA{}
			pos := clone.PixOffset(int(j), int(i))
			val.A = clone.Pix[pos+3]

			var _r, _g, _b float64
			wSum := 0.0 // sum of weight
			for iy := i - rs; iy < i+rs+1; iy++ {
				for ix := j - rs; ix < j+rs+1; ix++ {
					x := math.Min(float64(clone.Bounds().Max.X-1), math.Max(0, float64(ix)))
					y := math.Min(float64(clone.Bounds().Max.Y-1), math.Max(0, float64(iy)))

					dsq := float64((ix-j)*(ix-j) + (iy-i)*(iy-i))
					w := math.Exp(-dsq/(2*r*r)) / (math.Pi * 2 * r * r)

					pos := clone.PixOffset(int(x), int(y))
					_r += float64(clone.Pix[pos+0]) * w
					_g += float64(clone.Pix[pos+1]) * w
					_b += float64(clone.Pix[pos+2]) * w
					wSum += w
				}
			}
			val.R = uint8(math.Round(_r / wSum))
			val.G = uint8(math.Round(_g / wSum))
			val.B = uint8(math.Round(_b / wSum))

			dst.SetRGBA(int(j), int(i), val)
		}
	}

	return dst
}

// GaussianBlurUsingBox implements the convolution of box blur
func GaussianBlurUsingBox(src image.Image, r float64) *image.RGBA {
	clone := song2.CloneToRGBA(src)
	bxs := song2.BoxesForGauss(r, 3)

	dst := image.NewRGBA(src.Bounds())
	boxBlur2(clone, dst, (bxs[0]-1)/2)
	boxBlur2(dst, clone, (bxs[1]-1)/2)
	boxBlur2(clone, dst, (bxs[2]-1)/2)

	return dst
}

func boxBlur2(src, dst *image.RGBA, r int) {
	for i := src.Bounds().Min.Y; i < src.Bounds().Max.Y; i++ {
		for j := src.Bounds().Min.X; j < src.Bounds().Max.X; j++ {
			val := color.RGBA{}
			pos := src.PixOffset(int(j), int(i))
			val.A = uint8(src.Pix[pos+3])

			var _r, _g, _b int
			for iy := i - r; iy < i+r+1; iy++ {
				for ix := j - r; ix < j+r+1; ix++ {
					x := math.Min(float64(src.Bounds().Max.X-1), math.Max(0, float64(ix)))
					y := math.Min(float64(src.Bounds().Max.Y-1), math.Max(0, float64(iy)))

					pos := src.PixOffset(int(x), int(y))
					_r += int(src.Pix[pos+0])
					_g += int(src.Pix[pos+1])
					_b += int(src.Pix[pos+2])
				}
			}
			val.R = uint8(_r / ((r + r + 1) * (r + r + 1)))
			val.G = uint8(_g / ((r + r + 1) * (r + r + 1)))
			val.B = uint8(_b / ((r + r + 1) * (r + r + 1)))

			dst.SetRGBA(int(j), int(i), val)
		}
	}
}

// GaussianBlurHT implements blur using horizontal and total
func GaussianBlurHT(src image.Image, r float64) *image.RGBA {
	clone := song2.CloneToRGBA(src)
	bxs := song2.BoxesForGauss(r, 3)

	dst := image.NewRGBA(src.Bounds())
	boxBlurHT(clone, dst, (bxs[0]-1)/2)
	boxBlurHT(dst, clone, (bxs[1]-1)/2)
	boxBlurHT(clone, dst, (bxs[2]-1)/2)

	return dst
}

func boxBlurHT(src, dst *image.RGBA, r int) {
	for y := src.Bounds().Min.Y; y < src.Bounds().Max.Y; y++ {
		for x := src.Bounds().Min.X; x < src.Bounds().Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			dst.SetRGBA(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}

	boxBlur_H(dst, src, r)
	boxBlur_T(src, dst, r)
}

func boxBlur_H(src, dst *image.RGBA, r int) {
	for i := src.Bounds().Min.Y; i < src.Bounds().Max.Y; i++ {
		for j := src.Bounds().Min.X; j < src.Bounds().Max.X; j++ {
			val := color.RGBA{}
			pos := src.PixOffset(int(j), int(i))
			val.A = src.Pix[pos+3]

			var _r, _g, _b int
			for ix := j - r; ix < j+r+1; ix++ {
				x := math.Min(float64(src.Bounds().Max.X-1), math.Max(0, float64(ix)))

				pos := src.PixOffset(int(x), int(i))
				_r += int(src.Pix[pos+0])
				_g += int(src.Pix[pos+1])
				_b += int(src.Pix[pos+2])
			}
			val.R = uint8(_r / (r + r + 1))
			val.G = uint8(_g / (r + r + 1))
			val.B = uint8(_b / (r + r + 1))

			dst.SetRGBA(int(j), int(i), val)
		}
	}
}

func boxBlur_T(src, dst *image.RGBA, r int) {
	for i := src.Bounds().Min.Y; i < src.Bounds().Max.Y; i++ {
		for j := src.Bounds().Min.X; j < src.Bounds().Max.X; j++ {
			val := color.RGBA{}
			pos := src.PixOffset(int(j), int(i))
			val.A = src.Pix[pos+3]

			var _r, _g, _b int
			for iy := i - r; iy < i+r+1; iy++ {
				y := math.Min(float64(src.Bounds().Max.Y-1), math.Max(0, float64(iy)))

				pos := src.PixOffset(int(j), int(y))
				_r += int(src.Pix[pos+0])
				_g += int(src.Pix[pos+1])
				_b += int(src.Pix[pos+2])
			}
			val.R = uint8(_r / (r + r + 1))
			val.G = uint8(_g / (r + r + 1))
			val.B = uint8(_b / (r + r + 1))

			dst.SetRGBA(int(j), int(i), val)
		}
	}
}

func song2WithoutGoroutine(src image.Image, r float64) *image.RGBA {
	clone := song2.CloneToRGBA(src)
	dst := song2.CloneToRGBA(src)

	bxs := song2.BoxesForGauss(r, 3)

	for _, b := range bxs {
		song2.BoxBlurHorizontal(dst, clone, dst.Bounds().Min.Y, dst.Bounds().Max.Y, (b-1)/2)
		song2.BoxBlurTotal(clone, dst, src.Bounds().Min.X, src.Bounds().Max.X, (b-1)/2)
	}

	return dst
}
