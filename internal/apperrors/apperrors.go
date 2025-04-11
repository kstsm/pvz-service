package apperrors

import "errors"

var (
	ErrEmailNotFound              = errors.New("пользователь с таким email не найден")
	ErrEmailAlreadyExists         = errors.New("пользователь с таким email уже существует")
	ErrInvalidCredentials         = errors.New("неверный email или пароль")
	ErrNoActiveReception          = errors.New("нет активной приёмки для данного ПВЗ")
	ErrReceptionAlreadyClosed     = errors.New("приемка уже закрыта или не найдена")
	ErrNoProductToDelete          = errors.New("нет товаров для удаления")
	ErrReceptionAlreadyInProgress = errors.New("невозможно создать приёмку: предыдущая не закрыта")
)
