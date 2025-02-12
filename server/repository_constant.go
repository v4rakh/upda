package server

import (
	"gorm.io/gorm"
)

type constantRepository interface {
	findAll() ([]*Constant, error)
	findById(id string) (*Constant, error)
	findByKey(key string) (*Constant, error)
	create(key string, value string) (*Constant, error)
	update(id string, value string) (*Constant, error)
	delete(id string) (int64, error)
}

type constantDbRepo struct {
	db *gorm.DB
}

func newConstantDbRepo(db *gorm.DB) *constantDbRepo {
	return &constantDbRepo{
		db: db,
	}
}

func (r *constantDbRepo) findAll() ([]*Constant, error) {
	var e []*Constant
	var res *gorm.DB

	if res = r.db.Order("key asc").Find(&e); res.Error != nil {
		return nil, newServiceDatabaseError(res.Error)
	}

	return e, nil
}

func (r *constantDbRepo) findById(id string) (*Constant, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	var e Constant
	var res *gorm.DB

	if res = r.db.Find(&e, "id = ?", id); res.Error != nil {
		return nil, newServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, errorResourceNotFound
	}

	return &e, nil
}

func (r *constantDbRepo) findByKey(key string) (*Constant, error) {
	if key == "" {
		return nil, errorValidationNotBlank
	}

	var e Constant
	var res *gorm.DB

	if res = r.db.Find(&e, "key = ?", key); res.Error != nil {
		return nil, newServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, errorResourceNotFound
	}

	return &e, nil
}

func (r *constantDbRepo) create(key string, value string) (*Constant, error) {
	if key == "" || value == "" {
		return nil, errorValidationNotBlank
	}

	var e *Constant

	e = &Constant{
		Key:   key,
		Value: value,
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

func (r *constantDbRepo) update(id string, value string) (*Constant, error) {
	if id == "" || value == "" {
		return nil, errorValidationNotBlank
	}

	var err error
	var e *Constant

	if e, err = r.findById(id); err != nil {
		return nil, err
	}

	e.Value = value

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, errorDatabaseRowsExpected
	}

	return e, nil
}

func (r *constantDbRepo) delete(id string) (int64, error) {
	if id == "" {
		return 0, errorValidationNotBlank
	}

	var res *gorm.DB
	if res = r.db.Delete(&Constant{}, "id = ?", id); res.Error != nil {
		return 0, newServiceDatabaseError(res.Error)
	}
	return res.RowsAffected, nil
}
