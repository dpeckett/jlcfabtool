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

package placement

import (
	"fmt"
	"os"

	"github.com/dpeckett/jlcfabtool/csvx"
)

// Placement represents a component placement.
type Placement struct {
	Ref     string  `csv:"Ref"`
	Val     string  `csv:"Val"`
	Package string  `csv:"Package"`
	PosX    float64 `csv:"PosX"`
	PosY    float64 `csv:"PosY"`
	Rot     float64 `csv:"Rot"`
	Side    string  `csv:"Side"`
}

// LoadFromCSV loads KiCad component placements from a CSV file.
func LoadFromCSV(path string) ([]Placement, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer f.Close()

	placements, err := csvx.Unmarshal[Placement](f)
	if err != nil {
		return nil, fmt.Errorf("could not parse CSV: %w", err)
	}

	return placements, nil
}
