package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/adeynack/finances/model"
	"github.com/adeynack/finances/model/database"
	"github.com/adeynack/finances/model/query"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

var (
	filename            string
	bookOwnerEmail      string
	defaultCurrencyCode string
	autoDeleteBook      bool
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	readFlags()
	moneydanceExport, err := parseMoneydanceJson()
	if err != nil {
		log.Fatalf("error parsing the Moneydance JSON export: %v", err)
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("error opening database connection: %v", err)
	}

	book, err := updateAndSaveBook(db, moneydanceExport)
	if err != nil {
		log.Fatalf("error initializing book: %v", err)
	}

	importOrFail(db, moneydanceExport, book)
}

func readFlags() {
	flag.StringVar(&filename, "md-import-file", "./tmp/md.json", "path to the MoneyDance JSON file to import.")
	flag.StringVar(&bookOwnerEmail, "book-owner-email", "joe@example.com", "E-Mail of the user who will own the imported book.")
	flag.StringVar(&defaultCurrencyCode, "default-currency", "EUR", "Default currency of the imported book (ISO Code).")
	flag.BoolVar(&autoDeleteBook, "auto-delete-book", false, "Will delete any existing book with the same before re-importing it.")
	flag.Parse()
}

func updateAndSaveBook(db *gorm.DB, moneydanceExport *MoneydanceExport) (*model.Book, error) {
	rawBookDate, err := fromMdIntDate(moneydanceExport.Metadata.ExportDate)
	if err != nil {
		return nil, fmt.Errorf("updateAndSaveBook: %w", err)
	}
	bookDate := time.Time(rawBookDate)
	bookName := fmt.Sprintf("%s (%04d-%02d-%02d)", moneydanceExport.Metadata.FileName, bookDate.Year(), bookDate.Month(), bookDate.Day())

	q := query.Use(db)
	b, u := q.Book, q.User

	book, err := b.Where(b.Name.Eq(bookName)).FirstOrInit()
	if err != nil {
		return nil, fmt.Errorf("updateAndSaveBook: finding book candidate: %w", err)
	}
	if book.ID != "" {
		log.Printf("Book with name %q found\n", bookName)
		if autoDeleteBook {
			result, err := b.Delete(book)
			if err != nil || result.RowsAffected != 1 {
				return nil, fmt.Errorf("updateAndSaveBook: deleting book %q: %w", bookName, err)
			}
		} else {
			return nil, fmt.Errorf("A book with name %q already exists. Can only import from scratch (book must not already exist)", bookName)
		}
	}

	log.Printf("Create book with name %q owned by user with email %q\n", bookName, bookOwnerEmail)
	bookOwner, err := u.Where(u.Email.Eq(bookOwnerEmail)).First()
	if err != nil {
		return nil, fmt.Errorf("updateAndSaveBook: finding book owner with email %q: %w", bookOwnerEmail, err)
	}

	book.OwnerID = bookOwner.ID
	book.DefaultCurrencyIsoCode = model.ISOCurrency(defaultCurrencyCode)
	err = b.Save(book)
	if err != nil {
		return nil, fmt.Errorf("updateAndSaveBook: saving book: %w", err)
	}

	return book, nil
}

func importOrFail(db *gorm.DB, moneydanceExport *MoneydanceExport, book *model.Book) {
	query := query.Use(db)
	mdCurrenciesById := lo.KeyBy(moneydanceExport.MdItemsPerType["curr"], func(i MoneydanceItem) string { return i["id"] })
	registerIdByMdAcctid := make(map[string]string)
	registerIdByMdOldId := make(map[string]string)

	if err := newRegisterImport(query, moneydanceExport.MdItemsPerType, registerIdByMdAcctid, registerIdByMdOldId, book, mdCurrenciesById).ImportRegisters(); err != nil {
		log.Fatalf("error import registers: %v", err)
	}
}
