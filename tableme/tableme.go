package tableme

import (
	"fmt"
	"log"
	"os"
	//"strconv"
	"strings"
	//"github.com/davecgh/go-spew/spew"
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

func printTable(records [][]cell, headers []cell, colWidths []int) {
	var rightMargin string
	var nCells int
	var hdr header

	for _, row := range records {
		nCells = len(row)

		for j, cell := range row {
			hdr = headers[j].(header)

			if hdr.hideCol {
				if j == nCells-1 {
					fmt.Print("\n")
				}
			} else {
				if j == nCells-1 {
					rightMargin = "\n"
				} else {
					rightMargin = "  "
				}

				format := fmt.Sprintf("%%-%ds%s", colWidths[j], rightMargin)

				switch cell.(type) {
				case header:
					if hdr.hideHeaderText {
						fmt.Printf(format, strings.Repeat(" ", colWidths[j]))
					} else {
						fmt.Printf(format, cell.val())
					}
				default:
					fmt.Printf(format, cell.val())
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
	//var err error
	var headersMap map[string]header = make(map[string]header)

	for i, cell := range headerArgs {
		s := strings.Split(cell, ":")
		text := s[0]
		hideCol = false
		hideHeaderText = false
		shift = 0
		//err = nil

		if len(s) > 1 {
			opts := strings.Split(s[1], ",")

			for _, opt := range opts {
				//fmt.Printf("%s\n", opt)

				if opt == "." {
					hideCol = true
				} else if opt == "_" {
					hideHeaderText = true
					//} else if opt[0] == '{' || opt[0] == '}' {
					//num := opt[1:]

					//shift, err = strconv.Atoi(num)
					//if err != nil {
					//fmt.Printf("couldn't convert %s to number\n", num)
					//os.Exit(1)
					//}

					//if opt[0] == '{' {
					//shift = -shift
					//}
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

	//spew.Dump(headersMap)

	// reconcile shifts
	//shift = 0
	//var hdr header
	//for i := 0; i < len(headers); i++ {
	//hdr = headers[i]
	//shift = hdr.shift

	//if shift < 0 {
	//// swap i with it's left neighbor "shift" times
	//for j := 0; j < -shift; j++ {

	//}
	//} else if shift < 0 {

	//}
	//}

	return headers, len(headersMap)
}

func TableMe(headerArgs []string, records [][]string) error {
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
			//var l int

			//if cell != nil {
			l := len(cell)
			//} else {
			//l = 0
			//}

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

	printTable([][]cell{headers}, headers, colWidths)
	printTable(cells, headers, colWidths)

	return nil
}

//func TableMe(headerArgs []string, records [][]string) error {
//ptrRecords := make([][]*string, len(records))

//for i, row := range records {
//var rowSlice []*string = make([]*string, len(row))

//for j, cell := range row {
//rowSlice[j] = cell
//spew.Dump(rowSlice)
//}

//ptrRecords[i] = rowSlice
//}

//return TableMePtr(headerArgs, ptrRecords)
//}
