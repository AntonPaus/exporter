package storage

type Storage interface {
	GetCounter(mName string) (int64, error)
	GetGauge(mName string) (float64, error)
	UpdateCounter(mName string, mValue int64) (int64, error)
	UpdateGauge(mName string, mValue float64) (float64, error)
	HealthCheck() error
	Terminate()
}
