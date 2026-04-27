# jlcfabtool

A little CLI for working with JLCPCB. Yes there are plugins available for KiCad but I wanted something more unixy (eg. single purpose tools that consume and emit text).

## Installation

```shell
go build
```

## Features

- BOM conversion from KiCad CSV to JLCPCB CSV
- Placement conversion from KiCad CSV to JLCPCB CSV (including rotation and offset correction)

## Usage

### Convert a BOM from KiCad to JLCPCB

To convert a BOM from KiCad to JLCPCB, you need to export the BOM from KiCad as a CSV file. 
Then you can run the following command:

```shell
./jlcfabtool bom convert kicad-bom.csv
```

A new file `kicad-bom.jlcpcb.csv` will be created in the same directory as `kicad-bom.csv`.

### Convert Component Placements from KiCad to JLCPCB

To convert component placements from KiCad to JLCPCB, you need to export the 
component placements from KiCad as a CSV file. Then you can run the following command:

```shell
./jlcfabtool placement convert kicad-all-pos.csv
```

A new file `kicad-all-pos.jlcpcb.csv` will be created in the same directory as `kicad-all-pos.csv`.