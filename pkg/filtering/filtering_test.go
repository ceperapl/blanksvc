package filtering

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFilterRules_getRule(t *testing.T) {
	testTable := []struct {
		rules         filterRules
		fieldName     string
		expectedRule  *filterRule
		expectedError error
	}{
		{
			rules:         nil,
			fieldName:     "test1",
			expectedRule:  nil,
			expectedError: ErrFieldNotFound,
		},
		{
			rules: filterRules{
				{fieldName: "test_field1", alias: "test_alias1", fieldType: "string"},
				{fieldName: "test_field2", alias: "test_alias2", fieldType: "int"},
			},
			fieldName: "test_alias1",
			expectedRule: &filterRule{
				fieldName: "test_field1",
				alias:     "test_alias1",
				fieldType: "string",
			},
			expectedError: nil,
		},
		{
			rules: filterRules{
				{fieldName: "test_field1", alias: "test_alias1", fieldType: "string"},
				{fieldName: "test_field2", alias: "test_alias2", fieldType: "int"},
			},
			fieldName: "test_field2",
			expectedRule: &filterRule{
				fieldName: "test_field2",
				alias:     "test_alias2",
				fieldType: "int",
			},
			expectedError: nil,
		},
		{
			rules: filterRules{
				{fieldName: "test_field1", alias: "test_alias1", fieldType: "string"},
				{fieldName: "test_field2", alias: "test_alias2", fieldType: "int"},
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

func TestFilter_ToSQL(t *testing.T) {
	testTable := []struct {
		filter             *Filter
		argNumberStartWith int
		expectedSQL        string
		expectedValues     []interface{}
	}{
		{
			filter: &Filter{[]Condition{
				{FieldName: "name", Type: Equal, Value: "test"},
				{FieldName: "age", Type: In, Value: []int64{23, 24, 25}},
				{FieldName: "description", Type: Equal, Value: nil},
			}},
			argNumberStartWith: 1,
			expectedSQL:        "WHERE name = $1 AND age IN($2, $3, $4) AND description is null",
			expectedValues:     []interface{}{"test", int64(23), int64(24), int64(25)},
		},
		{
			filter: &Filter{[]Condition{
				{FieldName: "name", Type: Equal, Value: "test"},
				{FieldName: "age", Type: In, Value: []int64{23, 24, 25}},
			}},
			argNumberStartWith: 3,
			expectedSQL:        "WHERE name = $3 AND age IN($4, $5, $6)",
			expectedValues:     []interface{}{"test", int64(23), int64(24), int64(25)},
		},
		{
			filter: &Filter{[]Condition{
				{FieldName: "name", Type: NotEqual, Value: "test5"},
				{FieldName: "name", Type: Like, Value: "test"},
				{FieldName: "description", Type: NotEqual, Value: nil},
				{FieldName: "name", Type: In, Value: []string{"test6", "test7"}},
				{FieldName: "f", Type: LessThan, Value: float64(5.5)},
			}},
			argNumberStartWith: 2,
			expectedSQL:        "WHERE name <> $2 AND name LIKE $3 AND description is not null AND name IN($4, $5) AND f < $6",
			expectedValues:     []interface{}{"test5", "%test%", "test6", "test7", 5.5},
		},
		{
			filter:             nil,
			argNumberStartWith: 1,
			expectedSQL:        "",
			expectedValues:     nil,
		},
		{
			filter:             &Filter{[]Condition{}},
			argNumberStartWith: 1,
			expectedSQL:        "",
			expectedValues:     nil,
		},
	}
	for i, testCase := range testTable {
		t.Run(fmt.Sprintf("test case: %d", i), func(t *testing.T) {
			sql, values := testCase.filter.ToPostgresSQL(testCase.argNumberStartWith)
			assert.Equal(t, testCase.expectedSQL, sql)
			assert.Equal(t, testCase.expectedValues, values)
		})
	}
}

type Test struct {
	ID       string
	Name     string `filtering:"name:name_alias"`
	Age      int    `filtering:"age:age1"`
	Deadline string `filtering:"deadline,date"`
}

type TestInvalidTag struct {
	ID       string
	Name     string `filtering:"name:name_alias"`
	Age      int    `filtering:"invalid tag format"`
	Deadline string `filtering:"deadline,date"`
}

type TestFieldNamesNotUnique struct {
	ID       string
	Name     string `filtering:"name:name_alias"`
	Age      int    `filtering:"name:age1"`
	Deadline string `filtering:"deadline,date"`
}

type TestAliasesNotUnique struct {
	ID       string
	Name     string `filtering:"name:alias"`
	Age      int    `filtering:"age:alias"`
	Deadline string `filtering:"deadline,date"`
}

func TestNew(t *testing.T) {
	testTable := []struct {
		name           string
		filterString   string
		structure      interface{}
		expectedFilter *Filter
		expectedError  error
	}{
		{
			name:         "field names are used",
			filterString: "name~Ivan,,age==25",
			structure:    Test{},
			expectedFilter: &Filter{
				Conditions: []Condition{
					{FieldName: "name", Type: Like, Value: "Ivan"},
					{FieldName: "age", Type: Equal, Value: int64(25)},
				},
			},
			expectedError: nil,
		},
		{
			name:         "field names are used",
			filterString: "name_alias~Ivan,,age1==25",
			structure:    Test{},
			expectedFilter: &Filter{
				Conditions: []Condition{
					{FieldName: "name", Type: Like, Value: "Ivan"},
					{FieldName: "age", Type: Equal, Value: int64(25)},
				},
			},
			expectedError: nil,
		},
		{
			name:           "error while getting rules from structure tags",
			filterString:   "created_at:asc",
			structure:      TestInvalidTag{},
			expectedFilter: nil,
			expectedError:  ErrInvalidFilterTagFormat,
		},
		{
			name:           "empty filter string",
			filterString:   "",
			structure:      Test{},
			expectedFilter: nil,
			expectedError:  nil,
		},
		{
			name:           "error while getting condition from string",
			filterString:   "age~35",
			structure:      Test{},
			expectedFilter: nil,
			expectedError:  ErrInvalidCondition,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			filter, err := New(testCase.filterString, testCase.structure)
			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, testCase.expectedError))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedFilter, filter)
			}
		})
	}
}

func TestGetFilterRules(t *testing.T) {
	testTable := []struct {
		name               string
		structure          interface{}
		expectedFilterules filterRules
		expectedError      error
	}{
		{
			name:      "successful case",
			structure: Test{},
			expectedFilterules: filterRules{
				{fieldName: "name", fieldType: "string", alias: "name_alias"},
				{fieldName: "age", fieldType: "int", alias: "age1"},
				{fieldName: "deadline", fieldType: "date", alias: ""},
			},
			expectedError: nil,
		},
		{
			name:               "error structure type",
			structure:          "not a structure",
			expectedFilterules: nil,
			expectedError:      ErrStructType,
		},
		{
			name:               "invalid tag format",
			structure:          TestInvalidTag{},
			expectedFilterules: nil,
			expectedError:      ErrInvalidFilterTagFormat,
		},
		{
			name:               "field names are not unique",
			structure:          TestFieldNamesNotUnique{},
			expectedFilterules: nil,
			expectedError:      ErrFieldNamesNotUnique,
		},
		{
			name:               "aliases are not unique",
			structure:          TestAliasesNotUnique{},
			expectedFilterules: nil,
			expectedError:      ErrAliasesNotUnique,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			rules, err := getFilterRules(testCase.structure)
			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, testCase.expectedError))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedFilterules, rules)
			}
		})
	}
}

// nolint: funlen
func TestGetCondition(t *testing.T) {
	rules := filterRules{
		{fieldName: "name", fieldType: "string", alias: "name_alias"},
		{fieldName: "description", fieldType: "*string"},
		{fieldName: "age", fieldType: "int", alias: "age1"},
		{fieldName: "test_result", fieldType: "*int", alias: "testResult"},
		{fieldName: "f", fieldType: "float32"},
		{fieldName: "f_", fieldType: "*float64"},
		{fieldName: "joined_at", fieldType: "date", alias: "joinedAt"},
		{fieldName: "created_at", fieldType: "time.Time", alias: "createdAt"},
		{fieldName: "is_completed", fieldType: "bool", alias: "isCompleted"},
		{fieldName: "is_deleted", fieldType: "*bool", alias: "isDeleted"},
	}
	testTable := []struct {
		name              string
		condString        string
		expectedCondition *Condition
		expectedError     error
	}{
		{
			name:              "invalid condition format",
			condString:        "invalid condition",
			expectedCondition: nil,
			expectedError:     ErrInvalidConditionFormat,
		},
		{
			name:              "error while getting filtering rule by field name",
			condString:        "address~Gomel",
			expectedCondition: nil,
			expectedError:     ErrFieldNotFound,
		},
		{
			name:              "string 'equal' condition",
			condString:        "name==Ivan",
			expectedCondition: &Condition{FieldName: "name", Type: Equal, Value: "Ivan"},
			expectedError:     nil,
		},
		{
			name:              "string is null condition",
			condString:        "name==null",
			expectedCondition: &Condition{FieldName: "name", Type: Equal, Value: "null"},
			expectedError:     nil,
		},
		{
			name:              "*string is null condition",
			condString:        "description==null",
			expectedCondition: &Condition{FieldName: "description", Type: Equal, Value: nil},
			expectedError:     nil,
		},
		{
			name:              "string 'like' condition",
			condString:        "name~test",
			expectedCondition: &Condition{FieldName: "name", Type: Like, Value: "test"},
			expectedError:     nil,
		},
		{
			name:              "string 'in' condition",
			condString:        "name:test1,test2",
			expectedCondition: &Condition{FieldName: "name", Type: In, Value: []string{"test1", "test2"}},
			expectedError:     nil,
		},
		{
			name:              "invalid string condition",
			condString:        "name>asdf",
			expectedCondition: nil,
			expectedError:     ErrInvalidCondition,
		},
		{
			name:              "int 'equal' condition",
			condString:        "age1==50",
			expectedCondition: &Condition{FieldName: "age", Type: Equal, Value: int64(50)},
			expectedError:     nil,
		},
		{
			name:              "*int 'equal' condition",
			condString:        "testResult==NULL",
			expectedCondition: &Condition{FieldName: "test_result", Type: Equal, Value: nil},
			expectedError:     nil,
		},
		{
			name:              "int 'equal' condition - invalid value",
			condString:        "age1==50fds",
			expectedCondition: nil,
			expectedError:     ErrInvalidInteger,
		},
		{
			name:              "int 'less than or equal' condition",
			condString:        "age1<=50",
			expectedCondition: &Condition{FieldName: "age", Type: LessThanEqual, Value: int64(50)},
			expectedError:     nil,
		},
		{
			name:              "int 'less than or equal' condition - invalid value",
			condString:        "age1<=50asd",
			expectedCondition: nil,
			expectedError:     ErrInvalidInteger,
		},
		{
			name:              "int 'in' condition",
			condString:        "testResult:5,6,7",
			expectedCondition: &Condition{FieldName: "test_result", Type: In, Value: []int64{5, 6, 7}},
			expectedError:     nil,
		},
		{
			name:              "int 'in' condition - invalid value",
			condString:        "testResult:5,6,7r",
			expectedCondition: nil,
			expectedError:     ErrInvalidInteger,
		},
		{
			name:              "invalid int condition",
			condString:        "age~35",
			expectedCondition: nil,
			expectedError:     ErrInvalidCondition,
		},
		{
			name:              "float 'less than' condition",
			condString:        "f<1.5",
			expectedCondition: &Condition{FieldName: "f", Type: LessThan, Value: float64(1.5)},
			expectedError:     nil,
		},
		{
			name:              "float 'less than' condition - invalid value",
			condString:        "f<1.2aaff",
			expectedCondition: nil,
			expectedError:     ErrInvalidFloat,
		},
		{
			name:              "float 'equal' condition - invalid condition",
			condString:        "f==1.3",
			expectedCondition: nil,
			expectedError:     ErrInvalidCondition,
		},
		{
			name:              "*float 'equal' condition",
			condString:        "f_==null",
			expectedCondition: &Condition{FieldName: "f_", Type: Equal, Value: nil},
			expectedError:     nil,
		},
		{
			name:              "invalid float condition",
			condString:        "f~35",
			expectedCondition: nil,
			expectedError:     ErrInvalidCondition,
		},
		{
			name:              "field type is not supported",
			condString:        "createdAt==2022-01-01",
			expectedCondition: nil,
			expectedError:     ErrFieldTypeNotSupported,
		},
		{
			name:              "date 'equal' condition",
			condString:        "joinedAt==2022-01-01",
			expectedCondition: &Condition{FieldName: "joined_at", Type: Equal, Value: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)},
			expectedError:     nil,
		},
		{
			name:              "date 'equal' condition - invalid value",
			condString:        "joinedAt==invalid",
			expectedCondition: nil,
			expectedError:     ErrInvalidDate,
		},
		{
			name:              "invalid date condition",
			condString:        "joinedAt~2022-01-01",
			expectedCondition: nil,
			expectedError:     ErrInvalidCondition,
		},
		{
			name:              "bool 'equal' condition",
			condString:        "isCompleted==true",
			expectedCondition: &Condition{FieldName: "is_completed", Type: Equal, Value: true},
			expectedError:     nil,
		},
		{
			name:              "bool 'not equal' condition",
			condString:        "isDeleted!=null",
			expectedCondition: &Condition{FieldName: "is_deleted", Type: NotEqual, Value: nil},
			expectedError:     nil,
		},
		{
			name:              "bool 'equal' condition - invalid value",
			condString:        "isCompleted==that's right",
			expectedCondition: nil,
			expectedError:     ErrInvalidBool,
		},
		{
			name:              "invalid bool condition",
			condString:        "isCompleted>true",
			expectedCondition: nil,
			expectedError:     ErrInvalidCondition,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			condition, err := getCondition(testCase.condString, rules)
			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, testCase.expectedError))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedCondition, condition)
			}
		})
	}
}
