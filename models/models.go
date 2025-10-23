// Package models contains basic domain models
package models

import (
	"fmt"
)

type Passage struct {
	Track      []TPoint `json:"track"`
	LicenseNum string   `json:"licenseNumber"`

	Speeds  []float64      `json:"speeds"`
	Classes []VehicleClass `json:"classes"`
	Sides   []VehicleSide  `json:"sides"`
}
// не статичный класс ТС из-за помех? Мол, класс ТС может по-разному определяться на разных участках?

type TPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	T int     `json:"t"`
}

type VehicleClass int8

const (
	UndefinedClass VehicleClass = iota - 1
	Car
	Moto
	Bus
	Truck
)

type VehicleSide int8

const (
	UndefinedSide VehicleSide = iota - 1
	Front
	Read
)

type Offender struct {
	GRN string
}

type FatalError struct {
	Cause string
}

func (f *FatalError) Error() string {
	return fmt.Sprintf("fatal error: %s", f.Cause)
}
