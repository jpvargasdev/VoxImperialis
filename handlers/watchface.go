package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

// NewWatchfaceHandler returns a HandlerFunc that forwards watchface commands
// to the MachinusCronus API.
func NewWatchfaceHandler(baseURL string) HandlerFunc {
	baseURL = strings.TrimRight(baseURL, "/")
	apiURL := baseURL + "/api/v1/watchfaces"

	return func(cmd Command) string {
		if len(cmd.Args) < 1 {
			return watchfaceUsage()
		}

		action := strings.ToLower(cmd.Args[0])

		switch action {
		case "create":
			return watchfaceCreate(apiURL, cmd.Args[1:])
		case "list":
			return watchfaceList(apiURL)
		case "status":
			return watchfaceStatus(apiURL, cmd.Args[1:])
		case "iterate":
			return watchfaceIterate(apiURL, cmd.Args[1:])
		case "stop":
			return watchfaceStop(apiURL, cmd.Args[1:])
		case "delete":
			return watchfaceDelete(apiURL, cmd.Args[1:])
		default:
			return fmt.Sprintf("[error]\nunknown action '%s'\n\n%s", action, watchfaceUsage())
		}
	}
}

// watchface create <name> <prompt...>
func watchfaceCreate(apiURL string, args []string) string {
	if len(args) < 2 {
		return "[error]\nusage: watchface create <name> <prompt...>"
	}

	name := args[0]
	prompt := strings.Join(args[1:], " ")

	body := map[string]string{
		"name":   name,
		"prompt": prompt,
	}

	resp, err := doPost(apiURL, body)
	if err != nil {
		return fmt.Sprintf("[error]\n%v", err)
	}

	return fmt.Sprintf("[watchface]\naction:  create\nname:    %s\nstatus:  %s\npath:    %s",
		getString(resp, "name"),
		getString(resp, "status"),
		getString(resp, "repo_path"),
	)
}

// watchface list
func watchfaceList(apiURL string) string {
	resp, err := doGet(apiURL)
	if err != nil {
		return fmt.Sprintf("[error]\n%v", err)
	}

	watchfaces, ok := resp["watchfaces"].([]interface{})
	if !ok || len(watchfaces) == 0 {
		return "[watchface]\nno projects found"
	}

	var sb strings.Builder
	sb.WriteString("[watchface]\nprojects:\n")
	for _, w := range watchfaces {
		wf, ok := w.(map[string]interface{})
		if !ok {
			continue
		}
		sb.WriteString(fmt.Sprintf("  - %s (%s) iteration:%v\n",
			getString(wf, "name"),
			getString(wf, "status"),
			wf["iteration"],
		))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// watchface status <name>
func watchfaceStatus(apiURL string, args []string) string {
	if len(args) < 1 {
		return "[error]\nusage: watchface status <name>"
	}

	name := args[0]
	resp, err := doGet(apiURL + "/" + name + "/status")
	if err != nil {
		return fmt.Sprintf("[error]\n%v", err)
	}

	return fmt.Sprintf("[watchface]\nname:       %s\nstatus:     %s\niteration:  %v\nlast_commit: %s\npath:       %s",
		getString(resp, "name"),
		getString(resp, "status"),
		resp["iteration"],
		getString(resp, "last_commit"),
		getString(resp, "repo_path"),
	)
}

// watchface iterate <name> <feedback...>
func watchfaceIterate(apiURL string, args []string) string {
	if len(args) < 2 {
		return "[error]\nusage: watchface iterate <name> <feedback...>"
	}

	name := args[0]
	feedback := strings.Join(args[1:], " ")

	body := map[string]string{
		"feedback": feedback,
	}

	resp, err := doPost(apiURL+"/"+name+"/iterate", body)
	if err != nil {
		return fmt.Sprintf("[error]\n%v", err)
	}

	return fmt.Sprintf("[watchface]\naction:     iterate\nname:       %s\nstatus:     %s\niteration:  %v",
		getString(resp, "name"),
		getString(resp, "status"),
		resp["iteration"],
	)
}

// watchface stop <name>
func watchfaceStop(apiURL string, args []string) string {
	if len(args) < 1 {
		return "[error]\nusage: watchface stop <name>"
	}

	name := args[0]
	resp, err := doPost(apiURL+"/"+name+"/stop", nil)
	if err != nil {
		return fmt.Sprintf("[error]\n%v", err)
	}

	return fmt.Sprintf("[watchface]\naction:  stop\nname:    %s\nstatus:  %s",
		getString(resp, "name"),
		getString(resp, "status"),
	)
}

// watchface delete <name>
func watchfaceDelete(apiURL string, args []string) string {
	if len(args) < 1 {
		return "[error]\nusage: watchface delete <name>"
	}

	name := args[0]
	resp, err := doDelete(apiURL + "/" + name)
	if err != nil {
		return fmt.Sprintf("[error]\n%v", err)
	}

	return fmt.Sprintf("[watchface]\n%s", getString(resp, "message"))
}

func watchfaceUsage() string {
	return `[watchface]
usage:
  watchface create  <name> <prompt...>   create new project
  watchface list                         list all projects
  watchface status  <name>               check agent progress
  watchface iterate <name> <feedback...> send feedback to agent
  watchface stop    <name>               stop agent work
  watchface delete  <name>               remove project`
}

// --- HTTP helpers ---

func doGet(url string) (map[string]interface{}, error) {
	log.Printf("watchface: GET %s", url)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()
	return parseResponse(resp)
}

func doPost(url string, body interface{}) (map[string]interface{}, error) {
	log.Printf("watchface: POST %s", url)

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to encode request: %v", err)
		}
		reqBody = bytes.NewReader(data)
	}

	resp, err := httpClient.Post(url, "application/json", reqBody)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()
	return parseResponse(resp)
}

func doDelete(url string) (map[string]interface{}, error) {
	log.Printf("watchface: DELETE %s", url)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()
	return parseResponse(resp)
}

func parseResponse(resp *http.Response) (map[string]interface{}, error) {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v\nbody: %s", err, string(data))
	}

	if resp.StatusCode >= 400 {
		if errMsg, ok := result["error"].(string); ok {
			return nil, fmt.Errorf("%s", errMsg)
		}
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(data))
	}

	return result, nil
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		return fmt.Sprintf("%v", v)
	}
	return ""
}
