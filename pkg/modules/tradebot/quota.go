// Package tradebot package tradebot
package tradebot

import (
	"math"
	"time"

	"gitlab.tocraw.com/root/toc_trader/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/database"
	"gitlab.tocraw.com/root/toc_trader/pkg/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/traderecord"
)

// TradeQuota TradeQuota
var TradeQuota int64 = 1000000

const (
	// TradeFeeRatio TradeFeeRatio
	TradeFeeRatio float64 = 0.001425
	// FeeDiscount FeeDiscount
	FeeDiscount float64 = 0.35
	// TradeTaxRatio TradeTaxRatio
	TradeTaxRatio float64 = 0.0015
)

// InitStartUpQuota InitStartUpQuota
func InitStartUpQuota() {
	if time.Now().Day() == global.TradeDay.Day() {
		realOrderArr, err := traderecord.GetAllorderByDayTime(global.TradeDay, database.GetAgent())
		if err != nil {
			logger.GetLogger().Panic(err)
		}
		for _, v := range realOrderArr {
			if v.Status != 6 {
				continue
			}
			record := traderecord.TradeRecord{
				StockNum:  v.StockNum,
				StockName: global.AllStockNameMap.GetValueByKey(v.StockNum),
				Action:    v.Action,
				Price:     v.Price,
				Quantity:  v.Quantity,
				Status:    v.Status,
				OrderID:   v.OrderID,
				TradeTime: time.Now(),
			}
			if v.Action == 1 {
				switch {
				case !FilledBuyOrderMap.CheckStockExist(v.StockNum) && !FilledBuyLaterOrderMap.CheckStockExist(v.StockNum):
					if FilledSellFirstOrderMap.CheckStockExist(v.StockNum) {
						FilledBuyLaterOrderMap.Set(record)
						logger.GetLogger().Warnf("Filled Buy Later: %s", record.StockNum)
					} else {
						FilledBuyOrderMap.Set(record)
						TradeQuota -= GetStockBuyCost(v.Price, v.Quantity)
						logger.GetLogger().Warnf("Filled Buy: %s", record.StockNum)
					}
				case !FilledBuyLaterOrderMap.CheckStockExist(v.StockNum) && FilledBuyOrderMap.CheckStockExist(v.StockNum):
					FilledBuyLaterOrderMap.Set(record)
					logger.GetLogger().Warnf("Filled Buy Later: %s", record.StockNum)
				case FilledBuyLaterOrderMap.CheckStockExist(v.StockNum) && !FilledBuyOrderMap.CheckStockExist(v.StockNum):
					FilledBuyOrderMap.Set(record)
					TradeQuota -= GetStockBuyCost(v.Price, v.Quantity)
					logger.GetLogger().Warnf("Filled Buy: %s", record.StockNum)
				}
			} else if v.Action == 2 {
				switch {
				case !FilledSellOrderMap.CheckStockExist(v.StockNum) && !FilledSellFirstOrderMap.CheckStockExist(v.StockNum):
					if FilledBuyOrderMap.CheckStockExist(v.StockNum) {
						FilledSellOrderMap.Set(record)
						logger.GetLogger().Warnf("Filled Sell: %s", record.StockNum)
					} else {
						FilledSellFirstOrderMap.Set(record)
						TradeQuota -= GetStockBuyCost(v.Price, v.Quantity)
						logger.GetLogger().Warnf("Filled Sell First %s", record.StockNum)
					}
				case !FilledSellOrderMap.CheckStockExist(v.StockNum) && FilledSellFirstOrderMap.CheckStockExist(v.StockNum):
					FilledSellOrderMap.Set(record)
					logger.GetLogger().Warnf("Filled Sell: %s", record.StockNum)
				case FilledSellOrderMap.CheckStockExist(v.StockNum) && !FilledSellFirstOrderMap.CheckStockExist(v.StockNum):
					FilledSellFirstOrderMap.Set(record)
					TradeQuota -= GetStockBuyCost(v.Price, v.Quantity)
					logger.GetLogger().Warnf("Filled Sell First %s", record.StockNum)
				}
			}
		}
		logger.GetLogger().Warnf("Initial Quota: %d", TradeQuota)
	}
}

// FindUnfinishedStock FindUnfinishedStock
func FindUnfinishedStock() {
	filledBuyOrder := FilledBuyOrderMap.GetAllRecordMap()
	for stockNum, record := range filledBuyOrder {
		if !FilledSellOrderMap.CheckStockExist(stockNum) && !BuyOrderMap.CheckStockExist(stockNum) {
			BuyOrderMap.Set(record)
			logger.GetLogger().Warnf("Unfinished buy %s", record.StockNum)
		}
	}
	filledSellFirstOrder := FilledSellFirstOrderMap.GetAllRecordMap()
	for stockNum, record := range filledSellFirstOrder {
		if !FilledBuyLaterOrderMap.CheckStockExist(stockNum) && !SellFirstOrderMap.CheckStockExist(stockNum) {
			SellFirstOrderMap.Set(record)
			logger.GetLogger().Warnf("Unfinished sell first %s", record.StockNum)
		}
	}
}

// GetStockBuyCost GetStockBuyCost
func GetStockBuyCost(price float64, qty int64) int64 {
	return int64(math.Ceil(price*float64(qty)*1000) + math.Floor(price*float64(qty)*1000*TradeFeeRatio))
}

// GetStockSellCost GetStockSellCost
func GetStockSellCost(price float64, qty int64) int64 {
	return int64(math.Ceil(price*float64(qty)*1000) - math.Floor(price*float64(qty)*1000*TradeFeeRatio) - math.Floor(price*float64(qty)*1000*TradeTaxRatio))
}

// GetStockTradeFeeDiscount GetStockTradeFeeDiscount
func GetStockTradeFeeDiscount(price float64, qty int64) int64 {
	return int64(math.Floor(price*float64(qty)*1000*TradeFeeRatio) * (1 - FeeDiscount))
}
