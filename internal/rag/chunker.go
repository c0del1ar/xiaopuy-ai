package rag

import "strings"

// ChunkText splits plain text into bounded chunks while preferring paragraph boundaries.
func ChunkText(text string, maxChars, overlap int) []string {
	text = strings.TrimSpace(text)
	if text == "" || maxChars <= 0 {
		return nil
	}
	if overlap < 0 {
		overlap = 0
	}
	if overlap >= maxChars {
		overlap = maxChars / 5
	}

	paragraphs := strings.FieldsFunc(text, func(r rune) bool { return r == '\n' })
	var chunks []string
	var current string

	flush := func() {
		current = strings.TrimSpace(current)
		if current != "" {
			chunks = append(chunks, current)
		}
		current = ""
	}

	for _, paragraph := range paragraphs {
		paragraph = strings.TrimSpace(paragraph)
		if paragraph == "" {
			continue
		}
		if current == "" {
			if len(paragraph) <= maxChars {
				current = paragraph
				continue
			}
			chunks = append(chunks, splitLong(paragraph, maxChars, overlap)...)
			continue
		}

		candidate := current + "\n\n" + paragraph
		if len(candidate) <= maxChars {
			current = candidate
			continue
		}
		flush()
		if len(paragraph) <= maxChars {
			current = paragraph
		} else {
			chunks = append(chunks, splitLong(paragraph, maxChars, overlap)...)
		}
	}
	flush()
	return chunks
}

func splitLong(text string, maxChars, overlap int) []string {
	var result []string
	for start := 0; start < len(text); {
		end := start + maxChars
		if end >= len(text) {
			result = append(result, strings.TrimSpace(text[start:]))
			break
		}
		cut := end
		for i := end; i > start; i-- {
			if text[i-1] == ' ' {
				cut = i - 1
				break
			}
		}
		if cut <= start {
			cut = end
		}
		result = append(result, strings.TrimSpace(text[start:cut]))
		start = cut - overlap
		if start < 0 {
			start = 0
		}
		for start < len(text) && text[start] == ' ' {
			start++
		}
	}
	return result
}
