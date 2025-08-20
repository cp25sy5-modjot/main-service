package receipt

type Service struct {
    repo Repository
}

func NewService(repo Repository) *Service {
    return &Service{repo}
}

func (s *Service) Create(receipt *Receipt) error {
    return s.repo.Create(receipt)
}

func (s *Service) GetAll() ([]Receipt, error) {
    return s.repo.FindAll()
}

func (s *Service) GetByID(id uint) (*Receipt, error) {
    return s.repo.FindByID(id)
}

func (s *Service) Update(receipt *Receipt) error {
    return s.repo.Update(receipt)
}

func (s *Service) Delete(id uint) error {
    return s.repo.Delete(id)
}
