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

import "regexp"

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
