package storage

type Gauge float64
type Counter int64

type Storage interface {
	GetCounter(mName string) (Counter, error)
	GetGauge(mName string) (Gauge, error)
	UpdateCounter(mName string, mValue Counter) (Counter, error)
	UpdateGauge(mName string, mValue Gauge) (Gauge, error)
	HealthCheck() error
	Terminate()
}
