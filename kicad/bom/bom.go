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

package bom

import (
	"fmt"
	"os"

	"github.com/dpeckett/jlcfabtool/csvx"
)

// Entry represents a single row in a KiCad BOM CSV file.
type Entry struct {
	Reference string `csv:"Reference"`
	Value     string `csv:"Value"`
	Footprint string `csv:"Footprint"`
	Qty       int    `csv:"Qty"`
	LCSC      string `csv:"LCSC PN"`
}

// LoadFromCSV loads a KiCad BOM from a CSV file.
func LoadFromCSV(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer f.Close()

	entries, err := csvx.Unmarshal[Entry](f)
	if err != nil {
		return nil, fmt.Errorf("could not parse CSV: %w", err)
	}

	return entries, nil
}
