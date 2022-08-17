package service

type CommandService interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
	Count(value string) int
	Begin()
	Commit() error
	Rollback() error
}
