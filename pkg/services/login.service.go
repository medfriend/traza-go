package services

import (
     "encoding/json"
    "traza-go/pkg/entity"
    "traza-go/pkg/repository"
    "log"
)

type LoginService interface {
    Save(data []byte) error
}

type loginServiceImpl struct {
    loginRepository repository.LoginRepository
}

func NewLoginService(loginRepository repository.LoginRepository) LoginService {
    return &loginServiceImpl{
        loginRepository: loginRepository,
    }
}

type UsuarioIDMapper struct {
	UsuarioID uint `json:"usuario_id"`
}

func (s *loginServiceImpl) Save(data []byte) error {

	var mapper UsuarioIDMapper
	if err := json.Unmarshal(data, &mapper); err != nil {
		log.Printf("Error al extraer UsuarioID: %+v", err)
		return err
	}

	login := entity.Login{
		UsuarioID: &mapper.UsuarioID,
        Caducada: true,
	}

	return s.loginRepository.Save(&login)
}