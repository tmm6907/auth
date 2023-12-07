package interfaces

import "fmt"

type Validator func() (err error)

type Book struct {
	Name string
}

func (b *Book) IsValid() (err error) {
	if len([]rune(b.Name)) > 5 {
		return fmt.Errorf("name too long")
	}
	return nil
}

func IsValidator(validator Validator) {
	fmt.Println("Is Validator!")
}
