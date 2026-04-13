package fleet

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// CarGroupService implements business logic for car group operations.
type CarGroupService struct {
	repo CarGroupRepository
}

// NewCarGroupService constructs a CarGroupService backed by the given repository.
func NewCarGroupService(repo CarGroupRepository) *CarGroupService {
	return &CarGroupService{repo: repo}
}

// Create validates the input and creates a new car group.
func (s *CarGroupService) Create(ctx context.Context, input CreateCarGroupInput) (*CarGroup, error) {
	if err := validateCarGroupCreate(input); err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, input)
}

// List retrieves car groups matching the filter.
func (s *CarGroupService) List(ctx context.Context, filter ListCarGroupsFilter) ([]*CarGroup, error) {
	return s.repo.List(ctx, filter)
}

// GetByID retrieves a single car group by its identifier.
func (s *CarGroupService) GetByID(ctx context.Context, id uuid.UUID) (*CarGroup, error) {
	return s.repo.GetByID(ctx, id)
}

// Update validates the input and applies changes to an existing car group.
func (s *CarGroupService) Update(ctx context.Context, id uuid.UUID, input UpdateCarGroupInput) (*CarGroup, error) {
	if err := validateCarGroupUpdate(input); err != nil {
		return nil, err
	}
	return s.repo.Update(ctx, id, input)
}

// Delete soft-deletes a car group, rejecting the request if active vehicles reference it.
func (s *CarGroupService) Delete(ctx context.Context, id uuid.UUID, deletedBy string) error {
	has, err := s.repo.HasActiveVehicles(ctx, id)
	if err != nil {
		return fmt.Errorf("car group service: check active vehicles: %w", err)
	}
	if has {
		return fmt.Errorf("%w: car group still has active vehicles", ErrConflict)
	}
	return s.repo.Delete(ctx, id, deletedBy)
}

// validateCarGroupCreate checks required fields for a new car group.
func validateCarGroupCreate(input CreateCarGroupInput) error {
	if strings.TrimSpace(input.Name) == "" {
		return fmt.Errorf("%w: name is required", ErrValidation)
	}
	if strings.TrimSpace(input.CreatedBy) == "" {
		return fmt.Errorf("%w: created_by is required", ErrValidation)
	}
	return nil
}

// validateCarGroupUpdate checks that at least one field is being changed.
func validateCarGroupUpdate(input UpdateCarGroupInput) error {
	if input.Name != nil && strings.TrimSpace(*input.Name) == "" {
		return fmt.Errorf("%w: name must not be empty", ErrValidation)
	}
	return nil
}
