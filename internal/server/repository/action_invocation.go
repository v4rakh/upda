package repository

import (
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"gorm.io/gorm"
	"time"
)

type ActionInvocationRepository interface {
	Paginate(page int, pageSize int, orderBy string, order string) ([]*model.ActionInvocation, error)
	Count() (int64, error)
	Find(id string) (*model.ActionInvocation, error)
	FindAllByState(limit int, maxRetries int, state ...string) ([]*model.ActionInvocation, error)
	Create(eventId string, actionId string, state string) (*model.ActionInvocation, error)
	UpdateState(id string, state string) (*model.ActionInvocation, error)
	UpdateMessage(id string, message *string) (*model.ActionInvocation, error)
	UpdateRetryCount(id string, retryCount int) (*model.ActionInvocation, error)
	Delete(id string) (int64, error)
	DeleteByUpdatedAtBeforeAndStates(time time.Time, retryCount int, state ...string) (int64, error)
}

type ActionInvocationDbRepo struct {
	db *gorm.DB
}

func NewActionInvocationDbRepo(db *gorm.DB) *ActionInvocationDbRepo {
	return &ActionInvocationDbRepo{
		db: db,
	}
}

func (r *ActionInvocationDbRepo) Find(id string) (*model.ActionInvocation, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e model.ActionInvocation
	var res *gorm.DB
	if res = r.db.Find(&e, "id = ?", id); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, service_error.ErrResourceNotFound
	}

	return &e, nil
}

func (r *ActionInvocationDbRepo) Create(eventId string, actionId string, state string) (*model.ActionInvocation, error) {
	if eventId == "" || actionId == "" || state == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	e := &model.ActionInvocation{
		EventID:  eventId,
		ActionID: actionId,
		State:    state,
	}

	var res *gorm.DB
	if res = r.db.Create(&e); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, service_error.ErrDatabaseRowsExpected
	}

	return e, nil
}

func (r *ActionInvocationDbRepo) UpdateRetryCount(id string, retryCount int) (*model.ActionInvocation, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var err error
	var e *model.ActionInvocation

	if e, err = r.Find(id); err != nil {
		return nil, err
	}

	e.RetryCount = retryCount

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, service_error.ErrDatabaseRowsExpected
	}

	return e, nil
}

func (r *ActionInvocationDbRepo) UpdateState(id string, state string) (*model.ActionInvocation, error) {
	if id == "" || state == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var err error
	var e *model.ActionInvocation

	if e, err = r.Find(id); err != nil {
		return nil, err
	}

	e.State = state

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, service_error.ErrDatabaseRowsExpected
	}

	return e, nil
}

func (r *ActionInvocationDbRepo) UpdateMessage(id string, message *string) (*model.ActionInvocation, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var err error
	var e *model.ActionInvocation

	if e, err = r.Find(id); err != nil {
		return nil, err
	}

	e.Message = message

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, service_error.ErrDatabaseRowsExpected
	}

	return e, nil
}

func (r *ActionInvocationDbRepo) Delete(id string) (int64, error) {
	if id == "" {
		return 0, service_error.ErrValidationNotBlank
	}

	var res *gorm.DB
	if res = r.db.Delete(&model.ActionInvocation{}, "id = ?", id); res.Error != nil {
		return 0, service_error.NewServiceDatabaseError(res.Error)
	}

	return res.RowsAffected, nil
}

func (r *ActionInvocationDbRepo) Paginate(page int, pageSize int, orderBy string, order string) ([]*model.ActionInvocation, error) {
	if page == 0 {
		return nil, service_error.ErrValidationPageGreaterZero
	}
	if pageSize <= 0 {
		return nil, service_error.ErrValidationPageSizeGreaterZero
	}

	offset := (page - 1) * pageSize

	var e []*model.ActionInvocation
	var res *gorm.DB

	if orderBy != "" && order != "" {
		res = r.db.Order(orderBy + " " + order).Offset(offset).Limit(pageSize).Find(&e)
	} else {
		res = r.db.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&e)
	}

	if res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}

	return e, nil
}

func (r *ActionInvocationDbRepo) Count() (int64, error) {
	var c int64
	var res *gorm.DB

	if res = r.db.Model(&model.ActionInvocation{}).Count(&c); res.Error != nil {
		return 0, service_error.NewServiceDatabaseError(res.Error)
	}

	return c, nil
}

func (r *ActionInvocationDbRepo) FindAllByState(limit int, maxRetries int, state ...string) ([]*model.ActionInvocation, error) {
	if limit <= 0 {
		return nil, service_error.ErrValidationLimitGreaterZero
	}

	var e []*model.ActionInvocation

	if res := r.db.Model(&model.ActionInvocation{}).Scopes(allGetActionInvocationCriterion(state, maxRetries)).Order("created_at asc").Limit(limit).Find(&e); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}

	return e, nil
}

func (r *ActionInvocationDbRepo) DeleteByUpdatedAtBeforeAndStates(time time.Time, maxRetries int, state ...string) (int64, error) {
	if len(state) == 0 {
		return 0, service_error.ErrValidationNotEmpty
	}

	var res *gorm.DB
	if res = r.db.Where("retry_count >= ?", maxRetries).Where("state IN ?", state).Where("updated_at < ?", time).Delete(&model.ActionInvocation{}); res.Error != nil {
		return 0, service_error.NewServiceDatabaseError(res.Error)
	}

	return res.RowsAffected, nil
}

func criterionActonInvocationMaxRetries(maxRetries int) func(db *gorm.DB) *gorm.DB {
	if maxRetries > 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where("retry_count < ? ", maxRetries)
		}
	}

	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}

func criterionActionInvocationState(states []string) func(db *gorm.DB) *gorm.DB {
	if states != nil && len(states) > 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where("state IN (?)", states)
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}

func allGetActionInvocationCriterion(states []string, maxRetries int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(criterionActionInvocationState(states), criterionActonInvocationMaxRetries(maxRetries))
	}
}
