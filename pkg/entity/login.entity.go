package entity

import "time"

func (Login) TableName() string {
	return "trazabilidad-usuario-login"
}

type Login struct {
    LoginID      uint      `gorm:"primaryKey;autoIncrement;column:trazabilidad-usuario-login_id"`
    FechaIngreso  time.Time  `gorm:"autoCreateTime;column:fecha_ingreso" json:"fecha_ingreso"`
    UsuarioID    *uint     `gorm:"column:usuario_id" json:"UsuarioID"`
    FechaSalida  time.Time  `gorm:"autoCreateTime;column:fecha_salida" json:"fecha_salida"`
    Caducada     bool      `gorm:"column:caducada; json:"caducada"`
    Usuario      *User `gorm:"foreignKey:UsuarioID" json:"usuario,omitempty"`
}
