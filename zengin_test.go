package zengin

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseFile(t *testing.T) {

	var mockData = `
12100110999999ｷﾔﾂｼﾕ ﾀﾛｳ                               02249999               010               20999999
22606ｷﾔﾂｼﾕ ﾋｼﾑｹ     020ｵﾓﾀﾆ1              19876543ｷﾔﾂｼ ｼﾖｳｼﾞ                    00000000010                    0
22606ｷﾔﾂｼﾕ ﾋｼﾑｹ     030ｵﾓﾀﾆ2              29999999ｷﾔﾂｼ ﾊﾅｺ                      00000000020                    0
22606ｷﾔﾂｼﾕ ﾋｼﾑｹ     030ｵﾓﾀﾆ2              11234567ｷﾔﾂｼ ｼﾞﾛｳ                     00000000030200504001 000001    0
8000003000000000006
9
`

	expectedData := [][]string{
		{"1210011", "099999", "ｷﾔﾂｼﾕ ﾀﾛｳ", "02249999", "010", "20999999"},
	}

	reader := strings.NewReader(mockData)

	records, err := ParseFile(reader)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(records, expectedData) {
		t.Fatalf("expected %v, got %v", expectedData, records)
	}
}
