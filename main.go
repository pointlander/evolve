// Copyright 2025 The Evolve Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
)

// Prompt is a llm prompt
type Prompt struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// Query submits a query to the llm
func Query(query string) string {
	prompt := Prompt{
		Model:  "gpt-oss",
		Prompt: query,
	}
	data, err := json.Marshal(prompt)
	if err != nil {
		panic(err)
	}
	buffer := bytes.NewBuffer(data)
	response, err := http.Post("http://10.0.0.54:11434/api/generate", "application/json", buffer)
	if err != nil {
		panic(err)
	}
	reader, answer := bufio.NewReader(response.Body), ""
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		data := map[string]interface{}{}
		err = json.Unmarshal([]byte(line), &data)
		text := data["response"].(string)
		answer += text
	}
	return answer
}

// Variable is a variable
type Variable struct {
	Variable int
	Negation bool
}

// Clause is a 3sat clause
type Clause [3]Variable

// Clauses are a set of clauses
type Clauses []Clause

func main() {
	//fmt.Println(Query("Hello world!"))
	//fmt.Println(Query("How does llm based evolution work?"))
	//fmt.Println(Query("How to encode 3sat for gpt?"))
	sat := make(Clauses, 0, 8)
	sat = append(sat,
		Clause{{0, false}, {2, true}, {4, false}},
		Clause{{0, true}, {1, false}, {3, false}},
		Clause{{2, false}, {1, true}, {5, true}},
		Clause{{4, true}, {3, true}, {5, false}},
	)
	variables := make(map[int]bool)
	for _, clause := range sat {
		for _, variable := range clause {
			variables[variable.Variable] = true
		}
	}
	str := fmt.Sprintf("Problem: 3‑SAT instance with %d variables and %d clauses\n\n", len(variables), len(sat))
	list := make([]int, 0, 8)
	for variable := range variables {
		list = append(list, variable)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i] < list[j]
	})
	str += fmt.Sprintf("Variables: ")
	for i, variable := range list {
		if i == 0 {
			str += fmt.Sprintf("x%d", variable)
			continue
		}
		str += fmt.Sprintf(", x%d", variable)
	}
	str += fmt.Sprintf("\n\n")
	str += fmt.Sprintf("Clauses:\n")
	for i, clause := range sat {
		str += fmt.Sprintf("C%d: ", i)
		for ii, variable := range clause {
			if ii == 0 {
				if variable.Negation {
					str += fmt.Sprintf("¬")
				}
				str += fmt.Sprintf("x%d", variable.Variable)
				continue
			}
			str += fmt.Sprintf(" ∧ ")
			if variable.Negation {
				str += fmt.Sprintf("¬")
			}
			str += fmt.Sprintf("x%d", variable.Variable)
		}
		str += fmt.Sprintf("\n")
	}
	str += fmt.Sprintf("Given the 3‑SAT instance above, answer the following:\n\n")
	str += fmt.Sprintf("1. Is the instance satisfiable? (Yes/No)\n")
	str += fmt.Sprintf("2. If Yes, give a satisfying assignment in the form:\n")
	str += fmt.Sprintf("   ")
	for i, variable := range list {
		if i == 0 {
			str += fmt.Sprintf("x%d=%d", variable, i&1)
			continue
		}
		str += fmt.Sprintf(", x%d=%d", variable, i&1)
	}
	str += fmt.Sprintf("\n")
	str += fmt.Sprintf("   (use 1 for true, 0 for false)")
	fmt.Println(str)
	fmt.Println(Query(str))
}
