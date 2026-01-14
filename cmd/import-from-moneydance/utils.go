package main

import (
	"fmt"
	"math/bits"
	"strconv"
	"time"

	"github.com/adeynack/finances/model"
	"gorm.io/datatypes"
)

const (
	mdExternalSystem = "Moneydance"
)

func expiryDate(mdAccount MoneydanceItem) (*datatypes.Date, error) {
	var year, month int64
	var err error

	if rawYear, ok := mdAccount["exp_year"]; ok {
		if year, err = strconv.ParseInt(rawYear, 10, 32); err != nil {
			return nil, fmt.Errorf("error parsing year %q: %w", rawYear, err)
		}
	} else {
		return nil, nil
	}

	if rawMonth, ok := mdAccount["exp_month"]; ok {
		if month, err = strconv.ParseInt(rawMonth, 10, 32); err != nil {
			return nil, fmt.Errorf("error parsing month %q: %w", rawMonth, err)
		}
	} else {
		month = 1
	}

	t := datatypes.Date(time.Date(int(year), time.Month(month), 0, 0, 0, 0, 0, time.Local))
	return &t, nil
}

func fromMdIntDateStr(source string, defaultValue ...datatypes.Date) (datatypes.Date, error) {
	if source == "" && len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	intSource, err := strconv.ParseUint(source, 10, 64)
	if err != nil {
		return datatypes.Date{}, fmt.Errorf("parsing Moneydance integer-date from %q: %w", source, err)
	}

	return fromMdIntDate(intSource)
}

func fromMdIntDate(source uint64) (datatypes.Date, error) {
	source, day := bits.Div64(0, source, 100)
	year, month := bits.Div64(0, source, 100)

	t := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.Local)
	return datatypes.Date(t), nil
}

func fromMdUnixDate(source string, defaultValue ...model.ISODateTime) (model.ISODateTime, error) {
	if source == "" && len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	var result model.ISODateTime
	sourceInt, err := strconv.ParseInt(source, 10, 64)
	if err != nil {
		return result, fmt.Errorf("fromMdUnixDate: parsing date %q: %w", source, err)
	}

	result = model.ISODateTime(time.UnixMilli(sourceInt))
	return result, nil
}
