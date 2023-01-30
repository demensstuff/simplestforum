package repository

import (
	"simplestforum/internal/domain/entity"
	"simplestforum/internal/dto"
	"simplestforum/internal/infrastructure/dbmodel"
	"time"

	"github.com/gocraft/dbr"
)

// TopicRepository represents a Topic Repository.
type TopicRepository struct {
	*DBConn
}

// NewTopicRepository instantiates a TopicRepository.
func NewTopicRepository(db *DBConn) *TopicRepository {
	return &TopicRepository{db}
}

// Insert creates a new Topic entry in the database and returns a Topic object.
func (r *TopicRepository) Insert(sess entity.Session, e *entity.TopicAdd) (int64, error) {
	var topicID int64

	topic := dto.TopicAddToDB(e)

	err := r.Wrap(sess, func(tx Gateway) error {
		stmt := tx.InsertInto("topics").
			Returning("id")

		insertNotNil(stmt, topic)

		return stmt.Load(&topicID)
	})

	return topicID, err
}

// Update modifies an existing Topic entry.
func (r *TopicRepository) Update(sess entity.Session, e *entity.TopicEdit) error {
	topicUpdate, id := dto.TopicEditToDB(e)

	return r.Wrap(sess, func(tx Gateway) error {
		stmt := tx.Update("topics").
			Where("id = ?", id).
			Set("updated_at", time.Now())

		updateNotNil(stmt, topicUpdate)

		_, err := stmt.Exec()

		return err
	})
}

// Delete removes existing Topics (softly).
func (r *TopicRepository) Delete(sess entity.Session, ids ...int64) error {
	return r.Wrap(sess, func(tx Gateway) error {
		_, err := tx.Update("topics").
			Set("deleted_at", time.Now()).
			Where(dbr.Eq("id", ids)).
			Exec()

		return err
	})
}

// SelectByID returns a Section by its ID.
func (r *TopicRepository) SelectByID(sess entity.Session, id int64) (*entity.Topic, error) {
	var topic *dbmodel.Topic

	err := r.Wrap(sess, func(tx Gateway) error {
		return tx.Select("*").
			From("topics").
			Where("id = ?", id).
			LoadOne(&topic)
	})

	return dto.TopicFromDB(topic), err
}

// SelectAll returns all Topics.
func (r *TopicRepository) SelectAll(sess entity.Session, f *entity.TopicFilters, p *entity.Pagination, s *entity.TopicSort) ([]*entity.Topic, error) {
	var topics []*dbmodel.Topic

	err := r.Wrap(sess, func(tx Gateway) error {
		stmt := tx.Select("*").
			From("topics")
		conditions := []dbr.Builder{dbr.Eq("deleted_at", nil)}

		if f != nil {
			df := dto.TopicFiltersToDB(f)

			conditions = append(conditions, applyFilters(df)...)
		}

		if p != nil {
			stmt.Paginate(uint64(p.Page), uint64(p.Limit))
		}

		if s != nil {
			stmt.OrderDir(dto.SortColumnToDB(string(s.By)), s.Order == entity.SortOrderAsc)
		}

		_, err := stmt.Where(dbr.And(conditions...)).
			Load(&topics)

		return err
	})

	return dto.TopicsFromDB(topics), err
}
func (r *TopicRepository) IDsToDelete(sess entity.Session, e *entity.TopicDelete) ([]int64, error) {
	var ids []int64

	err := r.Wrap(sess, func(tx Gateway) error {
		stmt := tx.Select("id").
			From("topics")

		conditions := []dbr.Builder{dbr.Eq("deleted_at", nil)}

		if e != nil {
			df := dto.TopicDeleteToDB(e)

			conditions = append(conditions, applyFilters(df)...)
		}

		_, err := stmt.Where(dbr.And(conditions...)).
			Load(&ids)

		return err
	})

	return ids, err
}
