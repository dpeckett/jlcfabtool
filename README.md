# jlcfabtool

A little CLI for working with JLCPCB. Yes there are plugins available for KiCad 
but I wanted something more unixy (single purpose tools that consume and emit text).

## Installation

```shell
go build --tags "fts5" 
```

## Features

- BOM conversion from KiCad CSV to JLCPCB CSV
- BOM optimization (find JLC PCBA compatible/basic parts)
- Placement conversion from KiCad CSV to JLCPCB CSV (including rotation and offset correction)