package zengin

import (
	"reflect"
	"strings"
	"testing"
)

type test struct {
	name          string
	input         string
	expected      [][]string
	expectedError bool
}

func TestZenginToCSV(t *testing.T) {

	var tests = []test{
		{"InvalidFile", "invalid file", nil, true},
		{"EmptyFile", ``, [][]string{}, true},
		{"NoDataRecord", `12110110999999ｹﾝｼﾝ ﾀﾛｳ                                02242606               010               20999999
8000000000000000000
9`, [][]string{
			{"振込名義人", "金融機関コード", "支店コード", "科目", "口座番号", "口座名義人", "金額"},
		}, false},
		{"OneDataRecord", `12110110999999ｹﾝｼﾝ ﾀﾛｳ                                02242606               010               20999999
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    020ﾋﾖｳｺﾞ              19876543ｹﾝｼﾝ ｼﾖｳｼﾞ                    00000000010                    0
8000001000000000001
9`, [][]string{
			{"振込名義人", "金融機関コード", "支店コード", "科目", "口座番号", "口座名義人", "金額"},
			{"ｹﾝｼﾝ ﾀﾛｳ", "2606", "020", "1", "9876543", "ｹﾝｼﾝ ｼﾖｳｼﾞ", "1"},
		}, false},
		{"WrongDate", `12110110999999ｹﾝｼﾝ ﾀﾛｳ                                13322606               010               20999999
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    020ﾋﾖｳｺﾞ              19876543ｹﾝｼﾝ ｼﾖｳｼﾞ                    00000000010                    0
8000001000000000001
9`, [][]string{}, true},
		{"MultipleDataRecords", `12110110999999ｹﾝｼﾝ ﾀﾛｳ                                02242606               010               20999999
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    020ﾋﾖｳｺﾞ              19876543ｹﾝｼﾝ ｼﾖｳｼﾞ                    00000000010                    0
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              29999999ｹﾝｼﾝ ﾊﾅｺ                      00000000020                    0
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              11234567ｹﾝｼﾝ ｼﾞﾛｳ                     00000000030ﾏｲﾂｷﾌﾞﾝ             0Y
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              11234567ｹﾝｼﾝ ｼﾞﾛｳ                     00000000030200504001 000001    0
8000004000000000009
9`, [][]string{
			{"振込名義人", "金融機関コード", "支店コード", "科目", "口座番号", "口座名義人", "金額"},
			{"ｹﾝｼﾝ ﾀﾛｳ", "2606", "020", "1", "9876543", "ｹﾝｼﾝ ｼﾖｳｼﾞ", "1"},
			{"ｹﾝｼﾝ ﾀﾛｳ", "2606", "030", "2", "9999999", "ｹﾝｼﾝ ﾊﾅｺ", "2"},
			{"ｹﾝｼﾝ ﾀﾛｳ", "2606", "030", "1", "1234567", "ｹﾝｼﾝ ｼﾞﾛｳ", "3"},
			{"ｹﾝｼﾝ ﾀﾛｳ", "2606", "030", "1", "1234567", "ｹﾝｼﾝ ｼﾞﾛｳ", "3"},
		}, false},
		{"ErrorMissingTrailer", `12110110999999ｹﾝｼﾝ ﾀﾛｳ                                02242606               010               20999999
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    020ﾋﾖｳｺﾞ              19876543ｹﾝｼﾝ ｼﾖｳｼﾞ                    00000000010                    0
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              29999999ｹﾝｼﾝ ﾊﾅｺ                      00000000020                    0
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              11234567ｹﾝｼﾝ ｼﾞﾛｳ                     00000000030ﾏｲﾂｷﾌﾞﾝ             0Y
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              11234567ｹﾝｼﾝ ｼﾞﾛｳ                     00000000030200504001 000001    0
9`, [][]string{}, true},
		{"ErrorMissingHeader", `22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    020ﾋﾖｳｺﾞ              19876543ｹﾝｼﾝ ｼﾖｳｼﾞ                    00000000010                    0
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              29999999ｹﾝｼﾝ ﾊﾅｺ                      00000000020                    0
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              11234567ｹﾝｼﾝ ｼﾞﾛｳ                     00000000030ﾏｲﾂｷﾌﾞﾝ             0Y
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              11234567ｹﾝｼﾝ ｼﾞﾛｳ                     00000000030200504001 000001    0
8000004000000000009
9`, [][]string{}, true},
		{"ErrorMissingEnd", `12110110999999ｹﾝｼﾝ ﾀﾛｳ                                02242606               010               20999999
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    020ﾋﾖｳｺﾞ              19876543ｹﾝｼﾝ ｼﾖｳｼﾞ                    00000000010                    0
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              29999999ｹﾝｼﾝ ﾊﾅｺ                      00000000020                    0
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              11234567ｹﾝｼﾝ ｼﾞﾛｳ                     00000000030ﾏｲﾂｷﾌﾞﾝ             0Y
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              11234567ｹﾝｼﾝ ｼﾞﾛｳ                     00000000030200504001 000001    0
8000004000000000009
`, [][]string{}, true},
		{"ErrorWrongOrderDataFirst", `22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    020ﾋﾖｳｺﾞ              19876543ｹﾝｼﾝ ｼﾖｳｼﾞ                    00000000010                    0
12110110999999ｹﾝｼﾝ ﾀﾛｳ                                02242606               010               20999999
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              29999999ｹﾝｼﾝ ﾊﾅｺ                      00000000020                    0
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              11234567ｹﾝｼﾝ ｼﾞﾛｳ                     00000000030ﾏｲﾂｷﾌﾞﾝ             0Y
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              11234567ｹﾝｼﾝ ｼﾞﾛｳ                     00000000030200504001 000001    0
8000004000000000009
9`, [][]string{}, true},
		{"ErrorWrongOrderEndInTheMiddle", `12110110999999ｹﾝｼﾝ ﾀﾛｳ                                02242606               010               20999999
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    020ﾋﾖｳｺﾞ              19876543ｹﾝｼﾝ ｼﾖｳｼﾞ                    00000000010                    0
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              29999999ｹﾝｼﾝ ﾊﾅｺ                      00000000020                    0
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              11234567ｹﾝｼﾝ ｼﾞﾛｳ                     00000000030ﾏｲﾂｷﾌﾞﾝ             0Y
9                                                                                                           
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              11234567ｹﾝｼﾝ ｼﾞﾛｳ                     00000000030200504001 000001    0
8000004000000000009
`, [][]string{}, true},
		{"ErrorWrongOrderTrailerFirst", `8000004000000000009
12110110999999ｹﾝｼﾝ ﾀﾛｳ                                02242606               010               20999999
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    020ﾋﾖｳｺﾞ              19876543ｹﾝｼﾝ ｼﾖｳｼﾞ                    00000000010                    0
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              29999999ｹﾝｼﾝ ﾊﾅｺ                      00000000020                    0
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              11234567ｹﾝｼﾝ ｼﾞﾛｳ                     00000000030ﾏｲﾂｷﾌﾞﾝ             0Y
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              11234567ｹﾝｼﾝ ｼﾞﾛｳ                     00000000030200504001 000001    0
9`, [][]string{}, true},
		{"ErrorMultipleHeadersMissingTrailer", `12110110999999ｹﾝｼﾝ ﾀﾛｳ                                02242606               010               20999999
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    020ﾋﾖｳｺﾞ              19876543ｹﾝｼﾝ ｼﾖｳｼﾞ                    00000000010                    0
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              29999999ｹﾝｼﾝ ﾊﾅｺ                      00000000020                    0
12110110999999ｹﾝｼﾝ ﾀﾛ                                01142606               010               20999999
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              11234567ｹﾝｼﾝ ｼﾞﾛｳ                     00000000030ﾏｲﾂｷﾌﾞﾝ             0Y
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              11234567ｹﾝｼﾝ ｼﾞﾛｳ                     00000000030200504001 000001    0
8000004000000000009
9`, [][]string{}, true},
		{"ErrorMismatchAmount", `
12110110999999ｹﾝｼﾝ ﾀﾛｳ                                02242606               010               20999999                 
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    020ﾋﾖｳｺﾞ              19876543ｹﾝｼﾝ ｼﾖｳｼﾞ                    00000000010                    0        
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              29999999ｹﾝｼﾝ ﾊﾅｺ                      00000000020                    0        
8000002000000000003                                                                                                    
12110110999999ｹﾝｼﾝ ﾀﾛ                                 01142606               010               20999999                 
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              11234567ｹﾝｼﾝ ｼﾞﾛｳ                     00000000030ﾏｲﾂｷﾌﾞﾝ             0Y       
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              11234567ｹﾝｼﾝ ｼﾞﾛｳ                     00000000030200504001 000001    0        
8000002000000000009                                                                                                     
9`, [][]string{}, true},
		{"ErrorMismatchCount", `
12110110999999ｹﾝｼﾝ ﾀﾛｳ                                02242606               010               20999999                 
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    020ﾋﾖｳｺﾞ              19876543ｹﾝｼﾝ ｼﾖｳｼﾞ                    00000000010                    0        
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              29999999ｹﾝｼﾝ ﾊﾅｺ                      00000000020                    0        
8000003000000000003                                                                                                    
12110110999999ｹﾝｼﾝ ﾀﾛ                                 01142606               010               20999999                 
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              11234567ｹﾝｼﾝ ｼﾞﾛｳ                     00000000030ﾏｲﾂｷﾌﾞﾝ             0Y       
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              11234567ｹﾝｼﾝ ｼﾞﾛｳ                     00000000030200504001 000001    0        
8000002000000000006                                                                                                     
9`, [][]string{}, true},
		{"MultipleHeaderRecords", `
12110110999999ｹﾝｼﾝ ﾀﾛｳ                                02242606               010               20999999                 
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    020ﾋﾖｳｺﾞ              19876543ｹﾝｼﾝ ｼﾖｳｼﾞ                    00000000010                    0        
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              29999999ｹﾝｼﾝ ﾊﾅｺ                      00000000020                    0        
8000002000000000003                                                                                                    
12110110999999ｹﾝｼﾝ ﾀﾛ                                 01142606               010               20999999                 
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              11234567ｹﾝｼﾝ ｼﾞﾛｳ                     00000000030ﾏｲﾂｷﾌﾞﾝ             0Y       
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              11234567ｹﾝｼﾝ ｼﾞﾛｳ                     00000000030200504001 000001    0        
8000002000000000006                                                                                                     
9`, [][]string{
			{"振込名義人", "金融機関コード", "支店コード", "科目", "口座番号", "口座名義人", "金額"},
			{"ｹﾝｼﾝ ﾀﾛｳ", "2606", "020", "1", "9876543", "ｹﾝｼﾝ ｼﾖｳｼﾞ", "1"},
			{"ｹﾝｼﾝ ﾀﾛｳ", "2606", "030", "2", "9999999", "ｹﾝｼﾝ ﾊﾅｺ", "2"},
			{"ｹﾝｼﾝ ﾀﾛ", "2606", "030", "1", "1234567", "ｹﾝｼﾝ ｼﾞﾛｳ", "3"},
			{"ｹﾝｼﾝ ﾀﾛ", "2606", "030", "1", "1234567", "ｹﾝｼﾝ ｼﾞﾛｳ", "3"},
		}, false},
		{"MultipleHeaderRecordsWithBadSymbolAtTheEnd", `
12110110999999ｹﾝｼﾝ ﾀﾛｳ                                02242606               010               20999999                 
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    020ﾋﾖｳｺﾞ              19876543ｹﾝｼﾝ ｼﾖｳｼﾞ                    00000000010                    0        
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              29999999ｹﾝｼﾝ ﾊﾅｺ                      00000000020                    0        
8000002000000000003                                                                                                    
12110110999999ｹﾝｼﾝ ﾀﾛ                                 01142606               010               20999999                 
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              11234567ｹﾝｼﾝ ｼﾞﾛｳ                     00000000030ﾏｲﾂｷﾌﾞﾝ             0Y       
22606ﾋﾖｳｺﾞｹﾝｼﾝｸﾐ    030ｻﾝﾉﾐﾔ              11234567ｹﾝｼﾝ ｼﾞﾛｳ                     00000000030200504001 000001    0        
8000002000000000006                                                                                                     
9
b`, [][]string{
			{"振込名義人", "金融機関コード", "支店コード", "科目", "口座番号", "口座名義人", "金額"},
			{"ｹﾝｼﾝ ﾀﾛｳ", "2606", "020", "1", "9876543", "ｹﾝｼﾝ ｼﾖｳｼﾞ", "1"},
			{"ｹﾝｼﾝ ﾀﾛｳ", "2606", "030", "2", "9999999", "ｹﾝｼﾝ ﾊﾅｺ", "2"},
			{"ｹﾝｼﾝ ﾀﾛ", "2606", "030", "1", "1234567", "ｹﾝｼﾝ ｼﾞﾛｳ", "3"},
			{"ｹﾝｼﾝ ﾀﾛ", "2606", "030", "1", "1234567", "ｹﾝｼﾝ ｼﾞﾛｳ", "3"},
		}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := strings.NewReader(test.input)

			records, err := ToCSVJa(reader)
			if test.expectedError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(records, test.expected) {
					t.Fatalf("expected %v, got %v", test.expected, records)
				}
			}
		})
	}
}
