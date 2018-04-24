package counter

import (
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var (
	path, _     = os.Getwd()
	search_page = path + "/../assets/search.html"
)

//table which looks like an excel sheet
type Sheet struct {
	rows []Rows
}
type Rows struct {
	Columns []string
}

// input values from console
type Values struct {
	StartColumn string
	NumRows     string
	NumCols     string
	Inputarray  [][]string
}

// web portal
func Portal(w http.ResponseWriter, r *http.Request) {

	login_tpl := template.Must(template.ParseFiles(search_page))
	login_tpl.Execute(w, struct{ Success bool }{true})

}

// This section of the code gets the desired section of the input table according to the inputs provided
func Columnfinder(w http.ResponseWriter, r *http.Request) {
	var input Sheet
	var output Sheet
	var startcolumn int
	var rows int
	var columns int
	var data Values
	var returnvalue [][]string

	err := unmarshaldata(r.Body, &data)

	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		ReturnJSONAPIError(w, "the input is not valid")
		return
	}
	startcolumn, _ = strconv.Atoi(data.StartColumn)
	rows, _ = strconv.Atoi(data.NumRows)
	columns, _ = strconv.Atoi(data.NumCols)
	var row []Rows
	for _, columns := range data.Inputarray {
		temp := Rows{Columns: columns}
		row = append(row, temp)
	}
	input.rows = row
	// the logic to get the descired rows and columns based on the input
	if rows > len(input.rows) {
		w.WriteHeader(http.StatusBadRequest)
		ReturnJSONAPIError(w, "the number of rows entered is not valid")
		return
	}
	for i := 0; i <= rows-1; i++ {
		columndata := make([]string, 0)
		if startcolumn > len(input.rows[i].Columns) || (startcolumn-1+columns-1) > len(input.rows[i].Columns)-1 {
			w.WriteHeader(http.StatusBadRequest)
			ReturnJSONAPIError(w, "the columns entered is not valid")
			return
		} else {
			for j := startcolumn - 1; j <= startcolumn-1+columns-1; j++ {
				columndata = append(columndata, input.rows[i].Columns[j])
			}
			output.rows = append(output.rows, Rows{columndata})
		}

	}

	for _, v := range output.rows {
		returnvalue = append(returnvalue, v.Columns)
	}

	ReturnJSONAPISuccess(w, returnvalue)
}

//used for unmarshal
func unmarshaldata(Body io.ReadCloser, subject interface{}) error {
	u, _ := ioutil.ReadAll(Body)

	err := json.Unmarshal(u, &subject)
	return err
}

// return success output
func ReturnJSONAPISuccess(w http.ResponseWriter, extra interface{}) {
	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(struct {
		Status bool
		Msg    string
		Extra  interface{}
	}{Msg: "success", Status: true, Extra: extra})

	w.Write(j)
}

// return error output
func ReturnJSONAPIError(w http.ResponseWriter, extra interface{}) {
	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(struct {
		Status bool
		Msg    string
		Extra  interface{}
	}{Msg: "Failed", Status: false, Extra: extra})

	w.Write(j)
}
