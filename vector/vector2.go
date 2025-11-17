package vector

import (
	"fmt"
	"gamelib-go/number"
	"math"
)

type Vector2[T number.Number] struct {
	X T
	Y T
}

func NewVector2[T number.Number](x, y T) Vector2[T] {
	return Vector2[T]{X: x, Y: y}
}

func (v Vector2[T]) Clone() Vector2[T] {
	return Vector2[T]{X: v.X, Y: v.Y}
}

func (v Vector2[T]) Add(other Vector2[T]) Vector2[T] {
	return Vector2[T]{X: v.X + other.X, Y: v.Y + other.Y}
}

func (v Vector2[T]) Sub(other Vector2[T]) Vector2[T] {
	return Vector2[T]{X: v.X - other.X, Y: v.Y - other.Y}
}

func (v Vector2[T]) Scale(scale Vector2[T]) Vector2[T] {
	return Vector2[T]{X: v.X * scale.X, Y: v.Y * scale.Y}
}

func (v Vector2[T]) Dot(rhs Vector2[T]) float64 {
	return float64(v.X*rhs.X + v.Y*rhs.Y)
}

func (v Vector2[T]) SqrMagnitude() float64 {
	return float64(v.X*v.X + v.Y*v.Y)
}

func (v Vector2[T]) Magnitude() float64 {
	return math.Sqrt(v.SqrMagnitude())
}

func (v Vector2[T]) Normalized() Vector2[T] {
	mag := v.Magnitude()
	if mag < 1e-15 {
		return Vector2[T]{0, 0}
	}
	return Vector2[T]{X: T(float64(v.X) / mag), Y: T(float64(v.Y) / mag)}
}

func (v Vector2[T]) Distance(other Vector2[T]) float64 {
	dx := float64(v.X - other.X)
	dy := float64(v.Y - other.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

func (v Vector2[T]) Equals(other Vector2[T]) bool {
	return v.X == other.X && v.Y == other.Y
}

func (v Vector2[T]) String() string {
	return fmt.Sprintf("X: %v Y: %v", v.X, v.Y)
}

func Lerp[T number.Number](a, b Vector2[T], t float64) Vector2[float64] {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	return Vector2[float64]{
		X: float64(a.X) + (float64(b.X)-float64(a.X))*t,
		Y: float64(a.Y) + (float64(b.Y)-float64(a.Y))*t,
	}
}

func ClampMagnitude[T number.Number](vector Vector2[T], maxLength float64) Vector2[float64] {
	sqrMag := vector.SqrMagnitude()
	if sqrMag > maxLength*maxLength {
		mag := math.Sqrt(sqrMag)
		scale := maxLength / mag
		return Vector2[float64]{X: float64(vector.X) * scale, Y: float64(vector.Y) * scale}
	}
	return Vector2[float64]{X: float64(vector.X), Y: float64(vector.Y)}
}

func Angle[T number.Number](from, to Vector2[T]) float64 {
	m1 := from.Magnitude()
	m2 := to.Magnitude()
	if m1 < 1e-15 || m2 < 1e-15 {
		return 0
	}
	cos := from.Dot(to) / (m1 * m2)
	if cos > 1 {
		cos = 1
	} else if cos < -1 {
		cos = -1
	}
	return math.Acos(cos) * 180.0 / math.Pi
}

func SignedAngle[T number.Number](from, to Vector2[T]) float64 {
	angle := Angle(from, to)
	cross := float64(from.X)*float64(to.Y) - float64(from.Y)*float64(to.X)
	if cross >= 0 {
		return angle
	}
	return -angle
}

func Min[T number.Number](a, b Vector2[T]) Vector2[T] {
	return Vector2[T]{X: min(a.X, b.X), Y: min(a.Y, b.Y)}
}

func Max[T number.Number](a, b Vector2[T]) Vector2[T] {
	return Vector2[T]{X: max(a.X, b.X), Y: max(a.Y, b.Y)}
}

func InFOV[T number.Number](forward, dir Vector2[T], fov float64) bool {
	return math.Abs(SignedAngle(forward, dir)) <= fov/2
}

func InFOVDistance[T number.Number](currentPos, targetPos, forward Vector2[T], fov, maxDist float64) bool {
	if math.Abs(SignedAngle(forward, targetPos.Sub(currentPos).Normalized())) > fov/2 {
		return false
	}
	if float64(currentPos.Distance(targetPos)) > maxDist {
		return false
	}
	return true
}

func MoveTowards[T number.Number](current, target Vector2[T], maxDistanceDelta float64) Vector2[T] {
	dx := float64(target.X - current.X)
	dy := float64(target.Y - current.Y)
	sqrDist := dx*dx + dy*dy

	if sqrDist == 0 || sqrDist <= maxDistanceDelta*maxDistanceDelta {
		return Vector2[T]{X: target.X, Y: target.Y}
	}

	dist := math.Sqrt(sqrDist)
	scale := maxDistanceDelta / dist
	return Vector2[T]{X: current.X + T(dx*scale), Y: current.Y + T(dy*scale)}
}

func Reflect[T number.Number](inDirection, inNormal Vector2[T]) Vector2[T] {
	dot := inNormal.Dot(inDirection)
	return Vector2[T]{X: T(float64(inDirection.X) - 2*dot*float64(inNormal.X)), Y: T(float64(inDirection.Y) - 2*dot*float64(inNormal.Y))}
}
