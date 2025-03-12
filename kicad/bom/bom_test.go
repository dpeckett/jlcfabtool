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

package bom_test

import (
	"testing"

	"github.com/dpeckett/jlcfabtool/kicad/bom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFromCSV(t *testing.T) {
	entries, err := bom.LoadFromCSV("testdata/bom.csv")

	require.NoError(t, err)
	require.Len(t, entries, 6)

	// First record
	assert.Equal(t, "C1,C2,C3,C4,C7,C10", entries[0].Reference)
	assert.Equal(t, "100n", entries[0].Value)
	assert.Equal(t, "Capacitor_SMD:C_0603_1608Metric", entries[0].Footprint)
	assert.Equal(t, 6, entries[0].Qty)
	assert.Equal(t, "C14663", entries[0].LCSC)

	// Last record
	assert.Equal(t, "J2", entries[5].Reference)
	assert.Equal(t, "Boot Selection", entries[5].Value)
	assert.Equal(t, "Connector_PinHeader_2.54mm:PinHeader_1x03_P2.54mm_Vertical", entries[5].Footprint)
	assert.Equal(t, 1, entries[5].Qty)
	assert.Empty(t, entries[5].LCSC)
}
