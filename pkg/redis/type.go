package redis

func StringMembersToInterface(members []string) []interface{} {
	is := make([]interface{}, 0, len(members))
	for _, member := range members {
		is = append(is, member)
	}

	return is
}
