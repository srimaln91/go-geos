package geos

import (
	"runtime"
)

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -L/usr/local/lib -lgeos_c
#include <geos_c.h>

*/
import "C"

type coordSequence struct {
	c *C.GEOSCoordSequence
}

func (s *coordSequence) getSize() int {
	var size C.uint
	C.GEOSCoordSeq_getSize_r(ctxHandle, s.c, &size)

	return int(size)
}

func (s *coordSequence) setX(idx int, val float64) {
	C.GEOSCoordSeq_setX_r(ctxHandle, s.c, C.uint(idx), C.double(val))
}

func (s *coordSequence) setY(idx int, val float64) {
	C.GEOSCoordSeq_setY_r(ctxHandle, s.c, C.uint(idx), C.double(val))
}

func (s *coordSequence) setZ(idx int, val float64) {
	C.GEOSCoordSeq_setZ_r(ctxHandle, s.c, C.uint(idx), C.double(val))
}

func (s *coordSequence) getX(idx int) float64 {
	var val C.double
	i := C.GEOSCoordSeq_getX_r(ctxHandle, s.c, C.uint(idx), &val)
	if i == 0 {
		return 0.0
	}

	return float64(val)
}

func (s *coordSequence) getY(idx int) float64 {
	var val C.double
	i := C.GEOSCoordSeq_getY_r(ctxHandle, s.c, C.uint(idx), &val)
	if i == 0 {
		return 0.0
	}

	return float64(val)
}

func (s *coordSequence) getZ(idx int) float64 {
	var val C.double
	i := C.GEOSCoordSeq_getZ_r(ctxHandle, s.c, C.uint(idx), &val)
	if i == 0 {
		return 0.0
	}

	return float64(val)
}

func (s *coordSequence) toCoords() []Coord {
	var coords []Coord

	count := s.getSize()
	for i := 0; i < count; i++ {
		coord := Coord{s.getX(i), s.getY(i)}
		coords = append(coords, coord)
	}

	return coords
}

func (s *coordSequence) toCoordZs() []CoordZ {
	var coords []CoordZ

	count := s.getSize()
	for i := 0; i < count; i++ {
		coord := CoordZ{s.getX(i), s.getY(i), s.getZ(i)}
		coords = append(coords, coord)
	}

	return coords
}

func coordSeqFromC(c *C.GEOSCoordSequence) *coordSequence {
	coordSeq := &coordSequence{c: c}
	runtime.SetFinalizer(coordSeq, func(coordSeq *coordSequence) {
		C.GEOSCoordSeq_destroy_r(ctxHandle, coordSeq.c)
	})

	return coordSeq
}

func coordSeqFromCoords(coords []Coord) *coordSequence {
	size := len(coords)
	coordSeq := createCoordSeq(size, 2)

	for i := 0; i < size; i++ {
		coord := coords[i]
		coordSeq.setX(i, coord.X)
		coordSeq.setY(i, coord.Y)
	}

	return coordSeq
}

func coordSeqFromCoordZs(coords []CoordZ) *coordSequence {
	size := len(coords)
	coordSeq := createCoordSeq(size, 3)

	for i := 0; i < size; i++ {
		coord := coords[i]
		coordSeq.setX(i, coord.X)
		coordSeq.setY(i, coord.Y)
		coordSeq.setZ(i, coord.Z)
	}

	return coordSeq
}

func createCoordSeq(size, dims int) *coordSequence {
	c := C.GEOSCoordSeq_create_r(ctxHandle, C.uint(size), C.uint(dims))
	if c == nil {
		return nil
	}

	return coordSeqFromC(c)
}
