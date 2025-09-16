package service

import (
	"git.myservermanager.com/varakh/upda/internal/server/model"
	"git.myservermanager.com/varakh/upda/internal/server/repository"
	"git.myservermanager.com/varakh/upda/internal/server/service_error"
	"github.com/rs/zerolog/log"
)

type CommentService struct {
	repo repository.CommentRepository
}

func NewCommentService(r repository.CommentRepository) *CommentService {
	return &CommentService{
		repo: r,
	}
}

func (s *CommentService) GetByUpdateId(updateId string, page int, pageSize int) ([]*model.Comment, error) {
	if updateId == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	return s.repo.FindAllByUpdateId(updateId, page, pageSize)
}

func (s *CommentService) CountByUpdateId(updateId string) (int64, error) {
	if updateId == "" {
		return 0, service_error.ErrValidationNotBlank
	}
	return s.repo.CountByUpdateId(updateId)
}

func (s *CommentService) GetById(id string) (*model.Comment, error) {
	if id == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	e, err := s.repo.FindById(id)

	if err != nil {
		return nil, err
	}

	return e, nil
}

func (s *CommentService) Delete(id string) error {
	if id == "" {
		return service_error.ErrValidationNotBlank
	}

	e, err := s.GetById(id)
	if err != nil {
		return err
	}

	if _, err = s.repo.Delete(e.ID.String()); err != nil {
		return err
	}

	log.Info().Msgf("Deleted comment '%v'", id)

	return nil
}

func (s *CommentService) Create(author string, content string, update *model.Update) (*model.Comment, error) {
	if author == "" || content == "" {
		return nil, service_error.ErrValidationNotBlank
	}
	if update == nil {
		return nil, service_error.ErrValidationNotEmpty
	}

	var err error
	var e *model.Comment
	if e, err = s.repo.Create(author, content, update.ID.String()); err != nil {
		return nil, err
	} else {
		log.Info().Msg("Created comment")
		return e, nil
	}
}

func (s *CommentService) UpdateContent(id string, content string) (*model.Comment, error) {
	if id == "" || content == "" {
		return nil, service_error.ErrValidationNotBlank
	}

	var e *model.Comment
	var err error

	if e, err = s.GetById(id); err != nil {
		return nil, err
	}

	if e, err = s.repo.UpdateContent(id, content); err != nil {
		return nil, err
	}

	log.Info().Msgf("Modified comment '%v'", id)
	return e, nil
}
