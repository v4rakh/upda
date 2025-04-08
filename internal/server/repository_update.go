package server

import (
	"encoding/json"
	"git.myservermanager.com/varakh/upda/api"
	"gorm.io/gorm"
	"time"
)

type updateRepository interface {
	paginate(page int, pageSize int, orderBy string, order string, searchTerm string, searchIn string, state ...api.UpdateState) ([]*Update, error)
	count(searchTerm string, searchIn string, state ...api.UpdateState) (int64, error)
	findAll() ([]*Update, error)
	find(id string) (*Update, error)
	findBy(application string, provider string, host string) (*Update, error)
	create(application string, provider string, host string, version string, metadata interface{}) (*Update, error)
	update(id string, version string, metadata interface{}) (*Update, error)
	updateState(id string, state api.UpdateState) (*Update, error)
	delete(id string) (int64, error)
	deleteByUpdatedAtBeforeAndStates(time time.Time, state ...api.UpdateState) (int64, error)
}

type updateDbRepo struct {
	db *gorm.DB
}

func newUpdateDbRepo(db *gorm.DB) *updateDbRepo {
	return &updateDbRepo{
		db: db,
	}
}

func (r *updateDbRepo) findAll() ([]*Update, error) {
	var e []*Update
	var res *gorm.DB

	if res = r.db.Find(&e); res.Error != nil {
		return nil, newServiceDatabaseError(res.Error)
	}

	return e, nil
}

func (r *updateDbRepo) find(id string) (*Update, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	var e Update
	var res *gorm.DB

	if res = r.db.Find(&e, "id = ?", id); res.Error != nil {
		return nil, newServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, errorResourceNotFound
	}

	return &e, nil
}

func (r *updateDbRepo) findBy(application string, provider string, host string) (*Update, error) {
	if application == "" || provider == "" || host == "" {
		return nil, errorValidationNotBlank
	}

	var e *Update
	var res *gorm.DB

	if res = r.db.Find(&e, &Update{Application: application, Provider: provider, Host: host}).Limit(1); res.Error != nil {
		return nil, newServiceDatabaseError(res.Error)
	}

	if res.RowsAffected == 0 || e == nil {
		return nil, errorResourceNotFound

	}

	return e, nil
}

func (r *updateDbRepo) create(application string, provider string, host string, version string, metadata interface{}) (*Update, error) {
	if application == "" || provider == "" || host == "" || version == "" {
		return nil, errorValidationNotBlank
	}

	e := &Update{
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
		return nil, newServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, errorDatabaseRowsExpected
	}

	return e, nil
}

func (r *updateDbRepo) updateState(id string, state api.UpdateState) (*Update, error) {
	if id == "" || state == "" {
		return nil, errorValidationNotBlank
	}

	var err error
	var e *Update

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

func (r *updateDbRepo) update(id string, version string, metadata interface{}) (*Update, error) {
	if id == "" || version == "" {
		return nil, errorValidationNotBlank
	}

	var err error
	var e *Update

	if e, err = r.find(id); err != nil {
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
		return e, errorDatabaseRowsExpected
	}

	return e, nil
}

func (r *updateDbRepo) delete(id string) (int64, error) {
	if id == "" {
		return 0, errorValidationNotBlank
	}

	var res *gorm.DB
	if res = r.db.Delete(&Update{}, "id = ?", id); res.Error != nil {
		return 0, newServiceDatabaseError(res.Error)
	}
	return res.RowsAffected, nil
}

func (r *updateDbRepo) deleteByUpdatedAtBeforeAndStates(time time.Time, state ...api.UpdateState) (int64, error) {
	if len(state) == 0 {
		return 0, errorValidationNotEmpty
	}

	states := make([]string, 0, len(state))
	for _, i := range state {
		states = append(states, i.Value())
	}

	var res *gorm.DB
	if res = r.db.Where("state IN ?", states).Where("updated_at < ?", time).Delete(&Update{}); res.Error != nil {
		return 0, newServiceDatabaseError(res.Error)
	}

	return res.RowsAffected, nil
}

func (r *updateDbRepo) paginate(page int, pageSize int, orderBy string, order string, searchTerm string, searchIn string, state ...api.UpdateState) ([]*Update, error) {
	if page == 0 {
		return nil, errorValidationPageGreaterZero
	}
	if pageSize <= 0 {
		return nil, errorValidationPageSizeGreaterZero
	}

	offset := (page - 1) * pageSize

	var e []*Update

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
		return nil, newServiceDatabaseError(res.Error)
	}

	return e, nil
}

func (r *updateDbRepo) count(searchTerm string, searchIn string, state ...api.UpdateState) (int64, error) {
	var c int64

	states := make([]string, 0, len(state))
	if len(state) > 0 {
		for _, s := range state {
			states = append(states, s.Value())
		}
	}

	if res := r.db.Model(&Update{}).Scopes(allGetUpdateCriterion(searchTerm, searchIn, states)).Count(&c); res.Error != nil {
		return 0, newServiceDatabaseError(res.Error)
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
