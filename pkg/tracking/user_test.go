package tracking

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"reflect"
	"testing"
	"time"
)

func TestUser(t *testing.T) {
	c := NewUserClient("", "")

	if got := c.User.client; got == nil {
		t.Errorf("User service client is nil")
	}
	if got, want := c.User.path, "/users/v2"; got != want {
		t.Errorf("User service path is %v, want %v", got, want)
	}
}

func TestUser_ListDevices(t *testing.T) {
	client, mux, teardown := setupTestUserClient()
	defer teardown()

	mux.HandleFunc(path.Join(client.User.path, "/devices"), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testBody(t, r, "")

		w.WriteHeader(http.StatusOK)
	})

	err := client.User.ListDevices(context.Background())
	if err != nil {
		t.Errorf("User.ListDevices returned error: %v", err)
	}
}

func TestUser_Login(t *testing.T) {
	client, mux, teardown := setupTestUserClient()
	defer teardown()

	mux.HandleFunc(path.Join(client.User.path, "/login"), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testBody(t, r, `{"email":"email","password":"password"}`+"\n")

		fmt.Fprint(w, `{"userId":"userId","accessToken":"accessToken","expiresIn":86399,"refreshToken":"refreshToken"}`)
	})

	got, err := client.User.Login(context.Background(), "email", "password")
	if err != nil {
		t.Errorf("User.Login returned error: %v", err)
	}
	want := &UserToken{
		UserID:       "userId",
		AccessToken:  "accessToken",
		ExpiresIn:    Time{Time: time.Unix(0, 86399*int64(time.Millisecond))},
		RefreshToken: "refreshToken",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response is %v, want %v", got, want)
	}
}
