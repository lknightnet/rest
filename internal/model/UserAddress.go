package model

type UserAddress struct {
	ID            int
	Street        string
	Flat          *string
	House         string
	DoorphoneCode *string
	UserID        int
}

func (u *UserAddress) GetUserID() int {
	return u.UserID
}
