package main

// Sum возвращает сумму элементов.
func Sum(values ...int) int {
	var sum int
	for _, v := range values {
		sum += v
	}
	return sum
}

func main() {}
