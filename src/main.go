//Start redis server using redis-server redis.conf
//run the go file main.go
package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"math"
	"io/ioutil"
	"log"
	"github.com/gorilla/mux"
	"html/template"
	"os"
)

//	"github.com/go-redis/redis"

type City struct {
 Name string
 NoOfParkingSlots int64
 PicURL string
 Parking []Slot
 Longitude string
 Latitude string
}

type Slot struct {
 Name string
 DistanceFromYou float64
 ETA float64
 PicURL string
}

type MapBoxResp struct {
	Code   string `json:"code"`
	Routes []RoutesResp
}

type RoutesResp struct {
	Distance float64 `json:"distance"`
	Duration float64 `json:"duration"`
}

var city *City


//var client *redis.Client

var templates *template.Template
//var ctx = context.TODO()

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == ""{
	  port = "8000"
	}
	return ":" + port, nil
  }

func main(){
/*	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB: 0,
	})*/

	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}

	templates = template.Must(template.ParseGlob("templates/*.gohtml"))
	r := mux.NewRouter()
	r.HandleFunc("/", index_Show_City_Handler).Methods("GET")
	r.HandleFunc("/bhopal", index_Get_City_Info_Handler).Methods("POST")
	r.HandleFunc("/bangalore", index_Get_City_Info_Handler).Methods("POST")
	r.HandleFunc("/mumbai", index_Get_City_Info_Handler).Methods("POST")
	r.HandleFunc("/indore", index_Get_City_Info_Handler).Methods("POST")
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",fs))
	http.Handle("/",r)
	//http.ListenAndServe(":8080",nil)
	if err := http.ListenAndServe(addr, nil);
	err != nil {
		panic(err)
	  }
}

func index_Show_City_Handler (w http.ResponseWriter, r *http.Request) {
	/*cities, err := client.LRange(ctx,"cities", 0, 10).Result()
	if err != nil {
		return
	}*/
	templates.ExecuteTemplate(w, "index.gohtml", nil)
}

func index_Get_City_Info_Handler (w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("path", r.URL.Path)
	userLat := r.PostForm.Get("latitude")
	userLong := r.PostForm.Get("longitude")
	fmt.Println(userLat)
	fmt.Println(userLong)

	//fetch data from dbms and populate the city struct based on the path url

	//fetch all the parkings from the database and find the distance from currLat and Long
    
	
	if(r.URL.Path=="/bhopal"){
	
	//finding distance and ETA for all the parking lots in city and creating the struct
	var NoOfParkings int = 6
	latLongOfParkings := make([]string, NoOfParkings)
	latLongOfParkings =	[]string{"77.458059, 23.231247","77.338690, 23.270774","77.412542, 23.267715","77.407086, 23.263625","77.456512, 23.268870","77.436565,23.214282"}
	
	nameOfParkings := make([]string, NoOfParkings)
	nameOfParkings = []string{"Smart Parking Bhopal","New Market Parking Lot","Parking T.T Nagar","Bhopal Juntion Parking","Parking Sultania Road","Arera Colony Car Parking"}
	
	picURL := make([]string, NoOfParkings)
	picURL = []string{"/static/images/bhopalps1.jpg","/static/images/bhopalps2.jpg","/static/images/bhopalps3.jpg","/static/images/bhopalps4.jpg","/static/images/bhopalps5.jpg","/static/images/bhopalps5.jpg"}
	
	distanceFromUser := make([]float64, NoOfParkings)
	ETAOfParkings := make([]float64, NoOfParkings)

	city = &City{
		Name: "Bhopal",
		NoOfParkingSlots: int64(NoOfParkings), 
		PicURL: "/static/images/bhopal.jpg",
		Longitude: "77.408625",
		Latitude: "23.258660",
	}

	var arrOfSlots []Slot
	for i := 0; i < NoOfParkings; i++ {
			var result MapBoxResp
			resp, err := http.Get("https://api.mapbox.com/directions/v5/mapbox/driving/"+userLong+","+userLat+";"+latLongOfParkings[i]+"?access_token=pk.eyJ1Ijoic2hhc2hhbmstcHJha2FzaCIsImEiOiJja2ZjZzJmOGsxN3kxMnRvNXp1MDJxcGFzIn0.69B3TiOom1WXyot_FdTVHA")
			if err != nil{
				log.Fatalln(err)
			}

			defer resp.Body.Close()
			data, _ := ioutil.ReadAll(resp.Body)
			json.Unmarshal(data,&result)
			dis :=  math.Round(((result.Routes[0].Distance)/1000)*100)/100
			dur :=  math.Round(((result.Routes[0].Duration)/3600)*100)/100
			distanceFromUser[i] = dis
			ETAOfParkings[i] = dur
			slot:=Slot{
				Name: nameOfParkings[i], 
				DistanceFromYou: distanceFromUser[i],
				ETA: ETAOfParkings[i],
				PicURL: picURL[i],
			}
			arrOfSlots = append(arrOfSlots,slot)
	}
	city.Parking = arrOfSlots
		/*city = &City{
			Name: "Bhopal",
			NoOfParkingSlots: 10, 
			PicURL: "/static/images/bhopal.jpg",
			Parking: []Slot{
				Slot{
					Name: nameOfParkings[0], 
                    DistanceFromYou: distanceFromUser[0],
                    ETA: ETAOfParkings[0],
                    PicURL: picURL[0],
				},
				Slot{
					Name: nameOfParkings[1], 
					DistanceFromYou: distanceFromUser[1],
                    ETA: ETAOfParkings[1],
                    PicURL: picURL[1],
				},
				Slot{
					Name: nameOfParkings[2], 
					DistanceFromYou: distanceFromUser[2],
                    ETA: ETAOfParkings[2],
                    PicURL: picURL[2],
				},
				Slot{
					Name: nameOfParkings[3], 
					DistanceFromYou: distanceFromUser[3],
                    ETA: ETAOfParkings[3],
                    PicURL: picURL[3],
				},
				Slot{
					Name: nameOfParkings[4], 
					DistanceFromYou: distanceFromUser[4],
                    ETA: ETAOfParkings[4],
                    PicURL: picURL[4],
				},
			},
			Longitude: "77.408625",
			Latitude: "23.258660",
		}*/
		fmt.Println("Bhopal is reached")
	}

	if(r.URL.Path=="/bangalore"){
	//finding distance and ETA for all the parking lots in city and creating the struct
	var NoOfParkings int = 8
	latLongOfParkings := make([]string, NoOfParkings)
	latLongOfParkings =	[]string{"77.595903,13.051263","77.565560,13.031641","77.570756,13.014827","77.608237,12.982221","77.567669,12.909156","77.577522,12.894322","77.486266,13.085690","77.677213,12.846296"}
	
	nameOfParkings := make([]string, NoOfParkings)
	nameOfParkings = []string{"Tata Global Car Park","MSRIT Parking","Javanica Marg Parking","Commercial Street Parking","Student Parking Lot","B Block Visitors Parking","Acharya Car Parking","GHS Bike Parking Lot"}
	
	picURL := make([]string, NoOfParkings)
	picURL = []string{"/static/images/bhopalps1.jpg","/static/images/bhopalps2.jpg","/static/images/bhopalps3.jpg","/static/images/bhopalps4.jpg","/static/images/bhopalps5.jpg","/static/images/bhopalps5.jpg","/static/images/bhopalps5.jpg","/static/images/bhopalps5.jpg"}
	
	distanceFromUser := make([]float64, NoOfParkings)
	ETAOfParkings := make([]float64, NoOfParkings)

	city = &City{
		Name: "Bangalore",
		NoOfParkingSlots: int64(NoOfParkings), 
		PicURL: "/static/images/bangalore.jpg",
		Longitude: "77.580643", 
		Latitude: "12.972442",
	}

	var arrOfSlots []Slot
	for i := 0; i < NoOfParkings; i++ {
			var result MapBoxResp
			resp, err := http.Get("https://api.mapbox.com/directions/v5/mapbox/driving/"+userLong+","+userLat+";"+latLongOfParkings[i]+"?access_token=pk.eyJ1Ijoic2hhc2hhbmstcHJha2FzaCIsImEiOiJja2ZjZzJmOGsxN3kxMnRvNXp1MDJxcGFzIn0.69B3TiOom1WXyot_FdTVHA")
			if err != nil{
				log.Fatalln(err)
			}

			defer resp.Body.Close()
			data, _ := ioutil.ReadAll(resp.Body)
			json.Unmarshal(data,&result)
			dis :=  math.Round(((result.Routes[0].Distance)/1000)*100)/100
			dur :=  math.Round(((result.Routes[0].Duration)/3600)*100)/100
			distanceFromUser[i] = dis
			ETAOfParkings[i] = dur
			slot:=Slot{
				Name: nameOfParkings[i], 
				DistanceFromYou: distanceFromUser[i],
				ETA: ETAOfParkings[i],
				PicURL: picURL[i],
			}
			arrOfSlots = append(arrOfSlots,slot)
	}
	city.Parking = arrOfSlots
		fmt.Println("Bangalore is reached")
	}

	if(r.URL.Path=="/mumbai"){
	//finding distance and ETA for all the parking lots in city and creating the struct
	var NoOfParkings int = 6
	latLongOfParkings := make([]string, NoOfParkings)
	latLongOfParkings =	[]string{"77.458059, 23.231247","77.338690, 23.270774","77.412542, 23.267715","77.407086, 23.263625","77.456512, 23.268870","77.436565,23.214282"}
	
	nameOfParkings := make([]string, NoOfParkings)
	nameOfParkings = []string{"Smart Parking Bhopal","New Market Parking Lot","Parking T.T Nagar","Bhopal Juntion Parking","Parking Sultania Road","Arera Colony Car Parking"}
	
	picURL := make([]string, NoOfParkings)
	picURL = []string{"/static/images/bhopalps1.jpg","/static/images/bhopalps2.jpg","/static/images/bhopalps3.jpg","/static/images/bhopalps4.jpg","/static/images/bhopalps5.jpg","/static/images/bhopalps5.jpg"}
	
	distanceFromUser := make([]float64, NoOfParkings)
	ETAOfParkings := make([]float64, NoOfParkings)

	city = &City{
		Name: "Mumbai",
		NoOfParkingSlots: int64(NoOfParkings), 
		PicURL: "/static/images/mumbai.jpg",
		Longitude: "72.877426",
        Latitude: "19.076090",
	}

	var arrOfSlots []Slot
	for i := 0; i < NoOfParkings; i++ {
			var result MapBoxResp
			resp, err := http.Get("https://api.mapbox.com/directions/v5/mapbox/driving/"+userLong+","+userLat+";"+latLongOfParkings[i]+"?access_token=pk.eyJ1Ijoic2hhc2hhbmstcHJha2FzaCIsImEiOiJja2ZjZzJmOGsxN3kxMnRvNXp1MDJxcGFzIn0.69B3TiOom1WXyot_FdTVHA")
			if err != nil{
				log.Fatalln(err)
			}

			defer resp.Body.Close()
			data, _ := ioutil.ReadAll(resp.Body)
			json.Unmarshal(data,&result)
			dis :=  math.Round(((result.Routes[0].Distance)/1000)*100)/100
			dur :=  math.Round(((result.Routes[0].Duration)/3600)*100)/100
			distanceFromUser[i] = dis
			ETAOfParkings[i] = dur
			slot:=Slot{
				Name: nameOfParkings[i], 
				DistanceFromYou: distanceFromUser[i],
				ETA: ETAOfParkings[i],
				PicURL: picURL[i],
			}
			arrOfSlots = append(arrOfSlots,slot)
	}
	city.Parking = arrOfSlots

		fmt.Println("Mumbai is reached")
	}

	if(r.URL.Path=="/indore"){
			//finding distance and ETA for all the parking lots in city and creating the struct
	var NoOfParkings int = 6
	latLongOfParkings := make([]string, NoOfParkings)
	latLongOfParkings =	[]string{"77.458059, 23.231247","77.338690, 23.270774","77.412542, 23.267715","77.407086, 23.263625","77.456512, 23.268870","77.436565,23.214282"}
	
	nameOfParkings := make([]string, NoOfParkings)
	nameOfParkings = []string{"Smart Parking Bhopal","New Market Parking Lot","Parking T.T Nagar","Bhopal Juntion Parking","Parking Sultania Road","Arera Colony Car Parking"}
	
	picURL := make([]string, NoOfParkings)
	picURL = []string{"/static/images/bhopalps1.jpg","/static/images/bhopalps2.jpg","/static/images/bhopalps3.jpg","/static/images/bhopalps4.jpg","/static/images/bhopalps5.jpg","/static/images/bhopalps5.jpg"}
	
	distanceFromUser := make([]float64, NoOfParkings)
	ETAOfParkings := make([]float64, NoOfParkings)

	city = &City{
		Name: "Indore",
		NoOfParkingSlots: int64(NoOfParkings), 
		PicURL: "/static/images/indore.png",
		Longitude: "75.857727",
        Latitude: "22.719568",
	}

	var arrOfSlots []Slot
	for i := 0; i < NoOfParkings; i++ {
			var result MapBoxResp
			resp, err := http.Get("https://api.mapbox.com/directions/v5/mapbox/driving/"+userLong+","+userLat+";"+latLongOfParkings[i]+"?access_token=pk.eyJ1Ijoic2hhc2hhbmstcHJha2FzaCIsImEiOiJja2ZjZzJmOGsxN3kxMnRvNXp1MDJxcGFzIn0.69B3TiOom1WXyot_FdTVHA")
			if err != nil{
				log.Fatalln(err)
			}

			defer resp.Body.Close()
			data, _ := ioutil.ReadAll(resp.Body)
			json.Unmarshal(data,&result)
			dis :=  math.Round(((result.Routes[0].Distance)/1000)*100)/100
			dur :=  math.Round(((result.Routes[0].Duration)/3600)*100)/100
			distanceFromUser[i] = dis
			ETAOfParkings[i] = dur
			slot:=Slot{
				Name: nameOfParkings[i], 
				DistanceFromYou: distanceFromUser[i],
				ETA: ETAOfParkings[i],
				PicURL: picURL[i],
			}
			arrOfSlots = append(arrOfSlots,slot)
	}
	city.Parking = arrOfSlots
		fmt.Println("Indore is reached")
	}

	templates.ExecuteTemplate(w, "cityParkingInfo.gohtml", city)

	//comment := r.PostForm.Get("cities")
	//client.LPush(ctx,"comments",comment)
	//http.Redirect(w, r, "/", 302)
}