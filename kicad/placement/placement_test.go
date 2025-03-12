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

package placement_test

import (
	"testing"

	"github.com/dpeckett/jlcfabtool/kicad/placement"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFromCSV(t *testing.T) {
	placements, err := placement.LoadFromCSV("testdata/placements.csv")

	require.NoError(t, err)
	require.Len(t, placements, 6)

	// First placement.
	assert.Equal(t, "C1", placements[0].Ref)
	assert.Equal(t, "100n", placements[0].Val)
	assert.Equal(t, "C_0603_1608Metric", placements[0].Package)
	assert.Equal(t, 28.194, placements[0].PosX)
	assert.Equal(t, -173.26, placements[0].PosY)
	assert.Equal(t, 0.0, placements[0].Rot)
	assert.Equal(t, "top", placements[0].Side)

	// Last placement.
	assert.Equal(t, "X1", placements[5].Ref)
	assert.Equal(t, "50M", placements[5].Val)
	assert.Equal(t, "Oscillator_SMD_SeikoEpson_SG8002LB-4Pin_5.0x3.2mm", placements[5].Package)
	assert.Equal(t, 35.263, placements[5].PosX)
	assert.Equal(t, -158.02, placements[5].PosY)
	assert.Equal(t, 90.0, placements[5].Rot)
	assert.Equal(t, "top", placements[5].Side)
}
