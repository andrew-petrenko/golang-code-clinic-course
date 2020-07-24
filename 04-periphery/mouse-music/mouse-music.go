package main

// typedef unsigned char Uint8;
// void SineWave(void *userdata, Uint8 *stream, int len);
import "C"
import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
	"reflect"
	"time"
	"unsafe"
)

var toneHz float64 = 440
var volume float64 = 1

const (
	winTitle            = "Go-sdl2-audio-mouse"
	winWidth, winHeight = 800, 600
	pitchRange          = 800
	sampleHz            = 48000
)

// d - delta of phase changing
func dPhase() float64 {
	// 2 * PI = 1 oscillation. 440/48000 - oscillations per second
	return 2 * math.Pi * toneHz / sampleHz
}

//export SineWave
func SineWave(_ unsafe.Pointer, stream *C.Uint8, length C.int) {
	n := int(length)
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	buf := *(*[]C.Uint8)(unsafe.Pointer(&hdr))

	var phase float64
	for i := 0; i < n; i += 2 {
		phase += dPhase()
		sample := C.Uint8((math.Sin(phase) + 0.999999) * 128)
		buf[i] = sample
		buf[i+1] = sample
	}
}

func main() {
	err := sdl.Init(sdl.INIT_AUDIO)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	surface.FillRect(nil, 0)
	window.UpdateSurface()

	spec := &sdl.AudioSpec{
		Freq:     sampleHz,
		Format:   sdl.AUDIO_U8,
		Channels: 1,
		Samples:  2048,
		Callback: sdl.AudioCallback(C.SineWave),
	}

	err = sdl.OpenAudio(spec, nil)
	if err != nil {
		panic(err)
	}

	sdl.PauseAudio(true)

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseMotionEvent:
				sdl.LockAudio()
				toneHz = (float64(t.X)/winWidth)*pitchRange + 100
				volume = float64(t.Y) / winHeight
				sdl.UnlockAudio()
			case *sdl.MouseButtonEvent:
				if t.Button == sdl.BUTTON_LEFT {
					sdl.PauseAudio(t.State == sdl.RELEASED)
				}
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}
