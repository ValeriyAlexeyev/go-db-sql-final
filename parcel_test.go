package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func getTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	_, err = db.Exec(`
	CREATE TABLE parcel (
		number INTEGER PRIMARY KEY AUTOINCREMENT,
		client INTEGER,
		status TEXT,
		address TEXT,
		created_at TEXT
	);
	`)
	require.NoError(t, err)

	return db
}

func TestAddGetDelete(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	store := NewParcelStore(db)

	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	parcel.Number = id

	stored, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, parcel, stored)

	err = store.Delete(id)
	require.NoError(t, err)

	_, err = store.Get(id)
	require.Error(t, err)
}

func TestSetAddress(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	store := NewParcelStore(db)

	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)

	newAddress := "new test address"

	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	stored, err := store.Get(id)
	require.NoError(t, err)

	require.Equal(t, newAddress, stored.Address)
}

func TestSetStatus(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	store := NewParcelStore(db)

	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)

	err = store.SetStatus(id, ParcelStatusSent)
	require.NoError(t, err)

	stored, err := store.Get(id)
	require.NoError(t, err)

	require.Equal(t, ParcelStatusSent, stored.Status)
}

func TestGetByClient(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	store := NewParcelStore(db)

	client := randRange.Intn(100000)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}

	for i := range parcels {
		parcels[i].Client = client

		id, err := store.Add(parcels[i])
		require.NoError(t, err)

		parcels[i].Number = id
	}

	stored, err := store.GetByClient(client)
	require.NoError(t, err)

	require.Len(t, stored, len(parcels))

	for i := range stored {
		require.Equal(t, parcels[i], stored[i])
	}
}
