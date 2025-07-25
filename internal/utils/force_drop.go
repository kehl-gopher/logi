package utils

import (
	"context"
	"log"

	"github.com/uptrace/bun"
)

// drop tables incase of dirty database during migrations
func Droptables(db *bun.DB) {
	var tables []string

	err := db.NewSelect().
		Table("tablename").
		TableExpr("pg_tables").
		Where("schemaname = 'public'").
		Scan(context.Background(), &tables)

	if err != nil {
		log.Fatalf("failed to fetch tables: %v", err)
		return
	}

	for _, table := range tables {
		_, err := db.NewRaw("DROP TABLE IF EXISTS ? CASCADE", bun.Ident(table)).Exec(context.Background())
		if err != nil {
			log.Printf("failed to drop table %s: %v", table, err)
		} else {
			log.Printf("dropped table: %s", table)
		}
	}
}
