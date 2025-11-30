package factory

import (
	"database/sql"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	postgresRepo "github.com/EduGoGroup/edugo-api-administracion/internal/infrastructure/persistence/postgres/repository"
)

// postgresRepositoryFactory crea repositorios usando PostgreSQL
type postgresRepositoryFactory struct {
	db *sql.DB
}

// NewPostgresRepositoryFactory crea un factory para repositorios PostgreSQL
func NewPostgresRepositoryFactory(db *sql.DB) RepositoryFactory {
	return &postgresRepositoryFactory{db: db}
}

func (f *postgresRepositoryFactory) CreateSchoolRepository() repository.SchoolRepository {
	return postgresRepo.NewPostgresSchoolRepository(f.db)
}

func (f *postgresRepositoryFactory) CreateUserRepository() repository.UserRepository {
	return postgresRepo.NewPostgresUserRepository(f.db)
}

func (f *postgresRepositoryFactory) CreateAcademicUnitRepository() repository.AcademicUnitRepository {
	return postgresRepo.NewPostgresAcademicUnitRepository(f.db)
}

func (f *postgresRepositoryFactory) CreateUnitMembershipRepository() repository.UnitMembershipRepository {
	return postgresRepo.NewPostgresUnitMembershipRepository(f.db)
}

func (f *postgresRepositoryFactory) CreateUnitRepository() repository.UnitRepository {
	return postgresRepo.NewPostgresUnitRepository(f.db)
}

func (f *postgresRepositoryFactory) CreateSubjectRepository() repository.SubjectRepository {
	return postgresRepo.NewPostgresSubjectRepository(f.db)
}

func (f *postgresRepositoryFactory) CreateMaterialRepository() repository.MaterialRepository {
	return postgresRepo.NewPostgresMaterialRepository(f.db)
}

func (f *postgresRepositoryFactory) CreateStatsRepository() repository.StatsRepository {
	return postgresRepo.NewPostgresStatsRepository(f.db)
}

func (f *postgresRepositoryFactory) CreateGuardianRepository() repository.GuardianRepository {
	return postgresRepo.NewPostgresGuardianRepository(f.db)
}
