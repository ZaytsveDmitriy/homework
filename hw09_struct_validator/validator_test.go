package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require" //nolint:depguard
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	BadVilidator struct {
		BadDigits int `validate:"in:200,to,300"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			App{
				"12345",
			},
			nil,
		},
		{
			App{
				"1234",
			},
			ValidationErrors{ValidationError{Field: "Version", Err: ErrStringLenNotEqual}},
		},
		{
			Response{
				200,
				"jig",
			},
			nil,
		},
		{
			Response{
				103,
				"jig",
			},
			ValidationErrors{ValidationError{Field: "Code", Err: ErrIntNotContain}},
		},
		{
			User{
				"123456789012345678901234567890123456",
				"Alex",
				18,
				"worker@Worker.com",
				"stuff",
				[]string{"+7123456789", "+7123456789"},
				nil,
			},
			nil,
			// Place your code here.
		},
		{
			User{
				"13456789012345678901234567890123456",
				"Alex",
				4,
				"worker@Workercom",
				"stuffer",
				[]string{"+712345f678910", "+712345678910"},
				nil,
			},
			ValidationErrors{
				ValidationError{Field: "ID", Err: ErrStringLenNotEqual},
				ValidationError{Field: "Age", Err: ErrIntLess},
				ValidationError{Field: "Email", Err: ErrStringNotMatchRE},
				ValidationError{Field: "Role", Err: ErrStringNotContain},
				ValidationError{Field: "Phones", Err: ErrStringLenNotEqual},
			},
		},
		{
			Token{
				[]byte{12},
				[]byte{23},
				[]byte{12},
			},
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}
