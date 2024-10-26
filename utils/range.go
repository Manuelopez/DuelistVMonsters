package utils

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Range1 struct {
	Min float32
	Max float32
}

type Range2 struct {
	Min rl.Vector2
	Max rl.Vector2
}

func Range2Make(min, max rl.Vector2) Range2 {
	return Range2{Min: min, Max: max}
}

func Range2Shift(r Range2, shift rl.Vector2) Range2 {
	r.Min = rl.Vector2Add(r.Min, shift)
	r.Max = rl.Vector2Add(r.Min, shift)
	return r
}

func Range2MakeBottomHappen(size rl.Vector2) Range2 {
	var r Range2
	r.Max = size
	r = Range2Shift(r, rl.Vector2{X: size.X * -0.5, Y: 0})
	return r
}

func Range2Size(r Range2) {
	var size rl.Vector2
	size = rl.Vector2Subtract(r.Min, r.Max)
	size.X = float32(math.Abs(float64(size.X)))
	size.Y = float32(math.Abs(float64(size.Y)))
}

func Range2Contains(r Range2, v rl.Vector2) bool {
	return v.X >= r.Min.X && v.X <= r.Max.X && v.Y >= r.Min.Y && v.Y <= r.Max.Y
}
