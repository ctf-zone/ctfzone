package models_test

import (
	"math"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	udb "upper.io/db.v3"

	"github.com/ctf-zone/ctfzone/config"
	. "github.com/ctf-zone/ctfzone/models"
)

type ScoresList []*Score

func (s ScoresList) Len() int {
	return len(s)
}

func (s ScoresList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ScoresList) Less(i, j int) bool {
	if s[i].Score == 0 {
		return s[j].Score > 0 || (s[j].Score == 0 && s[i].Name > s[j].Name)
	} else {
		return (s[i].Score < s[j].Score) ||
			(s[i].Score == s[j].Score && s[i].UpdatedAt.After(*s[j].UpdatedAt)) ||
			(s[i].Score == s[j].Score && *s[i].UpdatedAt == *s[j].UpdatedAt && s[i].Name > s[j].Name)
	}
}

func Test_Scores_List_ClassicScoring(t *testing.T) {
	setup(t)
	defer teardown(t)

	scoresActual, _, err := db.ScoresList(
		&config.Scoring{Type: "classic"},
		ScoresPagination(Pagination{Count: 10}),
	)
	assert.NoError(t, err)

	scoresExpected := make([]*Score, 0)

	assert.NoError(t,
		upperDB.
			SelectFrom("users").
			Columns("id", "name", "extra").
			All(&scoresExpected),
	)

	for _, s := range scoresExpected {

		rows, err := upperDB.
			SelectFrom("solutions").
			Columns("points").
			Join("challenges").On("challenges.id = challenge_id").
			Where("user_id", s.ID).
			Query()
		require.NoError(t, err)

		defer rows.Close()

		for rows.Next() {
			var p int
			require.NoError(t, rows.Scan(&p))
			s.Score += p
		}
		require.NoError(t, rows.Err())

		row, err := upperDB.
			SelectFrom("solutions").
			Columns(udb.Raw("MAX(created_at) AS updated_at")).
			Where("user_id", s.ID).
			QueryRow()

		require.NoError(t, err)
		require.NoError(t, row.Scan(&s.UpdatedAt))
	}

	sort.Sort(sort.Reverse(ScoresList(scoresExpected)))

	for i, s := range scoresExpected {
		s.Rank = i + 1
	}

	assert.Equal(t, scoresExpected, scoresActual)
}

func Test_Scores_List_DynamicScoring(t *testing.T) {
	setup(t)
	defer teardown(t)

	min, max := 100, 1000
	coeff := 0.99

	c := &config.Scoring{
		Type: "dynamic",
		Dynamic: struct {
			Min   int     `json:"min"`
			Max   int     `json:"max"`
			Coeff float64 `json:"coeff"`
		}{
			Min:   100,
			Max:   1000,
			Coeff: 0.99,
		},
	}

	scoresActual, _, err := db.ScoresList(c,
		ScoresPagination(Pagination{Count: 10}),
	)
	assert.NoError(t, err)

	scoresExpected := make([]*Score, 0)

	assert.NoError(t,
		upperDB.
			SelectFrom("users").
			Columns("id", "name", "extra").
			All(&scoresExpected),
	)

	for _, s := range scoresExpected {

		solvedChallengesRows, err := upperDB.
			SelectFrom("solutions").
			Columns("challenge_id").
			Where("user_id", s.ID).
			Query()
		require.NoError(t, err)

		defer solvedChallengesRows.Close()

		for solvedChallengesRows.Next() {
			var challengeID int
			require.NoError(t, solvedChallengesRows.Scan(&challengeID))

			var n int
			countRow, err := upperDB.
				SelectFrom("solutions").
				Columns(udb.Raw("COUNT(*)")).
				Where("challenge_id", challengeID).
				QueryRow()
			require.NoError(t, err)
			require.NoError(t, countRow.Scan(&n))

			s.Score += int(float64(min) + float64(max-min)*math.Pow(coeff, float64(n-1)))
		}
		require.NoError(t, solvedChallengesRows.Err())

		updatedAtRow, err := upperDB.
			SelectFrom("solutions").
			Columns(udb.Raw("MAX(created_at)")).
			Where("user_id", s.ID).
			QueryRow()

		require.NoError(t, updatedAtRow.Scan(&s.UpdatedAt))
	}

	sort.Sort(sort.Reverse(ScoresList(scoresExpected)))

	for i, s := range scoresExpected {
		s.Rank = i + 1
	}

	assert.Equal(t, scoresExpected, scoresActual)
}
