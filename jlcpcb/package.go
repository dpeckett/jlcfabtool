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
	"regexp"
	"sort"

	"github.com/dpeckett/jlcfabtool/csvx"
)

// FootprintMapping maps a KiCad footprint to an LCSC category and package.
type FootprintMapping struct {
	FootprintPattern UnmarshallableRegexp `csv:"Footprint pattern"`
	Category         string               `csv:"LCSC Category"`
	Package          string               `csv:"LCSC Package"`
}

type UnmarshallableRegexp struct {
	*regexp.Regexp
	Length int
}

func (rx *UnmarshallableRegexp) UnmarshalText(text []byte) error {
	var err error
	rx.Regexp, err = regexp.Compile(string(text))
	rx.Length = len(text)
	return err
}

//go:embed kicad_footprints.csv
var footprintDBData []byte

var footprintDB []FootprintMapping

func init() {
	// Load footprint database
	var err error
	footprintDB, err = csvx.Unmarshal[FootprintMapping](bytes.NewReader(footprintDBData))
	if err != nil {
		panic(err)
	}
}

// FromKiCadFootprint returns the LCSC category and package for a KiCad footprint.
func FromKiCadFootprint(footprint string) (category, pkg string) {
	// Remove any library prefix if present.
	if i := bytes.IndexByte([]byte(footprint), ':'); i >= 0 {
		footprint = string(footprint[i+1:])
	}

	var matches []FootprintMapping
	for _, mapping := range footprintDB {
		if mapping.FootprintPattern.MatchString(footprint) {
			matches = append(matches, mapping)
		}
	}

	// Sort by length of the pattern (descending)
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].FootprintPattern.Length > matches[j].FootprintPattern.Length
	})

	// Return the first (best) match
	if len(matches) > 0 {
		return matches[0].Category, matches[0].Package
	}

	return "", ""
}
