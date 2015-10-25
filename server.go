package main

import (
    "fmt"
    "time"
    "errors"
    "strings"
    "io/ioutil"
    "net/http"
    "encoding/json"
    "gopkg.in/mgo.v2"
    "math/rand"
    "github.com/jmoiron/jsonq"
    "gopkg.in/mgo.v2/bson"
    "github.com/julienschmidt/httprouter"
)

type UserData struct {
    Id bson.ObjectId `json:"id" bson:"_id"`
    Name string `json:"name" bson:"name"`
    Address string `json:"address" bson:"address"`
    City string `json:"city" bson:"city"`
    State string `json:"state" bson:"state"`
    Zip string `json:"zip" bson:"zip"`
    Coordinate struct {
        Lat float64 `json:"lat" bson:"lat"`
        Lng float64 `json:"lng" bson:"lng"`
    } `json:"coordinate" bson:"coordinate"`
}

func getSession() *mgo.Session {
    //Connect to local mongo
    s, err := mgo.Dial("mongodb://kamlendrachauhan:cmpe273@ds045064.mongolab.com:45064/location_service")

    // Check if connection error, is mongo running?
    if err != nil {
        panic(err)
    }
    return s
}
//Get a Location - GET        /locations/{location_id}
func getLocations(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    location_id :=  p.ByName("location_id")

    if !bson.IsObjectIdHex(location_id) {
        rw.WriteHeader(404)
        return
    }

    original_loc_id := bson.ObjectIdHex(location_id)

    returnObj := UserData{}

    if err := getSession().DB("location_service").C("location").FindId(original_loc_id).One(&returnObj); err != nil {
        rw.WriteHeader(404)
        return
    }

    uj, _ := json.Marshal(returnObj)

    // Write content-type, statuscode, payload
    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(200)
    fmt.Fprintf(rw, "%s", uj)
}

//Create New Location - POST        /locations
func saveLocations(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    var u UserData
    URL := "http://maps.google.com/maps/api/geocode/json?address="
    //Populate the data in local object
    json.NewDecoder(req.Body).Decode(&u)

    //Randomly generated unique ID
   // u.Id = randomString(10)
    u.Id = bson.NewObjectId()

    URL = URL +u.Address+ " " + u.City + " " + u.State + " " + u.Zip+"&sensor=false"
    URL = strings.Replace(URL, " ", "+", -1)
    fmt.Println("URL "+ URL)

    //calling google map API
    response, err := http.Get(URL)
    if err != nil {
        return
    }
    defer response.Body.Close()

    resp := make(map[string]interface{})
    body, _ := ioutil.ReadAll(response.Body)
    err = json.Unmarshal(body, &resp)
    if err != nil {
        return
    }

    jq := jsonq.NewQuery(resp)
    status, err := jq.String("status")
    fmt.Println(status)
    if err != nil {
        return
    }
    if status != "OK" {
        err = errors.New(status)
        return
    }

    latitude, err := jq.Float("results" ,"0","geometry", "location", "lat")
   if err != nil {
       fmt.Println(err)
        return
    }
    longitude, err := jq.Float("results", "0","geometry", "location", "lng")
    if err != nil {
        fmt.Println(err)
        return
    }

    u.Coordinate.Lat = latitude
    u.Coordinate.Lng = longitude

    //Persisting Data
    getSession().DB("location_service").C("location").Insert(u)


    // Marshal provided interface into JSON structure
    uj, _ := json.Marshal(u)

    // Write content-type, status code, payload
    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(201)
    fmt.Fprintf(rw, "%s", uj)

}

//Create New Location - POST        /locations
func updateLocations(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    var u UserData
    location_id :=  p.ByName("location_id")

    URL := "http://maps.google.com/maps/api/geocode/json?address="

    //Populate the data in local object
    json.NewDecoder(req.Body).Decode(&u)

    URL = URL +u.Address+ " " + u.City + " " + u.State + " " + u.Zip+"&sensor=false"
    URL = strings.Replace(URL, " ", "+", -1)
    fmt.Println("URL "+ URL)

    //calling google map API
    response, err := http.Get(URL)
    if err != nil {
        return
    }
    defer response.Body.Close()

    resp := make(map[string]interface{})
    body, _ := ioutil.ReadAll(response.Body)
    err = json.Unmarshal(body, &resp)
    if err != nil {
        return
    }

    jq := jsonq.NewQuery(resp)
    status, err := jq.String("status")
    fmt.Println(status)
    if err != nil {
        return
    }
    if status != "OK" {
        err = errors.New(status)
        return
    }

    latitude, err := jq.Float("results" ,"0","geometry", "location", "lat")
    if err != nil {
        fmt.Println(err)
        return
    }
    longitude, err := jq.Float("results", "0","geometry", "location", "lng")
    if err != nil {
        fmt.Println(err)
        return
    }

    u.Coordinate.Lat = latitude
    u.Coordinate.Lng = longitude

    original_loc_id := bson.ObjectIdHex(location_id)
    var data = UserData{
        Address: u.Address,
        City: u.City,
        State: u.State,
        Zip: u.Zip,
    }
    //updateData := bson.M{ "$set": data}
    fmt.Println(data)
    //Persisting Data
    getSession().DB("location_service").C("location").Update(bson.M{"_id":original_loc_id }, bson.M{"$set": bson.M{ "address": u.Address,
        "city": u.City, "state": u.State,"zip": u.Zip, "coordinate.lat":u.Coordinate.Lat, "coordinate.lng":u.Coordinate.Lng}})

    returnObj := UserData{}

    //fetch the response data
    if err := getSession().DB("location_service").C("location").FindId(original_loc_id).One(&returnObj); err != nil {
        rw.WriteHeader(404)
        return
    }
    // Marshal provided interface into JSON structure
    uj, _ := json.Marshal(returnObj)

    // Write content-type, status code, payload
    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(201)
    fmt.Fprintf(rw, "%s", uj)

}

//Delete a Location - DELETE /locations/{location_id}
func removeLocations(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    location_id :=  p.ByName("location_id")

    if !bson.IsObjectIdHex(location_id) {
        rw.WriteHeader(404)
        return
    }

    original_loc_id := bson.ObjectIdHex(location_id)

    // Remove user
    if err := getSession().DB("location_service").C("location").RemoveId(original_loc_id); err != nil {
        rw.WriteHeader(404)
        return
    }

    rw.WriteHeader(200)
}

func main() {
    mux := httprouter.New()
    mux.GET("/locations/:location_id", getLocations)
    mux.POST("/locations", saveLocations)
    mux.PUT("/locations/:location_id", updateLocations)
    mux.DELETE("/locations/:location_id", removeLocations)
    rand.Seed( time.Now().UTC().UnixNano())

    server := http.Server{
        Addr:        "0.0.0.0:8880",
        Handler: mux,
    }
    server.ListenAndServe()
}