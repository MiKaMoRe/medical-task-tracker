package task

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/MiKaMoRe/medical-task-tracker/internal/db/txrunner"
	apperrors "github.com/MiKaMoRe/medical-task-tracker/internal/domain/errors"
	"github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task *task.Task) (*task.Task, error)
	UpdateTask(ctx context.Context, task *task.Task) (*task.Task, error)
	DeleteTask(ctx context.Context, id task.ID) error
	AddTaskTags(ctx context.Context, id task.ID, tags []string) (*task.Task, error)
	RemoveTaskTags(ctx context.Context, id task.ID, tags []string) (*task.Task, error)
	GetTask(ctx context.Context, id task.ID) (*task.Task, error)
	GetTasksForPeriod(ctx context.Context, from time.Time, to time.Time) ([]*task.Task, error)
	MarkTaskDone(ctx context.Context, id task.ID) error
	MarkTaskOccurrenceDone(ctx context.Context, id task.ID, occurrenceDate time.Time) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) conn(ctx context.Context) *gorm.DB {
	return txrunner.DB(ctx, r.db)
}

func (r *taskRepository) CreateTask(ctx context.Context, t *task.Task) (*task.Task, error) {
	createDTO := toCreateTaskDTO(t)
	taskDTO := toTaskDTO(createDTO)
	db := r.conn(ctx)

	if err := db.Omit("RecurrenceRule", "Tags").Create(taskDTO).Error; err != nil {
		return nil, err
	}

	if taskDTO.RecurrenceRule != nil {
		taskDTO.RecurrenceRule.TaskID = taskDTO.ID
		if err := db.
			Omit("WeeklyRule", "MonthlyRule", "YearlyRule", "BiweeklyRule", "ShiftRule", "ParityRule").
			Create(taskDTO.RecurrenceRule).Error; err != nil {
			return nil, err
		}
		if err := createRecurrenceTypedRule(db, taskDTO.RecurrenceRule); err != nil {
			return nil, err
		}
	}

	if err := syncTaskTags(db, taskDTO.ID, taskDTO.Tags); err != nil {
		return nil, err
	}

	return toDomainTask(taskDTO), nil
}

func (r *taskRepository) UpdateTask(ctx context.Context, t *task.Task) (*task.Task, error) {
	db := r.conn(ctx)
	createDTO := toCreateTaskDTO(t)
	taskDTO := toTaskDTO(createDTO)

	result := db.Model(&TaskDTO{}).
		Where("id = ?", t.ID.Int()).
		Updates(map[string]any{
			"name":          taskDTO.Name,
			"description":   taskDTO.Description,
			"date":          taskDTO.Date,
			"is_recurrence": taskDTO.IsRecurrence,
		})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, apperrors.NotFound("task not found")
	}

	if err := db.Where("task_id = ?", t.ID.Int()).Delete(&RecurrenceRuleDTO{}).Error; err != nil {
		return nil, err
	}
	if taskDTO.RecurrenceRule != nil {
		taskDTO.RecurrenceRule.TaskID = t.ID.Int()
		if err := db.
			Omit("WeeklyRule", "MonthlyRule", "YearlyRule", "BiweeklyRule", "ShiftRule", "ParityRule").
			Create(taskDTO.RecurrenceRule).Error; err != nil {
			return nil, err
		}
		if err := createRecurrenceTypedRule(db, taskDTO.RecurrenceRule); err != nil {
			return nil, err
		}
	}

	if err := syncTaskTags(db, t.ID.Int(), taskDTO.Tags); err != nil {
		return nil, err
	}

	return r.GetTask(ctx, t.ID)
}

func (r *taskRepository) DeleteTask(ctx context.Context, id task.ID) error {
	result := r.conn(ctx).Delete(&TaskDTO{}, id.Int())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.NotFound("task not found")
	}
	return nil
}

func (r *taskRepository) AddTaskTags(ctx context.Context, id task.ID, tags []string) (*task.Task, error) {
	db := r.conn(ctx)
	if err := ensureTaskExists(db, id); err != nil {
		return nil, err
	}

	for _, tagName := range uniqueStrings(tags) {
		tagDTO := TagDTO{Name: tagName}
		if err := db.Where("name = ?", tagName).FirstOrCreate(&tagDTO).Error; err != nil {
			return nil, err
		}
		if err := db.Exec(
			"INSERT INTO tags_tasks (task_id, tag_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
			id.Int(), tagDTO.ID,
		).Error; err != nil {
			return nil, err
		}
	}

	return r.GetTask(ctx, id)
}

func (r *taskRepository) RemoveTaskTags(ctx context.Context, id task.ID, tags []string) (*task.Task, error) {
	db := r.conn(ctx)
	if err := ensureTaskExists(db, id); err != nil {
		return nil, err
	}

	normalizedTags := uniqueStrings(tags)
	if len(normalizedTags) == 0 {
		return r.GetTask(ctx, id)
	}

	var blocked []string
	err := db.
		Table("tags AS t").
		Select("t.name").
		Joins("JOIN tags_tasks tt ON tt.tag_id = t.id").
		Where("tt.task_id = ? AND t.is_required = TRUE AND t.name IN ?", id.Int(), normalizedTags).
		Find(&blocked).Error
	if err != nil {
		return nil, err
	}
	if len(blocked) > 0 {
		return nil, apperrors.Conflict(
			fmt.Sprintf("cannot remove required tags: %s", strings.Join(blocked, ", ")),
		)
	}

	if err := db.Exec(`
		DELETE FROM tags_tasks AS tt
		USING tags AS t
		WHERE tt.tag_id = t.id
		  AND tt.task_id = ?
		  AND t.name IN ?
	`, id.Int(), normalizedTags).Error; err != nil {
		return nil, err
	}

	return r.GetTask(ctx, id)
}

func (r *taskRepository) GetTask(ctx context.Context, id task.ID) (*task.Task, error) {
	var taskDTO TaskDTO
	err := r.conn(ctx).
		Preload("RecurrenceRule").
		Preload("RecurrenceRule.WeeklyRule").
		Preload("RecurrenceRule.MonthlyRule").
		Preload("RecurrenceRule.YearlyRule").
		Preload("RecurrenceRule.BiweeklyRule").
		Preload("RecurrenceRule.ShiftRule").
		Preload("RecurrenceRule.ParityRule").
		Preload("Completions").
		Preload("Tags").
		Where("id = ?", id.Int()).
		First(&taskDTO).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("task not found")
		}
		return nil, err
	}
	return toDomainTask(&taskDTO), nil
}

func (r *taskRepository) GetTasksForPeriod(ctx context.Context, from time.Time, to time.Time) ([]*task.Task, error) {
	var taskDTOs []TaskDTO
	err := r.conn(ctx).
		Preload("RecurrenceRule").
		Preload("RecurrenceRule.WeeklyRule").
		Preload("RecurrenceRule.MonthlyRule").
		Preload("RecurrenceRule.YearlyRule").
		Preload("RecurrenceRule.BiweeklyRule").
		Preload("RecurrenceRule.ShiftRule").
		Preload("RecurrenceRule.ParityRule").
		Preload("Completions", "occurrence_date >= ? AND occurrence_date <= ?", from.Format(time.DateOnly), to.Format(time.DateOnly)).
		Preload("Tags").
		Where(`
			(tasks.is_recurrence = false AND tasks.date >= ? AND tasks.date <= ?)
			OR
			(tasks.is_recurrence = true AND tasks.date <= ? AND (recurrence_rules.end_date IS NULL OR recurrence_rules.end_date >= ?))
		`, from, to, to, from).
		Joins("LEFT JOIN recurrence_rules ON recurrence_rules.task_id = tasks.id").
		Find(&taskDTOs).Error
	if err != nil {
		return nil, err
	}

	tasks := make([]*task.Task, 0, len(taskDTOs))
	for i := range taskDTOs {
		tasks = append(tasks, toDomainTask(&taskDTOs[i]))
	}
	return tasks, nil
}

func (r *taskRepository) MarkTaskDone(ctx context.Context, id task.ID) error {
	result := r.conn(ctx).
		Model(&TaskDTO{}).
		Where("id = ? AND is_recurrence = false", id.Int()).
		Updates(map[string]any{
			"is_done":      true,
			"completed_at": time.Now().UTC(),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.NotFound("task not found")
	}
	return nil
}

func (r *taskRepository) MarkTaskOccurrenceDone(ctx context.Context, id task.ID, occurrenceDate time.Time) error {
	completion := TaskOccurrenceCompletionDTO{
		TaskID:         id.Int(),
		OccurrenceDate: occurrenceDate.UTC(),
		CompletedAt:    time.Now().UTC(),
	}

	return r.conn(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "task_id"}, {Name: "occurrence_date"}},
			DoUpdates: clause.AssignmentColumns([]string{"completed_at"}),
		}).
		Create(&completion).Error
}

type recurrenceTypedRuleCreator func(db *gorm.DB, recurrenceRule *RecurrenceRuleDTO) error

var recurrenceTypedRuleCreators = map[string]recurrenceTypedRuleCreator{
	string(task.RecurringTypeWeekly):   createWeeklyRule,
	string(task.RecurringTypeMonthly):  createMonthlyRule,
	string(task.RecurringTypeYearly):   createYearlyRule,
	string(task.RecurringTypeBiweekly): createBiweeklyRule,
	string(task.RecurringTypeShift):    createShiftRule,
	string(task.RecurringTypeParity):   createParityRule,
}

func createRecurrenceTypedRule(db *gorm.DB, recurrenceRule *RecurrenceRuleDTO) error {
	creator, ok := recurrenceTypedRuleCreators[recurrenceRule.RuleType]
	if !ok {
		return fmt.Errorf("unsupported recurrence rule type: %s", recurrenceRule.RuleType)
	}
	return creator(db, recurrenceRule)
}

func createWeeklyRule(db *gorm.DB, recurrenceRule *RecurrenceRuleDTO) error {
	if recurrenceRule.WeeklyRule == nil {
		return errors.New("weekly recurrence rule payload is required")
	}
	recurrenceRule.WeeklyRule.RecurrenceRuleID = recurrenceRule.ID
	return db.Create(recurrenceRule.WeeklyRule).Error
}

func createMonthlyRule(db *gorm.DB, recurrenceRule *RecurrenceRuleDTO) error {
	if recurrenceRule.MonthlyRule == nil {
		return errors.New("monthly recurrence rule payload is required")
	}
	recurrenceRule.MonthlyRule.RecurrenceRuleID = recurrenceRule.ID
	return db.Create(recurrenceRule.MonthlyRule).Error
}

func createYearlyRule(db *gorm.DB, recurrenceRule *RecurrenceRuleDTO) error {
	if recurrenceRule.YearlyRule == nil {
		return errors.New("yearly recurrence rule payload is required")
	}
	recurrenceRule.YearlyRule.RecurrenceRuleID = recurrenceRule.ID
	return db.Create(recurrenceRule.YearlyRule).Error
}

func createBiweeklyRule(db *gorm.DB, recurrenceRule *RecurrenceRuleDTO) error {
	if recurrenceRule.BiweeklyRule == nil {
		return errors.New("biweekly recurrence rule payload is required")
	}
	recurrenceRule.BiweeklyRule.RecurrenceRuleID = recurrenceRule.ID
	return db.Create(recurrenceRule.BiweeklyRule).Error
}

func createShiftRule(db *gorm.DB, recurrenceRule *RecurrenceRuleDTO) error {
	if recurrenceRule.ShiftRule == nil {
		return errors.New("shift recurrence rule payload is required")
	}
	recurrenceRule.ShiftRule.RecurrenceRuleID = recurrenceRule.ID
	return db.Create(recurrenceRule.ShiftRule).Error
}

func createParityRule(db *gorm.DB, recurrenceRule *RecurrenceRuleDTO) error {
	if recurrenceRule.ParityRule == nil {
		return errors.New("parity recurrence rule payload is required")
	}
	recurrenceRule.ParityRule.RecurrenceRuleID = recurrenceRule.ID
	return db.Create(recurrenceRule.ParityRule).Error
}

func syncTaskTags(db *gorm.DB, taskID int, tags []TagDTO) error {
	if err := ensureRequiredTagsNotRemoved(db, taskID, tagDTONames(tags)); err != nil {
		return err
	}

	if err := db.Exec("DELETE FROM tags_tasks WHERE task_id = ?", taskID).Error; err != nil {
		return err
	}

	for i := range tags {
		tagDTO := &tags[i]
		if err := db.Where("name = ?", tagDTO.Name).FirstOrCreate(tagDTO).Error; err != nil {
			return err
		}
		if err := db.Exec(
			"INSERT INTO tags_tasks (task_id, tag_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
			taskID, tagDTO.ID,
		).Error; err != nil {
			return err
		}
	}

	return nil
}

func ensureTaskExists(db *gorm.DB, id task.ID) error {
	var taskDTO TaskDTO
	result := db.Model(&TaskDTO{}).Select("id").Where("id = ?", id.Int()).Limit(1).Find(&taskDTO)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.NotFound("task not found")
	}
	return nil
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func tagDTONames(tags []TagDTO) []string {
	names := make([]string, len(tags))
	for i := range tags {
		names[i] = tags[i].Name
	}
	return names
}

func ensureRequiredTagsNotRemoved(db *gorm.DB, taskID int, nextTagNames []string) error {
	var currentRequired []string
	if err := db.
		Table("tags AS t").
		Select("t.name").
		Joins("JOIN tags_tasks tt ON tt.tag_id = t.id").
		Where("tt.task_id = ? AND t.is_required = TRUE", taskID).
		Find(&currentRequired).Error; err != nil {
		return err
	}
	if len(currentRequired) == 0 {
		return nil
	}

	nextSet := make(map[string]struct{}, len(nextTagNames))
	for _, name := range nextTagNames {
		nextSet[name] = struct{}{}
	}

	missing := make([]string, 0)
	for _, requiredTag := range currentRequired {
		if _, ok := nextSet[requiredTag]; !ok {
			missing = append(missing, requiredTag)
		}
	}
	if len(missing) > 0 {
		return apperrors.Conflict(
			fmt.Sprintf("cannot remove required tags: %s", strings.Join(missing, ", ")),
		)
	}

	return nil
}
