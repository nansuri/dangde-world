package persistence

import (
	"dangde-world/backend/internal/domain/activity"
	"dangde-world/backend/internal/domain/activitydata"
	"dangde-world/backend/internal/domain/assignment"
	"dangde-world/backend/internal/domain/category"
	"dangde-world/backend/internal/domain/user"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{db: db}
}

func (r UserRepository) FindByID(id uint) (user.User, error) {
	var model userModel
	if err := r.db.First(&model, id).Error; err != nil {
		return user.User{}, err
	}
	return toUserEntity(model), nil
}

func (r UserRepository) List(role string, parentID *uint) ([]user.User, error) {
	var models []userModel
	query := r.db.Order("id asc")
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	}
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	items := make([]user.User, 0, len(models))
	for _, model := range models {
		items = append(items, toUserEntity(model))
	}
	return items, nil
}

func (r UserRepository) CreateBatch(items []user.User) error {
	models := make([]userModel, 0, len(items))
	for _, item := range items {
		models = append(models, toUserModel(item))
	}
	return r.db.Create(&models).Error
}

func (r UserRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&userModel{}).Count(&count).Error
	return count, err
}

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return CategoryRepository{db: db}
}

func (r CategoryRepository) List(categoryType string) ([]category.Category, error) {
	var models []categoryModel
	query := r.db.Order("id asc")
	if categoryType != "" {
		query = query.Where("type = ?", categoryType)
	}
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	items := make([]category.Category, 0, len(models))
	for _, model := range models {
		items = append(items, toCategoryEntity(model))
	}
	return items, nil
}

func (r CategoryRepository) CreateBatch(items []category.Category) error {
	models := make([]categoryModel, 0, len(items))
	for _, item := range items {
		models = append(models, toCategoryModel(item))
	}
	return r.db.Create(&models).Error
}

type ActivityRepository struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return ActivityRepository{db: db}
}

func (r ActivityRepository) List(categoryID *uint) ([]activity.Activity, error) {
	var models []activityModel
	query := r.db.Preload("Category").Order("id asc")
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	items := make([]activity.Activity, 0, len(models))
	for _, model := range models {
		items = append(items, toActivityEntity(model))
	}
	return items, nil
}

func (r ActivityRepository) FindByID(id uint) (activity.Activity, error) {
	var model activityModel
	if err := r.db.Preload("Category").First(&model, id).Error; err != nil {
		return activity.Activity{}, err
	}
	return toActivityEntity(model), nil
}

func (r ActivityRepository) Create(item activity.Activity) (activity.Activity, error) {
	model := toActivityModel(item)
	if err := r.db.Create(&model).Error; err != nil {
		return activity.Activity{}, err
	}
	return r.FindByID(model.ID)
}

func (r ActivityRepository) Update(item activity.Activity) (activity.Activity, error) {
	model := toActivityModel(item)
	if err := r.db.Model(&activityModel{}).Where("id = ?", item.ID).Updates(&model).Error; err != nil {
		return activity.Activity{}, err
	}
	return r.FindByID(item.ID)
}

func (r ActivityRepository) CreateBatch(items []activity.Activity) error {
	models := make([]activityModel, 0, len(items))
	for _, item := range items {
		models = append(models, toActivityModel(item))
	}
	return r.db.Create(&models).Error
}

type AssignmentRepository struct {
	db *gorm.DB
}

func NewAssignmentRepository(db *gorm.DB) AssignmentRepository {
	return AssignmentRepository{db: db}
}

func (r AssignmentRepository) List(kidID, parentID *uint) ([]assignment.Assignment, error) {
	var models []assignmentModel
	query := r.db.Preload("Activity").Preload("Activity.Category").Preload("Kid").Preload("Parent").Order("id asc")
	if kidID != nil {
		query = query.Where("kid_id = ?", *kidID)
	}
	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	}
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	items := make([]assignment.Assignment, 0, len(models))
	for _, model := range models {
		items = append(items, toAssignmentEntity(model))
	}
	return items, nil
}

func (r AssignmentRepository) FindByID(id uint) (assignment.Assignment, error) {
	var model assignmentModel
	if err := r.db.Preload("Activity").Preload("Activity.Category").Preload("Kid").Preload("Parent").First(&model, id).Error; err != nil {
		return assignment.Assignment{}, err
	}
	return toAssignmentEntity(model), nil
}

func (r AssignmentRepository) Create(item assignment.Assignment) (assignment.Assignment, error) {
	model := toAssignmentModel(item)
	if err := r.db.Create(&model).Error; err != nil {
		return assignment.Assignment{}, err
	}
	return r.FindByID(model.ID)
}

func (r AssignmentRepository) Update(item assignment.Assignment) (assignment.Assignment, error) {
	model := toAssignmentModel(item)
	if err := r.db.Model(&assignmentModel{}).Where("id = ?", item.ID).Updates(&model).Error; err != nil {
		return assignment.Assignment{}, err
	}
	return r.FindByID(item.ID)
}

func (r AssignmentRepository) CreateBatch(items []assignment.Assignment) error {
	models := make([]assignmentModel, 0, len(items))
	for _, item := range items {
		models = append(models, toAssignmentModel(item))
	}
	return r.db.Create(&models).Error
}

type ActivityDataRepository struct {
	db *gorm.DB
}

func NewActivityDataRepository(db *gorm.DB) ActivityDataRepository {
	return ActivityDataRepository{db: db}
}

func (r ActivityDataRepository) List(activityID, assignmentID, kidID *uint, key string) ([]activitydata.Item, error) {
	var models []activityDataModel
	query := r.db.Order("id asc")
	if activityID != nil {
		query = query.Where("activity_id = ?", *activityID)
	}
	if assignmentID != nil {
		query = query.Where("assignment_id = ?", *assignmentID)
	}
	if kidID != nil {
		query = query.Where("kid_id = ?", *kidID)
	}
	if key != "" {
		query = query.Where("key = ?", key)
	}
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	items := make([]activitydata.Item, 0, len(models))
	for _, model := range models {
		items = append(items, toActivityDataEntity(model))
	}
	return items, nil
}

func (r ActivityDataRepository) FindByID(id uint) (activitydata.Item, error) {
	var model activityDataModel
	if err := r.db.First(&model, id).Error; err != nil {
		return activitydata.Item{}, err
	}
	return toActivityDataEntity(model), nil
}

func (r ActivityDataRepository) Create(item activitydata.Item) (activitydata.Item, error) {
	model := toActivityDataModel(item)
	if err := r.db.Create(&model).Error; err != nil {
		return activitydata.Item{}, err
	}
	return r.FindByID(model.ID)
}

func (r ActivityDataRepository) Update(item activitydata.Item) (activitydata.Item, error) {
	model := toActivityDataModel(item)
	if err := r.db.Model(&activityDataModel{}).Where("id = ?", item.ID).Updates(&model).Error; err != nil {
		return activitydata.Item{}, err
	}
	return r.FindByID(item.ID)
}

func (r ActivityDataRepository) Delete(id uint) error {
	return r.db.Delete(&activityDataModel{}, id).Error
}
