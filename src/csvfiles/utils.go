package csvfiles

import (
	"bytes"
	"encoding/csv"
)

func GetCSVData(p []byte) ([][]string, error) {
	b := bytes.NewBuffer(p)
	r := csv.NewReader(b)
	r.Comma = ';'
	r.FieldsPerRecord = -1
	r.TrimLeadingSpace = true
	return r.ReadAll()
}
