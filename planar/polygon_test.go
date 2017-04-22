package planar

import (
	"math"
	"testing"
)

func TestPolygonDistanceFrom(t *testing.T) {
	r1 := append(NewLineString(),
		NewPoint(0, 0),
		NewPoint(3, 0),
		NewPoint(3, 3),
		NewPoint(0, 3),
		NewPoint(0, 0),
	)

	r2 := append(NewLineString(),
		NewPoint(1, 1),
		NewPoint(2, 1),
		NewPoint(2, 2),
		NewPoint(1, 2),
		NewPoint(1, 1),
	)

	poly := append(NewPolygon(), r1, r2)

	cases := []struct {
		name   string
		point  Point
		result float64
	}{
		{
			name:   "outside",
			point:  NewPoint(-1, 2),
			result: 1,
		},
		{
			name:   "inside",
			point:  NewPoint(0.4, 2),
			result: 0,
		},
		{
			name:   "in hole",
			point:  NewPoint(1.3, 1.4),
			result: 0.3,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if d := poly.DistanceFrom(tc.point); math.Abs(d-tc.result) > epsilon {
				t.Errorf("incorrect distance: %v != %v", d, tc.result)
			}
		})
	}
}

func TestPolygonCentroid(t *testing.T) {
	cases := []struct {
		name   string
		points []Point
		result Point
	}{
		{
			name:   "triangle, cw",
			points: []Point{{0, 0}, {1, 3}, {2, 0}, {0, 0}},
			result: Point{1, 1},
		},
		{
			name:   "triangle, ccw",
			points: []Point{{0, 0}, {2, 0}, {1, 3}, {0, 0}},
			result: Point{1, 1},
		},
		{
			name:   "square, cw",
			points: []Point{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}},
			result: Point{0.5, 0.5},
		},
		{
			name:   "triangle, ccw",
			points: []Point{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
			result: Point{0.5, 0.5},
		},
		{
			name:   "redudent points",
			points: []Point{{0, 0}, {1, 0}, {2, 0}, {1, 3}, {0, 0}},
			result: Point{1, 1},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			poly := Polygon{LineString(tc.points)}
			if c := poly.Centroid(); !c.Equal(tc.result) {
				t.Errorf("wrong centroid: %v != %v", c, tc.result)
			}
		})
	}
}

func TestPolygonCentroidAdv(t *testing.T) {
	ls := append(NewLineString(),
		NewPoint(0, 0),
		NewPoint(0, 1),
		NewPoint(1, 1),
		NewPoint(1, 0.5),
		NewPoint(2, 0.5),
		NewPoint(2, 1),
		NewPoint(3, 1),
		NewPoint(3, 0),
		NewPoint(0, 0),
	)

	// +-+ +-+
	// | | | |
	// | +-+ |
	// |     |
	// +-----+
	poly := Polygon{ls}

	expected := NewPoint(1.5, 0.45)
	if c := poly.Centroid(); !c.Equal(expected) {
		t.Errorf("incorrect centroid: %v != %v", c, expected)
	}
}

func TestPolygonContains(t *testing.T) {
	ls := append(NewLineString(),
		NewPoint(0, 0),
		NewPoint(0, 1),
		NewPoint(1, 1),
		NewPoint(1, 0.5),
		NewPoint(2, 0.5),
		NewPoint(2, 1),
		NewPoint(3, 1),
		NewPoint(3, 0),
		NewPoint(0, 0),
	)

	// +-+ +-+
	// | | | |
	// | +-+ |
	// |     |
	// +-----+
	poly := Polygon{ls}

	cases := []struct {
		name   string
		point  Point
		result bool
	}{
		{
			name:   "in base",
			point:  Point{1.5, 0.25},
			result: true,
		},
		{
			name:   "in right tower",
			point:  Point{0.5, 0.75},
			result: true,
		},
		{
			name:   "in middle",
			point:  Point{1.5, 0.75},
			result: false,
		},
		{
			name:   "in left tower",
			point:  Point{2.5, 0.75},
			result: true,
		},
		{
			name:   "in tp middle",
			point:  Point{1.5, 1.0},
			result: false,
		},
		{
			name:   "above",
			point:  Point{2.5, 1.75},
			result: false,
		},
		{
			name:   "below",
			point:  Point{2.5, -1.75},
			result: false,
		},
		{
			name:   "left",
			point:  Point{-2.5, -0.75},
			result: false,
		},
		{
			name:   "right",
			point:  Point{3.5, 0.75},
			result: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			val := poly.Contains(tc.point)

			if val != tc.result {
				t.Errorf("wrong containment: %v != %v", val, tc.result)
			}
		})
	}

	// points should all be in
	for i, p := range poly[0] {
		if !poly.Contains(p) {
			t.Errorf("point index %d: should be inside", i)
		}
	}

	// on all the segments should be in.
	for i := 1; i < len(poly[0]); i++ {
		c := NewSegment(poly[0][i], poly[0][i-1]).Centroid()
		if !poly.Contains(c) {
			t.Errorf("index %d centroid: should be inside", i)
		}
	}

	// colinear with segments but outside
	for i := 1; i < len(poly[0]); i++ {
		p := NewSegment(poly[0][i], poly[0][i-1]).Interpolate(5)
		if poly.Contains(p) {
			t.Errorf("index %d centroid: should not be inside", i)
		}

		p = NewSegment(poly[0][i], poly[0][i-1]).Interpolate(-5)
		if poly.Contains(p) {
			t.Errorf("index %d centroid: should not be inside", i)
		}
	}
}

func TestPolygonArea(t *testing.T) {
	p1 := append(NewLineString(),
		NewPoint(0, 0),
		NewPoint(3, 0),
		NewPoint(3, 3),
		NewPoint(0, 3),
		NewPoint(0, 0),
	)

	p2 := append(NewLineString(),
		NewPoint(1, 1),
		NewPoint(2, 1),
		NewPoint(2, 2),
		NewPoint(1, 2),
		NewPoint(1, 1),
	)

	polygon := append(NewPolygon(), p1, p2)

	expected := 8.0
	if a := polygon.Area(); a != expected {
		t.Errorf("incorrect area: %v != %v", a, expected)
	}
}

func TestLineStringArea(t *testing.T) {
	cases := []struct {
		name   string
		points []Point
		result float64
	}{
		{
			name:   "simple box, ccw",
			points: []Point{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
			result: 1,
		},
		{
			name:   "simple box, cc",
			points: []Point{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}},
			result: 1,
		},
		{
			name:   "even number of points",
			points: []Point{{0, 0}, {1, 0}, {1, 1}, {0.4, 1}, {0, 1}, {0, 0}},
			result: 1,
		},
		{
			name:   "4 points",
			points: []Point{{0, 0}, {1, 0}, {1, 1}, {0, 0}},
			result: 0.5,
		},
		{
			name:   "6 points",
			points: []Point{{1, 1}, {2, 1}, {2, 1.5}, {2, 2}, {1, 2}, {1, 1}},
			result: 1.0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			val := lineStringArea(tc.points)

			if val != tc.result {
				t.Errorf("wrong area: %v != %v", val, tc.result)
			}
		})
	}
}

func TestPolygonWKT(t *testing.T) {
	r1 := append(NewLineString(),
		NewPoint(0, 0),
		NewPoint(1, 0),
		NewPoint(1, 1),
		NewPoint(0, 1),
		NewPoint(0, 0),
	)

	poly := Polygon{r1}
	expected := "POLYGON((0 0,1 0,1 1,0 1,0 0))"
	if w := poly.WKT(); w != expected {
		t.Errorf("incorrect wkt: %v", w)
	}

	r2 := append(NewLineString(),
		NewPoint(0.4, 0.4),
		NewPoint(0.6, 0.4),
		NewPoint(0.6, 0.6),
		NewPoint(0.4, 0.6),
		NewPoint(0.4, 0.4),
	)

	poly = Polygon{r1, r2}
	expected = "POLYGON((0 0,1 0,1 1,0 1,0 0),(0.4 0.4,0.6 0.4,0.6 0.6,0.4 0.6,0.4 0.4))"
	if w := poly.WKT(); w != expected {
		t.Errorf("incorrect wkt: %v", w)
	}

}