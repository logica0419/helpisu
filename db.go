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
	r time.Duration
	s chan struct{}
}

/*
NewDBDisconnectDetector 新たなDBDisconnectDetectorを作成
	durationSecは接続確認の実行間隔をs単位で指定して下さい
	pauseSecは`Pause()`してから検出を再開するまでの時間をs単位で指定して下さい
*/
func NewDBDisconnectDetector(durationSec, pauseSec int, db *sql.DB) *DBDisconnectDetector {
	return &DBDisconnectDetector{
		t: NewTicker(durationSec*1000, func() {
			log.Println("check DB connection")
			err := db.Ping()
			if err != nil {
				log.Panic("DB disconnected")
			}
		}),
		r: time.Second * time.Duration(pauseSec),
		s: make(chan struct{}),
	}
}

/*
Start DBからの切断の検出を開始
	必ずGoroutineとして実行して下さい
*/
func (d *DBDisconnectDetector) Start() {
	for {
		select {
		case <-d.s:
			return
		default:
			d.t.Start()
			time.Sleep(d.r)
		}
	}
}

/*
Pause DBからの切断の検出を一時的に停止
	検出はpauseSec秒後に再開します
*/
func (d *DBDisconnectDetector) Pause() {
	d.t.Stop()
}

// Stop DBからの切断の検出を完全に停止
func (d *DBDisconnectDetector) Stop() {
	d.t.Stop()
	d.s <- struct{}{}
}

// Reset 確認タイミングをリセット
func (d *DBDisconnectDetector) Reset() {
	d.t.Reset()
}
