package repository

import (
	"simplestforum/internal/domain/entity"
	"simplestforum/internal/dto"
	"simplestforum/internal/infrastructure/dbmodel"
	"time"

	"github.com/gocraft/dbr"
)

// PostRepository represents a Post Repository.
type PostRepository struct {
	*DBConn
}

// NewPostRepository instantiates a PostRepository.
func NewPostRepository(db *DBConn) *PostRepository {
	return &PostRepository{db}
}

// Insert creates a new Post entry in the database and returns a Post object.
func (r *PostRepository) Insert(sess entity.Session, e *entity.PostAdd) (int64, error) {
	var postID int64

	post := dto.PostAddToDB(e)

	err := r.Wrap(sess, func(tx Gateway) error {
		stmt := tx.InsertInto("posts").
			Returning("id")

		insertNotNil(stmt, post)

		return stmt.Load(&postID)
	})

	return postID, err
}

// Update modifies an existing Post entry.
func (r *PostRepository) Update(sess entity.Session, e *entity.PostEdit) error {
	postUpdate, id := dto.PostEditToDB(e)

	return r.Wrap(sess, func(tx Gateway) error {
		stmt := tx.Update("posts").
			Where("id = ?", id).
			Set("updated_at", time.Now())

		updateNotNil(stmt, postUpdate)

		_, err := stmt.Exec()

		return err
	})
}

// Delete removes existing Posts (softly).
func (r *PostRepository) Delete(sess entity.Session, ids ...int64) error {
	return r.Wrap(sess, func(tx Gateway) error {
		_, err := tx.Update("posts").
			Set("deleted_at", time.Now()).
			Where(dbr.Eq("id", ids)).
			Exec()

		return err
	})
}

// SelectByID returns a Section by its ID.
func (r *PostRepository) SelectByID(sess entity.Session, id int64) (*entity.Post, error) {
	var post *dbmodel.Post

	err := r.Wrap(sess, func(tx Gateway) error {
		return tx.Select("*").
			From("posts").
			Where("id = ?", id).
			LoadOne(&post)
	})

	return dto.PostFromDB(post), err
}

// SelectAll returns all Posts.
func (r *PostRepository) SelectAll(sess entity.Session, f *entity.PostFilters, p *entity.Pagination, s *entity.PostSort) ([]*entity.Post, error) {
	var posts []*dbmodel.Post

	err := r.Wrap(sess, func(tx Gateway) error {
		stmt := tx.Select("*").
			From("posts")
		conditions := []dbr.Builder{dbr.Eq("deleted_at", nil)}

		if f != nil {
			df := dto.PostFiltersToDB(f)

			conditions = append(conditions, applyFilters(df)...)
		}

		if p != nil {
			stmt.Paginate(uint64(p.Page), uint64(p.Limit))
		}

		if s != nil {
			stmt.OrderDir(dto.SortColumnToDB(string(s.By)), s.Order == entity.SortOrderAsc)
		}

		_, err := stmt.Where(dbr.And(conditions...)).
			Load(&posts)

		return err
	})

	return dto.PostsFromDB(posts), err
}
func (r *PostRepository) IDsToDelete(sess entity.Session, e *entity.PostDelete) ([]int64, error) {
	var ids []int64

	err := r.Wrap(sess, func(tx Gateway) error {
		stmt := tx.Select("id").
			From("posts")

		conditions := []dbr.Builder{dbr.Eq("deleted_at", nil)}

		if e != nil {
			df := dto.PostDeleteToDB(e)

			conditions = append(conditions, applyFilters(df)...)
		}

		_, err := stmt.Where(dbr.And(conditions...)).
			Load(&ids)

		return err
	})

	return ids, err
}
