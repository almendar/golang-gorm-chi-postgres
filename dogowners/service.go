package dogowners

type Service struct {
	storage *Storage
}

func NewService(storage *Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) GetDogowner(userID uint) (Owner, error) {

	user, err := s.storage.Owner(userID)
	if err != nil {
		return Owner{}, err
	}

	dogOwner := Owner{
		ID:   user.ID,
		Name: user.Name,
		Dogs: []Dog{},
	}
	for _, dog := range user.Dogs {
		dogOwner.Dogs = append(dogOwner.Dogs, Dog{
			ID:       dog.ID,
			Name:     dog.Name,
			Birthday: dog.Birthday,
		})
	}

	return dogOwner, nil
}

func (s *Service) SaveDog(ownerID uint, dog Dog) error {
	dogDB := DogDBModel{
		OwnerID:  ownerID,
		Name:     dog.Name,
		Birthday: dog.Birthday,
	}
	return s.storage.SaveDog(&dogDB)
}

func (s *Service) ListOwners() ([]Owner, error) {
	var users []OwnerDBModel
	err := s.storage.db.Preload("Dogs").Find(&users).Error
	if err != nil {
		return nil, err
	}

	var owners []Owner
	for _, user := range users {
		owner := Owner{
			ID:   user.ID,
			Name: user.Name,
			Dogs: []Dog{},
		}
		for _, dog := range user.Dogs {
			owner.Dogs = append(owner.Dogs, Dog{
				ID:       dog.ID,
				Name:     dog.Name,
				Birthday: dog.Birthday,
			})
		}
		owners = append(owners, owner)
	}

	return owners, nil
}
