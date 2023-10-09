package migrations

type InitDatabase struct {
	ID       string
	Migrate  any
	Rollback any
}
