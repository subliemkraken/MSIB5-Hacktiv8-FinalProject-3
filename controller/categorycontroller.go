package controller

import (
	"FinalProject3/config"
	"FinalProject3/dto"
	"FinalProject3/model"
	"FinalProject3/model/validation"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateCategory(c *gin.Context) {
	var category model.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := validation.ValidateStruct(category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := config.DB.Create(&category).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create category"})
		return
	}

	categorydto := dto.Category{
		ID:        category.ID,
		Type:      category.Type,
		CreatedAt: category.CreatedAt,
	}

	c.JSON(http.StatusCreated, categorydto)
}

func GetAllCategory(c *gin.Context) {
	var categories []model.Category

	err := config.DB.Preload("Task").Find(&categories).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var allCategories []dto.GetCategory

	for _, category := range categories {
		var categoryTasks []dto.GetCategoryTask

		for _, task := range category.Task {
			categoryTask := dto.GetCategoryTask{
				ID:          task.ID,
				Title:       task.Title,
				Description: task.Description,
				UserID:      task.UserID,
				CategoryID:  task.CategoryID,
				CreatedAt:   task.CreatedAt,
				UpdatedAt:   task.UpdatedAt,
			}
			categoryTasks = append(categoryTasks, categoryTask)
		}

		category := dto.GetCategory{
			ID:        category.ID,
			Type:      category.Type,
			UpdatedAt: category.UpdatedAt,
			CreatedAt: category.CreatedAt,
			Tasks:     categoryTasks,
		}

		allCategories = append(allCategories, category)
	}

	c.JSON(http.StatusOK, allCategories)
}

func PatchCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("categoryID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid categoryID"})
		return
	}
	var category model.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := validation.ValidateStruct(category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if config.DB.Model(&category).Where("id = ?", categoryID).Updates(&category).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Failed to update category"})
		return
	}

	userdto := dto.Category{
		ID:        uint(categoryID),
		Type:      category.Type,
		CreatedAt: category.UpdatedAt,
	}

	c.JSON(http.StatusOK, userdto)
}

func DeleteCategory(c *gin.Context) {

	categoryID, err := strconv.Atoi(c.Param("categoryID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid categoryID"})
		return
	}

	var category model.Category

	if err := config.DB.Where("id = ?", categoryID).Find(&category).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Category not found"})
		return
	}

	if err := config.DB.Delete(&category, categoryID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category has been succesfully deleted",
	})
}
