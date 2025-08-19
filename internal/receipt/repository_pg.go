package receipt

import "gorm.io/gorm"

type RepositoryPg struct {
    db *gorm.DB
}

func NewRepositoryPg(db *gorm.DB) *RepositoryPg {
    return &RepositoryPg{db}
}

func (r *RepositoryPg) Create(receipt *Receipt) error {
    return r.db.Create(receipt).Error
}

func (r *RepositoryPg) FindAll() ([]Receipt, error) {
    var receipts []Receipt
    err := r.db.Find(&receipts).Error
    return receipts, err
}

func (r *RepositoryPg) FindByID(id uint) (*Receipt, error) {
    var receipt Receipt
    err := r.db.First(&receipt, id).Error
    return &receipt, err
}

func (r *RepositoryPg) Update(receipt *Receipt) error {
    return r.db.Save(receipt).Error
}

func (r *RepositoryPg) Delete(id uint) error {
    return r.db.Delete(&Receipt{}, id).Error
}
