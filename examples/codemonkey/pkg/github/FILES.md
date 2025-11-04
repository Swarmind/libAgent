# examples/codemonkey/pkg/github/githubservice.go  
**Package / Component name**    
`githubservice`  
  
---  
  
### Imports  
```go  
import (  
	"context"  
	"net/http"  
	"os"  
	"strconv"  
  
	utility "github.com/Swarmind/libagent/examples/codemonkey/pkg/util"  
  
	"github.com/cbrgm/githubevents/v2/githubevents"  
	"github.com/google/go-github/v74/github"  
	"github.com/rs/zerolog/log"  
)  
```  
  
---  
  
## External data / input sources  
| Env variable | Purpose |  
|--------------|---------|  
| `WEBHOOK_SECRET_KEY` | Secret key for GitHub webhook authentication (used in `githubevents.New`) |  
| `WEBHOOK_ROUTE` | HTTP route where the webhook handler will be mounted (`http.HandleFunc`) |  
| `LISTEN_ADDR` | Address on which the HTTP server listens (`ListenAndServe`) |  
| `APP_ID` | Application ID for GitHub API (parsed into `GithubAPI.AppId`) |  
| `PRIVKEY_PATH` | Path to the private key file used by the GitHub client |  
  
---  
  
## TODOs  
No explicit `TODO:` comments are present in this file.  
  
---  
  
## Types  
  
| Type | Description |  
|------|-------------|  
| `GithubAPI` | Holds configuration for the GitHub API: application ID and path to a private key. |  
| `EventsService` | Encapsulates a `GithubAPI` instance and an event channel (`Ichan`) that receives parsed issue events. |  
| `IssueEvent` | Simple DTO containing the repository name and the text of an issue (title + body). |  
  
---  
  
## Functions  
  
### `ConstructGithubApi`  
* Reads `APP_ID` from environment, parses it as `int64`, logs fatal on error.  
* Reads `PRIVKEY_PATH` from environment, verifies file existence.  
* Returns a fully populated `GithubAPI` struct.  
  
### `StartWebhookHandler(port string)`  
* Creates a new GitHub events handler with the secret key (`WEBHOOK_SECRET_KEY`).  
* Registers an event callback for *any* issue event via `handle.OnIssuesEventAny(es.IssuesEventAnyHandler)`.  
* Mounts the handler on the HTTP route defined by `WEBHOOK_ROUTE` using `http.HandleFunc`.    
  The closure forwards incoming requests to the webhook handler and logs any error.  
* Starts an HTTP server listening on the address from `LISTEN_ADDR`; fatal log on failure.  
  
### `IssuesEventAnyHandler`  
* Callback invoked for every issue event received by the webhook.  
* Extracts the issue data (`event.GetIssue()`), builds a new `IssueEvent` containing:  
  * Repository full name (`event.GetRepo().GetFullName()`)  
  * Issue title + body concatenated with a newline.  
* Sends the constructed `IssueEvent` into the service’s channel (`es.Ichan`) and returns nil.  
  
### `GetEnv`  
* Helper that fetches an environment variable by key, logs fatal if empty, and returns its value.    
  Used throughout the file for configuration values.  
  
---  
  
**Summary**    
The `githubservice` package provides a lightweight wrapper around GitHub webhook events. It defines data structures (`GithubAPI`, `EventsService`, `IssueEvent`) and exposes two main entry points: `ConstructGithubApi` to build API credentials, and `StartWebhookHandler` to spin up an HTTP listener that forwards all issue events into a channel for further processing. The code relies on the external libraries `githubevents` (for webhook handling) and `go-github` (GitHub client), while configuration is driven by environment variables.  
  
