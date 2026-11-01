package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
)

type FunctionInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Signature   string `json:"signature"`
	Example     string `json:"example"`
}

var saumFunctions = []FunctionInfo{
	{Name: "saum.submit_item", Description: "Save one actionable item in this session's user-story plan. A key is idempotent only when the content is identical. Calls commit independently.", Signature: "saum.submit_item({key: string, title: string, description: string, acceptanceCriteria: string[]}): {id, sessionId, key, title, description, acceptanceCriteria}", Example: `return saum.submit_item({key:"login-form",title:"Deliver sign-in",description:"Let existing users sign in with their credentials.",acceptanceCriteria:["Valid credentials open the private area","Invalid credentials are rejected"]});`},
	{Name: "saum.complete_plan", Description: "Seal the plan after all items were saved. Requires at least one item. Further submissions are rejected. Does not implement or accept the work.", Signature: "saum.complete_plan(): {submitted: boolean}", Example: "return saum.complete_plan();"},
}

func searchFunctions(query string) []FunctionInfo {
	words := strings.Fields(strings.ToLower(query))
	results := []FunctionInfo{}
	for _, function := range saumFunctions {
		haystack := strings.ToLower(function.Name + " " + function.Description + " " + function.Signature)
		match := len(words) == 0
		for _, word := range words {
			if strings.Contains(haystack, word) {
				match = true
				break
			}
		}
		if match {
			results = append(results, function)
		}
	}
	return results
}
func (s *Store) executeCode(ctx context.Context, sessionID, code string) (any, error) {
	if len(code) == 0 || len(code) > 32000 {
		return nil, fmt.Errorf("code must be 1–32000 bytes")
	}
	session, err := s.session(sessionID)
	if err != nil {
		return nil, err
	}
	if session.Status != SessionRunning {
		return nil, ErrPlanClosed
	}
	request, _ := json.Marshal(map[string]string{"tool": "execute", "code": code})
	if err := s.appendEvent(sessionID, "tool", request); err != nil {
		return nil, err
	}
	result, err := s.runCode(ctx, sessionID, code)
	response := map[string]any{"tool": "execute", "result": result}
	if err != nil {
		response["error"] = err.Error()
	}
	data, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		return nil, marshalErr
	}
	if saveErr := s.appendEvent(sessionID, "tool", data); saveErr != nil {
		return nil, saveErr
	}
	return result, err
}
func (s *Store) runCode(parent context.Context, sessionID, code string) (any, error) {
	ctx, cancel := context.WithTimeout(parent, 2*time.Second)
	defer cancel()
	runtime := goja.New()
	runtime.SetParserOptions(parser.WithDisableSourceMaps)
	runtime.SetMaxCallStackSize(128)
	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			runtime.Interrupt(ctx.Err())
		case <-done:
		}
	}()
	defer close(done)
	namespace := runtime.NewObject()
	calls := 0
	check := func() {
		calls++
		if calls > 100 {
			panic(runtime.NewGoError(fmt.Errorf("execution call limit reached")))
		}
		if err := ctx.Err(); err != nil {
			panic(runtime.NewGoError(err))
		}
	}
	if err := namespace.Set("submit_item", func(call goja.FunctionCall) goja.Value {
		check()
		data, err := call.Argument(0).ToObject(runtime).MarshalJSON()
		if err != nil {
			panic(runtime.NewGoError(err))
		}
		if len(data) > 64000 {
			panic(runtime.NewGoError(fmt.Errorf("item is too large")))
		}
		var input ItemSubmission
		decoder := json.NewDecoder(strings.NewReader(string(data)))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&input); err != nil {
			panic(runtime.NewGoError(err))
		}
		item, err := s.submitItem(sessionID, input)
		if err != nil {
			panic(runtime.NewGoError(err))
		}
		output, _ := json.Marshal(item)
		value, err := runtime.RunString("(" + string(output) + ")")
		if err != nil {
			panic(runtime.NewGoError(err))
		}
		return value
	}); err != nil {
		return nil, err
	}
	if err := namespace.Set("complete_plan", func(call goja.FunctionCall) goja.Value {
		check()
		if err := s.completePlan(sessionID); err != nil {
			panic(runtime.NewGoError(err))
		}
		value, err := runtime.RunString("({submitted:true})")
		if err != nil {
			panic(runtime.NewGoError(err))
		}
		return value
	}); err != nil {
		return nil, err
	}
	// Only these session-scoped functions are bound. No filesystem, network, process, or Go objects are exposed.
	if err := runtime.Set("saum", namespace); err != nil {
		return nil, err
	}
	program, err := goja.Compile("saum-execute.js", code, false)
	if err != nil {
		program, err = goja.Compile("saum-execute.js", "(function(){\n"+code+"\n})()", false)
	}
	if err != nil {
		return nil, err
	}
	value, err := runtime.RunProgram(program)
	if err != nil {
		return nil, err
	}
	if goja.IsUndefined(value) || goja.IsNull(value) {
		return nil, nil
	}
	if err := runtime.Set("__saum_result", value); err != nil {
		return nil, err
	}
	serialized, err := runtime.RunString("JSON.stringify(__saum_result)")
	if err != nil {
		return nil, err
	}
	if goja.IsUndefined(serialized) {
		return nil, fmt.Errorf("return a JSON-serializable value")
	}
	data := serialized.String()
	if len(data) > 1<<20 {
		return nil, fmt.Errorf("execution result exceeds 1 MiB")
	}
	var result any
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}
	return result, nil
}
