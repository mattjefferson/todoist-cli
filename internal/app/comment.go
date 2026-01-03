package app

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/mattjefferson/todi/internal/todi"
)

type int64Slice []int64

func (s *int64Slice) String() string {
	values := make([]string, 0, len(*s))
	for _, v := range *s {
		values = append(values, strconv.FormatInt(v, 10))
	}
	return strings.Join(values, ",")
}

func (s *int64Slice) Set(value string) error {
	parsed, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err != nil {
		return err
	}
	*s = append(*s, parsed)
	return nil
}

func runComment(ctx context.Context, state *state, args []string) int {
	if len(args) == 0 {
		printCommentUsage(state.Out)
		return 2
	}
	switch args[0] {
	case "list":
		return runCommentList(ctx, state, args[1:])
	case "get":
		return runCommentGet(ctx, state, args[1:])
	case "add":
		return runCommentAdd(ctx, state, args[1:])
	case "update":
		return runCommentUpdate(ctx, state, args[1:])
	case "delete":
		return runCommentDelete(ctx, state, args[1:])
	case "-h", "--help", "help":
		printCommentUsage(state.Out)
		return 0
	default:
		writeLine(state.Err, "error: unknown comment command:", args[0])
		printCommentUsage(state.Err)
		return 2
	}
}

func runCommentList(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todi comment list", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var taskTitle string
	var taskID string
	var projectName string
	var projectID string
	var limit int
	var cursor string
	var all bool
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.StringVar(&taskTitle, "task", "", "Task title (exact match)")
	fs.StringVar(&taskID, "task-id", "", "Task ID")
	fs.StringVar(&projectName, "project", "", "Project title (exact match)")
	fs.StringVar(&projectID, "project-id", "", "Project ID")
	fs.IntVar(&limit, "limit", 50, "Max comments per page (1-200)")
	fs.StringVar(&cursor, "cursor", "", "Pagination cursor")
	fs.BoolVar(&all, "all", false, "Fetch all pages")
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printCommentUsage(state.Out)
		return 0
	}
	if len(fs.Args()) > 0 {
		writeLine(state.Err, "error: unexpected arguments")
		return 2
	}

	client, err := state.client()
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	key, value, err := resolveCommentScope(ctx, client, taskTitle, taskID, projectName, projectID)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}

	params := map[string]string{key: value}
	if limit > 0 {
		params["limit"] = strconv.Itoa(limit)
	}
	if cursor != "" {
		params["cursor"] = cursor
	}

	if all {
		comments, err := client.ListCommentsAll(ctx, params)
		if err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		if err := printComments(state.Out, comments, state.Mode); err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		return 0
	}

	comments, next, err := client.ListComments(ctx, params)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	if state.Mode == modeJSON {
		payload := map[string]any{"results": comments, "next_cursor": next}
		if err := printJSON(state.Out, payload); err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		return 0
	}
	if err := printComments(state.Out, comments, state.Mode); err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	return 0
}

func runCommentGet(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todi comment get", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printCommentUsage(state.Out)
		return 0
	}
	identifier := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if identifier == "" {
		writeLine(state.Err, "error: comment ID required")
		return 2
	}

	client, err := state.client()
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	comment, err := client.GetComment(ctx, identifier)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	if err := printComment(state.Out, comment, state.Mode); err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	return 0
}

func runCommentAdd(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todi comment add", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var taskTitle string
	var taskID string
	var projectName string
	var projectID string
	var notify int64Slice
	var uploadPath string
	var uploadName string
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.StringVar(&taskTitle, "task", "", "Task title (exact match)")
	fs.StringVar(&taskID, "task-id", "", "Task ID")
	fs.StringVar(&projectName, "project", "", "Project title (exact match)")
	fs.StringVar(&projectID, "project-id", "", "Project ID")
	fs.Var(&notify, "notify", "UID to notify (repeatable)")
	fs.StringVar(&uploadPath, "file", "", "Upload file attachment")
	fs.StringVar(&uploadName, "file-name", "", "Override upload file name")
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printCommentUsage(state.Out)
		return 0
	}
	content := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if content == "" {
		writeLine(state.Err, "error: content required")
		return 2
	}
	if uploadName != "" && uploadPath == "" {
		writeLine(state.Err, "error: --file-name requires --file")
		return 2
	}

	client, err := state.client()
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	key, value, err := resolveCommentScope(ctx, client, taskTitle, taskID, projectName, projectID)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}

	body := map[string]any{
		key:       value,
		"content": content,
	}
	if uploadPath != "" {
		uploadProjectID := ""
		if key == "project_id" {
			uploadProjectID = value
		}
		upload, _, err := client.UploadFile(ctx, uploadPath, uploadName, uploadProjectID)
		if err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		body["attachment"] = fileAttachmentFromUpload(upload)
	}
	if len(notify) > 0 {
		body["uids_to_notify"] = notify
	}

	comment, raw, err := client.CreateComment(ctx, body)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	if state.Mode == modeJSON {
		if err := printRawJSON(state.Out, raw); err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		return 0
	}
	if err := printComment(state.Out, comment, state.Mode); err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	return 0
}

func runCommentUpdate(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todi comment update", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var content string
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.StringVar(&content, "content", "", "Comment content")
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printCommentUsage(state.Out)
		return 0
	}
	identifier := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if identifier == "" {
		writeLine(state.Err, "error: comment ID required")
		return 2
	}
	if content == "" {
		writeLine(state.Err, "error: content required")
		return 2
	}

	client, err := state.client()
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	body := map[string]any{"content": content}
	comment, raw, err := client.UpdateComment(ctx, identifier, body)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	if state.Mode == modeJSON {
		if err := printRawJSON(state.Out, raw); err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		return 0
	}
	if comment.ID == "" {
		writeLine(state.Out, "ok")
		return 0
	}
	if err := printComment(state.Out, comment, state.Mode); err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	return 0
}

func runCommentDelete(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todi comment delete", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var force bool
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.BoolVar(&force, "force", false, "Skip confirmation")
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printCommentUsage(state.Out)
		return 0
	}
	identifier := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if identifier == "" {
		writeLine(state.Err, "error: comment ID required")
		return 2
	}

	if err := confirmDelete(state, "comment", identifier, force); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}

	client, err := state.client()
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	raw, err := client.DeleteComment(ctx, identifier)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	if state.Mode == modeJSON {
		if err := printRawJSON(state.Out, raw); err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		return 0
	}
	if _, err := fmt.Fprintf(state.Out, "deleted %s\n", identifier); err != nil {
		return 1
	}
	return 0
}

func resolveCommentScope(ctx context.Context, client *todi.Client, taskTitle, taskID, projectName, projectID string) (string, string, error) {
	if taskTitle != "" && taskID != "" {
		return "", "", fmt.Errorf("cannot use --task and --task-id together")
	}
	if projectName != "" && projectID != "" {
		return "", "", fmt.Errorf("cannot use --project and --project-id together")
	}
	hasTask := taskTitle != "" || taskID != ""
	hasProject := projectName != "" || projectID != ""
	if hasTask && hasProject {
		return "", "", fmt.Errorf("use either task or project, not both")
	}
	if taskID != "" {
		return "task_id", taskID, nil
	}
	if taskTitle != "" {
		task, err := client.FindTaskByContent(ctx, taskTitle)
		if err != nil {
			return "", "", err
		}
		if task.ID == "" {
			return "", "", fmt.Errorf("task not found: %s", taskTitle)
		}
		return "task_id", task.ID, nil
	}
	if projectID != "" {
		return "project_id", projectID, nil
	}
	if projectName != "" {
		projectIDValue, err := client.FindProjectIDByName(ctx, projectName)
		if err != nil {
			return "", "", err
		}
		return "project_id", projectIDValue, nil
	}
	return "", "", fmt.Errorf("task or project required")
}

func fileAttachmentFromUpload(upload todi.Upload) *todi.FileAttachment {
	return &todi.FileAttachment{
		FileURL:      upload.FileURL,
		FileName:     upload.FileName,
		FileType:     upload.FileType,
		FileSize:     upload.FileSize,
		ResourceType: upload.ResourceType,
		Image:        upload.Image,
		ImageWidth:   upload.ImageWidth,
		ImageHeight:  upload.ImageHeight,
		UploadState:  upload.UploadState,
	}
}
