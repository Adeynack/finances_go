package main

import (
	"log"

	"github.com/adeynack/finances/model"
	"github.com/adeynack/finances/model/database"
	"gorm.io/gen"
)

func main() {
	db, err := database.Connect()
	if err != nil && !database.FatalLogIfTrappedMigrationError(err) {
		log.Fatalln(err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: "./model/query",
		Mode:    gen.WithoutContext,
	})
	g.UseDB(db)
	g.ApplyBasic(model.All()...)
	g.Execute()
}
