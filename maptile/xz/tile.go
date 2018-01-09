package xz

import (
	"math/bits"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/internal/mercator"
	"github.com/paulmach/orb/maptile"
)

const maxZoom = 20

// Tile is an x, y, z tile in the xz ordering schema.
// See: http://www.dbs.ifi.lmu.de/Publikationen/Boehm/Ordering_99.pdf
type Tile uint64

// A Zoom is a strict type for a tile zoom level.
type Zoom uint32

// New creates a new tile with the given coordinates.
func New(xi, yi uint32, z Zoom) Tile {
	if z > maxZoom {
		panic("max zoom for this lib is 20")
	}

	x := uint64(xi)
	y := uint64(yi)

	var t Tile
	for i := Zoom(0); i < z; i++ {
		t |= Tile((x & (1 << i)) << (3*(maxZoom-z+i) - i + 1))
		t |= Tile((y & (1 << i)) << (3*(maxZoom-z+i) - i + 2))
		t |= 1 << (3*(maxZoom-z+i) + 0)
	}

	return t
}

// XYZ returns the x, y and z coord of the tile.
func (t Tile) XYZ() (uint32, uint32, Zoom) {
	var z Zoom
	var x, y uint32

	bit := Tile(1 << 57)
	for t&bit != 0 {
		x <<= 1
		if t&(bit<<1) != 0 {
			x |= 1
		}

		y <<= 1
		if t&(bit<<2) != 0 {
			y |= 1
		}

		bit >>= 3
		z++
	}

	return x, y, z
}

// X returns the horizontal/x-axis value of the tile.
func (t Tile) X() uint32 {
	var x uint32

	bit := Tile(1 << 57)
	for t&bit != 0 {
		x <<= 1
		if t&(bit<<1) != 0 {
			x |= 1
		}

		bit >>= 3
	}

	return x
}

// Y returns the vertical/y-axis value of the tile.
// 0 is at the top of the world.
func (t Tile) Y() uint32 {
	var y uint32

	bit := Tile(1 << 57)
	for t&bit != 0 {
		y <<= 1
		if t&(bit<<2) != 0 {
			y |= 1
		}

		bit >>= 3
	}

	return y
}

// Z returns the zoom of the tile.
func (t Tile) Z() Zoom {
	if t == 0 {
		return 0
	}

	c := bits.TrailingZeros64(uint64(t))
	return Zoom(maxZoom - c/3)
}

// At creates a tile for the point at the given zoom.
// Will create a valid tile for the zoom. Points outside
// the range lat [-85.0511, 85.0511] will be snapped to the
// max or min tile as appropriate.
func At(ll orb.Point, z Zoom) Tile {
	f := maptile.Fraction(ll, maptile.Zoom(z))
	x := uint32(f[0])
	y := uint32(f[1])

	if y >= 1<<z {
		y = (1 << z) - 1
	}

	return New(x, y, z)
}

// Bound returns the geo bound for the tile.
// An optional tileBuffer parameter can be passes to create a buffer
// around the bound in tile dimension. e.g. a tileBuffer of 1 would create
// a bound 9x the size of the tile, centered around the provided tile.
func (t Tile) Bound() orb.Bound {
	xi, yi, z := t.XYZ()
	x := float64(xi)
	y := float64(yi)

	lon1, lat1 := mercator.ToGeo(x, y, uint32(z))

	maxtiles := float64(uint32(1 << z))
	maxx := x + 2
	if maxx > maxtiles {
		maxx = maxtiles
	}

	maxy := y + 2
	if maxy > maxtiles {
		maxy = maxtiles
	}

	lon2, lat2 := mercator.ToGeo(maxx, maxy, uint32(z))

	return orb.Bound{
		Min: orb.Point{lon1, lat2},
		Max: orb.Point{lon2, lat1},
	}
}

// Center returns the center of the tile.
func (t Tile) Center() orb.Point {
	return t.Bound().Center()
}
