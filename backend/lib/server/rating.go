package server

import (
	"math"
	"sort"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func (db *DB) KickUpdateRating() error {
	err := db.UpdateRating()
	if err != nil {
		return errors.Wrapf(err, "Update Rating")
	}
	return nil
}

type EroRating struct {
	InitialRating float64
	K float64
}

var eroRating EroRating = EroRating{500, 32}

func IndexSort(values []int) []int {
	type IndexedValue struct {
		value int
		index int
	}
	indexedValues := make([]IndexedValue, len(values))
	for i, v := range values {
		indexedValues[i] = IndexedValue{v, i}
	}
	sort.SliceStable(indexedValues,
		func(i, j int) bool {return indexedValues[i].value < indexedValues[j].value})
	res := make([]int, len(values))
	for i, v := range indexedValues {
		res[i] = v.index
	}
	return res
}

func (e EroRating) Rating(old []float64, score []int) []float64 {
	res := make([]float64, len(old))
	copy(res, old)

	sorted := IndexSort(score)
	for ri := 0; ri <len(old); ri++ {
		i := sorted[ri]
		for rj :=ri+1; rj<len(old); rj++ {
			j := sorted[rj]
			if score[i] == score[j] {
				continue
			}
			win := i
			lose := j
			if score[i] < score[j] {
				win = j
				lose = i
			}
			winRate := 1 / (math.Pow(10, old[win] - old[lose])+ 1)
			newWin := res[win] + e.K * (1 - winRate)
			newLose := res[lose] + e.K * (-winRate)
			res[win] = newWin
			res[lose] = newLose
		}
	} 
	return res
}

type PlayoutResultForRating struct {
	PlayoutID int64
	CreatedAt time.Time
	AIs []int64
	Results []int
}
func (db *DB) FetchPlayoutsForRating() (PlayoutResultsForRating, error) {
	type Result struct {
		PlayoutID int64     `db:"playout_id"`
		Turn int64 `db:"turn"`
		AIID int64 `db:"ai_id"`
		Result int `db:"result"`
		CreatedAt time.Time `db:"created_at"`
	}
	var results []Result
	err := db.DB.Select(&results, `
	SELECT result_ai.playout_id, result_ai.turn, playout_ai.ai_id, result_ai.result, result.created_at
	FROM playout_result_ai AS result_ai
	INNER JOIN playout_ai
	ON playout_ai.playout_id = result_ai.playout_id AND playout_ai.turn = result_ai.turn
	INNER JOIN playout_result AS result
	ON result.playout_id = result_ai.playout_id
	WHERE result_ai.playout_id IN (
		/* not calclated playout id */
		SELECT result.playout_id
		FROM playout_result AS result
		INNER JOIN playout_ai
		ON result.playout_id = playout_ai.playout_id
		LEFT JOIN rate_ai
		ON playout_ai.ai_id = rate_ai.ai_id
		GROUP BY result.id
		HAVING SUM(CASE WHEN ISNULL(rate_ai.updated_at) THEN 1 WHEN result.created_at > rate_ai.updated_at THEN 1 ELSE 0 END) > 0
	)
	ORDER BY result_ai.playout_id
	`)
	if err != nil {
		return nil, errors.Wrapf(err, "Fetch results")
	}

	res := []PlayoutResultForRating{}
	var result *PlayoutResultForRating
	for _, r := range results {
		if result != nil && result.PlayoutID != r.PlayoutID {
			res = append(res, *result)
			result = nil
		}
		if result == nil {
			result = &PlayoutResultForRating{
				PlayoutID: r.PlayoutID,
				CreatedAt: r.CreatedAt,
				AIs: []int64{},
				Results: []int{},
			}
		}
		result.AIs = append(result.AIs, r.AIID)
		result.Results = append(result.Results, r.Result)
	}
	res = append(res, *result)
	return res, nil
}

type PlayoutResultsForRating []PlayoutResultForRating
func (p PlayoutResultsForRating) AIsForUpdate() []int64 {
	m := map[int64]int{}
	for _, playout := range p {
		for _, ai_id := range playout.AIs {
			m[ai_id] = 1
		}
	}
	var res []int64
	for ai  := range m {
		res = append(res, ai)
	}
	return res
}

type RateAI map[int64]float64

func (db *DB) FetchRateAIs(ai_ids []int64) (RateAI, error) {
	type Result struct{
		ID int64 `db:"ai_id"`
		Rate float64 `db:"rate"`
	}
	var res []Result
	query, args, err := sqlx.In(`SELECT ai_id, rate FROM rate_ai WHERE ai_id IN (?)`, ai_ids)
	query = db.DB.Rebind(query)
	if err != nil {
		return nil, errors.Wrapf(err, "sqlx.in stmt")
	}
	err = db.DB.Select(&res, query, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "select rate")
	}
	rating := RateAI{}
	for _, a := range ai_ids {
		rating[a] = eroRating.InitialRating
	}
	for _, r := range res {
		rating[r.ID] = r.Rate
	}
	return rating, nil
}


func (p PlayoutResultsForRating) CalculateRating(rate RateAI) (RateAI, error) {
	for _, r := range p {
		rates := []float64{}
		for _, ai := range r.AIs {
			if _, ok := rate[ai]; !ok {
				return nil, errors.Errorf("Missing rating %d in %v", ai, rate)
			}
			rates = append(rates, rate[ai])
		}
		rates = eroRating.Rating(rates, r.Results)
		for i, ai := range r.AIs {
			rate[ai] = rates[i]
		}
	}
	return rate, nil
}

func (db *DB) UpdateRateAIs(r RateAI) error{
	type RateDB struct {
		AIID int64 `db:"ai_id"`
		Rate float64 `db:"rate"`
	}
	update := []RateDB{}
	for ai, rate := range r {
		update = append(update, RateDB{ai, rate})
	}
	_, err := db.DB.NamedExec(`INSERT INTO rate_ai (ai_id, rate) VALUES(:ai_id, :rate)`, update)
	if err != nil {
		return errors.Wrapf(err, "Update rate")
	}
	return nil
}


//UpdateRating calcs all ratings
func (db *DB) UpdateRating() error {
	playouts, err := db.FetchPlayoutsForRating()
	if err != nil {
		return errors.Wrapf(err, "Fetch playouts")
	}
	rates, err := db.FetchRateAIs(playouts.AIsForUpdate())
	if err != nil {
		return errors.Wrapf(err, "Fetch Rates")
	}
	rates, err = playouts.CalculateRating(rates)
	if err != nil {
		return errors.Wrapf(err, "Caluclate Rates")
	}
	err = db.UpdateRateAIs(rates)
	if err != nil {
		return errors.Wrapf(err, "Update Rates")
	}
	return nil
}