package main

import (
	"fmt"
	"log"
	"os"

	"github.com/adeynack/finances/app/appenv"
	"github.com/adeynack/finances/app/utils"
	"github.com/adeynack/finances/model"
	"github.com/adeynack/finances/model/database"
	"github.com/adeynack/finances/model/query"
	"gorm.io/gorm"
)

func main() {
	db, err := database.Connect()
	if err != nil && !database.FatalLogIfTrappedMigrationError(err) {
		log.Fatal(err)
	}

	serverSecret := os.Getenv(appenv.EnvServerSecret)
	dbTruncateAll(db)
	if err != nil {
		log.Fatal(err)
	}
	users := dbSeedCreateUsers(db, serverSecret)
	if err != nil {
		log.Fatal(err)
	}
	_ = dbSeedCreateBooks(db, users)
	if err != nil {
		log.Fatal(err)
	}
}

func dbTruncateAll(db *gorm.DB) {
	tables, err := utils.GormTablesForModels(db, model.All()...)
	if err != nil {
		log.Fatal(err)
	}
	for _, table := range tables {
		err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %q CASCADE", table)).Error
		if err != nil {
			log.Fatalf("error truncating table %q: %v", table, err)
		}
	}
}

func dbSeedCreateUsers(db *gorm.DB, secret string) map[string]*model.User {
	u := query.Use(db).User
	userFixtures := map[string]*model.User{
		"joe": {
			Email:             "joe@example.com",
			DisplayName:       "Joe",
			EncryptedPassword: "joe",
		},
		"mary": {
			Email:             "mary@example.com",
			DisplayName:       "Proud Mary",
			EncryptedPassword: "mary",
		},
		"vlad": {
			Email:             "vlad@example.com",
			DisplayName:       "Vlad the Impaler",
			EncryptedPassword: "vlad",
		},
	}
	for fixtureName, user := range userFixtures {
		user.SetPassword(user.EncryptedPassword, secret)
		err := u.Create(user)
		if err != nil {
			log.Fatalf("error seeding user %q: %v", fixtureName, err)
		}
	}
	return userFixtures
}

func dbSeedCreateBooks(db *gorm.DB, users map[string]*model.User) map[string]*model.Book {
	b := query.Use(db).Book
	bookFixtures := map[string]*model.Book{
		"joe": {
			Name:                   "Joe's Book",
			Owner:                  users["joe"],
			DefaultCurrencyIsoCode: "EUR",
		},
		"foo": {
			Name:                   "Foo Inc.",
			Owner:                  users["joe"],
			DefaultCurrencyIsoCode: "USD",
		},
		"mary": {
			Name:                   "Book of Proud Mary",
			Owner:                  users["mary"],
			DefaultCurrencyIsoCode: "EUR",
		},
		"vlad": {
			Name:                   "The Financial Book of Vlad the Impaler",
			Owner:                  users["vlad"],
			DefaultCurrencyIsoCode: "RON",
		},
	}
	for fixtureName, book := range bookFixtures {
		err := b.Create(book)
		if err != nil {
			log.Fatalf("error seeding book %q: %v", fixtureName, err)
		}
	}
	return bookFixtures
}
