package geos

import (
	"math"
	"runtime"
	"unsafe"
)

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -L/usr/local/lib -lgeos_c
#include <geos_c.h>
#include <stdlib.h>
*/
import "C"

const (
	POINT              GeometryType = C.GEOS_POINT
	LINESTRING         GeometryType = C.GEOS_LINESTRING
	LINEARRING         GeometryType = C.GEOS_LINEARRING
	POLYGON            GeometryType = C.GEOS_POLYGON
	MULTIPOINT         GeometryType = C.GEOS_MULTIPOINT
	MULTILINESTRING    GeometryType = C.GEOS_MULTILINESTRING
	MULTIPOLYGON       GeometryType = C.GEOS_MULTIPOLYGON
	GEOMETRYCOLLECTION GeometryType = C.GEOS_GEOMETRYCOLLECTION

	CAP_ROUND  CapType = C.GEOSBUF_CAP_ROUND
	CAP_FLAT   CapType = C.GEOSBUF_CAP_FLAT
	CAP_SQUARE CapType = C.GEOSBUF_CAP_SQUARE

	JOIN_ROUND JoinType = C.GEOSBUF_JOIN_ROUND
	JOIN_MITRE JoinType = C.GEOSBUF_JOIN_MITRE
	JOIN_BEVEL JoinType = C.GEOSBUF_JOIN_BEVEL
)

type GeometryType int

type CapType int
type JoinType int

func (t GeometryType) ToString() string {
	switch t {
	case POINT:
		return "Point"
	case LINESTRING:
		return "LineString"
	case LINEARRING:
		return "LinearRing"
	case POLYGON:
		return "Polygon"
	case MULTIPOINT:
		return "MultiPoint"
	case MULTILINESTRING:
		return "MultiLineString"
	case MULTIPOLYGON:
		return "MultiPolygon"
	case GEOMETRYCOLLECTION:
		return "GeometryCollection"
	default:
		return "Unknown"
	}
}

type Geometry struct {
	c *C.GEOSGeometry
}

func (g *Geometry) Clone() *Geometry {
	c := C.GEOSGeom_clone_r(ctxHandle, g.c)
	return geomFromC(c)
}

func (g *Geometry) GetType() GeometryType {
	return GeometryType(C.GEOSGeomTypeId_r(ctxHandle, g.c))
}

func (g *Geometry) SetSRID(srid int) {
	C.GEOSSetSRID_r(ctxHandle, g.c, C.int(srid))
}

func (g *Geometry) GetSRID() int {
	return int(C.GEOSGetSRID_r(ctxHandle, g.c))
}

func (g *Geometry) GetNumGeometries() int {
	return int(C.GEOSGetNumGeometries_r(ctxHandle, g.c))
}

func (g *Geometry) GetGeometryN(n int) *Geometry {
	c := C.GEOSGetGeometryN_r(ctxHandle, g.c, C.int(n))
	return geomFromC(c)
}

// Only support Polygon
func (g *Geometry) GetNumInteriorRings() int {
	return int(C.GEOSGetNumInteriorRings_r(ctxHandle, g.c))
}

// Only support Polygon
func (g *Geometry) GetInteriorRingN(n int) *Geometry {
	c := C.GEOSGetInteriorRingN_r(ctxHandle, g.c, C.int(n))
	return geomFromC(c)
}

// Only support Polygon
func (g *Geometry) GetExteriorRing() *Geometry {
	c := C.GEOSGetExteriorRing_r(ctxHandle, g.c)
	return geomFromC(c)
}

func (g *Geometry) GetNumCoordinates() int {
	return int(C.GEOSGetNumCoordinates_r(ctxHandle, g.c))
}

// Only support Point
func (g *Geometry) GetXY() (float64, float64) {
	var x, y C.double
	C.GEOSGeomGetX_r(ctxHandle, g.c, &x)
	C.GEOSGeomGetY_r(ctxHandle, g.c, &y)

	return float64(x), float64(y)
}

// Only support LineString, LinearRing or Point
func (g *Geometry) GetCoords() []Coord {
	c := C.GEOSGeom_getCoordSeq_r(ctxHandle, g.c)
	coordSeq := coordSeqFromC(c)
	return coordSeq.toCoords()
}

// Only support LineString, LinearRing or Point
func (g *Geometry) GetCoordZs() []CoordZ {
	c := C.GEOSGeom_getCoordSeq_r(ctxHandle, g.c)
	coordSeq := coordSeqFromC(c)
	return coordSeq.toCoordZs()
}

func (g *Geometry) Area() float64 {
	var val C.double
	C.GEOSArea_r(ctxHandle, g.c, &val)
	return float64(val)
}

func (g *Geometry) Length() float64 {
	var val C.double
	C.GEOSLength_r(ctxHandle, g.c, &val)
	return float64(val)
}

func (g *Geometry) Distance(g2 *Geometry) float64 {
	var val C.double
	C.GEOSDistance_r(ctxHandle, g.c, g2.c, &val)
	return float64(val)
}

func (g *Geometry) HausdorffDistance(g2 *Geometry) float64 {
	var val C.double
	C.GEOSHausdorffDistance_r(ctxHandle, g.c, g2.c, &val)
	return float64(val)
}

func (g *Geometry) HausdorffDistanceDensify(g2 *Geometry, densifyFrac float64) float64 {
	var val C.double

	if densifyFrac > 1 {
		densifyFrac = 1
	} else if densifyFrac <= 0 {
		densifyFrac = math.SmallestNonzeroFloat64
	}

	C.GEOSHausdorffDistanceDensify_r(ctxHandle, g.c, g2.c, C.double(densifyFrac), &val)
	return float64(val)
}

func (g *Geometry) NearestPoints(g2 *Geometry) []Coord {
	c := C.GEOSNearestPoints_r(ctxHandle, g.c, g2.c)
	coordSeq := coordSeqFromC(c)
	return coordSeq.toCoords()
}

func (g *Geometry) NearestPointZs(g2 *Geometry) []CoordZ {
	c := C.GEOSNearestPoints_r(ctxHandle, g.c, g2.c)
	coordSeq := coordSeqFromC(c)
	return coordSeq.toCoordZs()
}

func (g *Geometry) Disjoint(g2 *Geometry) bool {
	flag := C.GEOSDisjoint_r(ctxHandle, g.c, g2.c)
	return flag == C.char(1)
}

func (g *Geometry) Touches(g2 *Geometry) bool {
	flag := C.GEOSTouches_r(ctxHandle, g.c, g2.c)
	return flag == C.char(1)
}

func (g *Geometry) Intersects(g2 *Geometry) bool {
	flag := C.GEOSIntersects_r(ctxHandle, g.c, g2.c)
	return flag == C.char(1)
}

func (g *Geometry) Crosses(g2 *Geometry) bool {
	flag := C.GEOSCrosses_r(ctxHandle, g.c, g2.c)
	return flag == C.char(1)
}

func (g *Geometry) Within(g2 *Geometry) bool {
	flag := C.GEOSWithin_r(ctxHandle, g.c, g2.c)
	return flag == C.char(1)
}

func (g *Geometry) Contains(g2 *Geometry) bool {
	flag := C.GEOSContains_r(ctxHandle, g.c, g2.c)
	return flag == C.char(1)
}

func (g *Geometry) Overlaps(g2 *Geometry) bool {
	flag := C.GEOSOverlaps_r(ctxHandle, g.c, g2.c)
	return flag == C.char(1)
}

func (g *Geometry) Equals(g2 *Geometry) bool {
	flag := C.GEOSEquals_r(ctxHandle, g.c, g2.c)
	return flag == C.char(1)
}

func (g *Geometry) EqualsExact(g2 *Geometry, tol float64) bool {
	flag := C.GEOSEqualsExact_r(ctxHandle, g.c, g2.c, C.double(tol))
	return flag == C.char(1)
}

func (g *Geometry) Covers(g2 *Geometry) bool {
	flag := C.GEOSCovers_r(ctxHandle, g.c, g2.c)
	return flag == C.char(1)
}

func (g *Geometry) CoveredBy(g2 *Geometry) bool {
	flag := C.GEOSCoveredBy_r(ctxHandle, g.c, g2.c)
	return flag == C.char(1)
}

func (g *Geometry) RelatePattern(g2 *Geometry, pattern string) bool {
	cs := C.CString(pattern)
	defer C.free(unsafe.Pointer(cs))

	flag := C.GEOSRelatePattern_r(ctxHandle, g.c, g2.c, cs)
	return flag == C.char(1)
}

func (g *Geometry) Relate(g2 *Geometry) string {
	return C.GoString(C.GEOSRelate_r(ctxHandle, g.c, g2.c))
}

func (g *Geometry) Normalize() {
	C.GEOSNormalize_r(ctxHandle, g.c)
}

func (g *Geometry) IsValid() bool {
	flag := C.GEOSisValid_r(ctxHandle, g.c)
	return flag == C.char(1)
}

func (g *Geometry) IsEmpty() bool {
	flag := C.GEOSisEmpty_r(ctxHandle, g.c)
	return flag == C.char(1)
}

func (g *Geometry) IsSimple() bool {
	flag := C.GEOSisSimple_r(ctxHandle, g.c)
	return flag == C.char(1)
}

func (g *Geometry) IsRing() bool {
	flag := C.GEOSisRing_r(ctxHandle, g.c)
	return flag == C.char(1)
}

func (g *Geometry) HasZ() bool {
	flag := C.GEOSHasZ_r(ctxHandle, g.c)
	return flag == C.char(1)
}

func (g *Geometry) IsClosed() bool {
	flag := C.GEOSisClosed_r(ctxHandle, g.c)
	return flag == C.char(1)
}

func (g *Geometry) ToWKT() string {
	writer := createWktWriter()
	return writer.write(g)
}

func (g *Geometry) ToWKB() []byte {
	writer := createWkbWriter()
	return writer.write(g)
}

func (g *Geometry) Buffer(width float64, quadsegs int, endCapStyle CapType, joinStyle JoinType, mitreLimit float64) *Geometry {
	c := C.GEOSBufferWithStyle_r(ctxHandle, g.c, C.double(width), C.int(quadsegs),
		C.int(endCapStyle), C.int(joinStyle), C.double(mitreLimit))
	return geomFromC(c)
}

// Only support LineString.
// Negative width for right side offset, positive width for left side offset.
func (g *Geometry) OffsetCurve(width float64, quadsegs int, joinStyle JoinType, mitreLimit float64) *Geometry {
	c := C.GEOSOffsetCurve_r(ctxHandle, g.c, C.double(width), C.int(quadsegs),
		C.int(joinStyle), C.double(mitreLimit))
	return geomFromC(c)
}

func CreateFromWKT(wkt string) *Geometry {
	reader := createWktReader()
	return reader.read(wkt)
}

func CreateFromWKB(wkb []byte) *Geometry {
	reader := createWkbReader()
	return reader.read(wkb)
}

func CreatePoint(x, y float64) *Geometry {
	s := createCoordSeq(1, 2)
	if s == nil || s.c == nil {
		return nil
	}

	s.setX(0, x)
	s.setY(0, y)

	c := C.GEOSGeom_createPoint_r(ctxHandle, s.c)
	if c == nil {
		return nil
	}

	return geomFromC(c)
}

func CreatePointZ(x, y, z float64) *Geometry {
	s := createCoordSeq(1, 3)
	if s == nil || s.c == nil {
		return nil
	}

	s.setX(0, x)
	s.setY(0, y)
	s.setZ(0, z)

	c := C.GEOSGeom_createPoint_r(ctxHandle, s.c)
	if c == nil {
		return nil
	}

	return geomFromC(c)
}

func CreateLinearRing(coords []Coord) *Geometry {
	coordSeq := coordSeqFromCoords(coords)
	c := C.GEOSGeom_createLinearRing_r(ctxHandle, coordSeq.c)

	return geomFromC(c)
}

func CreateLinearRingZ(coords []CoordZ) *Geometry {
	coordSeq := coordSeqFromCoordZs(coords)
	c := C.GEOSGeom_createLinearRing_r(ctxHandle, coordSeq.c)

	return geomFromC(c)
}

func CreateLineString(coords []Coord) *Geometry {
	coordSeq := coordSeqFromCoords(coords)
	c := C.GEOSGeom_createLineString_r(ctxHandle, coordSeq.c)

	return geomFromC(c)
}

func CreateLineStringZ(coords []CoordZ) *Geometry {
	coordSeq := coordSeqFromCoordZs(coords)
	c := C.GEOSGeom_createLineString_r(ctxHandle, coordSeq.c)

	return geomFromC(c)
}

func CreatePolygon(shell []Coord, holes ...[]Coord) *Geometry {
	shellGeom := CreateLinearRing(shell)
	if shellGeom == nil {
		return nil
	}
	shellC := shellGeom.c

	var holesPtr **C.GEOSGeometry
	var holeCs []*C.GEOSGeometry
	for i := range holes {
		holeGeom := CreateLinearRing(holes[i])
		holeCs = append(holeCs, holeGeom.c)
	}

	holeCount := len(holeCs)
	if holeCount > 0 {
		holesPtr = &holeCs[0]
	}

	c := C.GEOSGeom_createPolygon_r(ctxHandle, shellC, holesPtr, C.uint(holeCount))
	return geomFromC(c)
}

func CreatePolygonZ(shell []CoordZ, holes ...[]CoordZ) *Geometry {
	shellGeom := CreateLinearRingZ(shell)
	if shellGeom == nil {
		return nil
	}
	shellC := shellGeom.c

	var holesPtr **C.GEOSGeometry
	var holeCs []*C.GEOSGeometry
	for i := range holes {
		holeGeom := CreateLinearRingZ(holes[i])
		holeCs = append(holeCs, holeGeom.c)
	}

	holeCount := len(holeCs)
	if holeCount > 0 {
		holesPtr = &holeCs[0]
	}

	c := C.GEOSGeom_createPolygon_r(ctxHandle, shellC, holesPtr, C.uint(holeCount))
	return geomFromC(c)
}

func CreateMultiGeometry(geoms []*Geometry, geomType GeometryType) *Geometry {
	var geomsCs []*C.GEOSGeometry
	for i := range geoms {
		geom := geoms[i]
		thisType := geom.GetType()
		switch {
		case geomType == MULTIPOINT && thisType != POINT:
			{
				continue
			}
		case geomType == MULTILINESTRING && thisType != LINESTRING:
			{
				continue
			}
		case geomType == MULTIPOLYGON && thisType != POLYGON:
			{
				continue
			}
		}

		geomsCs = append(geomsCs, geom.c)
	}

	geomCount := len(geomsCs)
	if geomCount > 0 {
		geomsPtr := &geomsCs[0]
		c := C.GEOSGeom_createCollection_r(ctxHandle, C.int(geomType), geomsPtr, C.uint(geomCount))
		return geomFromC(c)
	}

	return nil
}

func geomFromC(c *C.GEOSGeometry) *Geometry {
	geom := &Geometry{c: c}
	runtime.SetFinalizer(geom, func(geom *Geometry) {
		C.GEOSGeom_destroy_r(ctxHandle, geom.c)
	})

	return geom
}