package main

import (
	"math"
	"sync"
)

func compare(images []*image) chan result {
	var wg sync.WaitGroup
	ch := make(chan result)

	for _, needle := range images {
		for _, haystack := range images {
			if needle.width > haystack.width {
				continue
			}

			if needle.height > haystack.height {
				continue
			}

			if needle.name == haystack.name {
				continue
			}

			wg.Add(1)
			go func(n, h *image, ch chan result) {
				defer wg.Done()
				comparePixels(n, h, ch)
			}(needle, haystack, ch)
		}
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}

func comparePixels(n, h *image, ch chan result) {
	var diff float64
	var diffReverse float64
	var wg sync.WaitGroup

	for i := 0; i < (h.width * h.height); i++ {
		x := i % h.width
		y := i / h.width

		if h.height-y < n.height {
			break
		}

		diff = pixelDiff(n.pixels[0], h.pixels[i])
		diffReverse = pixelDiff(n.pixels[n.width-1], h.pixels[i])

		if (diff < threshold) && (x+n.width < h.width) {
			wg.Add(1)
			go func(n, h *image, i int, ch chan result) {
				defer wg.Done()
				compareSequence(n, h, i, ch)
			}(n, h, i, ch)
		}

		if (diffReverse < threshold) && (x+n.width < h.width) {
			wg.Add(1)
			go func(n, h *image, i int, ch chan result) {
				defer wg.Done()
				compareSequenceReverse(n, h, i, ch)
			}(n, h, i, ch)
		}
	}
	wg.Wait()
}

func pixelDiff(n, h pixel) float64 {
	var diff float64
	diff += math.Abs(float64(int(n.r) - int(h.r)))
	diff += math.Abs(float64(int(n.g) - int(h.g)))
	diff += math.Abs(float64(int(n.b) - int(h.b)))
	diff += math.Abs(float64(int(n.a) - int(h.a)))
	return diff
}

func compareSequence(n, h *image, hIdx int, ch chan result) {
	compareIml(n, h, hIdx, ch, func(x int) int { return x })
}

func compareSequenceReverse(n, h *image, hIdx int, ch chan result) {
	compareIml(n, h, hIdx, ch, func(x int) int { return (n.width - 1) - (x % n.width) + (n.width * (x / n.width)) })
}

func compareIml(n, h *image, hIdx int, ch chan result, indexCalc func(int) int) {
	var counter int
	var accumulator uint64
	hStartPix := hIdx

	for i := 0; i < n.height*n.width; i++ {
		nIdx := indexCalc(i)

		newRow := (i%n.width == 0)
		notFirstRow := ((i / n.width) != 0)
		if newRow && notFirstRow {
			hIdx += (h.width - n.width)
		}

		diff := pixelDiff(n.pixels[nIdx], h.pixels[hIdx])

		if newRow && notFirstRow && counter < 10 {
			return
		}
		if newRow {
			counter = 0
		}
		if diff < threshold {
			counter++
		}

		hIdx++
		accumulator += uint64(diff)
	}

	avgDiff := int(accumulator / uint64(n.height*n.width))
	ch <- result{n, h, hStartPix, avgDiff}
}
