package user

type Service interface {
	FindAll() (*[]User, error)
	FindById(Id int) (*User, error)
	Create(userRequest UserRequest) (*User, error)
	Update(Id int, userRequest UserRequest) (*User, error)
	Delete(Id int) (*User, error)
}

type service struct {
	userRepository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) FindAll() (*[]User, error) {
	users, err := s.userRepository.FindAll()
	return users, err
}

func (s *service) FindById(Id int) (*User, error) {
	user, err := s.userRepository.FindById(Id)
	return user, err
}

func (s *service) Create(userRequest UserRequest) (*User, error) {
	user := User{
		Name:    userRequest.Name,
		Account: userRequest.Account,
		Phone:   userRequest.Phone,
	}
	newUser, err := s.userRepository.Create(user)
	return newUser, err
}

func (s *service) Update(Id int, userRequest UserRequest) (*User, error) {
	user, err := s.userRepository.FindById(Id)

	if err != nil {
		return nil, err
	}

	user.Name = userRequest.Name
	user.Account = userRequest.Account
	user.Phone = userRequest.Phone
	newUser, err := s.userRepository.Update(*user)
	return newUser, err
}

func (s *service) Delete(Id int) (*User, error) {
	user, err := s.userRepository.Delete(Id)
	return user, err
}
