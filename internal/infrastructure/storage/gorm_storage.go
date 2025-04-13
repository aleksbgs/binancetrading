package storage

import (
	"log"
	"time"

	"binancetrading/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// CandlestickModel represents the database model for a candlestick
type CandlestickModel struct {
	Symbol    string `gorm:"primaryKey"`
	Timestamp string `gorm:"primaryKey"`
	Open      string
	High      string
	Low       string
	Close     string
}

// GORMStorage implements the Storage interface using GORM
type GORMStorage struct {
	db *gorm.DB
}

func NewGORMStorage(dsn string) (*GORMStorage, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&CandlestickModel{}); err != nil {
		return nil, err
	}

	return &GORMStorage{db: db}, nil
}

func (s *GORMStorage) SaveCandlestick(symbol string, c *domain.OHLC) error {
	candle := CandlestickModel{
		Symbol:    symbol,
		Timestamp: c.Timestamp.Format(time.RFC3339),
		Open:      c.Open,
		High:      c.High,
		Low:       c.Low,
		Close:     c.Close,
	}

	if err := s.db.Save(&candle).Error; err != nil {
		return err
	}

	log.Printf("Saved candlestick for %s at %s", symbol, c.Timestamp)
	return nil
}

func (s *GORMStorage) GetCandlesticks(symbol, startTime, endTime string) ([]domain.OHLC, error) {
	var models []CandlestickModel
	if err := s.db.Where("symbol = ? AND timestamp BETWEEN ? AND ?", symbol, startTime, endTime).Find(&models).Error; err != nil {
		return nil, err
	}

	candlesticks := make([]domain.OHLC, len(models))
	for i, m := range models {
		ts, _ := time.Parse(time.RFC3339, m.Timestamp)
		candlesticks[i] = domain.OHLC{
			Symbol:    m.Symbol,
			Open:      m.Open,
			High:      m.High,
			Low:       m.Low,
			Close:     m.Close,
			Timestamp: ts,
		}
	}
	return candlesticks, nil
}

func (s *GORMStorage) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
