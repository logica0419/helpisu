package helpisu_test

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/logica0419/helpisu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringInt_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	randNum, err := rand.Int(rand.Reader, big.NewInt(2147483647))
	require.NoError(t, err)

	type fields struct {
		value int
	}

	type args struct {
		data []byte
	}

	tests := map[string]struct {
		fields    fields
		args      args
		assertion func(t *testing.T, si helpisu.StringInt, err error)
	}{
		"success": {
			fields: fields{
				value: 0,
			},
			args: args{
				data: []byte(fmt.Sprintf(`{"value":"%d"}`, randNum)),
			},
			assertion: func(t *testing.T, si helpisu.StringInt, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.Equal(t, int(randNum.Int64()), si.Value())
			},
		},
		"invalid type (int)": {
			fields: fields{
				value: 0,
			},
			args: args{
				data: []byte(fmt.Sprintf(`{"value":%d}`, randNum)),
			},
			assertion: func(t *testing.T, _ helpisu.StringInt, err error) {
				t.Helper()
				require.Error(t, err)
			},
		},
	}

	for name, tt := range tests {
		name, tt := name, tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			si := helpisu.NewStringInt(tt.fields.value)
			data := struct {
				Value helpisu.StringInt `json:"value"`
			}{
				Value: si,
			}
			err := json.Unmarshal(tt.args.data, &data)
			tt.assertion(t, data.Value, err)
		})
	}
}

func TestStringInt_MarshalJSON(t *testing.T) {
	t.Parallel()

	randNum, err := rand.Int(rand.Reader, big.NewInt(2147483647))
	require.NoError(t, err)

	type fields struct {
		value int
	}

	tests := map[string]struct {
		fields    fields
		assertion func(t *testing.T, data []byte, err error)
	}{
		"success": {
			fields: fields{
				value: int(randNum.Int64()),
			},
			assertion: func(t *testing.T, data []byte, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.Equal(t, fmt.Sprintf(`{"value":"%d"}`, randNum), string(data))
			},
		},
	}

	for name, tt := range tests {
		name, tt := name, tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			si := helpisu.NewStringInt(tt.fields.value)
			data, err := json.Marshal(struct {
				Value helpisu.StringInt `json:"value"`
			}{
				Value: si,
			})
			tt.assertion(t, data, err)
		})
	}
}

func TestStringInt_Scan(t *testing.T) {
	t.Parallel()

	randNum, err := rand.Int(rand.Reader, big.NewInt(2147483647))
	require.NoError(t, err)

	type fields struct {
		value int
	}

	type args struct {
		value any
	}

	tests := map[string]struct {
		fields    fields
		args      args
		assertion func(t *testing.T, si helpisu.StringInt, err error)
	}{
		"success": {
			fields: fields{
				value: int(randNum.Int64()),
			},
			args: args{
				value: int(randNum.Int64()),
			},
			assertion: func(t *testing.T, si helpisu.StringInt, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.Equal(t, int(randNum.Int64()), si.Value())
			},
		},
		"invalid type": {
			fields: fields{
				value: 0,
			},
			args: args{
				value: "",
			},
			assertion: func(t *testing.T, _ helpisu.StringInt, err error) {
				t.Helper()
				require.Error(t, err)
				assert.Equal(t, err, helpisu.ErrInvalidType)
			},
		},
	}

	for name, tt := range tests {
		name, tt := name, tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			si := helpisu.NewStringInt(tt.fields.value)
			err := si.Scan(tt.args.value)
			tt.assertion(t, si, err)
		})
	}
}

func TestStringInt_Value(t *testing.T) {
	t.Parallel()

	randNum, err := rand.Int(rand.Reader, big.NewInt(2147483647))
	require.NoError(t, err)

	type fields struct {
		value int
	}

	tests := map[string]struct {
		fields    fields
		assertion func(t *testing.T, value int)
	}{
		"success": {
			fields: fields{
				value: int(randNum.Int64()),
			},
			assertion: func(t *testing.T, value int) {
				t.Helper()
				assert.Equal(t, int(randNum.Int64()), value)
			},
		},
	}

	for name, tt := range tests {
		name, tt := name, tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			si := helpisu.NewStringInt(tt.fields.value)
			tt.assertion(t, si.Value())
		})
	}
}
