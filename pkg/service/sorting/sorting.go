package sorting

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	ErrSortStringIsEmpty = errors.New("sort string is empty")
	ErrInvalidSortFormat = errors.New("sort format is invalid")
	ErrInvalidSortOrder  = errors.New("sort order is invalid")
	ErrStructType        = errors.New("error structure type")
	ErrFieldNotFound     = errors.New("structure doesn't have the field")
)

type SortOrder string

const (
	AscendingOrder  SortOrder = "asc"
	DescendingOrder SortOrder = "desc"

	sortPartsCount = 2
)

type Sort struct {
	FieldName string
	Order     SortOrder
}

func NewSort(sortString string, structure interface{}) (*Sort, error) {
	var sort Sort
	if sortString == "" {
		return nil, ErrSortStringIsEmpty
	}
	sortParts := strings.Split(sortString, ":")
	if len(sortParts) != sortPartsCount {
		return nil, fmt.Errorf("%w: %s", ErrInvalidSortFormat, sortString)
	}
	fieldName := sortParts[0]
	order := sortParts[1]
	var err error
	if err = checkField(fieldName, structure); err != nil {
		return nil, err
	}
	sort.FieldName = fieldName
	if sort.Order, err = sortOrderFromString(order); err != nil {
		return nil, fmt.Errorf("getting sort order from string: %w", err)
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
		return "", fmt.Errorf("%w: %s", ErrInvalidSortOrder, s)
	}
}

func checkField(fieldName string, structure interface{}) error {
	structValue := reflect.TypeOf(structure)
	if structValue.Kind() != reflect.Struct {
		return ErrStructType
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
		return fmt.Errorf("%w: '%s'", ErrFieldNotFound, fieldName)
	}
	return nil
}

func replaceAbbr(abbr string) string {
	// nolint: gocritic
	switch abbr {
	case "id":
		return "ID"
	}
	return abbr
}
