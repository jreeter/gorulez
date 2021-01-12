package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type Rule struct {
	Condition *Condition
	Action    *Action
}

type Condition struct {
	All      []Condition
	Any      []Condition
	Name     string
	Operator string
	Value    interface{}
}

type Action struct {
	Name   string
	Params interface{}
}

type Target struct {
	Units string
}

type Comparator func(left interface{}, right interface{}) bool

func (c *Condition) getComparator(targetType string, evalOperator string) Comparator {

	Operators := make(map[string]map[string]Comparator)

	StringOperators := make(map[string]Comparator)
	StringOperators["EqualTo"] = func(a string, b string) bool { return a == b }
	StringOperators["LessThan"] = func(a string, b string) bool { return a < b }
	StringOperators["GreaterThan"] = func(a string, b string) bool { return a > b }
	StringOperators["NotEqual"] = func(a string, b string) bool { return a != b }

	Operators["string"] = StringOperators
	//Operators["...."] = ...

	return Operators[targetType][evalOperator]
}

func evaluate(c Condition, target interface{}) (bool, error) {

	all := len(c.All)
	any := len(c.Any)

	if all > 0 && any > 0 {
		return false, errors.New("all and Any cannot be set at same time")
	}

	// All conditions eval'd must be true, if not, return false
	if all > 0 {
		for _, condition := range c.All {
			result, err := evaluate(condition, target)
			if err != nil {
				return false, err
			}
			if !result {
				return false, nil
			}
		}
		return true, nil
	}

	if any > 0 {
		for _, condition := range c.Any {
			result, err := evaluate(condition, target)
			if err != nil {
				return false, err
			}
			if result {
				return true, nil
			}
		}
		return false, nil
	}

	// If we get here, we are no longer recursing through all or any, and now
	// need to eval the name, operator, and the value to our target.

	resolvedTarget := reflect.ValueOf(target)
	resolvedValue := resolvedTarget.FieldByName(c.Name).Interface().(string)
	comparator := c.getComparator(reflect.TypeOf(resolvedValue).String(), c.Operator)
	return comparator(c.Value.(string), resolvedValue), nil

}

func main() {
	str := `{"all":[{"Name":"Units", "Operator":"EqualTo", "Value":"100"}, {"Name":"Units", "Operator": "EqualTo", "Value":"5"}], "any": []}`
	res := Condition{}
	json.Unmarshal([]byte(str), &res)
	result, err := evaluate(res, Target{Units: "100"})
	if err != nil {
		fmt.Println("error evaluating condition")
	} else {
		fmt.Println(result)
	}
}
