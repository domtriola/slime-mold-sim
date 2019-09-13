package simulation

import (
	"math"
)

// Velocity represents the angular direction and speed of an entity
type Velocity struct {
	direction float64
	speed     float64
}

// NextPos finds the next theoretical (non-rounded) position
func NextPos(x float64, y float64, vel Velocity) (nextX float64, nextY float64) {
	xVel, yVel := LinearVelocity(vel)
	return x + xVel, y + yVel
}

// FloatPosToGridPos rounds a float position to a discrete position
func FloatPosToGridPos(x float64, y float64) (gridX int, gridY int) {
	return int(math.Round(x)), int(math.Round(y))
}

// Radians converts a degree angle to a radian angle
func Radians(angle float64) (radianAngle float64) {
	return angle * math.Pi / 180
}

// LinearVelocity calculates linear velocity from speed * direction
func LinearVelocity(vel Velocity) (xVel float64, yVel float64) {
	xVel = math.Cos(Radians(vel.direction))
	yVel = math.Sin(Radians(vel.direction))
	return vel.speed * xVel, vel.speed * yVel
}
