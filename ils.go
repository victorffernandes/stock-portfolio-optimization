package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

var cardinality = 40
var money = 10000.0

var restriction = 0.2
var min_restriction = 0.01
var offset = 1.0

// L = range(cardinality)
// s = [money/cardinality] * cardinality

func increment(lista []float64, indice int, valor float64) []float64 {
	restante_anterior := money - lista[indice]
	if lista[indice]+valor >= money*restriction {
		lista[indice] = money * restriction
	} else {
		lista[indice] = lista[indice] + valor
	}

	restante := money - lista[indice] - (float64(len(lista)-1) * float64(money) * float64(min_restriction))

	for i := range lista {
		if i != indice {
			lista[i] = ((float64(lista[i]) / float64(restante_anterior)) * float64(restante)) + (money * min_restriction)
		}
	}

	return lista
}

func fitness(s []float64, L []int) float64 {
	fit := 0.0

	for i := range s {
		fit += (float64(s[i]) / float64(money)) * float64(L[i])
	}
	// print("fitness: ", fit)
	return fit
}

func main() {
	L := make([]int, cardinality)
	s := make([]float64, cardinality)

	for i := range s {
		s[i] = money / float64(cardinality)
		L[i] = i
	}

	localSearchVND := func(s []float64, k int) ([]float64, float64) {
		best := make([]float64, cardinality)
		copy(best, s)
		best_fit := fitness(best, L)
		i := 0
		for i < cardinality {
			viz_ := best
			for j := 1; j < k+1; j++ {
				if float64(viz_[i])+offset*float64(j) <= restriction*money && float64(viz_[i])+offset*float64(j) >= min_restriction*money {
					viz_ = increment(viz_, i, offset*float64(j))
					fit := fitness(viz_, L)
					if fit > best_fit {
						best_fit = fit
						best = viz_
						i = 0
					}
				}
			}
			i++
		}
		return best, best_fit
	}

	// disturbanceSlice := func(s []float64, m float64) ([]float64, float64) {
	// 	dist := make([]float64, cardinality)
	// 	copy(dist, s)
	// 	initial := math.Floor(rand.Float64() * float64(cardinality))
	// 	final := math.Round(rand.Float64() * float64(cardinality))

	// 	//print("disturbance: indexes", int(math.Min(initial, final)), " : ", int(math.Max(initial, final)))
	// 	for i := int(math.Min(initial, final)); i < int(math.Max(initial, final)); i++ {
	// 		increment(dist, i, dist[i]*m)
	// 	}
	// 	dist_fit := fitness(dist, L)
	// 	// print("heree")
	// 	// print(dist_fit)
	// 	return dist, dist_fit
	// }

	disturbanceSlice := func(s []float64, disturbance_type int) ([]float64, float64) { // swap(2,1)
		dist := make([]float64, cardinality)
		copy(dist, s)
		r1 := int(math.Floor(rand.Float64() * float64(cardinality)))
		r2 := int(math.Floor(rand.Float64() * float64(cardinality)))

		if disturbance_type == 0 {
			// swap(1,1)
			tmp := dist[r1]
			dist[r1] = dist[len(s)-1]
			dist[len(s)-1] = tmp

			// swap(1,1)
			tmp = dist[r2]
			dist[r2] = dist[len(s)-2]
			dist[len(s)-2] = tmp
		} else {

		}

		dist_fit := fitness(dist, L)
		return dist, dist_fit
	}

	ils := func(s []float64, max int) ([]float64, float64) {
		best, best_fit := localSearchVND(s, 1)
		for j := 0; j < max; j++ {
			disturb, s1 := disturbanceSlice(best)
			fmt.Println("Disturbance Fitness: ", s1, best_fit)
			localsearch, s2 := localSearchVND(disturb, 10)

			if s1 > best_fit {
				best_fit = s1
				copy(best, disturb)
				fmt.Print("\nLocal Best ILS: ", best_fit, j)
			} else if s2 > best_fit {
				best_fit = s2
				copy(best, localsearch)
				fmt.Print("\nLocal Best ILS: ", best_fit, j)
			}
		}
		return best, best_fit
	}

	ils_parallel := func(s []float64, max int, disturbances int) ([]float64, float64) {
		best, best_fit := localSearchVND(s, 1)
		var wg sync.WaitGroup
		var mutex = &sync.Mutex{}

		for j := 0; j < max; j++ {
			for i := 0; i < disturbances; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					disturb, s1 := disturbanceSlice(best)
					fmt.Println("Disturbance Fitness Parallel: ", s1, best_fit)
					localsearch, s2 := localSearchVND(disturb, 10)

					mutex.Lock()
					if s1 > best_fit {
						best_fit = s1
						copy(best, disturb)
					} else if s2 > best_fit {
						best_fit = s2
						copy(best, localsearch)
					}
					mutex.Unlock()
					runtime.Gosched()
				}()

				wg.Wait()
			}
		}
		return best, best_fit
	}

	s1 := make([]float64, cardinality)
	s2 := make([]float64, cardinality)
	s3 := make([]float64, cardinality)
	copy(s1, s)
	copy(s2, s)
	copy(s3, s)
	s1_t := time.Now()
	s1, _ = localSearchVND(s1, 10)
	elapsed_s1 := time.Since(s1_t)

	rand.Seed(15)
	s2_t := time.Now()
	s2, _ = ils(s2, 1000)
	elapsed_s2 := time.Since(s2_t)

	rand.Seed(15)
	s3_t := time.Now()
	s3, _ = ils_parallel(s3, 30, 4)
	elapsed_s3 := time.Since(s3_t)

	fmt.Print("\n", fitness(s1, L))
	fmt.Print("\nelapsed s1: ", elapsed_s1)

	fmt.Print("\n", fitness(s2, L))
	fmt.Print("\nelapsed s2: ", elapsed_s2)

	fmt.Print("\n", fitness(s3, L))
	fmt.Print("\nelapsed s3: ", elapsed_s3)
	fmt.Print("\n copy here s3: ", elapsed_s2, elapsed_s3)

	sum_s1, sum_s2, sum_s3 := 0.0, 0.0, 0.0

	for i := range s3 {
		sum_s1 += s1[i]
		sum_s2 += s2[i]
		sum_s3 += s3[i]
	}
	fmt.Print("\n", sum_s1)
	fmt.Print("\n", sum_s2)
	fmt.Print("\n", sum_s3)

	// fmt.Print(fitness(s, L))
}
