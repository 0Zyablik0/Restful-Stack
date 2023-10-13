package models

type User struct {
	Id uint64 `gorm:"primaryKey"`
}

type UserStacks struct {
	UserId    uint64 `gorm:"primaryKey;references:User.Id"`
	StackSize uint64
}

type StackEntity struct {
	UserId   uint64 `gorm:"references:User.Id;primaryKey;autoIncrement:false"`
	StackPos uint64 `gorm:"primaryKey;autoIncrement:false"`
	Object   int64
}
