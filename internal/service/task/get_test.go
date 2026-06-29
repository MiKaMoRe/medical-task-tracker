package task

import (
	"context"
	"testing"
	"time"

	domaintask "github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
)

func TestTaskServiceGetTasks_FiltersDoneAndExpired(t *testing.T) {
	now := time.Now().UTC()
	from := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -2)
	to := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, time.UTC).AddDate(0, 0, 2)

	expiredDate := from.Add(12 * time.Hour)
	doneDate := from.AddDate(0, 0, 1).Add(12 * time.Hour)
	plannedDate := to.AddDate(0, 0, -1).Add(-12 * time.Hour)

	repo := &taskRepoMock{
		getTasksForPeriodFn: func(ctx context.Context, qFrom time.Time, qTo time.Time) ([]*domaintask.Task, error) {
			return []*domaintask.Task{
				{ID: domaintask.IDFromInt(1), Date: domaintask.Date(expiredDate), IsRecurring: false, IsDone: false},
				{ID: domaintask.IDFromInt(2), Date: domaintask.Date(doneDate), IsRecurring: false, IsDone: true},
				{ID: domaintask.IDFromInt(3), Date: domaintask.Date(plannedDate), IsRecurring: false, IsDone: false},
			}, nil
		},
	}
	svc := &TaskService{repo: repo}

	got, err := svc.GetTasks(context.Background(), from, to, []domaintask.TaskStatus{
		domaintask.TaskStatusDone,
		domaintask.TaskStatusExpired,
	})
	if err != nil {
		t.Fatalf("GetTasks returned error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 tasks after status filtering, got %d", len(got))
	}

	statuses := map[domaintask.TaskStatus]bool{}
	for _, task := range got {
		statuses[task.Status] = true
	}
	if !statuses[domaintask.TaskStatusDone] || !statuses[domaintask.TaskStatusExpired] {
		t.Fatalf("expected statuses to contain done and expired, got %+v", statuses)
	}
	if statuses[domaintask.TaskStatusPlanned] {
		t.Fatalf("expected planned status to be filtered out")
	}
}

func TestTaskServiceGetTasks_RejectsZeroPeriodBounds(t *testing.T) {
	svc := &TaskService{repo: &taskRepoMock{}}

	_, err := svc.GetTasks(context.Background(), time.Time{}, time.Now().UTC(), nil)
	assertErrorContains(t, err, "from and to are required")

	_, err = svc.GetTasks(context.Background(), time.Now().UTC(), time.Time{}, nil)
	assertErrorContains(t, err, "from and to are required")
}

func TestTaskServiceGetTasks_RejectsFromAfterTo(t *testing.T) {
	svc := &TaskService{repo: &taskRepoMock{}}
	from := mustUTC(2026, time.June, 27, 0, 0, 0)
	to := time.Date(2026, time.June, 26, 23, 59, 59, 0, time.UTC)

	_, err := svc.GetTasks(context.Background(), from, to, nil)
	assertErrorContains(t, err, "from must be less than or equal to to")
}
