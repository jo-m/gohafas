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
	Search     string        `arg:"positional,required" help:"search text" placeholder:"SEARCH-TEXT"`
}

// Description implements arg.Described in order
// to show a help text.
func (f *flags) Description() string {
	return "Search for public transport stop names.\n"
}

// compile time check
var _ arg.Described = (*flags)(nil)

func main() {
	f := flags{}
	arg.MustParse(&f)

	client := zvv.NewClient()

	ctx, cancel := context.WithTimeout(context.Background(), f.TimeOut)
	defer cancel()

	stops, err := client.GetStops(ctx, f.Search, f.MaxResults)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Results for '%s'\n", f.Search)
	for _, stop := range stops {
		fmt.Printf("%s %s\n", stop.TypeStr, stop.Name)
	}
}
