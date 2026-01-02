package todoist

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Client wraps Todoist API calls.
type Client struct {
	BaseURL string
	Token   string
	HTTP    *http.Client
	Verbose bool
}

// NewClient creates a Todoist API client.
func NewClient(baseURL, token string, verbose bool) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		Token:   token,
		HTTP: &http.Client{
			Timeout: 30 * time.Second,
		},
		Verbose: verbose,
	}
}

type listResponse[T any] struct {
	Results    []T    `json:"results"`
	NextCursor string `json:"next_cursor"`
}

// ListTasks fetches a page of tasks.
func (c *Client) ListTasks(ctx context.Context, params map[string]string) ([]Task, string, error) {
	var resp listResponse[Task]
	if err := c.get(ctx, "/api/v1/tasks", params, &resp); err != nil {
		return nil, "", err
	}
	return resp.Results, resp.NextCursor, nil
}

// ListTasksAll fetches all tasks across pages.
func (c *Client) ListTasksAll(ctx context.Context, params map[string]string) ([]Task, error) {
	if params == nil {
		params = map[string]string{}
	}
	params["limit"] = strconv.Itoa(200)
	var all []Task
	cursor := ""
	for {
		if cursor != "" {
			params["cursor"] = cursor
		}
		page, next, err := c.ListTasks(ctx, params)
		if err != nil {
			return nil, err
		}
		all = append(all, page...)
		if next == "" {
			break
		}
		cursor = next
	}
	return all, nil
}

// CreateTask creates a new task.
func (c *Client) CreateTask(ctx context.Context, body map[string]any) (Task, []byte, error) {
	var task Task
	raw, err := c.post(ctx, "/api/v1/tasks", body, &task)
	return task, raw, err
}

// UpdateTask updates an existing task.
func (c *Client) UpdateTask(ctx context.Context, id string, body map[string]any) (Task, []byte, error) {
	var task Task
	raw, err := c.post(ctx, "/api/v1/tasks/"+url.PathEscape(id), body, &task)
	return task, raw, err
}

// GetTask fetches a task by ID.
func (c *Client) GetTask(ctx context.Context, id string) (Task, error) {
	var task Task
	if err := c.get(ctx, "/api/v1/tasks/"+url.PathEscape(id), nil, &task); err != nil {
		return Task{}, err
	}
	return task, nil
}

// DeleteTask deletes a task by ID.
func (c *Client) DeleteTask(ctx context.Context, id string) ([]byte, error) {
	return c.delete(ctx, "/api/v1/tasks/"+url.PathEscape(id))
}

// CloseTask completes a task by ID.
func (c *Client) CloseTask(ctx context.Context, id string) ([]byte, error) {
	return c.post(ctx, "/api/v1/tasks/"+url.PathEscape(id)+"/close", nil, nil)
}

// ReopenTask reopens a completed task by ID.
func (c *Client) ReopenTask(ctx context.Context, id string) ([]byte, error) {
	return c.post(ctx, "/api/v1/tasks/"+url.PathEscape(id)+"/reopen", nil, nil)
}

// QuickAdd creates a task using Todoist quick-add syntax.
func (c *Client) QuickAdd(ctx context.Context, body map[string]any) (Task, []byte, error) {
	var resp struct {
		Task Task `json:"task"`
	}
	raw, err := c.post(ctx, "/api/v1/tasks/quick", body, &resp)
	return resp.Task, raw, err
}

// ListProjects fetches a page of projects.
func (c *Client) ListProjects(ctx context.Context, params map[string]string) ([]Project, string, error) {
	var resp listResponse[Project]
	if err := c.get(ctx, "/api/v1/projects", params, &resp); err != nil {
		return nil, "", err
	}
	return resp.Results, resp.NextCursor, nil
}

// ListProjectsAll fetches all projects across pages.
func (c *Client) ListProjectsAll(ctx context.Context) ([]Project, error) {
	params := map[string]string{"limit": "200"}
	var all []Project
	cursor := ""
	for {
		if cursor != "" {
			params["cursor"] = cursor
		}
		page, next, err := c.ListProjects(ctx, params)
		if err != nil {
			return nil, err
		}
		all = append(all, page...)
		if next == "" {
			break
		}
		cursor = next
	}
	return all, nil
}

// CreateProject creates a new project.
func (c *Client) CreateProject(ctx context.Context, body map[string]any) (Project, []byte, error) {
	var project Project
	raw, err := c.post(ctx, "/api/v1/projects", body, &project)
	return project, raw, err
}

// UpdateProject updates an existing project.
func (c *Client) UpdateProject(ctx context.Context, id string, body map[string]any) (Project, []byte, error) {
	var project Project
	raw, err := c.post(ctx, "/api/v1/projects/"+url.PathEscape(id), body, &project)
	return project, raw, err
}

// GetProject fetches a project by ID.
func (c *Client) GetProject(ctx context.Context, id string) (Project, error) {
	var project Project
	if err := c.get(ctx, "/api/v1/projects/"+url.PathEscape(id), nil, &project); err != nil {
		return Project{}, err
	}
	return project, nil
}

// DeleteProject deletes a project by ID.
func (c *Client) DeleteProject(ctx context.Context, id string) ([]byte, error) {
	return c.delete(ctx, "/api/v1/projects/"+url.PathEscape(id))
}

// FindProjectByName returns the project for a unique name match.
func (c *Client) FindProjectByName(ctx context.Context, name string) (Project, error) {
	projects, err := c.ListProjectsAll(ctx)
	if err != nil {
		return Project{}, err
	}
	matches := make([]Project, 0, 2)
	for _, project := range projects {
		if project.Name == name {
			matches = append(matches, project)
		}
	}
	if len(matches) == 0 {
		return Project{}, fmt.Errorf("project not found: %s", name)
	}
	if len(matches) > 1 {
		return Project{}, fmt.Errorf("project name not unique: %s", name)
	}
	return matches[0], nil
}

// FindProjectIDByName returns the project ID for a unique project name.
func (c *Client) FindProjectIDByName(ctx context.Context, name string) (string, error) {
	project, err := c.FindProjectByName(ctx, name)
	if err != nil {
		return "", err
	}
	return project.ID, nil
}

// FindTaskByContent returns the task for a unique content match.
func (c *Client) FindTaskByContent(ctx context.Context, title string) (Task, error) {
	tasks, err := c.ListTasksAll(ctx, nil)
	if err != nil {
		return Task{}, err
	}
	matches := make([]Task, 0, 2)
	for _, task := range tasks {
		if task.Content == title {
			matches = append(matches, task)
		}
	}
	if len(matches) == 0 {
		return Task{}, fmt.Errorf("task not found: %s", title)
	}
	if len(matches) > 1 {
		return Task{}, fmt.Errorf("task title not unique: %s", title)
	}
	return matches[0], nil
}

func (c *Client) get(ctx context.Context, path string, params map[string]string, out any) error {
	fullURL, err := c.url(path, params)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, out)
	return err
}

func (c *Client) post(ctx context.Context, path string, body map[string]any, out any) ([]byte, error) {
	fullURL, err := c.url(path, nil)
	if err != nil {
		return nil, err
	}
	var payload io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		payload = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, payload)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return c.do(req, out)
}

func (c *Client) delete(ctx context.Context, path string) ([]byte, error) {
	fullURL, err := c.url(path, nil)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fullURL, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req, nil)
}

func (c *Client) url(path string, params map[string]string) (string, error) {
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return "", err
	}
	u.Path = strings.TrimRight(u.Path, "/") + path
	if len(params) > 0 {
		q := u.Query()
		for key, value := range params {
			if value == "" {
				continue
			}
			q.Set(key, value)
		}
		u.RawQuery = q.Encode()
	}
	return u.String(), nil
}

func (c *Client) do(req *http.Request, out any) ([]byte, error) {
	req.Header.Set("Authorization", "Bearer "+c.Token)
	if c.Verbose {
		writef(os.Stderr, "%s %s\n", req.Method, req.URL.String())
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("api error: %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}
	if out != nil && len(body) > 0 {
		if err := json.Unmarshal(body, out); err != nil {
			return body, fmt.Errorf("decode response: %w", err)
		}
	}
	return body, nil
}
