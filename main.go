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

package main

import (
	_ "embed"
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/dpeckett/jlcfabtool/jlcpcb"
	"github.com/dpeckett/jlcfabtool/kicad/bom"
	"github.com/dpeckett/jlcfabtool/kicad/placement"
	"github.com/urfave/cli/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {
	app := &cli.App{
		Name:  "jlcfabtool",
		Usage: "A little CLI for working with JLCPCB.",
		Commands: []*cli.Command{
			{
				Name:  "bom",
				Usage: "Commands for working with BOMs.",
				Subcommands: []*cli.Command{
					{
						Name:      "convert",
						Usage:     "Convert a KiCad BOM into JLCPCB format.",
						ArgsUsage: "<file>",
						Action: func(c *cli.Context) error {
							return convertKiCadBOM(c.Args().First())
						},
					},
				},
			},
			{
				Name:  "placement",
				Usage: "Commands for working with component placements (CPL).",
				Subcommands: []*cli.Command{
					{
						Name:      "convert",
						Usage:     "Convert a KiCad component placements (CPL) into JLCPCB format.",
						ArgsUsage: "<file>",
						Action: func(c *cli.Context) error {
							return convertKiCadComponentPlacements(c.Args().First())
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("Error running app", slog.Any("error", err))
		os.Exit(1)
	}
}

func convertKiCadBOM(file string) error {
	slog.Info("Converting BOM", slog.Any("file", file))

	entries, err := bom.LoadFromCSV(file)
	if err != nil {
		return fmt.Errorf("error loading BOM: %w", err)
	}

	f, err := os.Create(strings.TrimSuffix(file, ".csv") + ".jlcpcb.csv")
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{"Comment", "Designator", "Footprint", "LCSC Part Number"}); err != nil {
		return fmt.Errorf("error writing header: %w", err)
	}

	for _, entry := range entries {
		if err := w.Write([]string{
			entry.Value,
			entry.Reference,
			entry.Footprint,
			entry.LCSC,
		}); err != nil {
			return fmt.Errorf("error writing record: %w", err)
		}
	}

	return nil
}

func convertKiCadComponentPlacements(file string) error {
	slog.Info("Converting component placement", slog.Any("file", file))

	placements, err := placement.LoadFromCSV(file)
	if err != nil {
		return fmt.Errorf("error loading component placements: %w", err)
	}

	f, err := os.Create(strings.TrimSuffix(file, ".csv") + ".jlcpcb.csv")
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{"Designator", "Mid X", "Mid Y", "Layer", "Rotation"}); err != nil {
		return fmt.Errorf("error writing header: %w", err)
	}

	titleCaser := cases.Title(language.English, cases.Compact)

	for _, placement := range placements {
		// Fixup differences between KiCad and JLCPCB rotations/placements.
		placement := jlcpcb.ApplyRotationCorrection(placement)

		if err := w.Write([]string{
			placement.Ref,
			strconv.FormatFloat(placement.PosX, 'f', -1, 64),
			strconv.FormatFloat(placement.PosY, 'f', -1, 64),
			titleCaser.String(placement.Side),
			strconv.FormatFloat(placement.Rot, 'f', -1, 64),
		}); err != nil {
			return fmt.Errorf("error writing record: %w", err)
		}
	}

	return nil
}
