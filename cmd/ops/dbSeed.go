package main

import (
	"fmt"
	"log"
	"os"

	"github.com/adeynack/finances/app/utils"
	"github.com/adeynack/finances/model"
	"github.com/adeynack/finances/model/database"
	"github.com/adeynack/finances/model/query"
	"gorm.io/gorm"
)

func dbSeed() {
	db, err := database.Connect()
	if err != nil && !database.FatalLogIfTrappedMigrationError(err) {
		log.Fatalln(err)
	}

	serverSalt := os.Getenv("SERVER_SALT")
	dbTruncateAll(db)
	users := dbSeedCreateUsers(db, serverSalt)
	dbSeedCreateBooks(db, users)
}

func dbTruncateAll(db *gorm.DB) {
	tables := []string{"books", "users"}
	for _, table := range tables {
		err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %q CASCADE", table)).Error
		if err != nil {
			panic(fmt.Errorf("error truncating table %q: %v", table, err))
		}
	}
}

func dbSeedCreateUsers(db *gorm.DB, salt string) map[string]*model.User {
	users := make(map[string]*model.User)
	user := query.Use(db).User

	add := func(key string, u *model.User) {
		u.SetPassword(u.EncryptedPassword, salt)
		users[key] = u
	}

	add("joe", &model.User{
		Email:             "joe@example.com",
		DisplayName:       "Joe",
		EncryptedPassword: "joe",
	})
	add("mary", &model.User{
		Email:             "mary@example.com",
		DisplayName:       "Proud Mary",
		EncryptedPassword: "mary",
	})
	add("vlad", &model.User{
		Email:             "vlad@example.com",
		DisplayName:       "Vlad the Impaler",
		EncryptedPassword: "vlad",
	})

	err := user.Create(utils.MapGetValues(users)...)
	if err != nil {
		panic(fmt.Errorf("error seeding users: %v", err))
	}

	return users
}

func dbSeedCreateBooks(db *gorm.DB, users map[string]*model.User) map[string]*model.Book {
	books := make(map[string]*model.Book)
	book := query.Use(db).Book

	books["joe"] = &model.Book{
		Name:                   "Joe's Book",
		Owner:                  users["joe"],
		DefaultCurrencyIsoCode: "EUR",
	}
	books["foo"] = &model.Book{
		Name:                   "Foo Inc.",
		Owner:                  users["joe"],
		DefaultCurrencyIsoCode: "USD",
	}
	books["mary"] = &model.Book{
		Name:                   "Book of Proud Mary",
		Owner:                  users["mary"],
		DefaultCurrencyIsoCode: "EUR",
	}
	books["vlad"] = &model.Book{
		Name:                   "The Financial Book of Vlad the Impaler",
		Owner:                  users["vlad"],
		DefaultCurrencyIsoCode: "RON",
	}

	err := book.Create(utils.MapGetValues(books)...)
	if err != nil {
		panic(fmt.Errorf("error seeding books: %v", err))
	}
	return books
}
