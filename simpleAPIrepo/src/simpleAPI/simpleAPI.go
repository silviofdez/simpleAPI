package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	//"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// ClinicalData Struct
type ClinicalData struct {
	Sample    string `json:"sample"`
	Condition string `json:"condition"`
	Age       string `json:"age"`
	Gender    string `json:"gender"`
	Race      string `json:"race"`
}

//Data Struct
type Data struct {
	Sample string `json:"sample"`
	M1     string `json:"m1` ///string in order to suppor NA
	M2     string `json:"m2`
	M3     string `json:"m3`
	M4     string `json:"m4`
	M5     string `json:"m5`
	M6     string `json:"m6`
	M7     string `json:"m7`
	M8     string `json:"m8`
	M9     string `json:"m9`
	M10    string `json:"m10`
	M11    string `json:"m11`
	M12    string `json:"m12`
	M13    string `json:"m13`
	M14    string `json:"m14`
	M15    string `json:"m15`
	M16    string `json:"m16`
	M17    string `json:"m17`
	M18    string `json:"m18`
}

//variable to map datasets - training
var trainingDataset = map[string]*Data{}

//variable to map datasets - validation
var validationDataset = map[string]*Data{}

//variable to map clinical datasets - training
var validationClinicalDataset = map[string]*ClinicalData{}

//variable to map clinical datasets - validation
var trainingClinicalDataset = map[string]*ClinicalData{}

//main loads data and starts api
func main() {
	///no really etl stage, just loading plain data
	log.Println("Loading data...")
	loadcvs("clinical/clinical_data_training")
	loadcvs("clinical/clinical_data_validation")
	loadcvs("data/training")
	loadcvs("data/validation")
	log.Println("Data loaded")
	///adding routers and muxes
	router := mux.NewRouter()
	http.Handle("/", &MyServer{router})
	router.HandleFunc("/data/{datatype}", handleDataSet).Methods("GET")
	router.HandleFunc("/data/{datatype}/{sample}", handleData).Methods("GET", "POST")
	router.HandleFunc("/clinicaldata/{datatype}", handleClinicalDataSet).Methods("GET")
	router.HandleFunc("/clinicaldata/{datatype}/{sample}", handleClinicalData).Methods("GET", "POST")
	log.Println("API listening...")
	http.ListenAndServe(":8080", nil)
	//http.ListenAndServe(":8080", handlers.CORS()(router))
}

type MyServer struct {
	r *mux.Router
}

func (s *MyServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, x-requested-with")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
}

//handler for specific data sample - create and read allowed, no update, no
//delete
func handleData(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	sample := vars["sample"]
	datatype := vars["datatype"]
	var ok bool
	var data *Data
	///two kinds of data, two vars, two origins...
	if datatype == "training" {
		data, ok = trainingDataset[sample]
	} else if datatype == "validation" {
		data, ok = validationDataset[sample]
	} else {
		ok = false
	}

	switch req.Method {
	case "GET":
		if !ok {
			res.WriteHeader(http.StatusNotFound)
			fmt.Fprint(res, string("Sample not found"))
		}
		outgoingJSON, error := json.Marshal(data)
		if error != nil {
			log.Println(error.Error())
			http.Error(res, error.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(res, string(outgoingJSON))
	case "POST":
		data := new(Data)
		decoder := json.NewDecoder(req.Body)
		error := decoder.Decode(&data)
		if error != nil {
			log.Println(error.Error())
			http.Error(res, error.Error(), http.StatusInternalServerError)
			return
		}
		if datatype == "training" {
			trainingDataset[sample] = data
			savecsv("data/training")
		} else if datatype == "validation" {
			validationDataset[sample] = data
			savecsv("data/validation")
		} else {
			res.WriteHeader(http.StatusNotFound)
			fmt.Fprint(res, string("Dataset not found"))
		}
		///naive approach, no data validation
		outgoingJSON, err := json.Marshal(data)
		if err != nil {
			log.Println(error.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusCreated)
		fmt.Fprint(res, string(outgoingJSON))
	}
}

//handler for generic data - create and read allowed, no update, no
//delete
func handleDataSet(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	datatype := vars["datatype"]
	var outgoingJSON []byte
	var error error
	if datatype == "training" {
		outgoingJSON, error = json.Marshal(trainingDataset)
	} else if datatype == "validation" {
		outgoingJSON, error = json.Marshal(validationDataset)
	} else {
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprint(res, string("Dataset not found"))
	}
	if error != nil {
		log.Println(error.Error())
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(res, string(outgoingJSON))
}

//handler for specific clinical data - create and read allowed, no
//update, no delete
func handleClinicalData(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	sample := vars["sample"]
	datatype := vars["datatype"]
	var ok bool
	var clinicaldata *ClinicalData
	if datatype == "training" {
		clinicaldata, ok = trainingClinicalDataset[sample]
	} else if datatype == "validation" {
		clinicaldata, ok = validationClinicalDataset[sample]
	} else {
		ok = false
	}

	switch req.Method {
	case "GET":
		if !ok {
			res.WriteHeader(http.StatusNotFound)
			fmt.Fprint(res, string("Sample not found"))
		}
		outgoingJSON, error := json.Marshal(clinicaldata)
		if error != nil {
			log.Println(error.Error())
			http.Error(res, error.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(res, string(outgoingJSON))
	case "POST":
		clinicaldata := new(ClinicalData)
		decoder := json.NewDecoder(req.Body)
		error := decoder.Decode(&clinicaldata)
		if error != nil {
			log.Println(error.Error())
			http.Error(res, error.Error(), http.StatusInternalServerError)
			return
		}
		if datatype == "training" {
			trainingClinicalDataset[sample] = clinicaldata
			savecsv("clinical/clinical_data_training")
		} else if datatype == "validation" {
			validationClinicalDataset[sample] = clinicaldata
			savecsv("clinical/clinical_data_validation")
		} else {
			res.WriteHeader(http.StatusNotFound)
			fmt.Fprint(res, string("Dataset not found"))
		}
		///again naive approach
		outgoingJSON, err := json.Marshal(clinicaldata)
		if err != nil {
			log.Println(error.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusCreated)
		fmt.Fprint(res, string(outgoingJSON))
	}
}

//handler for generic clinical data - create and read allowed, no update, no
//delete
func handleClinicalDataSet(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	datatype := vars["datatype"]
	var outgoingJSON []byte
	var error error
	if datatype == "training" {
		outgoingJSON, error = json.Marshal(trainingClinicalDataset)
	} else if datatype == "validation" {
		outgoingJSON, error = json.Marshal(validationClinicalDataset)
	} else {
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprint(res, string("Dataset not found"))
	}
	if error != nil {
		log.Println(error.Error())
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(res, string(outgoingJSON))
}

//read csv and load data into vars
func loadcvs(filename string) {
	csvfile, err := os.Open("dataset/" + filename + ".csv")
	if err != nil {
		log.Println(err)
		return
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)

	reader.FieldsPerRecord = -1 // see the Reader struct information below

	rawCSVdata, err := reader.ReadAll()

	if err != nil {
		log.Println(err)
	}

	///clinical data
	if filename == "clinical/clinical_data_training" || filename == "clinical/clinical_data_validation" {
		for _, each := range rawCSVdata {
			tempClinicalData := new(ClinicalData)
			//skipping headers
			if each[0] != "Sample" {
				tempClinicalData.Sample = each[0]
				tempClinicalData.Condition = each[1]
				tempClinicalData.Age = each[2]
				tempClinicalData.Gender = each[3]
				tempClinicalData.Race = each[4]
				if filename == "clinical/clinical_data_training" {
					trainingClinicalDataset[tempClinicalData.Sample] = tempClinicalData
				} else if filename == "clinical/clinical_data_validation" {
					validationClinicalDataset[tempClinicalData.Sample] = tempClinicalData
				}
			}
		}
	}
	///data
	if filename == "data/training" || filename == "data/validation" {
		for _, each := range rawCSVdata {
			tempData := new(Data)
			//skipping headers
			if each[0] != "Sample" {
				tempData.Sample = each[0]
				tempData.M1 = each[1]
				tempData.M2 = each[2]
				tempData.M3 = each[3]
				tempData.M4 = each[4]
				tempData.M5 = each[5]
				tempData.M6 = each[6]
				tempData.M7 = each[7]
				tempData.M8 = each[8]
				tempData.M9 = each[9]
				tempData.M10 = each[10]
				tempData.M11 = each[11]
				tempData.M12 = each[12]
				tempData.M13 = each[13]
				tempData.M14 = each[14]
				tempData.M15 = each[15]
				tempData.M16 = each[16]
				tempData.M17 = each[17]
				tempData.M18 = each[18]
				if filename == "data/training" {
					trainingDataset[tempData.Sample] = tempData
				} else if filename == "data/validation" {
					validationDataset[tempData.Sample] = tempData
				}
			}
		}
	}
}

//write csv - naive approach, write whole files
func savecsv(filename string) {
	file, err := os.Create("dataset/" + filename + ".csv")
	checkError("Cannot create file", err)
	defer file.Close()

	workingClinicalDataset := map[string]*ClinicalData{}
	workingDataset := map[string]*Data{}
	if filename == "clinical/clinical_data_training" {
		workingClinicalDataset = trainingClinicalDataset
	} else if filename == "clinical/clinical_data_validation" {
		workingClinicalDataset = validationClinicalDataset
	} else {
		workingClinicalDataset = nil
	}
	if filename == "data/validation" {
		workingDataset = validationDataset
	} else if filename == "data/training" {
		workingDataset = trainingDataset
	} else {
		workingDataset = nil
	}

	writer := csv.NewWriter(file)

	valuecsv := []string{}
	if workingClinicalDataset != nil {
		valuecsv = []string{"Sample", "Condition", "Age", "Gender", "Race"}
		err = writer.Write(valuecsv)
		checkError("Cannot write to file", err)
		for _, value := range workingClinicalDataset {
			valuecsv = []string{value.Sample, value.Condition, value.Age, value.Gender, value.Race}
			err := writer.Write(valuecsv)
			checkError("Cannot write to file", err)
		}
	}

	if workingDataset != nil {
		///header
		valuecsv = []string{"Sample", "M1", "M2", "M3", "M4", "M5", "M6", "M7", "M8", "M9", "M10", "M11", "M12", "M13", "M14", "M15", "M16", "M17", "M18"}
		err = writer.Write(valuecsv)
		checkError("Cannot write to file", err)
		for _, value := range workingDataset {
			//skipping headers??
			if value.Sample != "Sample" {
				valuecsv = []string{value.Sample, value.M1, value.M2, value.M3, value.M4, value.M5, value.M6, value.M7, value.M8, value.M9, value.M10, value.M11, value.M12, value.M13, value.M14, value.M15, value.M16, value.M17, value.M18}
				err := writer.Write(valuecsv)
				checkError("Cannot write to file", err)
			}
		}
	}
	defer writer.Flush()
}

//aux func for error checking
func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
