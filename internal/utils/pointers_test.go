package utils_test

import (
	"go_project_template/internal/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestFromPointer(t *testing.T) {
	test := uuid.NewString()
	pointer := utils.ToPointer(test)
	result := utils.FromPointer(pointer)
	require.Equal(t, test, result)

	var sample *string
	result = utils.FromPointer(sample)
	require.Equal(t, "", result)
}
