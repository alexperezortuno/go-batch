package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RecordsProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "batch_records_processed_total",
		Help: "Total records read from the file",
	})
	RecordsInserted = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "batch_records_inserted_total",
		Help: "Total records successfully inserted",
	})
	RecordsInvalid = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "batch_records_invalid_total",
		Help: "Total invalid records",
	})
)

func InitMetrics() {
	prometheus.MustRegister(RecordsProcessed, RecordsInserted, RecordsInvalid)
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":2112", nil)
		if err != nil {
			return
		}
	}()
}
