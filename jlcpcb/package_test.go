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

func TestFromKiCadFootprint(t *testing.T) {
	// No library prefix.
	category, pkg := jlcpcb.FromKiCadFootprint("Oscillator_SMD_SeikoEpson_SG8002LB-4Pin_5.0x3.2mm")

	assert.Equal(t, "Crystals/Oscillators/Resonators", category)
	assert.Equal(t, "SMD5032", pkg)

	// With library prefix.
	category, pkg = jlcpcb.FromKiCadFootprint("Capacitor_SMD:C_0603_1608Metric")

	assert.Equal(t, "Capacitors", category)
	assert.Equal(t, "0603", pkg)

	// Unknown footprint.
	category, pkg = jlcpcb.FromKiCadFootprint("Unknown_Footprint")

	assert.Equal(t, "", category)
	assert.Equal(t, "", pkg)
}
