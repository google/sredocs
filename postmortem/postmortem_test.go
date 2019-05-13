package postmortem

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

type table struct {
	name  string
	value string
}

func TestPM01(t *testing.T) {
	// Order must match Fields order for PM01.
	tables := []table{
		{"shortlink", "http://example.com/yep"},
		{"lastupdated", "2019-04-01"},
		{"teamname", "Kitchen Sink SRE"},
	}

	testParser(t, "testdata/pm01.txt", tables)
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

func toCSV(s string) ([][]string, error) {
	r := csv.NewReader(strings.NewReader(s))
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	// 1st entry is the header with field names.
	if len(records) != 2 {
		return nil, fmt.Errorf("Multiple PMs in a single file?")
	}
	return records, nil
}
