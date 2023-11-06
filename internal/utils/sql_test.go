package utils_test

import (
	"go_project_template/internal/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

type Fruit struct {
	name   string
	amount int
}

func TestGenerateInsertSQL(t *testing.T) {
	t.Parallel()

	result, params := utils.GenerateInsertSQL("fruits", map[string]any{
		"name": "amount",
	})
	require.Equal(t, "INSERT INTO fruits (name) VALUES ($1)", result)
	require.Len(t, params, 1)
	require.Equal(t, "amount", params[0])

	result, params = utils.GenerateInsertSQL("fruits", map[string]any{
		"name":  "amount",
		"count": 1,
	})
	require.Len(t, params, 2)
	if result == "INSERT INTO fruits (name, count) VALUES ($1, $2)" {
		require.Equal(t, "amount", params[0])
		require.Equal(t, 1, params[1])
	} else if result == "INSERT INTO fruits (count, name) VALUES ($1, $2)" {
		require.Equal(t, "amount", params[1])
		require.Equal(t, 1, params[0])
	} else {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestGenerateBulkInsertSQL(t *testing.T) {
	res, params := utils.GenerateBulkInsertSQL[Fruit]("sample", utils.PQParamPlaceholder, []Fruit{
		{"Apple", 10},
		{"Pear", 100},
		{"Cherry", 36},
		{"Banana", 4},
		{"Apricot", 99},
	}, func(entity Fruit) map[string]any {
		return map[string]any{
			"name":   entity.name,
			"amount": entity.amount,
		}
	})
	t.Log(res)
	valid := res == "INSERT INTO sample (name,amount) VALUES ($1,$2),($3,$4),($5,$6),($7,$8),($9,$10)" ||
		res == "INSERT INTO sample (amount,name) VALUES ($1,$2),($3,$4),($5,$6),($7,$8),($9,$10)"
	require.True(t, valid)
	require.Len(t, params, 10)
}
