package migrations

type Migration interface {
	Migrate() error
}
