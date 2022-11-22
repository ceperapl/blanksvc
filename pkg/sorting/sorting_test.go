package sorting

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSortRule_toSort(t *testing.T) {
	testTable := []struct {
		rule         *sortRule
		sortOrder    SortOrder
		expectedSort *Sort
	}{
		{
			rule:         nil,
			sortOrder:    AscendingOrder,
			expectedSort: nil,
		},
		{
			rule: &sortRule{
				fieldName: "test_field",
			},
			sortOrder:    DescendingOrder,
			expectedSort: &Sort{Name: "test_field", Order: DescendingOrder},
		},
	}

	for i, testCase := range testTable {
		t.Run(fmt.Sprintf("test case: %d", i), func(t *testing.T) {
			sort := testCase.rule.toSort(testCase.sortOrder)
			assert.Equal(t, testCase.expectedSort, sort)
		})
	}
}

func TestSortRules_getDefault(t *testing.T) {
	testTable := []struct {
		rules           sortRules
		expectedDefault *sortRule
	}{
		{
			rules:           nil,
			expectedDefault: nil,
		},
		{
			rules: sortRules{
				{fieldName: "test_field1", alias: "test_alias1", isDefault: false, defaultOrder: AscendingOrder},
				{fieldName: "test_field2", alias: "test_alias2", isDefault: true, defaultOrder: DescendingOrder},
			},
			expectedDefault: &sortRule{
				fieldName:    "test_field2",
				alias:        "test_alias2",
				isDefault:    true,
				defaultOrder: DescendingOrder,
			},
		},
	}

	for i, testCase := range testTable {
		t.Run(fmt.Sprintf("test case: %d", i), func(t *testing.T) {
			defaultSort := testCase.rules.getDefault()
			assert.Equal(t, testCase.expectedDefault, defaultSort)
		})
	}
}

func TestSortRules_getRule(t *testing.T) {
	testTable := []struct {
		rules         sortRules
		fieldName     string
		expectedRule  *sortRule
		expectedError error
	}{
		{
			rules:         nil,
			fieldName:     "test1",
			expectedRule:  nil,
			expectedError: ErrFieldNotFound,
		},
		{
			rules: sortRules{
				{fieldName: "test_field1", alias: "test_alias1", isDefault: false, defaultOrder: AscendingOrder},
				{fieldName: "test_field2", alias: "test_alias2", isDefault: true, defaultOrder: DescendingOrder},
			},
			fieldName: "test_alias1",
			expectedRule: &sortRule{
				fieldName:    "test_field1",
				alias:        "test_alias1",
				isDefault:    false,
				defaultOrder: AscendingOrder,
			},
			expectedError: nil,
		},
		{
			rules: sortRules{
				{fieldName: "test_field1", alias: "test_alias1", isDefault: false, defaultOrder: AscendingOrder},
				{fieldName: "test_field2", alias: "test_alias2", isDefault: true, defaultOrder: DescendingOrder},
			},
			fieldName: "test_field2",
			expectedRule: &sortRule{
				fieldName:    "test_field2",
				alias:        "test_alias2",
				isDefault:    true,
				defaultOrder: DescendingOrder,
			},
			expectedError: nil,
		},
		{
			rules: sortRules{
				{fieldName: "test_field1", alias: "test_alias1", isDefault: false, defaultOrder: AscendingOrder},
				{fieldName: "test_field2", alias: "test_alias2", isDefault: true, defaultOrder: DescendingOrder},
			},
			fieldName:     "test_field3",
			expectedRule:  nil,
			expectedError: ErrFieldNotFound,
		},
	}
	for i, testCase := range testTable {
		t.Run(fmt.Sprintf("test case: %d", i), func(t *testing.T) {
			rule, err := testCase.rules.getRule(testCase.fieldName)
			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, testCase.expectedError))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedRule, rule)
			}
		})
	}
}

func TestSort_ToSQL(t *testing.T) {
	testTable := []struct {
		sort            *Sort
		expectedSQLSort string
	}{
		{
			sort:            &Sort{Name: "name", Order: AscendingOrder},
			expectedSQLSort: "ORDER BY name asc",
		},
		{
			sort:            &Sort{Name: "created_at", Order: DescendingOrder},
			expectedSQLSort: "ORDER BY created_at desc",
		},
		{
			sort:            nil,
			expectedSQLSort: "",
		},
	}

	for i, testCase := range testTable {
		t.Run(fmt.Sprintf("test case: %d", i), func(t *testing.T) {
			sqlSort := testCase.sort.ToSQL()
			assert.Equal(t, testCase.expectedSQLSort, sqlSort)
		})
	}
}

type Test struct {
	ID        string
	Name      string    `sorting:"name:name_alias, default:asc"`
	Age       int       `sorting:"age:age1"`
	CreatedAt time.Time `sorting:"created_at"`
}

type TestInvalidTag struct {
	ID          string
	Name        string    `sorting:"name:name_alias, default:asc"`
	Age         int       `sorting:"age:age1"`
	FailedField string    `sorting:"invalid tag format"`
	CreatedAt   time.Time `sorting:"created_at"`
}

type TestFieldNamesNotUnique struct {
	ID          string
	Name        string    `sorting:"name:name_alias, default:asc"`
	Age         int       `sorting:"name:age1"`
	FailedField string    `sorting:"invalid tag format"`
	CreatedAt   time.Time `sorting:"created_at"`
}

type TestAliasesNotUnique struct {
	ID          string
	Name        string    `sorting:"name:alias, default:asc"`
	Age         int       `sorting:"age:alias"`
	FailedField string    `sorting:"invalid tag format"`
	CreatedAt   time.Time `sorting:"created_at"`
}

type TestManyDefaults struct {
	ID        string
	Name      string    `sorting:"name:name_alias, default:asc"`
	Age       int       `sorting:"age:age1"`
	CreatedAt time.Time `sorting:"created_at,default"`
}

type TestDefaultSortOrder struct {
	ID        string
	Name      string    `sorting:"name:name_alias,default"`
	Age       int       `sorting:"age:age1"`
	CreatedAt time.Time `sorting:"created_at"`
}

func TestNew(t *testing.T) {
	testTable := []struct {
		name          string
		sortString    string
		structure     interface{}
		expectedSort  *Sort
		expectedError error
	}{
		{
			name:          "field name is used",
			sortString:    "name:desc",
			structure:     Test{},
			expectedSort:  &Sort{Name: "name", Order: DescendingOrder},
			expectedError: nil,
		},
		{
			name:          "alias is used",
			sortString:    "age1:desc",
			structure:     Test{},
			expectedSort:  &Sort{Name: "age", Order: DescendingOrder},
			expectedError: nil,
		},
		{
			name:          "default is used",
			sortString:    "",
			structure:     Test{},
			expectedSort:  &Sort{Name: "name", Order: AscendingOrder},
			expectedError: nil,
		},
		{
			name:          "error while getting rules from structure tags",
			sortString:    "created_at:asc",
			structure:     TestInvalidTag{},
			expectedSort:  nil,
			expectedError: ErrInvalidSortTagFormat,
		},
		{
			name:          "invalid sort string format",
			sortString:    "invalid_string",
			structure:     Test{},
			expectedSort:  nil,
			expectedError: ErrInvalidSortStringFormat,
		},
		{
			name:          "there is no rule for the field name",
			sortString:    "updated_at:asc",
			structure:     Test{},
			expectedSort:  nil,
			expectedError: ErrFieldNotFound,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			sort, err := New(testCase.sortString, testCase.structure)
			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, testCase.expectedError))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedSort, sort)
			}
		})
	}
}

func TestGetSortRules(t *testing.T) {
	testTable := []struct {
		name              string
		structure         interface{}
		expectedSortRules sortRules
		expectedError     error
	}{
		{
			name:      "successful case",
			structure: Test{},
			expectedSortRules: sortRules{
				{fieldName: "name", alias: "name_alias", isDefault: true, defaultOrder: "asc"},
				{fieldName: "age", alias: "age1", isDefault: false, defaultOrder: "asc"},
				{fieldName: "created_at", alias: "", isDefault: false, defaultOrder: "asc"},
			},
			expectedError: nil,
		},
		{
			name:              "error structure type",
			structure:         "not a structure",
			expectedSortRules: nil,
			expectedError:     ErrStructType,
		},
		{
			name:              "invalid tag format",
			structure:         TestInvalidTag{},
			expectedSortRules: nil,
			expectedError:     ErrInvalidSortTagFormat,
		},
		{
			name:              "field names are not unique",
			structure:         TestFieldNamesNotUnique{},
			expectedSortRules: nil,
			expectedError:     ErrFieldNamesNotUnique,
		},
		{
			name:              "aliases are not unique",
			structure:         TestAliasesNotUnique{},
			expectedSortRules: nil,
			expectedError:     ErrAliasesNotUnique,
		},
		{
			name:              "defaults cannot be more than one",
			structure:         TestManyDefaults{},
			expectedSortRules: nil,
			expectedError:     ErrDefaultsMoreThanOne,
		},
		{
			name:      "default sort order",
			structure: TestDefaultSortOrder{},
			expectedSortRules: sortRules{
				{fieldName: "name", alias: "name_alias", isDefault: true, defaultOrder: "asc"},
				{fieldName: "age", alias: "age1", isDefault: false, defaultOrder: "asc"},
				{fieldName: "created_at", alias: "", isDefault: false, defaultOrder: "asc"},
			},
			expectedError: nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			rules, err := getSortRules(testCase.structure)
			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, testCase.expectedError))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedSortRules, rules)
			}
		})
	}
}
