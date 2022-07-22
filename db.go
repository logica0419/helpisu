package helpisu

import (
	"database/sql"
	"log"
	"time"
)

// WaitDBStartUp DBの起動を待機
func WaitDBStartUp(db *sql.DB) {
	for {
		err := db.Ping()
		if err == nil {
			break
		}

		time.Sleep(time.Second * 2)
	}
}

// DBDisconnectDetector DBから切断されるとアプリを強制終了する検出器
type DBDisconnectDetector struct {
	t *Ticker
}

/*
NewDBDisconnectDetector 新たなDBDisconnectDetectorを作成
	durationSecは接続確認の実行間隔をs単位で指定して下さい
*/
func NewDBDisconnectDetector(durationSec int, db *sql.DB) *DBDisconnectDetector {
	return &DBDisconnectDetector{
		t: NewTicker(durationSec*1000, func() {
			err := db.Ping()
			if err != nil {
				log.Panic("DB disconnected")
			}
		}),
	}
}

/*
Start DBからの切断の検出を開始
	必ずGoroutineとして実行して下さい
*/
func (d *DBDisconnectDetector) Start() {
	d.t.Start()
}

// Stop DBからの切断の検出を停止
func (d *DBDisconnectDetector) Stop() {
	d.t.Stop()
}

// Reset 確認タイミングをリセット
func (d *DBDisconnectDetector) Reset() {
	d.t.Reset()
}
