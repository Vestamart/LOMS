package repository

import (
	"context"
	"encoding/json"
	"github.com/vestamart/loms/internal/domain"
	"github.com/vestamart/loms/internal/localErr"
	"os"
)

// Error

type SKUID = uint32

type StocksRepository = map[SKUID]domain.StocksItem

type InMemoryStocksRepository struct {
	stocksRepository StocksRepository
}

// Создание репозитория и считываение данных в json файла
func NewInMemoryStocksRepositoryFromFile() (*InMemoryStocksRepository, error) {

	file, err := os.Open("stock-data.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var jsonStocks []struct {
		SKU        uint32 `json:"sku"`
		TotalCount uint32 `json:"total_count"`
		Reserved   uint32 `json:"reserved"`
	}

	if err = json.NewDecoder(file).Decode(&jsonStocks); err != nil {
		return nil, err
	}

	repo := &InMemoryStocksRepository{
		stocksRepository: make(StocksRepository, len(jsonStocks)),
	}

	for _, item := range jsonStocks {
		repo.stocksRepository[item.SKU] = domain.StocksItem{
			TotalCount: item.TotalCount,
			Reserved:   item.Reserved,
		}
	}

	return repo, nil
}

func (r *InMemoryStocksRepository) Reserve(_ context.Context, sku SKUID, count uint32) error {
	v, ok := r.stocksRepository[sku]
	if !ok {
		return localErr.SKUNotExistErr
	}

	v.Reserved += count
	if v.Reserved > v.TotalCount {
		return localErr.ItemNotEnoughErr
	}

	r.stocksRepository[sku] = v
	return nil
}

func (r *InMemoryStocksRepository) ReserveRemove(_ context.Context, sku SKUID, count uint32) error {
	v, ok := r.stocksRepository[sku]
	if !ok {
		return localErr.SKUNotExistErr
	}

	v.TotalCount -= count
	v.Reserved -= count
	r.stocksRepository[sku] = v
	return nil
}

func (r *InMemoryStocksRepository) ReserveCancel(_ context.Context, skus map[SKUID]uint32) error {
	for sku, count := range skus {
		v, ok := r.stocksRepository[sku]
		if !ok {
			return localErr.SKUNotExistErr
		}
		v.Reserved -= count
		r.stocksRepository[sku] = v
	}
	return nil
}

func (r *InMemoryStocksRepository) GetBySKU(_ context.Context, sku SKUID) (uint32, error) {
	v, ok := r.stocksRepository[sku]
	if !ok {
		return 0, localErr.SKUNotExistErr
	}

	return v.TotalCount - v.Reserved, nil
}

func (r *InMemoryStocksRepository) RollbackReserve(_ context.Context, skus map[SKUID]uint32) error {
	if skus == nil {
		return nil
	}
	for sku, count := range skus {
		stock := r.stocksRepository[sku]

		stock.Reserved -= count
		r.stocksRepository[sku] = stock
	}
	return nil
}
