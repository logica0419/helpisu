package helpisu

import (
	"context"
	"database/sql"
	"log"
	"time"
)

const dbStartUpWaitSec = 2

// WaitDBStartUp DBの起動を待機
func WaitDBStartUp(ctx context.Context, db *sql.DB) {
	for {
		err := db.PingContext(ctx)
		if err == nil {
			break
		}

		time.Sleep(dbStartUpWaitSec * time.Second)
	}
}

// DBDisconnectDetector DBから切断されるとアプリを強制終了する検出器
type DBDisconnectDetector struct {
	db []*sql.DB
	t  *Ticker
	r  time.Duration
	s  chan struct{}
	st bool
}

const msInSec = 1000

/*
NewDBDisconnectDetector 新たなDBDisconnectDetectorを作成

	durationSecは接続確認の実行間隔をs単位で指定して下さい
	pauseSecは`Pause()`してから検出を再開するまでの時間をs単位で指定して下さい
*/
func NewDBDisconnectDetector(ctx context.Context, durationSec, pauseSec int) *DBDisconnectDetector {
	det := DBDisconnectDetector{
		db: make([]*sql.DB, 0),
		t:  nil,
		r:  time.Second * time.Duration(pauseSec),
		s:  make(chan struct{}),
		st: false,
	}

	// nolint:gomnd
	det.t = NewTicker(durationSec*msInSec, func() {
		for _, db := range det.db {
			err := db.PingContext(ctx)
			if err != nil {
				log.Panic("DB disconnected")
			}
		}
	})

	return &det
}

// RegisterDB DBをDBDisconnectDetectorに登録
func (d *DBDisconnectDetector) RegisterDB(db *sql.DB) {
	d.db = append(d.db, db)
}

/*
Start DBからの切断の検出を開始

	必ずGoroutineとして実行して下さい
*/
func (d *DBDisconnectDetector) Start() {
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
