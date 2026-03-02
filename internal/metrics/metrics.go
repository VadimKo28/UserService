package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Registry exposes only app-specific metrics (no Go/process defaults).
var Registry = prometheus.NewRegistry()

var (
	usersCreatedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "users_created_total",
			Help: "Total number of created users per day.",
		},
		[]string{"date"},
	)
	subscriptionsCreatedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "subscriptions_created_total",
			Help: "Total number of created subscriptions per day.",
		},
		[]string{"date"},
	)
)

func init() {
	Registry.MustRegister(usersCreatedTotal, subscriptionsCreatedTotal)
}

func IncUsersCreated(at time.Time) {
	usersCreatedTotal.WithLabelValues(formatDay(at)).Inc()
}

func IncSubscriptionsCreated(at time.Time) {
	subscriptionsCreatedTotal.WithLabelValues(formatDay(at)).Inc()
}

func formatDay(at time.Time) string {
	return at.UTC().Format("2006-01-02")
}
