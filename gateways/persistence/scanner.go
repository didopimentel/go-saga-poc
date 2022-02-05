package persistence

type scanner interface {
	Scan(dest ...interface{}) error
}
