/*
Purpose:
- MGRS/UTMREF <-> UTM <-> Lon Lat

Description:
- testing

Releases:
- v0.1.0 - 2019/05/09 : initial release

Author:
- Klaus Tockloth

Remarks:
- NN
*/

package coco

import (
	"fmt"
	"log"
	"testing"
)

func TestUTM_ToLL(t *testing.T) {

	var tests = []struct {
		utm UTM   // in
		ll  LL    // out
		err error // out
	}{
		// positive tests
		{UTM{ZoneNumber: 32, ZoneLetter: 'U', Easting: 399000, Northing: 5757000}, LL{Lon: 7.53023117, Lat: 51.95451906}, nil},
		{UTM{ZoneNumber: 32, ZoneLetter: 'U', Easting: 574126, Northing: 5815291}, LL{Lon: 10.09155526, Lat: 52.48272900}, nil},
		// negative tests
		{UTM{ZoneNumber: 132, ZoneLetter: 'U', Easting: 574126, Northing: 5815291}, LL{}, fmt.Errorf("invalid zone number, zone number = 132")},
	}

	for _, test := range tests {
		ll, err := test.utm.ToLL()
		function := fmt.Sprintf("utm = %s, ToLL()", test.utm)
		got := fmt.Sprintf("%s %v", ll, err)
		want := fmt.Sprintf("%s %v", test.ll, test.err)
		if got != want {
			t.Errorf("\n%s -> %s != %s\n", function, got, want)
		}
	}
}

func TestLL_ToUTM(t *testing.T) {

	var tests = []struct {
		ll  LL  // in
		utm UTM // out
	}{
		// positive tests
		{LL{Lon: 7.53, Lat: 51.95}, UTM{ZoneNumber: 32, ZoneLetter: 'U', Easting: 398973, Northing: 5756497}},
		{LL{Lon: -1.908445, Lat: 52.482728}, UTM{ZoneNumber: 30, ZoneLetter: 'U', Easting: 574125, Northing: 5815290}},
		{LL{Lon: -43.932663, Lat: -19.887495}, UTM{ZoneNumber: 23, ZoneLetter: 'K', Easting: 611733, Northing: 7800614}},
		// negative tests
		// nothing to do here
	}

	for _, test := range tests {
		utm := test.ll.ToUTM()
		function := fmt.Sprintf("ll = %s, ToUTM()", test.ll)
		got := utm.String()
		want := test.utm.String()
		if got != want {
			t.Errorf("\n%s -> %s != %s\n", function, got, want)
		}
	}
}

func TestUTM_ToMGRS(t *testing.T) {

	var tests = []struct {
		utm      UTM    // in
		accuracy int    // in
		mgrs     string // out
	}{
		// positive tests
		{UTM{ZoneNumber: 32, ZoneLetter: 'U', Easting: 398973, Northing: 5756497}, 1, "32ULC9897356497"},
		{UTM{ZoneNumber: 32, ZoneLetter: 'U', Easting: 398973, Northing: 5756497}, 10, "32ULC98975649"},
		{UTM{ZoneNumber: 32, ZoneLetter: 'U', Easting: 398973, Northing: 5756497}, 100, "32ULC989564"},
		{UTM{ZoneNumber: 32, ZoneLetter: 'U', Easting: 398973, Northing: 5756497}, 1000, "32ULC9856"},
		{UTM{ZoneNumber: 32, ZoneLetter: 'U', Easting: 398973, Northing: 5756497}, 10000, "32ULC95"},
		{UTM{ZoneNumber: 23, ZoneLetter: 'K', Easting: 611733, Northing: 7800614}, 1, "23KPU1173300614"},
		// negative tests
		// nothing to do here
	}

	for _, test := range tests {
		mgrs := test.utm.ToMGRS(test.accuracy)
		function := fmt.Sprintf("utm = %s, ToMGRS(%d)", test.utm, test.accuracy)
		got := string(mgrs)
		want := test.mgrs
		if got != want {
			t.Errorf("\n%s -> %s != %s\n", function, got, want)
		}
	}
}

func TestMGRS_ToUTM(t *testing.T) {

	var tests = []struct {
		mgrs     MGRS  // in
		utm      UTM   // out
		accuracy int   // out
		err      error // out
	}{
		// positive tests
		{"32ULC9897356497", UTM{ZoneNumber: 32, ZoneLetter: 'U', Easting: 398973, Northing: 5756497}, 1, nil},
		{"32ULC98975649", UTM{ZoneNumber: 32, ZoneLetter: 'U', Easting: 398970, Northing: 5756490}, 10, nil},
		{"32ULC989564", UTM{ZoneNumber: 32, ZoneLetter: 'U', Easting: 398900, Northing: 5756400}, 100, nil},
		{"32ULC9856", UTM{ZoneNumber: 32, ZoneLetter: 'U', Easting: 398000, Northing: 5756000}, 1000, nil},
		{"32ULC95", UTM{ZoneNumber: 32, ZoneLetter: 'U', Easting: 390000, Northing: 5750000}, 10000, nil},
		{"23KPU1173300614", UTM{ZoneNumber: 23, ZoneLetter: 'K', Easting: 611733, Northing: 7800614}, 1, nil},
		{"18TWL9334507672", UTM{ZoneNumber: 18, ZoneLetter: 'T', Easting: 593345, Northing: 4507672}, 1, nil},
		{"10SGJ0683244683", UTM{ZoneNumber: 10, ZoneLetter: 'S', Easting: 706832, Northing: 4344683}, 1, nil},
		{"31UGT0037304554", UTM{ZoneNumber: 31, ZoneLetter: 'U', Easting: 700373, Northing: 5704554}, 1, nil},
		// negative tests
		{"", UTM{}, 0, fmt.Errorf("invalid empty mgrs string")},
	}

	for _, test := range tests {
		utm, accuracy, err := test.mgrs.ToUTM()
		function := fmt.Sprintf("mgrs = %s, ToUTM()", test.mgrs)
		got := fmt.Sprintf("%s %d %v", utm, accuracy, err)
		want := fmt.Sprintf("%s %d %v", test.utm, test.accuracy, test.err)
		if got != want {
			t.Errorf("\n%s -> %s != %s\n", function, got, want)
		}
	}
}

func TestLL_ToMGRS(t *testing.T) {

	var tests = []struct {
		ll       LL     // in
		accuracy int    // in
		mgrs     string // out
		err      error  // out
	}{
		// positive tests
		{LL{Lon: 7.53, Lat: 51.95}, 1, "32ULC9897356497", nil},
		{LL{Lon: 7.53, Lat: 51.95}, 100, "32ULC989564", nil},
		{LL{Lon: -43.932663, Lat: -19.887495}, 1, "23KPU1173300614", nil},
		// negative tests
		{LL{Lon: 188.53, Lat: 51.95}, 100, "", fmt.Errorf("invalid longitude, lon = 188.53")},
		{LL{Lon: -188.53, Lat: 51.95}, 100, "", fmt.Errorf("invalid longitude, lon = -188.53")},
		{LL{Lon: 7.53, Lat: 99.95}, 100, "", fmt.Errorf("invalid latitude, lat = 99.95")},
		{LL{Lon: 7.53, Lat: -99.95}, 100, "", fmt.Errorf("invalid latitude, lat = -99.95")},
		{LL{Lon: 7.53, Lat: 88.95}, 100, "", fmt.Errorf("polar regions below 80°S and above 84°N not supported, lat = 88.95")},
		{LL{Lon: 7.53, Lat: -88.95}, 100, "", fmt.Errorf("polar regions below 80°S and above 84°N not supported, lat = -88.95")},
	}

	for _, test := range tests {
		mgrs, err := test.ll.ToMGRS(test.accuracy)
		function := fmt.Sprintf("ll = %s, ll.ToMGRS(%d)", test.ll, test.accuracy)
		got := fmt.Sprintf("%s %v", mgrs, err)
		want := fmt.Sprintf("%s %v", test.mgrs, test.err)
		if got != want {
			t.Errorf("\n%s -> %s != %s\n", function, got, want)
		}
	}
}

func TestMGRS_ToLL(t *testing.T) {

	var tests = []struct {
		mgrs     MGRS  // in
		ll       LL    // out
		accuracy int   // out
		err      error // out
	}{
		// positive tests
		{"32ULC9897356497", LL{Lon: 7.52998627, Lat: 51.94999316}, 1, nil},
		{"33UXP04", LL{Lon: 16.34592696, Lat: 48.20534841}, 10000, nil},
		{"11SPA7234911844", LL{Lon: -115.08209766, Lat: 36.23612346}, 1, nil},
		{"23KPU1173300614", LL{Lon: -43.93266429, Lat: -19.88749831}, 1, nil},
		{"31UGT03734554", LL{Lon: 5.95633528, Lat: 51.82349008}, 10, nil},
		// negative tests
		{"32ULC9897356497CORRUPT", LL{}, 0, fmt.Errorf("error <uneven number of digits, mgrs = 32ULC9897356497CORRUPT> at mgrs.ToUTM()")},
	}

	for _, test := range tests {
		ll, accuracy, err := test.mgrs.ToLL()
		function := fmt.Sprintf("mgrs = %s, mgrs.ToLL()", test.mgrs)
		got := fmt.Sprintf("%s %d %v", ll, accuracy, err)
		want := fmt.Sprintf("%s %d %v", test.ll, test.accuracy, test.err)
		if got != want {
			t.Errorf("\n%s -> %s != %s\n", function, got, want)
		}
	}
}

func ExampleUTM_ToLL() {

	utm := UTM{ZoneNumber: 23, ZoneLetter: 'K', Easting: 611733, Northing: 7800614}
	ll, err := utm.ToLL()
	if err != nil {
		log.Fatalf("error <%v> at utm.ToLL()", err)
	}
	fmt.Printf("%s -> %s\n", utm, ll)
	// Output:
	// 23 K 611733 7800614 -> -43.93266429 -19.88749831
}

func ExampleLL_ToUTM() {

	ll := LL{Lon: -115.08209766, Lat: 36.23612346}
	utm := ll.ToUTM()
	fmt.Printf("%s -> %s\n", ll, utm)
	// Output:
	// -115.08209766 36.23612346 -> 11 S 672349 4011843
}

func ExampleUTM_ToMGRS() {

	utm := UTM{ZoneNumber: 31, ZoneLetter: 'U', Easting: 700373, Northing: 5704554}
	accuracy := 1 // meters
	mgrs := utm.ToMGRS(accuracy)
	fmt.Printf("%s -> %s\n", utm, mgrs)
	// Output:
	// 31 U 700373 5704554 -> 31UGT0037304554
}

func ExampleMGRS_ToUTM() {

	mgrs := MGRS("32ULC989564")
	utm, accuracy, err := mgrs.ToUTM()
	if err != nil {
		log.Fatalf("error <%v> at mgrs.ToUTM()", err)
	}
	fmt.Printf("%s -> %s (accuracy %d meters)\n", mgrs, utm, accuracy)
	// Output:
	// 32ULC989564 -> 32 U 398900 5756400 (accuracy 100 meters)
}

func ExampleLL_ToMGRS() {

	ll := LL{Lon: -88.53, Lat: 51.95}
	accuracy := 10 // meters
	mgrs, err := ll.ToMGRS(accuracy)
	if err != nil {
		log.Fatalf("error <%v> at ll.ToMGRS()", err)
	}
	fmt.Printf("%s -> %s (accuracy %d meters)\n", ll, mgrs, accuracy)
	// Output:
	// -88.53000000 51.95000000 -> 16UCC94855658 (accuracy 10 meters)
}

func ExampleMGRS_ToLL() {

	mgrs := MGRS("11SPA7234911844")
	ll, accuracy, err := mgrs.ToLL()
	if err != nil {
		log.Fatalf("error <%v> at mgrs.ToLL()", err)
	}
	fmt.Printf("%s (with accuracy %d meters) -> %s\n", mgrs, accuracy, ll)
	// Output:
	// 11SPA7234911844 (with accuracy 1 meters) -> -115.08209766 36.23612346
}
