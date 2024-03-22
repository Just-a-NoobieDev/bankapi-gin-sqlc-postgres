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

func createRandomEntry(account db.Account) db.Entry {
	return db.Entry{
		ID: util.RandomInt(1, 1000),
		AccountID: account.ID,
		Amount: util.RandomAmount(),
	}
}

func TestGetEntriesByAccountAPI(t *testing.T) {
	
	n := 10
	account := randomAccount()

	entries := make([]db.Entry, n)
	for i := 0; i < n; i++ {
		entries[i] = createRandomEntry(account)
	}

	type Query struct {
		Id int64
		Page int32
		Size int32
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
				Id: account.ID,
				Page: 1,
				Size: int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetEntriesParams{
					AccountID: account.ID,
					Limit: int32(n),
					Offset: 0,
				}

				store.EXPECT().
					GetEntries(gomock.Any(), arg).
					Times(1).
					Return(entries, nil)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
				requireBodyMatchEntries(t, rec.Body, entries)
			},
		},
		{
			name: "InternalError",
			query: Query{
				Id: account.ID,
				Page: 1,
				Size: int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetEntriesParams{
					AccountID: account.ID,
					Limit: int32(n),
					Offset: 0,
				}

				store.EXPECT().
					GetEntries(gomock.Any(), arg).
					Times(1).
					Return([]db.Entry{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			name: "InvalidID",
			query: Query{
				Id: 0,
				Page: 1,
				Size: int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "InvalidPage",
			query: Query{
				Id: account.ID,
				Page: 0,
				Size: int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "InvalidSize",
			query: Query{
				Id: account.ID,
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

			url := fmt.Sprintf("/api/v1/entry?id=%d&page=%d&size=%d", tc.query.Id, tc.query.Page, tc.query.Size)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(t, rec)
		})
	}

}


func requireBodyMatchEntries(t *testing.T, body *bytes.Buffer, entries []db.Entry) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotEntries []db.Entry
	err = json.Unmarshal(data, &gotEntries)
	require.NoError(t, err)
	require.Equal(t, entries, gotEntries)
}