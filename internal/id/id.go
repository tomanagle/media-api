package id

type ID interface {
	New() (string, error)
}
