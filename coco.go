/*
Purpose:
- MGRS/UTMREF <-> UTM <-> Lon Lat

Description:
- Package for converting coordinates between WGS84 Lon Lat, UTM and MGRS/UTMREF.

Releases:
- v0.1.0 - 2019/05/09 : initial release
- v0.2.0 - 2019/05/10 : coord formatting changed
- v0.2.1 - 2019/05/12 : redundant comments removed

Author:
- Klaus Tockloth

Copyright and license:
- Copyright (c) 2019 Klaus Tockloth
- MIT license

Permission is hereby granted, free of charge, to any person obtaining a copy of this software
and associated documentation files (the Software), to deal in the Software without restriction,
including without limitation the rights to use, copy, modify, merge, publish, distribute,
sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or
substantial portions of the Software.

The software is provided 'as is', without warranty of any kind, express or implied, including
but not limited to the warranties of merchantability, fitness for a particular purpose and
noninfringement. In no event shall the authors or copyright holders be liable for any claim,
damages or other liability, whether in an action of contract, tort or otherwise, arising from,
out of or in connection with the software or the use or other dealings in the software.

Contact (eMail):
- freizeitkarte@googlemail.com

Remarks:
- This library is a partial port from "github.com/proj4js/mgrs" (JavaScript).
- Build library:
  go install
- Test library:
  go test
  go test -cover
  go test -coverprofile=c.out + go tool cover -html=c.out
- Document library:
  godoc
  view document in browser (http://localhost:6060)

Links:
- https://github.com/proj4js/mgrs
- https://gist.github.com/tmcw/285949
*/

/*
Package coco (coordinate conversion) provides methods for converting coordinates between WGS84 Lon Lat, UTM and MGRS/UTMREF.

Supported conversions:
  utm.ToLL()   : converts from UTM to LL
  utm.ToMGRS() : converts from UTM to MGRS
  ll.ToUTM()   : converts from LL to UTM
  ll.ToMGRS()  : converts from LL to MGRS
  mgrs.ToUTM() : converts from MGRS to UTM
  mgrs.ToLL()  : converts from MGRS to LL

Data objects:
  UTM  : ZoneNumber ZoneLetter Easting Northing
  LL   : Latitude Longitude
  MGRS : String

Abbreviations:
  Lat    : Latitude
  Lon    : Longitude
  MGRS   : Military Grid Reference System (same as UTMREF)
  UTM    : Universal Transverse Mercator
  UTMREF : UTM Reference System (same as MGRS)
  WGS84  : World Geodetic System 1984 (same as EPSG:4326)
*/
package coco

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// UTM defines coordinate in Universal Transverse Mercator
type UTM struct {
	ZoneNumber int
	ZoneLetter byte
	Easting    float64
	Northing   float64
}

/*
String returns stringified UTM object.
*/
func (utm UTM) String() string {

	return fmt.Sprintf("%d%c %.0f %.0f", utm.ZoneNumber, utm.ZoneLetter, utm.Easting, utm.Northing)
}

// LL defines coordinate in Longitude / Latitude
type LL struct {
	Lat float64
	Lon float64
}

/*
String returns stringified LL object (order according to ISO-6709, precision 0.11 meter).

*/
func (ll LL) String() string {

	return fmt.Sprintf("%.6f %.6f", ll.Lat, ll.Lon)
}

// MGRS defines cordinate in MGRS/UTMREF
type MGRS string

// setOriginColumnLetters defines the column letters (for easting) of the lower left value, per set.
const setOriginColumnLetters = "AJSAJS"

// setOriginRowLetters defines the row letters (for northing) of the lower left value, per set.
const setOriginRowLetters = "AFAFAF"

// character constants
const (
	charA = 65 // character 'A'
	charI = 73 // character 'I'
	charO = 79 // character 'O'
	charV = 86 // character 'V'
	charZ = 90 // character 'Z'
)

/*
ToMGRS converts Lon Lat to MGRS.
accuracy holds the wanted accuracy in meters. Possible values are 1, 10, 100, 1000 or 10000 meters.
*/
func (ll LL) ToMGRS(accuracy int) (MGRS, error) {

	if ll.Lon < -180 || ll.Lon > 180 {
		return "", fmt.Errorf("invalid longitude, lon = %v", ll.Lon)
	}
	if ll.Lat < -90 || ll.Lat > 90 {
		return "", fmt.Errorf("invalid latitude, lat = %v", ll.Lat)
	}
	if ll.Lat < -80 || ll.Lat > 84 {
		return "", fmt.Errorf("polar regions below 80°S and above 84°N not supported, lat = %v", ll.Lat)
	}

	utm := ll.ToUTM()
	mgrs := utm.ToMGRS(accuracy)

	return mgrs, nil
}

/*
ToLL converts MGRS/UTMREF to Lon Lat.
*/
func (mgrs MGRS) ToLL() (LL, int, error) {

	utm, accuracy, err := mgrs.ToUTM()
	if err != nil {
		return LL{}, 0, fmt.Errorf("error <%v> at mgrs.ToUTM()", err)
	}

	ll, err := utm.ToLL()
	if err != nil {
		return LL{}, 0, fmt.Errorf("error <%v> at utm.ToLL(), utm = %#v", err, utm)
	}

	return ll, accuracy, nil
}

/*
degToRad converts from degrees to radians.
del holds the angle in degrees.
*/
func degToRad(deg float64) float64 {

	return (deg * (math.Pi / 180.0))
}

/*
radToDeg converts from radians to degrees.
rad holds the angle in radians.
*/
func radToDeg(rad float64) float64 {

	return (180.0 * (rad / math.Pi))
}

/*
ToUTM converts Lon Lat to UTM.
*/
func (ll LL) ToUTM() UTM {

	Lat := ll.Lat
	Long := ll.Lon
	a := 6378137.0           //ellip.radius;
	eccSquared := 0.00669438 //ellip.eccsq;
	k0 := 0.9996
	LatRad := degToRad(Lat)
	LongRad := degToRad(Long)

	ZoneNumber := 0 // (int)
	ZoneNumber = int(math.Floor((Long+180)/6) + 1)

	// make sure the longitude 180.00 is in Zone 60
	if Long == 180 {
		ZoneNumber = 60
	}

	// Special zone for Norway
	if Lat >= 56.0 && Lat < 64.0 && Long >= 3.0 && Long < 12.0 {
		ZoneNumber = 32
	}

	// special zones for Svalbard
	if Lat >= 72.0 && Lat < 84.0 {
		if Long >= 0.0 && Long < 9.0 {
			ZoneNumber = 31
		} else if Long >= 9.0 && Long < 21.0 {
			ZoneNumber = 33
		} else if Long >= 21.0 && Long < 33.0 {
			ZoneNumber = 35
		} else if Long >= 33.0 && Long < 42.0 {
			ZoneNumber = 37
		}
	}

	LongOrigin := (ZoneNumber-1)*6 - 180 + 3 // +3 puts origin in middle of zone
	LongOriginRad := degToRad(float64(LongOrigin))

	eccPrimeSquared := eccSquared / (1 - eccSquared)

	N := a / math.Sqrt(1-eccSquared*math.Sin(LatRad)*math.Sin(LatRad))
	T := math.Tan(LatRad) * math.Tan(LatRad)
	C := eccPrimeSquared * math.Cos(LatRad) * math.Cos(LatRad)
	A := math.Cos(LatRad) * (LongRad - LongOriginRad)

	M := a * ((1-eccSquared/4-3*eccSquared*eccSquared/64-5*eccSquared*eccSquared*eccSquared/256)*LatRad - (3*eccSquared/8+3*eccSquared*eccSquared/32+45*eccSquared*eccSquared*eccSquared/1024)*math.Sin(2*LatRad) + (15*eccSquared*eccSquared/256+45*eccSquared*eccSquared*eccSquared/1024)*math.Sin(4*LatRad) - (35*eccSquared*eccSquared*eccSquared/3072)*math.Sin(6*LatRad))

	UTMEasting := (k0*N*(A+(1-T+C)*A*A*A/6.0+(5-18*T+T*T+72*C-58*eccPrimeSquared)*A*A*A*A*A/120.0) + 500000.0)

	UTMNorthing := (k0 * (M + N*math.Tan(LatRad)*(A*A/2+(5-T+9*C+4*C*C)*A*A*A*A/24.0+(61-58*T+T*T+600*C-330*eccPrimeSquared)*A*A*A*A*A*A/720.0)))
	if Lat < 0.0 {
		UTMNorthing += 10000000.0 // 10000000 meters offset for southern hemisphere
	}

	utm := UTM{}
	utm.ZoneNumber = ZoneNumber
	utm.ZoneLetter = getLetterDesignator(Lat)
	utm.Easting = math.Trunc(UTMEasting)
	utm.Northing = math.Trunc(UTMNorthing)

	return utm
}

/*
ToLL converts UTM to Lon Lat.
*/
func (utm UTM) ToLL() (LL, error) {

	zoneNumber := utm.ZoneNumber
	zoneLetter := utm.ZoneLetter
	UTMEasting := utm.Easting
	UTMNorthing := utm.Northing

	// check the ZoneNummber is valid
	if zoneNumber < 0 || zoneNumber > 60 {
		return LL{}, fmt.Errorf("invalid zone number, zone number = %v", zoneNumber)
	}

	k0 := 0.9996
	a := 6378137.0           //ellip.radius;
	eccSquared := 0.00669438 //ellip.eccsq;
	e1 := (1 - math.Sqrt(1-eccSquared)) / (1 + math.Sqrt(1-eccSquared))

	// remove 500,000 meters offset for longitude
	x := UTMEasting - 500000.0
	y := UTMNorthing

	// We must know somehow if we are in the Northern or Southern hemisphere, this is the only time we use the letter.
	// So even if the Zone letter isn't exactly correct it should indicate the hemisphere correctly.
	if zoneLetter < 'N' {
		y -= 10000000.0 // remove 10,000,000 meters offset used
		// for southern hemisphere
	}

	// there are 60 zones with zone 1 being at West -180 to -174
	LongOrigin := (zoneNumber-1)*6 - 180 + 3 // +3 puts origin in middle of zone

	eccPrimeSquared := (eccSquared) / (1 - eccSquared)

	M := y / k0
	mu := M / (a * (1 - eccSquared/4 - 3*eccSquared*eccSquared/64 - 5*eccSquared*eccSquared*eccSquared/256))

	phi1Rad := mu + (3*e1/2-27*e1*e1*e1/32)*math.Sin(2*mu) + (21*e1*e1/16-55*e1*e1*e1*e1/32)*math.Sin(4*mu) + (151*e1*e1*e1/96)*math.Sin(6*mu)

	N1 := a / math.Sqrt(1-eccSquared*math.Sin(phi1Rad)*math.Sin(phi1Rad))
	T1 := math.Tan(phi1Rad) * math.Tan(phi1Rad)
	C1 := eccPrimeSquared * math.Cos(phi1Rad) * math.Cos(phi1Rad)
	R1 := a * (1 - eccSquared) / math.Pow(1-eccSquared*math.Sin(phi1Rad)*math.Sin(phi1Rad), 1.5)
	D := x / (N1 * k0)

	lat := phi1Rad - (N1*math.Tan(phi1Rad)/R1)*(D*D/2-(5+3*T1+10*C1-4*C1*C1-9*eccPrimeSquared)*D*D*D*D/24+(61+90*T1+298*C1+45*T1*T1-252*eccPrimeSquared-3*C1*C1)*D*D*D*D*D*D/720)
	lat = radToDeg(lat)

	lon := (D - (1+2*T1+C1)*D*D*D/6 + (5-2*C1+28*T1-3*C1*C1+8*eccPrimeSquared+24*T1*T1)*D*D*D*D*D/120) / math.Cos(phi1Rad)
	lon = float64(LongOrigin) + radToDeg(lon)

	ll := LL{}
	ll.Lat = lat
	ll.Lon = lon

	return ll, nil
}

/*
getLetterDesignator calculates the MGRS letter designator for the given latitude.
lat holds lat the latitude in WGS84 to get the letter designator for.
*/
func getLetterDesignator(lat float64) byte {

	// This is here as an error flag to show that the Latitude is outside MGRS limits
	LetterDesignator := 'Z'

	if (84 >= lat) && (lat >= 72) {
		LetterDesignator = 'X'
	} else if (72 > lat) && (lat >= 64) {
		LetterDesignator = 'W'
	} else if (64 > lat) && (lat >= 56) {
		LetterDesignator = 'V'
	} else if (56 > lat) && (lat >= 48) {
		LetterDesignator = 'U'
	} else if (48 > lat) && (lat >= 40) {
		LetterDesignator = 'T'
	} else if (40 > lat) && (lat >= 32) {
		LetterDesignator = 'S'
	} else if (32 > lat) && (lat >= 24) {
		LetterDesignator = 'R'
	} else if (24 > lat) && (lat >= 16) {
		LetterDesignator = 'Q'
	} else if (16 > lat) && (lat >= 8) {
		LetterDesignator = 'P'
	} else if (8 > lat) && (lat >= 0) {
		LetterDesignator = 'N'
	} else if (0 > lat) && (lat >= -8) {
		LetterDesignator = 'M'
	} else if (-8 > lat) && (lat >= -16) {
		LetterDesignator = 'L'
	} else if (-16 > lat) && (lat >= -24) {
		LetterDesignator = 'K'
	} else if (-24 > lat) && (lat >= -32) {
		LetterDesignator = 'J'
	} else if (-32 > lat) && (lat >= -40) {
		LetterDesignator = 'H'
	} else if (-40 > lat) && (lat >= -48) {
		LetterDesignator = 'G'
	} else if (-48 > lat) && (lat >= -56) {
		LetterDesignator = 'F'
	} else if (-56 > lat) && (lat >= -64) {
		LetterDesignator = 'E'
	} else if (-64 > lat) && (lat >= -72) {
		LetterDesignator = 'D'
	} else if (-72 > lat) && (lat >= -80) {
		LetterDesignator = 'C'
	}

	return byte(LetterDesignator)
}

/*
ToMGRS converts UTM to MGRS/UTMREF.
accuracy holds the wanted accuracy in meters. Possible values are 1, 10, 100, 1000 or 10000 meters.
*/
func (utm UTM) ToMGRS(accuracy int) MGRS {

	// meters to number of digits
	switch accuracy {
	case 1:
		accuracy = 5
	case 10:
		accuracy = 4
	case 100:
		accuracy = 3
	case 1000:
		accuracy = 2
	case 10000:
		accuracy = 1
	default:
		accuracy = 5
	}

	// prepend with leading zeroes
	seasting := "00000" + fmt.Sprintf("%.0f", utm.Easting)
	snorthing := "00000" + fmt.Sprintf("%.0f", utm.Northing)

	mgrs := fmt.Sprintf("%d%s%s%s%s",
		utm.ZoneNumber,
		string(utm.ZoneLetter),
		get100kID(utm.Easting, utm.Northing, utm.ZoneNumber),
		seasting[len(seasting)-5:len(seasting)-5+accuracy],
		snorthing[len(snorthing)-5:len(snorthing)-5+accuracy])

	return MGRS(mgrs)
}

/*
get100kID gets the two letter 100k designator for a given UTM easting, northing and zone number value.
*/
func get100kID(easting, northing float64, zoneNumber int) string {

	setParm := get100kSetForZone(zoneNumber)
	setColumn := int(math.Floor(easting / 100000))
	setRow := int(math.Floor(northing/100000)) % 20

	return getLetter100kID(setColumn, setRow, setParm)
}

/*
get100kSetForZone gets the MGRS 100K set for a given UTM zone number.
*/
func get100kSetForZone(i int) int {

	// UTM zones are grouped, and assigned to one of a group of 6 sets.
	numberOf100kSets := 6

	setParm := i % numberOf100kSets
	if setParm == 0 {
		setParm = numberOf100kSets
	}

	return setParm
}

/*
getLetter100kID gets the two-letter MGRS 100k designator given information translated from the UTM northing, easting and zone number.
column holds the column index as it relates to the MGRS 100k set spreadsheet, created from the UTM easting. Values are 1-8.
row holds the row index as it relates to the MGRS 100k set spreadsheet, created from the UTM northing value. Values are from 0-19.
parm holds the set block, as it relates to the MGRS 100k set spreadsheet, created from the UTM zone. Values are from 1-60.
*/
func getLetter100kID(column, row, parm int) string {

	// colOrigin and rowOrigin are the letters at the origin of the set
	index := parm - 1
	colOrigin := setOriginColumnLetters[index]
	rowOrigin := setOriginRowLetters[index]

	// colInt and rowInt are the letters to build to return
	colInt := int(colOrigin) + column - 1
	rowInt := int(rowOrigin) + row
	rollover := false

	if colInt > charZ {
		colInt = colInt - charZ + charA - 1
		rollover = true
	}

	if colInt == charI || (colOrigin < charI && colInt > charI) || ((colInt > charI || colOrigin < charI) && rollover) {
		colInt++
	}

	if colInt == charO || (colOrigin < charO && colInt > charO) || ((colInt > charO || colOrigin < charO) && rollover) {
		colInt++
		if colInt == charI {
			colInt++
		}
	}

	if colInt > charZ {
		colInt = colInt - charZ + charA - 1
	}

	if rowInt > charV {
		rowInt = rowInt - charV + charA - 1
		rollover = true
	} else {
		rollover = false
	}

	if ((rowInt == charI) || ((rowOrigin < charI) && (rowInt > charI))) || (((rowInt > charI) || (rowOrigin < charI)) && rollover) {
		rowInt++
	}

	if ((rowInt == charO) || ((rowOrigin < charO) && (rowInt > charO))) || (((rowInt > charO) || (rowOrigin < charO)) && rollover) {
		rowInt++
		if rowInt == charI {
			rowInt++
		}
	}

	if rowInt > charV {
		rowInt = rowInt - charV + charA - 1
	}

	twoLetter := string(colInt) + string(rowInt)
	return twoLetter
}

/*
ToUTM converts MGRS/UTMREF to UTM.
*/
func (mgrs MGRS) ToUTM() (UTM, int, error) {

	mgrsTmp := string(mgrs)
	if mgrs == "" {
		return UTM{}, 0, fmt.Errorf("invalid empty mgrs string")
	}

	mgrsTmp = strings.ToUpper(mgrsTmp)

	sb := ""
	i := 0

	// get Zone number
	re := regexp.MustCompile("[A-Z]")
	for !re.MatchString(string(mgrsTmp[i])) {
		if i >= 2 {
			return UTM{}, 0, fmt.Errorf("bad conversion, mgrs = %s", mgrs)
		}
		sb += string(mgrsTmp[i])
		i++
	}

	zoneNumberTmp, err := strconv.ParseInt(sb, 10, 0)
	if err != nil {
		return UTM{}, 0, fmt.Errorf("error <%v> at strconv.ParseInt(), string = %v", err, sb)
	}
	zoneNumber := int(zoneNumberTmp)

	// A good MGRS string has to be 4-5 digits long, ##AAA/#AAA at least.
	if i == 0 || i+3 > len(mgrsTmp) {
		return UTM{}, 0, fmt.Errorf("bad conversion, mgrs = %s", mgrs)
	}

	zoneLetter := mgrsTmp[i]
	i++

	// Should we check the zone letter here? Why not.
	if zoneLetter <= 'A' || zoneLetter == 'B' || zoneLetter == 'Y' || zoneLetter >= 'Z' || zoneLetter == 'I' || zoneLetter == 'O' {
		return UTM{}, 0, fmt.Errorf("zone letter %v not handled, mgrs = %s", zoneLetter, mgrs)
	}

	hunK := mgrsTmp[i : i+2]
	i += 2

	set := get100kSetForZone(zoneNumber)

	east100k, err := getEastingFromChar(hunK[0], set)
	if err != nil {
		return UTM{}, 0, fmt.Errorf("error <%v> at getEastingFromChar()", err)
	}

	north100k, err := getNorthingFromChar(hunK[1], set)
	if err != nil {
		return UTM{}, 0, fmt.Errorf("error <%v> at getNorthingFromChar()", err)
	}

	// We have a bug where the northing may be 2000000 too low. How do we know when to roll over?
	minNorthing, err := getMinNorthing(zoneLetter)
	if err != nil {
		return UTM{}, 0, fmt.Errorf("error <%v> at getMinNorthing()", err)
	}

	for north100k < minNorthing {
		north100k += 2000000
	}

	// calculate the char index for easting/northing separator
	remainder := len(mgrsTmp) - i

	if remainder%2 != 0 {
		return UTM{}, 0, fmt.Errorf("uneven number of digits, mgrs = %s", mgrs)
	}

	sep := remainder / 2

	sepEasting := 0.0
	sepNorthing := 0.0
	accuracy := 0.0
	if sep > 0 {
		accuracy = 100000.0 / math.Pow(10, float64(sep))

		sepEastingString := mgrsTmp[i : i+sep]
		tmpEasting, err := strconv.ParseFloat(sepEastingString, 64)
		if err != nil {
			return UTM{}, 0, fmt.Errorf("error <%v> at strconv.ParseFloat(), easting string = %v", err, sepEastingString)
		}
		sepEasting = tmpEasting * accuracy

		sepNorthingString := mgrsTmp[i+sep:]
		tmpNorthing, _ := strconv.ParseFloat(sepNorthingString, 64)
		if err != nil {
			return UTM{}, 0, fmt.Errorf("error <%v> at strconv.ParseFloat(), northing string = %v", err, sepNorthingString)
		}
		sepNorthing = tmpNorthing * accuracy
	}

	easting := sepEasting + east100k
	northing := sepNorthing + north100k

	utm := UTM{}
	utm.ZoneNumber = zoneNumber
	utm.ZoneLetter = zoneLetter
	utm.Easting = easting
	utm.Northing = northing

	return utm, int(accuracy), nil
}

/*
getEastingFromChar gets the easting value that should be added to the other, secondary easting value.
e holds the first letter from a two-letter MGRS 100k zone.
set holds the MGRS table set for the zone number.
*/
func getEastingFromChar(e byte, set int) (float64, error) {

	// colOrigin is the letter at the origin of the set for the column
	curCol := setOriginColumnLetters[set-1]
	eastingValue := 100000.0
	rewindMarker := false

	for curCol != e {
		curCol++
		if curCol == charI {
			curCol++
		}
		if curCol == charO {
			curCol++
		}
		if curCol > charZ {
			if rewindMarker {
				return -1.0, fmt.Errorf("bad character: %v", e)
			}
			curCol = charA
			rewindMarker = true
		}
		eastingValue += 100000.0
	}

	return eastingValue, nil
}

/*
getNorthingFromChar gets the northing value that should be added to the other, secondary northing value.
n holds the second letter of the MGRS 100k zone.
set holds the MGRS table set number, which is dependent on the UTM zone number.
Remark: You have to remember that Northings are determined from the equator, and the vertical
cycle of letters mean a 2000000 additional northing meters. This happens
approx. every 18 degrees of latitude. This method does *NOT* count any
additional northings. You have to figure out how many 2000000 meters need
to be added for the zone letter of the MGRS coordinate.
*/
func getNorthingFromChar(n byte, set int) (float64, error) {

	if n > 'V' {
		return 0.0, fmt.Errorf("invalid northing, char = %v", n)
	}

	// rowOrigin is the letter at the origin of the set for the column
	curRow := setOriginRowLetters[set-1]
	northingValue := 0.0
	rewindMarker := false

	for curRow != byte(n) {
		curRow++
		if curRow == charI {
			curRow++
		}
		if curRow == charO {
			curRow++
		}
		// fixing a bug making whole application hang in this loop when 'n' is a wrong character
		if curRow > charV {
			if rewindMarker { // making sure that this loop ends
				return -1.0, fmt.Errorf("bad character, char = %v", n)
			}
			curRow = charA
			rewindMarker = true
		}
		northingValue += 100000.0
	}

	return northingValue, nil
}

/*
getMinNorthing gets the minimum northing value of a MGRS zone.
zoneLetter holds the MGRS zone to get the min northing for.
*/
func getMinNorthing(zoneLetter byte) (float64, error) {

	var northing float64

	switch zoneLetter {
	case 'C':
		northing = 1100000.0
	case 'D':
		northing = 2000000.0
	case 'E':
		northing = 2800000.0
	case 'F':
		northing = 3700000.0
	case 'G':
		northing = 4600000.0
	case 'H':
		northing = 5500000.0
	case 'J':
		northing = 6400000.0
	case 'K':
		northing = 7300000.0
	case 'L':
		northing = 8200000.0
	case 'M':
		northing = 9100000.0
	case 'N':
		northing = 0.0
	case 'P':
		northing = 800000.0
	case 'Q':
		northing = 1700000.0
	case 'R':
		northing = 2600000.0
	case 'S':
		northing = 3500000.0
	case 'T':
		northing = 4400000.0
	case 'U':
		northing = 5300000.0
	case 'V':
		northing = 6200000.0
	case 'W':
		northing = 7000000.0
	case 'X':
		northing = 7900000.0
	default:
		northing = -1.0
	}

	if northing >= 0.0 {
		return northing, nil
	}

	return northing, fmt.Errorf("Invalid zone letter: %v", zoneLetter)
}
