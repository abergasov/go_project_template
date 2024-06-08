package utils_test

import (
	"context"
	"fmt"
	"go_project_template/internal/utils"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_FuncRepeater(t *testing.T) {
	t.Run("should execute one by one", func(t *testing.T) {
		cnt := int32(0)
		fn := func() (interface{}, error) {
			atomic.AddInt32(&cnt, 1)
			time.Sleep(100 * time.Millisecond)
			require.Equal(t, int32(1), atomic.LoadInt32(&cnt))
			atomic.AddInt32(&cnt, -1)
			return []string{""}, nil
		}
		_, err := utils.NewFuncRepeater(fn).
			WithErrMsg("sample").
			WithTriesTimeout(1 * time.Second).
			WithLinearTimeout(true).
			WithMaxTries(10).Run()
		require.NoError(t, err)
	})
	t.Run("should stop when max retries reached", func(t *testing.T) {
		cnt := 0
		fn := func() (interface{}, error) {
			cnt++
			return nil, nil
		}

		tries := 10
		res, err := utils.NewFuncRepeater(fn).WithMaxTries(tries).Run()
		require.ErrorIs(t, err, utils.ErrorRepeatableFuncResNil)
		require.Nil(t, res)
		require.Equal(t, cnt, tries)
	})
	t.Run("should work with bool return type", func(t *testing.T) {
		cnt := 0
		fn := func() (bool, error) {
			cnt++
			return true, nil
		}

		tries := 10
		res, err := utils.NewFuncRepeater(fn).WithMaxTries(tries).Run()
		require.NoError(t, err)
		require.Equal(t, true, res)
		require.Equal(t, cnt, 1)
	})

	t.Run("should switch to fallback when max retries reached and retry fallbackRetries more", func(t *testing.T) {
		cnt, fbCnt := 0, 0
		fn := func() (interface{}, error) {
			cnt++
			return nil, nil
		}
		fbFn := func() (interface{}, error) {
			fbCnt++
			return nil, nil
		}

		tries, fbTries := 10, 5
		res, err := utils.NewFuncRepeater(fn).WithFallback(fbFn).WithMaxTries(tries).WithFallbackTries(fbTries).Run()
		require.ErrorIs(t, err, utils.ErrorRepeatableFuncResNil)
		require.Equal(t, res, nil)
		require.Equal(t, cnt, tries)
		require.Equal(t, fbCnt, fbTries)
	})

	t.Run("should stop when context is done", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cnt, cancelIdx := 0, 107
		fn := func() (interface{}, error) {
			if cnt == cancelIdx {
				cancel()
			}
			cnt++
			return nil, nil
		}

		res, err := utils.NewFuncRepeater(fn).WithCtx(ctx).Run()
		require.ErrorIs(t, err, utils.ErrorRepeatableFuncResNil)
		require.Equal(t, res, nil)
		require.Equal(t, cnt, cancelIdx+1)
	})

	t.Run("should stop when exitErrors encountered", func(t *testing.T) {
		tstErr := fmt.Errorf("test err")
		cnt, cancelIdx := 0, 117
		fn := func() (interface{}, error) {
			if cnt == cancelIdx {
				return nil, tstErr
			}
			cnt++
			return nil, nil
		}

		res, err := utils.NewFuncRepeater(fn).WithExitErrors(tstErr).Run()
		require.ErrorIs(t, err, tstErr)
		require.Equal(t, res, nil)
		require.Equal(t, cnt, cancelIdx)
	})

	t.Run("should return non-nil result", func(t *testing.T) {
		expectedRes := []byte{0x10, 0x20, 0x30}
		cnt, cancelIdx := 0, 107
		fn := func() ([]byte, error) {
			if cnt == cancelIdx {
				return expectedRes, nil
			}
			cnt++
			return nil, nil
		}

		res, err := utils.NewFuncRepeater(fn).Run()
		require.NoError(t, err)
		require.Equal(t, expectedRes, res)
		require.Equal(t, cnt, cancelIdx)
	})

	t.Run("should return result from pointer return type", func(t *testing.T) {
		data := uuid.NewString()
		fn := func() (*string, error) {
			return utils.ToPointer(data), nil
		}
		res, err := utils.NewFuncRepeater(fn).Run()
		require.NoError(t, err)
		require.Equal(t, data, *res)
	})
	t.Run("should return result from struct pointer return type", func(t *testing.T) {
		type Tc struct {
			Data  string
			Data2 string
		}
		data := Tc{
			Data:  uuid.NewString(),
			Data2: uuid.NewString(),
		}
		fn := func() (*Tc, error) {
			return &Tc{
				Data:  data.Data,
				Data2: data.Data2,
			}, nil
		}
		res, err := utils.NewFuncRepeater(fn).Run()
		require.NoError(t, err)
		require.Equal(t, data, *res)
	})
	t.Run("should return ok in case of array of bytes", func(t *testing.T) {
		fn := func() (common.Hash, error) {
			return common.HexToHash("0xb62b8c1cc42a0e90f047b7dc88850a2facfcd4e8e87e874bf0cfd28c6adb7af2"), nil
		}
		res, err := utils.NewFuncRepeater(fn).Run()
		require.NoError(t, err)
		t.Logf("res: %v", res.String())
	})
	t.Run("should return ok in case of array of bytes return err", func(t *testing.T) {
		fn := func() (txHash common.Hash, err error) {
			return txHash, fmt.Errorf("test error")
		}
		res, err := utils.NewFuncRepeater(fn).WithMaxTries(3).WithFallbackTries(3).Run()
		require.Error(t, err)
		require.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000000", res.String())
	})
	t.Run("should return ok in case of uin32 primitive", func(t *testing.T) {
		fn := func() (uint32, error) {
			return 1, nil
		}
		res, err := utils.NewFuncRepeater(fn).WithMaxTries(3).WithFallbackTries(3).Run()
		require.NoError(t, err)
		require.Equal(t, uint32(1), res)
	})
}
