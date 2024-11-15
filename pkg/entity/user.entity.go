package entity

import (
	"time"
)

func (User) TableName() string {
	return "usuario"
}

type User struct {
	UsuarioID          uint       `gorm:"primaryKey;autoIncrement;column:usuario_id" json:"usuario_id"`
	Usuario            uint       `gorm:"column:usuario" json:"usuario"`
	Nombre1            string     `gorm:"size:100;column:nombre_1" json:"nombre_1"`
	Nombre2            string     `gorm:"size:100;column:nombre_2" json:"nombre_2"`
	ApellidoPaterno    string     `gorm:"size:100;column:apellido_paterno" json:"apellido_paterno"`
	ApellidoMaterno    string     `gorm:"size:100;column:apellido_materno" json:"apellido_materno"`
	Clave              string     `gorm:"not null;column:clave" json:"clave"`
	Email              string     `gorm:"size:100;unique;not null;column:email" json:"email"`
	CambiarClave       bool       `gorm:"default:false;column:cambiar_clave" json:"cambiar_clave"`
	FechaCambioClave   time.Time  `gorm:"column:fecha_cambio_clave" json:"fecha_cambio_clave"`
	FechaCreacion      time.Time  `gorm:"autoCreateTime;column:fecha_creacion" json:"fecha_creacion"`
	FechaRetiro        *time.Time `gorm:"column:fecha_retiro" json:"fecha_retiro"`
	Activo             bool       `gorm:"default:true;column:activo" json:"activo"`
	TiempoValidezToken string     `gorm:"column:tiempo_valides_token" json:"tiempo_valides_token"`
}
