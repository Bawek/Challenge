package person_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/izymalhaw/go-crud/yishakterefe/internal/domain"
	"github.com/izymalhaw/go-crud/yishakterefe/internal/ports"
)

type PersonService struct {
	PersonRepo ports.PersonRepository
}

func NewPersonSvc(personRepo ports.PersonRepository) *PersonService {
	return &PersonService{PersonRepo: personRepo}
}

func (s *PersonService) CreatePerson(ctx context.Context, person domain.Person) error {
	return s.PersonRepo.CreatePerson(ctx, person)
}

func (s *PersonService) GetAllPersons(ctx context.Context, limit, offset int) ([]domain.Person, error) {
	return s.PersonRepo.GetAllPersons(ctx, limit, offset)
}

func (s *PersonService) GetPerson(ctx context.Context, id uuid.UUID) (domain.Person, error) {
	data, err := s.PersonRepo.GetPerson(ctx, id)
	if err != nil {
		return domain.Person{}, err
	}
	return domain.Person{
		Id:      data.Id,
		Name:    data.Name,
		Age:     data.Age,
		Hobbies: data.Hobbies,
	}, nil
}

func (s *PersonService) UpdatePerson(ctx context.Context, person domain.Person) error {
	return s.PersonRepo.UpdatePerson(ctx, person)
}

func (s *PersonService) DeletePerson(ctx context.Context, id uuid.UUID) error {
	return s.PersonRepo.DeletePerson(ctx, id)
}
