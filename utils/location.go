package utils

import (
	"errors"
	"xg/entity"
)

var (
	ErrLocationNotFound    = errors.New("location not found")
	ErrLocationInfoInvalid = errors.New("location info invalid")
)

type LocationResponse struct {
	Status   string             `json:"status"`
	Info     string             `json:"info"`
	Count    string             `json:"count"`
	GeoCodes []*LocationGeocode `json:"geocodes"`
}
type LocationGeocode struct {
	FormattedAddress string `json:"formatted_address"`
	Country          string `json:"country"`
	Province         string `json:"province"`
	City             string `json:"city"`
	District         string `json:"district"`
	AdCode           string `json:"adcode"`
	Location         string `json:"location"`
	Level            string `json:"level"`
}

func GetAddressLocation(addr string) (*entity.Coordinate, error) {
	// locationKey := conf.Get().AMapKey
	// path := fmt.Sprintf("https://restapi.amap.com/v3/geocode/geo?address=%v&output=JSON&key=%v", addr, locationKey)
	// resp, err := http.Get(path)
	// if err != nil {
	// 	return nil, err
	// }
	// defer resp.Body.Close()
	// resBody, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return nil, err
	// }

	// locationResp := new(LocationResponse)
	// err = json.Unmarshal(resBody, locationResp)
	// if err != nil {
	// 	return nil, err
	// }

	// if locationResp.Status != "1" {
	// 	return nil, errors.New(locationResp.Info)
	// }
	// count := ParseInt(locationResp.Count)
	// if count < 1 || len(locationResp.GeoCodes) < 1 {
	// 	return nil, ErrLocationNotFound
	// }

	// ret := ParseFloats(locationResp.GeoCodes[0].Location)
	// if len(ret) < 2 {
	// 	return nil, ErrLocationInfoInvalid
	// }

	// log.Info.Printf("Get address: %#v, resp: %#v", addr, locationResp.GeoCodes[0])
	// return &entity.Coordinate{
	// 	Longitude: ret[0],
	// 	Latitude:  ret[1],
	// }, nil
	return &entity.Coordinate{Longitude: 0, Latitude: 0}, nil
}
