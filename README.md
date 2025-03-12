# jlcfabtool

A little CLI for working with JLCPCB. Yes there are plugins available for KiCad but I wanted something more unixy (eg. single purpose tools that consume and emit text).

## Installation

```shell
go build --tags "fts5" 
```

## Features

- BOM conversion from KiCad CSV to JLCPCB CSV
- BOM optimization (find JLC PCBA compatible/basic parts)
- Placement conversion from KiCad CSV to JLCPCB CSV (including rotation and offset correction)

## Usage

### Convert a BOM from KiCad to JLCPCB

To convert a BOM from KiCad to JLCPCB, you need to export the BOM from KiCad as a CSV file. 
Then you can run the following command:

```shell
./jlcfabtool bom convert kicad-bom.csv
```

A new file `kicad-bom.jlcpcb.csv` will be created in the same directory as `kicad-bom.csv`.

### Optimize a BOM for JLCPCB

To optimize a BOM for JLCPCB (provide suggestions for basic and recommended parts), 
you need to export the BOM from KiCad as a CSV file. Then you can run the following command:

```shell
./jlcfabtool bom optimize kicad.csv
```

A text report `recommended_parts.txt` will be created in the current directory.

### Convert Component Placements from KiCad to JLCPCB

To convert component placements from KiCad to JLCPCB, you need to export the 
component placements from KiCad as a CSV file. Then you can run the following command:

```shell
./jlcfabtool placement convert kicad-all-pos.csv
```

A new file `kicad-all-pos.jlcpcb.csv` will be created in the same directory as `kicad-all-pos.csv`.