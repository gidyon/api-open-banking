package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gidyon/api-open-banking/internal/model/bank"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func createBankCustomerHandlerV1(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	bankID := params.ByName("bankID")

	url := fmt.Sprintf("%s/obp/v3.1.0/banks/%s/customers", BASEURL, bankID)

	res, err := http.Post(url, r.Header.Get("content-type"), r.Body)
	if err != nil {
		http.Error(w, "failed request: "+err.Error(), res.StatusCode)
		return
	}

	resBytes := make([]byte, 0)

	_, err = res.Body.Read(resBytes)
	if err != nil {
		http.Error(w, "failed to read response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// We write the response
	w.Write(resBytes)
}

func createBankCustomerHandlerV2(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	bankID := params.ByName("bankID")

	url := fmt.Sprintf("%s/obp/v3.1.0/banks/%s/customers", BASEURL, bankID)

	// Let's forward the request to their api and return response as is.
	http.Redirect(w, r, url, http.StatusContinue)
}

// adds a bank resource to user list of banks
func addBankHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// 1. We get the bank from api
	// 2. Unmarshal it
	// 2. Add it to the user list of banks

	bankID := params.ByName("bankID")

	url := fmt.Sprintf("%s/obp/v3.1.0/banks/%s", BASEURL, bankID)

	res, err := http.Get(url)
	if err != nil {
		http.Error(w, "couldn't get bank resource: "+err.Error(), res.StatusCode)
		return
	}

	appBank := &bank.Bank{}

	// We write the response but in json format
	err = json.NewDecoder(r.Body).Decode(appBank)
	if err != nil {
		http.Error(w, "couldn't decode bank payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = bank.AddBank(params.ByName("userID"), appBank)
	if err != nil {
		http.Error(w, "couldn't add bank to user list of banks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the response
	fmt.Fprint(w, "bank added to user list of banks")
}

// get banks resource the user is registered with
func getBanksHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID := params.ByName("userID")

	banks, err := bank.GetBanks(userID)
	if err != nil {
		http.Error(w, "failed to get user banks: "+err.Error(), http.StatusNotFound)
		return
	}

	// We write the response but in json format
	err = json.NewEncoder(w).Encode(&banks)
	if err != nil {
		http.Error(w, "failed encode user banks: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// get single bank resource the user is registered with
func getBankHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID := params.ByName("userID")
	bankID := params.ByName("bankID")

	bank, err := bank.GetBank(userID, bankID)
	if err != nil {
		http.Error(w, "failed to get user bank: "+err.Error(), http.StatusNotFound)
		return
	}

	// We write the response but in json format
	err = json.NewEncoder(w).Encode(bank)
	if err != nil {
		http.Error(w, "failed encode user bank: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
