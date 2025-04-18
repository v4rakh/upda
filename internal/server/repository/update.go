package repository

import (
	"encoding/json"
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"gorm.io/gorm"
	"time"
)

type UpdateRepository interface {
	Paginate(page int, pageSize int, orderBy string, order string, searchTerm string, searchIn string, state ...api.UpdateState) ([]*model.Update, error)
	Count(searchTerm string, searchIn string, state ...api.UpdateState) (int64, error)
	FindAll() ([]*model.Update, error)
	Find(id string) (*model.Update, error)
	FindBy(application string, provider string, host string) (*model.Update, error)
	Create(application string, provider string, host string, version string, metadata interface{}) (*model.Update, error)
	Update(id string, version string, metadata interface{}) (*model.Update, error)
	UpdateState(id string, state api.UpdateState) (*model.Update, error)
	Delete(id string) (int64, error)
	DeleteByUpdatedAtBeforeAndStates(time time.Time, state ...api.UpdateState) (int64, error)
}

type UpdateDbRepo struct {
	db *gorm.DB
}

func NewUpdateDbRepo(db *gorm.DB) *UpdateDbRepo {
	return &UpdateDbRepo{
		db: db,
	}
}

func (r *UpdateDbRepo) FindAll() ([]*model.Update, error) {
	var e []*model.Update
	var res *gorm.DB

	if res = r.db.Find(&e); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}

	return e, nil
}

func (r *UpdateDbRepo) Find(id string) (*model.Update, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e model.Update
	var res *gorm.DB

	if res = r.db.Find(&e, "id = ?", id); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, service_error.ErrResourceNotFound
	}

	return &e, nil
}

func (r *UpdateDbRepo) FindBy(application string, provider string, host string) (*model.Update, error) {
	if application == "" || provider == "" || host == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e *model.Update
	var res *gorm.DB

	if res = r.db.Find(&e, &model.Update{Application: application, Provider: provider, Host: host}).Limit(1); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}

	if res.RowsAffected == 0 || e == nil {
		return nil, service_error.ErrResourceNotFound

	}

	return e, nil
}

func (r *UpdateDbRepo) Create(application string, provider string, host string, version string, metadata interface{}) (*model.Update, error) {
	if application == "" || provider == "" || host == "" || version == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	e := &model.Update{
		Application: application,
		Provider:    provider,
		Host:        host,
		Version:     version,
		State:       api.UpdateStatePending.Value(),
	}

	if metadata != nil {
		b, err := json.Marshal(metadata)
		if err != nil {
			return nil, err
		}
		e.Metadata = b
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

func (r *UpdateDbRepo) UpdateState(id string, state api.UpdateState) (*model.Update, error) {
	if id == "" || state == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var err error
	var e *model.Update

	if e, err = r.Find(id); err != nil {
		return nil, err
	}

	e.State = state.Value()

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, service_error.ErrDatabaseRowsExpected
	}

	return e, nil
}

func (r *UpdateDbRepo) Update(id string, version string, metadata interface{}) (*model.Update, error) {
	if id == "" || version == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var err error
	var e *model.Update

	if e, err = r.Find(id); err != nil {
		return nil, err
	}

	if metadata != nil {
		var b []byte
		if b, err = json.Marshal(metadata); err != nil {
			return nil, err
		}
		e.Metadata = b
	}

	e.Version = version

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, service_error.ErrDatabaseRowsExpected
	}

	return e, nil
}

func (r *UpdateDbRepo) Delete(id string) (int64, error) {
	if id == "" {
		return 0, service_error.ErrValidationNotBlank
	}

	var res *gorm.DB
	if res = r.db.Delete(&model.Update{}, "id = ?", id); res.Error != nil {
		return 0, service_error.NewServiceDatabaseError(res.Error)
	}
	return res.RowsAffected, nil
}

func (r *UpdateDbRepo) DeleteByUpdatedAtBeforeAndStates(time time.Time, state ...api.UpdateState) (int64, error) {
	if len(state) == 0 {
		return 0, service_error.ErrValidationNotEmpty
	}

	states := make([]string, 0, len(state))
	for _, i := range state {
		states = append(states, i.Value())
	}

	var res *gorm.DB
	if res = r.db.Where("state IN ?", states).Where("updated_at < ?", time).Delete(&model.Update{}); res.Error != nil {
		return 0, service_error.NewServiceDatabaseError(res.Error)
	}

	return res.RowsAffected, nil
}

func (r *UpdateDbRepo) Paginate(page int, pageSize int, orderBy string, order string, searchTerm string, searchIn string, state ...api.UpdateState) ([]*model.Update, error) {
	if page == 0 {
		return nil, service_error.ErrValidationPageGreaterZero
	}
	if pageSize <= 0 {
		return nil, service_error.ErrValidationPageSizeGreaterZero
	}

	offset := (page - 1) * pageSize

	var e []*model.Update

	if orderBy == "" {
		orderBy = "updated_at"
	}
	if order == "" {
		order = "desc"
	}

	states := make([]string, 0, len(state))
	if len(state) > 0 {
		for _, s := range state {
			states = append(states, s.Value())
		}
	}

	if res := r.db.Scopes(allGetUpdateCriterion(searchTerm, searchIn, states)).Order(orderBy + " " + order).Offset(offset).Limit(pageSize).Find(&e); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}

	return e, nil
}

func (r *UpdateDbRepo) Count(searchTerm string, searchIn string, state ...api.UpdateState) (int64, error) {
	var c int64

	states := make([]string, 0, len(state))
	if len(state) > 0 {
		for _, s := range state {
			states = append(states, s.Value())
		}
	}

	if res := r.db.Model(&model.Update{}).Scopes(allGetUpdateCriterion(searchTerm, searchIn, states)).Count(&c); res.Error != nil {
		return 0, service_error.NewServiceDatabaseError(res.Error)
	}

	return c, nil
}

func criterionUpdateSearch(searchTerm string, searchIn string) func(db *gorm.DB) *gorm.DB {
	if searchTerm == "" || searchIn == "" {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}
	switch searchIn {
	case "host":
		return func(db *gorm.DB) *gorm.DB {
			return db.Where("host LIKE ?", "%"+searchTerm+"%")
		}
	case "provider":
		return func(db *gorm.DB) *gorm.DB {
			return db.Where("provider LIKE ?", "%"+searchTerm+"%")
		}
	default:
		return func(db *gorm.DB) *gorm.DB {
			return db.Where("application LIKE ?", "%"+searchTerm+"%")
		}
	}
}

func criterionUpdateState(states []string) func(db *gorm.DB) *gorm.DB {
	if states != nil && len(states) > 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where("state IN (?)", states)
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}

func allGetUpdateCriterion(searchTerm string, searchIn string, states []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(criterionUpdateSearch(searchTerm, searchIn), criterionUpdateState(states))
	}
}
