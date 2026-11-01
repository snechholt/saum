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
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Signature   string      `json:"signature"`
	Example     string      `json:"example"`
	Kind        SessionKind `json:"-"`
}

var saumFunctions = []FunctionInfo{
	{Kind: ImplementationSession, Name: "saum.work_item.run_check", Description: "Queue a verification command in a separate snapshot checkout with normal host permissions, networking, and caches. Poll checks() for its recorded result; no model call is made.", Signature: "saum.work_item.run_check({command:string[]}): {id,status}", Example: `return saum.work_item.run_check({command:["go","test","./..."]});`},
	{Kind: ImplementationSession, Name: "saum.work_item.checks", Description: "Get recorded checks and requiredCommands from earlier approved items. Rerun those regressions against the current source before requesting review.", Signature: "saum.work_item.checks(): {checks:Check[],requiredCommands:string[][]}", Example: "return saum.work_item.checks();"},
	{Kind: ImplementationSession, Name: "saum.work_item.submit_acceptance_criteria", Description: "Submit verification evidence for ONE acceptance criterion of the assigned work item, in criterion order. This is not human approval.", Signature: "saum.work_item.submit_acceptance_criteria({criterionIndex:number,summary:string,evidence:string[],checkIds:string[]}): {submitted:boolean}", Example: `return saum.work_item.submit_acceptance_criteria({criterionIndex:0,summary:"Valid login works",evidence:["Recorded auth check passed"],checkIds:["id-from-run_check"]});`},
	{Kind: ImplementationSession, Name: "saum.work_item.request_review", Description: "Request a factory-agent executive summary after every criterion has evidence. Then stop. Only human approval of that summary unlocks the next item.", Signature: "saum.work_item.request_review(): {requested:boolean}", Example: "return saum.work_item.request_review();"},
	{Kind: ImplementationSession, Name: "saum.work_item.report_blocker", Description: "Report the exact missing decision or dependency blocking this assigned item, then stop.", Signature: "saum.work_item.report_blocker({reason:string}): {blocked:boolean}", Example: `return saum.work_item.report_blocker({reason:"Which verified recovery channel should be used?"});`},
	{Kind: SummarySession, Name: "saum.work_item.submit_executive_summary", Description: "The factory main agent submits a human-readable executive summary, verification assessment, and risks for the assigned work item. Does not approve it.", Signature: "saum.work_item.submit_executive_summary({summary:string,verification:string,risks:string[]}): {submitted:boolean}", Example: `return saum.work_item.submit_executive_summary({summary:"Users can sign in securely.",verification:"Observed 4 passing auth tests in the goal transcript.",risks:["Recovery is a later item."]});`},
	{Name: "saum.submit_item", Description: "Save one actionable item in this session's user-story plan. A key is idempotent only when the content is identical. Calls commit independently.", Signature: "saum.submit_item({key: string, title: string, description: string, acceptanceCriteria: string[]}): {id, sessionId, key, title, description, acceptanceCriteria}", Example: `return saum.submit_item({key:"login-form",title:"Deliver sign-in",description:"Let existing users sign in with their credentials.",acceptanceCriteria:["Valid credentials open the private area","Invalid credentials are rejected"]});`},
	{Name: "saum.update_item", Description: "Replace an existing item in the current working plan. Provide every field, preserving its key. Never changes an earlier submitted plan.", Signature: "saum.update_item({key: string, title: string, description: string, acceptanceCriteria: string[]}): Item", Example: `return saum.update_item({key:"login-form",title:"Deliver sign-in",description:"Use the revised behavior.",acceptanceCriteria:["Revised behavior works"]});`},
	{Name: "saum.remove_item", Description: "Remove an item by key from the current working plan only.", Signature: "saum.remove_item({key: string}): {removed: boolean}", Example: `return saum.remove_item({key:"unnecessary-item"});`},
	{Name: "saum.resolve_feedback", Description: "Record how a feedback comment assigned to this revision was addressed. Every assigned comment must be resolved before completing the plan. Human approval is still required.", Signature: "saum.resolve_feedback({commentId: string, resolution: string}): {resolved: boolean}", Example: `return saum.resolve_feedback({commentId:"id-from-feedback",resolution:"Updated the item and its acceptance criteria to match the requested behavior."});`},
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
	workItemNamespace := runtime.NewObject()
	if err := namespace.Set("work_item", workItemNamespace); err != nil {
		return nil, err
	}
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
	bind := func(name string, handler func([]byte) (any, error)) error {
		target := namespace
		if strings.HasPrefix(name, "work_item.") {
			target = workItemNamespace
			name = strings.TrimPrefix(name, "work_item.")
		}
		return target.Set(name, func(call goja.FunctionCall) goja.Value {
			check()
			data := []byte("{}")
			if !goja.IsUndefined(call.Argument(0)) {
				var err error
				data, err = call.Argument(0).ToObject(runtime).MarshalJSON()
				if err != nil {
					panic(runtime.NewGoError(err))
				}
			}
			if len(data) > 64000 {
				panic(runtime.NewGoError(fmt.Errorf("function input is too large")))
			}
			result, err := handler(data)
			if err != nil {
				panic(runtime.NewGoError(err))
			}
			output, err := json.Marshal(result)
			if err != nil {
				panic(runtime.NewGoError(err))
			}
			value, err := runtime.RunString("(" + string(output) + ")")
			if err != nil {
				panic(runtime.NewGoError(err))
			}
			return value
		})
	}
	decode := func(data []byte, input any) error {
		decoder := json.NewDecoder(strings.NewReader(string(data)))
		decoder.DisallowUnknownFields()
		return decoder.Decode(input)
	}
	handlers := []struct {
		name string
		call func([]byte) (any, error)
	}{

		{"work_item.run_check", func(data []byte) (any, error) {
			var input CheckRequest
			if err := decode(data, &input); err != nil {
				return nil, err
			}
			return s.requestVerification(sessionID, input)
		}},
		{"work_item.checks", func(data []byte) (any, error) { return s.checkReport(sessionID) }},
		{"work_item.submit_acceptance_criteria", func(data []byte) (any, error) {
			var input CriterionInput
			if err := decode(data, &input); err != nil {
				return nil, err
			}
			err := s.submitAcceptanceCriterion(sessionID, input)
			return map[string]bool{"submitted": err == nil}, err
		}},
		{"work_item.request_review", func(data []byte) (any, error) {
			err := s.requestWorkItemReview(sessionID)
			return map[string]bool{"requested": err == nil}, err
		}},
		{"work_item.report_blocker", func(data []byte) (any, error) {
			var input struct {
				Reason string `json:"reason"`
			}
			if err := decode(data, &input); err != nil {
				return nil, err
			}
			err := s.blockWorkItem(sessionID, input.Reason)
			return map[string]bool{"blocked": err == nil}, err
		}},
		{"work_item.submit_executive_summary", func(data []byte) (any, error) {
			var input ExecutiveSummaryInput
			if err := decode(data, &input); err != nil {
				return nil, err
			}
			err := s.submitExecutiveSummary(sessionID, input)
			return map[string]bool{"submitted": err == nil}, err
		}},
		{"submit_item", func(data []byte) (any, error) {
			var input ItemSubmission
			if err := decode(data, &input); err != nil {
				return nil, err
			}
			return s.submitItem(sessionID, input)
		}},
		{"update_item", func(data []byte) (any, error) {
			var input ItemSubmission
			if err := decode(data, &input); err != nil {
				return nil, err
			}
			return s.updateItem(sessionID, input)
		}},
		{"remove_item", func(data []byte) (any, error) {
			var input struct {
				Key string `json:"key"`
			}
			if err := decode(data, &input); err != nil {
				return nil, err
			}
			err := s.removeItem(sessionID, input.Key)
			return map[string]bool{"removed": err == nil}, err
		}},
		{"resolve_feedback", func(data []byte) (any, error) {
			var input FeedbackResolution
			if err := decode(data, &input); err != nil {
				return nil, err
			}
			err := s.resolveFeedback(sessionID, input)
			return map[string]bool{"resolved": err == nil}, err
		}},
		{"complete_plan", func(data []byte) (any, error) {
			err := s.completePlan(sessionID)
			return map[string]bool{"submitted": err == nil}, err
		}},
	}
	for _, handler := range handlers {
		if err := bind(handler.name, handler.call); err != nil {
			return nil, err
		}
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
