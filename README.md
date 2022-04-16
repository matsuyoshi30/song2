# Song2

Fast (linear time) implementation of the Gaussian Blur algorithm in Go.

Original algorithm taken from http://blog.ivank.net/fastest-gaussian-blur.html, and use goroutine.


## Install & Usage

### Package

Download and compile from this repository.

```sh
go get -u github.com/matsuyoshi30/song2
```

And import as a package, call the API `song2.GaussianBlur(src, blurRadius)`.

```go
package main

import (
    "fmt"
    "image"
    "image/png"
    "os"

    "github.com/matsuyoshi30/song2"
)

func main() {
    file, err := os.Open("./input.png")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer file.Close()

    img, _, err := image.Decode(file)
    if err != nil {
        fmt.Println(err)
        return
    }

    blurred := song2.GaussianBlur(img, 3.0)

    out, err := os.Create("./output.png")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer out.Close()

    if err := png.Encode(out, blurred); err != nil {
        fmt.Println(err)
        return
    }
}
```

### CLI tool

Clone this repository, and `go install`.

```sh
git clone https://github.com/matsuyoshi30/song2
cd song2/cmd/song2 && go install
```

You can use `song2` as a cli tool.

```sh
Usage:
  song2 [FLAGS] [FILE]

FLAGS:
  -o  Write output image to specifig filepath [default: blurred.png]
  -r  Radius [default: 3.0]

Author:
  matsuyoshi30 <sfbgwm30@gmail.com>
```


## Example

`song2 -o assets/blurred.png assets/sample.png`

|  Original                      |  Blured                      |
| ------------------------------ | ---------------------------- |
| ![original](assets/sample.png) | ![blurred](assets/blurred.png) |

This image is from http://sipi.usc.edu/database/database.php?volume=misc&image=10#top


## Benchmark

I also implemented another algorithms of [this post](http://blog.ivank.net/fastest-gaussian-blur.html) in test package (AnotherAlgorithm1-3).

```sh
% go test -bench . -benchmem
goos: darwin
goarch: amd64
pkg: github.com/matsuyoshi30/song2
cpu: Intel(R) Core(TM) i5-5287U CPU @ 2.90GHz
BenchmarkGaussianBlur1-4                	       1	12656993105 ns/op	 2097280 B/op	       4 allocs/op
BenchmarkGaussianBlur2-4                	       1	1649545490 ns/op	 2097304 B/op	       5 allocs/op
BenchmarkGaussianBlur3-4                	       4	 263524278 ns/op	 2097310 B/op	       5 allocs/op
BenchmarkGaussianBlur-4                 	      18	  56368998 ns/op	 2097309 B/op	       5 allocs/op
BenchmarkStackblur-4                    	      40	  28339048 ns/op	 3146083 B/op	  524301 allocs/op
BenchmarkBildBlur-4                     	      24	  48113312 ns/op	 4245167 B/op	      22 allocs/op
BenchmarkGaussianBlurUsingGoroutine-4   	      39	  29797842 ns/op	 2097535 B/op	      11 allocs/op
```

reference [stackblur-go](https://github.com/esimov/stackblur-go) and [bild](https://github.com/anthonynsimon/bild).


## LICENSE

[MIT](./LICENSE)


## Author

[matsuyoshi30](https://twitter.com/matsuyoshi30)
