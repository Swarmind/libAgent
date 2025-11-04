# examples/codemonkey/cmd/codemonkey/main.go  
# Package/Component    
**main**  
  
## Imports  
- `fmt` – standard Go formatting package.  
- `github.com/Swarmind/libagent/examples/codemonkey/pkg/executor` – provides command execution utilities.  
- `github.com/Swarmind/libagent/examples/codemonkey/pkg/planner` – contains planning logic for GitHub issues.  
  
## External Data / Input Sources  
| Source | Description |  
|--------|-------------|  
| `githubservice.EventsService` (commented) | A channel of issue events that will be consumed to build a task. |  
| `reviewer.GatherInfo` (commented) | Builds a task from an issue’s text and repository name. |  
| `planner.PlanCLIExecutor` (active) | Generates a plan based on a CLI‑style description string. |  
  
## TODOs  
- Implement the commented “issue flow” section to consume real GitHub events.  
- Add the missing imports for `githubservice` and `reviewer` once those parts are activated.  
  
## Summary of Major Code Parts  
  
### 1. Issue Flow (commented)  
```go  
/* es := &githubservice.EventsService{  
    GithubAPI: githubservice.ConstructGithubApi(),  
    Ichan:     make(chan githubservice.IssueEvent, 10),  
}  
  
go es.StartWebhookHandler(utility.GetEnv("LISTEN_ADDR"))  
  
for issue := range es.Ichan {  
    fmt.Printf("Got issue: %s\n", issue.RepoName)  
  
    task := reviewer.GatherInfo(issue.IssueText, issue.RepoName)  
    task = util.RemoveThinkTag(task)  
    fmt.Println("Reviewer result: ", task)  
}  
*/  
```  
*Purpose*: Set up a GitHub event service, start a webhook handler, and iterate over incoming issues to build tasks. The code is currently commented out but shows the intended flow.  
  
### 2. Test Stuff (commented)  
```go  
/* task := reviewer.GatherInfo("Change hello message to Can I haz cheeseburger?", "Hellper")  
task = util.RemoveThinkTag(task)  
plan := planner.PlanGitHelper(task)  
fmt.Println("Planner result: ", plan)  
*/  
```  
*Purpose*: Quick manual test of the `reviewer` and `planner` components. It is also commented out.  
  
### 3. Active Plan Execution  
```go  
plan := planner.PlanCLIExecutor(`Find current OS version and OS type (windows/linux/android)`)  
fmt.Println(plan)  
executer.ExecuteCommands(plan)  
```  
*Purpose*: The only active code in `main`.    
1. Calls `planner.PlanCLIExecutor` with a string that describes the desired action: detecting the current operating system and its type.    
2. Prints the resulting plan to stdout.    
3. Passes the plan to `executer.ExecuteCommands`, which will run the generated commands.  
  
The flow is straightforward: generate a plan from a CLI description, print it, then execute it.  
  
---  
  
