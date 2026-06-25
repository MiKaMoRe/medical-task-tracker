package task

type CreateTaskRequest struct {
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Date        string               `json:"date"`
	IsRecurring bool                 `json:"is_recurring"`
	Tags        []string             `json:"tags"`
	Recurring   *CreateRecurringTask `json:"recurring,omitempty"`
}

type CreateRecurringTask struct {
	RecurringType string `json:"recurring_type"`
	EndDate       string `json:"end_date"`
}

type CreateTaskResponse struct {
	ID string `json:"id"`
}

type GetTaskRequest struct {
	ID string `json:"id"`
}

type GetTaskResponse struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Date        string   `json:"date"`
	DueDate     string   `json:"due_date"`
	IsRecurring bool     `json:"is_recurring"`
	Tags        []string `json:"tags"`
}
