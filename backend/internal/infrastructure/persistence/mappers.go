package persistence

import (
	"dangde-world/backend/internal/domain/activity"
	"dangde-world/backend/internal/domain/activitydata"
	"dangde-world/backend/internal/domain/assignment"
	"dangde-world/backend/internal/domain/category"
	"dangde-world/backend/internal/domain/user"
)

func toUserEntity(model userModel) user.User {
	return user.User{
		ID:            model.ID,
		Name:          model.Name,
		Role:          model.Role,
		Avatar:        model.Avatar,
		ParentID:      model.ParentID,
		PreferredLang: model.PreferredLang,
	}
}

func toUserModel(entity user.User) userModel {
	return userModel{
		ID:            entity.ID,
		Name:          entity.Name,
		Role:          entity.Role,
		Avatar:        entity.Avatar,
		ParentID:      entity.ParentID,
		PreferredLang: entity.PreferredLang,
	}
}

func toCategoryEntity(model categoryModel) category.Category {
	return category.Category{
		ID:       model.ID,
		Name:     model.Name,
		Slug:     model.Slug,
		Type:     model.Type,
		ParentID: model.ParentID,
	}
}

func toCategoryModel(entity category.Category) categoryModel {
	return categoryModel{
		ID:       entity.ID,
		Name:     entity.Name,
		Slug:     entity.Slug,
		Type:     entity.Type,
		ParentID: entity.ParentID,
	}
}

func toActivityEntity(model activityModel) activity.Activity {
	return activity.Activity{
		ID:          model.ID,
		Title:       model.Title,
		Description: model.Description,
		Language:    model.Language,
		Difficulty:  model.Difficulty,
		AgeGroup:    model.AgeGroup,
		CategoryID:  model.CategoryID,
		Category:    toCategoryEntity(model.Category),
		Prompt:      model.Prompt,
		Icon:        model.Icon,
		HTMLCode:    model.HTMLCode,
		CSSCode:     model.CSSCode,
		JSCode:      model.JSCode,
	}
}

func toActivityModel(entity activity.Activity) activityModel {
	return activityModel{
		ID:          entity.ID,
		Title:       entity.Title,
		Description: entity.Description,
		Language:    entity.Language,
		Difficulty:  entity.Difficulty,
		AgeGroup:    entity.AgeGroup,
		CategoryID:  entity.CategoryID,
		Prompt:      entity.Prompt,
		Icon:        entity.Icon,
		HTMLCode:    entity.HTMLCode,
		CSSCode:     entity.CSSCode,
		JSCode:      entity.JSCode,
	}
}

func toActivityDataEntity(model activityDataModel) activitydata.Item {
	return activitydata.Item{
		ID:           model.ID,
		ActivityID:   model.ActivityID,
		AssignmentID: model.AssignmentID,
		KidID:        model.KidID,
		Key:          model.Key,
		Value:        model.Value,
	}
}

func toActivityDataModel(entity activitydata.Item) activityDataModel {
	return activityDataModel{
		ID:           entity.ID,
		ActivityID:   entity.ActivityID,
		AssignmentID: entity.AssignmentID,
		KidID:        entity.KidID,
		Key:          entity.Key,
		Value:        entity.Value,
	}
}

func toAssignmentEntity(model assignmentModel) assignment.Assignment {
	return assignment.Assignment{
		ID:         model.ID,
		ParentID:   model.ParentID,
		Parent:     toUserEntity(model.Parent),
		KidID:      model.KidID,
		Kid:        toUserEntity(model.Kid),
		ActivityID: model.ActivityID,
		Activity:   toActivityEntity(model.Activity),
		Status:     model.Status,
		Progress:   model.Progress,
	}
}

func toAssignmentModel(entity assignment.Assignment) assignmentModel {
	return assignmentModel{
		ID:         entity.ID,
		ParentID:   entity.ParentID,
		KidID:      entity.KidID,
		ActivityID: entity.ActivityID,
		Status:     entity.Status,
		Progress:   entity.Progress,
	}
}
