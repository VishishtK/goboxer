package goboxer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type UserStatus string

const (
	UserStatusActive                 UserStatus = "active"
	UserStatusInactive               UserStatus = "inactive"
	UserStatusCannotDeleteEdit       UserStatus = "cannot_delete_edit"
	UserStatusCannotDeleteEditUpload UserStatus = "cannot_delete_edit_upload"
)

func (us *UserStatus) String() string {
	if us == nil {
		return "<nil>"
	}
	return string(*us)
}

func (us *UserStatus) MarshalJSON() ([]byte, error) {
	if us == nil {
		return []byte("null"), nil
	} else {
		return []byte(`"` + us.String() + `"`), nil
	}
}

type UserRole string

func (ur *UserRole) String() string {
	if ur == nil {
		return "<nil>"
	}
	return string(*ur)
}

func (ur *UserRole) MarshalJSON() ([]byte, error) {
	if ur == nil {
		return []byte("null"), nil
	} else {
		return []byte(`"` + ur.String() + `"`), nil
	}
}

const (
	UserRoleAdmin   UserRole = "admin"
	UserRoleCoAdmin UserRole = "coadmin"
	UserRoleUser    UserRole = "user"
)

type EnterpriseType string

func (et *EnterpriseType) String() string {
	if et == nil {
		return "<nil>"
	}
	return string(*et)
}

func (et *EnterpriseType) MarshalJSON() ([]byte, error) {
	if et == nil {
		return []byte("null"), nil
	} else {
		return []byte(`"` + et.String() + `"`), nil
	}
}

const (
	EnterpriseTypeEnterprise EnterpriseType = "enterprise"
	EnterpriseTypeUser       EnterpriseType = "user"
)

type Enterprise struct {
	Type EnterpriseType
	Id   string
	Name string
}

type User struct {
	UserGroupMini
	apiInfo                       *apiInfo            `json:"-"`
	CreatedAt                     *time.Time          `json:"created_at,omitempty"`
	ModifiedAt                    *time.Time          `json:"modified_at,omitempty"`
	Language                      *string             `json:"language,omitempty"`
	Timezone                      *string             `json:"timezone,omitempty"`
	SpaceAmount                   int64               `json:"space_amount,omitempty"`
	SpaceUsed                     int64               `json:"space_used,omitempty"`
	MaxUploadSize                 int                 `json:"max_upload_size,omitempty"`
	Status                        *UserStatus         `json:"status,omitempty"`
	JobTitle                      *string             `json:"job_title,omitempty"`
	Phone                         *string             `json:"phone,omitempty"`
	Address                       *string             `json:"address,omitempty"`
	AvatarUrl                     *string             `json:"avatar_url,omitempty"`
	Role                          *UserRole           `json:"role,omitempty"`
	TrackingCodes                 []map[string]string `json:"tracking_codes,omitempty"`
	CanSeeManagedUsers            *bool               `json:"can_see_managed_users,omitempty"`
	IsSyncEnabled                 *bool               `json:"is_sync_enabled,omitempty"`
	IsExternalCollabRestricted    *bool               `json:"is_external_collab_restricted,omitempty"`
	IsExemptFromDeviceLimits      *bool               `json:"is_exempt_from_device_limits,omitempty"`
	IsExemptFromLoginVerification *bool               `json:"is_exempt_from_login_verification,omitempty"`
	Enterprise                    *Enterprise         `json:"enterprise,omitempty"`
	MyTags                        *[]string           `json:"my_tags,omitempty"`
	Hostname                      *string             `json:"hostname,omitempty"`
	IsPlatformAccessOnly          *bool               `json:"is_platform_access_only,omitempty"`
	ExternalAppUserId             *string             `json:"external_app_user_id,omitempty"`

	changeFlag uint64
}

const (
	cUserName uint64 = 1 << (iota)
	cUserRole
	cUserLanguage
	cUserIsSyncEnabled
	cUserJobTitle
	cUserPhone
	cUserAddress
	cUserSpaceAmount
	cUserTrackingCodes
	cUserCanSeeMangedUsers
	cUserTimezone
	cUserIsExemptFromDeviceLimits
	cUserIsExemptFromLoginVerification
	cUserIsExternalCollabRestricted
	cUserStatus
	cUserIsPasswordResetRequired
	cUserIsPlatformAccessOnly
)

func (u *User) ResourceType() BoxResourceType {
	return UserResource
}

func NewUser(api *ApiConn) *User {
	return &User{
		apiInfo: &apiInfo{api: api},
	}
}

var UserAllFields = []string{
	"type", "id", "name", "login", "created_at", "modified_at",
	"language", "timezone", "space_amount", "space_used", "max_upload_size",
	"status", "job_title", "phone", "address", "avatar_url", "role",
	"tracking_codes", "can_see_managed_users", "is_sync_enabled",
	"is_external_collab_restricted", "is_exempt_from_device_limits",
	"is_exempt_from_login_verification", "enterprise",
	"my_tags", "hostname", "is_platform_access_only", "external_app_user_id",
}

// Get Current User
//
// Get information about the user who is currently logged in (i.e. the user for whom this access token was generated).
// https://developer.box.com/reference#get-the-current-users-information
func (u *User) GetCurrentUserReq(fields []string) *Request {
	var url string
	var query string

	url = fmt.Sprintf("%s%s", u.apiInfo.api.BaseURL, "users/me")
	if fieldsParams := BuildFieldsQueryParams(fields); fieldsParams != "" {
		query = fmt.Sprintf("?%s", fieldsParams)
	}

	return NewRequest(u.apiInfo.api, url+query, GET, nil, nil)
}

// Get User
//
// Get information about a user in the enterprise. Requires enterprise administration authorization.
// https://developer.box.com/reference#users
func (u *User) GetCurrentUser(fields []string) (*User, error) {

	req := u.GetCurrentUserReq(fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusOK {
		return nil, newApiStatusError(resp.Body)
	}

	r := &User{apiInfo: &apiInfo{api: u.apiInfo.api}}
	err = UnmarshalJSONWrapper(resp.Body, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Get User
//
// Get information about a user in the enterprise. Requires enterprise administration authorization.
// https://developer.box.com/reference#users
func (u *User) GetUserReq(userId string, fields []string) *Request {
	var url string
	var query string

	url = fmt.Sprintf("%s%s%s", u.apiInfo.api.BaseURL, "users/", userId)
	if fieldsParams := BuildFieldsQueryParams(fields); fieldsParams != "" {
		query = fmt.Sprintf("?%s", fieldsParams)
	}

	return NewRequest(u.apiInfo.api, url+query, GET, nil, nil)
}
func (u *User) GetUser(userId string, fields []string) (*User, error) {

	req := u.GetUserReq(userId, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusOK {
		return nil, newApiStatusError(resp.Body)
	}

	r := &User{apiInfo: &apiInfo{api: u.apiInfo.api}}
	err = UnmarshalJSONWrapper(resp.Body, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// TODO Get User Avatar

func (u *User) CreateUserReq(fields []string) *Request {
	var url string
	var query string
	url = fmt.Sprintf("%s%s", u.apiInfo.api.BaseURL, "users")
	if fieldsParams := BuildFieldsQueryParams(fields); fieldsParams != "" {
		query = fmt.Sprintf("?%s", fieldsParams)
	}

	data := &User{}
	data.Login = u.Login
	data.Name = u.Name
	if u.changeFlag&cUserRole == cUserRole {
		data.Role = u.Role
	}
	if u.changeFlag&cUserLanguage == cUserLanguage {
		data.Language = u.Language
	}
	if u.changeFlag&cUserIsSyncEnabled == cUserIsSyncEnabled {
		data.IsSyncEnabled = u.IsSyncEnabled
	}
	if u.changeFlag&cUserJobTitle == cUserJobTitle {
		data.JobTitle = u.JobTitle
	}
	if u.changeFlag&cUserPhone == cUserPhone {
		data.Phone = u.Phone
	}
	if u.changeFlag&cUserAddress == cUserAddress {
		data.Address = u.Address
	}
	if u.changeFlag&cUserSpaceAmount == cUserSpaceAmount {
		data.SpaceAmount = u.SpaceAmount
	}
	if u.changeFlag&cUserTrackingCodes == cUserTrackingCodes {
		data.TrackingCodes = u.TrackingCodes
	}
	if u.changeFlag&cUserCanSeeMangedUsers == cUserCanSeeMangedUsers {
		data.CanSeeManagedUsers = u.CanSeeManagedUsers
	}
	if u.changeFlag&cUserTimezone == cUserTimezone {
		data.Timezone = u.Timezone
	}
	if u.changeFlag&cUserIsExemptFromDeviceLimits == cUserIsExemptFromDeviceLimits {
		data.IsExemptFromDeviceLimits = u.IsExemptFromDeviceLimits
	}
	if u.changeFlag&cUserIsExemptFromLoginVerification == cUserIsExemptFromLoginVerification {
		data.IsExemptFromLoginVerification = u.IsExemptFromLoginVerification
	}
	if u.changeFlag&cUserIsExternalCollabRestricted == cUserIsExternalCollabRestricted {
		data.IsExternalCollabRestricted = u.IsExternalCollabRestricted
	}
	if u.changeFlag&cUserStatus == cUserStatus {
		data.Status = u.Status
	}
	bodyBytes, _ := json.Marshal(data)

	return NewRequest(u.apiInfo.api, url+query, POST, nil, bytes.NewReader(bodyBytes))
}
func (u *User) CreateUser(fields []string) (*User, error) {

	req := u.CreateUserReq(fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusCreated {
		return nil, newApiStatusError(resp.Body)
	}

	r := &User{apiInfo: &apiInfo{api: u.apiInfo.api}}
	err = UnmarshalJSONWrapper(resp.Body, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (u *User) SetLogin(login string) *User {
	u.Login = &login
	return u
}
func (u *User) SetName(name string) *User {
	u.Name = &name
	u.changeFlag |= cUserName
	return u
}
func (u *User) SetRole(role UserRole) *User {
	u.Role = &role
	u.changeFlag |= cUserRole
	return u
}
func (u *User) SetLanguage(language string) *User {
	u.Language = &language
	u.changeFlag |= cUserLanguage
	return u
}
func (u *User) SetIsSyncEnabled(isSyncEnabled bool) *User {
	u.IsSyncEnabled = &isSyncEnabled
	u.changeFlag |= cUserIsSyncEnabled
	return u
}
func (u *User) SetJobTitle(jobTitle string) *User {
	u.JobTitle = &jobTitle
	u.changeFlag |= cUserJobTitle
	return u
}
func (u *User) SetPhone(phone string) *User {
	u.Phone = &phone
	u.changeFlag |= cUserPhone
	return u
}
func (u *User) SetAddress(address string) *User {
	u.Address = &address
	u.changeFlag |= cUserAddress
	return u
}
func (u *User) SetSpaceAmount(spaceAmount int64) *User {
	u.SpaceAmount = spaceAmount
	u.changeFlag |= cUserSpaceAmount
	return u
}
func (u *User) SetTrackingCodes(trackingCodes []map[string]string) *User {
	u.TrackingCodes = trackingCodes
	u.changeFlag |= cUserTrackingCodes
	return u
}
func (u *User) SetCanSeeManagedUsers(canSeeManagedUsers bool) *User {
	u.CanSeeManagedUsers = &canSeeManagedUsers
	u.changeFlag |= cUserCanSeeMangedUsers
	return u
}
func (u *User) SetTimezone(timezone string) *User {
	u.Timezone = &timezone
	u.changeFlag |= cUserTimezone
	return u
}
func (u *User) SetIsExemptFromDeviceLimits(isExemptFromDeviceLimits bool) *User {
	u.IsExemptFromDeviceLimits = &isExemptFromDeviceLimits
	u.changeFlag |= cUserIsExemptFromDeviceLimits
	return u
}

func (u *User) SetIsExemptFromLoginVerification(isExemptFromLoginVerification bool) *User {
	u.IsExemptFromLoginVerification = &isExemptFromLoginVerification
	u.changeFlag |= cUserIsExemptFromLoginVerification
	return u
}
func (u *User) SetIsExternalCollabRestricted(isExternalCollabRestricted bool) *User {
	u.IsExternalCollabRestricted = &isExternalCollabRestricted
	u.changeFlag |= cUserIsExternalCollabRestricted
	return u
}
func (u *User) SetStatus(status UserStatus) *User {
	u.Status = &status
	u.changeFlag |= cUserStatus
	return u
}

func (u *User) UpdateUserReq(userId string, fields []string) *Request {
	var url string
	url = fmt.Sprintf("%s%s%s?%s", u.apiInfo.api.BaseURL, "users/", userId, BuildFieldsQueryParams(fields))

	bodyBytes, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}

	return NewRequest(u.apiInfo.api, url, PUT, nil, bytes.NewReader(bodyBytes))
}
func (u *User) UpdateUser(userId string, fields []string) (*User, error) {

	req := u.UpdateUserReq(userId, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusOK {
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to update user info"))
		return nil, err
	}

	r := &User{apiInfo: &apiInfo{api: u.apiInfo.api}}
	err = json.Unmarshal(resp.Body, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (u *User) CreateAppUserReq(fields []string) *Request {
	var url string
	url = fmt.Sprintf("%s%s?%s", u.apiInfo.api.BaseURL, "users/", BuildFieldsQueryParams(fields))

	b := true
	u.IsPlatformAccessOnly = &b
	bodyBytes, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}

	return NewRequest(u.apiInfo.api, url, POST, nil, bytes.NewReader(bodyBytes))
}
func (u *User) CreateAppUser(fields []string) (*User, error) {

	req := u.CreateAppUserReq(fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusCreated {
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to create app user"))
		return nil, err
	}

	r := &User{apiInfo: &apiInfo{api: u.apiInfo.api}}
	err = json.Unmarshal(resp.Body, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (u *User) DeleteUserReq(userId string, notify bool, force bool) *Request {
	var url string
	url = fmt.Sprintf("%s%s%s?notify=%t&force=%t", u.apiInfo.api.BaseURL, "users/", userId, notify, force)

	return NewRequest(u.apiInfo.api, url, DELETE, nil, nil)
}
func (u *User) DeleteUser(userId string, notify bool, force bool) error {

	req := u.DeleteUserReq(userId, notify, force)
	resp, err := req.Send()
	if err != nil {
		return err
	}

	if resp.ResponseCode != http.StatusNoContent {
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to delete user"))
		return err
	}

	return nil
}

func (u *User) GetEnterpriseUsersReq(filterTerm string, offset int32, limit int32, fields []string) *Request {
	var url string
	url = fmt.Sprintf("%s%s?filter_term=%s&offset=%d&limit=%d&%s", u.apiInfo.api.BaseURL, "users", filterTerm, offset, limit, BuildFieldsQueryParams(fields))

	return NewRequest(u.apiInfo.api, url, GET, nil, nil)
}
func (u *User) GetEnterpriseUsers(filterTerm string, offset int32, limit int32, fields []string) (outUsers []*User, outOffset int, outLimit int, outTotalCount int, err error) {

	req := u.GetEnterpriseUsersReq(filterTerm, offset, limit, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, 0, 0, 0, err
	}

	if resp.ResponseCode != http.StatusOK {
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to get enterprise users info"))
		return nil, 0, 0, 0, err
	}
	users := struct {
		TotalCount int     `json:"total_count"`
		Entries    []*User `json:"entries"`
		Offset     int     `json:"offset"`
		Limit      int     `json:"limit"`
	}{}
	err = json.Unmarshal(resp.Body, &users)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	for _, user := range users.Entries {
		user.apiInfo = &apiInfo{api: u.apiInfo.api}
	}
	return users.Entries, users.Offset, users.Limit, users.TotalCount, nil
}
