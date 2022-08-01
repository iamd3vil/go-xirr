package xirr

import (
	"fmt"
	"math"
	"sort"
	"time"
)

const MaxError float64 = 1e-10
const MaxComputeWithGuessIterations uint32 = 50

type (
	Payment struct {
		Date   time.Time
		Amount float64
	}

	Payments []Payment
)

func (p Payments) Len() int {
	return len(p)
}

func (p Payments) Less(i, j int) bool {
	return p[i].Date.Before(p[j].Date)
}

func (p Payments) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
	return
}

func Compute(payments Payments) (float64, error) {
	if err := validate(payments); err != nil {
		return 0, err
	}

	// Sort by date.
	sort.Sort(payments)

	var (
		rate  = computeWithGuess(payments, 0.1)
		guess = -0.99
	)

	for guess < 1 && (math.IsNaN(rate) || math.IsInf(rate, 0)) {
		rate = computeWithGuess(payments, guess)
		guess += 0.01
	}

	return rate, nil
}

func computeWithGuess(payments []Payment, guess float64) float64 {
	var (
		r = guess
		e = 1.0
	)

	for i := 0; i < int(MaxComputeWithGuessIterations); i++ {
		if e <= MaxError {
			return r
		}

		r1 := r - xirr(payments, r)/dxirr(payments, r)
		e = math.Abs(r1 - r)
		r = r1
	}

	return math.NaN()
}

func xirr(payments []Payment, rate float64) float64 {
	var result float64
	for _, p := range payments {
		exp := getExp(p, payments[0])
		result += p.Amount / math.Pow(1+rate, exp)
	}

	return result
}

func dxirr(payments []Payment, rate float64) float64 {
	var result float64
	for _, p := range payments {
		exp := getExp(p, payments[0])
		result -= p.Amount * exp / math.Pow(1+rate, exp+1)
	}

	return result
}

func validate(payments []Payment) error {
	var (
		positive, negative bool
	)

	for _, p := range payments {
		if p.Amount > 0 {
			positive = true
			break
		}
	}

	for _, p := range payments {
		if p.Amount < 0 {
			negative = true
			break
		}
	}

	if positive && negative {
		return nil
	} else {
		return fmt.Errorf("invalid payments")
	}
}

func getExp(p, p0 Payment) float64 {
	d := p.Date.Sub(p0.Date).Hours() / 24
	return d / 365.0
}
