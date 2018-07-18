package main

import (
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
    "encoding/json"
    "net/http"
    "strings"
)
var db *sql.DB
var err error
func main() {	
	db, err = sql.Open("mysql", "test:password@tcp(localhost:3306)/test")
	 
	http.HandleFunc("/select/", httpSelect)
	http.HandleFunc("/insert/", httpInsert)
	http.ListenAndServe(":8080", nil)

   	db.Close()
}

func getJSON(queryString string) ([]byte, error) {


    
    if(err != nil){
        return []byte(""), err
    }   
    rows, err := db.Query(queryString)
    
    if err != nil {
    	return []byte(""), err
    }
    
    defer rows.Close()
    cols, err := rows.Columns()
    
    if err != nil {
    	return []byte(""), err
    }
    
    count := len(cols)
    tableData := make([]map[string]interface{}, 0)
    values := make([]interface{}, count)
    valuePointers := make([]interface{}, count)
    
    for rows.Next() {
    	for i := 0; i < count; i++ {
    		valuePointers[i] = &values[i]
    	}
    	rows.Scan(valuePointers...)
    	entry := make(map[string]interface{})
    	for i, col := range cols {
    		var v interface{}
    		val := values[i]
    		b, ok := val.([]byte)
    		if ok {
    			v = string(b)
    		} else {
    			v = val
    		}
    		entry[col] = v
    	}
    	
    	tableData = append(tableData, entry)
    }
    jsonData, err := json.Marshal(tableData)
    if err != nil {
    	return []byte(""), err
    }
    fmt.Println(string(jsonData))
    return jsonData, nil
}

func httpSelect(w http.ResponseWriter, r *http.Request) {
	
	query := r.URL.Path
	query = strings.TrimPrefix(query, "/select/")
	message, err := getJSON(query)
	
	if err != nil {
		panic(err.Error)
	}
	w.Write(message)
}

func httpInsert(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Path
	query = strings.TrimPrefix(query, "/insert/")
	res, err := db.Exec(query)
	if err != nil {
		w.Write([]byte("Error"))
	} else {
		aff, err := res.RowsAffected()
		fmt.Println(aff)
		if err != nil {
			return
		}
		w.Write([]byte("Success"))
	}
}
