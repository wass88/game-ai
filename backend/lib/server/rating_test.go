package server

import (
	"testing"
	"time"
)


func TestEroRating(t *testing.T) {
	r := EroRating{K:32}.Rating([]float64{1500,1500},[]int{1, 0})
	t.Logf("Rate %v", r)
	if r[0] <= r[1] {
		t.Fatalf("R0 must be bigger than R1, %f <= %f", r[0], r[1])
	}
	r = EroRating{K:32}.Rating([]float64{1500,1500,1500,1500},[]int{1, 2, 1, 0})
	t.Logf("Rate %v", r)
	if !(r[1] > r[0] && r[1] > r[2] && r[0] > r[3] && r[2] > r[3]) {
		t.Fatalf("Expeted [1] > [0] & [1] > [2] & [0] > [3] & [2] > [3]. %v", r)
	}
}

func TestFetchPlayoutsForRating(t *testing.T) {
	db := mockPlayoutResultDB()
	res, err := db.FetchPlayoutsForRating()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", res)
	ais := res.AIsForUpdate()
	t.Logf("%v", ais)
}

func TestFetchRatingAIs(t *testing.T) {
	db :=mockPlayoutResultDB()
	rating, err := db.FetchRateAIs([]int64{1, 2})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", rating)
}

func TestCalclateRating(t *testing.T) {
	playouts := PlayoutResultsForRating {
		{1, time.Time{}, []int64{1, 2}, []int{-1, 1}},
	}
	rating := RateAI{1: 1500, 2: 1500}
	r, err := playouts.CalculateRating(rating)
	if err != nil {
		t.Fatal(err)
	}
	if r[1] >= r[2] {
		t.Fatal("R[1] should be worse than R[2]", r[1], r[2])
	}
	t.Logf("%v", r)
}

func TestUpdateRateAIs(t *testing.T) {
	r1 := 1200.0
	r2 := 1600.0
	rates := RateAI{1: r1, 2:r2}
	db := mockPlayoutResultDB()
	err := db.UpdateRateAIs(rates)
	if err != nil {
		t.Fatal(err)
	}
	res, err := db.FetchRateAIs([]int64{1, 2})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", res)
	eps := 0.01
	if res[1] - r1 > eps {
		t.Fatalf("res[0]:%f shuold be %f", res[1], r1)
	}
	if res[2] - r2 > eps {
		t.Fatalf("res[0]:%f shuold be %f", res[2], r2)
	}
}

func TestUpdateRating(t *testing.T) {
	db := mockPlayoutResultDB()
	err := db.UpdateRating()
	if err != nil {
		t.Fatal(err)
	}
}