package filtering

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	ErrFilterStringIsEmpty    = errors.New("sort string is empty")
	ErrInvalidConditionType   = errors.New("invalid condition type")
	ErrInvalidConditionFormat = errors.New("invalid condition format")
	ErrInvalidCondition       = errors.New("invalid condition")
	ErrFieldTypeNotSupported  = errors.New("the type of the field is not supported")
	ErrStructType             = errors.New("error structure type")
	ErrFieldNotFound          = errors.New("structure doesn't have the field")
	ErrInvalidFilterTagFormat = errors.New("filter rule in tag is invalid")
	ErrFieldNamesNotUnique    = errors.New("field names are not unique")
	ErrAliasesNotUnique       = errors.New("aliases are not unique")
	ErrInvalidInteger         = errors.New("invalid integer value")
	ErrInvalidFloat           = errors.New("invalid float value")
	ErrInvalidDate            = errors.New("invalid date value")
	ErrInvalidBool            = errors.New("invalid bool value")
)

const (
	GreaterThan      ConditionType = ">"
	LessThan         ConditionType = "<"
	GreaterThanEqual ConditionType = ">="
	LessThanEqual    ConditionType = "<="
	NotEqual         ConditionType = "!="
	Equal            ConditionType = "=="
	In               ConditionType = ":"
	Like             ConditionType = "~"

	nullValue = "null"

	filteringTagName = "filtering"
	filterTagRegex   = `^(\w+)(?::(\w+))?(?:,(date))?$`
	conditionRegex   = `^(\w+)(>=|<=|>|<|==|!=|:|~)(.+)$`
)

type ConditionType string

type filterRule struct {
	fieldName string
	fieldType string
	alias     string
}

type filterRules []filterRule

func (f filterRules) getRule(name string) (*filterRule, error) {
	// search the name in aliasses
	for _, rule := range f {
		if rule.alias == name {
			return &rule, nil
		}
	}
	// search the name in field names
	for _, rule := range f {
		if rule.fieldName == name {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("%w: %s", ErrFieldNotFound, name)
}

type Condition struct {
	FieldName string
	Type      ConditionType
	Value     interface{}
}

type Filter struct {
	Conditions []Condition
}

// nolint: funlen, gocognit
func (f *Filter) ToPostgresSQL(argNumberStartWith int) (string, []interface{}) {
	if f == nil {
		return "", nil
	}
	if len(f.Conditions) == 0 {
		return "", nil
	}

	sqlString := "WHERE"
	var values []interface{}
	argNumber := argNumberStartWith
	for i, condition := range f.Conditions {
		if i > 0 {
			sqlString = fmt.Sprintf("%s AND", sqlString)
		}

		switch condition.Type {
		case Like:
			sqlString = fmt.Sprintf("%s %s LIKE $%d", sqlString, condition.FieldName, argNumber)
			values = append(values, fmt.Sprintf("%%%s%%", condition.Value))
			argNumber++
		case In:
			var valuesStr string
			switch t := condition.Value.(type) {
			case []string:
				for i, value := range t {
					if i > 0 {
						valuesStr = fmt.Sprintf("%s, ", valuesStr)
					}
					valuesStr = fmt.Sprintf("%s$%d", valuesStr, argNumber)
					values = append(values, value)
					argNumber++
				}
			case []int64:
				for i, value := range t {
					if i > 0 {
						valuesStr = fmt.Sprintf("%s, ", valuesStr)
					}
					valuesStr = fmt.Sprintf("%s$%d", valuesStr, argNumber)
					values = append(values, value)
					argNumber++
				}
			}
			sqlString = fmt.Sprintf("%s %s IN(%s)", sqlString, condition.FieldName, valuesStr)
		case GreaterThan, GreaterThanEqual, LessThan, LessThanEqual:
			sqlString = fmt.Sprintf("%s %s %s $%d", sqlString, condition.FieldName, condition.Type, argNumber)
			values = append(values, condition.Value)
			argNumber++
		case Equal:
			if condition.Value == nil {
				sqlString = fmt.Sprintf("%s %s is null", sqlString, condition.FieldName)
			} else {
				sqlString = fmt.Sprintf("%s %s = $%d", sqlString, condition.FieldName, argNumber)
				values = append(values, condition.Value)
				argNumber++
			}
		case NotEqual:
			if condition.Value == nil {
				sqlString = fmt.Sprintf("%s %s is not null", sqlString, condition.FieldName)
			} else {
				sqlString = fmt.Sprintf("%s %s <> $%d", sqlString, condition.FieldName, argNumber)
				values = append(values, condition.Value)
				argNumber++
			}
		}
	}
	return sqlString, values
}

// New creates Filter from filterString and structure
// filterString is received from a user
// filterString can consist of filter conditions separated by double commas, like
// 		 		name~Ivan,,age<=50,,job:developer,tester,,createdAt>=2022-01-01
//              updatedAt==null
//              deadline:2022-01-02,2022-02-01,,completed!=NULL
// If the filterString is empty, nil Filter will be returned
//
// Filtering can be configured via structure tags
// for example,
//              `filtering:"name:alias"`
// The value is a database field name and its alias.
// Alias is optional and can be used to protect the database.
//
// If there is no filtering tag in the passed structure, then this field is excluded from filtering.
//
// Also filtering depends on the types of structure fields.

func New(filterString string, structure interface{}) (*Filter, error) {
	if filterString == "" {
		// nolint: nilnil
		return nil, nil
	}
	rules, err := getFilterRules(structure)
	if err != nil {
		return nil, fmt.Errorf("get filter rules from structure tags, %w", err)
	}
	var filter Filter
	conditionsStr := strings.Split(filterString, ",,")
	for _, conditionStr := range conditionsStr {
		condition, err := getCondition(conditionStr, rules)
		if err != nil {
			return nil, fmt.Errorf("get condition from string, %w", err)
		}
		if condition != nil {
			filter.Conditions = append(filter.Conditions, *condition)
		}
	}

	return &filter, nil
}

func getFilterRules(structure interface{}) (filterRules, error) {
	// the passed argument must be a struct
	structValue := reflect.TypeOf(structure)
	if structValue.Kind() != reflect.Struct {
		return nil, ErrStructType
	}
	var rules filterRules
	rgx := regexp.MustCompile(filterTagRegex)
	fieldNames := make(map[string]struct{})
	aliasses := make(map[string]struct{})
	for i := 0; i < structValue.NumField(); i++ {
		// search for filter tag
		tagValue, ok := structValue.Field(i).Tag.Lookup(filteringTagName)
		if !ok {
			continue
		}
		// parse filtering tag
		matches := rgx.FindStringSubmatch(tagValue)
		if len(matches) == 0 {
			return nil, fmt.Errorf("%w: '%s'", ErrInvalidFilterTagFormat, tagValue)
		}
		fieldName, alias, fieldType := matches[1], matches[2], matches[3]
		// fieldname uniqueness check
		if _, ok := fieldNames[fieldName]; ok {
			return nil, fmt.Errorf("%w: '%s'", ErrFieldNamesNotUnique, fieldName)
		}
		fieldNames[fieldName] = struct{}{}
		// alias uniqueness check
		if _, ok := aliasses[alias]; alias != "" && ok {
			return nil, fmt.Errorf("%w: '%s'", ErrAliasesNotUnique, alias)
		}
		aliasses[alias] = struct{}{}
		// getting field type
		if fieldType == "" {
			fieldType = structValue.Field(i).Type.String()
		}

		rule := filterRule{
			fieldName: fieldName,
			alias:     alias,
			fieldType: fieldType,
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func getCondition(condString string, rules filterRules) (*Condition, error) {
	rgx := regexp.MustCompile(conditionRegex)
	matches := rgx.FindStringSubmatch(condString)
	if len(matches) == 0 {
		return nil, fmt.Errorf("%w: '%s'", ErrInvalidConditionFormat, condString)
	}
	fieldName, condType, condValueStr := matches[1], ConditionType(matches[2]), matches[3]
	rule, err := rules.getRule(fieldName)
	if err != nil {
		return nil, fmt.Errorf("getting filtering rule by field name: %w", err)
	}

	var condition *Condition

	switch rule.fieldType {
	case "string", "*string":
		condition, err = getStringCondition(*rule, condType, condValueStr)
		if err != nil {
			return nil, fmt.Errorf("getting string condition: %w, %s", err, rule.fieldType)
		}
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64",
		"*int", "*int8", "*int16", "*int32", "*int64", "*uint", "*uint8", "*uint16", "*uint32", "*uint64":
		condition, err = getIntCondition(*rule, condType, condValueStr)
		if err != nil {
			return nil, fmt.Errorf("getting int condition: %w, %s", err, rule.fieldType)
		}
	case "float32", "float64",
		"*float32", "*float64":
		condition, err = getFloatCondition(*rule, condType, condValueStr)
		if err != nil {
			return nil, fmt.Errorf("getting float condition: %w, %s", err, rule.fieldType)
		}
	case "bool", "*bool":
		condition, err = getBoolCondition(*rule, condType, condValueStr)
		if err != nil {
			return nil, fmt.Errorf("getting bool condition: %w, %s", err, rule.fieldType)
		}
	case "date":
		condition, err = getDateCondition(*rule, condType, condValueStr)
		if err != nil {
			return nil, fmt.Errorf("getting date condition: %w, %s", err, rule.fieldType)
		}
	default:
		return nil, fmt.Errorf("%w: '%s'", ErrFieldTypeNotSupported, rule.fieldType)
	}
	return condition, nil
}

func getStringCondition(rule filterRule, condType ConditionType, condValueStr string) (*Condition, error) {
	var condValue interface{}

	switch condType {
	case Like:
		condValue = condValueStr
	case In:
		strSlice := strings.Split(condValueStr, ",")
		condValue = strSlice
	case Equal, NotEqual:
		if strings.HasPrefix(rule.fieldType, "*") && strings.ToLower(condValueStr) == nullValue {
			condValue = nil
		} else {
			condValue = condValueStr
		}
	default:
		return nil, fmt.Errorf("%w: invalid sign '%s' for type '%s'", ErrInvalidCondition, condType, rule.fieldType)
	}

	return &Condition{
		FieldName: rule.fieldName,
		Type:      condType,
		Value:     condValue,
	}, nil
}

func getIntCondition(rule filterRule, condType ConditionType, condValueStr string) (*Condition, error) {
	var condValue interface{}
	var err error

	switch condType {
	case LessThan, GreaterThan, LessThanEqual, GreaterThanEqual:
		// nolint: gomnd
		intValue, parseErr := strconv.ParseInt(condValueStr, 10, 0)
		if parseErr != nil {
			return nil, fmt.Errorf("%w: %s", ErrInvalidInteger, parseErr.Error())
		}
		condValue = intValue
	case In:
		strSlice := strings.Split(condValueStr, ",")
		intSlice := make([]int64, len(strSlice))
		for i, str := range strSlice {
			// nolint: gomnd
			intSlice[i], err = strconv.ParseInt(str, 10, 0)
			if err != nil {
				return nil, fmt.Errorf("%w: %s", ErrInvalidInteger, err.Error())
			}
		}
		condValue = intSlice
	case Equal, NotEqual:
		if strings.HasPrefix(rule.fieldType, "*") && strings.ToLower(condValueStr) == nullValue {
			condValue = nil
		} else {
			// nolint: gomnd
			intValue, parseErr := strconv.ParseInt(condValueStr, 10, 0)
			if parseErr != nil {
				return nil, fmt.Errorf("%w: %s", ErrInvalidInteger, parseErr.Error())
			}
			condValue = intValue
		}
	default:
		return nil, fmt.Errorf("%w: invalid sign '%s' for type '%s'", ErrInvalidCondition, condType, rule.fieldType)
	}

	return &Condition{
		FieldName: rule.fieldName,
		Type:      condType,
		Value:     condValue,
	}, nil
}

func getFloatCondition(rule filterRule, condType ConditionType, condValueStr string) (*Condition, error) {
	var condValue interface{}

	switch condType {
	case LessThan, GreaterThan, LessThanEqual, GreaterThanEqual:
		// nolint: gomnd
		floatValue, parseErr := strconv.ParseFloat(condValueStr, 32)
		if parseErr != nil {
			return nil, fmt.Errorf("%w: %s", ErrInvalidFloat, parseErr.Error())
		}
		condValue = floatValue
	case Equal, NotEqual:
		if strings.HasPrefix(rule.fieldType, "*") && strings.ToLower(condValueStr) == nullValue {
			condValue = nil
		} else {
			return nil, fmt.Errorf("%w: invalid sign '%s' for type '%s'", ErrInvalidCondition, condType, rule.fieldType)
		}
	default:
		return nil, fmt.Errorf("%w: invalid sign '%s' for type '%s'", ErrInvalidCondition, condType, rule.fieldType)
	}

	return &Condition{
		FieldName: rule.fieldName,
		Type:      condType,
		Value:     condValue,
	}, nil
}

func getBoolCondition(rule filterRule, condType ConditionType, condValueStr string) (*Condition, error) {
	var condValue interface{}
	var err error

	switch condType {
	case Equal, NotEqual:
		if strings.HasPrefix(rule.fieldType, "*") && strings.ToLower(condValueStr) == nullValue {
			condValue = nil
		} else {
			condValue, err = strconv.ParseBool(condValueStr)
			if err != nil {
				return nil, fmt.Errorf("%w: %s", ErrInvalidBool, err.Error())
			}
		}
	default:
		return nil, fmt.Errorf("%w: invalid sign '%s' for type '%s'", ErrInvalidCondition, condType, rule.fieldType)
	}

	return &Condition{
		FieldName: rule.fieldName,
		Type:      condType,
		Value:     condValue,
	}, nil
}

func getDateCondition(rule filterRule, condType ConditionType, condValueStr string) (*Condition, error) {
	var condValue interface{}
	var err error

	switch condType {
	case Equal, GreaterThan, LessThan, GreaterThanEqual, LessThanEqual, NotEqual:
		const dateLayout = "2006-01-02"
		condValue, err = time.Parse(dateLayout, condValueStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrInvalidDate, err.Error())
		}
	default:
		return nil, fmt.Errorf("%w: invalid sign '%s' for type '%s'", ErrInvalidCondition, condType, rule.fieldType)
	}

	return &Condition{
		FieldName: rule.fieldName,
		Type:      condType,
		Value:     condValue,
	}, nil
}
