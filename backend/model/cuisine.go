package model

import "time"

type Cuisine struct {
	ID        uint      `json:"id" gorm:"primaryKey"`  //主キーになる
	Title     string    `json:"title" gorm:"not null"` //空の値を許可しない
	IconURL   *string   `json:"icon_url"`
	URL       string    `json:"url"`
	Comment   string    `json:"comment"` // コメント追加
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `json:"user" gorm:"foreignKey:UserId; constraint:OnDelete:CASCADE"` //userを削除したときにuserに紐づいている料理も消去される
	UserId    uint      `json:"user_id" gorm:"not null"`
}

type CuisineResponse struct {
	ID        uint      `json:"id" gorm:"primaryKey"`  //主キーになる
	Title     string    `json:"title" gorm:"not null"` //空の値を許可しない
	IconURL   *string   `json:"icon_url"`
	URL       string    `json:"url"`
	Comment   string    `json:"comment"` // コメント追加
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserId    uint      `json:"user_id"`
}
