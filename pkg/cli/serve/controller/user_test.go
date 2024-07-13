package serve

import (
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/controller/dto"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/utils"
	"fmt"
	"testing"
)

type TestCase[T any] struct {
	name  string
	valid bool
	value T
}

func NewUser() dto.User {
	return dto.User{
		Email:       "mahdieh@gmail.com",
		FirstName:   "M",
		LastName:    "MM",
		Age:         20,
		Gender:      "women",
		Password:    "123456789",
		PhoneNumber: "09177585872",
		InfoType:    "trainer",
	}
}
func NewUser2() dto.User {
	return dto.User{
		Email:       "",
		FirstName:   "",
		LastName:    "MM",
		Age:         20,
		Password:    "",
		PhoneNumber: "123",
		InfoType:    "trainer",
	}
}

func TestValidation(t *testing.T) {
	testCases := []TestCase[dto.User]{
		{"ValidUser", true, NewUser()},
		{"InvalidUser", false, NewUser2()},
	}

	for i := range testCases {
		t.Run(testCases[i].name, func(t *testing.T) {
			err := utils.ValidateUser(&testCases[i].value)
			fmt.Println(testCases[i].valid, err)
			if err == nil && !testCases[i].valid {
				t.Errorf("Validation passed for invalid user %v", testCases[i].name)
			} else if err != nil && testCases[i].valid {
				t.Errorf("Validation failed for valid user %v: %v", testCases[i].name, err)
			}
		})
	}
}
