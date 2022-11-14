package filtering

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	ErrFilterStringIsEmpty    = errors.New("sort string is empty")
	ErrInvalidConditionType   = errors.New("invalid condition type")
	ErrInvalidConditionFormat = errors.New("invalid condition format")
	ErrInvalidCondition       = errors.New("invalid condition")
	ErrFieldTypeNotSupported  = errors.New("the type of the field is not supported")
	ErrStructType             = errors.New("error structure type")
	ErrFieldNotFound          = errors.New("structure doesn't have the field")
)

type Filter struct {
	Conditions []Condition
}

type ConditionType string

const (
	GreaterThan      ConditionType = ">"
	LessThan         ConditionType = "<"
	GreaterThanEqual ConditionType = ">="
	LessThanEqual    ConditionType = "<="
	NotEqual         ConditionType = "!="
	Equal            ConditionType = "=="
	In               ConditionType = ":"
	Like             ConditionType = "~"
)

type Condition struct {
	FieldName string
	Type      ConditionType
	Value     interface{}
}

func condTypeFromStr(s string) (ConditionType, error) {
	switch s {
	case ">":
		return GreaterThan, nil
	case "<":
		return LessThan, nil
	case ">=":
		return GreaterThanEqual, nil
	case "<=":
		return LessThanEqual, nil
	case "!=":
		return NotEqual, nil
	case "==":
		return Equal, nil
	case ":":
		return In, nil
	case "~":
		return Like, nil
	default:
		return "", fmt.Errorf("%w: %s", ErrInvalidConditionType, s)
	}
}

func NewFilter(filterString string, structure interface{}) (*Filter, error) {
	if filterString == "" {
		return nil, ErrFilterStringIsEmpty
	}
	var filter Filter
	conditions := strings.Split(filterString, ",,")
	for _, condString := range conditions {
		condition, err := getCondition(condString, structure)
		if err != nil {
			return nil, err
		}
		if condition != nil {
			filter.Conditions = append(filter.Conditions, *condition)
		}
	}
	return &filter, nil
}

// nolint: funlen
func getCondition(condString string, structure interface{}) (*Condition, error) {
	conditionRegex := `^([a-zA-Z0-9]+)(>=|<=|>|<|==|!=|:|~)(.+)$`
	rgx := regexp.MustCompile(conditionRegex)
	matches := rgx.FindStringSubmatch(condString)
	if len(matches) == 0 {
		return nil, fmt.Errorf("%w: '%s'", ErrInvalidConditionFormat, condString)
	}

	fieldName, condSign, condValueStr := matches[1], matches[2], matches[3]
	var condValue interface{}
	fieldType, err := getFieldType(fieldName, structure)
	if err != nil {
		return nil, err
	}
	condType, err := condTypeFromStr(condSign)
	if err != nil {
		return nil, err
	}
	switch fieldType {
	case "string", "*string":
		switch condType {
		case Like:
			condValue = condValueStr
		case In:
			strSlice := strings.Split(condValueStr, ",")
			condValue = strSlice
		case Equal, NotEqual:
			if fieldType == "*string" && condValueStr == "NULL" || condValueStr == "null" {
				condValue = nil
			} else {
				condValue = condValueStr
			}
		default:
			return nil, fmt.Errorf("%w: invalid sign '%s' for type '%s'", ErrInvalidCondition, condSign, fieldType)
		}
	case "int", "bool", "float64":
		switch condType {
		case Equal, NotEqual:
			condValue = condValueStr
		default:
			return nil, fmt.Errorf("%w: invalid sign '%s' for type '%s'", ErrInvalidCondition, condSign, fieldType)
		}
	case "time.Time":
		switch condType {
		case Equal, GreaterThan, LessThan:
			const dateLayout = "2006-01-02"
			condValue, err = time.Parse(dateLayout, condValueStr)
			if err != nil {
				return nil, fmt.Errorf("parse time: %w", err)
			}
		default:
			return nil, fmt.Errorf("%w: invalid sign '%s' for type '%s'", ErrInvalidCondition, condSign, fieldType)
		}
	default:
		return nil, fmt.Errorf("%w: '%s'", ErrFieldTypeNotSupported, fieldType)
	}
	return &Condition{
		FieldName: fieldName,
		Type:      condType,
		Value:     condValue,
	}, nil
}

func getFieldType(fieldName string, structure interface{}) (string, error) {
	structValue := reflect.TypeOf(structure)
	if structValue.Kind() != reflect.Struct {
		return "", ErrStructType
	}
	caser := cases.Title(language.English)
	for i := 0; i < structValue.NumField(); i++ {
		if structValue.Field(i).Name == caser.String(replaceAbbr(fieldName)) {
			fieldType := structValue.Field(i).Tag.Get("type")
			if fieldType != "" {
				return fieldType, nil
			}
			return structValue.Field(i).Type.String(), nil
		}
	}
	return "", fmt.Errorf("%w: '%s'", ErrFieldNotFound, fieldName)
}

func replaceAbbr(abbr string) string {
	// nolint: gocritic
	switch abbr {
	case "id":
		return "ID"
	}
	return abbr
}
