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

// Sprint 2 - Implementados
func (f *mockRepositoryFactory) CreateAcademicUnitRepository() repository.AcademicUnitRepository {
	return mockRepo.NewMockAcademicUnitRepository()
}

func (f *mockRepositoryFactory) CreateUnitMembershipRepository() repository.UnitMembershipRepository {
	return mockRepo.NewMockUnitMembershipRepository()
}

func (f *mockRepositoryFactory) CreateUnitRepository() repository.UnitRepository {
	return mockRepo.NewMockUnitRepository()
}

func (f *mockRepositoryFactory) CreateSubjectRepository() repository.SubjectRepository {
	return mockRepo.NewMockSubjectRepository()
}

// Sprint 3 - Implementados (100% completo)
func (f *mockRepositoryFactory) CreateMaterialRepository() repository.MaterialRepository {
	return mockRepo.NewMockMaterialRepository()
}

func (f *mockRepositoryFactory) CreateStatsRepository() repository.StatsRepository {
	return mockRepo.NewMockStatsRepository()
}

func (f *mockRepositoryFactory) CreateGuardianRepository() repository.GuardianRepository {
	return mockRepo.NewMockGuardianRepository()
}
