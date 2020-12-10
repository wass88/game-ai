package server

import (
	"database/sql"
	"math"
	"sort"

	"github.com/pkg/errors"
)

type EroRating struct {
	InitialRating float64
	K float64
}

var eroRating EroRating = EroRating{1500, 32}

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
	if e.K <= 0 {
		panic("ero.K must be > 0")
	}
	res := make([]float64, len(old))
	copy(res, old)

	sorted := IndexSort(score)
	for ri := 0; ri < len(old); ri++ {
		i := sorted[ri]
		for rj := ri+1; rj < len(old); rj++ {
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
			winRate := 1.0 / (math.Pow(10, (old[lose] - old[win])/400) + 1)
			diff := e.K * winRate
			res[win] += diff
			res[lose] -= diff
		}
	} 
	return res
}

func (p *PlayoutID) FetchRated() (bool, error) {
	res := []bool{}
	err := p.DB.DB.Select(&res, `SELECT rated FROM playout WHERE id = ?`, p.ID)
	if err != nil {
		return false, errors.Wrapf(err, "Fetch Rated")
	}
	return res[0], nil
}

func (p *PlayoutID) FetchLatestRate() ([]float64, bool, error) {
	type Result struct {
		AIID int64 `db:"ai_id"`
		Rate sql.NullFloat64 `db:"rate"`
	}
	var res []Result
	err := p.DB.DB.Select(&res, `SELECT playout_ai.ai_id, rate.rate
	FROM playout
	INNER JOIN playout_ai ON playout_ai.playout_id = playout.id
	LEFT JOIN (
	SELECT o_playout_ai.ai_id, o_result_ai.rate, o_playout.game_id
	  FROM playout AS o_playout
	  INNER JOIN playout_ai AS o_playout_ai ON o_playout_ai.playout_id = o_playout.id
	  INNER JOIN playout_result_ai AS o_result_ai ON o_result_ai.turn = o_playout_ai.turn AND o_result_ai.playout_id = o_playout_ai.playout_id
	  WHERE NOT EXISTS(
		 SELECT 1 FROM playout_ai AS t_playout_ai
		   INNER JOIN playout_result_ai AS t_result_ai ON t_result_ai.turn = t_playout_ai.turn AND t_result_ai.playout_id = t_playout_ai.playout_id
		   INNER JOIN playout AS t_playout ON t_playout.id = t_playout_ai.playout_id
		   WHERE o_playout.game_id = t_playout.game_id
			 AND o_playout_ai.ai_id = t_playout_ai.ai_id
			 AND o_result_ai.created_at <= t_result_ai.created_at
			 AND o_result_ai.id < t_result_ai.id
		 )
	) AS rate ON rate.ai_id = playout_ai.ai_id AND playout.game_id = rate.game_id
	WHERE playout.id = ?
	ORDER BY playout_ai.turn
	`, p.ID)
	if err != nil {
		return nil, false, errors.Wrapf(err, "Fetch Latest Rate")
	}
	rate := []float64{}
	ais := map[int64]interface{}{}
	for _, r := range res {
		t := eroRating.InitialRating
		if r.Rate.Valid {
			t = r.Rate.Float64
		}
		rate = append(rate, t)

		ais[r.AIID] = nil
	}

	selfMatch := false
	if len(ais) < len(res) {
		selfMatch = true
	}
	return rate, selfMatch, nil
}
