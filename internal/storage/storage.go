package storage

type Storage interface {
	Get(mType, mName string) (any, error)
	Update(mType, mName string, mValue any) (any, error)
	HealthCheck() error
	Terminate()
}
