package helper

import (
	"math/rand"
	"strings"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateInitialName(name string) (res string) {

	names := strings.Fields(name)

	res = ""

	if len(names) > 2 {
		res = string(names[0][0]) + string(names[len(names)-1][0])
	} else {
		for _, element := range names {
			if element != "" {
				res += strings.ToUpper(element[0:1])
			}
		}
	}

	return res
}

