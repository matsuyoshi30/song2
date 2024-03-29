package song2

// https://www.youtube.com/watch?v=SSbBvKaM6sk

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"runtime"
	"sync"
)

func GaussianBlur(src image.Image, r float64) *image.RGBA {
	clone := CloneToRGBA(src)
	dst := CloneToRGBA(src)

	bxs := BoxesForGauss(r, 3)

	for _, b := range bxs {
		boxBlur(clone, dst, (b-1)/2)
	}

	return dst
}

type Direction int

const (
	dirX Direction = iota
	dirY
)

func boxBlur(src, dst *image.RGBA, r int) {
	height := src.Bounds().Max.Y - src.Bounds().Min.Y
	width := src.Bounds().Max.X - src.Bounds().Min.X

	boxBlurParallel(dirX, height, dst, src, r)
	boxBlurParallel(dirY, width, src, dst, r)
}

func boxBlurParallel(d Direction, length int, src, dst *image.RGBA, r int) {
	procs := runtime.NumCPU()
	ps := length / procs

	var wg sync.WaitGroup
	for length > 0 {
		start := length - ps
		if start < 0 {
			start = 0
		}
		end := length

		length -= ps
		wg.Add(1)
		go func() {
			defer wg.Done()
			switch d {
			case dirX:
				BoxBlurHorizontal(src, dst, src.Bounds().Min.Y+start, src.Bounds().Min.Y+end, r)
			case dirY:
				BoxBlurTotal(src, dst, src.Bounds().Min.X+start, src.Bounds().Min.X+end, r)
			}
		}()
	}

	wg.Wait()
}

func BoxBlurHorizontal(src, dst *image.RGBA, start, end, r int) {
	fr := float64(r)
	iarr := 1.0 / (fr + fr + 1.0)

	for i := start; i < end; i++ {
		ti := src.Bounds().Min.X
		li := ti
		ri := ti + r

		fvpos := src.PixOffset(ti, i)
		lvpos := src.PixOffset(src.Bounds().Max.X-1, i)

		fvr := int(src.Pix[fvpos+0])
		fvg := int(src.Pix[fvpos+1])
		fvb := int(src.Pix[fvpos+2])
		fva := int(src.Pix[fvpos+3])

		val_r := fvr * (r + 1)
		val_g := fvg * (r + 1)
		val_b := fvb * (r + 1)
		val_a := fva * (r + 1)

		for j := 0; j < r; j++ {
			pos := src.PixOffset(ti+j, i)
			val_r += int(src.Pix[pos+0])
			val_g += int(src.Pix[pos+1])
			val_b += int(src.Pix[pos+2])
			val_a += int(src.Pix[pos+3])
		}

		for j := 0; j <= r; j++ {
			pos := src.PixOffset(ri, i)
			ri++

			val_r += int(src.Pix[pos+0]) - fvr
			val_g += int(src.Pix[pos+1]) - fvg
			val_b += int(src.Pix[pos+2]) - fvb
			val_a += int(src.Pix[pos+3]) - fva

			_r := uint8(math.Round(float64(val_r) * iarr))
			_g := uint8(math.Round(float64(val_g) * iarr))
			_b := uint8(math.Round(float64(val_b) * iarr))
			_a := uint8(math.Round(float64(val_a) * iarr))

			dst.SetRGBA(ti, i, color.RGBA{_r, _g, _b, _a})
			ti++
		}

		for j := r + 1; j < src.Bounds().Max.X-r; j++ {
			ripos := src.PixOffset(ri, i)
			ri++

			lipos := src.PixOffset(li, i)
			li++

			val_r += int(src.Pix[ripos+0]) - int(src.Pix[lipos+0])
			val_g += int(src.Pix[ripos+1]) - int(src.Pix[lipos+1])
			val_b += int(src.Pix[ripos+2]) - int(src.Pix[lipos+2])
			val_a += int(src.Pix[ripos+3]) - int(src.Pix[lipos+3])

			_r := uint8(math.Round(float64(val_r) * iarr))
			_g := uint8(math.Round(float64(val_g) * iarr))
			_b := uint8(math.Round(float64(val_b) * iarr))
			_a := uint8(math.Round(float64(val_a) * iarr))

			dst.SetRGBA(ti, i, color.RGBA{_r, _g, _b, _a})
			ti++
		}

		for j := src.Bounds().Max.X - r; j < src.Bounds().Max.X; j++ {
			pos := src.PixOffset(li, i)
			li++

			val_r += int(src.Pix[lvpos+0]) - int(src.Pix[pos+0])
			val_g += int(src.Pix[lvpos+1]) - int(src.Pix[pos+1])
			val_b += int(src.Pix[lvpos+2]) - int(src.Pix[pos+2])
			val_a += int(src.Pix[lvpos+3]) - int(src.Pix[pos+3])

			_r := uint8(math.Round(float64(val_r) * iarr))
			_g := uint8(math.Round(float64(val_g) * iarr))
			_b := uint8(math.Round(float64(val_b) * iarr))
			_a := uint8(math.Round(float64(val_a) * iarr))

			dst.SetRGBA(ti, i, color.RGBA{_r, _g, _b, _a})
			ti++
		}
	}
}

func BoxBlurTotal(src, dst *image.RGBA, start, end, r int) {
	fr := float64(r)
	iarr := 1.0 / (fr + fr + 1.0)

	for i := start; i < end; i++ {
		ti := src.Bounds().Min.Y
		li := ti
		ri := ti + r

		fvpos := src.PixOffset(i, ti)
		lvpos := src.PixOffset(i, src.Bounds().Max.Y-1)

		fvr := int(src.Pix[fvpos+0])
		fvg := int(src.Pix[fvpos+1])
		fvb := int(src.Pix[fvpos+2])
		fva := int(src.Pix[fvpos+3])

		val_r := fvr * (r + 1)
		val_g := fvg * (r + 1)
		val_b := fvb * (r + 1)
		val_a := fva * (r + 1)

		for j := 0; j < r; j++ {
			pos := src.PixOffset(i, ti+j)
			val_r += int(src.Pix[pos+0])
			val_g += int(src.Pix[pos+1])
			val_b += int(src.Pix[pos+2])
			val_a += int(src.Pix[pos+3])
		}

		for j := 0; j <= r; j++ {
			pos := src.PixOffset(i, ri)
			ri++

			val_r += int(src.Pix[pos+0]) - fvr
			val_g += int(src.Pix[pos+1]) - fvg
			val_b += int(src.Pix[pos+2]) - fvb
			val_a += int(src.Pix[pos+3]) - fva

			_r := uint8(math.Round(float64(val_r) * iarr))
			_g := uint8(math.Round(float64(val_g) * iarr))
			_b := uint8(math.Round(float64(val_b) * iarr))
			_a := uint8(math.Round(float64(val_a) * iarr))

			dst.SetRGBA(i, ti, color.RGBA{_r, _g, _b, _a})
			ti++
		}

		for j := r + 1; j < src.Bounds().Max.Y-r; j++ {
			ripos := src.PixOffset(i, ri)
			ri++

			lipos := src.PixOffset(i, li)
			li++

			val_r += int(src.Pix[ripos+0]) - int(src.Pix[lipos+0])
			val_g += int(src.Pix[ripos+1]) - int(src.Pix[lipos+1])
			val_b += int(src.Pix[ripos+2]) - int(src.Pix[lipos+2])
			val_a += int(src.Pix[ripos+3]) - int(src.Pix[lipos+3])

			_r := uint8(math.Round(float64(val_r) * iarr))
			_g := uint8(math.Round(float64(val_g) * iarr))
			_b := uint8(math.Round(float64(val_b) * iarr))
			_a := uint8(math.Round(float64(val_a) * iarr))

			dst.SetRGBA(i, ti, color.RGBA{_r, _g, _b, _a})
			ti++
		}

		for j := src.Bounds().Max.Y - r; j < src.Bounds().Max.Y; j++ {
			pos := src.PixOffset(i, li)
			li++

			val_r += int(src.Pix[lvpos+0]) - int(src.Pix[pos+0])
			val_g += int(src.Pix[lvpos+1]) - int(src.Pix[pos+1])
			val_b += int(src.Pix[lvpos+2]) - int(src.Pix[pos+2])
			val_a += int(src.Pix[lvpos+3]) - int(src.Pix[pos+3])

			_r := uint8(math.Round(float64(val_r) * iarr))
			_g := uint8(math.Round(float64(val_g) * iarr))
			_b := uint8(math.Round(float64(val_b) * iarr))
			_a := uint8(math.Round(float64(val_a) * iarr))

			dst.SetRGBA(i, ti, color.RGBA{_r, _g, _b, _a})
			ti++
		}
	}
}

// BoxesForGauss
func BoxesForGauss(sigma float64, n int) []int { // standard deviation, number of boxes
	nf := float64(n)

	wIdeal := math.Sqrt(12.0*sigma*sigma/nf + 1.0)
	wl := int(math.Floor(wIdeal))
	if wl%2 == 0 {
		wl--
	}
	wu := wl + 2

	mIdeal := (12.0*sigma*sigma - float64(n*wl*wl+4*n*wl+3*n)) / float64(-4*wl-4)
	m := math.Round(mIdeal)

	sizes := make([]int, n)
	for i := 0; i < n; i++ {
		if float64(i) < m {
			sizes[i] = wl
		} else {
			sizes[i] = wu
		}
	}

	return sizes
}

// CloneToRGBA
func CloneToRGBA(src image.Image) *image.RGBA {
	b := src.Bounds()
	dst := image.NewRGBA(b)
	draw.Draw(dst, b, src, b.Min, draw.Src)
	return dst
}
