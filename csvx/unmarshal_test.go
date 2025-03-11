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

package csvx_test

import (
	"strings"
	"testing"
	"time"

	"github.com/dpeckett/jlcfabtool/csvx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Define a struct using custom types
type Person struct {
	Name      string    `csv:"name"`
	Age       int       `csv:"age"`
	Birthdate time.Time `csv:"birthdate"`
	Active    bool      `csv:"active"`
}

func TestUnmarshalCSV(t *testing.T) {
	// Sample CSV input
	csvData := `Name,Age,Birthdate,Active
John Doe,30,"1994-06-15T00:00:00Z",true
Jane Smith,25,"1998-09-10T00:00:00Z",false
`

	reader := strings.NewReader(csvData)

	people, err := csvx.UnmarshalCSV[Person](reader)

	require.NoError(t, err)
	require.Len(t, people, 2)

	// First record
	assert.Equal(t, "John Doe", people[0].Name)
	assert.Equal(t, 30, people[0].Age)
	assert.Equal(t, time.Date(1994, 6, 15, 0, 0, 0, 0, time.UTC), people[0].Birthdate)
	assert.Equal(t, true, people[0].Active)

	// Second record
	assert.Equal(t, "Jane Smith", people[1].Name)
	assert.Equal(t, 25, people[1].Age)
	assert.Equal(t, time.Date(1998, 9, 10, 0, 0, 0, 0, time.UTC), people[1].Birthdate)
	assert.Equal(t, false, people[1].Active)
}
