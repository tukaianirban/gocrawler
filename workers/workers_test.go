package workers

import (
	"testing"
	"runtime"
)

func TestWorkerPool_GetWorkerToken(t *testing.T) {

	maxWorkerCount := runtime.NumCPU()
	testPool := NewWorkerPool(1, maxWorkerCount)

	for i:=0; i<maxWorkerCount; i++ {
		t.Logf("received worker token: %d", testPool.GetWorkerToken())
	}
	t.Logf("all tokens drained from pool")

	t.Logf("next token request: %d", testPool.GetWorkerToken())
	t.Logf("returning more tokens than obtained ...")
	for i:=0; i<maxWorkerCount + 2; i++ {
		testPool.ReturnWorkerToken()
	}

	t.Logf("drain all tokens again...")
	for i:=0; i<maxWorkerCount; i++ {
		t.Logf("received worker token: %d", testPool.GetWorkerToken())
	}
	t.Logf("all tokens drained from pool")
}