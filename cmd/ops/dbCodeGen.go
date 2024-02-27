package main

import (
	"log"

	"github.com/adeynack/finances/database"
	"github.com/adeynack/finances/model"
	"gorm.io/gen"
)

func dbCodeGen() {
	db, err := database.Connect()
	if err != nil && !database.FatalLogIfTrappedMigrationError(err) {
		log.Fatalln(err)
	}
	g := gen.NewGenerator(gen.Config{
		OutPath: "./model/query",
		Mode:    gen.WithoutContext, // | gen.WithQueryInterface | gen.WithDefaultQuery,
	})
	g.UseDB(db)
	g.ApplyBasic(model.Models()...)
	g.Execute()
}
