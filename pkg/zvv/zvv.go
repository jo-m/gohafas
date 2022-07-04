package zvv

import "github.com/jo-m/gohafas/pkg/hafas"

const (
	BaseURL  = "https://zvv.hafas.cloud/bin"
	TimeZone = "Europe/Zurich"
)

func NewClient() *hafas.Client {
	ret, err := hafas.NewClient(BaseURL, TimeZone)
	if err != nil {
		panic(err)
	}
	return ret
}
