package todi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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

// ListActivities fetches a page of activity log entries.
func (c *Client) ListActivities(ctx context.Context, params map[string]string) ([]Activity, string, error) {
	var resp listResponse[Activity]
	if err := c.get(ctx, "/api/v1/activities", params, &resp); err != nil {
		return nil, "", err
	}
	return resp.Results, resp.NextCursor, nil
}

// ListActivitiesAll fetches all activity log entries across pages.
func (c *Client) ListActivitiesAll(ctx context.Context, params map[string]string) ([]Activity, error) {
	if params == nil {
		params = map[string]string{}
	}
	params["limit"] = strconv.Itoa(100)
	var all []Activity
	cursor := ""
	for {
		if cursor != "" {
			params["cursor"] = cursor
		}
		page, next, err := c.ListActivities(ctx, params)
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

// GetUserInfo fetches the authenticated user.
func (c *Client) GetUserInfo(ctx context.Context) (User, error) {
	var user User
	if err := c.get(ctx, "/api/v1/user", nil, &user); err != nil {
		return User{}, err
	}
	return user, nil
}

// ListSections fetches a page of sections.
func (c *Client) ListSections(ctx context.Context, params map[string]string) ([]Section, string, error) {
	var resp listResponse[Section]
	if err := c.get(ctx, "/api/v1/sections", params, &resp); err != nil {
		return nil, "", err
	}
	return resp.Results, resp.NextCursor, nil
}

// ListSectionsAll fetches all sections across pages.
func (c *Client) ListSectionsAll(ctx context.Context, params map[string]string) ([]Section, error) {
	if params == nil {
		params = map[string]string{}
	}
	params["limit"] = strconv.Itoa(200)
	var all []Section
	cursor := ""
	for {
		if cursor != "" {
			params["cursor"] = cursor
		}
		page, next, err := c.ListSections(ctx, params)
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

// CreateSection creates a new section.
func (c *Client) CreateSection(ctx context.Context, body map[string]any) (Section, []byte, error) {
	var section Section
	raw, err := c.post(ctx, "/api/v1/sections", body, &section)
	return section, raw, err
}

// UpdateSection updates an existing section.
func (c *Client) UpdateSection(ctx context.Context, id string, body map[string]any) (Section, []byte, error) {
	var section Section
	raw, err := c.post(ctx, "/api/v1/sections/"+url.PathEscape(id), body, &section)
	return section, raw, err
}

// GetSection fetches a section by ID.
func (c *Client) GetSection(ctx context.Context, id string) (Section, error) {
	var section Section
	if err := c.get(ctx, "/api/v1/sections/"+url.PathEscape(id), nil, &section); err != nil {
		return Section{}, err
	}
	return section, nil
}

// DeleteSection deletes a section by ID.
func (c *Client) DeleteSection(ctx context.Context, id string) ([]byte, error) {
	return c.delete(ctx, "/api/v1/sections/"+url.PathEscape(id))
}

// ListComments fetches a page of comments.
func (c *Client) ListComments(ctx context.Context, params map[string]string) ([]Comment, string, error) {
	var resp listResponse[Comment]
	if err := c.get(ctx, "/api/v1/comments", params, &resp); err != nil {
		return nil, "", err
	}
	return resp.Results, resp.NextCursor, nil
}

// ListCommentsAll fetches all comments across pages.
func (c *Client) ListCommentsAll(ctx context.Context, params map[string]string) ([]Comment, error) {
	if params == nil {
		params = map[string]string{}
	}
	params["limit"] = strconv.Itoa(200)
	var all []Comment
	cursor := ""
	for {
		if cursor != "" {
			params["cursor"] = cursor
		}
		page, next, err := c.ListComments(ctx, params)
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

// CreateComment creates a new comment.
func (c *Client) CreateComment(ctx context.Context, body map[string]any) (Comment, []byte, error) {
	var comment Comment
	raw, err := c.post(ctx, "/api/v1/comments", body, &comment)
	return comment, raw, err
}

// UpdateComment updates an existing comment.
func (c *Client) UpdateComment(ctx context.Context, id string, body map[string]any) (Comment, []byte, error) {
	var comment Comment
	raw, err := c.post(ctx, "/api/v1/comments/"+url.PathEscape(id), body, &comment)
	return comment, raw, err
}

// GetComment fetches a comment by ID.
func (c *Client) GetComment(ctx context.Context, id string) (Comment, error) {
	var comment Comment
	if err := c.get(ctx, "/api/v1/comments/"+url.PathEscape(id), nil, &comment); err != nil {
		return Comment{}, err
	}
	return comment, nil
}

// DeleteComment deletes a comment by ID.
func (c *Client) DeleteComment(ctx context.Context, id string) ([]byte, error) {
	return c.delete(ctx, "/api/v1/comments/"+url.PathEscape(id))
}

// UploadFile uploads a file for use in comments.
func (c *Client) UploadFile(ctx context.Context, path, name, projectID string) (Upload, []byte, error) {
	var upload Upload
	file, err := os.Open(path)
	if err != nil {
		return Upload{}, nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			return
		}
	}()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if projectID != "" {
		if err := writer.WriteField("project_id", projectID); err != nil {
			return Upload{}, nil, err
		}
	}
	if name != "" {
		if err := writer.WriteField("file_name", name); err != nil {
			return Upload{}, nil, err
		}
	}

	fileName := name
	if fileName == "" {
		fileName = filepath.Base(path)
	}
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return Upload{}, nil, err
	}
	if _, err := io.Copy(part, file); err != nil {
		return Upload{}, nil, err
	}
	if err := writer.Close(); err != nil {
		return Upload{}, nil, err
	}

	fullURL, err := c.url("/api/v1/uploads", nil)
	if err != nil {
		return Upload{}, nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, body)
	if err != nil {
		return Upload{}, nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	raw, err := c.do(req, &upload)
	return upload, raw, err
}

// DeleteUpload deletes a file upload by file URL.
func (c *Client) DeleteUpload(ctx context.Context, fileURL string) ([]byte, error) {
	params := map[string]string{"file_url": fileURL}
	return c.deleteWithParams(ctx, "/api/v1/uploads", params)
}

// ListLabels fetches a page of labels.
func (c *Client) ListLabels(ctx context.Context, params map[string]string) ([]Label, string, error) {
	var resp listResponse[Label]
	if err := c.get(ctx, "/api/v1/labels", params, &resp); err != nil {
		return nil, "", err
	}
	return resp.Results, resp.NextCursor, nil
}

// ListLabelsAll fetches all labels across pages.
func (c *Client) ListLabelsAll(ctx context.Context, params map[string]string) ([]Label, error) {
	if params == nil {
		params = map[string]string{}
	}
	params["limit"] = strconv.Itoa(200)
	var all []Label
	cursor := ""
	for {
		if cursor != "" {
			params["cursor"] = cursor
		}
		page, next, err := c.ListLabels(ctx, params)
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

// CreateLabel creates a new label.
func (c *Client) CreateLabel(ctx context.Context, body map[string]any) (Label, []byte, error) {
	var label Label
	raw, err := c.post(ctx, "/api/v1/labels", body, &label)
	return label, raw, err
}

// UpdateLabel updates an existing label.
func (c *Client) UpdateLabel(ctx context.Context, id string, body map[string]any) (Label, []byte, error) {
	var label Label
	raw, err := c.post(ctx, "/api/v1/labels/"+url.PathEscape(id), body, &label)
	return label, raw, err
}

// GetLabel fetches a label by ID.
func (c *Client) GetLabel(ctx context.Context, id string) (Label, error) {
	var label Label
	if err := c.get(ctx, "/api/v1/labels/"+url.PathEscape(id), nil, &label); err != nil {
		return Label{}, err
	}
	return label, nil
}

// DeleteLabel deletes a label by ID.
func (c *Client) DeleteLabel(ctx context.Context, id string) ([]byte, error) {
	return c.delete(ctx, "/api/v1/labels/"+url.PathEscape(id))
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

// FindSectionByName returns the section for a unique name match.
func (c *Client) FindSectionByName(ctx context.Context, name, projectID string) (Section, error) {
	params := map[string]string{}
	if projectID != "" {
		params["project_id"] = projectID
	}
	sections, err := c.ListSectionsAll(ctx, params)
	if err != nil {
		return Section{}, err
	}
	matches := make([]Section, 0, 2)
	for _, section := range sections {
		if section.Name == name {
			matches = append(matches, section)
		}
	}
	if len(matches) == 0 {
		return Section{}, fmt.Errorf("section not found: %s", name)
	}
	if len(matches) > 1 {
		return Section{}, fmt.Errorf("section name not unique: %s", name)
	}
	return matches[0], nil
}

// FindSectionIDByName returns the section ID for a unique name match.
func (c *Client) FindSectionIDByName(ctx context.Context, name, projectID string) (string, error) {
	section, err := c.FindSectionByName(ctx, name, projectID)
	if err != nil {
		return "", err
	}
	return section.ID, nil
}

// FindLabelByName returns the label for a unique name match.
func (c *Client) FindLabelByName(ctx context.Context, name string) (Label, error) {
	labels, err := c.ListLabelsAll(ctx, nil)
	if err != nil {
		return Label{}, err
	}
	matches := make([]Label, 0, 2)
	for _, label := range labels {
		if label.Name == name {
			matches = append(matches, label)
		}
	}
	if len(matches) == 0 {
		return Label{}, fmt.Errorf("label not found: %s", name)
	}
	if len(matches) > 1 {
		return Label{}, fmt.Errorf("label name not unique: %s", name)
	}
	return matches[0], nil
}

// FindLabelIDByName returns the label ID for a unique name match.
func (c *Client) FindLabelIDByName(ctx context.Context, name string) (string, error) {
	label, err := c.FindLabelByName(ctx, name)
	if err != nil {
		return "", err
	}
	return label.ID, nil
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

func (c *Client) deleteWithParams(ctx context.Context, path string, params map[string]string) ([]byte, error) {
	fullURL, err := c.url(path, params)
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
