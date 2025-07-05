package helpisu

import (
	"encoding/json"
	"errors"
	"strconv"
)

/*
StringInt DBにはint型に、jsonにはstring型として認識される特殊な型

	PrimaryKeyをランダムなstringからauto incrementなintに変換する時などに使います
	中身の値を使いたいときは、Value()メソッドを使用して下さい
*/
type StringInt struct { // nolint:recvcheck
	value int
}

// ErrInvalidType 不正な型が渡されたときに返すエラー
var ErrInvalidType = errors.New("StringInt.Scan: invalid type")

// NewStringInt 新たなStringIntを作成
func NewStringInt(value int) StringInt {
	return StringInt{
		value: value,
	}
}

// UnmarshalJSON json.Unmarshalerの実装
func (si *StringInt) UnmarshalJSON(data []byte) error {
	var str string

	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	si.value, err = strconv.Atoi(str)
	if err != nil {
		return err
	}

	return nil
}

// MarshalJSON json.Marshalerの実装
func (si StringInt) MarshalJSON() ([]byte, error) {
	str := strconv.Itoa(si.value)

	return json.Marshal(str)
}

// Scan sql.Scannerの実装
func (si *StringInt) Scan(value any) error {
	var ok bool

	si.value, ok = value.(int)
	if !ok {
		return ErrInvalidType
	}

	return nil
}

// Value sql.Valuerの実装
func (si StringInt) Value() int {
	return si.value
}
