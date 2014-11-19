/*
The converter package handles conversion between csv and json
*/
package converter

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/morcmarc/csvtoolkit/inferer"
	"github.com/morcmarc/csvtoolkit/utils"
)

type Converter struct {
	input  *os.File
	output *os.File
}

// Returns a new Converter for the given input and output
func NewConverter(csvInput *os.File, jsonOutput *os.File) *Converter {
	converter := &Converter{
		input:  csvInput,
		output: jsonOutput,
	}
	return converter
}

// Processes the input and writes converted objects onto the output
func (c *Converter) Run() {
	cReader := utils.NewDefaultCSVReader(c.input)

	fields, err := cReader.Read()
	if err != nil {
		log.Fatalf("Could not read input: %s", err)
	}
	typeMap, err := inferer.Infer(cReader, fields, 10)
	if err != nil {
		log.Fatalf("Could not infer types: %s", err)
	}

	cReader.Reset()
	cReader.Read()

	r := NewRecords(fields, typeMap)

	c.output.WriteString("[")
	firstItem := true
	for {
		line, err := cReader.Read()
		if err == io.EOF {
			break
		}
		if !firstItem {
			c.output.WriteString(",")
		} else {
			firstItem = false
		}

		j, err := json.Marshal(r.Convert(line))
		if err != nil {
			log.Fatalf("Failed encoding json: %s", err)
		}
		c.output.Write(j)
	}
	c.output.WriteString("]")
}

func getNewCsvReader(in *os.File) *csv.Reader {
	return csv.NewReader(in)
}
