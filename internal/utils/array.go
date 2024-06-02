package utils

func FindMissing[T any](original []T, new []T) (missing []T) {
	mapOrigin := make(map[any]bool)

	for _, v := range original {
		mapOrigin[v] = true
	}

	for _, v := range new {
		if _, ok := mapOrigin[v]; !ok {
			missing = append(missing, v)
		}
	}

	return missing
}
