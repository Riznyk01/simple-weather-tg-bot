package repository

type UserControlMemoryStorage struct {
	memoryStor *MemoryStorage
}

func NewUserControlMemoryStorage(memoryStor *MemoryStorage) *UserControlMemoryStorage {
	return &UserControlMemoryStorage{
		memoryStor: memoryStor,
	}
}

func (uc *UserControlMemoryStorage) BanUser(userId int64) error {
	currentData := uc.memoryStor.data[userId]
	currentData.Ban = true
	uc.memoryStor.data[userId] = currentData
	return nil
}

func (uc *UserControlMemoryStorage) UnbanUser(userId int64) error {
	currentData := uc.memoryStor.data[userId]
	currentData.Ban = false
	uc.memoryStor.data[userId] = currentData
	return nil
}
