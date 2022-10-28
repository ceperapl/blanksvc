package sorting

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type SortOrder string

const (
	AscendingOrder  SortOrder = "asc"
	DescendingOrder SortOrder = "desc"
)

type Sort struct {
	FieldName string
	Order     SortOrder
}

func NewSort(sortString string, structure interface{}) (*Sort, error) {
	var sort Sort
	if sortString == "" {
		return nil, nil
	}
	sortParts := strings.Split(sortString, ":")
	if len(sortParts) != 2 {
		return nil, fmt.Errorf("wrong format of sort string: %s", sortString)
	}
	fieldName := sortParts[0]
	order := sortParts[1]
	var err error
	if err := checkField(fieldName, structure); err != nil {
		return nil, err
	}
	sort.FieldName = fieldName
	if sort.Order, err = sortOrderFromString(order); err != nil {
		return nil, err
	}
	return &sort, nil
}

func sortOrderFromString(s string) (SortOrder, error) {
	switch s {
	case "asc":
		return AscendingOrder, nil
	case "desc":
		return DescendingOrder, nil
	default:
		return "", fmt.Errorf("failed to get sort order by string: %s", s)
	}
}

func checkField(fieldName string, structure interface{}) error {
	structValue := reflect.TypeOf(structure)
	if structValue.Kind() != reflect.Struct {
		return errors.New("error structure type")
	}
	fieldExists := false
	caser := cases.Title(language.English)
	for i := 0; i < structValue.NumField(); i++ {
		if structValue.Field(i).Name == caser.String(replaceAbbr(fieldName)) {
			fieldExists = true
			break
		}
	}
	if !fieldExists {
		return fmt.Errorf("structure doesn't have the field '%s'", fieldName)
	}
	return nil
}

func replaceAbbr(abbr string) string {
	switch abbr {
	case "id":
		return "ID"
	}
	return abbr
}
