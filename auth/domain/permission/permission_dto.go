package permission

import "blogfa/auth/model"

// ToList, change Permission list into string list
func ToList(Permissions []*model.Permission) []string {
	response := make([]string, len(Permissions))
	for i, v := range Permissions {
		response[i] = v.Name
	}

	return response
}
