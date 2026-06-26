package task

import (
	"time"

	"github.com/MiKaMoRe/medical-task-tracker/internal/repository"
)

type TaskDTO struct {
	repository.Base
	Name         string     `gorm:"column:name;type:varchar(255);not null"`
	Description  string     `gorm:"column:description;type:text"`
	Date         time.Time  `gorm:"column:date;type:timestamptz;not null"`
	IsRecurrence bool       `gorm:"column:is_recurrence;not null;default:false"`
	IsDone       bool       `gorm:"column:is_done;not null;default:false"`
	CompletedAt  *time.Time `gorm:"column:completed_at;type:timestamptz"`

	RecurrenceRule *RecurrenceRuleDTO            `gorm:"foreignKey:TaskID"`
	Completions    []TaskOccurrenceCompletionDTO `gorm:"foreignKey:TaskID"`
	Tags           []TagDTO                      `gorm:"many2many:tags_tasks;joinForeignKey:task_id;joinReferences:tag_id"`
}

func (TaskDTO) TableName() string {
	return "tasks"
}

type RecurrenceRuleDTO struct {
	repository.Base
	TaskID   int        `gorm:"column:task_id;not null;uniqueIndex"`
	RuleType string     `gorm:"column:rule_type;type:recurrence_rule_type;not null"`
	EndDate  *time.Time `gorm:"column:end_date;type:date"`

	WeeklyRule   *RecurrenceWeeklyRuleDTO   `gorm:"foreignKey:RecurrenceRuleID"`
	MonthlyRule  *RecurrenceMonthlyRuleDTO  `gorm:"foreignKey:RecurrenceRuleID"`
	YearlyRule   *RecurrenceYearlyRuleDTO   `gorm:"foreignKey:RecurrenceRuleID"`
	BiweeklyRule *RecurrenceBiweeklyRuleDTO `gorm:"foreignKey:RecurrenceRuleID"`
	ShiftRule    *RecurrenceShiftRuleDTO    `gorm:"foreignKey:RecurrenceRuleID"`
	ParityRule   *RecurrenceParityRuleDTO   `gorm:"foreignKey:RecurrenceRuleID"`
}

func (RecurrenceRuleDTO) TableName() string {
	return "recurrence_rules"
}

type RecurrenceWeeklyRuleDTO struct {
	repository.Base
	RecurrenceRuleID int `gorm:"column:recurrence_rule_id;not null;uniqueIndex"`
	WeekDay          int `gorm:"column:week_day;not null"`
}

func (RecurrenceWeeklyRuleDTO) TableName() string {
	return "recurrence_weekly_rules"
}

type RecurrenceMonthlyRuleDTO struct {
	repository.Base
	RecurrenceRuleID int `gorm:"column:recurrence_rule_id;not null;uniqueIndex"`
	MonthDay         int `gorm:"column:month_day;not null"`
}

func (RecurrenceMonthlyRuleDTO) TableName() string {
	return "recurrence_monthly_rules"
}

type RecurrenceYearlyRuleDTO struct {
	repository.Base
	RecurrenceRuleID int `gorm:"column:recurrence_rule_id;not null;uniqueIndex"`
	Month            int `gorm:"column:month;not null"`
	Day              int `gorm:"column:day;not null"`
}

func (RecurrenceYearlyRuleDTO) TableName() string {
	return "recurrence_yearly_rules"
}

type RecurrenceBiweeklyRuleDTO struct {
	repository.Base
	RecurrenceRuleID int  `gorm:"column:recurrence_rule_id;not null;uniqueIndex"`
	IsOdd            bool `gorm:"column:is_odd;not null"`
	WeekDay          int  `gorm:"column:week_day;not null"`
}

func (RecurrenceBiweeklyRuleDTO) TableName() string {
	return "recurrence_biweekly_rules"
}

type RecurrenceShiftRuleDTO struct {
	repository.Base
	RecurrenceRuleID  int `gorm:"column:recurrence_rule_id;not null;uniqueIndex"`
	NumberOfTaskDays  int `gorm:"column:number_of_task_days;not null"`
	NumberOfShiftDays int `gorm:"column:number_of_shift_days;not null"`
}

func (RecurrenceShiftRuleDTO) TableName() string {
	return "recurrence_shift_rules"
}

type RecurrenceParityRuleDTO struct {
	repository.Base
	RecurrenceRuleID int  `gorm:"column:recurrence_rule_id;not null;uniqueIndex"`
	IsEven           bool `gorm:"column:is_even;not null"`
}

func (RecurrenceParityRuleDTO) TableName() string {
	return "recurrence_parity_rules"
}

type TagDTO struct {
	repository.Base
	Name       string `gorm:"column:name;type:varchar(255);not null;uniqueIndex"`
	IsRequired bool   `gorm:"column:is_required;not null;default:false"`
}

type TaskOccurrenceCompletionDTO struct {
	repository.Base
	TaskID         int       `gorm:"column:task_id;not null;primaryKey"`
	OccurrenceDate time.Time `gorm:"column:occurrence_date;type:date;not null;primaryKey"`
	CompletedAt    time.Time `gorm:"column:completed_at;type:timestamptz;not null"`
}

func (TaskOccurrenceCompletionDTO) TableName() string {
	return "task_occurrence_completions"
}

func (TagDTO) TableName() string {
	return "tags"
}

type CreateTaskDTO struct {
	Title         string
	Description   string
	Date          time.Time
	IsRecurring   bool
	Tags          []string
	RecurringRule *CreateRecurringRuleDTO
}

type CreateRecurringRuleDTO struct {
	RecurringType string
	EndDate       *time.Time
	WeeklyRule    *CreateRecurrenceWeeklyRuleDTO
	MonthlyRule   *CreateRecurrenceMonthlyRuleDTO
	YearlyRule    *CreateRecurrenceYearlyRuleDTO
	BiweeklyRule  *CreateRecurrenceBiweeklyRuleDTO
	ShiftRule     *CreateRecurrenceShiftRuleDTO
	ParityRule    *CreateRecurrenceParityRuleDTO
}

type CreateRecurrenceWeeklyRuleDTO struct {
	WeekDay int
}

type CreateRecurrenceMonthlyRuleDTO struct {
	MonthDay int
}

type CreateRecurrenceYearlyRuleDTO struct {
	Month int
	Day   int
}

type CreateRecurrenceBiweeklyRuleDTO struct {
	IsOdd   bool
	WeekDay int
}

type CreateRecurrenceShiftRuleDTO struct {
	NumberOfTaskDays  int
	NumberOfShiftDays int
}

type CreateRecurrenceParityRuleDTO struct {
	IsEven bool
}
