package server

import (
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
)

// parseNatsURLs parses the nats urls
func parseNatsURLs(s []string) ([]*url.URL, error) {
	urls := []*url.URL{}

	for _, a := range s {
		if len(a) == 0 {
			continue
		}
		if !strings.HasPrefix(a, "nats-route://") {
			a = fmt.Sprintf("nats-route://%s", a)
		}
		u, err := url.Parse(a)
		if err != nil {
			return nil, err
		}

		iprecords, err := net.LookupIP(u.Hostname())
		if err == nil {
			for _, rec := range iprecords {
				fmt.Println("IPZ", rec.String())
			}
		} else {
			log.Error().Err(err).Msgf("failed to lookup IP")
		}

		urls = append(urls, u)
	}

	return urls, nil
}
