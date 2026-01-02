package app

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/mattjefferson/todoist-cli/internal/todoist"
)

type outputMode int

const (
	modeHuman outputMode = iota
	modeJSON
	modePlain
)

func printTasks(out io.Writer, tasks []todoist.Task, mode outputMode) error {
	switch mode {
	case modeJSON:
		payload := map[string]any{"results": tasks}
		return printJSON(out, payload)
	case modePlain:
		for _, task := range tasks {
			if _, err := fmt.Fprintf(out, "%s\t%s\t%s\n", task.ID, task.Content, dueSummary(task)); err != nil {
				return err
			}
		}
		return nil
	default:
		w := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
		if _, err := fmt.Fprintln(w, "ID\tCONTENT\tDUE"); err != nil {
			return err
		}
		for _, task := range tasks {
			if _, err := fmt.Fprintf(w, "%s\t%s\t%s\n", task.ID, task.Content, dueSummary(task)); err != nil {
				return err
			}
		}
		return w.Flush()
	}
}

func printProjects(out io.Writer, projects []todoist.Project, mode outputMode) error {
	switch mode {
	case modeJSON:
		payload := map[string]any{"results": projects}
		return printJSON(out, payload)
	case modePlain:
		for _, project := range projects {
			if _, err := fmt.Fprintf(out, "%s\t%s\n", project.ID, project.Name); err != nil {
				return err
			}
		}
		return nil
	default:
		w := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
		if _, err := fmt.Fprintln(w, "ID\tNAME"); err != nil {
			return err
		}
		for _, project := range projects {
			if _, err := fmt.Fprintf(w, "%s\t%s\n", project.ID, project.Name); err != nil {
				return err
			}
		}
		return w.Flush()
	}
}

func printLabels(out io.Writer, labels []todoist.Label, mode outputMode) error {
	switch mode {
	case modeJSON:
		payload := map[string]any{"results": labels}
		return printJSON(out, payload)
	case modePlain:
		for _, label := range labels {
			if _, err := fmt.Fprintf(out, "%s\t%s\t%s\t%d\t%t\n", label.ID, label.Name, label.Color, label.Order, label.IsFavorite); err != nil {
				return err
			}
		}
		return nil
	default:
		w := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
		if _, err := fmt.Fprintln(w, "ID\tNAME\tCOLOR\tORDER\tFAVORITE"); err != nil {
			return err
		}
		for _, label := range labels {
			if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%t\n", label.ID, label.Name, label.Color, label.Order, label.IsFavorite); err != nil {
				return err
			}
		}
		return w.Flush()
	}
}

func printComments(out io.Writer, comments []todoist.Comment, mode outputMode) error {
	switch mode {
	case modeJSON:
		payload := map[string]any{"results": comments}
		return printJSON(out, payload)
	case modePlain:
		for _, comment := range comments {
			if _, err := fmt.Fprintf(out, "%s\t%s\t%s\n", comment.ID, comment.Content, comment.PostedAt); err != nil {
				return err
			}
		}
		return nil
	default:
		w := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
		if _, err := fmt.Fprintln(w, "ID\tCONTENT\tPOSTED"); err != nil {
			return err
		}
		for _, comment := range comments {
			if _, err := fmt.Fprintf(w, "%s\t%s\t%s\n", comment.ID, comment.Content, comment.PostedAt); err != nil {
				return err
			}
		}
		return w.Flush()
	}
}

func printTask(out io.Writer, task todoist.Task, mode outputMode) error {
	switch mode {
	case modeJSON:
		return printJSON(out, task)
	case modePlain:
		_, err := fmt.Fprintf(out, "%s\t%s\t%s\n", task.ID, task.Content, dueSummary(task))
		return err
	default:
		_, err := fmt.Fprintf(out, "ID: %s\nContent: %s\nDue: %s\n", task.ID, task.Content, dueSummary(task))
		return err
	}
}

func printProject(out io.Writer, project todoist.Project, mode outputMode) error {
	switch mode {
	case modeJSON:
		return printJSON(out, project)
	case modePlain:
		_, err := fmt.Fprintf(out, "%s\t%s\n", project.ID, project.Name)
		return err
	default:
		_, err := fmt.Fprintf(out, "ID: %s\nName: %s\n", project.ID, project.Name)
		return err
	}
}

func printLabel(out io.Writer, label todoist.Label, mode outputMode) error {
	switch mode {
	case modeJSON:
		return printJSON(out, label)
	case modePlain:
		_, err := fmt.Fprintf(out, "%s\t%s\t%s\t%d\t%t\n", label.ID, label.Name, label.Color, label.Order, label.IsFavorite)
		return err
	default:
		_, err := fmt.Fprintf(out, "ID: %s\nName: %s\nColor: %s\nOrder: %d\nFavorite: %t\n", label.ID, label.Name, label.Color, label.Order, label.IsFavorite)
		return err
	}
}

func printComment(out io.Writer, comment todoist.Comment, mode outputMode) error {
	switch mode {
	case modeJSON:
		return printJSON(out, comment)
	case modePlain:
		_, err := fmt.Fprintf(out, "%s\t%s\t%s\n", comment.ID, comment.Content, comment.PostedAt)
		return err
	default:
		_, err := fmt.Fprintf(out, "ID: %s\nContent: %s\nPosted: %s\n", comment.ID, comment.Content, comment.PostedAt)
		return err
	}
}

func printJSON(out io.Writer, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(out, string(data))
	return err
}

func printRawJSON(out io.Writer, raw []byte) error {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" {
		_, err := fmt.Fprintln(out, "null")
		return err
	}
	_, err := fmt.Fprintln(out, trimmed)
	return err
}

func dueSummary(task todoist.Task) string {
	if task.Due == nil {
		return ""
	}
	if task.Due.Date != "" {
		return task.Due.Date
	}
	if task.Due.Datetime != "" {
		return task.Due.Datetime
	}
	return task.Due.String
}
