package shows

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"lan-chat/admin"
	"lan-chat/admin/dbErrors"
	"lan-chat/httpErrors"
	"lan-chat/logger"
	"lan-chat/utils"
	"net/http"
	"strconv"
	"strings"
)

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

func listShows(w http.ResponseWriter, r *http.Request) {
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
		var filterQuery strings.Builder
		var filters []interface{}

		if showFilter.Name != "" {
			filters = append(filters, showFilter.Name)
			filterQuery.WriteString(fmt.Sprintf("name=$%s ", strconv.Itoa(len(filters))))
		}
		if showFilter.Type != "" {
			filters = append(filters, showFilter.Type)
			filterQuery.WriteString(fmt.Sprintf("AND typ=$%s ", strconv.Itoa(len(filters))))
		}
		if showFilter.TotalEps != nil {
			filters = append(filters, *showFilter.TotalEps)
			filterQuery.WriteString(fmt.Sprintf("AND totalEps>$%s", strconv.Itoa(len(filters))))
		}
		fQ := strings.TrimPrefix(filterQuery.String(), "AND ")
		query := baseQuery + fQ + ";"
		logger.InfoLog.Println(query)
		rows, err = admin.Db.Query(query, filters...)
	} else {
		rows, err = admin.Db.Query(baseQuery)
	}

	return rows, err

}

func listShow(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(utils.GetField(r, 0), 10, 64)

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
	id, _ := strconv.ParseInt(utils.GetField(r, 0), 10, 64)
	var showInfo Show
	err := json.NewDecoder(r.Body).Decode(&showInfo)
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
	var filterQuery strings.Builder
	filters := []interface{}{*show.Id}

	if show.Name != "" {
		filters = append(filters, show.Name)
		filterQuery.WriteString(fmt.Sprintf("name=$%s,", strconv.Itoa(len(filters))))
	}
	if show.Type != "" {
		filters = append(filters, show.Type)
		filterQuery.WriteString(fmt.Sprintf("typ=$%s,", strconv.Itoa(len(filters))))
	}
	if show.TotalEps != nil {
		filters = append(filters, *show.TotalEps)
		filterQuery.WriteString(fmt.Sprintf("totalEps=$%s", strconv.Itoa(len(filters))))
	}
	fQ := strings.TrimSuffix(filterQuery.String(), ",")
	query := fmt.Sprintf(baseQuery, fQ)

	_, err := admin.Db.Exec(query, filters...)
	if err != nil && dbErrors.InternalServerError(err) {
		logger.ErrorLog.Println(err)
		httpErrors.InternalServerError(w)
		return
	}
	_, _ = w.Write([]byte("show updated successfully"))
}

func deleteShow(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(utils.GetField(r, 0), 10, 64)
	_, err := admin.Db.Query("DELETE FROM lan_show.shows WHERE id=$1", id)
	if err != nil && dbErrors.InternalServerError(err) {
		logger.ErrorLog.Println(err)
		httpErrors.InternalServerError(w)
		return
	}
	_, _ = w.Write([]byte("show deleted successfully"))

}
