package charter

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"testing"
)

type table struct {
	name  string
	value string
}

func TestCharter01(t *testing.T) {
	// Order must match Fields order for Charter01.
	tables := []table{
		{"lastupdated", "2019-04-01"},
		// TODO(stratus): Include all fields.
	}

	testParser(t, "testdata/charter01.txt", tables)
}

func testParser(t *testing.T, pmfile string, tables []table) {
	b, err := ioutil.ReadFile(pmfile)
	if err != nil {
		t.Fatalf("Can't open test data")
	}
	s, err := Parse(Fields, b)
	if err != nil {
		t.Fatal(err)
	}
	records, err := toCSV(s)
	if err != nil {
		t.Fatal(err)
	}

	for i, m := range tables {
		if m.value != records[1][i] {
			t.Errorf("Expected %s, got %s.", m.value, records[1][i])
		}
	}
}

func toCSV(b *bytes.Buffer) ([][]string, error) {
	r := csv.NewReader(bytes.NewReader(b.Bytes()))
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	// 1st entry is the header with field names.
	if len(records) != 2 {
		return nil, fmt.Errorf("Multiple charters in a single file?")
	}
	return records, nil
}
