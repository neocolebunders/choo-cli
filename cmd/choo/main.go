package main

import (
	"cmp"
	"fmt"
	"os"
	"time"

	"choo-cli/internal/irail"
)

// iRail's fuzzy matcher handles most names but rejects "nord"/"noord" on their own.
var aliases = map[string]string{
	"nord":  "Brussel-Noord",
	"noord": "Brussel-Noord",
}

func station(name string) string {
	return cmp.Or(aliases[name], name)
}

var green, red, yellow, bold, dim, reset = "\033[32m", "\033[31m", "\033[33m", "\033[1m", "\033[2m", "\033[0m"

func init() {
	// disable colors when output is piped or redirected
	if fi, err := os.Stdout.Stat(); err != nil || fi.Mode()&os.ModeCharDevice == 0 {
		green, red, yellow, bold, dim, reset = "", "", "", "", "", ""
	}
}

func delayTag(d time.Duration) string {
	if d > 0 {
		mins := int((d + time.Minute - 1) / time.Minute) // ceil: 30s → +1, not +0
		return fmt.Sprintf(" %s+%d%s", red, mins, reset)
	}
	return ""
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: choo <from> <to>   e.g. trains nord diest")
		os.Exit(1)
	}

	conns, err := irail.Connections(station(os.Args[1]), station(os.Args[2]))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, c := range conns[:min(3, len(conns))] {
		note := dim + "direct" + reset
		if c.Transfers > 0 {
			note = fmt.Sprintf("%s%d transfer(s)%s", yellow, c.Transfers, reset)
		}
		if c.Departure.Canceled {
			note = red + "CANCELED" + reset
		}
		fmt.Printf("%s%s%s%s  plat %s%-3s%s  %s→ %s  arr%s %s%s  %s\n",
			green, c.Departure.Time.Format("15:04"), reset, delayTag(c.Departure.Delay),
			bold, c.Departure.Platform, reset,
			dim, c.Arrival.Station, reset, c.Arrival.Time.Format("15:04"), delayTag(c.Arrival.Delay), note)
	}
}
