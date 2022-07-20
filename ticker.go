package helpisu

import "time"

// Ticker 一定時間ごとに関数を実行するタイマー
type Ticker struct {
	d time.Duration
	t *time.Ticker
	f func()
	s chan struct{}
}

/*
NewTicker 新たなTickerを作成
	durationSecはタイマーの実行間隔を秒で指定して下さい
*/
func NewTicker(durationSec int, callback func()) *Ticker {
	return &Ticker{
		d: time.Duration(durationSec) * time.Second,
		f: callback,
		s: make(chan struct{}),
	}
}

/*
Start タイマーを開始
	必ずGoroutineとして実行して下さい
*/
func (t *Ticker) Start() {
	t.t = time.NewTicker(t.d)
	defer t.t.Stop()
	defer close(t.s)

	for {
		select {
		case <-t.t.C:
			go t.f()
		case <-t.s:
			return
		}
	}
}

// Stop タイマーを停止
func (t *Ticker) Stop() {
	t.s <- struct{}{}
}

// Reset タイマーをリセット
func (t *Ticker) Reset() {
	if t.t != nil {
		t.t.Reset(t.d)
	}
}
