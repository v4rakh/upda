package server

import (
	"gorm.io/gorm"
)

type secretRepository interface {
	findAll() ([]*Secret, error)
	findById(id string) (*Secret, error)
	findByKey(key string) (*Secret, error)
	create(key string, value string) (*Secret, error)
	update(id string, value string) (*Secret, error)
	delete(id string) (int64, error)
}

type secretDbRepo struct {
	db *gorm.DB
}

func newSecretDbRepo(db *gorm.DB) *secretDbRepo {
	return &secretDbRepo{
		db: db,
	}
}

func (r *secretDbRepo) findAll() ([]*Secret, error) {
	var e []*Secret
	var res *gorm.DB

	if res = r.db.Order("key asc").Find(&e); res.Error != nil {
		return nil, newServiceDatabaseError(res.Error)
	}

	return e, nil
}

func (r *secretDbRepo) findById(id string) (*Secret, error) {
	if id == "" {
		return nil, errorValidationNotBlank
	}

	var e Secret
	var res *gorm.DB

	if res = r.db.Find(&e, "id = ?", id); res.Error != nil {
		return nil, newServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, errorResourceNotFound
	}

	return &e, nil
}

func (r *secretDbRepo) findByKey(key string) (*Secret, error) {
	if key == "" {
		return nil, errorValidationNotBlank
	}

	var e Secret
	var res *gorm.DB

	if res = r.db.Find(&e, "key = ?", key); res.Error != nil {
		return nil, newServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, errorResourceNotFound
	}

	return &e, nil
}

func (r *secretDbRepo) create(key string, value string) (*Secret, error) {
	if key == "" || value == "" {
		return nil, errorValidationNotBlank
	}

	var e *Secret

	e = &Secret{
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

func (r *secretDbRepo) update(id string, value string) (*Secret, error) {
	if id == "" || value == "" {
		return nil, errorValidationNotBlank
	}

	var err error
	var e *Secret

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

func (r *secretDbRepo) delete(id string) (int64, error) {
	if id == "" {
		return 0, errorValidationNotBlank
	}

	var res *gorm.DB
	if res = r.db.Delete(&Secret{}, "id = ?", id); res.Error != nil {
		return 0, newServiceDatabaseError(res.Error)
	}
	return res.RowsAffected, nil
}
