package plist

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func BenchmarkXMLDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		var bval interface{}
		buf := bytes.NewReader([]byte(plistValueTreeAsXML))
		b.StartTimer()
		decoder := NewDecoder(buf)
		decoder.Decode(bval)
		b.StopTimer()
	}
}

func BenchmarkBplistDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		var bval interface{}
		buf := bytes.NewReader(plistValueTreeAsBplist)
		b.StartTimer()
		decoder := NewDecoder(buf)
		decoder.Decode(bval)
		b.StopTimer()
	}
}

func TestLaxDecode(t *testing.T) {
	var plistValueTreeStringsOnlyAsXML string = xmlPreamble + `<plist version="1.0"><dict><key>intarray</key><array><string>1</string><string>8</string><string>16</string><string>32</string><string>64</string><string>2</string><string>9</string><string>17</string><string>33</string><string>65</string></array><key>floats</key><array><string>32</string><string>64</string></array><key>booleans</key><array><string>1</string><string>0</string></array><key>strings</key><array><string>Hello, ASCII</string><string>Hello, 世界</string></array><key>data</key><data>AQIDBA==</data><key>string</key><string>2013-11-27T00:34:00Z</string></dict></plist>`
	d := EverythingTestData{}
	buf := bytes.NewReader([]byte(plistValueTreeStringsOnlyAsXML))
	decoder := NewDecoder(buf)
	decoder.lax = true
	err := decoder.Decode(&d)
	if err != nil {
		t.Error(err.Error())
	}
	t.Logf("%#v", d)
}

func TestDecode(t *testing.T) {
	var failed bool
	for _, test := range tests {
		if test.SkipDecode {
			continue
		}

		failed = false

		t.Logf("Testing Decode (%s)", test.Name)

		d := test.DecodeData
		if d == nil {
			d = test.Data
		}

		testData := reflect.ValueOf(test.Data)
		if !testData.IsValid() || isEmptyInterface(testData) {
			continue
		}
		if testData.Kind() == reflect.Ptr || testData.Kind() == reflect.Interface {
			testData = testData.Elem()
		}
		//typ := testData.Type()

		var err error
		var bval interface{}
		var xval interface{}
		var val interface{}

		if test.ExpectedBin != nil {
			bval = reflect.New(testData.Type()).Interface()
			buf := bytes.NewReader(test.ExpectedBin)
			decoder := NewDecoder(buf)
			err = decoder.Decode(bval)
			vt := reflect.ValueOf(bval)
			if vt.Kind() == reflect.Ptr || vt.Kind() == reflect.Interface {
				vt = vt.Elem()
				bval = vt.Interface()
			}
			val = bval
			if !reflect.DeepEqual(d, bval) {
				failed = true
			}
		}

		if !test.SkipDecodeXML && test.ExpectedXML != "" {
			xval = reflect.New(testData.Type()).Interface()
			buf := bytes.NewReader([]byte(test.ExpectedXML))
			decoder := NewDecoder(buf)
			err = decoder.Decode(xval)
			vt := reflect.ValueOf(xval)
			if vt.Kind() == reflect.Ptr || vt.Kind() == reflect.Interface {
				vt = vt.Elem()
				xval = vt.Interface()
			}
			val = xval
			if !reflect.DeepEqual(d, xval) {
				failed = true
			}
		}

		if bval != nil && xval != nil {
			if !reflect.DeepEqual(bval, xval) {
				t.Log("Binary and XML decoding yielded different values.")
				t.Log("Binary:", bval)
				t.Log("XML   :", xval)
				failed = true
			}
		}

		if failed {
			t.Log("Expected:", d)

			if err == nil {
				t.Log("Received:", val)
			} else {
				t.Log("   Error:", err)
			}
			t.Log("FAILED")
			t.Fail()
		}
	}
}

func TestInterfaceDecode(t *testing.T) {
	var xval interface{}
	buf := bytes.NewReader([]byte{98, 112, 108, 105, 115, 116, 48, 48, 214, 1, 13, 17, 21, 25, 27, 2, 14, 18, 22, 26, 28, 88, 105, 110, 116, 97, 114, 114, 97, 121, 170, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 16, 1, 16, 8, 16, 16, 16, 32, 16, 64, 16, 2, 16, 9, 16, 17, 16, 33, 16, 65, 86, 102, 108, 111, 97, 116, 115, 162, 15, 16, 34, 66, 0, 0, 0, 35, 64, 80, 0, 0, 0, 0, 0, 0, 88, 98, 111, 111, 108, 101, 97, 110, 115, 162, 19, 20, 9, 8, 87, 115, 116, 114, 105, 110, 103, 115, 162, 23, 24, 92, 72, 101, 108, 108, 111, 44, 32, 65, 83, 67, 73, 73, 105, 0, 72, 0, 101, 0, 108, 0, 108, 0, 111, 0, 44, 0, 32, 78, 22, 117, 76, 84, 100, 97, 116, 97, 68, 1, 2, 3, 4, 84, 100, 97, 116, 101, 51, 65, 184, 69, 117, 120, 0, 0, 0, 8, 21, 30, 41, 43, 45, 47, 49, 51, 53, 55, 57, 59, 61, 68, 71, 76, 85, 94, 97, 98, 99, 107, 110, 123, 142, 147, 152, 157, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 29, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 166})
	decoder := NewDecoder(buf)
	err := decoder.Decode(&xval)
	if err != nil {
		t.Log("Error:", err)
		t.Fail()
	}
}

func ExampleDecoder_Decode() {
	type sparseBundleHeader struct {
		InfoDictionaryVersion string `plist:"CFBundleInfoDictionaryVersion"`
		BandSize              uint64 `plist:"band-size"`
		BackingStoreVersion   int    `plist:"bundle-backingstore-version"`
		DiskImageBundleType   string `plist:"diskimage-bundle-type"`
		Size                  uint64 `plist:"size"`
	}

	buf := bytes.NewReader([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
	<dict>
		<key>CFBundleInfoDictionaryVersion</key>
		<string>6.0</string>
		<key>band-size</key>
		<integer>8388608</integer>
		<key>bundle-backingstore-version</key>
		<integer>1</integer>
		<key>diskimage-bundle-type</key>
		<string>com.apple.diskimage.sparsebundle</string>
		<key>size</key>
		<integer>4398046511104</integer>
	</dict>
</plist>`))

	var data sparseBundleHeader
	decoder := NewDecoder(buf)
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)

	// Output: {6.0 8388608 1 com.apple.diskimage.sparsebundle 4398046511104}
}
