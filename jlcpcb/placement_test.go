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

package jlcpcb_test

import (
	"testing"

	"github.com/dpeckett/jlcfabtool/jlcpcb"
	"github.com/dpeckett/jlcfabtool/kicad/placement"
	"github.com/stretchr/testify/assert"
)

func TestApplyRotationCorrection(t *testing.T) {
	p := jlcpcb.ApplyRotationCorrection(placement.Placement{
		Ref:     "J1",
		Val:     "",
		Package: "IDC-Header_2x10_P2.54mm_Vertical",
		PosX:    30.988,
		PosY:    -89.44,
		Rot:     90.0,
		Side:    "top",
	})

	assert.Equal(t, "J1", p.Ref)
	assert.Equal(t, "", p.Val)
	assert.Equal(t, "IDC-Header_2x10_P2.54mm_Vertical", p.Package)
	assert.InDelta(t, 42.538, p.PosX, 0.000001)
	assert.InDelta(t, -88.19, p.PosY, 0.000001)
	assert.InDelta(t, 0.0, p.Rot, 0.000001)
	assert.Equal(t, "top", p.Side)

	p = jlcpcb.ApplyRotationCorrection(placement.Placement{
		Ref:     "U1",
		Val:     "ADS131M02IRUKR",
		Package: "WQFN-20-1EP_3x3mm_P0.4mm_EP1.7x1.7mm",
		PosX:    77.9817,
		PosY:    -174.2624,
		Rot:     0.0,
		Side:    "top",
	})

	assert.Equal(t, "U1", p.Ref)
	assert.Equal(t, "ADS131M02IRUKR", p.Val)
	assert.Equal(t, "WQFN-20-1EP_3x3mm_P0.4mm_EP1.7x1.7mm", p.Package)
	assert.InDelta(t, 77.9817, p.PosX, 0.000001)
	assert.InDelta(t, -174.2624, p.PosY, 0.000001)
	assert.InDelta(t, 90.0, p.Rot, 0.000001)
	assert.Equal(t, "top", p.Side)
}

func TestApplyRotationCorrectionSkipsValueSpecificMismatch(t *testing.T) {
	p := jlcpcb.ApplyRotationCorrection(placement.Placement{
		Ref:     "U1",
		Val:     "SOME_OTHER_PART",
		Package: "WQFN-20-1EP_3x3mm_P0.4mm_EP1.7x1.7mm",
		PosX:    77.9817,
		PosY:    -174.2624,
		Rot:     0.0,
		Side:    "top",
	})

	assert.Equal(t, "U1", p.Ref)
	assert.Equal(t, "SOME_OTHER_PART", p.Val)
	assert.Equal(t, "WQFN-20-1EP_3x3mm_P0.4mm_EP1.7x1.7mm", p.Package)
	assert.InDelta(t, 77.9817, p.PosX, 0.000001)
	assert.InDelta(t, -174.2624, p.PosY, 0.000001)
	assert.InDelta(t, 0.0, p.Rot, 0.000001)
	assert.Equal(t, "top", p.Side)
}

func TestApplyRotationCorrectionNoMatch(t *testing.T) {
	p := jlcpcb.ApplyRotationCorrection(placement.Placement{
		Ref:     "U99",
		Val:     "Unknown_Value",
		Package: "Unknown_Package",
		PosX:    10.0,
		PosY:    20.0,
		Rot:     270.0,
		Side:    "bottom",
	})

	assert.Equal(t, "U99", p.Ref)
	assert.Equal(t, "Unknown_Value", p.Val)
	assert.Equal(t, "Unknown_Package", p.Package)
	assert.InDelta(t, 10.0, p.PosX, 0.000001)
	assert.InDelta(t, 20.0, p.PosY, 0.000001)
	assert.InDelta(t, 270.0, p.Rot, 0.000001)
	assert.Equal(t, "bottom", p.Side)
}
