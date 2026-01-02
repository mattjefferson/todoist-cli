package todoist

// Task represents a Todoist task.
type Task struct {
	ID          string   `json:"id"`
	Content     string   `json:"content"`
	Description string   `json:"description,omitempty"`
	ProjectID   string   `json:"project_id"`
	Labels      []string `json:"labels"`
	Priority    int      `json:"priority"`
	Due         *Due     `json:"due"`
}

// Due represents a task due date or datetime.
type Due struct {
	Date     string `json:"date"`
	Datetime string `json:"datetime"`
	String   string `json:"string"`
	Timezone string `json:"timezone"`
}

// Project represents a Todoist project.
type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Label represents a Todoist label.
type Label struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color,omitempty"`
	Order      int    `json:"order,omitempty"`
	IsFavorite bool   `json:"is_favorite,omitempty"`
}

// Comment represents a Todoist comment.
type Comment struct {
	ID             string              `json:"id"`
	ProjectID      string              `json:"project_id,omitempty"`
	TaskID         string              `json:"task_id,omitempty"`
	PostedUID      int64               `json:"posted_uid,omitempty"`
	Content        string              `json:"content"`
	FileAttachment *FileAttachment     `json:"file_attachment,omitempty"`
	UIDsToNotify   []int64             `json:"uids_to_notify,omitempty"`
	IsDeleted      bool                `json:"is_deleted,omitempty"`
	PostedAt       string              `json:"posted_at,omitempty"`
	Reactions      map[string][]string `json:"reactions,omitempty"`
}

// FileAttachment represents a comment file attachment.
type FileAttachment struct {
	FileName     string `json:"file_name,omitempty"`
	FileType     string `json:"file_type,omitempty"`
	FileURL      string `json:"file_url,omitempty"`
	ResourceType string `json:"resource_type,omitempty"`
}
