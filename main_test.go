package main

import "testing"

import "github.com/gorilla/sessions"
import "fmt"
import "github.com/markbates/goth/providers/twitter"

func TestSession(t *testing.T) {

	bytes := "MTQ4MTA1NDEzNXxEdi1CQkFFQ180SUFBUkFCRUFBQUtmLUNBQUVHYzNSeWFXNW5EQkVBRDE5bmIzUm9hV05mYzJWemMybHZiZ1p6ZEhKcGJtY01BZ0FBfIaSMyMO0UNKQQtI0WroAIwvxOy0d9HP3wCX8xahO1uH"

	ses := sessions.Session{}
	provider := twitter.Provider{}
	s, _ := provider.UnmarshalSession(bytes)
	fmt.Printf("%+v")
}
