// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Cuisine struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	IconURL   *string `json:"icon_url,omitempty"`
	URL       string  `json:"url"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	User      *User   `json:"user"`
	UserID    int     `json:"user_id"`
}

type CuisineInput struct {
	Title   string  `json:"title"`
	IconURL *string `json:"icon_url,omitempty"`
	URL     string  `json:"url"`
	UserID  int     `json:"user_id"`
}

type Mutation struct {
}

type Query struct {
}

type User struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	IconURL  *string `json:"icon_url,omitempty"`
}

type UserInput struct {
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	IconURL  *string `json:"icon_url,omitempty"`
}
