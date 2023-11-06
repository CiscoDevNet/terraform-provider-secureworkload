// +build all integrationtests

package secureworkload

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	rolesAPIURL                = os.Getenv("SECUREWORKLOAD_API_URL")
	rolesAPIKey                = os.Getenv("SECUREWORKLOAD_API_KEY")
	rolesAPISecret             = os.Getenv("SECUREWORKLOAD_API_SECRET")
	rolesRootScopeAppId        = os.Getenv("SECUREWORKLOAD_APP_SCOPE_ID")
	rolesDefaultRootScopeAppId = os.Getenv("SECUREWORKLOAD_ROOT_SCOPE_APP_ID")
	rolesDefaultClientConfig   = Config{
		APIKey:                 rolesAPIKey,
		APISecret:              rolesAPISecret,
		APIURL:                 rolesAPIURL,
		DisableTLSVerification: false,
	}
	createApplicationParams = CreateApplicationRequest{
		AppScopeId:         rolesDefaultRootScopeAppId,
		Name:               fmt.Sprintf("test+%d", time.Now().UnixNano()),
		Description:        "TestAddingRoleToScopeAndUser",
		AlternateQueryMode: true,
		StrictValidation:   true,
		Primary:            true,
		Clusters: []Cluster{
			Cluster{
				Id:          "ClusterA",
				Name:        "ClusterA",
				Description: "A Cluster.",
				Nodes: []Node{
					{
						Name:      "ClusterA Node1",
						IPAddress: "10.0.0.1",
					},
				},
				ConsistentUUID: "ClusterA",
			},
		},
		Filters: []PolicyFilter{
			PolicyFilter{
				Id:   "FilterA",
				Name: "DisplayedClusterName",
				Query: []byte(`{
                      "type": "eq",
                      "field": "ip",
                      "value": "10.0.0.1"
                        }`),
			},
		},
		AbsolutePolicies: []Policy{
			Policy{
				ConsumerFilterId: rolesDefaultRootScopeAppId,
				ProviderFilterId: rolesDefaultRootScopeAppId,
				Action:           "ALLOW",
				Layer4NetworkPolicies: []Layer4NetworkPolicy{Layer4NetworkPolicy{
					PortRange: [2]int{80, 80}},
				},
			},
		},
		DefaultPolicies: []Policy{
			Policy{
				ConsumerFilterId: rolesDefaultRootScopeAppId,
				ProviderFilterId: rolesDefaultRootScopeAppId,
				Action:           "DENY",
				Layer4NetworkPolicies: []Layer4NetworkPolicy{{
					PortRange: [2]int{8080, 8080}},
				},
			},
		},
		CatchAllAction: "DENY",
	}
)

func deleteUsers(usersToDelete []User, client Client, t *testing.T) {
	for _, user := range usersToDelete {
		err := client.DeleteUser(user.Id)
		if err != nil {
			t.Errorf("Error deleting User: %s", err)
		}
	}
}
func deleteRoles(rolesToDelete []Role, client Client, t *testing.T) {
	for _, role := range rolesToDelete {
		err := client.DeleteRole(role.Id)
		if err != nil {
			t.Errorf("Error %s deleting role %s with config %+v", err, role.Id, rolesDefaultClientConfig)
		}
	}
}

func deleteApplications(applicationsToDelete []Application, client Client, t *testing.T) {
	for _, application := range applicationsToDelete {
		err := client.DeleteApplication(application.Id)
		if err != nil {
			t.Errorf("Error deleting Application: %s", err)
		}
		if !Await(func() bool {
			_, err := client.DescribeApplication(DescribeApplicationRequest{ApplicationId: application.Id})
			// easy way to check error status?
			if err == nil {
				return false
			}
			return true
		}, 5) {
			t.Errorf("Could not delete application %s", application.Id)
		}
	}
}

// This never works, due to how long it takes an application to be removed from a scope.
// It's taken over a minute multiple times
func deleteScopes(scopesToDelete []Scope, client Client, t *testing.T) {
	// for _, scope := range scopesToDelete {
	// 	if !Await(func() bool {
	// 		err := client.DeleteScope(scope.Id)
	// 		if err != nil {
	// 			return false
	// 		}
	// 		return true
	// 	}, 3) {
	// 		t.Logf("Could not delete scope %s", scope.Id)
	// 	}
	// }
}

func TestListRoles(t *testing.T) {
	client, err := New(rolesDefaultClientConfig)
	if err != nil {
		t.Fatalf("Error %s creating client with config %+v", err, rolesDefaultClientConfig)
	}
	_, err = client.ListRoles()
	if err != nil {
		t.Fatalf("Error %s getting roles with client %+v", err, client)
	}
}

func TestCreateAndDescribeRole(t *testing.T) {
	client, err := New(rolesDefaultClientConfig)
	if err != nil {
		t.Fatalf("Error %s creating client with config %+v", err, rolesDefaultClientConfig)
	}

	createRoleParams := CreateRoleRequest{
		Name:        fmt.Sprintf("GoSDKTest %d", time.Now().UnixNano()),
		Description: "TestCreateAndDescribeRole",
		AppScopeId:  rolesRootScopeAppId,
	}

	createdRole, err := client.CreateRole(createRoleParams)
	if err != nil {
		t.Fatalf("Error %s creating role with params %+v and config %+v", err, createRoleParams, rolesDefaultClientConfig)
	}

	fetchedRole, err := client.GetRole(createdRole.Id)
	if err != nil {
		t.Fatalf("Error %s getting role %s with config %+v", err, createdRole.Id, rolesDefaultClientConfig)
	}
	if fetchedRole.Id != createdRole.Id {
		t.Fatalf("Expected fetched role id to be %s, got %s", createdRole.Id, fetchedRole.Id)
	}

	err = client.DeleteRole(createdRole.Id)
	if err != nil {
		t.Fatalf("Error %s deleting role %s with config %+v", err, createdRole.Id, rolesDefaultClientConfig)
	}

	roles, err := client.ListRoles()
	if err != nil {
		t.Fatalf("Error %s listing roles with config %+v", err, rolesDefaultClientConfig)
	}

	for _, role := range roles {
		if role.Id == createdRole.Id {
			t.Errorf("Role %+v should have been deleted but wasn't", role)
		}
	}
}

func TestAddingRoleToScopeAndUser(t *testing.T) {
	rolesToDelete := []Role{}
	usersToDelete := []User{}
	applicationsToDelete := []Application{}
	scopesToDelete := []Scope{}

	client, err := New(rolesDefaultClientConfig)
	if err != nil {
		t.Fatalf("Error %s creating client with config %+v", err, rolesDefaultClientConfig)
	}
	t.Cleanup(func() {
		deleteRoles(rolesToDelete, client, t)
		deleteUsers(usersToDelete, client, t)
		deleteApplications(applicationsToDelete, client, t)
		deleteScopes(scopesToDelete, client, t)
	})

	createScopeParams := CreateScopeRequest{
		ShortName:        fmt.Sprintf("GoSDKTest %d", time.Now().UnixNano()),
		Description:      "TestAddingRoleToScopeAndUser",
		ParentAppScopeId: rolesRootScopeAppId,
		ShortQuery: ShortQuery{
			Type:  "eq",
			Field: "ip",
			Value: "10.0.0.1",
		},
	}
	createdScope, err := client.CreateScope(createScopeParams)
	if err != nil {
		t.Error(err)
	}

	scopesToDelete = append(scopesToDelete, createdScope)

	createRoleParams := CreateRoleRequest{
		Name:        fmt.Sprintf("GoSDKTest %d", time.Now().UnixNano()),
		Description: "TestAddingRoleToScopeAndUser",
		AppScopeId:  rolesRootScopeAppId,
	}

	createdRole, err := client.CreateRole(createRoleParams)
	if err != nil {
		t.Fatalf("Error %s creating role with params %+v and config %+v", err, createRoleParams, rolesDefaultClientConfig)
	}
	rolesToDelete = append(rolesToDelete, createdRole)

	newParams := createApplicationParams
	newParams.AppScopeId = createdScope.Id

	createdApplication, err := client.CreateApplication(newParams)
	if err != nil {
		t.Fatal(err)
	}
	applicationsToDelete = append(applicationsToDelete, createdApplication)

	badScopeAccessParams := GiveScopeAccessToRoleRequest{
		RoleId:     createdRole.Id,
		AppScopeId: createdApplication.AppScopeId,
		Ability:    "MALICIOUS_ACTOR",
	}
	emptyScopeAccessParams := GiveScopeAccessToRoleRequest{
		RoleId:     createdRole.Id,
		AppScopeId: createdApplication.AppScopeId,
		Ability:    "",
	}
	goodScopeAccessParams := GiveScopeAccessToRoleRequest{
		RoleId:     createdRole.Id,
		AppScopeId: createdApplication.AppScopeId,
		Ability:    "developer",
	}

	// verify AccessType validation
	_, err = client.GiveScopeAccessToRole(badScopeAccessParams)
	if err == nil {
		t.Fatal("Expected validation of bad ability to fail, it did not", err, createdApplication.Id, createdRole)
	}
	_, err = client.GiveScopeAccessToRole(emptyScopeAccessParams)
	if err == nil {
		t.Fatal("Expected validation of empty ability to fail, it did not", err, createdApplication.Id, createdRole)
	}

	// happy path
	fetchedRole, err := client.GetRole(createdRole.Id)
	fmt.Println(fmt.Sprintf("Happy Path - Fetched role:\n%+v", fetchedRole))
	_, err = client.GiveScopeAccessToRole(goodScopeAccessParams)
	if err != nil {
		t.Fatalf("Error giving scope %s access to role %+v\n\nCaused by: %s\n", createdApplication.Id, createdRole, err)
	}

	createUserParams := CreateUserRequest{
		Email:      fmt.Sprintf("test+%d@example.com", time.Now().UnixNano()),
		FirstName:  "Eric",
		LastName:   "Hibbs",
		AppScopeId: rolesRootScopeAppId,
	}
	createdUser, err := client.CreateUser(createUserParams)
	if err != nil {
		t.Fatal(err)
	}

	usersToDelete = append(usersToDelete, createdUser)

	// verify addition and removal of role from a user
	var userWithRole, userWithoutRole User

	if !Await(func() bool {
		userWithRole, err = client.AddRoleToUser(AddRoleToUserRequest{
			UserId: createdUser.Id,
			RoleId: createdRole.Id,
		})
		if err != nil {
			t.Fatalf("Error %s adding role %+v to user %+v", err, createdRole, userWithRole)
		}
		return doesUserHaveRole(userWithRole, createdRole)
	}, 3) {
		t.Fatalf("Expected user to have role %s, but user only had [%s]", createdRole.Id, strings.Join(userWithRole.RoleIds, ", "))
	}

	if !Await(func() bool {
		userWithoutRole, err = client.RemoveRoleFromUser(RemoveRoleFromUserRequest{
			RoleId: createdRole.Id,
			UserId: createdUser.Id,
		})
		if err != nil {
			t.Fatalf("Error %s removing role %s from user %+v", err, createdRole.Id, userWithoutRole)
		}
		return !doesUserHaveRole(userWithoutRole, createdRole)
	}, 3) {
		t.Fatalf("Role %s should have been removed from the user but wasn't", createdRole.Id)
	}

}
func TestCreatingScopedRole(t *testing.T) {
	rolesToDelete := []Role{}
	applicationsToDelete := []Application{}
	scopesToDelete := []Scope{}

	client, err := New(rolesDefaultClientConfig)
	if err != nil {
		t.Fatalf("Error %s creating client with config %+v", err, rolesDefaultClientConfig)
	}

	t.Cleanup(func() {
		deleteRoles(rolesToDelete, client, t)
		deleteApplications(applicationsToDelete, client, t)
		deleteScopes(scopesToDelete, client, t)
	})

	createScopeParams := CreateScopeRequest{
		ShortName:        fmt.Sprintf("GoSDKTest %d", time.Now().UnixNano()),
		Description:      "TestAddingRoleToScopeAndUser",
		ParentAppScopeId: rolesRootScopeAppId,
		ShortQuery: ShortQuery{
			Type:  "eq",
			Field: "ip",
			Value: "10.0.0.1",
		},
	}
	createdScope, err := client.CreateScope(createScopeParams)
	if err != nil {
		t.Error(err)
	}

	scopesToDelete = append(scopesToDelete, createdScope)

	newParams := createApplicationParams
	newParams.AppScopeId = createdScope.Id
	createdApplication, err := client.CreateApplication(newParams)
	if err != nil {
		t.Fatal(err)
	}
	applicationsToDelete = append(applicationsToDelete, createdApplication)

	createScopedRoleParams := CreateScopedRoleRequest{
		Name:                fmt.Sprintf("GoSDKTest %d", time.Now().UnixNano()),
		Description:         "TestCreatingScopedRole",
		AppScopeId:          rolesRootScopeAppId,
		AbilitiesAppScopeId: createdApplication.AppScopeId,
		Ability:             "SCOPE_READ",
	}

	createScopedRoleResponse, err := client.CreateScopedRole(createScopedRoleParams)
	if err != nil {
		t.Fatalf("Error %s creating scoped role with params %+v and config %+v", err, createScopedRoleParams, rolesDefaultClientConfig)
	}
	rolesToDelete = append(rolesToDelete, createScopedRoleResponse.Role)

	if createScopedRoleResponse.AppScopeId != createdApplication.AppScopeId {
		t.Fatal("Expected role to be assigned to created application scope id, but it wasn't")
	}
}

func TestCreateScopedRoleForUsers(t *testing.T) {
	rolesToDelete := []Role{}
	usersToDelete := []User{}
	applicationsToDelete := []Application{}
	scopesToDelete := []Scope{}

	client, err := New(rolesDefaultClientConfig)
	if err != nil {
		t.Fatalf("Error %s creating client with config %+v", err, rolesDefaultClientConfig)
	}

	t.Cleanup(func() {
		deleteRoles(rolesToDelete, client, t)
		deleteUsers(usersToDelete, client, t)
		deleteApplications(applicationsToDelete, client, t)
		deleteScopes(scopesToDelete, client, t)
	})

	// scope creation
	createScopeParams := CreateScopeRequest{
		ShortName:        fmt.Sprintf("GoSDKTest %d", time.Now().UnixNano()),
		Description:      "TestAddingRoleToScopeAndUser",
		ParentAppScopeId: rolesRootScopeAppId,
		ShortQuery: ShortQuery{
			Type:  "eq",
			Field: "ip",
			Value: "10.0.0.1",
		},
	}
	createdScope, err := client.CreateScope(createScopeParams)
	if err != nil {
		t.Error(err)
	}

	scopesToDelete = append(scopesToDelete, createdScope)

	// application creation
	newParams := createApplicationParams
	newParams.AppScopeId = createdScope.Id
	createdApplication, err := client.CreateApplication(newParams)
	if err != nil {
		t.Fatal(err)
	}
	applicationsToDelete = append(applicationsToDelete, createdApplication)

	// user creation
	user1Params := CreateUserRequest{
		Email:      fmt.Sprintf("test+%d@example.com", time.Now().UnixNano()),
		FirstName:  "Eric",
		LastName:   "Hibbs",
		AppScopeId: rolesRootScopeAppId,
	}
	user2Params := CreateUserRequest{
		Email:     fmt.Sprintf("test+%d@example.com", time.Now().UnixNano()),
		FirstName: "Eric",
		LastName:  "Hibbs",
		// TODO: why won't this work with createdScope.id, and does that ruin the use case?
		AppScopeId: rolesRootScopeAppId,
	}

	user1, err := client.CreateUser(user1Params)
	if err != nil {
		t.Fatal(err)
	}

	user2, err := client.CreateUser(user2Params)
	if err != nil {
		t.Fatal(err)
	}
	usersToDelete = append(usersToDelete, user1, user2)

	// role creation
	createScopedRoleForUsersParams := CreateScopedRoleForUsersRequest{
		CreateScopedRoleRequest: CreateScopedRoleRequest{
			Name:                fmt.Sprintf("GoSDKTest %d", time.Now().UnixNano()),
			Description:         "TestCreatingScopedRoleForUsers",
			AppScopeId:          rolesRootScopeAppId,
			AbilitiesAppScopeId: createdApplication.AppScopeId,
			Ability:             "SCOPE_READ",
		},
		Users: []string{user1.Id, user2.Id},
	}

	createScopedRoleForUsersResponse, err := client.CreateScopedRoleForUsers(createScopedRoleForUsersParams)
	if err != nil {
		t.Fatalf("Error %s creating scoped role for users with params %+v and config %+v", err, createScopedRoleForUsersParams, rolesDefaultClientConfig)
	}
	createdRole := createScopedRoleForUsersResponse.Role
	rolesToDelete = append(rolesToDelete, createdRole)

	// validation
	describedUser1 := createScopedRoleForUsersResponse.Users[0]
	describedUser2 := createScopedRoleForUsersResponse.Users[1]

	if !doesUserHaveRole(describedUser1, createdRole) {
		t.Errorf("User 1 expected to have role %s, but only had [%s]", createdRole.Id, strings.Join(describedUser1.RoleIds, ", "))
	}
	if !doesUserHaveRole(describedUser2, createdRole) {
		t.Errorf("User 1 expected to have role %s, but only had [%s]", createdRole.Id, strings.Join(describedUser2.RoleIds, ", "))
	}
}

func doesUserHaveRole(user User, role Role) bool {
	for _, userRole := range user.RoleIds {
		if role.Id == userRole {
			return true
		}
	}
	return false
}
