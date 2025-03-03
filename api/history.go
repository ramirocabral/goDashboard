package api

import (
    "net/http"
    "time"
)

func (app *app) cpuHistoryHandler (w http.ResponseWriter, r *http.Request){

    start := r.URL.Query().Get("start")
    end := r.URL.Query().Get("end")
    interval := r.URL.Query().Get("interval")

    endTime, _ := time.Parse(time.RFC3339, end)
    startTime, _ := time.Parse(time.RFC3339, start)

    data, err := app.store.ReadCpuStats(startTime, endTime, interval)

    if err != nil{
	app.internalServerError(w, r, err)
	return
    }

    err = writeJSON(w, http.StatusOK, data)
    if err != nil{
	app.internalServerError(w, r, err)
	return
    }
}

func (app *app) ioHistoryHandler (w http.ResponseWriter, r *http.Request){

	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	interval := r.URL.Query().Get("interval")

	endTime, _ := time.Parse(time.RFC3339, end)
	startTime, _ := time.Parse(time.RFC3339, start)

	data, err := app.store.ReadIOStats(startTime, endTime, interval)

	if err != nil{
	    app.internalServerError(w, r, err)
	    return
	}

	err = writeJSON(w, http.StatusOK, data)
	if err != nil{
	    app.internalServerError(w, r, err)
	    return
	}
}

func (app *app) memoryHistoryHandler (w http.ResponseWriter, r *http.Request){

    start := r.URL.Query().Get("start")    
    end := r.URL.Query().Get("end")
    interval := r.URL.Query().Get("interval")

    endTime, _ := time.Parse(time.RFC3339, end)
    startTime, _ := time.Parse(time.RFC3339, start)

    data, err := app.store.ReadMemoryStats(startTime, endTime, interval)

    if err != nil{
	app.internalServerError(w, r, err)
	return
    }

    err = writeJSON(w, http.StatusOK, data)
    if err != nil{
	app.internalServerError(w, r, err)
	return
    }
}

func (app *app) networkHistoryHandler (w http.ResponseWriter, r *http.Request){

    start := r.URL.Query().Get("start")
    end := r.URL.Query().Get("end")
    interval := r.URL.Query().Get("interval")

    endTime, _ := time.Parse(time.RFC3339, end)
    startTime, _ := time.Parse(time.RFC3339, start)

    data, err := app.store.ReadNetworkStats(startTime, endTime, interval)

    if err != nil{
	app.internalServerError(w, r, err)
	return
    }

    err = writeJSON(w, http.StatusOK, data)
    if err != nil{
	app.internalServerError(w, r, err)
	return
    }
}

