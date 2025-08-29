package receipt

type Usecase struct {
    repo Repository
}

func NewUsecase(repo Repository) *Usecase {
    return &Usecase{repo}
}

func (u *Usecase) Create(receipt *Receipt) error {
    return u.repo.Create(receipt)
}

func (u *Usecase) GetAll() ([]Receipt, error) {
    return u.repo.FindAll()
}

func (u *Usecase) GetByID(id uint) (*Receipt, error) {
    return u.repo.FindByID(id)
}

func (u *Usecase) Update(receipt *Receipt) error {
    return u.repo.Update(receipt)
}

func (u *Usecase) Delete(id uint) error {
    return u.repo.Delete(id)
}
