# Envelope

Envelope is a lightweight Go package designed to decode environment variables into structs without relying on any external dependencies.

## Usage
Using Envelope in your project is straightforward. Here's an example:

```go
type Environment struct {
    Avocado string `envelope:"AVOCADO,required"`
}

func main() {
    myEnv := new(Environment)

    if err := envelope.Decode(myEnv); err != nil {
        panic(err)
    }
}
```

## Properties
To configure your struct with Envelope, you need to define the envelope tag as follows: ``envelope:"YOUR_ENV_NAME"``

| Property    | Description | Usage   |
| --------    | -------    | ------- |
| Required    | Adding this to the envelope tag validates the variable as required. It will return an error if it is not found.    | ``envelope:"PORT,required"`` |
| Default | Set a default value for an environment variable. If it is not found, this value will be used as a replacement.     | ``envelope:"PORT,default:8080"`` |

Feel free to use Envelope to effortlessly decode your environment variables into your Go structs.
