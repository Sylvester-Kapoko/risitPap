package printer

import (
	"fmt"
	"strings"
)

func center(b *strings.Builder, text string, width int) {
	if len(text) >= width {
		b.WriteString(text[:width])
		b.WriteByte('\n')
		return
	}
	padding := width - len(text)
	left := padding / 2
	b.WriteString(strings.Repeat(" ", left))
	b.WriteString(text)
	b.WriteByte('\n')
}

func separator(b *strings.Builder, width int) {
	b.WriteString(strings.Repeat("-", width))
	b.WriteByte('\n')
}

func itemLine(b *strings.Builder, name, qty, unit, total string, width int) {
	nameWidth := width - 16
	if len(name) > nameWidth {
		name = name[:nameWidth-1] + "."
	}
	fmt.Fprintf(b, "%-*s %3s %6s %6s\n", nameWidth, name, qty, unit, total)
}

func moneyLine(b *strings.Builder, label, amount string, width int) {
	spaces := width - len(label) - len(amount)
	if spaces < 1 {
		spaces = 1
	}
	b.WriteString(label)
	b.WriteString(strings.Repeat(" ", spaces))
	b.WriteString(amount)
	b.WriteByte('\n')
}
