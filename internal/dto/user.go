package dto

import (
	"simplestforum/internal/delivery/api/apimodel"
	"simplestforum/internal/domain/entity"
	"simplestforum/internal/infrastructure/dbmodel"
)

func UserAddToDB(e *entity.UserAdd) *dbmodel.User {
	if e == nil {
		return nil
	}

	return &dbmodel.User{
		Nickname: e.Nickname,
		ShowInfo: e.ShowInfo,
		Password: e.Password,
	}
}

func UserInfoToDB(e *entity.UserInfo, userID int64) *dbmodel.UserInfo {
	if e == nil {
		return nil
	}

	return &dbmodel.UserInfo{
		UserID:    userID,
		Phone:     e.Phone,
		Email:     e.Email,
		FirstName: e.FirstName,
		LastName:  e.LastName,
	}
}

func UserFromDB(user *dbmodel.User) *entity.User {
	if user == nil {
		return nil
	}

	e := &entity.User{
		ID:        user.ID,
		Nickname:  user.Nickname,
		ShowInfo:  user.ShowInfo,
		Rank:      user.Rank,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if user.Level == nil {
		e.Level = entity.UserLevelNone
	} else {
		e.Level = entity.UserLevel(*user.Level)
	}

	if user.Restriction == nil {
		e.Restriction = entity.UserRestrictionNone
	} else {
		e.Restriction = entity.UserRestriction(*user.Restriction)
	}

	return e
}

func UserInfoFromDB(u *dbmodel.UserInfo) *entity.UserInfo {
	if u == nil {
		return nil
	}

	return &entity.UserInfo{
		Phone:     u.Phone,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

func UsersWithInfoFromDB(u []*dbmodel.UserWithInfo) []*entity.User {
	if u == nil {
		return nil
	}

	e := make([]*entity.User, len(u))

	for i, user := range u {
		e[i] = UserFromDB(&user.User)
		e[i].UserInfo = UserInfoFromDB(&user.UserInfo)
	}

	return e
}

func UserAddFromRest(u *apimodel.AddUserInput) *entity.UserAdd {
	if u == nil {
		return nil
	}

	return &entity.UserAdd{
		Nickname: u.Nickname,
		Password: u.Password,
		ShowInfo: u.ShowInfo,
		UserInfo: &entity.UserInfo{
			Phone:     u.Phone,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
		},
	}
}

func UserEditFromRest(u *apimodel.EditUserInput, userID int64) *entity.UserEdit {
	if u == nil {
		return nil
	}

	e := &entity.UserEdit{
		Nickname:    u.Nickname,
		Password:    u.Password,
		ShowInfo:    u.ShowInfo,
		Rank:        u.Rank,
		CountTopics: u.CountTopics,
		CountPosts:  u.CountPosts,
		Level:       (*entity.UserLevel)(u.Level),
		Restriction: (*entity.UserRestriction)(u.Restriction),
		UserInfo: &entity.UserInfo{
			Phone:     u.Phone,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
		},
	}

	if u.ID == nil {
		e.ID = userID
	} else {
		e.ID = *u.ID
	}

	return e
}

func UserToRest(e *entity.User) *apimodel.User {
	if e == nil {
		return nil
	}

	return &apimodel.User{
		ID:          e.ID,
		Nickname:    e.Nickname,
		ShowInfo:    e.ShowInfo,
		Rank:        e.Rank,
		Level:       apimodel.UserLevel(e.Level),
		Restriction: apimodel.UserRestriction(e.Restriction),
		UserInfo:    UserInfoToRest(e.UserInfo),
		CountTopics: e.CountTopics,
		CountPosts:  e.CountPosts,
		Topics:      TopicsToRest(e.Topics),
		Posts:       PostsToRest(e.Posts),
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func UserInfoToRest(e *entity.UserInfo) *apimodel.UserInfo {
	if e == nil {
		return nil
	}

	return &apimodel.UserInfo{
		Phone:     e.Phone,
		Email:     e.Email,
		FirstName: e.FirstName,
		LastName:  e.LastName,
	}
}

func UsersToRest(e []*entity.User) []*apimodel.User {
	if e == nil {
		return nil
	}

	u := make([]*apimodel.User, len(e))

	for i, user := range e {
		u[i] = UserToRest(user)
	}

	return u
}

func UserFiltersFromRest(u *apimodel.UserFilters) *entity.UserFilters {
	if u == nil {
		return nil
	}

	return &entity.UserFilters{
		IDs:            u.Ids,
		RankFrom:       u.RankFrom,
		RankTo:         u.RankTo,
		Level:          (*entity.UserLevel)(u.Level),
		Restriction:    (*entity.UserRestriction)(u.Restriction),
		CountPostsFrom: u.CountPostsFrom,
		CountPostsTo:   u.CountPostsTo,
	}
}

func UserSortFromRest(u *apimodel.UserSort) *entity.UserSort {
	if u == nil {
		return nil
	}

	return &entity.UserSort{
		By:    entity.UserSortBy(u.By),
		Order: entity.SortOrder(u.Order),
	}
}

func UserEditToDB(e *entity.UserEdit) (*dbmodel.UserUpdate, int64) {
	if e == nil {
		return nil, 0
	}

	userUpdate := &dbmodel.UserUpdate{
		Nickname:    e.Nickname,
		Password:    e.Password,
		ShowInfo:    e.ShowInfo,
		Rank:        e.Rank,
		CountTopics: e.CountTopics,
		CountPosts:  e.CountPosts,
	}

	var ptr *string

	if e.Level != nil {
		if *e.Level != entity.UserLevelNone {
			ptr = (*string)(e.Level)
		}

		userUpdate.Level = &ptr
	}

	if e.Restriction != nil {
		if *e.Restriction != entity.UserRestrictionNone {
			ptr = (*string)(e.Restriction)
		}

		userUpdate.Restriction = &ptr
	}

	return userUpdate, e.ID
}

func UserFiltersToDB(e *entity.UserFilters) *dbmodel.UserFilters {
	if e == nil {
		return nil
	}

	return &dbmodel.UserFilters{
		IDs:            e.IDs,
		RankFrom:       e.RankFrom,
		RankTo:         e.RankTo,
		Level:          (*string)(e.Level),
		Restriction:    (*string)(e.Restriction),
		CountPostsFrom: e.CountPostsFrom,
		CountPostsTo:   e.CountPostsTo,
	}
}
