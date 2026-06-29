package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInt(t *testing.T) {
	var stuff struct {
		I int `env:"APP_INTERVAL" default:"100"`
	}
	err := Apply(&stuff)
	assert.NoError(t, err)
	assert.Equal(t, stuff.I, 100)
}

func TestInt64(t *testing.T) {
	var stuff struct {
		I int64 `env:"APP_INTERVAL" default:"100"`
	}
	err := Apply(&stuff)
	assert.NoError(t, err)
	assert.Equal(t, int(stuff.I), 100)
}

func TestInt32(t *testing.T) {
	var stuff struct {
		I int32 `env:"APP_INTERVAL" default:"100"`
	}
	err := Apply(&stuff)
	assert.NoError(t, err)
	assert.Equal(t, int(stuff.I), 100)
}

func TestInt16(t *testing.T) {
	var stuff struct {
		I int16 `env:"APP_INTERVAL" default:"100"`
	}
	err := Apply(&stuff)
	assert.NoError(t, err)
	assert.Equal(t, int(stuff.I), 100)
}

func TestInt8(t *testing.T) {
	var stuff struct {
		I int16 `env:"APP_INTERVAL" default:"100"`
	}
	err := Apply(&stuff)
	assert.NoError(t, err)
	assert.Equal(t, int(stuff.I), 100)
}

func TestFloat64(t *testing.T) {
	var stuff struct {
		F float64 `env:"APP_INTERVAL" default:"10.3"`
	}
	err := Apply(&stuff)
	assert.NoError(t, err)
	assert.Equal(t, stuff.F, float64(10.3))
}

func TestFloat32(t *testing.T) {
	var stuff struct {
		F float32 `env:"APP_INTERVAL" default:"10.3"`
	}
	err := Apply(&stuff)
	assert.NoError(t, err)
	assert.Equal(t, stuff.F, float32(10.3))
}

func TestString(t *testing.T) {
	var stuff struct {
		S string `env:"APP_NAME" default:"hello"`
	}
	err := Apply(&stuff)
	assert.NoError(t, err)
	assert.Equal(t, stuff.S, "hello")
}

func TestBoolTrueVariantT(t *testing.T) {
	var stuff struct {
		S bool `env:"APP_ENABLE" default:"t"`
	}
	err := Apply(&stuff)
	assert.NoError(t, err)
	assert.Equal(t, stuff.S, true)
}

func TestBoolTrueVariantTrue(t *testing.T) {
	var stuff struct {
		S bool `env:"APP_ENABLE" default:"tRue"`
	}
	err := Apply(&stuff)
	assert.NoError(t, err)
	assert.Equal(t, stuff.S, true)
}

func TestBoolTrueVariant1(t *testing.T) {
	var stuff struct {
		S bool `env:"APP_ENABLE" default:"1"`
	}
	err := Apply(&stuff)
	assert.NoError(t, err)
	assert.Equal(t, stuff.S, true)
}

func TestBoolFalseVariantF(t *testing.T) {
	var stuff struct {
		S bool `env:"APP_ENABLE" default:"F"`
	}
	err := Apply(&stuff)
	assert.NoError(t, err)
	assert.Equal(t, stuff.S, false)
}

func TestBoolFalseVariantFalse(t *testing.T) {
	var stuff struct {
		S bool `env:"APP_ENABLE" default:"False"`
	}
	err := Apply(&stuff)
	assert.NoError(t, err)
	assert.Equal(t, stuff.S, false)
}

func TestBoolFalseVariant0(t *testing.T) {
	var stuff struct {
		S bool `env:"APP_ENABLE" default:"0"`
	}
	err := Apply(&stuff)
	assert.NoError(t, err)
	assert.Equal(t, stuff.S, false)
}

func TestBoolFalseInvalid(t *testing.T) {
	var stuff struct {
		S bool `env:"APP_ENABLE" default:"boof"`
	}
	err := Apply(&stuff)
	assert.Error(t, err)
}

func TestDurationWeek(t *testing.T) {
	var stuff struct {
		Dur time.Duration `env:"APP_LENGTH" default:"3w"`
	}
	err := Apply(&stuff)
	assert.NoError(t, err)
	assert.Equal(t, stuff.Dur, 3*7*time.Hour*24)
}

func TestDurationDay(t *testing.T) {
	var stuff struct {
		Dur time.Duration `env:"APP_LENGTH" default:"2d"`
	}
	err := Apply(&stuff)
	assert.NoError(t, err)
	assert.Equal(t, stuff.Dur, 2*time.Hour*24)
}

func TestDurationHour(t *testing.T) {
	var stuff struct {
		Dur time.Duration `env:"APP_LENGTH" default:"2h"`
	}
	err := Apply(&stuff)
	assert.NoError(t, err)
	assert.Equal(t, stuff.Dur, 2*time.Hour)
}

func TestDurationMinute(t *testing.T) {
	var stuff struct {
		Dur time.Duration `env:"APP_LENGTH" default:"2m"`
	}
	err := Apply(&stuff)
	assert.NoError(t, err)
	assert.Equal(t, stuff.Dur, 2*time.Minute)
}

func TestDurationSecond(t *testing.T) {
	var stuff struct {
		Dur time.Duration `env:"APP_LENGTH" default:"2s"`
	}
	err := Apply(&stuff)
	assert.NoError(t, err)
	assert.Equal(t, stuff.Dur, 2*time.Second)
}

func TestDurationMillisecond(t *testing.T) {
	var stuff struct {
		Dur time.Duration `env:"APP_LENGTH" default:"2ms"`
	}
	err := Apply(&stuff)
	assert.NoError(t, err)
	assert.Equal(t, stuff.Dur, 2*time.Millisecond)
}

func TestDurationMicrosecond(t *testing.T) {
	var stuff struct {
		Dur time.Duration `env:"APP_LENGTH" default:"2us"`
	}
	err := Apply(&stuff)
	assert.NoError(t, err)
	assert.Equal(t, stuff.Dur, 2*time.Microsecond)
}

func TestDurationNanosecond(t *testing.T) {
	var stuff struct {
		Dur time.Duration `env:"APP_LENGTH" default:"2ns"`
	}
	err := Apply(&stuff)
	assert.NoError(t, err)
	assert.Equal(t, stuff.Dur, 2*time.Nanosecond)
}

func TestNoDefault(t *testing.T) {
	var stuff struct {
		F string `env:"APP_SOMETHING"`
	}
	err := Apply(&stuff, WithErrorOnMissing())
	assert.Error(t, err)
}

func TestNoDefaultHasValue(t *testing.T) {
	t.Setenv("APP_TEST_NAME", "hello")
	var stuff struct {
		TestName string `env:"APP_TEST_NAME"`
	}
	err := Apply(&stuff, WithErrorOnMissing())
	assert.NoError(t, err)
	assert.Equal(t, stuff.TestName, "hello")
}

func TestEncodedByteArray(t *testing.T) {
	var stuff struct {
		RealData []byte `env:"APP_SOME_ENCODED_STRING" default:"c2lsbHkgd2FiYml0" encoded:"true"`
	}
	err := Apply(&stuff)
	assert.NoError(t, err)
	assert.Equal(t, stuff.RealData, []byte("silly wabbit"))
}

func TestNonEncodedByteArray(t *testing.T) {
	var stuff struct {
		RealData []byte `env:"APP_SOME_ENCODED_STRING" default:"c2lsbHkgd2FiYml0"`
	}
	err := Apply(&stuff)
	assert.NoError(t, err)
	assert.Equal(t, stuff.RealData, []byte("c2lsbHkgd2FiYml0"))
}

func TestSet(t *testing.T) {
	var stuff struct {
		RealData Set `env:"APP_MY_ITEMS" default:"abc,def"`
	}
	err := Apply(&stuff)
	assert.NoError(t, err)
	assert.True(t, stuff.RealData.Contains("abc"))
	assert.True(t, stuff.RealData.Contains("def"))
}

func TestEmptySet(t *testing.T) {
	var stuff struct {
		RealData Set `env:"APP_MY_ITEMS"`
	}
	err := Apply(&stuff)
	assert.NoError(t, err)
	assert.False(t, stuff.RealData.Contains("abc"))
}
