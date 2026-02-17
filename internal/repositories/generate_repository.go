package repositories

import (
	"fmt"
	"golang-api/internal/models"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

type GenerateRepository struct {
	db *gorm.DB
}

func NewGenerateRepository(db *gorm.DB) *GenerateRepository {
	return &GenerateRepository{db: db}
}

func (r *GenerateRepository) FindByAlias(alias string) (*models.Generate, error) {
	var generate models.Generate
	err := r.db.Unscoped().Where("alias = ?", alias).First(&generate).Error
	if err != nil {
		return nil, err
	}
	return &generate, nil
}

func (r *GenerateRepository) Update(generate *models.Generate) error {
	return r.db.Save(generate).Error
}

func (r *GenerateRepository) GetCode(alias string, isNotPreview bool) string {
	gen, err := r.FindByAlias(alias)
	if err != nil || gen == nil {
		return fmt.Sprintf("ER-%05d", rand.Intn(90000)+10000)
	}

	now := time.Now()
	date := now.Format("060102")
	separator := gen.Separator

	separatorTime, err := time.Parse("060102", separator)
	if err != nil {
		separatorTime = now
	}
	separatorStr := separatorTime.Format("060102")

	if gen.Queue == 9999 || date[:4] != separatorStr[:4] {
		gen.Queue = 1
		gen.Separator = date
	}

	queue := fmt.Sprintf("%s%04d%s", date[:4], gen.Queue, date[4:6])

	if gen.Prefix != nil && *gen.Prefix != "" {
		queue = *gen.Prefix + queue
	}
	if gen.Suffix != nil && *gen.Suffix != "" {
		queue = queue + *gen.Suffix
	}

	if isNotPreview {
		gen.Queue += 1
		r.Update(gen)
	}

	return queue
}
