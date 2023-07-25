package processor

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type PrometheusMetrics struct {
	Registry              *prometheus.Registry
	PacketObservedCounter *prometheus.CounterVec
	PacketRelayedCounter  *prometheus.CounterVec
	LatestHeightGauge     *prometheus.GaugeVec
	WalletBalance         *prometheus.GaugeVec
	FeesSpent             *prometheus.GaugeVec
	BlockQueryFailure     *prometheus.CounterVec
	ClientExpiration      *prometheus.GaugeVec
}

func (m *PrometheusMetrics) AddPacketsObserved(path, chain, channel, port, eventType string, count int) {
	m.PacketObservedCounter.WithLabelValues(path, chain, channel, port, eventType).Add(float64(count))
}

func (m *PrometheusMetrics) IncPacketsRelayed(path, chain, channel, port, eventType string) {
	m.PacketRelayedCounter.WithLabelValues(path, chain, channel, port, eventType).Inc()
}

func (m *PrometheusMetrics) SetLatestHeight(chain string, height int64) {
	m.LatestHeightGauge.WithLabelValues(chain).Set(float64(height))
}

func (m *PrometheusMetrics) SetWalletBalance(chain, gasPrice, key, address, denom string, balance float64) {
	m.WalletBalance.WithLabelValues(chain, gasPrice, key, address, denom).Set(balance)
}

func (m *PrometheusMetrics) SetFeesSpent(chain, gasPrice, key, address, denom string, amount float64) {
	m.FeesSpent.WithLabelValues(chain, gasPrice, key, address, denom).Set(amount)
}

func (m *PrometheusMetrics) SetClientExpiration(pathName, chain, clientID, trustingPeriod string, timeToExpiration time.Duration) {
	m.ClientExpiration.WithLabelValues(pathName, chain, clientID, trustingPeriod).Set(timeToExpiration.Seconds())
}

func (m *PrometheusMetrics) IncBlockQueryFailure(chain, err string) {
	m.BlockQueryFailure.WithLabelValues(chain, err).Inc()
}

func NewPrometheusMetrics() *PrometheusMetrics {
	packetLabels := []string{"path", "chain", "channel", "port", "type"}
	heightLabels := []string{"chain"}
	blockQueryFailureLabels := []string{"chain", "type"}
	walletLabels := []string{"chain", "gas_price", "key", "address", "denom"}
	clientExpirationLables := []string{"path_name", "chain", "client_id", "trusting_period"}
	registry := prometheus.NewRegistry()
	registerer := promauto.With(registry)
	return &PrometheusMetrics{
		Registry: registry,
		PacketObservedCounter: registerer.NewCounterVec(prometheus.CounterOpts{
			Name: "cosmos_relayer_observed_packets",
			Help: "The total number of observed packets",
		}, packetLabels),
		PacketRelayedCounter: registerer.NewCounterVec(prometheus.CounterOpts{
			Name: "cosmos_relayer_relayed_packets",
			Help: "The total number of relayed packets",
		}, packetLabels),
		LatestHeightGauge: registerer.NewGaugeVec(prometheus.GaugeOpts{
			Name: "cosmos_relayer_chain_latest_height",
			Help: "The current height of the chain",
		}, heightLabels),
		WalletBalance: registerer.NewGaugeVec(prometheus.GaugeOpts{
			Name: "cosmos_relayer_wallet_balance",
			Help: "The current balance for the relayer's wallet",
		}, walletLabels),
		FeesSpent: registerer.NewGaugeVec(prometheus.GaugeOpts{
			Name: "cosmos_relayer_fees_spent",
			Help: "The amount of fees spent from the relayer's wallet",
		}, walletLabels),
		BlockQueryFailure: registerer.NewCounterVec(prometheus.CounterOpts{
			Name: "cosmos_relayer_block_query_errors_total",
			Help: "The total number of block query failures. The failures are separated into two catagories: 'RPC Client' and 'IBC Header'",
		}, blockQueryFailureLabels),
		ClientExpiration: registerer.NewGaugeVec(prometheus.GaugeOpts{
			Name: "cosmos_relayer_client_expiration_seconds",
			Help: "Seconds until the client expires",
		}, clientExpirationLables),
	}
}
