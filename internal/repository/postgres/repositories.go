package postgres

import "database/sql"

type Repositories struct {
	Test   *TestRepo
	DBTest *DBTestRepo
	Users  *UsersRepo
	Orders *OrdersRepo
	Docs   *DocumentsRepo
}

func NewRepositories(db *sql.DB) Repositories {
	return Repositories{
		Test:   &TestRepo{db: db},
		DBTest: &DBTestRepo{db: db},
		Users:  &UsersRepo{db: db},
		Orders: &OrdersRepo{db: db},
		Docs:   &DocumentsRepo{db: db},
	}
}
