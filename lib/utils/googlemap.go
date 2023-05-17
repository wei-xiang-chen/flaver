package utils

import (
	"context"
	"flaver/globals"
	"sync"

	"googlemaps.github.io/maps"
)

var (
	googleMapClient     *maps.Client
	googleMapClientMux  sync.Mutex
	googleMapClientOnce sync.Once
)

type IMapUtil interface {
	GetPlaceDetails(placeID string) (*maps.PlaceDetailsResult, error)
}

func newGoogleMapClient() *maps.Client {
	client, err := maps.NewClient(maps.WithAPIKey(globals.GetViper().GetString("google_map.api_key"))); 
	if err != nil {
		globals.GetLogger().Errorf("[google map client startup] error: ", err)
	}
	
	return client
}

func getGoogleMapClient() *maps.Client {
	googleMapClientMux.Lock()
	defer googleMapClientMux.Unlock()
	googleMapClientOnce.Do(func() {
		if googleMapClient == nil {
			googleMapClient = newGoogleMapClient()
		}
	})

	return googleMapClient
}

type GoogleMapUtil struct {
}

func (this *GoogleMapUtil) GetPlaceDetails(placeID string) (*maps.PlaceDetailsResult, error) {
	ctx := context.Background()
	if placeDetails, err := getGoogleMapClient().PlaceDetails(ctx, &maps.PlaceDetailsRequest{PlaceID: placeID}); err != nil {
		return nil, err
	} else {
		return &placeDetails, nil
	}
}
