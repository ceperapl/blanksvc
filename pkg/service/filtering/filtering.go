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
		return "", fmt.Errorf("invalid condition type: %s", s)
	}
}

func NewFilter(filterString string, structure interface{}) (*Filter, error) {
	if filterString == "" {
		return nil, nil
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

func getCondition(condString string, structure interface{}) (*Condition, error) {
	conditionRegex := `^([a-zA-Z0-9]+)(>=|<=|>|<|==|!=|:|~)(.+)$`
	rgx := regexp.MustCompile(conditionRegex)
	matches := rgx.FindStringSubmatch(condString)
	if len(matches) == 0 {
		return nil, fmt.Errorf("wrong format of condition: %s", condString)
	}

	fieldName := matches[1]
	condSign := matches[2]
	condValueStr := matches[3]
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
	case "off":
		return nil, nil
	case "string":
		switch condType {
		case Like:
			condValue = condValueStr
		case In:
			strSlice := strings.Split(condValueStr, ",")
			condValue = strSlice
		case Equal, NotEqual:
			condValue = condValueStr
		default:
			return nil, fmt.Errorf("invalid condition for string type: %v", condSign)
		}
	case "*string":
		switch condType {
		case Like:
			condValue = condValueStr
		case In:
			strSlice := strings.Split(condValueStr, ",")
			condValue = strSlice
		case Equal, NotEqual:
			if condValueStr == "NULL" || condValueStr == "null" {
				condValue = nil
			} else {
				condValue = condValueStr
			}
		default:
			return nil, fmt.Errorf("invalid condition for string type: %v", condSign)
		}
	case "int":
	case "bool":
	case "float64":
	case "time.Time":
		switch condType {
		case Equal, GreaterThan, LessThan:
			const dateLayout = "2006-01-02"
			condValue, err = time.Parse(dateLayout, condValueStr)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("invalid condition '%v' for string type", condSign)
		}
	default:
		return nil, fmt.Errorf("the type of field '%s' is not supported", fieldType)
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
		return "", errors.New("error structure type")
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
	return "", fmt.Errorf("structure doesn't have the field '%s'", fieldName)
}

func replaceAbbr(abbr string) string {
	switch abbr {
	case "id":
		return "ID"
	}
	return abbr
}
