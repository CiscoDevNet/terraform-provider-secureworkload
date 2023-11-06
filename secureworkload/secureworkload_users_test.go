// +build all integrationtests

package secureworkload

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var (
	usersAPIURL              = os.Getenv("SECUREWORKLOAD_API_URL")
	usersAPIKey              = os.Getenv("SECUREWORKLOAD_API_KEY")
	usersAPISecret           = os.Getenv("SECUREWORKLOAD_API_SECRET")
	usersDefaultAppScopeId   = os.Getenv("SECUREWORKLOAD_APP_SCOPE_ID")
	usersDefaultClientConfig = Config{
		APIKey:                 usersAPIKey,
		APISecret:              usersAPISecret,
		APIURL:                 usersAPIURL,
		DisableTLSVerification: false,
	}
)

func TestListUsers(t *testing.T) {
	client, err := New(usersDefaultClientConfig)
	if err != nil {
		t.Errorf("Error %s creating client with config %+v", err, usersDefaultClientConfig)
	}
	_, err = client.ListUsers(ListUsersRequest{})
	if err != nil {
		t.Errorf("Error %s listing users with client %+v", err, client)
	}
}

func TestDeleteUserDeletesCreatedUser(t *testing.T) {
	client, err := New(usersDefaultClientConfig)
	if err != nil {
		t.Errorf("Error %s creating client with config %+v", err, usersDefaultClientConfig)
	}
	createUserParams := CreateUserRequest{
		Email:      fmt.Sprintf("test+%d@example.com", time.Now().UnixNano()),
		FirstName:  "Levi",
		LastName:   "Schoen",
		AppScopeId: usersDefaultAppScopeId,
	}
	createdUser, err := client.CreateUser(createUserParams)
	if err != nil {
		t.Error(err)
	}
	err = client.DeleteUser(createdUser.Id)
	if err != nil {
		t.Error(err)
	}
	users, err := client.ListUsers(ListUsersRequest{})
	if err != nil {
		t.Errorf("Error %s listing users with client %+v", err, client)
	}
	for _, user := range users {
		if user.Id == createdUser.Id {
			t.Errorf("User %+v should be deleted but wasn't", user)
		}
	}
}

func TestDescribeUserDescribesCreatedUser(t *testing.T) {
	client, err := New(defaultClientConfig)
	if err != nil {
		t.Errorf("Error %s creating client with config %+v", err, defaultClientConfig)
	}
	createUserParams := CreateUserRequest{
		Email:      fmt.Sprintf("test+%d@example.com", time.Now().UnixNano()),
		FirstName:  "Levi",
		LastName:   "Schoen",
		AppScopeId: usersDefaultAppScopeId,
	}
	createdUser, err := client.CreateUser(createUserParams)
	if err != nil {
		t.Error(err)
	}
	describedUser, err := client.DescribeUser(createdUser.Id)
	if describedUser.Email != createUserParams.Email {
		t.Errorf("Expected described user email to be %s, got %s",
			createUserParams.Email, describedUser.Email)
	}
	if err != nil {
		t.Error(err)
	}
	err = client.DeleteUser(createdUser.Id)
	if err != nil {
		t.Error(err)
	}
	users, err := client.ListUsers(ListUsersRequest{})
	if err != nil {
		t.Errorf("Error %s listing users with client %+v", err, client)
	}
	for _, user := range users {
		if user.Id == createdUser.Id {
			t.Errorf("User %+v should be deleted but wasn't", user)
		}
	}
}

func TestEnableUserEnablesDisabledUser(t *testing.T) {
	client, err := New(usersDefaultClientConfig)
	if err != nil {
		t.Errorf("Error %s creating client with config %+v", err, usersDefaultClientConfig)
	}
	createUserParams := CreateUserRequest{
		Email:      fmt.Sprintf("test+%d@example.com", time.Now().UnixNano()),
		FirstName:  "Levi",
		LastName:   "Schoen",
		AppScopeId: usersDefaultAppScopeId,
	}
	createdUser, err := client.CreateUser(createUserParams)
	if err != nil {
		t.Error(err)
	}
	err = client.DeleteUser(createdUser.Id)
	if err != nil {
		t.Error(err)
	}
	users, err := client.ListUsers(ListUsersRequest{})
	if err != nil {
		t.Errorf("Error %s listing users with client %+v", err, client)
	}
	for _, user := range users {
		if user.Id == createdUser.Id {
			t.Errorf("User %+v should be deleted but wasn't", user)
		}
	}
	reenabledUser, err := client.EnableUser(createdUser.Id)
	if err != nil {
		t.Errorf("Error %s re-enabling user %+v with client %+v", err, createdUser, client)
	}
	users, err = client.ListUsers(ListUsersRequest{})
	if err != nil {
		t.Errorf("Error %s re listing users with client %+v", err, client)
	}
	var userReenabled bool
	for _, user := range users {
		if user.Id == reenabledUser.Id {
			userReenabled = true
			break
		}
	}
	if !userReenabled {
		t.Errorf("Expected user %+v\n to be in list of enabled users %+v\n", reenabledUser, users)
	}
	err = client.DeleteUser(reenabledUser.Id)
	if err != nil {
		t.Error(err)
	}
	users, err = client.ListUsers(ListUsersRequest{})
	if err != nil {
		t.Errorf("Error %s listing users with client %+v", err, client)
	}
	for _, user := range users {
		if user.Id == reenabledUser.Id {
			t.Errorf("User %+v should be deleted but wasn't", user)
		}
	}
}
