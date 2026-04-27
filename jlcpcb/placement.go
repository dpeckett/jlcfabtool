/*
 * Copyright (C) 2026 Damian Peckett <damian@pecke.tt>.
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
	"github.com/dpeckett/jlcfabtool/kicad/placement"
)

// RotationCorrection defines how to adjust placement for a component.
type RotationCorrection struct {
	PackagePattern UnmarshallableRegexp `csv:"Package pattern"`
	ValuePattern   UnmarshallableRegexp `csv:"Value pattern"`
	Rotation       float64              `csv:"Rotation"`
	CenterX        float64              `csv:"Center X"`
	CenterY        float64              `csv:"Center Y"`
}

// specificity determines how specific a rule is (used for sorting).
func (c RotationCorrection) specificity() int {
	score := c.PackagePattern.Length

	// Value match makes it much more specific
	if c.ValuePattern.Regexp != nil {
		score += 10000 + c.ValuePattern.Length
	}

	return score
}

//go:embed kicad_rotations.csv
var rotationDBData []byte

var rotationDB []RotationCorrection

func init() {
	var err error
	rotationDB, err = csvx.Unmarshal[RotationCorrection](bytes.NewReader(rotationDBData))
	if err != nil {
		panic(err)
	}
}

// ApplyRotationCorrection applies a rotation correction based on package and optional value.
func ApplyRotationCorrection(p placement.Placement) *placement.Placement {
	slog.Info(
		"Checking for rotation correction",
		slog.String("package", p.Package),
		slog.String("value", p.Val),
	)

	var matches []RotationCorrection
	for _, correction := range rotationDB {
		// Package must match
		if correction.PackagePattern.Regexp == nil {
			continue
		}
		if !correction.PackagePattern.MatchString(p.Package) {
			continue
		}

		// If value pattern exists, it must match
		if correction.ValuePattern.Regexp != nil &&
			!correction.ValuePattern.MatchString(p.Val) {
			continue
		}

		matches = append(matches, correction)
	}

	if len(matches) == 0 {
		return &p
	}

	slog.Info("Applying rotation correction", slog.String("package", p.Package))

	// Pick most specific match
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].specificity() > matches[j].specificity()
	})

	correction := matches[0]

	// Apply center offset + rotation
	rotatedX, rotatedY := rotatePoint(
		p.PosX+correction.CenterX,
		p.PosY+correction.CenterY,
		p.PosX,
		p.PosY,
		p.Rot,
	)

	finalRotation := clampRotation(p.Rot + correction.Rotation)

	slog.Debug(
		"Rotation correction applied",
		slog.Float64("originalX", p.PosX),
		slog.Float64("originalY", p.PosY),
		slog.Float64("correctedX", rotatedX),
		slog.Float64("correctedY", rotatedY),
		slog.Float64("originalRotation", p.Rot),
		slog.Float64("finalRotation", finalRotation),
	)

	return &placement.Placement{
		Ref:     p.Ref,
		Package: p.Package,
		Val:     p.Val,
		PosX:    rotatedX,
		PosY:    rotatedY,
		Rot:     finalRotation,
		Side:    p.Side,
	}
}

// rotatePoint rotates a point (x, y) around origin (x0, y0) by theta degrees.
func rotatePoint(x, y, x0, y0, theta float64) (float64, float64) {
	thetaRad := theta * (math.Pi / 180.0)

	xPrime := x - x0
	yPrime := y - y0

	xRotated := xPrime*math.Cos(thetaRad) - yPrime*math.Sin(thetaRad)
	yRotated := xPrime*math.Sin(thetaRad) + yPrime*math.Cos(thetaRad)

	xNew := xRotated + x0
	yNew := yRotated + y0

	return xNew, yNew
}

// clampRotation clamps angle to [0, 360).
func clampRotation(angle float64) float64 {
	for angle < 0 {
		angle += 360
	}
	for angle >= 360 {
		angle -= 360
	}
	return angle
}
