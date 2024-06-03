package utils_test

import (
	"go_project_template/internal/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashSHA512(t *testing.T) {
	table := map[string]string{
		"Hello, World!": "374d794a95cdcfd8b35993185fef9ba368f160d8daf432d08ba9f1ed1e5abe6cc69291e0fa2fe0006a52570ef18c19def4e617c33ce52ef0a6e5fbe318cb0387",
		"Hello, 123!":   "c242f458c7473bb510c3d273e593f4d16ea4d45a1b763308ae4e780c5e6175afdfc3b2735ea81a9a9b84a8b5515749d5d443f641d7aed13236295303ecf420b5",
		"Hello, 123":    "84df6bdafdaa325beeaa4dedf46e6519e350ac7c9936d44d5b1de84359572d3d7047bece9f25dbb12876b9f307bb994f0df737b87757a0081583f3b23b7d4a4b",
		"Привет, 123":   "12d0223b0f141766b6629870a77af593f5958e5df8e9210b6facb0c4583b3552ed4447ca40ff9193dfe665fe528b9ff50ebe949d545e66fd5f9f23a8672d4b2e",
		"Это тест":      "86b29139f4aabc93cc861dd1400b2b65a8af6e62c095cb09ad905f5bdc3a6ef8315f69dee2a5d59b6d8c65c7c62171437f10d975f831c66a6018a37646495a86",
	}
	for src, res := range table {
		require.Equal(t, res, utils.HashSHA512([]byte(src)))
	}
}

func TestHashSHA256(t *testing.T) {
	table := map[string]string{
		"Hello, World!": "dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f",
		"Hello, 123!":   "52ce7f8d6d1f955e01b896c8fb38de421b4f9d0a2978fb1e3a3c9f3a6efa80ff",
		"Hello, 123":    "30b6bfae65bce9ae9ab1cef925407ddc3bcc3ee3ccbb4991619a4d7cd0c72675",
		"Привет, 123":   "cbb149fb48cac87fe20a177d260e2bd25d8f712a847af73e769dea07349fd127",
		"Это тест":      "6ef00d219391d63e93b7f130f46776ab449fd629ec461c8be3d8385234e2c6a2",
	}
	for src, res := range table {
		require.Equal(t, res, utils.HashSHA256([]byte(src)))
	}
}
