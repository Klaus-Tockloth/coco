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
LL   : Longitude Latitude
MGRS : String
```

## Abbreviations

``` TXT
Lon    : Longitude
Lat    : Latitude
MGRS   : Military Grid Reference System (same as UTMREF)
UTM    : Universal Transverse Mercator
UTMREF : UTM Reference System (same as MGRS)
WGS84  : World Geodetic System 1984 (same as EPSG:4326)
```

## Remarks

* This package is a partial port from github.com/proj4js/mgrs (JavaScript).
* See utility github.com/Klaus-Tockloth/coordconv for standalone program.
