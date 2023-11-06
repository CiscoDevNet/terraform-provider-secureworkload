// +build all integrationtests

package secureworkload

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var (
	tagsAPIURL               = os.Getenv("SECUREWORKLOAD_API_URL")
	tagsAPIKey               = os.Getenv("SECUREWORKLOAD_API_KEY")
	tagsAPISecret            = os.Getenv("SECUREWORKLOAD_API_SECRET")
	tagsDefaultRootScopeName = os.Getenv("SECUREWORKLOAD_ROOT_SCOPE_APP_NAME")
	tagsDefaultClientConfig  = Config{
		APIKey:                 tagsAPIKey,
		APISecret:              tagsAPISecret,
		APIURL:                 tagsAPIURL,
		DisableTLSVerification: false,
	}
)

func TestDeleteTagDeletesCreatedTag(t *testing.T) {
	client, err := New(tagsDefaultClientConfig)
	if err != nil {
		t.Errorf("Error %s creating client with config %+v", err, tagsDefaultClientConfig)
	}
	tagIp := "10.0.0.1"
	createTagParams := CreateTagRequest{
		RootScopeName: tagsDefaultRootScopeName,
		Ip:            tagIp,
		Attributes: map[string]interface{}{
			"environment": "test",
			"app_name":    fmt.Sprintf("TestDeleteTagDeletesCreatedTag%d", time.Now().Unix()),
		},
	}
	_, err = client.CreateTag(createTagParams)
	if err != nil {
		t.Error(err)
	}
	deleteTagsRequest := DeleteTagRequest{
		RootAppScopeName: tagsDefaultRootScopeName,
		Ip:               tagIp}
	err = client.DeleteTag(deleteTagsRequest)
	if err != nil {
		t.Errorf("error %s deleting tag: %+v\n", err, deleteTagsRequest)
	}
}

func TestDescribeTagDescribesCreatedTag(t *testing.T) {
	client, err := New(tagsDefaultClientConfig)
	if err != nil {
		t.Errorf("Error %s creating client with config %+v", err, tagsDefaultClientConfig)
	}
	tagIp := "10.0.0.1"
	createTagParams := CreateTagRequest{
		RootScopeName: tagsDefaultRootScopeName,
		Ip:            tagIp,
		Attributes: map[string]interface{}{
			"Environment": "test",
			"app_name":    fmt.Sprintf("TestDeleteTagDeletesCreatedTag%d", time.Now().Unix()),
		},
	}
	_, err = client.CreateTag(createTagParams)
	if err != nil {
		t.Error(err)
	}
	tagAttributesTemplate := make(map[string]string)
	describeTagRequest := DescribeTagRequest{
		RootAppScopeName: tagsDefaultRootScopeName,
		Ip:               tagIp}
	err = client.DescribeTag(describeTagRequest, &tagAttributesTemplate)
	if err != nil {
		t.Log(err)
	}
	for key, value := range createTagParams.Attributes {
		if tagAttributesTemplate[key] != value || tagAttributesTemplate[key] == "" {
			t.Errorf("expected tag attribute %s to equal %s, got %s", key, value, tagAttributesTemplate[key])
		}
	}
	deleteTagsRequest := DeleteTagRequest{
		RootAppScopeName: tagsDefaultRootScopeName,
		Ip:               tagIp,
	}
	err = client.DeleteTag(deleteTagsRequest)
	if err != nil {
		t.Errorf("error %s deleting tag: %+v\n", err, deleteTagsRequest)
	}
}
