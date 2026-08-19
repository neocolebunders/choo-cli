# trains

Tiny CLI for the next 3 Belgian trains between two stations, using the [iRail API](https://docs.irail.be). No dependencies, pure Go stdlib.

```
$ trains nord diest
13:52  plat 6    → Diest  arr 14:34  direct
14:08 +2 plat 3  → Diest  arr 15:09 +3  1 transfer(s)
14:52  plat 6    → Diest  arr 15:34  direct
```

- Green departure time, red `+N` delay in minutes (after departure and/or arrival time), red `CANCELED`, yellow transfer count.
- Colors auto-disable when output is piped.
- Station names are fuzzy-matched by iRail (`leuven`, `brussel noord`, ...); `nord`/`noord` are local aliases for Brussel-Noord.

## Install

```
go install ./cmd/trains
```

Make sure `~/go/bin` is on your `PATH`.

## Usage

```
trains <from> <to>
```
