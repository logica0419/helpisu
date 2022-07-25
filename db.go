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
	db *sql.DB
	d  int
	t  *Ticker
	r  time.Duration
	s  chan struct{}
	st bool
}

/*
NewDBDisconnectDetector 新たなDBDisconnectDetectorを作成
	durationSecは接続確認の実行間隔をs単位で指定して下さい
	pauseSecは`Pause()`してから検出を再開するまでの時間をs単位で指定して下さい
*/
func NewDBDisconnectDetector(durationSec, pauseSec int) *DBDisconnectDetector {
	return &DBDisconnectDetector{
		d:  durationSec,
		r:  time.Second * time.Duration(pauseSec),
		s:  make(chan struct{}),
		st: false,
	}
}

// RegisterDB DBをDBDisconnectDetectorに登録
func (d *DBDisconnectDetector) RegisterDB(db *sql.DB) {
	d.db = db
	d.t = NewTicker(d.d*1000, func() {
		log.Println("DB check")
		err := d.db.Ping()
		if err != nil {
			log.Panic("DB disconnected")
		}
	})
}

/*
Start DBからの切断の検出を開始
	DBが登録されてない場合はPanicします
	必ずGoroutineとして実行して下さい
*/
func (d *DBDisconnectDetector) Start() {
	if d.db == nil {
		log.Panic("DB not registered")
	}

	d.st = true
	d.t.Start()

	for {
		select {
		case <-d.s:
			d.st = false
			return
		default:
			time.Sleep(d.r)
			d.t.Start()
		}
	}
}

/*
Pause DBからの切断の検出を一時的に停止
	検出は`pauseSec`秒後に再開します
*/
func (d *DBDisconnectDetector) Pause() {
	if !d.st {
		return
	}

	d.t.Stop()
}

// Stop DBからの切断の検出を完全に停止
func (d *DBDisconnectDetector) Stop() {
	if !d.st {
		return
	}

	d.t.Stop()
	d.s <- struct{}{}
}

// Reset 確認タイミングをリセット
func (d *DBDisconnectDetector) Reset() {
	if !d.st {
		return
	}

	d.t.Reset()
}
