package utils

func ArrayCompare[T any](original []T, new []T) (all []T, missing []T) {
	mapOrigin := make(map[any]bool)

	for _, v := range original {
		mapOrigin[v] = true
		all = append(all, v)
	}

	for _, v := range new {
		if _, ok := mapOrigin[v]; !ok {
			missing = append(missing, v)
			all = append(all, v)
		}
	}

	return all, missing
}
