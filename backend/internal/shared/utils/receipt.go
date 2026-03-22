package utils

import (
	"fmt"
	"strings"
)

// ESC/POS Command constants
const (
	EscBoldOn    = "\x1B\x45\x01"
	EscBoldOff   = "\x1B\x45\x00"
	EscAlignCenter = "\x1B\x61\x01"
	EscAlignLeft   = "\x1B\x61\x00"
	EscAlignRight  = "\x1B\x61\x02"
	EscCut         = "\x1D\x56\x01"
	EscInit        = "\x1B\x40"
)

type ReceiptBuilder struct {
	content strings.Builder
}

func NewReceiptBuilder() *ReceiptBuilder {
	b := &ReceiptBuilder{}
	b.content.WriteString(EscInit)
	return b
}

func (b *ReceiptBuilder) AddLine(text string) *ReceiptBuilder {
	b.content.WriteString(text + "\n")
	return b
}

func (b *ReceiptBuilder) AddCentered(text string) *ReceiptBuilder {
	b.content.WriteString(EscAlignCenter + text + "\n" + EscAlignLeft)
	return b
}

func (b *ReceiptBuilder) AddBold(text string) *ReceiptBuilder {
	b.content.WriteString(EscBoldOn + text + EscBoldOff + "\n")
	return b
}

func (b *ReceiptBuilder) AddKeyValue(key, value string) *ReceiptBuilder {
	// Simple 32-column layout
	space := 32 - len(key) - len(value)
	if space < 1 {
		space = 1
	}
	b.content.WriteString(key + strings.Repeat(" ", space) + value + "\n")
	return b
}

func (b *ReceiptBuilder) AddDivider() *ReceiptBuilder {
	b.content.WriteString(strings.Repeat("-", 32) + "\n")
	return b
}

func (b *ReceiptBuilder) Cut() string {
	b.content.WriteString(EscCut)
	return b.content.String()
}

func (b *ReceiptBuilder) String() string {
	return b.content.String()
}

// FormatPrice is a helper for receipt printing
func FormatPrice(val float64) string {
	return fmt.Sprintf("%.2f", val)
}
