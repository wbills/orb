package xz

import (
	"fmt"
	"math"
	"testing"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/internal/mercator"
)

func TestNew(t *testing.T) {
	for z := Zoom(10); z <= 20; z++ {
		t.Run(fmt.Sprintf("zoom %d", z), func(t *testing.T) {

			for _, city := range mercator.Cities {
				t1 := At(orb.Point{city[1], city[0]}, z)

				t2 := New(t1.X(), t1.Y(), t1.Z())
				if t1 != t2 {
					t.Errorf("incorrect tile: %v %v %v", t2.X(), t2.Y(), t2.Z())
				}

				x, y, z := t1.XYZ()
				t2 = New(x, y, z)
				if t1 != t2 {
					t.Errorf("incorrect tile: %v %v %v", t2.X(), t2.Y(), t2.Z())
				}
			}
		})
	}
}

func TestZ(t *testing.T) {
	for z := Zoom(0); z <= 20; z++ {
		t.Run(fmt.Sprintf("zoom %d", z), func(t *testing.T) {
			if v := New(0, 0, z).Z(); v != z {
				t.Errorf("incorrect zoom: %v", v)
			}
		})
	}
}

func TestAt(t *testing.T) {
	tile := At(orb.Point{0, 0}, 20)
	if b := tile.Bound(); b.Top() != 0 || b.Left() != 0 {
		t.Errorf("incorrect tile bound: %v", b)
	}

	// specific case
	if tile := At(orb.Point{-87.65005229999997, 41.850033}, 20); tile.X() != 268988 || tile.Y() != 389836 {
		t.Errorf("projection incorrect: %v %v %v", tile.X(), tile.Y(), tile.Z())
	}

	for _, city := range mercator.Cities {
		tile := At(orb.Point{city[1], city[0]}, 20)
		c := tile.Center()

		if math.Abs(c[1]-city[0]) > 1e-2 {
			t.Errorf("latitude miss match: %f != %f", c[1], city[0])
		}

		if math.Abs(c[0]-city[1]) > 1e-2 {
			t.Errorf("longitude miss match: %f != %f", c[0], city[1])
		}
	}

	// test polar regions
	if tile := At(orb.Point{0, 89.9}, 18); tile.Y() != 0 {
		t.Errorf("top of the world error: %d != %d", tile.Y(), 0)
	}

	if tile := At(orb.Point{0, -89.9}, 18); tile.Y() != (1<<18)-1 {
		t.Errorf("bottom of the world error: %d != %d", tile.Y(), (1<<18)-1)
	}
}

func TestTileBound(t *testing.T) {
	bound := New(7, 8, 9).Bound()

	level := Zoom(9 + 5) // we're testing point +5 zoom, in same tile
	factor := uint32(5)

	// edges should be within the bound
	p := New(7<<factor+1, 8<<factor+1, level).Center()
	if !bound.Contains(p) {
		t.Errorf("should contain point")
	}

	p = New(7<<factor-2, 8<<factor-2, level).Center()
	if bound.Contains(p) {
		t.Errorf("should not contain point")
	}

	// over one
	p = New(8<<factor-2, 9<<factor-2, level).Center()
	if !bound.Contains(p) {
		t.Errorf("should contain point")
	}

	p = New(8<<factor+1, 9<<factor+1, level).Center()
	if !bound.Contains(p) {
		t.Errorf("should contain point")
	}

	p = New(8<<factor+2, 9<<factor+2, level).Center()
	if !bound.Contains(p) {
		t.Errorf("should contain point")
	}

	// over two
	p = New(9<<factor+1, 10<<factor+1, level).Center()
	if bound.Contains(p) {
		t.Errorf("should not contain point")
	}

	expected := orb.Bound{Min: orb.Point{-180, -85.05112877980659}, Max: orb.Point{180, 85.05112877980659}}
	if b := New(0, 0, 0).Bound(); !b.Equal(expected) {
		t.Errorf("should be full earth, got %v", b)
	}
}
