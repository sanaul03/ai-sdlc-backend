package cargroup

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Service provides business logic for car groups.
type Service interface {
	Create(ctx context.Context, req CreateRequest, createdBy string) (CarGroup, error)
	List(ctx context.Context, filter ListFilter) ([]CarGroup, error)
	GetByID(ctx context.Context, id uuid.UUID) (CarGroup, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateRequest, updatedBy string) (CarGroup, error)
	Delete(ctx context.Context, id uuid.UUID, deletedBy string) error
}

type service struct {
	repo Repository
}

// NewService returns a Service backed by the provided Repository.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req CreateRequest, createdBy string) (CarGroup, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return CarGroup{}, fmt.Errorf("%w: name is required", ErrInvalidInput)
	}

	now := time.Now().UTC()
	cg := CarGroup{
		ID:           uuid.New(),
		Name:         name,
		Description:  req.Description,
		SizeCategory: req.SizeCategory,
		CreatedAt:    now,
		UpdatedAt:    now,
		CreatedBy:    createdBy,
		UpdatedBy:    createdBy,
		Deleted:      false,
	}

	created, err := s.repo.Create(ctx, cg)
	if err != nil {
		return CarGroup{}, err
	}
	return created, nil
}

func (s *service) List(ctx context.Context, filter ListFilter) ([]CarGroup, error) {
	return s.repo.List(ctx, filter)
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (CarGroup, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, req UpdateRequest, updatedBy string) (CarGroup, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return CarGroup{}, err
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return CarGroup{}, fmt.Errorf("%w: name cannot be empty", ErrInvalidInput)
		}
		existing.Name = name
	}
	if req.Description != nil {
		existing.Description = req.Description
	}
	if req.SizeCategory != nil {
		existing.SizeCategory = req.SizeCategory
	}

	existing.UpdatedAt = time.Now().UTC()
	existing.UpdatedBy = updatedBy

	return s.repo.Update(ctx, existing)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID, deletedBy string) error {
	has, err := s.repo.HasActiveVehicles(ctx, id)
	if err != nil {
		return err
	}
	if has {
		return ErrHasVehicles
	}
	return s.repo.SoftDelete(ctx, id, deletedBy)
}
