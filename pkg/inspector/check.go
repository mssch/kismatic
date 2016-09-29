package inspector

// A Check implements a workflow that validates a condition. If an error
// occurs while running the check, it returns false and the error. If the
// check is able to successfully determine the condition, it returns true
// if satisfied, or false otherwise.
type Check interface {
	Check() (bool, error)
}

// A ClosableCheck implements a long-running check workflow that requires closing
type ClosableCheck interface {
	Check
	Close() error
}
