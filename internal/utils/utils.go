package utils

import "strconv"

func CheckNumberObject(text string) (string, bool) {

	num, err := strconv.Atoi(text)
	if err != nil {
		return "Номер объекта введен некорректно!", false
	}

	if num < 1 || num > 9999 {
		return "Номер объекта введен некорректно!", false
	}
	return "", true
}
