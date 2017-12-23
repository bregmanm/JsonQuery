# JsonQuery
Here is an implementation of query information from text in JSON format. The standard way of marshal/unmarshall JSON text from/to GO structures is to define the structure with exactly same format as JSON. Example from https://gobyexample.com/json:
type Response2 struct {
    Page   int      `json:"page"`
    Fruits []string `json:"fruits"`
}
 str := `{"page": 1, "fruits": ["apple", "peach"]}`
 res := Response2{}
 json.Unmarshal([]byte(str), &res)
 
But what about arbitrary JSON texts?

The simple example:
	{
		"CEO":{"name":"John","Salary":10000},
		"Secretary":{"name":"Evelina","Salary":2000},
		"Others":[{"Group1":
			{"name":"Fabian","Salary":3000}},
			{"Group2":
			{"name":"Gabriel","Salary":3500}}
		]
	}

The pair "Salary":value is placed in various levels of JSON. The task is to collect all such pairs without defining the appropriate structure.
I discovered one interesting example from the same http page:
    var dat map[string]interface{}
    if err := json.Unmarshal(byt, &dat); err != nil {
	panic(err)
    }
    fmt.Println(dat) 
    
Is it possible to unmarshal arbitrary JSON text using this way? The answer is YES!
In fact design pattern Builder was implemented. The example of typical JSON query is:

result := New().SetSourceJsonText(sourceJson).SetKeyFilter("Salary").SetLimit(1).Query()

Method New() creates the query builder, SetSourceJsonText(sourceJson) loads the source JSON text string, SetKeyFilter("Salary") filters pairs "Salary":value from any JSON's level, SetLimit(1) stops processing when one pair is found and Query() executes process of search. 
The result is the slice [](string,interface{}) which contains matched pairs key:value. SetFilter is used instead of SetKeyFilter for arbitrary search. In this case user have to implement a function with arguments (key,value) and return true when pair is match. 
This function need be passed as an argument of Set Filter.
For other details please see the contents of JsonQuery.go. In case of any problems please contact with Mikhail Bregman bregmanm@mail.ru
