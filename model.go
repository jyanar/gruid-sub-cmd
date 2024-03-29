package main

import (
	"context"
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/anaseto/gruid"
)

type model struct {
	grid       gruid.Grid // The drawing grid.
	startedSub bool       // Whether we've started the Sub loop
	startedCmd bool       // Whether we've started the Cmd loop
}

func NewModel(gd gruid.Grid) *model {
	return &model{
		grid:       gd,
		startedSub: false,
		startedCmd: false,
	}
}

func (m *model) randomizeGrid() {
	it := m.grid.Iterator()
	possible_color_values := []uint{0, 25, 50, 75, 100, 125, 150, 175, 200, 225, 250}
	for it.Next() {
		fg := gruid.Color(possible_color_values[rand.Intn(len(possible_color_values))])
		bg := gruid.Color(possible_color_values[rand.Intn(len(possible_color_values))])
		it.SetCell(gruid.Cell{Rune: 'a' + rune(rand.Intn(26)), Style: gruid.Style{Fg: fg, Bg: bg}})

		//! Stress test
		//! Warning: Memory will grow indefinitely. This is because gruid-sdl maintains
		//! textures in a map of type map[gruid.Cell]*sdl.Texture. Thus, if we're
		//! generating every possible gruid.Cell->sdl.Texture combination, that's
		//! (26 runes)*(4,294,967,295 fg colors)*(4,294,967,295 bg colors) = 4.8e20
		//! unique key-texture pairs.
		// fg := gruid.Color(rand.Intn(256))
		// bg := gruid.Color(rand.Intn(256))
		// it.SetCell(gruid.Cell{
		// 	Rune:  'a' + rune(rand.Intn(26)),
		// 	Style: gruid.Style{Fg: fg, Bg: bg},
		// })
	}
}

type subMsg int
type cmdMsg int

// Update implements gruid.Model.update. It handles keyboard and mouse input
// messages and updates the model in response to them.
func (m *model) Update(msg gruid.Msg) gruid.Effect {

	switch msg := msg.(type) {

	case gruid.MsgInit:
		m.randomizeGrid()

	case subMsg:
		log.Println("Processing subMsg!")
		m.randomizeGrid()

	case cmdMsg:
		log.Println("Processing cmdMsg!")
		m.randomizeGrid()
		return m.myCmd()

	case gruid.MsgKeyDown:
		switch msg.Key {

		case "s":
			if !m.startedSub {
				log.Println("Starting Sub loop...")
				m.startedSub = true
				return m.mySub()
			}

		case "c":
			if !m.startedCmd {
				log.Println("Starting the Cmd loop...")
				m.startedCmd = true
				return m.myCmd()
			}

		case gruid.KeyEscape:
			return gruid.End()

		}
	}
	return nil
}

// Recurring event
func (m *model) mySub() gruid.Sub {
	return func(ctx context.Context, ch chan<- gruid.Msg) {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				ch <- subMsg(1)
			case <-ctx.Done():
				return
			}
		}
	}
}

// Single event
func (m *model) myCmd() gruid.Cmd {
	return func() gruid.Msg {
		t := time.NewTimer(200 * time.Millisecond)
		<-t.C
		return cmdMsg(1)
	}
}

func (m *model) Draw() gruid.Grid {
	return m.grid
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: /pkg/runtime/#MemStats
	log.Printf("\tAlloc = %v MiB", bToMb(m.Alloc))
	log.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	log.Printf("\tSys = %v MiB", bToMb(m.Sys))
	log.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
