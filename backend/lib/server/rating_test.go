package server

import (
	"testing"
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
		t.Fatalf("Expected [1] > [0] & [1] > [2] & [0] > [3] & [2] > [3]. %v", r)
	}
	r = EroRating{K:32}.Rating([]float64{1500,1500},[]int{1, 1})
	t.Logf("Rate %v", r)
	if r[0] - r[1] > 1e-5 {
		t.Fatalf("Expected %f == %f", r[0], r[1])
	}
	r = EroRating{K:32}.Rating([]float64{2812,140},[]int{-14, 14})
	t.Logf("Rate %v", r)
	if r[0] >= 2812 {
		t.Fatalf("Expected %f > 2812, %f", r[0], r[1])
	}
}


func TestFetchLatestRate(t *testing.T) {
	db := mockPlayoutResultDB()
	playout := PlayoutID{1, db}
	rate, selfMatch, err := playout.FetchLatestRate()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v %v", rate, selfMatch)
	if selfMatch {
		t.Fatalf("selfMatch should be false")
	}

	playout = PlayoutID{2, db}
	rate, selfMatch, err = playout.FetchLatestRate()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v %v", rate, selfMatch)
	if !selfMatch {
		t.Fatalf("selfMatch should be true")
	}
}