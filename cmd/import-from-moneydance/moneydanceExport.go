package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type MoneydanceItem map[string]string

type MoneydanceMetadata struct {
	FileName   string `json:"file_name"`
	ExportDate uint64 `json:"export_date"`
}

type MoneydanceExport struct {
	Metadata       MoneydanceMetadata
	MdItemsPerType map[string][]MoneydanceItem
}

func parseMoneydanceJson() (*MoneydanceExport, error) {
	log.Printf("Opening Moneydance JSON export at %q\n", filename)
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("parseMoneydanceJson: opening file: %w", err)
	}

	log.Println("Parsing Moneydance JSON export")
	dec := json.NewDecoder(file)
	t, err := dec.Token()
	if err != nil {
		return nil, fmt.Errorf("parseMoneydanceJson: reading the first token: %w", err)
	}
	delim, ok := t.(json.Delim)
	if !ok || delim.String() != "{" {
		return nil, fmt.Errorf("parseMoneydanceJson: expected object start, got %q", t)
	}

	moneydanceExport := &MoneydanceExport{
		MdItemsPerType: make(map[string][]MoneydanceItem),
	}
	return moneydanceExport, parseMdJsonRootObject(dec, moneydanceExport)
}

func parseMdJsonRootObject(dec *json.Decoder, moneydanceExport *MoneydanceExport) error {
	for dec.More() {
		t, err := dec.Token()
		if err != nil {
			return fmt.Errorf("parseMdJsonRootObject: reading token: %w", err)
		}
		stringToken, ok := t.(string)
		if !ok {
			return fmt.Errorf("parseMdJsonRootObject: expecting token %q to be string", t)
		}

		switch stringToken {
		case "metadata":
			err = parseMetadata(dec, moneydanceExport)
		case "all_items":
			err = parseAllItems(dec, moneydanceExport)
		default:
			log.Printf("Ignoring root attribute %q\n", stringToken)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func parseMetadata(dec *json.Decoder, moneydanceExport *MoneydanceExport) error {
	err := dec.Decode(&moneydanceExport.Metadata)
	if err != nil {
		return fmt.Errorf("parseMetadata: decoding metadata: %w", err)
	}
	return nil
}

func parseAllItems(dec *json.Decoder, moneydanceExport *MoneydanceExport) error {
	t, err := dec.Token()
	if err != nil {
		return fmt.Errorf("parseAllItems: reading token in \"all_items\": %w", err)
	}
	delim, ok := t.(json.Delim)
	if !ok || delim.String() != "[" {
		return fmt.Errorf("parseAllItems: expecting \"all_items\" to be an array")
	}

	for dec.More() {
		var item MoneydanceItem
		err := dec.Decode(&item)
		if err != nil {
			return fmt.Errorf("parseAllItems: reading item: %w", err)
		}
		err = parseItem(item, moneydanceExport)
		if err != nil {
			return fmt.Errorf("parseAllItems: parsing item: %w", err)
		}
	}

	return nil
}

func parseItem(item MoneydanceItem, moneydanceExport *MoneydanceExport) error {
	objType, ok := item["obj_type"]
	if !ok || objType == "" {
		return fmt.Errorf("parseItem: item does not have \"obj_type\" attribute: %v", item)
	}

	items, ok := moneydanceExport.MdItemsPerType[objType]
	if !ok {
		items = make([]MoneydanceItem, 0)
	}
	moneydanceExport.MdItemsPerType[objType] = append(items, item)

	return nil
}
