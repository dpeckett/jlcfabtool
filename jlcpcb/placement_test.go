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

package jlcpcb_test

import (
	"testing"

	"github.com/dpeckett/jlcfabtool/jlcpcb"
	"github.com/stretchr/testify/assert"
)

func TestApplyRotationCorrection(t *testing.T) {
	// A component with a rotation and position correction.
	x, y, rotation := jlcpcb.ApplyRotationCorrection("IDC-Header_2x10_P2.54mm_Vertical", 30.988, -89.44, 90.0)

	assert.Equal(t, 42.538, x)
	assert.Equal(t, -88.19, y)
	assert.Equal(t, 0.0, rotation)

	// Just a rotation correction.
	x, y, rotation = jlcpcb.ApplyRotationCorrection("SOT-223-3_TabPin2", 55.042, -135.89, 0.0)

	assert.Equal(t, 55.042, x)
	assert.Equal(t, -135.89, y)
	assert.Equal(t, 180.0, rotation)
}
