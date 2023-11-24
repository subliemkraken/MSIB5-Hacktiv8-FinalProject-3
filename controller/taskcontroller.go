package controller

import (
	"FinalProject3/config"
	"FinalProject3/dto"
	"FinalProject3/model"
	"FinalProject3/model/validation"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GetIdToken(c *gin.Context) uint {

	userData, _ := c.Get("userData")
	tokenClaims := userData.(jwt.MapClaims)

	userID := uint(tokenClaims["id"].(float64))

	return userID
}

func GetAllTask(c *gin.Context) {
	var tasks []model.Task
	err := config.DB.Preload("User").Find(&tasks).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var allTasks []dto.GetTask

	for _, data := range tasks {

		task := dto.GetTask{
			ID:          data.ID,
			Title:       data.Title,
			Status:      data.Status,
			Description: data.Description,
			UserID:      data.UserID,
			CategoryID:  data.CategoryID,
			CreatedAt:   data.CreatedAt,
			User: []dto.GetUserTask{
				{
					ID:        data.User.ID,
					Email:     data.User.Email,
					Full_name: data.User.Full_name,
				},
			},
		}

		allTasks = append(allTasks, task)
	}

	c.JSON(http.StatusOK, allTasks)
}

func CreateTask(c *gin.Context) {
	var (
		task       model.Task
		category   model.Category
		newtaskreq model.NewTaskRequest
	)

	if err := c.ShouldBindJSON(&newtaskreq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	task.UserID = GetIdToken(c)
	categoryID := newtaskreq.CategoryID

	if err := config.DB.First(&category, categoryID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Category not found"})
		return
	}

	if err := validation.ValidateStruct(newtaskreq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	task.Title = newtaskreq.Title
	task.Description = newtaskreq.Description
	task.CategoryID = newtaskreq.CategoryID
	task.Status = false

	if err := config.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	taskdto := dto.Task{
		ID:          task.ID,
		Title:       task.Title,
		Status:      task.Status,
		Description: task.Description,
		UserID:      task.UserID,
		CategoryID:  task.CategoryID,
		CreatedAt:   task.CreatedAt,
	}

	c.JSON(http.StatusCreated, taskdto)
}

func PutTask(c *gin.Context) {
	var (
		task    model.Task
		puttask model.PutRequest
	)

	userID := GetIdToken(c)

	taskID, err := strconv.Atoi(c.Param("taskID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid taskID"})
		return
	}

	if err := c.ShouldBindJSON(&puttask); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := config.DB.First(&task, taskID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Task not found"})
		return
	}

	if task.UserID != userID {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "This task does not belong to you"})
		return
	}

	if err := validation.ValidateStruct(puttask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	task.Title = puttask.Title
	task.Description = puttask.Description

	if config.DB.Model(&task).Where("id = ?", taskID).Updates(&task).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Failed to update task"})
		return
	}

	taskdto := dto.UpdateTask{
		ID:          task.ID,
		Title:       task.Title,
		Status:      task.Status,
		Description: task.Description,
		UserID:      task.UserID,
		CategoryID:  task.CategoryID,
		UpdatedAt:   task.UpdatedAt,
	}

	c.JSON(http.StatusOK, taskdto)
}

func PatchStat(c *gin.Context) {
	var (
		task         model.Task
		patchstatask model.PatchStatusRequest
	)

	userID := GetIdToken(c)

	taskID, err := strconv.Atoi(c.Param("taskID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid taskID"})
		return
	}

	if err := c.ShouldBindJSON(&patchstatask); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := config.DB.First(&task, taskID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Task not found"})
		return
	}

	if task.UserID != userID {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "This task does not belong to you"})
		return
	}

	if err := validation.ValidateStruct(patchstatask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	task.Status = patchstatask.Status

	if config.DB.Model(&task).Where("id = ?", taskID).Updates(&task).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Failed to update task"})
		return
	}

	taskdto := dto.UpdateTask{
		ID:          task.ID,
		Title:       task.Title,
		Status:      task.Status,
		Description: task.Description,
		UserID:      task.UserID,
		CategoryID:  task.CategoryID,
		UpdatedAt:   task.UpdatedAt,
	}

	c.JSON(http.StatusOK, taskdto)
}
func PatchCatId(c *gin.Context) {
	var (
		task       model.Task
		patchcatid model.PatchCatIdRequest
	)

	userID := GetIdToken(c)

	taskID, err := strconv.Atoi(c.Param("taskID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid taskID"})
		return
	}

	if err := c.ShouldBindJSON(&patchcatid); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := config.DB.First(&task, taskID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Task not found"})
		return
	}

	if task.UserID != userID {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "This task does not belong to you"})
		return
	}

	if err := validation.ValidateStruct(patchcatid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	task.CategoryID = patchcatid.CategoryID

	if config.DB.Model(&task).Where("id = ?", taskID).Updates(&task).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Failed to update task"})
		return
	}

	taskdto := dto.UpdateTask{
		ID:          task.ID,
		Title:       task.Title,
		Status:      task.Status,
		Description: task.Description,
		UserID:      task.UserID,
		CategoryID:  task.CategoryID,
		UpdatedAt:   task.UpdatedAt,
	}

	c.JSON(http.StatusOK, taskdto)
}

func DeleteTask(c *gin.Context) {
	var task model.Task
	userID := GetIdToken(c)

	taskID, err := strconv.Atoi(c.Param("taskID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid taskID"})
		return
	}

	if err := config.DB.Where("id = ?", taskID).Find(&task).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Taskid not found"})
		return
	}

	if userID != task.UserID {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "this task not belong to you"})
		return
	}

	if err := config.DB.Delete(&task, taskID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete task"})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task has been succesfully deleted",
	})

}
