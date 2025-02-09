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

## Pre-registered Validation Tags

GOOV comes with several pre-registered validation rules that are ready to use without calling `AddRule`. These are common validation rules that are automatically registered when you create a new validator instance using `validator.New()`.

### Basic Pre-registered Rules

```go
v := validator.New() // These rules are automatically registered

type User struct {
    // All these tags work without explicit AddRule calls
    ID          int     `validate:"required"`
    Email       string  `validate:"email"`
    URL         string  `validate:"url"`
    IP          string  `validate:"ip"`
    UUID        string  `validate:"uuid"`
    JSON        string  `validate:"json"`
    Phone       string  `validate:"phone"`
}
```

### Pre-registered Rules List

1. **Data Presence**
   - `required` - Field must not be empty

2. **String Formats**
   - `email` - Validates email format
   - `url` - Validates URL format
   - `ip` - Validates IPv4 or IPv6 address
   - `uuid` - Validates UUID format
   - `json` - Validates JSON string format
   - `phone` - Validates phone number format

3. **Numeric Validations**
   - `positive` - Number must be positive
   - `negative` - Number must be negative
   - `min` - Minimum value check
   - `max` - Maximum value check

4. **Collection Validations**
   - `unique` - All elements must be unique
   - `min` - Minimum number of elements
   - `max` - Maximum number of elements
   - `dive` - Validates elements of a slice/array

### Using Pre-registered Rules with Parameters

Some pre-registered rules accept parameters. Here's how to use them:

```go
type Product struct {
    // Numeric parameters
    Age         int     `validate:"min=18"`
    Score       float64 `validate:"range=0:100"`
    
    // String length
    Name        string  `validate:"length=10"`      // Exact length
    Description string  `validate:"length:5:100"`   // Min:Max length
    
    // Collection size
    Tags        []string `validate:"min=1,max=5"`   // Between 1 and 5 elements
}
```

### Combining Pre-registered Rules

You can combine multiple pre-registered rules in a single tag:

```go
type User struct {
    Email    string   `validate:"required,email"`           // Must be present and valid email
    Tags     []string `validate:"required,min=1,unique"`    // Required, non-empty, unique elements
    Age      int      `validate:"required,min=0,max=150"`   // Required, between 0 and 150
}
```

### Important Notes

1. Pre-registered rules are designed for common validation scenarios
2. They are optimized for performance and memory usage
3. You can override pre-registered rules by calling `AddRule` with the same name
4. Custom validation rules still need to be registered using `AddRule`

## Understanding AddRule and Tags

### The Relationship Between AddRule and Tags

The validation system in GOOV works through two complementary components:

1. **AddRule**: Registers validation rules with the validator instance
2. **Tags**: Applies those registered rules to struct fields

Before a tag can be used in a struct, its corresponding rule must be registered using `AddRule`. Here's how it works:

```go
// 1. First, register the rule
v.AddRule("length:3:20", &rules.Length{Min: 3, Max: 20})

// 2. Then use it in struct tags
type User struct {
    Username string `validate:"length:3:20"`
}
```

### When to Use AddRule

You need to use `AddRule` in these scenarios:
- When initializing your validator instance
- When defining parameterized rules (e.g., `length:3:20`, `min=18`)
- When adding custom validation rules
- When setting up conditional validations

Example of registering different types of rules:
```go
v := validator.New()

// Simple rules
v.AddRule("required", rules.Required{})

// Parameterized rules
v.AddRule("length:3:20", &rules.Length{Min: 3, Max: 20})
v.AddRule("min=18", &rules.Min{Value: 18})

// Conditional rules
v.AddRule("required_if", &rules.If{
    Field: "ContactMethod",
    Then:  rules.Required{},
})
```

### When to Use Tags

Tags are used in your struct definitions to:
- Apply registered validation rules to specific fields
- Combine multiple validation rules
- Define the validation flow for your data model

Example of using tags:
```go
type User struct {
    // Single rule
    IsActive  bool    `validate:"required"`
    
    // Multiple rules
    Username  string  `validate:"required,length:3:20"`
    
    // Conditional validation
    Phone     string  `validate:"required_if=ContactType phone"`
    
    // Nested validation
    Addresses []Address `validate:"required,dive"`
}
```

Remember: Any validation rule used in a tag must first be registered using `AddRule`. If you use a tag that hasn't been registered, the validator will return an error.

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
