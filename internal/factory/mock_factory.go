package factory

import (
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	mockRepo "github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/persistence/mock/repository"
)

// mockRepositoryFactory crea repositorios usando implementaciones mock en memoria
type mockRepositoryFactory struct{}

// NewMockRepositoryFactory crea un factory para repositorios mock
func NewMockRepositoryFactory() RepositoryFactory {
	return &mockRepositoryFactory{}
}

func (f *mockRepositoryFactory) CreateSchoolRepository() repository.SchoolRepository {
	return mockRepo.NewMockSchoolRepository()
}

func (f *mockRepositoryFactory) CreateUserRepository() repository.UserRepository {
	return mockRepo.NewMockUserRepository()
}

// Sprint 2 - Pendientes
func (f *mockRepositoryFactory) CreateAcademicUnitRepository() repository.AcademicUnitRepository {
	return nil // TODO: Implementar en Sprint 2
}

func (f *mockRepositoryFactory) CreateUnitMembershipRepository() repository.UnitMembershipRepository {
	return nil // TODO: Implementar en Sprint 2
}

func (f *mockRepositoryFactory) CreateUnitRepository() repository.UnitRepository {
	return nil // TODO: Implementar en Sprint 2
}

func (f *mockRepositoryFactory) CreateSubjectRepository() repository.SubjectRepository {
	return nil // TODO: Implementar en Sprint 2
}

// Sprint 3 - Pendientes
func (f *mockRepositoryFactory) CreateMaterialRepository() repository.MaterialRepository {
	return nil // TODO: Implementar en Sprint 3
}

func (f *mockRepositoryFactory) CreateStatsRepository() repository.StatsRepository {
	return nil // TODO: Implementar en Sprint 3
}

func (f *mockRepositoryFactory) CreateGuardianRepository() repository.GuardianRepository {
	return nil // TODO: Implementar en Sprint 3
}
