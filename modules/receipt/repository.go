package receipt

type Repository interface {
    Create(receipt *Receipt) error
    FindAll() ([]Receipt, error)
    FindByID(id uint) (*Receipt, error)
    Update(receipt *Receipt) error
    Delete(id uint) error
}
