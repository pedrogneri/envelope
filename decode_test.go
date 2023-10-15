package envelope

import (
	"os"
	"reflect"
	"testing"
)

type Environment struct {
	DatabaseAddr string `envelope:"DATABASE_ADDR,required"`
	DatabasePort int    `envelope:"DATABASE_PORT,required"`
}

type EnvironmentWithDefaults struct {
	DatabasePort int `envelope:"DATABASE_PORT,default:3000"`
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

func assertEqual(t *testing.T, expected interface{}, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v, received %v", expected, actual)
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("unexpected error %s", err.Error())
	}
}

func assertEqualError(t *testing.T, expectedErrMessage string, err error) {
	if err.Error() != expectedErrMessage {
		t.Errorf(`expected "%s", receive "%s"`, expectedErrMessage, err.Error())
	}
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

		assertEqual(t, expectedEnv, env)
		assertNoError(t, err)
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

		assertEqual(t, expectedEnv, env)
		assertNoError(t, err)
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

		assertEqual(t, expectedEnv, env)
		assertNoError(t, err)
	})

	t.Run("should return error if a required env was not found", func(t *testing.T) {
		defer clear()

		os.Setenv("DATABASE_ADDR", "localhost")

		env := new(NestedStruct)
		err := Decode(env)

		expectedErrorMsg := `missing a required field "DATABASE_PORT"`

		assertEqualError(t, expectedErrorMsg, err)
	})

	t.Run("should return error if failed on convert env to struct type", func(t *testing.T) {
		defer clear()

		os.Setenv("DATABASE_ADDR", "localhost")
		os.Setenv("DATABASE_PORT", "invalid port")

		env := new(Environment)
		err := Decode(env)

		expectedErrorMsg := `error converting value from "DATABASE_PORT" field into int`

		assertEqualError(t, expectedErrorMsg, err)
	})

	t.Run("should return aggregated errors", func(t *testing.T) {
		defer clear()

		os.Setenv("DATABASE_PORT", "invalid port")

		env := new(Environment)
		err := Decode(env)

		expectedErrorMsg := `missing a required field "DATABASE_ADDR"; error converting value from "DATABASE_PORT" field into int`

		assertEqualError(t, expectedErrorMsg, err)
	})

	t.Run("should use default value if env not set", func(t *testing.T) {
		defer clear()

		env := new(EnvironmentWithDefaults)
		err := Decode(env)

		expectedEnv := &EnvironmentWithDefaults{
			DatabasePort: 3000,
		}

		assertEqual(t, expectedEnv, env)
		assertNoError(t, err)
	})
}
