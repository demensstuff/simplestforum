package repository

import (
	"simplestforum/internal/domain/entity"
	"simplestforum/internal/dto"
	"simplestforum/internal/infrastructure/dbmodel"
	"time"

	"github.com/gocraft/dbr"
)

// SectionRepository represents a Section Repository.
type SectionRepository struct {
	*DBConn
}

// NewSectionRepository instantiates a SectionRepository.
func NewSectionRepository(db *DBConn) *SectionRepository {
	return &SectionRepository{db}
}

// Insert creates a new Section entry in the database and returns a Section object.
func (r *SectionRepository) Insert(sess entity.Session, e *entity.SectionAdd) (*entity.Section, error) {
	section := dto.SectionAddToDB(e)

	err := r.Wrap(sess, func(tx Gateway) error {
		stmt := tx.InsertInto("sections").
			Returning("created_at", "updated_at")

		insertNotNil(stmt, section)

		return stmt.Load(&section)
	})

	return dto.SectionFromDB(section), err
}

// Update modifies an existing User entry.
func (r *SectionRepository) Update(sess entity.Session, e *entity.SectionEdit) error {
	sectionUpdate, id := dto.SectionEditToDB(e)

	return r.Wrap(sess, func(tx Gateway) error {
		stmt := tx.Update("sections").
			Where("id = ?", id).
			Set("updated_at", time.Now())

		updateNotNil(stmt, sectionUpdate)

		_, err := stmt.Exec()

		return err
	})
}

// Delete removes an existing Section (softly).
func (r *SectionRepository) Delete(sess entity.Session, id int64) error {
	return r.Wrap(sess, func(tx Gateway) error {
		_, err := tx.Update("sections").
			Set("deleted_at", time.Now()).
			Where("id = ?", id).
			Exec()

		return err
	})
}

// SelectByID returns a Section by its ID.
func (r *SectionRepository) SelectByID(sess entity.Session, id int64) (*entity.Section, error) {
	var section *dbmodel.Section

	err := r.Wrap(sess, func(tx Gateway) error {
		return tx.Select("*").
			From("sections").
			Where("id = ?", id).
			LoadOne(&section)
	})

	return dto.SectionFromDB(section), err
}

// SelectAll returns all Sections.
func (r *SectionRepository) SelectAll(sess entity.Session, f *entity.SectionFilters, p *entity.Pagination, s *entity.SectionSort) ([]*entity.Section, error) {
	var sections []*dbmodel.Section

	err := r.Wrap(sess, func(tx Gateway) error {
		stmt := tx.Select("*").
			From("sections")
		conditions := []dbr.Builder{dbr.Eq("deleted_at", nil)}

		if f != nil {
			df := dto.SectionFiltersToDB(f)

			conditions = append(conditions, applyFilters(df)...)
		}

		if p != nil {
			stmt.Paginate(uint64(p.Page), uint64(p.Limit))
		}

		if s != nil {
			stmt.OrderDir(dto.SortColumnToDB(string(s.By)), s.Order == entity.SortOrderAsc)
		}

		_, err := stmt.Where(dbr.And(conditions...)).
			Load(&sections)

		return err
	})

	return dto.SectionsFromDB(sections), err
}
