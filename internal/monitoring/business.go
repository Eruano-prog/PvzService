package monitoring

import (
	"AvitoPvz/internal/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type PrometheusBusinessMetrics struct {
	pvzCreated        prometheus.Counter
	receptionsCreated prometheus.Counter
	productsAdded     prometheus.Counter
}

func NewPrometheusMetrics() service.BusinessMetrics {
	return &PrometheusBusinessMetrics{
		pvzCreated: promauto.NewCounter(prometheus.CounterOpts{
			Name: "pvz_created_total",
			Help: "Total number of created PVZ",
		}),
		receptionsCreated: promauto.NewCounter(prometheus.CounterOpts{
			Name: "receptions_created_total",
			Help: "Total number of created order receptions",
		}),
		productsAdded: promauto.NewCounter(prometheus.CounterOpts{
			Name: "products_added_total",
			Help: "Total number of added products",
		}),
	}
}

func (m *PrometheusBusinessMetrics) IncPVZCreated() {
	m.pvzCreated.Inc()
}

func (m *PrometheusBusinessMetrics) IncReceptionCreated() {
	m.receptionsCreated.Inc()
}

func (m *PrometheusBusinessMetrics) IncProductAdded() {
	m.productsAdded.Inc()
}
