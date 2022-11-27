# Go API client for HAFAS

Currently, only stop search and stationboard endpoints are implemented.
The only API which has been tested is the [ZVV](https://www.zvv.ch/) one.

It has the base URL <https://zvv.hafas.cloud/bin> and a nice web
interface at <https://online.fahrplan.zvv.ch/bin/stboard.exe>.

## HAFAS

Read more about it here:

- <https://de.wikipedia.org/wiki/HAFAS>
- <https://administrator.de/knowledge/hafas-fahrplanauskunft-api-sammlung-177145.html>
- <https://gist.github.com/derhuerst/2b7ed83bfa5f115125a5>

## Usage and example binary

See the example CLI utility in `cmd/stationboard` for example usage.

```
$ go run ./cmd/stationboard/ --help
List next departures for a public transport stop.

Usage: stationboard [--timeout T] [--max-results N] [STOP]

Positional arguments:
  STOP                   stop name, use cmd/findStops to find stops

Options:
  --timeout T, -t T      API timeout [default: 5s]
  --max-results N, -n N
                         max number of results [default: 5]
  --help, -h             display this help and exit

$ go run ./cmd/stationboard/ 'Zürich HB'
Next departures for Zürich HB
🚊 Tram 13 to Zürich, Hardturm
   In 0 minute(s) (27.11.22 15:38)
🚆 InterCity 3 to Chur
   In 0 minute(s) (27.11.22 15:38)
   Platform: 8
🚆 RegioExpress  to Aarau
   In 0 minute(s) (27.11.22 15:38)
   Platform: 15
🚆 S-Bahn S4 to Langnau-Gattikon
   In 0 minute(s) (27.11.22 15:38)
   Platform: 21
🚊 Tram 4 to Zürich, Bahnhof Tiefenbrunnen
   In 0 minute(s) (27.11.22 15:38)
   Real time: 1 (0+1) platform
```
