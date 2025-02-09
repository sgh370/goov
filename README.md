# GOOV - Go Object Validator

GOOV is a powerful, flexible, and easy-to-use validation library for Go that provides struct-level validation through tags. It offers a wide range of built-in validation rules and supports custom validation logic.

## Features

- Tag-based validation
- Extensive set of built-in validators
- Conditional validation
- Cross-field validation
- Slice and map validation
- Custom validation rules
- Multiple error handling
- DNS validation for emails
- Nested struct validation

## Installation

```bash
go get github.com/sgh370/goov
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/sgh370/goov/validator"
    "github.com/sgh370/goov/validator/rules"
)

type User struct {
    Username    string   `validate:"required,length:3:20"`
    Email       string   `validate:"required,email"`
    Age         int      `validate:"required,min=18"`
    Interests   []string `validate:"required,min=1,dive,required"`
}

func main() {
    // Create a new validator instance
    v := validator.New()

    // Add validation rules
    v.AddRule("required", rules.Required{})
    v.AddRule("length:3:20", &rules.Length{Min: 3, Max: 20})
    v.AddRule("min=18", &rules.Min{Value: 18})
    v.AddRule("email", &rules.EmailDNS{})

    // Create a user to validate
    user := User{
        Username:  "johndoe",
        Email:     "john@example.com",
        Age:       25,
        Interests: []string{"coding", "reading"},
    }

    // Validate the user
    if err := v.Validate(user); err != nil {
        fmt.Printf("Validation failed: %v\n", err)
        return
    }
    fmt.Println("User is valid!")
}
```

## Available Validation Rules

### Common Rules
- `required` - Ensures a value is not empty
- `length:min:max` - Validates string/slice/map length
- `min=value` - Validates minimum numeric value
- `max=value` - Validates maximum numeric value
- `oneof=value1 value2` - Ensures value is one of the specified options
- `email` - Validates email format with DNS check
- `url` - Validates URL format
- `ip` - Validates IP address format
- `uuid` - Validates UUID format
- `json` - Validates JSON string format

### Numeric Rules
- `positive` - Ensures number is positive
- `range=min:max` - Validates number within range

### Collection Rules
- `unique` - Ensures slice elements are unique
- `contains=value` - Checks if slice/array contains value
- `each=rule` - Applies rule to each element
- `dive` - Validates nested slice elements

### Conditional Rules
- `required_if=field value` - Conditional required validation
- `if=field then=rule` - Conditional rule application
- `unless=field then=rule` - Inverse conditional validation

### Advanced Rules
- `password` - Validates password complexity
- `creditcard` - Validates credit card numbers
- `phone` - Validates phone numbers
- `semver` - Validates semantic version strings
- `domain` - Validates domain names
- `port` - Validates port numbers

## Custom Validation Rules

You can create custom validation rules by implementing the `Rule` interface:

```go
type MyRule struct{}

func (r MyRule) Validate(value interface{}) error {
    // Your validation logic here
    return nil
}

// Add your custom rule
v.AddRule("myrule", MyRule{})
```

## Cross-Field Validation

GOOV supports validation based on other field values:

```go
type Form struct {
    PaymentType string `validate:"required,oneof=card bank"`
    CardNumber  string `validate:"required_if=PaymentType card"`
    BankAccount string `validate:"required_if=PaymentType bank"`
}
```

## Error Handling

GOOV provides detailed error messages for validation failures. You can use `ValidateAll` to get all validation errors at once:

```go
errors := v.ValidateAll(user)
for _, err := range errors {
    fmt.Printf("Validation error: %v\n", err)
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
