package tableme

import (
	"bytes"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"log"
	"os"
	"strings"
)

type cell interface {
	val() string
}

type header struct {
	text           string
	shift          int
	hideCol        bool
	hideHeaderText bool
	index          int
}

func (h header) val() string {
	return h.text
}

type data struct {
	text string
}

func (d data) val() string {
	return d.text
}

func bufferTable(buff *bytes.Buffer, records [][]cell, headers []cell, colWidths []int) {
	var rightMargin string
	var nCells int
	var hdr header

	for _, row := range records {
		nCells = len(row)

		for j, cell := range row {
			hdr = headers[j].(header)

			if hdr.hideCol {
				if j == nCells-1 {
					buff.WriteString(fmt.Sprintf("\n"))
				}
			} else {
				if j == nCells-1 {
					rightMargin = "\n"
				} else {
					rightMargin = "  "
				}

				colWidth := colWidths[j]

				// dont leave trailing spaces on the last column.
				if j == nCells-1 {
					colWidth = len(cell.val())
				}

				format := fmt.Sprintf("%%-%ds%s", colWidth, rightMargin)

				switch cell.(type) {
				case header:
					if hdr.hideHeaderText {
						buff.WriteString(fmt.Sprintf(format, strings.Repeat(" ", colWidths[j])))
					} else {
						buff.WriteString(fmt.Sprintf(format, cell.val()))
					}
				default:
					buff.WriteString(fmt.Sprintf(format, cell.val()))
				}
			}
		}
	}
}

func parseHeaders(headerArgs []string, colCount int) ([]cell, int) {
	var headers []cell = make([]cell, 0, colCount)
	var hideCol bool
	var hideHeaderText bool
	var shift int
	var headersMap map[string]header = make(map[string]header)

	for i, cell := range headerArgs {
		s := strings.Split(cell, ":")
		text := s[0]
		hideCol = false
		hideHeaderText = false
		shift = 0

		if len(s) > 1 {
			opts := strings.Split(s[1], ",")

			for _, opt := range opts {

				if opt == "." {
					hideCol = true
				} else if opt == "_" {
					hideHeaderText = true
				} else {
					fmt.Printf("unknown option: %s\n", opt)
					os.Exit(1)
				}
			}
		}

		newHeader := header{
			hideCol:        hideCol,
			hideHeaderText: hideHeaderText,
			shift:          shift,
			text:           text,
			index:          i,
		}

		if existingHeader, ok := headersMap[newHeader.text]; !ok {
			headers = append(headers, newHeader)
			headersMap[newHeader.text] = newHeader
		} else {
			headers[existingHeader.index] = newHeader
		}
	}

	return headers, len(headersMap)
}

func TableMe(headerArgs []string, records [][]string, noHeaders bool) []byte {
	//spew.Dump(records)
	if len(records) == 0 {
		os.Exit(0)
	}

	colCount := len(records[0])
	colWidths := make([]int, len(records[0]))

	// validate row sizes
	for _, row := range records {
		l := len(row)
		if l != colCount {
			log.Fatal("all rows must be the same width")
			os.Exit(1)
		}
	}

	// convert headers to structs
	var headers []cell
	var uniqueCount int
	headers, uniqueCount = parseHeaders(headerArgs, colCount)

	if uniqueCount != len(records[0]) {
		log.Fatal("header count must be same as row width")
		os.Exit(1)
	}

	// get colWidths
	for j, cell := range headers {
		l := len(cell.val())
		hdr, _ := cell.(header)

		if l > colWidths[j] && !hdr.hideHeaderText {
			colWidths[j] = l
		}
	}

	for _, row := range records {
		for j, cell := range row {
			l := len(cell)

			if l > colWidths[j] {
				colWidths[j] = l
			}
		}
	}

	// convert records to structs
	var cells [][]cell = make([][]cell, 0, len(records))

	for _, row := range records {
		rowCells := make([]cell, 0, colCount)
		for _, cell := range row {
			rowCells = append(rowCells, data{text: cell})
		}
		cells = append(cells, rowCells)
	}

	var buffer *bytes.Buffer = &bytes.Buffer{}
	if !noHeaders {
		bufferTable(buffer, [][]cell{headers}, headers, colWidths)
	}
	bufferTable(buffer, cells, headers, colWidths)

	return buffer.Bytes()
}
