package deku

import "strings"

type Writer struct {
	content *string
	indent int
}

func NewWriter() *Writer {
	initial := ""
	return &Writer{
		content: &initial,
		indent: 1,
	}
}

func (writer *Writer) AddLine(content string) *Writer {
	*writer.content += content + "\n"
	return writer
}

func (writer *Writer) AddIndentedLine(content string) *Writer {
	writer.AddLine(strings.Repeat("  ", writer.indent) + content)

	return &Writer{
		indent: writer.indent + 1,
		content: writer.content,
	}
}

func (writer *Writer) ToString() string {
	content := *writer.content
	if strings.HasSuffix(content, "\n") {
		return content[:len(content)-len("\n")]
	}
	return content
}

func (writer *Writer) Append(content string) *Writer {
	*writer.content += content
	return writer
}