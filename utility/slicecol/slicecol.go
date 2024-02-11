package slicecol

//i番目の要素を削除する(順番を保証しない)
func RemoveFast[T any](s []T, i int) []T {
	if i < 0 {
		panic("[slicecol.RemoveFast] i must be positive")
	}

	s[i] = s[len(s)-1]
	s = s[:len(s)-1]
	return s
}

//i番目の要素を削除する(順番を保証する)
func Remove[T any](s []T, i int) []T {
	if i < 0 {
		panic("[slicecol.Remove] i must be positive")
	}
	s = s[:i+int(copy(s[i:], s[i+1:]))]
	return s
}

//i番目の要素からj番目までの要素の順序を逆転する
func Reverse[T any](s []T, i, j int) []T {

	if i < 0 || j < 0 {
		panic("[slicecol.Reverse] i and j must be positive")
	}

	n := int(len(s))
	if j >= n {
		panic("[slicecol.Reverse] j must be less than length of s")
	}

	if i >= j {
		return s
	}

	for i < j {
		s[i], s[j] = s[j], s[i]
		i++
		j--
	}
	return s
}
