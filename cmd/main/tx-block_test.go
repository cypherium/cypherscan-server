package main_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func (m *MockedRepo) GetBlock(number int64) (*repo.TxBlock, error) {
	args := m.Called(number)
	return args.Get(0).(*repo.TxBlock), args.Error(1)
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

func (m *MockedBlocksFetcher) GetLatestBlockNumber() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
func (m *MockedBlocksFetcher) GetLatestKeyBlockNumber() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
func TestGetTxBlocksWithoutAnyInDb(t *testing.T) {
	req, _ := http.NewRequest("GET", "/tx-blocks?p=1&pagesize=3", nil)
	mockedRepo := new(MockedRepo)
	mockedRepo.On("GetBlocks", &repo.BlockSearchContdition{Scenario: 1, StartWith: 11, PageSize: 3}).Return([]repo.TxBlock{}, nil)
	mockedWsServer := new(MockedWebSocketServer)
	mockedBlocksFetcher := new(MockedBlocksFetcher)
	mockedBlocksFetcher.On("GetLatestBlockNumber").Return(int64(11), nil)
	mockedBlocksFetcher.On("BlockHeadersByNumbers", []int64{11, 10, 9}).Return([]*types.Header{
		&types.Header{Number: big.NewInt(9), Time: big.NewInt(time.Now().Unix())},
		&types.Header{Number: big.NewInt(11), Time: big.NewInt(time.Now().Unix())},
		&types.Header{Number: big.NewInt(10), Time: big.NewInt(time.Now().Unix())},
	}, nil)
	app := main.NewApp(mockedRepo, mockedWsServer, mockedBlocksFetcher, "")
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	mockedBlocksFetcher.AssertNotCalled(t, "BlockHeadersByNumbers")
	assert.Equal(t, rr.Code, http.StatusOK)

	var m main.ResponseOfGetBlocks
	json.Unmarshal(rr.Body.Bytes(), &m)
	assert.Equal(t, 3, len(m.Blocks))
	assert.Equal(t, int64(11), m.Blocks[0].Number)
	assert.Equal(t, int64(10), m.Blocks[1].Number)
	assert.Equal(t, int64(9), m.Blocks[2].Number)
}

func TestGetTxBlocksWithFirtPageAllInDb(t *testing.T) {
	mockedRepo := new(MockedRepo)
	mockedRepo.On("GetBlocks", &repo.BlockSearchContdition{Scenario: 1, StartWith: 12, PageSize: 3}).Return([]repo.TxBlock{
		{Number: 12},
		{Number: 11},
		{Number: 10},
	}, nil)
	mockedWsServer := new(MockedWebSocketServer)
	mockedBlocksFetcher := new(MockedBlocksFetcher)
	mockedBlocksFetcher.On("GetLatestBlockNumber").Return(int64(12), nil)
	mockedBlocksFetcher.On("BlockHeadersByNumbers").Return([]*types.Header{}, nil)
	app := main.NewApp(mockedRepo, mockedWsServer, mockedBlocksFetcher, "")

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tx-blocks?p=1&pagesize=3", nil)
	app.Router.ServeHTTP(rr, req)
	mockedBlocksFetcher.AssertNotCalled(t, "BlockHeadersByNumbers") // already got from db, no need to call blockchain

	assert.Equal(t, rr.Code, http.StatusOK)

	var m main.ResponseOfGetBlocks
	json.Unmarshal(rr.Body.Bytes(), &m)
	assert.Equal(t, int64(12), m.Total)
	assert.Equal(t, int64(12), m.Blocks[0].Number)
	assert.Equal(t, int64(11), m.Blocks[1].Number)
	assert.Equal(t, int64(10), m.Blocks[2].Number)
}

func TestGetTxBlocksFirstPageWithSomeInDb(t *testing.T) {
	req, _ := http.NewRequest("GET", "/tx-blocks?p=1&pagesize=3", nil)
	mockedRepo := new(MockedRepo)
	mockedRepo.On("GetBlocks", &repo.BlockSearchContdition{Scenario: 1, StartWith: 12, PageSize: 3}).Return([]repo.TxBlock{
		{Number: 12},
		{Number: 10},
	}, nil)

	mockedWsServer := new(MockedWebSocketServer)

	mockedBlocksFetcher := new(MockedBlocksFetcher)
	mockedBlocksFetcher.On("GetLatestBlockNumber").Return(int64(12), nil)
	mockedBlocksFetcher.On("BlockHeadersByNumbers", []int64{11}).Return([]*types.Header{
		&types.Header{Number: big.NewInt(11), Time: big.NewInt(time.Now().Unix())},
	}, nil)

	app := main.NewApp(mockedRepo, mockedWsServer, mockedBlocksFetcher, "")
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)

	var m main.ResponseOfGetBlocks
	json.Unmarshal(rr.Body.Bytes(), &m)
	assert.Equal(t, 3, len(m.Blocks))
	assert.Equal(t, int64(12), m.Blocks[0].Number)
	assert.Equal(t, int64(11), m.Blocks[1].Number)
	assert.Equal(t, int64(10), m.Blocks[2].Number)
}
func TestGetTxBlocksSecondPageWithSomeInDb(t *testing.T) {
	req, _ := http.NewRequest("GET", "/tx-blocks?p=2&pagesize=3", nil)
	mockedRepo := new(MockedRepo)
	mockedRepo.On("GetBlocks", &repo.BlockSearchContdition{Scenario: 1, StartWith: 9, PageSize: 3}).Return([]repo.TxBlock{
		{Number: 8},
	}, nil)

	mockedWsServer := new(MockedWebSocketServer)

	mockedBlocksFetcher := new(MockedBlocksFetcher)
	mockedBlocksFetcher.On("GetLatestBlockNumber").Return(int64(12), nil)
	mockedBlocksFetcher.On("BlockHeadersByNumbers", []int64{9, 7}).Return([]*types.Header{
		&types.Header{Number: big.NewInt(9), Time: big.NewInt(time.Now().Unix())},
		&types.Header{Number: big.NewInt(7), Time: big.NewInt(time.Now().Unix())},
	}, nil)

	app := main.NewApp(mockedRepo, mockedWsServer, mockedBlocksFetcher, "")
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)

	var m main.ResponseOfGetBlocks
	json.Unmarshal(rr.Body.Bytes(), &m)
	assert.Equal(t, 3, len(m.Blocks))
	assert.Equal(t, int64(9), m.Blocks[0].Number)
	assert.Equal(t, int64(8), m.Blocks[1].Number)
	assert.Equal(t, int64(7), m.Blocks[2].Number)
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
