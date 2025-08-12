package main

type Position struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type XYZData struct {
	Tags []string `json:"tags"`
}

/*
func (app *application) addressSearchHandler(w http.ResponseWriter, r *http.Request) {
	headers := make(http.Header)

	var input struct {
		Search  string `json:"search"`
		Viewbox string `json:"viewbox"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	res, errors := extended.SearchOsm(input.Search, input.Viewbox)

	featureCollection := geojson.NewFeatureCollection()

	for key := range res {

		lat64, _ := strconv.ParseFloat(res[key].Lat, 64)
		lon64, _ := strconv.ParseFloat(res[key].Lon, 64)

		pos := Position{lat64, lon64}

		geo := app.fillGeoJSON(strconv.FormatInt(int64(res[key].OsmID), 10), "loc", pos, envelope{"place_id": strconv.FormatInt(int64(res[key].PlaceID), 10), "type": res[key].Type, "osm_type": res[key].OsmType, "display": res[key].DisplayName, "importance": res[key].Importance, "address": res[key].Address, "extratags": res[key].Extratags, "boundingbox": res[key].Boundingbox, "svg": res[key].Svg})
		featureCollection.AddFeature(geo)

	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"query": input.Search, "results": featureCollection, "errors": errors}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

*/
