package main

import (
	"fmt"
	"encoding/json"
	"reflect"
)

type Rule struct {
	Condition *Condition
	Action    *Action
}

type Condition struct {
	All  []Condition
	Any	 []Condition
	Name 	 string
  	Operator string
	Value	 interface{}
}

type Action struct {
	Name	string
	Params	interface{}
}

type Target struct {
	Units string
}


func EvaluateConditionRecursively(c Condition, target interface{}) bool {

	StringOperators := make(map[string]func(string,string)bool)
	StringOperators["EqualTo"] = func(a string, b string) bool { return a == b }

	all := len(c.All)
	any := len(c.Any)

	// Can't have both any and all... for now?
	if all > 0 && any > 0 {
		return false
	}

	// All conditions eval'd must be true, if not, return false
	if all > 0 {
		for _, condition := range c.All {
			result := EvaluateConditionRecursively(condition, target)
			if !result { return false }
		}
		return true
	}

	if any > 0 {
		for _, condition := range c.Any {
			result := EvaluateConditionRecursively(condition, target)
			if result { return true }
		}
		return false
	}

	// If we get here, we are no longer recursing through all or any, and now
	// need to eval the name, operator, and the value to our target.

	resolvedTarget := reflect.ValueOf(target)
	resolvedValue	:= resolvedTarget.FieldByName(c.Name).Interface().(string)
	comparator	:= StringOperators[c.Operator]
	return comparator(c.Value.(string), resolvedValue)

}

func main() {
	str := `{"any":[{"Name":"Units", "Operator":"EqualTo", "Value":"100"}, {"Name":"Units", "Operator": "EqualTo", "Value":"5"}]}`
	res := Condition{}
	json.Unmarshal([]byte(str), &res)
	result := EvaluateConditionRecursively(res, Target{Units: "100"})
	fmt.Println(result)
}

