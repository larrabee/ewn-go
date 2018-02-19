package ewn

import "testing"

func TestStripOutput(t *testing.T) {
	type testString struct {
		inString  string
		text      string
		maxLen    int
		outString string
	}
	testStrings := []testString{
		{"HhjasidiewekncxcPosidsdmandhqwejnnbasdhj", "END", 15, "HhjasidiewekEND"},
		{"HhjasidiewekncxcPosidsdmandhqwejnnbasdhj", "END", 40, "HhjasidiewekncxcPosidsdmandhqwejnnbasdhj"},
		{"ФввыыдйцуфывфЛыдззщщйцуььвфьься", "ЦЦЦ", 15, "ФввыыдйцуфывЦЦЦ"},
		{"ФввыыдйцуфывфЛыдззщщйцуььвфьься", "ЦЦЦ", 31, "ФввыыдйцуфывфЛыдззщщйцуььвфьься"},
		{"ФввыыдйцуфывфЛыдззщщйцуььвфьься", "ЦЦЦ", 3, "ЦЦЦ"},
		{"ФввыыдйцуфывфЛыдззщщйцуььвфьься", "ЦЦЦ", 0, ""},
		{"ФввыыдйцуфывфЛыдззщщйцуььвфьься", "ЦЦЦ", 1, "Ц"},
		{"ФввыыдйцуфывфЛыдззщщйцуььвфьься", "ЦЦЦ", 4, "ФЦЦЦ"},
	}
	for _, test := range testStrings {
		realOut := stripOutput(test.inString, test.maxLen, test.text)
		if realOut != test.outString {
			t.Error(
				"For", test,
				"expected", test.outString,
				"got", realOut,
			)
		}
	}
}

func TestConnStrToStruct(t *testing.T) {
	type testString struct {
		conString string
		outHost   string
		outPort   int
	}
	testStrings := []testString{
		{"google.com", "google.com", 25},
		{"test.google.com", "test.google.com", 25},
		{"test.google.com:587", "test.google.com", 587},
		{"test.google.com:65535", "test.google.com", 65535},
		{"test.google.com:asd", "test.google.com", 0},
	}

	for _, test := range testStrings {
		realOut, _ := connStrToStruct(test.conString)
		if (realOut.host != test.outHost) || (realOut.port != test.outPort) {
			t.Error(
				"For", test,
				"expected", test.outHost, test.outPort,
				"got", realOut,
			)
		}
	}
}
