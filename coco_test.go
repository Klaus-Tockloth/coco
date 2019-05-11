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
- https://github.com/chrisveness/geodesy/blob/master/test/utm-mgrs-tests.js
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
		{UTM{ZoneNumber: 32, ZoneLetter: 'U', Easting: 399000, Northing: 5757000}, LL{Lat: 51.954519, Lon: 7.530231}, nil},
		{UTM{ZoneNumber: 32, ZoneLetter: 'U', Easting: 574126, Northing: 5815291}, LL{Lat: 52.482729, Lon: 10.091555}, nil},
		// test set from chris veness
		{UTM{ZoneNumber: 31, ZoneLetter: 'N', Easting: 166021, Northing: 0}, LL{Lat: 0.0, Lon: -0.000004}, nil},
		{UTM{ZoneNumber: 31, ZoneLetter: 'N', Easting: 277438, Northing: 110597}, LL{Lat: 0.999991, Lon: 0.999998}, nil},
		{UTM{ZoneNumber: 30, ZoneLetter: 'M', Easting: 722561, Northing: 9889402}, LL{Lat: -1.0, Lon: -1.000007}, nil},
		{UTM{ZoneNumber: 31, ZoneLetter: 'N', Easting: 448251, Northing: 5411943}, LL{Lat: 48.858293, Lon: 2.294488}, nil},    // eiffel tower
		{UTM{ZoneNumber: 56, ZoneLetter: 'H', Easting: 334873, Northing: 6252266}, LL{Lat: -33.857001, Lon: 151.214998}, nil}, // sidney o/h
		{UTM{ZoneNumber: 18, ZoneLetter: 'N', Easting: 323394, Northing: 4307395}, LL{Lat: 38.897694, Lon: -77.036503}, nil},  // white house
		{UTM{ZoneNumber: 23, ZoneLetter: 'K', Easting: 683466, Northing: 7460687}, LL{Lat: -22.951904, Lon: -43.210602}, nil}, // rio christ
		{UTM{ZoneNumber: 32, ZoneLetter: 'N', Easting: 297508, Northing: 6700645}, LL{Lat: 60.391347, Lon: 5.324893}, nil},    // bergen
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
		{LL{Lat: 51.95, Lon: 7.53}, UTM{ZoneNumber: 32, ZoneLetter: 'U', Easting: 398973, Northing: 5756497}},
		{LL{Lat: 52.482728, Lon: -1.908445}, UTM{ZoneNumber: 30, ZoneLetter: 'U', Easting: 574125, Northing: 5815290}},
		{LL{Lat: -19.887495, Lon: -43.932663}, UTM{ZoneNumber: 23, ZoneLetter: 'K', Easting: 611733, Northing: 7800614}},
		{LL{Lat: 60.0, Lon: 4.0}, UTM{ZoneNumber: 32, ZoneLetter: 'V', Easting: 221288, Northing: 6661953}},  // Norway 31->32
		{LL{Lat: 75.0, Lon: 8.0}, UTM{ZoneNumber: 31, ZoneLetter: 'X', Easting: 644293, Northing: 8329692}},  // Svalbard 32->31
		{LL{Lat: 75.0, Lon: 10.0}, UTM{ZoneNumber: 33, ZoneLetter: 'X', Easting: 355706, Northing: 8329692}}, // Svalbard 32->33
		{LL{Lat: 75.0, Lon: 10.0}, UTM{ZoneNumber: 33, ZoneLetter: 'X', Easting: 355706, Northing: 8329692}}, // Svalbard 34->33
		{LL{Lat: 75.0, Lon: 22.0}, UTM{ZoneNumber: 35, ZoneLetter: 'X', Easting: 355706, Northing: 8329692}}, // Svalbard 34->35
		{LL{Lat: 75.0, Lon: 32.0}, UTM{ZoneNumber: 35, ZoneLetter: 'X', Easting: 644293, Northing: 8329692}}, // Svalbard 36->35
		{LL{Lat: 75.0, Lon: 34.0}, UTM{ZoneNumber: 37, ZoneLetter: 'X', Easting: 355706, Northing: 8329692}}, // Svalbard 36->37
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
		{"30NYF6799300000", UTM{ZoneNumber: 30, ZoneLetter: 'N', Easting: 767993, Northing: 0}, 1, nil},
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
		{LL{Lat: 51.95, Lon: 7.53}, 1, "32ULC9897356497", nil},
		{LL{Lat: 51.95, Lon: 7.53}, 100, "32ULC989564", nil},
		{LL{Lat: -19.887495, Lon: -43.932663}, 1, "23KPU1173300614", nil},
		{LL{Lat: 0.0, Lon: -0.592328}, 1, "30NYF6799300000", nil},
		// negative tests
		{LL{Lat: 51.95, Lon: 188.53}, 100, "", fmt.Errorf("invalid longitude, lon = 188.53")},
		{LL{Lat: 51.95, Lon: -188.53}, 100, "", fmt.Errorf("invalid longitude, lon = -188.53")},
		{LL{Lat: 99.95, Lon: 7.53}, 100, "", fmt.Errorf("invalid latitude, lat = 99.95")},
		{LL{Lat: -99.95, Lon: 7.53}, 100, "", fmt.Errorf("invalid latitude, lat = -99.95")},
		{LL{Lat: 88.95, Lon: 7.53}, 100, "", fmt.Errorf("polar regions below 80°S and above 84°N not supported, lat = 88.95")},
		{LL{Lat: -88.95, Lon: 7.53}, 100, "", fmt.Errorf("polar regions below 80°S and above 84°N not supported, lat = -88.95")},
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
		{"32ULC9897356497", LL{Lat: 51.949993, Lon: 7.529986}, 1, nil},
		{"33UXP04", LL{Lat: 48.205348, Lon: 16.345927}, 10000, nil},
		{"11SPA7234911844", LL{Lat: 36.236123, Lon: -115.082098}, 1, nil},
		{"23KPU1173300614", LL{Lat: -19.887498, Lon: -43.932664}, 1, nil},
		{"31UGT03734554", LL{Lat: 51.823490, Lon: 5.956335}, 10, nil},
		{"30NYF6799300000", LL{Lat: 0.0, Lon: -0.592328}, 1, nil},
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
	// 23K 611733 7800614 -> -19.887498 -43.932664
}

func ExampleLL_ToUTM() {

	ll := LL{Lon: -115.08209766, Lat: 36.23612346}
	utm := ll.ToUTM()
	fmt.Printf("%s -> %s\n", ll, utm)
	// Output:
	// 36.236123 -115.082098 -> 11S 672349 4011843
}

func ExampleUTM_ToMGRS() {

	utm := UTM{ZoneNumber: 31, ZoneLetter: 'U', Easting: 700373, Northing: 5704554}
	accuracy := 1 // meters
	mgrs := utm.ToMGRS(accuracy)
	fmt.Printf("%s -> %s\n", utm, mgrs)
	// Output:
	// 31U 700373 5704554 -> 31UGT0037304554
}

func ExampleMGRS_ToUTM() {

	mgrs := MGRS("32ULC989564")
	utm, accuracy, err := mgrs.ToUTM()
	if err != nil {
		log.Fatalf("error <%v> at mgrs.ToUTM()", err)
	}
	fmt.Printf("%s -> %s (accuracy %d meters)\n", mgrs, utm, accuracy)
	// Output:
	// 32ULC989564 -> 32U 398900 5756400 (accuracy 100 meters)
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
	// 51.950000 -88.530000 -> 16UCC94855658 (accuracy 10 meters)
}

func ExampleMGRS_ToLL() {

	mgrs := MGRS("11SPA7234911844")
	ll, accuracy, err := mgrs.ToLL()
	if err != nil {
		log.Fatalf("error <%v> at mgrs.ToLL()", err)
	}
	fmt.Printf("%s (with accuracy %d meters) -> %s\n", mgrs, accuracy, ll)
	// Output:
	// 11SPA7234911844 (with accuracy 1 meters) -> 36.236123 -115.082098
}
