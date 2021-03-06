// Package streamtick package streamtick
package streamtick

import (
	"gorm.io/gorm"
)

// StreamTick StreamTick
type StreamTick struct {
	gorm.Model      `json:"-" swaggerignore:"true"`
	StockNum        string  `gorm:"column:stock_num;index:idx_streamtick"`
	TimeStamp       int64   `gorm:"column:timestamp;index:idx_streamtick"`
	Open            float64 `gorm:"column:open"`
	AvgPrice        float64 `gorm:"column:avg_price"`
	Close           float64 `gorm:"column:close"`
	High            float64 `gorm:"column:high"`
	Low             float64 `gorm:"column:low"`
	Amount          float64 `gorm:"column:amount"`
	AmountSum       float64 `gorm:"column:amount_sum"`
	Volume          int64   `gorm:"column:volume"`
	VolumeSum       int64   `gorm:"column:volume_sum"`
	TickType        int64   `gorm:"column:tick_type"`
	ChgType         int64   `gorm:"column:chg_type"`
	PriceChg        float64 `gorm:"column:price_chg"`
	PctChg          float64 `gorm:"column:pct_chg"`
	BidSideTotalVol int64   `gorm:"column:bid_side_total_vol"`
	AskSideTotalVol int64   `gorm:"column:ask_side_total_vol"`
	BidSideTotalCnt int64   `gorm:"column:bid_side_total_cnt"`
	AskSideTotalCnt int64   `gorm:"column:ask_side_total_cnt"`
	Suspend         int64   `gorm:"column:suspend"`
	Simtrade        int64   `gorm:"column:simtrade"`
}

// Tabler Tabler
type Tabler interface {
	TableName() string
}

// TableName TableName
func (StreamTick) TableName() string {
	return "tick_stream"
}

// PtrArr PtrArr
type PtrArr []*StreamTick

// GetOutInRatio GetOutInRatio
func (c *PtrArr) GetOutInRatio() float64 {
	var outSum, inSum int64
	for _, v := range *c {
		switch v.TickType {
		case 0:
			continue
		case 1:
			outSum += v.Volume
		case 2:
			inSum += v.Volume
		}
	}
	return float64(outSum) / float64(outSum+inSum)
}

// GetOutSum GetOutSum
func (c *PtrArr) GetOutSum() int64 {
	var outSum int64
	for _, v := range *c {
		switch v.TickType {
		case 0:
			continue
		case 1:
			outSum += v.Volume
		}
	}
	return outSum
}

// GetInSum GetInSum
func (c *PtrArr) GetInSum() int64 {
	var inSum int64
	for _, v := range *c {
		switch v.TickType {
		case 0:
			continue
		case 2:
			inSum += v.Volume
		}
	}
	return inSum
}

// GetLastClose GetLastClose
func (c *PtrArr) GetLastClose() float64 {
	var tmp []*StreamTick = *c
	return tmp[len(tmp)-1].Close
}

// GetAllCloseArr GetAllCloseArr
func (c *PtrArr) GetAllCloseArr() []float64 {
	var tmp []*StreamTick = *c
	var closeArr []float64
	for _, v := range tmp {
		closeArr = append(closeArr, v.Close)
	}
	return closeArr
}

// GetTotalTime GetTotalTime
func (c *PtrArr) GetTotalTime() float64 {
	tmp := *c
	return float64(tmp[len(tmp)-1].TimeStamp-tmp[0].TimeStamp) / 1000 / 1000 / 1000
}

// PtrArrArr PtrArrArr
type PtrArrArr []PtrArr

// Get Get
func (c *PtrArrArr) Get() []PtrArr {
	return *c
}

// GetCount GetCount
func (c *PtrArrArr) GetCount() int {
	return len(*c)
}

// GetLastNRow GetLastNRow
func (c *PtrArrArr) GetLastNRow(n int) []PtrArr {
	var tmp []PtrArr = *c
	var ans []PtrArr
	for i := len(*c) - 1; i > len(*c)-1-n; i-- {
		ans = append(ans, tmp[i])
	}
	return ans
}

// Append Append
func (c *PtrArrArr) Append(data PtrArr) {
	*c = append(*c, data)
}

// ClearAll ClearAll
func (c *PtrArrArr) ClearAll() {
	*c = []PtrArr{}
}

// GetCloseDiff GetCloseDiff
func (c *PtrArrArr) GetCloseDiff() float64 {
	var tmp []PtrArr = *c
	first := tmp[0][0].Close
	last := tmp[len(tmp)-1][len(tmp[len(tmp)-1])-1].Close
	return last - first
}

// GetLastClose GetLastClose
func (c *PtrArrArr) GetLastClose() float64 {
	var tmp []PtrArr = *c
	last := tmp[len(tmp)-1][len(tmp[len(tmp)-1])-1].Close
	return last
}

// GetLastTick GetLastTick
func (c *PtrArrArr) GetLastTick() *StreamTick {
	var tmp []PtrArr = *c
	last := tmp[len(tmp)-1][len(tmp[len(tmp)-1])-1]
	return last
}
