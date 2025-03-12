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

package main

import (
	_ "embed"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dpeckett/jlcfabtool/jlcpcb"
	"github.com/dpeckett/jlcfabtool/jlcpcb/partsdb"
	"github.com/dpeckett/jlcfabtool/kicad/bom"
	"github.com/dpeckett/jlcfabtool/kicad/placement"
	_ "github.com/mattn/go-sqlite3"
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
					{
						Name:      "optimize",
						Usage:     "Optimize a BOM by suggesting the best parts from the JLCPCB parts database.",
						ArgsUsage: "<bom file>",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "output",
								Aliases: []string{"o"},
								Usage:   "Path to the reccommended parts report",
								Value:   "recommended_parts.txt",
							},
						},
						Action: func(c *cli.Context) error {
							db, err := partsdb.Open()
							if err != nil {
								return err
							}
							defer db.Close()

							entries, err := bom.LoadFromCSV(c.Args().First())
							if err != nil {
								return err
							}

							reportWriter, err := os.Create(c.String("output"))
							if err != nil {
								return fmt.Errorf("error creating report file: %w", err)
							}
							defer reportWriter.Close()

							fmt.Fprintln(reportWriter, "BOM Optimization Report")
							fmt.Fprintf(reportWriter, "========================\n\n")

							for _, entry := range entries {
								slog.Info("Searching for preferred part",
									slog.String("designator", entry.Reference),
									slog.String("value", entry.Value))

								var filters []partsdb.Filter
								category, pkg := jlcpcb.FromKiCadFootprint(entry.Footprint)
								if category != "" {
									filters = append(filters, partsdb.Filter{Category: category})
								}
								if pkg != "" {
									filters = append(filters, partsdb.Filter{Package: pkg})
								}

								// Perform search for best matches
								results, err := db.Search(strconv.Quote(entry.Value)+"*", filters...)
								if err != nil {
									return err
								}

								// Log results to the report file
								fmt.Fprintf(reportWriter, "Component: %s (%s)\n", entry.Reference, entry.Value)
								fmt.Fprintf(reportWriter, "Footprint: %s\n", entry.Footprint)

								// Display the top 5 candidates
								count := 0
								for _, result := range results {
									partType := "Extended"
									if result.Basic == 1 {
										partType = "Basic"
									} else if result.Preferred == 1 {
										partType = "Preferred"
									}

									fmt.Fprintf(reportWriter, "  Candidate %d: LCSC# %d - %s - %s [%s]\n",
										count+1, result.LCSC, result.Manufacturer, result.Description, partType)

									if count++; count >= 5 {
										break
									}
								}

								fmt.Fprintln(reportWriter)
							}

							slog.Info("Report generated", slog.String("path", c.String("report")))
							return nil
						},
					},
					{
						Name:  "download-partsdb",
						Usage: "Download the JLCPCB parts database.",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "output",
								Usage: "Output file",
								Value: "jlcpcb-parts.db",
							},
						},
						Action: func(c *cli.Context) error {
							tmpDir, err := os.MkdirTemp("", "jlcpcb-partsdb")
							if err != nil {
								return fmt.Errorf("error creating temporary directory: %w", err)
							}
							defer os.RemoveAll(tmpDir)

							slog.Info("Downloading parts database, this may take a little while...")

							resp, err := http.Get(partsdb.URL)
							if err != nil {
								return fmt.Errorf("error downloading parts database: %w", err)
							}

							csvPath := filepath.Join(tmpDir, "jlcpcb-components-basic-preferred.csv")
							f, err := os.Create(csvPath)
							if err != nil {
								return fmt.Errorf("error creating temporary file: %w", err)
							}

							_, err = io.Copy(f, resp.Body)
							_ = f.Close()
							_ = resp.Body.Close()
							if err != nil {
								return fmt.Errorf("error writing temporary file: %w", err)
							}

							err = partsdb.CreateFromCSV(c.String("output"), csvPath)
							if err != nil {
								return err
							}

							return nil
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
		x := placement.PosX
		y := placement.PosY
		rotation := placement.Rot

		// Fixup differences between KiCad and JLCPCB rotations/placements.
		x, y, rotation = jlcpcb.ApplyRotationCorrection(placement.Package, x, y, rotation)

		if err := w.Write([]string{
			placement.Ref,
			strconv.FormatFloat(x, 'f', -1, 64),
			strconv.FormatFloat(y, 'f', -1, 64),
			titleCaser.String(placement.Side),
			strconv.FormatFloat(rotation, 'f', -1, 64),
		}); err != nil {
			return fmt.Errorf("error writing record: %w", err)
		}
	}

	return nil
}
