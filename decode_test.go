package envelope

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Environment struct {
	DatabaseAddr string `envelope:"DATABASE_ADDR,required"`
	DatabasePort int    `envelope:"DATABASE_PORT,required"`
}

type EmbeddedStruct struct {
	Environment
}

type NestedStruct struct {
	Environment Environment
}

func clear() {
	os.Clearenv()
}

func TestDecode(t *testing.T) {
	t.Run("should decode struct", func(t *testing.T) {
		defer clear()

		os.Setenv("DATABASE_ADDR", "localhost")
		os.Setenv("DATABASE_PORT", "5432")

		env := new(Environment)
		err := Decode(env)

		expectedEnv := &Environment{
			DatabaseAddr: "localhost",
			DatabasePort: 5432,
		}

		assert.Equal(t, expectedEnv, env)
		assert.NoError(t, err)
	})

	t.Run("should decode embedded struct", func(t *testing.T) {
		defer clear()

		os.Setenv("DATABASE_ADDR", "localhost")
		os.Setenv("DATABASE_PORT", "5432")

		env := new(EmbeddedStruct)
		err := Decode(env)

		expectedEnv := &EmbeddedStruct{
			Environment: Environment{
				DatabaseAddr: "localhost",
				DatabasePort: 5432,
			},
		}

		assert.Equal(t, expectedEnv, env)
		assert.NoError(t, err)
	})

	t.Run("should decode nested struct", func(t *testing.T) {
		defer clear()

		os.Setenv("DATABASE_ADDR", "localhost")
		os.Setenv("DATABASE_PORT", "5432")

		env := new(NestedStruct)
		err := Decode(env)

		expectedEnv := &NestedStruct{
			Environment: Environment{
				DatabaseAddr: "localhost",
				DatabasePort: 5432,
			},
		}

		assert.Equal(t, expectedEnv, env)
		assert.NoError(t, err)
	})

	t.Run("should return error if a required env was not found", func(t *testing.T) {
		defer clear()

		os.Setenv("DATABASE_ADDR", "localhost")

		env := new(NestedStruct)
		err := Decode(env)

		expectedErrorMsg := `missing a required field "DATABASE_PORT"`

		assert.EqualError(t, err, expectedErrorMsg)
	})

	t.Run("should return error if failed on convert env to struct type", func(t *testing.T) {
		defer clear()

		os.Setenv("DATABASE_ADDR", "localhost")
		os.Setenv("DATABASE_PORT", "invalid port")

		env := new(Environment)
		err := Decode(env)

		expectedErrorMsg := `error converting value from "DATABASE_PORT" field into int`

		assert.EqualError(t, err, expectedErrorMsg)
	})

	t.Run("should return aggregated errors", func(t *testing.T) {
		defer clear()

		os.Setenv("DATABASE_PORT", "invalid port")

		env := new(Environment)
		err := Decode(env)

		expectedErrorMsg := `missing a required field "DATABASE_ADDR"; error converting value from "DATABASE_PORT" field into int`

		assert.EqualError(t, err, expectedErrorMsg)
	})
}
