package stemcell

type Importer interface {
	ImportFromPath(string, Props) (Stemcell, error)
}

type Finder interface {
	Find(string) (Stemcell, bool, error)
}

type Stemcell interface {
	ID() string

	RootFSPath() string
	Stack() string

	Delete() error
}
