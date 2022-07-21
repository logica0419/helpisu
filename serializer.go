package helpisu

import (
	"github.com/bytedance/sonic/decoder"
	"github.com/bytedance/sonic/encoder"
	"github.com/labstack/echo/v4"
)

// sonicSerializer sonicを用いたecho用Jsonシリアライザ
type sonicSerializer struct{}

// NewSonicSerializer sonicを用いたecho用Jsonシリアライザを作成
func NewSonicSerializer() echo.JSONSerializer {
	return &sonicSerializer{}
}

// Serialize シリアライズ
func (s *sonicSerializer) Serialize(c echo.Context, i interface{}, indent string) error {
	enc := encoder.NewStreamEncoder(c.Response())
	return enc.Encode(i)
}

// Deserialize デシリアライズ
func (s *sonicSerializer) Deserialize(c echo.Context, i interface{}) error {
	return decoder.NewStreamDecoder(c.Request().Body).Decode(i)
}
