# coco

Package for **co**nverting **co**ordinates between WGS84 Lon Lat, UTM and MGRS/UTMREF.

[![GoDoc](https://godoc.org/github.com/Klaus-Tockloth/coco?status.svg)](https://godoc.org/github.com/Klaus-Tockloth/coco)

## Supported conversions

``` TXT
utm.ToLL()   : converts from UTM to LL
utm.ToMGRS() : converts from UTM to MGRS
ll.ToUTM()   : converts from LL to UTM
ll.ToMGRS()  : converts from LL to MGRS
mgrs.ToUTM() : converts from MGRS to UTM
mgrs.ToLL()  : converts from MGRS to LL
```

## Data objects

``` TXT
UTM  : ZoneNumber ZoneLetter Easting Northing
LL   : Latitude Longitude
MGRS : String
```

## Abbreviations

``` TXT
Lat    : Latitude
Lon    : Longitude
MGRS   : Military Grid Reference System (same as UTMREF)
UTM    : Universal Transverse Mercator
UTMREF : UTM Reference System (same as MGRS)
WGS84  : World Geodetic System 1984 (same as EPSG:4326)
```

## Remarks

* Partial ported from JavaScript [mgrs](https://github.com/proj4js/mgrs) library.
* See utility [coordconv](https://github.com/Klaus-Tockloth/coordconv) for standalone program.
