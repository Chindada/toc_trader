// Package streamtickprocess package streamtickprocess
package streamtickprocess

import (
	"github.com/markcheno/go-quote"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/analyzestreamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/streamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tickanalyze"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tradebot"
	"gitlab.tocraw.com/root/toc_trader/tools/common"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
)

// TickProcess TickProcess
func TickProcess(lastClose float64, cond simulationcond.AnalyzeCondition, ch chan *streamtick.StreamTick, saveCh chan []*streamtick.StreamTick) {
	var sellChan, buyLaterChan chan *streamtick.StreamTick
	if lastClose == 0 {
		return
	}
	analyzeTickChan := make(chan *analyzestreamtick.AnalyzeStreamTick)
	go tradebot.TradeAgent(analyzeTickChan, cond)

	if global.TradeSwitch.Buy {
		buyLaterChan = make(chan *streamtick.StreamTick)
		go tradebot.BuyLaterBot(buyLaterChan, cond)
	}
	if global.TradeSwitch.Sell {
		sellChan = make(chan *streamtick.StreamTick)
		go tradebot.SellBot(sellChan, cond)
	}
	var input quote.Quote
	var unSavedTicks streamtick.PtrArrArr
	var tmpArr streamtick.PtrArr
	var lastSaveLastClose, openChangeRatio float64
	for {
		tick := <-ch
		if openChangeRatio == 0 {
			openChangeRatio = common.Round((tick.Open - lastClose), 2)
		}
		tmpArr = append(tmpArr, tick)
		switch {
		case tradebot.FilledBuyOrderMap.CheckStockExist(tick.StockNum):
			sellChan <- tick
		case tradebot.FilledSellFirstOrderMap.CheckStockExist(tick.StockNum):
			buyLaterChan <- tick
		}

		if tmpArr.GetTotalTime() < cond.TicksPeriodThreshold {
			continue
		}
		if tmpArr.GetTotalTime() > cond.TicksPeriodLimit {
			unSavedTicks.ClearAll()
		}
		unSavedTicks.Append(tmpArr)
		saveCh <- tmpArr
		tmpArr = []*streamtick.StreamTick{}

		if unSavedTicks.GetCount() >= cond.TicksPeriodCount {
			var outSum, inSum int64
			var totalTime float64
			data := unSavedTicks.Get()
			for _, v := range data {
				input.Close = append(input.Close, v.GetAllCloseArr()...)
				outSum += v.GetOutSum()
				inSum += v.GetInSum()
				totalTime += v.GetTotalTime()
			}
			if len(input.Close) < int(cond.HistoryCloseCount) {
				unSavedTicks.ClearAll()
				continue
			} else {
				input.Close = input.Close[len(input.Close)-int(cond.HistoryCloseCount):]
			}
			rsi, err := tickanalyze.GenerateRSI(input)
			if err != nil {
				logger.Logger.Errorf("TickProcess Stock: %s, Err: %s", tick.StockNum, err)
				continue
			}

			closeDiff := common.Round((unSavedTicks.GetLastClose() - lastSaveLastClose), 2)
			if lastSaveLastClose == 0 {
				closeDiff = 0
			}
			lastSaveLastClose = unSavedTicks.GetLastClose()
			unSavedTicksInOutRatio := common.Round((100 * (float64(outSum) / float64(outSum+inSum))), 2)
			analyze := analyzestreamtick.AnalyzeStreamTick{
				TimeStamp:        tick.TimeStamp,
				StockNum:         tick.StockNum,
				Close:            tick.Close,
				OpenChangeRatio:  openChangeRatio,
				CloseChangeRatio: tick.PctChg,
				OutSum:           outSum,
				InSum:            inSum,
				OutInRatio:       unSavedTicksInOutRatio,
				TotalTime:        totalTime,
				CloseDiff:        closeDiff,
				Open:             tick.Open,
				AvgPrice:         tick.AvgPrice,
				High:             tick.High,
				Low:              tick.Low,
				Rsi:              rsi,
				Volume:           outSum + inSum,
			}
			analyzeTickChan <- &analyze
			unSavedTicks.ClearAll()
		}
	}
}

// SaveStreamTicks SaveStreamTicks
func SaveStreamTicks(saveCh chan []*streamtick.StreamTick) {
	for {
		unSavedTicks := <-saveCh
		if len(unSavedTicks) != 0 {
			if err := streamtick.InsertMultiRecord(unSavedTicks, global.GlobalDB); err != nil {
				logger.Logger.Error(err)
				continue
			}
		}
	}
}
