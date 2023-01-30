package repository

import (
	"simplestforum/internal/domain/entity"
	"simplestforum/internal/dto"
	"simplestforum/internal/infrastructure/dbmodel"
	"time"

	"github.com/gocraft/dbr"
)

// UserRepository represents a User Repository.
type UserRepository struct {
	*DBConn
}

// NewUserRepository instantiates a UserRepository.
func NewUserRepository(db *DBConn) *UserRepository {
	return &UserRepository{db}
}

// Insert creates a new User entry in the database and returns a User object.
func (r *UserRepository) Insert(sess entity.Session, e *entity.UserAdd) (*entity.User, error) {
	user := dto.UserAddToDB(e)

	err := r.Wrap(sess, func(tx Gateway) error {
		stmt := tx.InsertInto("users").
			Returning("id", "rank", "created_at", "updated_at")

		insertNotNil(stmt, user)

		return stmt.Load(&user)
	})

	return dto.UserFromDB(user), err
}

// InsertInfo creates a new entry in the 'users_info' database.
func (r *UserRepository) InsertInfo(sess entity.Session, e *entity.UserInfo, userID int64) error {
	userInfo := dto.UserInfoToDB(e, userID)

	return r.Wrap(sess, func(tx Gateway) error {
		stmt := tx.InsertInto("users_info")

		insertNotNil(stmt, userInfo)

		_, err := stmt.Exec()

		return err
	})
}

// Update modifies an existing User entry.
func (r *UserRepository) Update(sess entity.Session, e *entity.UserEdit) error {
	userUpdate, id := dto.UserEditToDB(e)

	return r.Wrap(sess, func(tx Gateway) error {
		stmt := tx.Update("users").
			Where("id = ?", id).
			Set("updated_at", time.Now())

		updateNotNil(stmt, userUpdate)

		_, err := stmt.Exec()

		return err
	})
}

// UpdateInfo modifies an existing UserInfo entry.
func (r *UserRepository) UpdateInfo(sess entity.Session, e *entity.UserInfo, userID int64) error {
	userInfo := dto.UserInfoToDB(e, userID)

	return r.Wrap(sess, func(tx Gateway) error {
		stmt := tx.Update("users_info").
			Where("user_id", userInfo.UserID)

		updateNotNil(stmt, userInfo)

		_, err := stmt.Exec()

		return err
	})
}

// Delete removes an existing User (softly).
func (r *UserRepository) Delete(sess entity.Session, id int64) error {
	return r.Wrap(sess, func(tx Gateway) error {
		_, err := tx.Update("users").
			Set("deleted_at", time.Now()).
			Where("id = ?", id).
			Exec()

		return err
	})
}

// SelectAll returns all Users.
func (r *UserRepository) SelectAll(sess entity.Session, f *entity.UserFilters, p *entity.Pagination, s *entity.UserSort) ([]*entity.User, error) {
	return r.selectAll(sess, f, p, s, false)
}

// SelectAllWithInfo returns all Users with UserInfo.
func (r *UserRepository) SelectAllWithInfo(sess entity.Session, f *entity.UserFilters, p *entity.Pagination, s *entity.UserSort) ([]*entity.User, error) {
	return r.selectAll(sess, f, p, s, true)
}

// SelectByID returns a User by its ID.
func (r *UserRepository) SelectByID(sess entity.Session, id int64) (*entity.User, error) {
	var user *dbmodel.User

	err := r.Wrap(sess, func(tx Gateway) error {
		return tx.Select("*").
			From("users").
			Where("id = ?", id).
			LoadOne(&user)
	})

	return dto.UserFromDB(user), err
}

// SelectByNicknameWithPassword returns a User and its encrypted password by its nickname.
func (r *UserRepository) SelectByNicknameWithPassword(sess entity.Session, nickname string) (*entity.User, string, error) {
	var user *dbmodel.User

	err := r.Wrap(sess, func(tx Gateway) error {
		err := tx.Select("*").
			From("users").
			Where("nickname = ?", nickname).
			LoadOne(&user)

		return err
	})

	if err != nil {
		return nil, "", err
	}

	return dto.UserFromDB(user), user.Password, nil
}

func (r *UserRepository) selectAll(sess entity.Session, f *entity.UserFilters, p *entity.Pagination, s *entity.UserSort, withInfo bool) ([]*entity.User, error) {
	var users []*dbmodel.UserWithInfo

	err := r.Wrap(sess, func(tx Gateway) error {
		stmt := tx.Select("*").
			From("users")
		conditions := []dbr.Builder{dbr.Eq("deleted_at", nil)}

		if withInfo {
			stmt.Join("users_info", "users.id = users_info.user_id")
		}

		if f != nil {
			df := dto.UserFiltersToDB(f)

			conditions = append(conditions, applyFilters(df)...)
		}

		if p != nil {
			stmt.Paginate(uint64(p.Page), uint64(p.Limit))
		}

		if s != nil {
			stmt.OrderDir(dto.SortColumnToDB(string(s.By)), s.Order == entity.SortOrderAsc)
		}

		_, err := stmt.Where(dbr.And(conditions...)).
			Load(&users)

		return err
	})

	return dto.UsersWithInfoFromDB(users), err
}
