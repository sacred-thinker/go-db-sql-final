package main

import (
	"database/sql"

	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestAddGetDelete(t *testing.T) {

	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	parcel.Number = id

	storedParcel, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, parcel.Number, storedParcel.Number, "Parcel Number does not match")
	assert.Equal(t, parcel.Client, storedParcel.Client, "Parcel Client does not match")
	assert.Equal(t, parcel.Status, storedParcel.Status, "Parcel Status does not match")
	assert.Equal(t, parcel.Address, storedParcel.Address, "Parcel Address does not match")
	assert.Equal(t, parcel.CreatedAt, storedParcel.CreatedAt, "Parcel CreatedAt does not match")

	// storedParcelById, err := store.Get(id)
	// require.NoError(t, err)

	//assert.False(t, reflect.DeepEqual(parcel, storedParcelById), "Stored parcel by ID does not match the original parcel")

	err = store.Delete(parcel.Number)
	require.NoError(t, err)

	_, err = store.Get(parcel.Number)
	require.Error(t, err)

	assert.ErrorIs(t, err, sql.ErrNoRows, "Expected error to be sql.ErrNoRows")

}

func TestSetAddress(t *testing.T) {

	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	assert.NotEmpty(t, id)

	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	storedParcel, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, newAddress, storedParcel.Address)
}

func TestSetStatus(t *testing.T) {

	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	assert.NotEmpty(t, id)

	err = store.SetStatus(id, ParcelStatusDelivered)
	require.NoError(t, err)

	storedParcel, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, ParcelStatusDelivered, storedParcel.Status)
}

func TestGetByClient(t *testing.T) {

	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		assert.NotEmpty(t, id)

		parcels[i].Number = id

		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	assert.Len(t, storedParcels, len(parcelMap))
	// assert.Equal(t, len(parcelMap), len(storedParcels))

	for _, parcel := range storedParcels {
		_, ok := parcelMap[parcel.Number]
		assert.True(t, ok)
		assert.Equal(t, parcel, parcelMap[parcel.Number])
	}
}
