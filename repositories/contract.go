package repositories

// Contract interface is standardized for the repository
// If it wants to be a worker from events
type Contract interface {
	DoJob(eventString string) error
}
