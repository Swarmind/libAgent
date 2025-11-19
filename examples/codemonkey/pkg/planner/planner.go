package planner

import (
	"context"
	"encoding/json"

	"os"

	"github.com/Swarmind/libagent/pkg/config"
	"github.com/Swarmind/libagent/pkg/tools"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const lePromptGithelper = `Role: You are an Instruction Synthesis Agent. Your task is to transform a code review summary into a precise, executable action plan for a command-executor agent.

Input: You will be given a "Reviewer Result" text block containing an issue summary, desired outcome, relevant information, affected files, and code analysis.

Output Instructions: Transform the input into a structured guide using the following exact template. Do not deviate from this structure.
text

### **EXECUTION PLAN**

**1. OBJECTIVE:**
[Concise, one-sentence description of the goal, copied from 'Desired Outcome'.]

**2. CONTEXT:**
[Bulleted list summarizing the 'Relevant Information'. Rephrase for clarity and brevity. This helps the executor understand the *why*.]

**3. AFFECTED FILES:**
[List the full file paths, one per line, exactly as provided in the 'Affected Files' section.]

**4. ACTION: REPLACE CODE BLOCK**
- **File:** [The primary file to edit, e.g., "dialog.go"]
- **Search for the following exact lines:**

[Paste the exact code lines from the 'Code Analysis' section that need to be changed. Include the line comment (// >) if present.]
text

- **Replace with:**

[Provide the exact new code lines as specified in the 'Code Analysis' or 'Desired Outcome'.]
text


**5. VERIFICATION:**
- Run a syntax check or linter specific to the project's language (e.g., "go fmt <filepath>", "python -m py_compile <filepath>").
- If syntax check fails, return the instruction with error text added at the last new line of it.

Example Input:
text

Reviewer result:  Issue Summary: The "hello" message in the Hellper bot is hardcoded and needs modification.
Desired Outcome: Replace the existing "hello" message with "Can I haz cheeseburger?".
Relevant Information:
- The "hello" message is defined in "var msgTemplates" within the "command" package.
- The code is part of a Telegram bot interacting with an AI endpoint.
Affected Files:
- "dialog.go"
- "lib/bot/dialog/dialog.go"
Code Analysis:
// > Line in "dialog.go":
"hello": "Hey, this bot is working with LocalAI node! Please input your local-ai api_key 🐱",
// Replace the value of the "hello" key with "Can I haz cheeseburger?".

Example Output using the Template:
text

### **EXECUTION PLAN**

**1. OBJECTIVE:**
Replace the existing "hello" message with "Can I haz cheeseburger?".

**2. CONTEXT:**
- The hello message is a hardcoded string in a message template variable.
- The change is for a Telegram bot's command package.

**3. AFFECTED FILES:**
dialog.go
lib/bot/dialog/dialog.go

**4. ACTION: REPLACE CODE BLOCK**
- **File:** dialog.go
- **Search for the following exact lines:**

"hello": "Hey, this bot is working with LocalAI node! Please input your local-ai api_key 🐱",
text

- **Replace with:**

"hello": "Can I haz cheeseburger?",


**5. VERIFICATION:**
- Run a syntax check on the modified file according to the language used
`

var lePromptCLI = `You are an AI command generation assistant specialized in creating executable CLI command sequences. Your task is to analyze the given objective and output ONLY the CLI commands needed to accomplish it.

RULES:
1. Output ONLY valid Unix/Linux CLI commands, nothing else
2. Separate multiple commands with newlines (\n)
3. Use proper command sequencing with && when commands depend on each other
4. Include all necessary flags and options for precise execution
5. Ensure commands are safe, non-destructive, and non-interactive
6. Use standard Unix/Linux commands that are widely available
7. If verification is needed, include appropriate check commands
8. Escape special characters properly for shell execution
9. Do not use any programming languages or references to them, only CLI.

OUTPUT FORMAT:
- Output ONLY the raw commands, separated by newlines
- No explanations, no step numbers, no descriptions
- Multiple related commands can be chained with && on one line
- Each discrete operation should be on its own line

EXAMPLE 1:
For objective "list files in current directory":
ls -la

EXAMPLE 2:
For objective "find all .txt files and count lines":
find . -name "*.txt" -type f
xargs wc -l

EXAMPLE 3:
For objective "create directory and file":
mkdir -p new_directory && cd new_directory && touch new_file.txt

Generate ONLY the CLI commands needed for the following objective: 

`

func PlanGitHelper(review string) string {

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("new config")
	}

	ctx := context.Background()

	toolsToWhitelist := []string{
		tools.ReWOOToolDefinition.Name,
	}

	toolsExecutor, err := tools.NewToolsExecutor(ctx, cfg, tools.WithToolsWhitelist(toolsToWhitelist...))
	if err != nil {
		log.Fatal().Err(err).Msg("new tools executor")
	}
	defer func() {
		if err := toolsExecutor.Cleanup(); err != nil {
			log.Fatal().Err(err).Msg("tools executor cleanup")
		}
	}()

	rewooQuery := tools.ReWOOToolArgs{
		Query: (lePromptGithelper + "\n" + review),
	}
	rewooQueryBytes, err := json.Marshal(rewooQuery)
	if err != nil {
		log.Fatal().Err(err).Msg("json marhsal rewooQuery")
	}

	result, err := toolsExecutor.CallTool(ctx,
		tools.ReWOOToolDefinition.Name,
		string(rewooQueryBytes),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("rewoo tool call")
	}

	if result == "" {
		log.Fatal().Msg("main empty result")
	}

	return result

}

func PlanCLIExecutor(task string) string {

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("new config")
	}

	ctx := context.Background()

	toolsToWhitelist := []string{
		tools.ReWOOToolDefinition.Name,
	}

	toolsExecutor, err := tools.NewToolsExecutor(ctx, cfg, tools.WithToolsWhitelist(toolsToWhitelist...))
	if err != nil {
		log.Fatal().Err(err).Msg("new tools executor")
	}
	defer func() {
		if err := toolsExecutor.Cleanup(); err != nil {
			log.Fatal().Err(err).Msg("tools executor cleanup")
		}
	}()

	rewooQuery := tools.ReWOOToolArgs{
		Query: (lePromptCLI + "\n" + task),
	}
	rewooQueryBytes, err := json.Marshal(rewooQuery)
	if err != nil {
		log.Fatal().Err(err).Msg("json marhsal rewooQuery")
	}

	result, err := toolsExecutor.CallTool(ctx,
		tools.ReWOOToolDefinition.Name,
		string(rewooQueryBytes),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("rewoo tool call")
	}

	if result == "" {
		log.Fatal().Msg("main empty result")
	}

	return result

}
