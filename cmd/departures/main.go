package main

import (
	"context"
	"fmt"
	"html"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/jo-m/gohafas/pkg/departures"
	"github.com/jo-m/gohafas/pkg/zvv"
)

type flags struct {
	TimeOut time.Duration `arg:"-t,--timeout" default:"5s" help:"API timeout" placeholder:"T"`
}

// Description implements arg.Described in order
// to show a help text.
func (f *flags) Description() string {
	return "List next departures for multiple public transport stops and directions.\n"
}

// compile time check
var _ arg.Described = (*flags)(nil)

var stops []departures.Stop = []departures.Stop{
	{
		Name: "Oerlikon",
		LineDirections: map[string][]string{
			"*": {},
		},
		TimeToStop: 8 * time.Minute,
	},
	{
		Name: "Zürich, Hagenholz",
		LineDirections: map[string][]string{
			"787": {"Zürich, Bahnhof Oerlikon"},
			"781": {"Zürich, Bahnhof Oerlikon"},
		},
		TimeToStop: 3 * time.Minute,
	},
}

func main() {
	f := flags{}
	arg.MustParse(&f)

	client := zvv.NewClient()

	ctx, cancel := context.WithTimeout(context.Background(), f.TimeOut)
	defer cancel()

	conns, err := departures.Compile(ctx, client, stops, 10)
	if err != nil {
		panic(err)
	}

	for _, conn := range conns {
		if conn.Cancelled {
			continue
		}
		fmt.Printf("%s %s | von: %s\n", conn.Product.Emoji(), conn.Product.String(), html.UnescapeString(conn.Locations[0].Location.Name))
		loc := conn.MainLocation
		fmt.Printf("   In %s minute(s) (%s)\n", loc.Countdown, loc.DateTime())
		if loc.Platform != "" {
			fmt.Printf("   Platform: %s\n", loc.Platform)
		}
		if rt := loc.RealTime; rt.HasRealTime && (rt.IsDelayed || rt.IsPlatformChanged) {
			fmt.Printf("   Real time: %s (%s+%s) platform %s\n", rt.Countdown, loc.Countdown, rt.Delay, rt.Platform)
		}
	}
}
