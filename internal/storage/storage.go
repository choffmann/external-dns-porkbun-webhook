package storage

type PorkbunRepository interface {
}

type Repository struct {
	Porkbun PorkbunRepository
}
