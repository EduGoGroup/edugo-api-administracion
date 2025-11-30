package factory

import (
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
)

// RepositoryFactory es el factory para crear repositorios
// Implementa el patrón Abstract Factory para soportar múltiples
// implementaciones (PostgreSQL, Mock, MongoDB futuro, etc.)
type RepositoryFactory interface {
	CreateSchoolRepository() repository.SchoolRepository
	CreateUserRepository() repository.UserRepository
	CreateAcademicUnitRepository() repository.AcademicUnitRepository
	CreateUnitMembershipRepository() repository.UnitMembershipRepository
	CreateUnitRepository() repository.UnitRepository
	CreateSubjectRepository() repository.SubjectRepository
	CreateMaterialRepository() repository.MaterialRepository
	CreateStatsRepository() repository.StatsRepository
	CreateGuardianRepository() repository.GuardianRepository
}
