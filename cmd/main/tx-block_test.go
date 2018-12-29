package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cypherium/CypherTestNet/go-cypherium/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/ron-liu/cypherscan-server/cmd/main"
	"gitlab.com/ron-liu/cypherscan-server/internal/repo"
)

type MockedRepo struct {
	mock.Mock
}

func (m *MockedRepo) GetBlocks(condition *repo.BlockSearchContdition) ([]repo.TxBlock, error) {
	args := m.Called(condition)
	return args.Get(0).([]repo.TxBlock), args.Error(1)
}
func (m *MockedRepo) GetKeyBlocks(condition *repo.BlockSearchContdition) ([]repo.KeyBlock, error) {
	args := m.Called(condition)
	return args.Get(0).([]repo.KeyBlock), args.Error(1)
}

func (m *MockedRepo) GetTransactions(condition *repo.TransactionSearchCondition) ([]repo.Transaction, error) {
	args := m.Called(condition)
	return args.Get(0).([]repo.Transaction), args.Error(1)
}

type MockedWebSocketServer struct {
	mock.Mock
}

func (m *MockedWebSocketServer) ServeWebsocket(w http.ResponseWriter, r *http.Request) {
	args := m.Called(w, r)
	fmt.Printf("%v", args)
}

type MockedBlocksFetcher struct {
	mock.Mock
}

func (m *MockedBlocksFetcher) BlockHeadersByNumbers(numbers []int64) ([]*types.Header, error) {
	args := m.Called(numbers)
	return args.Get(0).([]*types.Header), args.Error(1)
}
func (m *MockedBlocksFetcher) KeyBlocksByNumbers(numbers []int64) ([]*types.KeyBlock, error) {
	args := m.Called(numbers)
	return args.Get(0).([]*types.KeyBlock), args.Error(1)
}

func TestGetTxBlocks(t *testing.T) {
	req, _ := http.NewRequest("GET", "/tx-blocks/11?pagesize=5", nil)
	mockedRepo := new(MockedRepo)
	mockedRepo.On("GetBlocks", &repo.BlockSearchContdition{Scenario: 1, StartWith: 11, PageSize: 5}).Return([]repo.TxBlock{}, nil)

	mockedWsServer := new(MockedWebSocketServer)

	mockedBlocksFetcher := new(MockedBlocksFetcher)
	mockedBlocksFetcher.On("BlockHeadersByNumbers", mock.Anything).Return([]*types.Header{}, nil)

	app := main.NewApp(mockedRepo, mockedWsServer, mockedBlocksFetcher, "")

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)
	mockedBlocksFetcher.AssertCalled(t, "BlockHeadersByNumbers", []int64{11, 10, 9, 8, 7})

}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
