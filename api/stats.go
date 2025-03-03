package api

import (
	"net/http"
)

func (app *app) smartHandler(w http.ResponseWriter, r *http.Request){

    stats, err := app.statsManager.GetSMART()

    if err != nil{
	app.internalServerError(w, r, err)
	return
    }

    err = writeJSON(w, http.StatusOK, stats)  

    if err != nil{
	app.internalServerError(w, r, err)
	return
    }
}

func (app *app) hostHandler(w http.ResponseWriter, r *http.Request){

    stats, err := app.statsManager.GetHost()

    if err != nil{
	app.internalServerError(w, r, err)
	return
    }

    err = writeJSON(w, http.StatusOK, stats)  

    if err != nil{
	app.internalServerError(w, r, err)
	return
    }
}

func (app *app) diskHandler(w http.ResponseWriter, r *http.Request){
    
    stats, err := app.statsManager.GetDisk()

    if err != nil{
	app.internalServerError(w, r, err)
	return
    }

    err = writeJSON(w, http.StatusOK, stats)  

    if err != nil{
	app.internalServerError(w, r, err)
	return
    }
}
