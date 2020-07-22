package main

import (
	stdimage "image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
	"strconv"
)

func mkImg(r result) {
	n := r.needle
	h := r.haystack

	of, err := os.Open(dir + h.name)
	if err != nil {
		log.Fatalln("Failed to open file: ", err.Error())
	}
	defer of.Close()

	nf, err := os.Create(n.name + "_IN_" + h.name + "_DIFFERENCE_" + strconv.Itoa(r.avgDiff) + "_Y_X_" + strconv.Itoa(r.hIdx/(h.width)) + "_" + strconv.Itoa(r.hIdx%(h.width)) + " .jpg")
	if err != nil {
		log.Fatalln("Failed to create a file: ", err.Error())
	}
	defer nf.Close()

	dof, _ := jpeg.Decode(of)
	m := stdimage.NewRGBA(dof.Bounds())

	hIdx := r.hIdx

	for i := 0; i < n.height*n.width; i++ {
		hX := i % (h.width)
		hY := i / (h.width)
		r := uint8(h.pixels[i].r)
		g := uint8(h.pixels[i].g)
		b := uint8(h.pixels[i].b)
		a := uint8(h.pixels[i].a)
		m.Set(hX, hY, color.RGBA{
			R: r,
			G: g,
			B: b,
			A: a,
		})
	}

	for i := 0; i < n.height*n.width; i++ {
		diff := pixelDiff(n.pixels[i], h.pixels[i])

		var d uint8
		a := uint8(math.Floor(255 * .5))

		if diff < threshold {
			d = uint8(255 * (1 - (diff - threshold)))
		} else {
			d = 0
		}

		nX := i % (n.width)
		nY := i / (n.width)
		if nX < 5 || ((n.width - nX) < 5) {
			d = 160
		}
		if nY == 0 || (nY == (n.height - 1)) {
			d = 160
		}

		hX := hIdx % (h.width)
		hY := hIdx / (h.width)
		m.Set(hX, hY, color.RGBA{
			R: d,
			G: d,
			B: d,
			A: a,
		})

		if ((i + 1) % n.width) == 0 {
			hIdx += (h.width - n.width)
		}
		hIdx++
	}

	jpeg.Encode(nf, m, nil)
}
