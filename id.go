package helpisu

import (
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

/*
NewUUID 新たなUUIDを作成

	36文字のstring形式で返します
*/
func NewUUID() string {
	return uuid.Must(uuid.NewRandom()).String()
}

/*
NewULID 新たなULIDを作成

	26文字のstring形式で返します
*/
func NewULID() string {
	return ulid.Make().String()
}
