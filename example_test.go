package traversal_test

import (
	"encoding/json"
	"os"

	tr "github.com/dougfort/traversal"
)

type testType struct {
	A string
	B int32
	C bool
}

func ExampleStart() {
	data, _ := json.Marshal(&testType{A: "a", B: 43, C: true})

	tr.Start(data)
}

func ExampleEnd() {
	data, _ := json.Marshal(&testType{A: "a", B: 43, C: true})

	tr.Start(data).End(os.Stdout)

	// Output: {"A":"a","B":43,"C":true}
}

func ExampleObjectKey() {
	data, _ := json.Marshal(&testType{A: "a", B: 43, C: true})

	tr.Start(data).ObjectKey("B").End(os.Stdout)

	// Output: 43
}

func ExampleArraySingleton() {
	data, _ := json.Marshal([]testType{{A: "a", B: 43, C: true}})

	tr.Start(data).ArraySingleton().End(os.Stdout)

	// Output: {"A":"a","B":43,"C":true}
}

func ExampleArrayPredicate() {
	data, _ := json.Marshal([]testType{
		{A: "a", B: 43, C: true},
		{A: "a", B: 41, C: true},
		{A: "a", B: 43, C: true},
	})
	predicate := func(r json.RawMessage) bool {
		m, err := tr.GetMapFromRawMessage(r)
		if err != nil {
			return false
		}
		n, err := tr.GetInt32FromRawMessage(m["B"])
		if err != nil {
			return false
		}
		return n == 41
	}

	tr.Start(data).ArrayPredicate(predicate).End(os.Stdout)

	// Output: {"A":"a","B":41,"C":true}
}

func ExampleSelector() {
	data, _ := json.Marshal([]testType{
		{A: "a", B: 43, C: true},
		{A: "a", B: 41, C: true},
		{A: "a", B: 43, C: true},
	})
	predicate := func(r json.RawMessage) bool {
		m, err := tr.GetMapFromRawMessage(r)
		if err != nil {
			return false
		}
		n, err := tr.GetInt32FromRawMessage(m["B"])
		if err != nil {
			return false
		}
		return n == 41
	}

	tr.Start(data).ArrayPredicate(predicate).End(os.Stdout)

	// Output: {"A":"a","B":41,"C":true}
}
