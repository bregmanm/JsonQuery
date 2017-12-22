package main


import (
	"fmt"
	"encoding/json"
	"reflect"
)


// Input JSON data is decoded here
var jsonTopMap = make(map[string]interface{})
var jsonTopArray = make([]interface{},0)

// Pair key:value in the result
type Entry struct {
	Key string
	Value interface{}
}

func switchType(value interface{}, jqb *jsonQueryDescriptor, result *[]Entry) {
	switch value.(type) {
	case []interface {}: // array
		recursiveQueryArrayValues(value.([]interface {}), jqb, result)
	case map[string]interface {}: // nested json
		recursiveQueryKeysValues(value.(map[string]interface {}), jqb, result)
	}
}

// Top-map JSON
func recursiveQueryKeysValues(dat map[string]interface{}, jqb *jsonQueryDescriptor, result *[]Entry) {

	for key, value := range dat {
		if jqb.limitExceed {
			break
		}
		if jqb.funcFilter(key, value) { // Pair key:value is math
			*result = append(*result, Entry{key, value})
			if jqb.limit > 0 && len(*result) >= jqb.limit { // stop on limit
				jqb.limitExceed = true
				break
			}
		}
		switchType(value, jqb, result)
	}
}

// Top-array JSON
func recursiveQueryArrayValues(arrDat []interface{}, jqb *jsonQueryDescriptor, result *[]Entry) {
	arr := reflect.ValueOf(arrDat)
	for i := 0; i < arr.Len() && !jqb.limitExceed; i++ {
		value := arr.Index(i)
		switchType(value, jqb, result)
	}
}

// Section of query builder
type JsonQueryBuilder interface {
	SetSourceJsonText(string) JsonQueryBuilder // set the source json text
	SetLimit(int) JsonQueryBuilder            // set the limit of key:value result entries, 0 - unlimited
	SetFilter(elemFilter) JsonQueryBuilder   // set the function which filters pair key:value
	SetKeyFilter(string) JsonQueryBuilder          // set the key to be filtered
	Query() []Entry                         // execute query and return map element:value from JSON
												// and result of searching: true when any JSON element was match
}

// Function which filter pair key:value
type elemFilter func (string, interface{}) bool

type jsonQueryDescriptor struct {
	json string // source text in JSON format
	limit int // limit of query results, 0 - collect all results
	limitExceed bool // flag of exceeding limit
	funcFilter elemFilter // function which filters pair key:value
	jsonTopMap map[string]interface{}
	jsonTopArray []interface{}

}

// Set the source json text
func (jqb *jsonQueryDescriptor) SetSourceJsonText(sourceJsonText string) JsonQueryBuilder {
	jqb.json = sourceJsonText
	return jqb
}

// Set stop flag on first match
func (jqb *jsonQueryDescriptor) SetLimit(limit int) JsonQueryBuilder {
	jqb.limit = limit
	return jqb
}

// Set function which filter pair key:value
func (jqb *jsonQueryDescriptor) SetFilter(funcFilter elemFilter) JsonQueryBuilder {
	jqb.funcFilter = funcFilter
	return jqb
}

// Set function which filter by key
func (jqb *jsonQueryDescriptor) SetKeyFilter(keyValue string) JsonQueryBuilder {
	jqb.funcFilter = func(s string, i interface{}) bool {
		return s == keyValue
	}
	return jqb
}

// Execute query
func (jqb *jsonQueryDescriptor) Query() []Entry {
	result := make([]Entry, 0)
	// Is JSON top-stringmap?
	jsonByteRepr := []byte(jqb.json)
	if err := json.Unmarshal(jsonByteRepr, &(jqb.jsonTopMap)); err != nil {
		// Is JSON top-array?
		if err1 := json.Unmarshal(jsonByteRepr, &(jqb.jsonTopArray)); err1 != nil {
			panic(err1) // Wrong JSON format
		} else {
			recursiveQueryArrayValues(jqb.jsonTopArray, jqb, &result)
		}
	} else {
		recursiveQueryKeysValues(jqb.jsonTopMap, jqb, &result)
	}
	return result
}

func New() JsonQueryBuilder {
	return &jsonQueryDescriptor{}
}

func main() {

	sourceJson := `
	{
		"CEO":{"name":"John","Salary":10000},
		"Secretary":{"name":"Evelina","Salary":2000},
		"Marketing":{"Group1":
			{"name":"Fabian","Salary":3000}},
		"Accounting":{"Group2":
			{"name":"Gabriel","Salary":3500}}
	}
	`
	result := New().SetSourceJsonText(sourceJson).SetKeyFilter("Salary").SetLimit(1).Query()
	fmt.Println("Result of query", result)


}

