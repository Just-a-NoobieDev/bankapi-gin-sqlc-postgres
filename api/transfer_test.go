package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/Just-A-NoobieDev/bankapi-gin-sqlc/db/mock"
	db "github.com/Just-A-NoobieDev/bankapi-gin-sqlc/db/sqlc"
	"github.com/Just-A-NoobieDev/bankapi-gin-sqlc/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func randomTransfer(account1 db.Account, account2 db.Account) db.Transfer {
	return db.Transfer{
		ID:           util.RandomInt(1, 1000),
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Amount:       util.RandomAmount(),
	}
}

func TestCreateTransferAPI(t *testing.T) {

	amount := int64(10)

	account1 := randomAccount()
	account2 := randomAccount()
	account3 := randomAccount()

	account1.Currency = "USD"
	account2.Currency = "USD"
	account3.Currency = "EUR"

	testCases := []struct {
		name        string
		body        string
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: fmt.Sprintf(`{"from_account_id": %d, "to_account_id": %d, "amount": 10, "currency": "USD"}`, account1.ID, account2.ID),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)

				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID: account2.ID,
					Amount: amount,
				}

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(arg)).
					Times(1)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
			},
		},
		{
			name: "FromAccountNotFound",
			body: fmt.Sprintf(`{"from_account_id": %d, "to_account_id": %d, "amount": 10, "currency": "USD"}`, account1.ID, account2.ID),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, rec.Code)
			},
		},
		{
			name: "ToAccountNotFound",
			body: fmt.Sprintf(`{"from_account_id": %d, "to_account_id": %d, "amount": 10, "currency": "USD"}`, account1.ID, account2.ID),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, rec.Code)
			},
		},
		{
			name: "FromAccountCurrencyMismatch",
			body: fmt.Sprintf(`{"from_account_id": %d, "to_account_id": %d, "amount": 10, "currency": "USD"}`, account1.ID, account3.ID),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account3.ID)).Times(1).Return(account3, nil)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "ToAccountCurrencyMismatch",
			body: fmt.Sprintf(`{"from_account_id": %d, "to_account_id": %d, "amount": 10, "currency": "USD"}`, account1.ID, account3.ID),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account3.ID)).Times(1).Return(account3, nil)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "NegativeAmount",
			body: fmt.Sprintf(`{"from_account_id": %d, "to_account_id": %d, "amount": -10, "currency": "USD"}`, account1.ID, account2.ID),
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "Get Account Error",
			body: fmt.Sprintf(`{"from_account_id": %d, "to_account_id": %d, "amount": 10, "currency": "USD"}`, account1.ID, account2.ID),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(1).Return(db.Account{}, sql.ErrConnDone)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			name: "TransferTx Error",
			body: fmt.Sprintf(`{"from_account_id": %d, "to_account_id": %d, "amount": 10, "currency": "USD"}`, account1.ID, account2.ID),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(1).Return(db.TransferTxResult{}, sql.ErrTxDone)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			name: "InvalidCurrency",
			body: fmt.Sprintf(`{"from_account_id": %d, "to_account_id": %d, "amount": 10, "currency": "AAA"}`, account1.ID, account2.ID),
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			rec := httptest.NewRecorder()

			url := "/api/v1/transfers"
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(tc.body))
			require.NoError(t, err)

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(t, rec)
		})
	}
}

func TestGetTransfersByAccountAPI(t *testing.T) {
	account1 := randomAccount()
	account2 := randomAccount()

	n := 5
	transfers := make([]db.Transfer, n)
	for i := 0; i < n; i++ {
		transfers[i] = randomTransfer(account1, account2)
	}

	type Query struct {
		Id int64
		Page int
		Size int
	}

	testCases := []struct {
		name string
		query Query
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				Id: account1.ID,
				Page: 1,
				Size: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetTransfersByAccountParams{
					ID: account1.ID,
					Off: 0,
					Size: int32(n),
				}

				store.EXPECT().
					GetTransfersByAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(transfers, nil)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
				require.Len(t, transfers, n)
				requireBodyMatchTransfers(t, rec.Body, transfers)
			},
		},
		{
			name: "InternalError",
			query: Query{
				Id: account1.ID,
				Page: 1,
				Size: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTransfersByAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Transfer{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			name: "InvalidId",
			query: Query{
				Id: 0,
				Page: 1,
				Size: n,
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "InvalidPage",
			query: Query{
				Id: account1.ID,
				Page: 0,
				Size: n,
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "InvalidSize",
			query: Query{
				Id: account1.ID,
				Page: 1,
				Size: 0,
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			rec := httptest.NewRecorder()

			url := fmt.Sprintf("/api/v1/transfers?id=%d&page=%d&size=%d", tc.query.Id, tc.query.Page, tc.query.Size)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(t, rec)
		})
	}
}

func TestGetTransferByIdAPI(t *testing.T) {
	account1 := randomAccount()
	account2 := randomAccount()

	transfer := randomTransfer(account1, account2)

	testCases := []struct {
		name string
		id int64
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			id: transfer.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTransfer(gomock.Any(), gomock.Eq(transfer.ID)).
					Times(1).
					Return(transfer, nil)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
				requireBodyMatchTransfer(t, rec.Body, transfer)
			},
		},
		{
			name: "InternalError",
			id: transfer.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTransfer(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Transfer{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			name: "NotFound",
			id: transfer.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTransfer(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Transfer{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, rec.Code)
			},
		},
		{
			name: "InvalidId",
			id: 0,
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			rec := httptest.NewRecorder()

			url := fmt.Sprintf("/api/v1/transfers/%d", tc.id)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(t, rec)
		})
	}
}

func requireBodyMatchTransfer(t *testing.T, body *bytes.Buffer, transfer db.Transfer) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTransfer db.Transfer
	err = json.Unmarshal(data, &gotTransfer)
	require.NoError(t, err)
	require.Equal(t, transfer, gotTransfer)
}

func requireBodyMatchTransfers(t *testing.T, body *bytes.Buffer, transfers []db.Transfer) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTransfers []db.Transfer
	err = json.Unmarshal(data, &gotTransfers)
	require.NoError(t, err)
	require.Equal(t, transfers, gotTransfers)
}