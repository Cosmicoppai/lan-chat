package shows

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"lan-chat/admin"
	"lan-chat/admin/dbErrors"
	"lan-chat/httpErrors"
	"lan-chat/logger"
	"lan-chat/middleware"
	"net/http"
	"strconv"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.HandlerFunc(listShow).ServeHTTP(w, r)
	case http.MethodPost:
		middleware.AuthMiddleware(http.HandlerFunc(createShow)).ServeHTTP(w, r)
	case http.MethodPut:
		middleware.AuthMiddleware(http.HandlerFunc(updateShow)).ServeHTTP(w, r)
	case http.MethodDelete:
		middleware.AuthMiddleware(http.HandlerFunc(deleteShow)).ServeHTTP(w, r)

	}

}

func createShow(w http.ResponseWriter, r *http.Request) {
	var show Show
	err := json.NewDecoder(r.Body).Decode(&show)

	if err != nil {
		httpErrors.UnProcessableEntry(w)
		return
	}
	if show.Type == "" || show.Name == "" || show.TotalEps == nil {
		httpErrors.BadRequest(w, "One or more fields are absent.")
		return
	}
	_, err = admin.Db.Exec("INSERT INTO lan_show.shows (name, totaleps, typ) VALUES ($1, $2, $3)", show.Name, show.TotalEps, show.Type)
	if err != nil {
		logger.ErrorLog.Println("Error while creating show: ", err)
		if dbErrors.IntegrityViolation(err) {
			httpErrors.UnProcessableEntry(w)
			return
		}
		httpErrors.InternalServerError(w)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("Show Created Successfully"))

}

func ListShows(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpErrors.MethodNotAllowed(w)
		return
	}
	var showFilter ShowFilter
	FilterMap := map[string]string{}

	queryParams := r.URL.Query()
	for key, value := range queryParams { // convert map[string][]string to map[string]string
		FilterMap[key] = value[0]
	}
	data, _ := json.Marshal(FilterMap)
	err := json.Unmarshal(data, &showFilter)
	if err != nil {
		logger.ErrorLog.Println(err)
		httpErrors.UnProcessableEntry(w, "One or more invalid query parameter")
		return
	}
	var rows *sql.Rows

	if (ShowFilter{}) != showFilter {
		rows, err = _listFilter(showFilter, "SELECT * FROM lan_show.shows WHERE ", true)
	} else {
		rows, err = _listFilter(showFilter, "SELECT * from lan_show.shows;", false)
	}

	if err != nil && dbErrors.InternalServerError(err) {
		logger.ErrorLog.Println(err)
		httpErrors.InternalServerError(w)
		return
	}
	_listShowHelper(w, rows)
}

func _listFilter(showFilter ShowFilter, baseQuery string, filter bool) (*sql.Rows, error) {
	var (
		rows *sql.Rows
		err  error
	)
	if filter {
		filterQuery := ""
		var filters []interface{}

		if showFilter.Name != "" {
			filters = append(filters, showFilter.Name)
			filterQuery += fmt.Sprintf("name=$%s ", strconv.Itoa(len(filters)))
		}
		if showFilter.Type != "" {
			filters = append(filters, showFilter.Type)
			filterQuery += fmt.Sprintf("AND typ=$%s ", strconv.Itoa(len(filters)))
		}
		if showFilter.TotalEps != nil {
			filters = append(filters, *showFilter.TotalEps)
			filterQuery += fmt.Sprintf("AND totalEps>$%s", strconv.Itoa(len(filters)))
		}
		filterQuery = strings.TrimPrefix(filterQuery, "AND ")
		query := baseQuery + filterQuery + ";"
		logger.InfoLog.Println(query)
		rows, err = admin.Db.Query(query, filters...)
	} else {
		rows, err = admin.Db.Query(baseQuery)
	}

	logger.ErrorLog.Println(err)

	return rows, err

}

func listShow(w http.ResponseWriter, r *http.Request) {
	id, err := getId(w, r)
	if err != nil {
		return
	}

	rows, err := admin.Db.Query("SELECT * FROM lan_show.shows WHERE id=$1;", id)
	if err != nil && dbErrors.InternalServerError(err) {
		logger.ErrorLog.Println(err)
		httpErrors.InternalServerError(w)
		return
	}
	_listShowHelper(w, rows)

}

func _listShowHelper(w http.ResponseWriter, rows *sql.Rows) {
	defer rows.Close()
	var (
		shows []Show
		show  Show
	)

	for rows.Next() {
		err := rows.Scan(&show.Id, &show.Name, &show.TotalEps, &show.Type)
		if err != nil {
			logger.ErrorLog.Println("Error while scanning rows for shows: ", err)
		}
		shows = append(shows, show)
	}
	err := rows.Err()
	if err != nil {
		httpErrors.InternalServerError(w)
		return
	}
	if len(shows) == 0 {
		httpErrors.NotFound(w, "No Result found")
		return
	}
	w.Header().Set("Content-type", "application/json")
	_ = json.NewEncoder(w).Encode(shows)

}

func updateShow(w http.ResponseWriter, r *http.Request) {
	id, err := getId(w, r)
	if err != nil {
		return
	}
	var showInfo Show
	err = json.NewDecoder(r.Body).Decode(&showInfo)
	if err != nil {
		httpErrors.BadRequest(w, "invalid json")
		return
	}
	if showInfo.Id != nil {
		httpErrors.UnProcessableEntry(w, "")
	}
	showInfo.Id = &id
	_updateShow(w, showInfo)

}

func _updateShow(w http.ResponseWriter, show Show) {
	baseQuery := "UPDATE lan_show.shows SET %s WHERE id=$1;"
	filterQuery := ""
	filters := []interface{}{*show.Id}

	if show.Name != "" {
		filters = append(filters, show.Name)
		filterQuery += fmt.Sprintf("name=$%s,", strconv.Itoa(len(filters)))
	}
	if show.Type != "" {
		filters = append(filters, show.Type)
		filterQuery += fmt.Sprintf("typ=$%s,", strconv.Itoa(len(filters)))
	}
	if show.TotalEps != nil {
		filters = append(filters, *show.TotalEps)
		filterQuery += fmt.Sprintf("totalEps=$%s", strconv.Itoa(len(filters)))
	}
	filterQuery = strings.TrimSuffix(filterQuery, ",")
	query := fmt.Sprintf(baseQuery, filterQuery)

	_, err := admin.Db.Exec(query, filters...)
	if err != nil && dbErrors.InternalServerError(err) {
		logger.ErrorLog.Println(err)
		httpErrors.InternalServerError(w)
		return
	}
	_, _ = w.Write([]byte("show updated successfully"))
}

func deleteShow(w http.ResponseWriter, r *http.Request) {
	id, err := getId(w, r)
	if err != nil {
		return
	}
	_, err = admin.Db.Query("DELETE FROM lan_show.shows WHERE id=$1", id)
	if err != nil && dbErrors.InternalServerError(err) {
		logger.ErrorLog.Println(err)
		httpErrors.InternalServerError(w)
		return
	}
	_, _ = w.Write([]byte("show deleted successfully"))

}

func getId(w http.ResponseWriter, r *http.Request) (int64, error) {
	uri := strings.Trim(r.RequestURI, "/")
	pathParams := strings.Split(uri, "/")
	if len(pathParams) < 2 {
		httpErrors.BadRequest(w, "id not present")
		return 0, errors.New("id not present")
	}
	id, err := strconv.ParseInt(pathParams[1], 10, 64)
	if err != nil {
		httpErrors.BadRequest(w, "Invalid Id")
		return 0, err
	}
	return id, nil
}
