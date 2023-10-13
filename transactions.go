package main

import (
	"restful_stack/models"

	"gorm.io/gorm"
)

func TransactionStackPush(
	tx *gorm.DB,
	user *models.User,
	stackRequest *models.StackPushRequest,
	stackResponse *models.StackPushResponse,
) error {
	stack := &models.UserStacks{UserId: user.Id}
	if err := tx.Transaction(func(tx *gorm.DB) error {
		return TransactionGetStack(tx, user, stack)
	}); err != nil {
		return err
	}

	stackResponse.StackSize = stack.StackSize

	stackEntity := models.StackEntity{
		UserId:   user.Id,
		StackPos: stack.StackSize,
		Object:   stackRequest.Value,
	}
	stack.StackSize += 1

	tx.SavePoint("sp")
	if err := tx.Create(&stackEntity).Error; err != nil {
		return err
	}

	if err := tx.Save(&stack).Error; err != nil {
		tx.RollbackTo("sp")
		return err
	}
	stackResponse.StackSize = stack.StackSize
	stackResponse.Code = StatusOK
	stackResponse.Err = ""

	return nil
}

func TransactionStackPushRange(
	tx *gorm.DB,
	user *models.User,
	stackRequest *models.StackPushRangeRequest,
	stackResponse *models.StackPushRangeResponse,
) error {
	stack := &models.UserStacks{UserId: user.Id}
	if err := tx.Transaction(func(tx *gorm.DB) error {
		return TransactionGetStack(tx, user, stack)
	}); err != nil {
		return err
	}
	stackResponse.StackSize = stack.StackSize

	stackEntities := make([]*models.StackEntity, 0, len(stackRequest.Values))
	for i := 0; i < len(stackRequest.Values); i++ {
		stackEntities = append(stackEntities, &models.StackEntity{
			UserId:   user.Id,
			StackPos: stack.StackSize + uint64(i),
			Object:   stackRequest.Values[i],
		})
	}
	stack.StackSize += uint64(len(stackRequest.Values))

	tx.SavePoint("sp")
	if err := tx.Create(stackEntities).Error; err != nil {
		return err
	}

	if err := tx.Save(&stack).Error; err != nil {
		tx.RollbackTo("sp")
		return err
	}

	stackResponse.StackSize = stack.StackSize
	stackResponse.Code = StatusOK
	stackResponse.Err = ""

	return nil
}

func TransactionStackTop(
	tx *gorm.DB,
	user *models.User,
	response *models.StackTopResponse,
) error {
	stack := &models.UserStacks{UserId: user.Id}

	err := tx.Transaction(func(tx *gorm.DB) error {
		return TransactionGetStack(tx, user, stack)
	})
	if err != nil {
		return err
	}
	if stack.StackSize == 0 {
		response.StackSize = 0
		response.Code = StatusError
		response.Err = "Stack is empty"
		return nil
	}
	response.StackSize = stack.StackSize
	stackEntity := &models.StackEntity{
		UserId:   user.Id,
		StackPos: stack.StackSize - 1,
	}
	tx.SavePoint("sp")
	if err := tx.First(stackEntity).Error; err != nil {
		tx.RollbackTo("sp")
		return err
	}
	response.Code = StatusOK
	response.Value = stackEntity.Object
	response.Err = ""
	return nil
}

func TransactionStackPop(
	tx *gorm.DB,
	user *models.User,
	response *models.StackPopResponse,
) error {
	stack := &models.UserStacks{UserId: user.Id}

	err := tx.Transaction(func(tx *gorm.DB) error {
		return TransactionGetStack(tx, user, stack)
	})
	if err != nil {
		return err
	}
	if stack.StackSize == 0 {
		response.StackSize = 0
		response.Code = StatusError
		response.Err = "Stack is empty"
		return nil
	}

	stackEntity := &models.StackEntity{
		UserId:   user.Id,
		StackPos: stack.StackSize - 1,
	}
	response.StackSize = stack.StackSize
	tx.SavePoint("sp")
	if err := tx.First(stackEntity).Error; err != nil {
		tx.RollbackTo("sp")
		return err
	}

	if err := tx.Delete(stackEntity).Error; err != nil {
		tx.RollbackTo("sp")
		return err
	}
	stack.StackSize--
	if err := tx.Save(stack).Error; err != nil {
		tx.RollbackTo("sp")
		return err
	}
	response.StackSize--
	response.Code = StatusOK
	response.Value = stackEntity.Object
	response.Err = ""
	return nil
}

func TransactionGetStack(tx *gorm.DB,
	user *models.User,
	stack *models.UserStacks,
) error {
	var stacks []models.UserStacks
	if err := tx.Find(&stacks, stack.UserId).Error; err != nil {
		return err
	}

	if len(stacks) == 0 {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		if err := tx.Create(stack).Error; err != nil {
			return err
		}
	} else {
		*stack = stacks[0]
	}
	return nil

}
