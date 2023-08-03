package prometheus

import (
	"context"

	"github.com/alexfalkowski/go-service/version"
	"github.com/go-redis/cache/v8"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
)

// Register for prometheus.
func Register(lc fx.Lifecycle, cache *cache.Cache, version version.Version) {
	collector := NewStatsCollector(cache, version)
	// TODO
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			prometheus.MustRegister(collector)

			return nil
		}
	})
}
