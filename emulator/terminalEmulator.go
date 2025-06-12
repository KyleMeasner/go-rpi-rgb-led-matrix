package emulator

import (
	"fmt"
	"image/color"
	"time"
)

type TerminalEmulator struct {
	Width            int
	Height           int
	Pixels           []color.Color
	ShowRefreshRate  bool
	RenderTimestamps []time.Time
}

func NewTerminalEmulator(width, height int, showRefreshRate bool) *TerminalEmulator {
	return &TerminalEmulator{
		Width:            width,
		Height:           height,
		Pixels:           make([]color.Color, width*height),
		ShowRefreshRate:  showRefreshRate,
		RenderTimestamps: []time.Time{},
	}
}

func (t *TerminalEmulator) Geometry() (int, int) {
	return t.Width, t.Height
}

func (t *TerminalEmulator) At(position int) color.Color {
	if t.Pixels[position] == nil {
		return color.Black
	}
	return t.Pixels[position]
}

func (t *TerminalEmulator) Set(position int, c color.Color) {
	t.Pixels[position] = color.RGBAModel.Convert(c)
}

func (t *TerminalEmulator) Apply([]color.Color) error {
	defer func() {
		t.Pixels = make([]color.Color, t.Height*t.Width)
		fmt.Print("\033[39m\033[49m") // Reset colors
	}()

	fmt.Println()

	for row := 0; row < t.Height; row++ {
		for col := 0; col < t.Width; col++ {
			r, g, b, _ := t.At(col + (row * t.Width)).RGBA()
			fmt.Printf("\033[38;2;%d;%d;%dm⬤ ", r, g, b)
		}
		fmt.Println()
	}

	t.RenderTimestamps = append(t.RenderTimestamps, time.Now())
	if len(t.RenderTimestamps) > 60 {
		t.RenderTimestamps = t.RenderTimestamps[1:]
	}
	if t.ShowRefreshRate && len(t.RenderTimestamps) == 60 {
		timeElapsed := t.RenderTimestamps[59].Sub(t.RenderTimestamps[0])
		adjustment := 1 / timeElapsed.Seconds()
		refreshRate := (60 / timeElapsed.Seconds()) * adjustment
		fmt.Printf("\033[39m\033[49m%.2f FPS", refreshRate)
	}

	fmt.Printf("\033[%dA", t.Height+1) // Move cursor up t.Height rows
	return nil
}

func (t *TerminalEmulator) Render() error {
	return t.Apply(t.Pixels)
}

func (t *TerminalEmulator) Close() error {
	return nil
}
