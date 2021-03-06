package graph

import "lyrical-app/database"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Db *database.DB
}

func NewResolver(db *database.DB) *Resolver {
	return &Resolver{
		Db: db,
	}
}
