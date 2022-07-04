package departures

/*
Package departures returns a list of next departures from a list of
stops and directions, plus walking distances.
Use it to get a list of next departures from a given starting point,
e.g. your office or flat.
*/

import (
	"context"
	"html"
	"sort"
	"time"

	"github.com/jo-m/gohafas/pkg/hafas"
)

type APIClient interface {
	StBoard(ctx context.Context, stopName, stopID string, ts time.Time, nResults int) ([]hafas.Conn, error)
	FindStop(ctx context.Context, query string) (*hafas.Stop, error)
}

// compile time check
var _ APIClient = (*hafas.Client)(nil)

// Stop is a public transport stop to list departures for.
type Stop struct {
	// stop name, e.g. "Zürich, Brunau/Mutschellenstr."
	Name string
	// map key: line name, e.g. "72" or "*" for all
	// values: list of line terminals, e.g.
	//   ["Zürich, Milchbuck", "Zürich, Albisriederplatz"]
	//   empty list means all directions
	LineDirections map[string][]string
	// time it takes you to walk to this stop from your starting point
	TimeToStop time.Duration
}

func connMatches(conn hafas.Conn, lineDirections map[string][]string) bool {
	direction := html.UnescapeString(conn.Product.Direction)

	// now, check if any line & direction matches this conn
	for line, directions := range lineDirections {
		// does the line match?
		if line != "*" && conn.Product.Line != line {
			continue
		}

		// check if a direction matches
		if len(directions) == 0 {
			return true
		}
		for _, d := range directions {
			if d == direction {
				return true
			}
		}
	}

	return false
}

func filterConns(conns []hafas.Conn, stop Stop) []hafas.Conn {
	ret := []hafas.Conn{}

	for _, conn := range conns {
		if connMatches(conn, stop.LineDirections) &&
			conn.MainLocation.BestCountdown() >= stop.TimeToStop {
			ret = append(ret, conn)
		}
	}

	return ret
}

// Compile compiles a list of departures for a given
// list of stops, lines, and directions.
func Compile(ctx context.Context, client APIClient, stops []Stop, nResults int) ([]hafas.Conn, error) {
	now := time.Now()

	allConns := []hafas.Conn{}
	for _, stopSearch := range stops {
		stop, err := client.FindStop(ctx, stopSearch.Name)
		if err != nil {
			return nil, err
		}

		conns, err := client.StBoard(ctx, stop.Name, stop.ID, now, nResults*3)
		if err != nil {
			return nil, err
		}

		allConns = append(allConns, filterConns(conns, stopSearch)...)
	}

	sort.Slice(allConns, func(i, j int) bool {
		return allConns[i].MainLocation.Countdown < allConns[j].MainLocation.Countdown
	})

	if len(allConns) < nResults {
		return allConns, nil
	}

	return allConns[:nResults], nil
}
