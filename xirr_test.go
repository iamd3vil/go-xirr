package xirr

import (
	"encoding/csv"
	"io"
	"math"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestRandom(t *testing.T) {
	payments, err := loadPayments("samples/random.csv")
	if err != nil {
		t.Fatalf("error loading payments: %v", err)
	}

	actual, err := Compute(payments)
	if err != nil {
		t.Fatalf("error computing xirr: %v", err)
	}

	cErr := math.Abs(actual - 0.1266061750083439787)
	if cErr > MaxError {
		t.Fatalf("expected error to be less than %v but the error is %v", MaxError, cErr)
	}
}

func loadPayments(fpath string) (Payments, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rdr := csv.NewReader(f)
	rdr.ReuseRecord = true

	var payments Payments

	for {
		rec, err := rdr.Read()
		if err != nil {
			if err == io.EOF {
				return payments, nil
			}
			return nil, err
		}

		amount, err := strconv.ParseFloat(rec[0], 64)
		if err != nil {
			return nil, err
		}

		date, err := time.Parse("02/01/06", rec[1])
		if err != nil {
			return nil, err
		}

		payments = append(payments, Payment{
			Date:   date,
			Amount: amount,
		})
	}
}
