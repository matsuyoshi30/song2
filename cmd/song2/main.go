package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/matsuyoshi30/song2"
)

var (
	output = flag.String("o", "blured.png", "Write output image to specific filepath")
	radius = flag.Float64("r", 3.0, "Radius")

	name = "song2"
)

const (
	exitCodeOK = iota
	exitCodeErr
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage:
  %s [FLAGS] [FILE]

FLAGS:
  -o  Write output image to specifig filepath [default: blured.png]
  -r  Radius [default: 3.0]

Author:
  matsuyoshi30 <sfbgwm30@gmail.com>
`, name)
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "USAGE: %s [FLAGS] [FILE]\n", name)
		return
	}

	os.Exit(run(args[0]))
}

func run(src string) int {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitCodeErr
	}

	file, err := os.Open(filepath.Join(pwd, src))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitCodeErr
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitCodeErr
	}

	blured := song2.GaussianBlur(img, *radius)

	out, err := os.Create(filepath.Join(pwd, *output))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitCodeErr
	}
	defer out.Close()

	if err := png.Encode(out, blured); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitCodeErr
	}

	return exitCodeOK
}
