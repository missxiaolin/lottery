package services

import (
	"lottery/dao"
	"lottery/datasource"
	"lottery/models"
)

type GiftService interface {
	GetAll() []models.LtGift
	CountAll() int64
	//Search(country string) []models.LtGift
	Get(id int, useCache bool) *models.LtGift
	Delete(id int) error
	Update(data *models.LtGift, columns []string) error
	Create(data *models.LtGift) error
}

type giftService struct {
	dao *dao.GiftDao
}

func NewGiftService() GiftService {
	return &giftService{
		dao: dao.NewGiftDao(datasource.InstanceDbMaster()),
	}
}

func (s *giftService) GetAll() []models.LtGift {
	return nil
}
func (s *giftService) CountAll() int64 {
	return s.dao.CountAll()
}

//Search(country string) []models.LtGift
func (s *giftService) Get(id int, useCache bool) *models.LtGift{
	return s.dao.Get(id)
}

func (s *giftService) Delete(id int) error {
	return s.dao.Delete(id)
}

func (s *giftService) Update(data *models.LtGift, columns []string) error {
	return s.dao.Update(data, columns)
}

func (s *giftService) Create(data *models.LtGift) error{
	return s.dao.Create(data)
}