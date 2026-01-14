package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/adeynack/finances/app/utils"
	"github.com/adeynack/finances/model"
	"github.com/adeynack/finances/model/query"
	"github.com/samber/lo"
	"gorm.io/datatypes"
)

type registerImport struct {
	query                *query.Query
	mdItemsPerType       map[string][]MoneydanceItem
	registerIdByMdAcctid map[string]string
	registerIdByMdOldId  map[string]string
	book                 *model.Book
	mdCurrenciesById     map[string]MoneydanceItem

	mdAccountsByParentId map[string][]MoneydanceItem
}

func newRegisterImport(
	query *query.Query,
	mdItemsPerType map[string][]MoneydanceItem,
	registerIdByMdAcctid map[string]string,
	registerIdByMdOldId map[string]string,
	book *model.Book,
	mdCurrenciesById map[string]MoneydanceItem,
) *registerImport {
	return &registerImport{
		query:                query,
		mdItemsPerType:       mdItemsPerType,
		registerIdByMdAcctid: registerIdByMdAcctid,
		registerIdByMdOldId:  registerIdByMdOldId,
		book:                 book,
		mdCurrenciesById:     mdCurrenciesById,
	}
}

func (i *registerImport) ImportRegisters() error {
	log.Println("Importing registers (Moneydance accounts)")

	i.mdAccountsByParentId = lo.GroupBy(i.mdItemsPerType["acct"], func(i MoneydanceItem) string { return i["parentid"] })

	rootAccounts := i.mdAccountsByParentId[""]
	switch c := len(rootAccounts); {
	case c > 1:
		return fmt.Errorf("ImportRegisters: more than one root account found")
	case c == 0:
		return fmt.Errorf("ImportRegisters: no root account found")
	}
	rootAccount := rootAccounts[0]

	firstLevelAccountsPerType := lo.GroupBy(i.mdAccountsByParentId[rootAccount["id"]], func(a MoneydanceItem) string { return a["type"] })
	if err := i.importAccountBatchesByTypePriorities(firstLevelAccountsPerType, "i", "e"); err != nil {
		return fmt.Errorf("ImportRegisters: %w", err)
	}

	return nil
}

func (i *registerImport) importAccountBatchesByTypePriorities(firstLevelAccountsPerType map[string][]MoneydanceItem, types ...string) error {
	// Import the account by type priority.
	for _, typeToImport := range types {
		if accountsToImport, ok := firstLevelAccountsPerType[typeToImport]; ok {
			if err := i.importAccountBatch(accountsToImport); err != nil {
				return fmt.Errorf("importAccountBatchesByTypePriorities: accounts of type %q: %w", typeToImport, err)
			}
			delete(firstLevelAccountsPerType, typeToImport)
		}
	}

	// Import the rest
	for accType, accountsToImport := range firstLevelAccountsPerType {
		if err := i.importAccountBatch(accountsToImport); err != nil {
			return fmt.Errorf("importAccountBatchesByTypePriorities: accounts of type %q: %w", accType, err)
		}
	}

	return nil
}

func (i *registerImport) importAccountBatch(accountsToImport []MoneydanceItem) error {
	log.Printf("Importing Moneydance accounts of type %q\n", accountsToImport[0]["type"])
	for _, mdAccount := range accountsToImport {
		if err := i.importAccountRecursively(nil, mdAccount); err != nil {
			return fmt.Errorf("importAccountBatch: %w", err)
		}
	}
	return nil
}

func (i *registerImport) importAccountRecursively(parentRegister *model.Register, mdAccount MoneydanceItem) error {
	log.Printf("Importing Moneydance %q type account %q (ID %q)\n", mdAccount["type"], mdAccount["name"], mdAccount["id"])
	register, err := i.createRegister(mdAccount, parentRegister)
	if err != nil {
		return fmt.Errorf("importAccountRecursively: %w", err)
	}

	i.registerIdByMdAcctid[mdAccount["id"]] = register.ID
	i.registerIdByMdOldId[mdAccount["old_id"]] = register.ID

	if err = model.CreateImportOrigin(i.query.ImportOrigin.UnderlyingDB(), register, mdExternalSystem, mdAccount["id"]); err != nil {
		return fmt.Errorf("importAccountRecursively: %w", err)
	}

	if err = i.importChildAccounts(mdAccount, register); err != nil {
		return fmt.Errorf("importAccountRecursively: %w", err)

	}

	return nil
}

func (i *registerImport) extractRegisterType(mdAccountType string) (string, error) {
	switch mdAccountType {
	case "a":
		return "Asset", nil
	case "b":
		return "Bank", nil
	case "c":
		return "Card", nil
	case "e":
		return "Expense", nil
	case "i":
		return "Income", nil
	case "l":
		return "Liability", nil
	case "o":
		return "Loan", nil
	case "v":
		return "Investment", nil
	default:
		return "", fmt.Errorf("Unknown account type code %q", mdAccountType)
	}
}

func (i *registerImport) createRegister(mdAccount MoneydanceItem, parentRegister *model.Register) (*model.Register, error) {
	register := new(model.Register)

	if registerType, err := i.extractRegisterType(mdAccount["type"]); err == nil {
		register.Type = registerType
	} else {
		return nil, fmt.Errorf("createRegister: parsing \"type\": %w", err)
	}

	if dt, err := fromMdUnixDate(mdAccount["creation_date"], model.Now()); err == nil {
		register.CreatedAt = dt
	} else {
		return nil, fmt.Errorf("createRegister: parsing \"creation_date\" from MD account %q: %w", mdAccount["id"], err)
	}

	register.Name = mdAccount["name"]
	register.BookID = i.book.ID

	if parentRegister != nil {
		register.ParentID = parentRegister.ID
	}

	if date, err := fromMdIntDateStr(mdAccount["date_created"], datatypes.Date(i.book.CreatedAt)); err == nil {
		register.StartsAt = date
	} else {
		return nil, fmt.Errorf("createRegister: %w", err)
	}

	if currency, err := i.extractCurrencyISOCode(mdAccount); err == nil {
		register.CurrentcyIsoCode = currency
	} else {
		return nil, fmt.Errorf("createRegister: %w", err)
	}

	if initialBalance, err := strconv.Atoi(mdAccount["sbal"]); err == nil {
		register.InitialBalance = int64(initialBalance)
	} else {
		return nil, fmt.Errorf("createRegister: parsing initial balance from %q: %w", mdAccount["mbal"], err)
	}

	if isInactive, ok := mdAccount["is_inactive"]; !ok || isInactive != "y" {
		register.Active = true
	}

	if defaultCategoryId, err := i.extractDefaultCategoryId(mdAccount); err == nil {
		register.DefaultCategoryID = defaultCategoryId
	} else {
		return nil, fmt.Errorf("createRegister: getting default category: %w", err)
	}
	register.Notes = mdAccount["comment"]

	if err := i.assignAccountInformation(mdAccount, register); err != nil {
		return nil, fmt.Errorf("createRegister: assigning account information: %w", err)
	}

	if err := i.query.Register.Save(register); err != nil {
		return nil, fmt.Errorf("createRegister: saving register to database: %w", err)
	}

	return register, nil
}

func (i *registerImport) importChildAccounts(mdAccount MoneydanceItem, parentRegister *model.Register) error {
	accounts, ok := i.mdAccountsByParentId[mdAccount["id"]]
	if !ok {
		return nil
	}
	for _, childAccount := range accounts {
		if err := i.importAccountRecursively(parentRegister, childAccount); err != nil {
			return fmt.Errorf("importChildAccounts: importing MD account %q recursively: %w", childAccount["id"], err)
		}
	}
	return nil
}

func (i *registerImport) extractCurrencyISOCode(mdAccount MoneydanceItem) (model.ISOCurrency, error) {
	currId, ok := mdAccount["currid"]
	if !ok {
		return model.ISOCurrency(""), fmt.Errorf("Moneydance account did not have a \"currid\" attribute")
	}
	mdCurrency, ok := i.mdCurrenciesById[currId]
	if !ok {
		return model.ISOCurrency(""), fmt.Errorf("Moneydance currency for ID %q not found", currId)
	}
	result, ok := mdCurrency["currid"]
	if !ok {
		return model.ISOCurrency(""), fmt.Errorf("Moneydance currency with ID %q does not have attribute \"currid\"", currId)
	}
	return model.ISOCurrency(result), nil
}

func (i *registerImport) extractDefaultCategoryId(mdAccount MoneydanceItem) (string, error) {
	defaultCategory, ok := mdAccount["default_category"]
	if !ok || defaultCategory == "" {
		return "", nil
	}
	registerId, ok := i.registerIdByMdOldId[defaultCategory]
	if !ok {
		return "", fmt.Errorf("No register with old ID %q found", defaultCategory)
	}
	return registerId, nil
}

func (i *registerImport) assignAccountInformation(mdAccount MoneydanceItem, register *model.Register) error {
	switch register.Type {
	case model.RegisterTypeBank:
		return i.assignBankAccountInformation(mdAccount, register)
	case model.RegisterTypeCard:
		return i.assignCardAccountInformation(mdAccount, register)
	case model.RegisterTypeInvestment:
		return i.assignInvestmentAccountInformation(mdAccount, register)
	}
	return nil
}

func (i *registerImport) assignBankAccountInformation(mdAccount MoneydanceItem, register *model.Register) error {
	register.AccountNumber = utils.MapGetOr(mdAccount, "bank_account_number", "")

	return nil
}

func (i *registerImport) assignCardAccountInformation(mdAccount MoneydanceItem, register *model.Register) error {
	register.InstitutionName = utils.MapGetOr(mdAccount, "bank_name", "")

	if rawApr, ok := mdAccount["apr"]; ok {
		if apr, err := strconv.ParseFloat(rawApr, 32); err == nil {
			register.AnnualInterestRate = float32(apr)
		} else {
			return fmt.Errorf("unable to parse annual interest rate value %q: %w", rawApr, err)
		}
	}

	if rawCreditLimit, ok := mdAccount["credit_limit"]; ok {
		if creditLimit, err := strconv.ParseInt(rawCreditLimit, 10, 64); err == nil {
			register.CreditLimit = int64(creditLimit)
		} else {
			return fmt.Errorf("unable to parse credit limit value %q: %w", rawCreditLimit, err)
		}
	}

	register.CardNumber = utils.MapGetOr(mdAccount, "bank_account_number", "")

	if expiryDate, err := expiryDate(mdAccount); err == nil {
		register.ExpiresAt = expiryDate
	} else {
		return err
	}

	return nil
}

func (i *registerImport) assignInvestmentAccountInformation(mdAccount MoneydanceItem, register *model.Register) error {
	register.AccountNumber = utils.MapGetOr(mdAccount, "invst_account_number", "")
	return nil
}
