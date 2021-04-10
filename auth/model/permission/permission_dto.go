package permission

func ToList(Permissions []*Permission) []string {
	response := make([]string, len(Permissions))
	for _, v := range Permissions {
		response = append(response, v.Name)
	}

	return response
}
