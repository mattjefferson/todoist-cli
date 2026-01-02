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

// User represents the currently authenticated Todoist user.
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

// Section represents a Todoist section.
type Section struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id,omitempty"`
	ProjectID    string `json:"project_id"`
	AddedAt      string `json:"added_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
	ArchivedAt   string `json:"archived_at,omitempty"`
	Name         string `json:"name"`
	SectionOrder int    `json:"section_order,omitempty"`
	IsArchived   bool   `json:"is_archived,omitempty"`
	IsDeleted    bool   `json:"is_deleted,omitempty"`
	IsCollapsed  bool   `json:"is_collapsed,omitempty"`
}

// Label represents a Todoist label.
type Label struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color,omitempty"`
	Order      int    `json:"order,omitempty"`
	IsFavorite bool   `json:"is_favorite,omitempty"`
}

// Upload represents a Todoist upload response.
type Upload struct {
	FileURL      string `json:"file_url"`
	FileName     string `json:"file_name"`
	FileSize     int    `json:"file_size"`
	FileType     string `json:"file_type"`
	ResourceType string `json:"resource_type,omitempty"`
	Image        string `json:"image,omitempty"`
	ImageWidth   int    `json:"image_width,omitempty"`
	ImageHeight  int    `json:"image_height,omitempty"`
	UploadState  string `json:"upload_state,omitempty"`
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

// Activity represents an activity log entry.
type Activity struct {
	ID              string         `json:"id"`
	EventType       string         `json:"event_type"`
	ObjectType      string         `json:"object_type"`
	ObjectID        string         `json:"object_id"`
	ParentProjectID string         `json:"parent_project_id,omitempty"`
	ParentItemID    string         `json:"parent_item_id,omitempty"`
	InitiatorID     string         `json:"initiator_id,omitempty"`
	EventDate       string         `json:"event_date,omitempty"`
	ExtraData       map[string]any `json:"extra_data,omitempty"`
}

// FileAttachment represents a comment file attachment.
type FileAttachment struct {
	FileName     string `json:"file_name,omitempty"`
	FileType     string `json:"file_type,omitempty"`
	FileURL      string `json:"file_url,omitempty"`
	FileSize     int    `json:"file_size,omitempty"`
	UploadState  string `json:"upload_state,omitempty"`
	ResourceType string `json:"resource_type,omitempty"`
	Image        string `json:"image,omitempty"`
	ImageWidth   int    `json:"image_width,omitempty"`
	ImageHeight  int    `json:"image_height,omitempty"`
	TnS          []any  `json:"tn_s,omitempty"`
	TnM          []any  `json:"tn_m,omitempty"`
	TnL          []any  `json:"tn_l,omitempty"`
}
