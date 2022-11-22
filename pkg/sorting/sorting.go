package sorting

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
)

var (
	ErrInvalidSortStringFormat = errors.New("sort string format is invalid")
	ErrInvalidSortTagFormat    = errors.New("sort rule in tag is invalid")
	ErrInvalidSortOrder        = errors.New("sort order is invalid")
	ErrStructType              = errors.New("error structure type")
	ErrFieldNotFound           = errors.New("structure doesn't have the field")
	ErrDefaultsMoreThanOne     = errors.New("defaults cannot be more than one")
	ErrFieldNamesNotUnique     = errors.New("field names are not unique")
	ErrAliasesNotUnique        = errors.New("aliases are not unique")
)

type SortOrder string

const (
	sortingTagName = "sorting"

	AscendingOrder  SortOrder = "asc"
	DescendingOrder SortOrder = "desc"
	DefaultOrder              = AscendingOrder

	sortTagRegex    = `^(\w+)(?::(\w+))?(?:,\s?(default)(?::(asc|desc))?)?$`
	sortStringRegex = `^(\w+):(asc|desc)$`
)

type sortRule struct {
	fieldName    string
	alias        string
	isDefault    bool
	defaultOrder SortOrder
}

func (s *sortRule) toSort(order SortOrder) *Sort {
	if s == nil {
		return nil
	}
	sort := Sort{
		Name:  s.fieldName,
		Order: order,
	}
	return &sort
}

type sortRules []sortRule

// getDefault returns the first default found
// the check for at most one default must be done earlier
func (s sortRules) getDefault() *sortRule {
	for _, rule := range s {
		if rule.isDefault {
			return &rule
		}
	}
	return nil
}

func (s sortRules) getRule(name string) (*sortRule, error) {
	// search the name in aliasses
	for _, rule := range s {
		if rule.alias == name {
			return &rule, nil
		}
	}
	// search the name in field names
	for _, rule := range s {
		if rule.fieldName == name {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("%w: %s", ErrFieldNotFound, name)
}

type Sort struct {
	Name  string
	Order SortOrder
}

func (s *Sort) ToSQL() string {
	if s == nil {
		return ""
	}
	return fmt.Sprintf("ORDER BY %s %s", s.Name, s.Order)
}

// New creates Sort from sortString and structure
// sortString is received from a user
// sortString must consist of the name of the structure field and the order separated by colon
// for example, name:asc
// the order can take the values 'asc' and 'desc'
// If the sortString is empty, the default sort will be set, or nothing if no default is specified
//
// Sorting can be configured via structure tags
// for example,
//              `sorting:"name:alias, default:asc"`
//              `sorting:"name:alias,default:desc"`
//              `sorting:"name:alias,default"`
//              `sorting:"name:alias"`
//              `sorting:"name, default:desc"`
//              `sorting:"name,default"`
// The first value is a database field name and its alias.
// Alias is optional and can be used to protect the database.
//
// The second value is optional. It defines the sorting by default and order.
// If Order is not specified, then it is equal to 'asc'.
//
// If there is no sorting tag in the passed structure, then this field is excluded from sorting.
// If there are no sorting tags in the passed structure,
// then sorting cannot be specified and will rely on the database.
func New(sortString string, structure interface{}) (*Sort, error) {
	rules, err := getSortRules(structure)
	if err != nil {
		return nil, fmt.Errorf("get sort rules from structure tags, %w", err)
	}
	// if sort string is empty return default Sort or nothing if no default is specified
	if sortString == "" {
		defaultRule := rules.getDefault()
		return defaultRule.toSort(defaultRule.defaultOrder), err
	}

	// parse sort string received from user
	rgx := regexp.MustCompile(sortStringRegex)
	matches := rgx.FindStringSubmatch(sortString)
	if len(matches) == 0 {
		return nil, fmt.Errorf("%w: '%s'", ErrInvalidSortStringFormat, sortString)
	}
	fieldName, order := matches[1], SortOrder(matches[2])
	rule, err := rules.getRule(fieldName)
	if err != nil {
		return nil, fmt.Errorf("getting sorting rule by field name: %w", err)
	}

	return rule.toSort(order), nil
}

func getSortRules(structure interface{}) (sortRules, error) {
	// the passed argument must be a struct
	structValue := reflect.TypeOf(structure)
	if structValue.Kind() != reflect.Struct {
		return nil, ErrStructType
	}
	var rules sortRules
	rgx := regexp.MustCompile(sortTagRegex)
	fieldNames := make(map[string]struct{})
	aliasses := make(map[string]struct{})
	defaultCount := 0
	for i := 0; i < structValue.NumField(); i++ {
		// search for sort tag
		tagValue, ok := structValue.Field(i).Tag.Lookup(sortingTagName)
		if !ok {
			continue
		}
		// parse sorting tag
		matches := rgx.FindStringSubmatch(tagValue)
		if len(matches) == 0 {
			return nil, fmt.Errorf("%w: '%s'", ErrInvalidSortTagFormat, tagValue)
		}
		fieldName, alias, isDefault, defaultOrder := matches[1], matches[2], matches[3] == "default", SortOrder(matches[4])
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
		// only one field can be default
		if isDefault {
			if defaultCount != 1 {
				defaultCount++
			} else {
				return nil, ErrDefaultsMoreThanOne
			}
		}
		// default order can be empty string
		if defaultOrder == "" {
			defaultOrder = DefaultOrder
		}
		rule := sortRule{
			fieldName:    fieldName,
			alias:        alias,
			isDefault:    isDefault,
			defaultOrder: defaultOrder,
		}
		rules = append(rules, rule)
	}
	return rules, nil
}
