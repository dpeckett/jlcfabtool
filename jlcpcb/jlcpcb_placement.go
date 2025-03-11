/*
 * Copyright (C) 2025 Damian Peckett <damian@pecke.tt>.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package jlcpcb

import (
	"bytes"
	_ "embed"
	"log/slog"
	"math"
	"sort"

	"github.com/dpeckett/jlcfabtool/csvx"
)

type RotationCorrection struct {
	FootprintPattern UnmarshallableRegexp `csv:"Footprint pattern"`
	// Rotation is the number of degrees that the component has been rotated
	// relative to the JLCPCB canonical orientation.
	Rotation float64 `csv:"Rotation"`
	// CenterX and CenterY are the coordinates of the center of the component
	// relative to KiCad's origin. JLCPCB always uses the center/midpoint of
	// the component for placement.
	CenterX float64 `csv:"Center X"`
	CenterY float64 `csv:"Center Y"`
}

//go:embed kicad_rotations.csv
var rotationDBData []byte

var rotationDB []RotationCorrection

func init() {
	// Load rotation database
	var err error
	rotationDB, err = csvx.UnmarshalCSV[RotationCorrection](bytes.NewReader(rotationDBData))
	if err != nil {
		panic(err)
	}
}

// ApplyRotationCorrection applies a rotation correction to a component's position and rotation.
func ApplyRotationCorrection(pkg string, x, y, rotation float64) (float64, float64, float64) {
	var matches []RotationCorrection
	for _, correction := range rotationDB {
		if correction.FootprintPattern.MatchString(pkg) {
			matches = append(matches, correction)
		}
	}

	if len(matches) > 0 {
		slog.Info("Applying rotation correction", slog.String("package", pkg))

		// Find the best matching pattern
		sort.Slice(matches, func(i, j int) bool {
			return matches[i].FootprintPattern.Length > matches[j].FootprintPattern.Length
		})
		correction := matches[0]

		// Calculate the corrected position and rotation
		rotatedX, rotatedY := rotatePoint(x+correction.CenterX, y+correction.CenterY, x, y, rotation)
		finalRotation := clampRotation(rotation + correction.Rotation)

		slog.Debug("Rotation correction applied",
			slog.Float64("originalX", x), slog.Float64("originalY", y),
			slog.Float64("correctedX", rotatedX), slog.Float64("correctedY", rotatedY),
			slog.Float64("originalRotation", rotation), slog.Float64("finalRotation", finalRotation))

		return rotatedX, rotatedY, finalRotation
	}

	// No correction found, return original values
	return x, y, rotation
}

// rotatePoint rotates a point (x, y) around an origin (x0, y0) by theta degrees.
func rotatePoint(x, y, x0, y0, theta float64) (float64, float64) {
	// Convert theta to radians
	thetaRad := theta * (math.Pi / 180.0)

	// Translate point to origin
	xPrime := x - x0
	yPrime := y - y0

	// Apply rotation transformation
	xRotated := xPrime*math.Cos(thetaRad) - yPrime*math.Sin(thetaRad)
	yRotated := xPrime*math.Sin(thetaRad) + yPrime*math.Cos(thetaRad)

	// Translate back
	xNew := xRotated + x0
	yNew := yRotated + y0

	return xNew, yNew
}

// clampRotation clamps an angle to the range [0, 360).
func clampRotation(angle float64) float64 {
	for angle < 0 {
		angle += 360
	}

	for angle >= 360 {
		angle -= 360
	}

	return angle
}
