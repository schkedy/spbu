package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const maxGoroutines = 10
const minSizeForParallel = 1000

// ParallelQuickSort сортирует срез data параллельно с использованием компаратора cmp и дженерик T
func ParallelQuickSort[T any](data []T, cmp func(a, b T) bool) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxGoroutines) // семафор для ограничения числа горутин

	var quicksort func(lo, hi int)
	quicksort = func(lo, hi int) {
		if lo >= hi {
			return
		}
		// Опорный элемент (используем последний среза)
		pivot := data[hi]
		i := lo
		for j := lo; j < hi; j++ {
			if cmp(data[j], pivot) {
				data[i], data[j] = data[j], data[i]
				i++
			}
		}
		data[i], data[hi] = data[hi], data[i]
		leftSize := i - 1 - lo + 1
		rightSize := hi - (i + 1) + 1
		// Функция запуск подзадач с контролем горутин
		run := func(f func()) {
			if (leftSize > minSizeForParallel) || (rightSize > minSizeForParallel) {
				sem <- struct{}{} // захват ресурса
				wg.Add(1)
				// тут происходит распаралеливание
				go func() {
					defer wg.Done()
					f()
					<-sem // освобождение ресурса
				}()
			} else {
				f()
			}
		}

		run(func() { quicksort(lo, i-1) })
		run(func() { quicksort(i+1, hi) })
	}

	quicksort(0, len(data)-1)
	wg.Wait()
}

func generateRandomSlice(n int, max int) []int {
	rand.Seed(time.Now().UnixNano())
	result := make([]int, n)
	for i := 0; i < n; i++ {
		result[i] = rand.Intn(max)
	}
	return result
}

func main() {
	n := 300      // длина среза
	maxVal := 400 // максимальное значение (не включая)

	data := generateRandomSlice(n, maxVal)
	// fmt.Println("Сгенерированный срез:", data)

	start := time.Now()
	ParallelQuickSort(data, func(a, b int) bool { return a < b })
	elapsed := time.Since(start)

	fmt.Println("Отсортированный массив:", data)
	fmt.Printf("Время сортировки: %s\n", elapsed)
}
