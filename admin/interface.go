package admin

type Videos interface { // Interfaces are named collection of method signature
	save() error
	list() ([]Video, error)
}
