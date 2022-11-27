package main

import (
	"context"
	"fmt"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/jo-m/gohafas/pkg/zvv"
)

type flags struct {
	TimeOut    time.Duration `arg:"-t,--timeout" default:"5s" help:"API timeout" placeholder:"T"`
	MaxResults int           `arg:"-n,--max-results" default:"5" help:"max number of results" placeholder:"N"`
	Stop       string        `arg:"positional" default:"Zürich HB" help:"stop name, use cmd/findStops to find stops" placeholder:"STOP"`
}

// Description implements arg.Described in order
// to show a help text.
func (f *flags) Description() string {
	return "List next departures for a public transport stop.\n"
}

// compile time check
var _ arg.Described = (*flags)(nil)

func main() {
	f := flags{}
	arg.MustParse(&f)

	client := zvv.NewClient()

	ctx, cancel := context.WithTimeout(context.Background(), f.TimeOut)
	defer cancel()

	stop, err := client.FindStop(ctx, f.Stop)
	if err != nil {
		panic(err)
	}

	conns, err := client.StBoard(ctx, stop.Name, stop.ID, time.Now(), f.MaxResults)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Next departures for %s\n", stop.Name)
	for _, conn := range conns {
		fmt.Printf("%s %s\n", conn.Product.Emoji(), conn.Product.String())
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
