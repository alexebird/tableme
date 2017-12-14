package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/alexebird/tableme/tableme"
	//"github.com/davecgh/go-spew/spew"
	"gopkg.in/urfave/cli.v1"
)

func readCsv(inputFile io.Reader) [][]string {
	var buf bytes.Buffer

	if _, err := io.Copy(&buf, inputFile); err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(strings.NewReader(buf.String()))
	buf.Reset()

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return records
}

func main() {
	app := cli.NewApp()
	app.Name = "tableme"

	app.Action = func(c *cli.Context) error {
		records := readCsv(os.Stdin)
		bites := tableme.TableMe(c.Args(), records)
		buff := bytes.NewBuffer(bites)
		fmt.Print(buff.String())
		return nil
	}

	app.Run(os.Args)
}
