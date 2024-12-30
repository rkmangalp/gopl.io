// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 41.

//!+

package tempconv

type Feet float64
type Meters float64
type Pounds float64
type Kilogram float64

// CToF converts a Celsius temperature to Fahrenheit.
func CToF(c Celsius) Fahrenheit { return Fahrenheit(c*9/5 + 32) }

// FToC converts a Fahrenheit temperature to Celsius.
func FToC(f Fahrenheit) Celsius { return Celsius((f - 32) * 5 / 9) }

//FtoM converts Feets to meters in length
func FtoM(f Feet) Meters { return Meters(f * 0.3048) }

//MtoF converts Feets to meters in length
func MtoF(m Meters) Feet { return Feet(m / 0.3048) }

//KgtoP converts kilograms to pounds in weight
func KgtoP(kg Kilogram) Pounds { return Pounds(kg * 2.20462) }

//PtoKg converts pounds to kilograms in weight
func PtoKg(p Pounds) Kilogram { return Kilogram(p * 3.28084) }

//!-
