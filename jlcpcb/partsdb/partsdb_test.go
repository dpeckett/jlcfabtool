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

package partsdb_test

import (
	"path/filepath"
	"testing"

	"github.com/dpeckett/jlcfabtool/jlcpcb/partsdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPartsDB(t *testing.T) {
	tempDir := t.TempDir()

	dbPath := filepath.Join(tempDir, "jlcpcb-parts.db")

	err := partsdb.CreateFromCSV(dbPath, "testdata/jlcpcb-components-basic-preferred.csv")
	require.NoError(t, err)

	db, err := partsdb.Open(dbPath)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})

	parts, err := db.Search("100nF", partsdb.Filter{Package: "0805"})
	require.NoError(t, err)

	assert.Len(t, parts, 2)
	assert.Equal(t, 28233, parts[0].LCSC)
	assert.Equal(t, 49678, parts[1].LCSC)
}
