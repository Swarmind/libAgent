package rewoo

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/Swarmind/libagent/internal/tools"

	graph "github.com/JackBekket/langgraphgo/graph/stategraph"
	"github.com/rs/zerolog/log"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

const (
	GraphPlanName  = "plan"
	GraphToolName  = "tool"
	GraphSolveName = "solve"
)

const PromptGetPlan = `For the following task, make plans that can solve the problem step by step. For each plan, indicate
which external tool together with tool input to retrieve evidence. You can store the evidence into a
variable #E that can be called by later tools. (Plan, #E1, Plan, #E2, Plan, ...)
Each step is context isolated and need to be explicitly provided with evidence variable or task context details if needed.
You need to strictly stick to the output format as in the example below, as your output will be parsed using regex match for the future use.

Example input:
	List of tools:
	(1) search[json: {"query": "string"}]: Worker that searches results from Duckduckgo. Useful when you need to find short
	and succinct answers about a specific topic. The input should be a search query.
	(2) LLM[string]: A pretrained LLM like yourself. Useful when you need to act with general
	world knowledge and common sense. Prioritize it when you are confident in solving the problem
	yourself. Input can be any instruction.

	Task: Thomas, Toby, and Rebecca worked a total of 157 hours in one week. Thomas worked x
	hours. Toby worked 10 hours less than twice what Thomas worked, and Rebecca worked 8 hours
	less than Toby. How many hours did Rebecca work?

Example output:
	Plan: Given Thomas worked x hours, translate the problem into algebraic expressions and solve
	with Wolfram Alpha. #E1 = WolframAlpha[{"query": "Solve x + (2x − 10) + ((2x − 10) − 8) = 157"}]
	Plan: Find out the number of hours Thomas worked. #E2 = LLM[What is x, given #E1]
	Plan: Calculate the number of hours Rebecca worked. #E3 = Calculator[{"query": "(2 ∗ #E2 − 10) − 8"}]

Begin! 
Describe your plans with rich details. Each Plan should be followed by only one #E.

`

const PromptSolver = `Solve the following task or problem. To solve the problem, we have made step-by-step Plan and \
retrieved corresponding Evidence to each Plan. Use them with caution since long evidence might \
contain irrelevant information.

%s

Now solve the question or task according to provided Evidence above. Respond with the answer
directly with no extra words.

Task: %s
Response:`

const PromptLLMTool = `Do not include any introductory phrases or explanations. Task: %s`
const PromptCallTool = `You will be provided with tool name and arguments for the call.
	Use tool description and plan to decide how to resolve the provided arguments into the tool call schema.
	Try to sanitize arguments, resolve possible string concatenation.
	Plan: %s
	Tool name: %s
	Tool description: %s
	Arguments: %s
`

type ReWOO struct {
	LLM           *openai.LLM
	ToolsExecutor *tools.ToolsExecutor
}

type State struct {
	Task       string
	PlanString string
	Steps      []Step
	Results    map[string]string
	Result     string
}

type Step struct {
	Plan      string
	Name      string
	Tool      string
	ToolInput string
}

var StepPattern *regexp.Regexp = regexp.MustCompile(
	`Plan:\s*(.+)\s*(#E\d+)\s*=\s*(\w+)\s*\[([^\]]+)\]`,
)

func (r ReWOO) InitializeGraph() (*graph.Runnable, error) {
	workflowGraph := graph.NewStateGraph()

	workflowGraph.AddNode("plan", r.GetPlan)
	workflowGraph.AddNode("tool", r.ToolExecution)
	workflowGraph.AddNode("solve", r.Solve)
	workflowGraph.AddEdge("plan", "tool")
	workflowGraph.AddEdge("solve", graph.END)
	workflowGraph.AddConditionalEdge("tool", r.Route)
	workflowGraph.SetEntryPoint("plan")
	return workflowGraph.Compile()
}

func (r ReWOO) GetPlan(ctx context.Context, s interface{}) (interface{}, error) {
	state := s.(State)
	task := state.Task

	response, err := r.LLM.GenerateContent(ctx,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeHuman,
				fmt.Sprintf(
					"%s\nList of tools:\n%s\nTask:\n```\n%s```",
					PromptGetPlan,
					r.ToolsExecutor.ToolsPromptDesc(),
					task,
				),
			)},
	)
	if err != nil {
		return s, err
	}

	result := response.Choices[0].Content
	matches := StepPattern.FindAllStringSubmatch(result, -1)
	if matches == nil {
		return s, fmt.Errorf("empty plan matches")
	}

	sortedKeys := []string{}
	// using map approach, as think models can double the step, and the last match is preferred.
	stepMap := map[string]Step{}
	for _, m := range matches {
		stepMap[m[2]] = Step{
			// m[0] - full match,
			Plan:      m[1],
			Name:      m[2],
			Tool:      m[3],
			ToolInput: m[4],
		}
		if !slices.Contains(sortedKeys, m[2]) {
			sortedKeys = append(sortedKeys, m[2])
		}
	}
	for _, key := range sortedKeys {
		state.Steps = append(state.Steps, stepMap[key])
	}

	state.PlanString = result

	log.Debug().
		Interface("state.Steps", state.Steps).
		Msg("ReWOO: GetPlan")

	return state, nil
}

func (r ReWOO) Solve(ctx context.Context, s interface{}) (interface{}, error) {
	state := s.(State)

	plan := ""
	for _, step := range state.Steps {
		for stepName, result := range state.Results {
			step.ToolInput = strings.ReplaceAll(step.ToolInput, stepName, result)
			step.Name = strings.ReplaceAll(step.Name, stepName, result)
		}
		plan += fmt.Sprintf(
			"Plan: %s\n%s = %s[%s]\n",
			step.Plan,
			step.Name,
			step.Tool,
			step.ToolInput,
		)
	}
	response, err := r.LLM.GenerateContent(ctx,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeHuman,
				fmt.Sprintf(PromptSolver, plan, state.Task),
			)},
	)
	if err != nil {
		return state, err
	}

	state.Result = response.Choices[0].Content
	log.Debug().
		Str("state.Result", state.Result).
		Msg("ReWOO: Solve")

	return state, nil
}

func (r ReWOO) ToolExecution(ctx context.Context, s interface{}) (interface{}, error) {
	state := s.(State)

	step := state.Steps[getCurrentTask(state)]

	for stepName, result := range state.Results {
		step.ToolInput = strings.ReplaceAll(step.ToolInput, stepName, result)
	}

	prompt := fmt.Sprintf(PromptLLMTool, step.ToolInput)
	options := []llms.CallOption{}
	content := ""
	if step.Tool != "LLM" {
		toolDesc := ""
		for _, tool := range r.ToolsExecutor.Tools {
			if tool.Definition.Name == step.Tool {
				options = append(options, llms.WithTools([]llms.Tool{{
					Type:     "function",
					Function: &tool.Definition,
				},
				}))
				toolDesc = tool.Definition.Description
			}
		}

		prompt = fmt.Sprintf(
			PromptCallTool,
			step.Plan,
			step.Tool,
			toolDesc,
			step.ToolInput,
		)
	}

	log.Debug().
		Str("name", step.Name).
		Str("tool", step.Tool).
		Str("prompt", prompt).
		Msg("ReWOO: ToolExecution pre-GenerateContent")

	response, err := r.LLM.GenerateContent(ctx,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeHuman,
				prompt,
			)},
		options...,
	)
	if err != nil {
		return state, err
	}
	content = response.Choices[0].Content
	if toolContent := r.ToolsExecutor.ProcessToolCalls(
		ctx, response.Choices[0].ToolCalls,
	); toolContent != "" {
		content = toolContent
	}
	log.Debug().
		Str("name", step.Name).
		Str("tool", step.Tool).
		Str("prompt", prompt).
		Str("content", content).
		Msg("ReWOO: ToolExecution")

	if len(state.Results) == 0 {
		state.Results = map[string]string{}
	}
	jsonSafeContent, err := json.Marshal(content)
	if err != nil {
		return state, err
	}

	state.Results[step.Name] = string(jsonSafeContent)
	return state, nil
}

func (_ ReWOO) Route(ctx context.Context, state interface{}) string {
	if getCurrentTask(state.(State)) == -1 {
		return GraphSolveName
	} else {
		return GraphToolName
	}
}

func getCurrentTask(state State) int {
	if len(state.Results) == len(state.Steps) {
		return -1
	}
	return len(state.Results)
}
