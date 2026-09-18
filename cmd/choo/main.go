package main

import (
	"cmp"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/neocolebunders/choo-cli/internal/irail"
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

// dateParam turns "tomorrow" or "DD/MM" into iRail's ddmmyy format.
func dateParam(d string) (string, error) {
	switch d {
	case "":
		return "", nil
	case "tomorrow":
		return time.Now().AddDate(0, 0, 1).Format("020106"), nil
	}
	t, err := time.Parse("2/1", d)
	if err != nil {
		return "", fmt.Errorf("bad date %q (use DD/MM or 'tomorrow')", d)
	}
	return fmt.Sprintf("%02d%02d%s", t.Day(), int(t.Month()), time.Now().Format("06")), nil
}

// timeParam turns "14:30" into iRail's hhmm format.
func timeParam(t string) (string, error) {
	if t == "" {
		return "", nil
	}
	tm, err := time.Parse("15:04", t)
	if err != nil {
		return "", fmt.Errorf("bad time %q (use HH:MM)", t)
	}
	return tm.Format("1504"), nil
}

func main() {
	n := flag.Int("n", 3, "number of connections to show")
	t := flag.String("t", "", "departure time, e.g. 14:30")
	d := flag.String("d", "", "date: DD/MM or 'tomorrow'")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: choo [-n 5] [-t 14:30] [-d tomorrow] <from> <to>   e.g. choo nord diest")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}

	date, err := dateParam(*d)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	tm, err := timeParam(*t)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	conns, err := irail.Connections(station(flag.Arg(0)), station(flag.Arg(1)), date, tm)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, c := range conns[:min(*n, len(conns))] {
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
