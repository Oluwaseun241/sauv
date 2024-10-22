package cmd

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
)

var (
	subtleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	progressEmpty = subtleStyle.Render(progressEmptyChar)
	ramp          = makeRampStyles("#B14FFF", "#00FFA3", progressBarWidth)
)

func progressbar(percent float64) string {
	w := float64(progressBarWidth)

	fullSize := int(math.Round(w * percent))
	var fullCells string
	for i := 0; i < fullSize; i++ {
		fullCells += ramp[i].Render(progressFullChar)
	}

	emptySize := int(w) - fullSize
	emptyCells := strings.Repeat(progressEmpty, emptySize)

	return fmt.Sprintf("%s%s %3.0f", fullCells, emptyCells, math.Round(percent*100))
}

func makeRampStyles(colorA, colorB string, steps float64) (s []lipgloss.Style) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0.0; i < steps; i++ {
		c := cA.BlendLuv(cB, i/steps)
		s = append(s, lipgloss.NewStyle().Foreground(lipgloss.Color(colorToHex(c))))
	}
	return
}

func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}

func PerformBackupAndMigration(oldDBUrl, newDBUrl string) error {
	fmt.Printf("Performing migration from %s to %s\n", oldDBUrl, newDBUrl)
	totalSteps := 100
	for i := 1; i <= totalSteps; i++ {
		// Simulate each step of migration with a delay
		time.Sleep(15 * time.Millisecond)

		// Calculate progress percentage
		progress := float64(i) / float64(totalSteps)

		// Display the progress bar
		fmt.Printf("\r%s", progressbar(progress)+"%")
	}
	fmt.Println("\nMigration completed successfully!")
	return nil
}
