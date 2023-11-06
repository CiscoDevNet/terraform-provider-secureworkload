package secureworkload

import (
	"fmt"
	"net/http"

	// "github.com/secureworkload-exchange/terraform-go-sdk/signer"
	"terraform-provider-secureworkload/secureworkload/signer"
)

var (
	TagsAPIV1BasePath = fmt.Sprintf("%s/inventory/tags", SecureWorkloadAPIV1BasePath)
)

// Tag wraps annotations for tagging flows and inventory items in a root scope on the SecureWorkload appliance.
type Tag struct {
	// IPv4/IPv6 address or subnet.
	Ip string `json:"ip"`
	// Key/value map for tagging matching flows and inventory items.
	Attributes map[string]interface{} `json:"attributes"`
}

// CreateTagRequest wraps parameters for making a request to create a tag.
type CreateTagRequest struct {
	RootScopeName string
	Ip            string                 `json:"ip"`
	Attributes    map[string]interface{} `json:"attributes"`
}

// CreateTag creates a tags with the specified
// params, returning the created tags and error
// (if any).
func (c Client) CreateTag(params CreateTagRequest) (Tag, error) {
	var tag Tag
	url := c.Config.APIURL + TagsAPIV1BasePath + fmt.Sprintf("/%s", params.RootScopeName)
	request, err := signer.CreateJSONRequest(http.MethodPost, url, params)
	if err != nil {
		return tag, err
	}
	err = c.Do(request, nil)
	if err != nil {
		return tag, err
	}
	tag = Tag{
		Ip:         params.Ip,
		Attributes: params.Attributes,
	}
	return tag, nil
}

// DescribeTagRequest wraps the parameters for a
// describeTag request
type DescribeTagRequest struct {
	RootAppScopeName string
	Ip               string
}

// DescribeTag describes a tag by id
// returning the tag and error (if any).
func (c Client) DescribeTag(params DescribeTagRequest, attributesTemplate *map[string]string) error {
	url := c.Config.APIURL + TagsAPIV1BasePath + fmt.Sprintf("/%s?ip=%s", params.RootAppScopeName, params.Ip)
	request, err := signer.CreateJSONRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	err = c.Do(request, attributesTemplate)
	return err
}

// DeleteTagRequest wraps the parameters for a
// deleteTag request
type DeleteTagRequest struct {
	RootAppScopeName string
	Ip               string `json:"ip"`
}

// DeleteTag deletes all tag for a given ip under a root scope
// returning error if any
func (c Client) DeleteTag(params DeleteTagRequest) error {
	url := c.Config.APIURL + TagsAPIV1BasePath + fmt.Sprintf("/%s", params.RootAppScopeName)
	request, err := signer.CreateJSONRequest(http.MethodDelete, url, params)
	if err != nil {
		return err
	}
	return c.Do(request, nil)
}
