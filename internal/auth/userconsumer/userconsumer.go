package userconsumer

type UserConsumer interface {
	CheckUser(username, password string) bool
	CreateUser(username, password string) error
}

type userConsumer struct {
}

func New() UserConsumer {
	return &userConsumer{}
}

func (u *userConsumer) CheckUser(username, password string) bool {

}

func (u *userConsumer) CreateUser(username, password string) error {
	return nil
}
