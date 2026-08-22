package rag

import "strings"

// ContextBuilder turns retrieved chunks into bounded, clearly attributed model context.
type ContextBuilder struct {
	MaxChars int
}

func (b ContextBuilder) Build(chunks []Chunk) string {
	if len(chunks) == 0 {
		return ""
	}
	maxChars := b.MaxChars
	if maxChars <= 0 {
		maxChars = 12000
	}

	var out strings.Builder
	for i, chunk := range chunks {
		block := "[Knowledge " + itoa(i+1) + "]\n"
		if chunk.Metadata["source"] != "" {
			block += "Source: " + chunk.Metadata["source"] + "\n"
		}
		if chunk.Metadata["title"] != "" {
			block += "Title: " + chunk.Metadata["title"] + "\n"
		}
		block += "Content: " + strings.TrimSpace(chunk.Content) + "\n\n"
		if out.Len()+len(block) > maxChars {
			break
		}
		out.WriteString(block)
	}
	return strings.TrimSpace(out.String())
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
