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

package csvx

import (
	"encoding"
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

// Unmarshal reads a CSV file from an io.Reader and unmarshals it into a slice of structs.
func Unmarshal[T any](r io.Reader) ([]T, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true
	reader.Comment = '#'

	// Read headers
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Normalize headers to lowercase for case-insensitive matching
	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[strings.ToLower(h)] = i
	}

	var results []T

	// Process rows
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV row: %w", err)
		}

		var item T
		itemVal := reflect.ValueOf(&item).Elem()
		itemType := itemVal.Type()

		// Map CSV fields to struct fields
		for i := 0; i < itemVal.NumField(); i++ {
			field := itemVal.Field(i)
			fieldType := itemType.Field(i)

			// Get column name from struct tag
			tag := fieldType.Tag.Get("csv")
			if tag == "" {
				continue // Skip fields without a CSV tag
			}

			// Lookup CSV column index
			colIdx, exists := headerMap[strings.ToLower(tag)]
			if !exists {
				continue // Skip if column is not found
			}

			rawValue := record[colIdx]

			// If field implements encoding.TextUnmarshaler, use it
			if field.CanAddr() {
				if unmarshaler, ok := field.Addr().Interface().(encoding.TextUnmarshaler); ok {
					if err := unmarshaler.UnmarshalText([]byte(rawValue)); err != nil {
						return nil, fmt.Errorf("failed to unmarshal field %s: %w", fieldType.Name, err)
					}
					continue
				}
			}

			// Is the field empty?
			if rawValue == "" {
				continue
			}

			// Convert to primitive types
			switch field.Kind() {
			case reflect.String:
				field.SetString(rawValue)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if intValue, err := strconv.ParseInt(rawValue, 10, 64); err == nil {
					field.SetInt(intValue)
				} else {
					return nil, fmt.Errorf("invalid int for field %s: %s", fieldType.Name, rawValue)
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if uintValue, err := strconv.ParseUint(rawValue, 10, 64); err == nil {
					field.SetUint(uintValue)
				} else {
					return nil, fmt.Errorf("invalid uint for field %s: %s", fieldType.Name, rawValue)
				}
			case reflect.Float32, reflect.Float64:
				if floatValue, err := strconv.ParseFloat(rawValue, 64); err == nil {
					field.SetFloat(floatValue)
				} else {
					return nil, fmt.Errorf("invalid float for field %s: %s", fieldType.Name, rawValue)
				}
			case reflect.Bool:
				if boolValue, err := strconv.ParseBool(strings.ToLower(rawValue)); err == nil {
					field.SetBool(boolValue)
				} else {
					return nil, fmt.Errorf("invalid bool for field %s: %s", fieldType.Name, rawValue)
				}
			default:
				return nil, fmt.Errorf("unsupported field type: %s", fieldType.Name)
			}
		}

		results = append(results, item)
	}

	return results, nil
}
