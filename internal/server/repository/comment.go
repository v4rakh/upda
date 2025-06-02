package repository

import (
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"gorm.io/gorm"
)

type CommentRepository interface {
	FindById(id string) (*model.Comment, error)
	FindAllByUpdateId(updateId string, page int, pageSize int) ([]*model.Comment, error)
	CountByUpdateId(updateId string) (int64, error)
	Create(author string, content string, updateId string) (*model.Comment, error)
	UpdateContent(id string, content string) (*model.Comment, error)
	Delete(id string) (int64, error)
}

type CommentDbRepo struct {
	db *gorm.DB
}

func NewCommentDbRepo(db *gorm.DB) *CommentDbRepo {
	return &CommentDbRepo{
		db: db,
	}
}

func (r *CommentDbRepo) FindById(id string) (*model.Comment, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e model.Comment
	var res *gorm.DB
	if res = r.db.Find(&e, "id = ?", id); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, service_error.ErrResourceNotFound
	}

	return &e, nil
}

func (r *CommentDbRepo) Create(author string, content string, updateId string) (*model.Comment, error) {
	if author == "" || content == "" || updateId == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	e := &model.Comment{
		Author:   author,
		Content:  content,
		UpdateID: updateId,
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

func (r *CommentDbRepo) UpdateContent(id string, content string) (*model.Comment, error) {
	if id == "" || content == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var err error
	var e *model.Comment

	if e, err = r.FindById(id); err != nil {
		return nil, err
	}

	e.Content = content

	var res *gorm.DB
	if res = r.db.Save(&e); res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return e, service_error.ErrDatabaseRowsExpected
	}

	return e, nil
}

func (r *CommentDbRepo) Delete(id string) (int64, error) {
	if id == "" {
		return 0, service_error.ErrValidationNotBlank
	}

	var res *gorm.DB
	if res = r.db.Delete(&model.Comment{}, "id = ?", id); res.Error != nil {
		return 0, service_error.NewServiceDatabaseError(res.Error)
	}

	return res.RowsAffected, nil
}

func (r *CommentDbRepo) FindAllByUpdateId(updateId string, page int, pageSize int) ([]*model.Comment, error) {
	if updateId == "" {
		return nil, service_error.ErrValidationLimitGreaterZero
	}
	if page == 0 {
		return nil, service_error.ErrValidationPageGreaterZero
	}
	if pageSize <= 0 {
		return nil, service_error.ErrValidationPageSizeGreaterZero
	}
	offset := (page - 1) * pageSize

	var e []*model.Comment

	if res := r.db.Model(&model.Comment{}).
		Scopes(allGetCommentCriterion(updateId)).
		Order("created_at desc").
		Offset(offset).
		Limit(pageSize).
		Find(&e); res.Error != nil {
		return nil, service_error.NewServiceDatabaseError(res.Error)
	}

	return e, nil
}

func (r *CommentDbRepo) CountByUpdateId(updateId string) (int64, error) {
	if updateId == "" {
		return 0, service_error.ErrValidationLimitGreaterZero
	}
	var c int64

	if res := r.db.Model(&model.Comment{}).Scopes(allGetCommentCriterion(updateId)).Count(&c); res.Error != nil {
		return 0, service_error.NewServiceDatabaseError(res.Error)
	}

	return c, nil
}

func criterionCommentUpdateId(updateId string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("update_id = ?", updateId)
	}
}

func allGetCommentCriterion(updateId string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(criterionCommentUpdateId(updateId))
	}
}
