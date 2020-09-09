package candle

import (
	"app/infrastructure"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type candleInfraStruct struct {
	ProductCode string
	Duration    time.Duration
	Time        time.Time
	Open        float64
	Close       float64
	High        float64
	Low         float64
	Volume      float64
}

func NewCandle(productCode string, duration time.Duration, timeDate time.Time, open, close, high, low, volume float64) *candleInfraStruct {
	return &candleInfraStruct{
		ProductCode: productCode,
		Duration:    duration,
		Time:        timeDate,
		Open:        open,
		Close:       close,
		High:        high,
		Low:         low,
		Volume:      volume,
	}
}

// テーブルネームを取得する関数
func GetCandleTableName(productCode string, duration time.Duration) string {
	return fmt.Sprintf("%s_%s", productCode, duration)
}

// テーブルネームを取得するメソッド
func (c *candleInfraStruct) TableName() string {
	return GetCandleTableName(c.ProductCode, c.Duration)
}

// テーブルを空にする
func Truncate() error {
	cmd1 := fmt.Sprintf("TRUNCATE %s", "FX_BTC_JPY_1h0m0s")
	cmd2 := fmt.Sprintf("TRUNCATE %s", "FX_BTC_JPY_1m0s")
	truncate1, err1 := infrastructure.DB.Prepare(cmd1)
	truncate2, err2 := infrastructure.DB.Prepare(cmd2)
	if err1 != nil {
		log.Println(err1)
	}
	if err2 != nil {
		log.Println(err2)
	}
	_, err1 = truncate1.Exec()
	_, err2 = truncate2.Exec()
	if err1 != nil {
		log.Println(err1)
	}
	if err2 != nil {
		log.Println(err2)
	}
	return nil
}

// キャンドル情報を追加する
func (c *candleInfraStruct) Insert() error {
	cmd := fmt.Sprintf("INSERT INTO %s (time, open, close, high, low, volume) VALUES (?, ?, ?, ?, ?, ?)", c.TableName())
	ins, err := infrastructure.DB.Prepare(cmd)
	if err != nil {
		log.Println(err)
	}
	// jst, _ := time.LoadLocation("Asia/Tokyo")
	_, err = ins.Exec(c.Time, c.Open, c.Close, c.High, c.Low, c.Volume)
	if err != nil {
		log.Println(err)
	}
	return nil
}

// キャンドル情報を更新する
func (c *candleInfraStruct) Save() error {
	cmd := fmt.Sprintf("UPDATE %s SET open = ?, close = ?, high = ?, low = ?, volume = ? WHERE time = ?", c.TableName())
	upd, err := infrastructure.DB.Prepare(cmd)
	// jst, _ := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Println(err)
	}
	upd.Exec(c.Open, c.Close, c.High, c.Low, c.Volume, c.Time)
	return nil
}

// キャンドル情報を取得する
func SelectOne(productCode string, duration time.Duration, dateTime time.Time) *candleInfraStruct {
	tableName := GetCandleTableName(productCode, duration)
	cmd := fmt.Sprintf("SELECT time, open, close, high, low, volume FROM  %s WHERE time = ?", tableName)
	var candle candleInfraStruct
	err := infrastructure.DB.QueryRow(cmd, dateTime).Scan(&candle.Time, &candle.Open, &candle.Close, &candle.High, &candle.Low, &candle.Volume)
	if err != nil {
		return nil
	}
	return NewCandle(productCode, duration, candle.Time, candle.Open, candle.Close, candle.High, candle.Low, candle.Volume)
}

// キャンドル情報を全て取得する
func SelectAll(productCode string, duration time.Duration, limit int) *sql.Rows {
	tableName := GetCandleTableName(productCode, duration)
	cmd := fmt.Sprintf(`SELECT * FROM (
		SELECT time, open, close, high, low, volume FROM %s ORDER BY time DESC LIMIT ?
		) AS candle ORDER BY time ASC`, tableName)
	rows, err := infrastructure.DB.Query(cmd, limit)
	if err != nil {
		return nil
	}
	return rows
}
