package server

import (
	"git.myservermanager.com/varakh/upda/api"
	"gorm.io/gorm"
	"time"
)

type ActionInvocationRepository interface {
	paginate(page int, pageSize int, orderBy string, order string) ([]*ActionInvocation, error)
	count() (int64, error)
	find(id string) (*ActionInvocation, error)
	findAllByState(limit int, maxRetries int, state ...api.ActionInvocationState) ([]*ActionInvocation, error)
	create(eventId string, actionId string, state api.ActionInvocationState) (*ActionInvocation, error)
	updateState(id string, state api.ActionInvocationState) (*ActionInvocation, error)
	updateMessage(id string, message *string) (*ActionInvocation, error)
	updateRetryCount(id string, retryCount int) (*ActionInvocation, error)
	delete(id string) (int64, error)
	deleteByUpdatedAtBeforeAndStates(time time.Time, retryCount int, state ...api.ActionInvocationState) (int64, error)
}

type actionInvocationDbRepo struct {
	db *gorm.DB
}

func newActionInvocationDbRepo(db *gorm.DB) *actionInvocationDbRepo {
	return &actionInvocationDbRepo{
		db: db,
	}
}

func (r *actionInvocationDbRepo) find(id string) (*ActionInvocation, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	var e ActionInvocation
	var res *gorm.DB
	if res = r.db.Find(&e, "id = ?", id); res.Error != nil {
		return nil, newServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, errorResourceNotFound
	}

	return &e, nil
}

func (r *actionInvocationDbRepo) create(eventId string, actionId string, state api.ActionInvocationState) (*ActionInvocation, error) {
	if eventId == "" || actionId == "" || state == "" {
		return nil, errorValidationNotBlank
	}

	e := &ActionInvocation{
		EventID:  eventId,
		ActionID: actionId,
		State:    state.Value(),
	}

	var res *gorm.DB
	if res = r.db.Create(&e); res.Error != nil {
		return nil, newServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, errorDatabaseRowsExpected
	}

	return e, nil
}

func (r *actionInvocationDbRepo) updateRetryCount(id string, retryCount int) (*ActionInvocation, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	var err error
	var e *ActionInvocation

	if e, err = r.find(id); err != nil {
		return nil, err
	}

	e.RetryCount = retryCount

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, errorDatabaseRowsExpected
	}

	return e, nil
}

func (r *actionInvocationDbRepo) updateState(id string, state api.ActionInvocationState) (*ActionInvocation, error) {
	if id == "" || state == "" {
		return nil, errorValidationNotBlank
	}

	var err error
	var e *ActionInvocation

	if e, err = r.find(id); err != nil {
		return nil, err
	}

	e.State = state.Value()

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, errorDatabaseRowsExpected
	}

	return e, nil
}

func (r *actionInvocationDbRepo) updateMessage(id string, message *string) (*ActionInvocation, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	var err error
	var e *ActionInvocation

	if e, err = r.find(id); err != nil {
		return nil, err
	}

	e.Message = message

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, errorDatabaseRowsExpected
	}

	return e, nil
}

func (r *actionInvocationDbRepo) delete(id string) (int64, error) {
	if id == "" {
		return 0, errorValidationNotBlank
	}

	var res *gorm.DB
	if res = r.db.Delete(&ActionInvocation{}, "id = ?", id); res.Error != nil {
		return 0, newServiceDatabaseError(res.Error)
	}

	return res.RowsAffected, nil
}

func (r *actionInvocationDbRepo) paginate(page int, pageSize int, orderBy string, order string) ([]*ActionInvocation, error) {
	if page == 0 {
		return nil, errorValidationPageGreaterZero
	}
	if pageSize <= 0 {
		return nil, errorValidationPageSizeGreaterZero
	}

	offset := (page - 1) * pageSize

	var e []*ActionInvocation
	var res *gorm.DB

	if orderBy != "" && order != "" {
		res = r.db.Order(orderBy + " " + order).Offset(offset).Limit(pageSize).Find(&e)
	} else {
		res = r.db.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&e)
	}

	if res.Error != nil {
		return nil, newServiceDatabaseError(res.Error)
	}

	return e, nil
}

func (r *actionInvocationDbRepo) count() (int64, error) {
	var c int64
	var res *gorm.DB

	if res = r.db.Model(&ActionInvocation{}).Count(&c); res.Error != nil {
		return 0, newServiceDatabaseError(res.Error)
	}

	return c, nil
}

func (r *actionInvocationDbRepo) findAllByState(limit int, maxRetries int, state ...api.ActionInvocationState) ([]*ActionInvocation, error) {
	if limit <= 0 {
		return nil, errorValidationLimitGreaterZero
	}

	var e []*ActionInvocation

	states := translateActionInvocationState(state...)

	if res := r.db.Model(&ActionInvocation{}).Scopes(allGetActionInvocationCriterion(states, maxRetries)).Order("created_at asc").Limit(limit).Find(&e); res.Error != nil {
		return nil, newServiceDatabaseError(res.Error)
	}

	return e, nil
}

func (r *actionInvocationDbRepo) deleteByUpdatedAtBeforeAndStates(time time.Time, maxRetries int, state ...api.ActionInvocationState) (int64, error) {
	if len(state) == 0 {
		return 0, errorValidationNotEmpty
	}

	states := translateActionInvocationState(state...)

	var res *gorm.DB
	if res = r.db.Where("retry_count >= ?", maxRetries).Where("state IN ?", states).Where("updated_at < ?", time).Delete(&ActionInvocation{}); res.Error != nil {
		return 0, newServiceDatabaseError(res.Error)
	}

	return res.RowsAffected, nil
}

func translateActionInvocationState(state ...api.ActionInvocationState) []string {
	states := make([]string, 0, len(state))
	if len(state) > 0 {
		for _, s := range state {
			states = append(states, s.Value())
		}
	}

	return states
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
