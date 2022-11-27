package repositories

type Contract interface {
	DoJob(eventString string) error
}
