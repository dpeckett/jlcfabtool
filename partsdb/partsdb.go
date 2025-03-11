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

package partsdb

import (
	"fmt"
	"os"

	"github.com/dpeckett/jlcfabtool/csvx"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Where to retrieve the CSV file from.
const URL = "https://cdfer.github.io/jlcpcb-parts-database/jlcpcb-components-basic-preferred.csv"

// Part represents a JLCPCB part.
type Part struct {
	LCSC            int    `csv:"lcsc" db:"lcsc"`
	CategoryID      int    `csv:"category_id" db:"category_id"`
	Category        string `csv:"category" db:"category"`
	Subcategory     string `csv:"subcategory" db:"subcategory"`
	Mfr             string `csv:"mfr" db:"mfr"`
	Package         string `csv:"package" db:"package"`
	Joints          int    `csv:"joints" db:"joints"`
	Manufacturer    string `csv:"manufacturer" db:"manufacturer"`
	Basic           int    `csv:"basic" db:"basic"`
	Preferred       int    `csv:"preferred" db:"preferred"`
	Description     string `csv:"description" db:"description"`
	Datasheet       string `csv:"datasheet" db:"datasheet"`
	Stock           int    `csv:"stock" db:"stock"`
	LastOnStock     int    `csv:"last_on_stock" db:"last_on_stock"`
	Price           string `csv:"price" db:"price"`
	Extra           string `csv:"extra" db:"extra"`
	AssemblyProcess string `csv:"Assembly Process" db:"assembly_process"`
	MinOrderQty     int    `csv:"Min Order Qty" db:"min_order_qty"`
	AttritionQty    int    `csv:"Attrition Qty" db:"attrition_qty"`
}

type PartsDB struct {
	db *sqlx.DB
}

// CreateFromCSV creates a new PartsDB from a CSV file.
func CreateFromCSV(dbPath, csvPath string) (*PartsDB, error) {
	db, err := sqlx.Connect("sqlite3", dbPath+"?_journal_mode=WAL&_synchronous=NORMAL")
	if err != nil {
		return nil, fmt.Errorf("could not open database: %w", err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS parts (
		lcsc INTEGER PRIMARY KEY,
		category_id INTEGER,
		category TEXT,
		subcategory TEXT,
		mfr TEXT,
		package TEXT,
		joints INTEGER,
		manufacturer TEXT,
		basic INTEGER,
		preferred INTEGER,
		description TEXT,
		datasheet TEXT,
		stock INTEGER,
		last_on_stock INTEGER,
		price TEXT,
		extra TEXT,
		assembly_process TEXT,
		min_order_qty INTEGER,
		attrition_qty INTEGER
	);
	CREATE VIRTUAL TABLE IF NOT EXISTS parts_fts USING fts5(
		lcsc UNINDEXED, description, mfr, manufacturer, category, subcategory
	);
	CREATE TRIGGER IF NOT EXISTS parts_ai AFTER INSERT ON parts BEGIN
		INSERT INTO parts_fts(lcsc, description, mfr, manufacturer, category, subcategory)
		VALUES (NEW.lcsc, NEW.description, NEW.mfr, NEW.manufacturer, NEW.category, NEW.subcategory);
	END;
	CREATE TRIGGER IF NOT EXISTS parts_ad AFTER DELETE ON parts BEGIN
		DELETE FROM parts_fts WHERE lcsc = OLD.lcsc;
	END;
	CREATE TRIGGER IF NOT EXISTS parts_au AFTER UPDATE ON parts BEGIN
		DELETE FROM parts_fts WHERE lcsc = OLD.lcsc;
		INSERT INTO parts_fts(lcsc, description, mfr, manufacturer, category, subcategory)
		VALUES (NEW.lcsc, NEW.description, NEW.mfr, NEW.manufacturer, NEW.category, NEW.subcategory);
	END;
	`

	if _, err = db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("could not create schema: %w", err)
	}

	// Load CSV data
	f, err := os.Open(csvPath)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("could not open CSV file: %w", err)
	}
	defer f.Close()

	parts, err := csvx.UnmarshalCSV[Part](f)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("could not parse CSV: %w", err)
	}

	// Insert into database using sqlx
	insertSQL := `INSERT OR IGNORE INTO parts (
		lcsc, category_id, category, subcategory, mfr, package, joints, manufacturer,
		basic, preferred, description, datasheet, stock, last_on_stock, price,
		extra, assembly_process, min_order_qty, attrition_qty
	) VALUES (:lcsc, :category_id, :category, :subcategory, :mfr, :package, :joints, :manufacturer,
		:basic, :preferred, :description, :datasheet, :stock, :last_on_stock, :price,
		:extra, :assembly_process, :min_order_qty, :attrition_qty)`

	_, err = db.NamedExec(insertSQL, parts)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("could not insert records: %w", err)
	}

	return &PartsDB{db: db}, nil
}

// Open opens an existing PartsDB.
func Open(dbPath string) (*PartsDB, error) {
	db, err := sqlx.Connect("sqlite3", dbPath+"?_journal_mode=WAL&_synchronous=NORMAL")
	if err != nil {
		return nil, fmt.Errorf("could not open database: %w", err)
	}
	return &PartsDB{db: db}, nil
}

func (pdb *PartsDB) Close() error {
	return pdb.db.Close()
}

// Filter represents a filter to apply to the search query.
type Filter struct {
	Category string
	Package  string
}

// Search searches the database for parts matching the query and filters.
func (pdb *PartsDB) Search(query string, filters ...Filter) ([]Part, error) {
	args := []any{query}
	selectSQL := `SELECT p.* FROM parts p JOIN parts_fts fts ON p.lcsc = fts.lcsc WHERE parts_fts MATCH ?`

	for _, filter := range filters {
		if filter.Category != "" {
			selectSQL += ` AND p.category = ?`
			args = append(args, filter.Category)
		}

		if filter.Package != "" {
			selectSQL += ` AND p.package = ?`
			args = append(args, filter.Package)
		}
	}

	// Order by basic parts first, then by preferred parts
	selectSQL += ` ORDER BY p.basic DESC, p.preferred DESC`

	var parts []Part
	if err := pdb.db.Select(&parts, selectSQL, args...); err != nil {
		return nil, fmt.Errorf("could not search database: %w", err)
	}

	return parts, nil
}
