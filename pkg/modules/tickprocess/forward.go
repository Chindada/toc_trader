// Package tickprocess package tickprocess
package tickprocess

import (
	"github.com/markcheno/go-quote"
	"gitlab.tocraw.com/root/toc_trader/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/analyzestreamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/streamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tickanalyze"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tradebot"
	"gitlab.tocraw.com/root/toc_trader/pkg/utils"
)

// ForwardTickProcess ForwardTickProcess
func ForwardTickProcess(lastClose float64, cond simulationcond.AnalyzeCondition, ch chan *streamtick.StreamTick, saveCh chan []*streamtick.StreamTick) {
	var tradeSwitch bool
	var input quote.Quote
	var unSavedTicks streamtick.PtrArrArr
	var tmpArr streamtick.PtrArr
	var lastSaveLastClose, openChangeRatio float64
	var sellChan chan *streamtick.StreamTick
	if lastClose == 0 {
		return
	}
	analyzeTickChan := make(chan *analyzestreamtick.AnalyzeStreamTick)
	go tradebot.BuyAgent(analyzeTickChan)
	if global.TradeSwitch.Sell {
		sellChan = make(chan *streamtick.StreamTick)
		go tradebot.SellBot(sellChan, cond, &input.Close)
	}
	var shouldSell bool
	for {
		tick := <-ch
		if !tradeSwitch {
			tradeSwitch = MissingTicksStatus.CheckByStockNum(tick.StockNum)
		}
		if openChangeRatio == 0 && tick.Open != 0 {
			openChangeRatio = utils.Round((tick.Open - lastClose), 2)
		}
		tmpArr = append(tmpArr, tick)
		if tradebot.BuyOrderMap.CheckStockExist(tick.StockNum) && tradebot.FilledBuyOrderMap.CheckStockExist(tick.StockNum) {
			if !shouldSell {
				input.Close = []float64{}
				shouldSell = true
			}
			sellChan <- tick
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
			} else if cond.TrimHistoryCloseCount {
				input.Close = input.Close[len(input.Close)-int(cond.HistoryCloseCount):]
			}
			rsi, err := tickanalyze.GenerateRSI(input)
			if err != nil {
				logger.GetLogger().Errorf("GenerateRSI at StreamTickProcess Stock: %s, Err: %s", tick.StockNum, err)
				continue
			}
			closeDiff := utils.Round((unSavedTicks.GetLastClose() - lastSaveLastClose), 2)
			if lastSaveLastClose == 0 {
				closeDiff = 0
			}
			lastSaveLastClose = unSavedTicks.GetLastClose()
			unSavedTicksInOutRatio := utils.Round((100 * (float64(outSum) / float64(outSum+inSum))), 2)
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
			if tradeSwitch {
				analyzeTickChan <- &analyze
			}
			unSavedTicks.ClearAll()
		}
	}
}
