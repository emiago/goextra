package mysqlhelp

import (
	"database/sql"
	"fmt"
	"strconv"
)

func TableExists(db *sql.DB, table string) (bool, error) {
	var tname string
	q := fmt.Sprintf("SHOW TABLES LIKE '%s'", table)
	err := db.QueryRow(q).Scan(&tname)

	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	}
	return true, nil
}

func GetTableCommentRev(db *sql.DB, table string, revStr string) (float64, error) {
	q := fmt.Sprintf("SHOW TABLE STATUS LIKE '%s'", table)
	res, err := ScanIntoMap(db, q)
	if err != nil {
		return 0, err
	}

	if len(res) == 0 {
		return 0, fmt.Errorf("No result for table status %s", table)
	}

	row := res[0]
	commentVal, ok := row["Comment"]
	if !ok {
		return 0, fmt.Errorf("No comment for table status %s", table)
	}

	comment, ok := commentVal.(string)
	if !ok {
		return 0, fmt.Errorf("No comment for table status %s", table)
	}

	//Including space in comment
	if len(comment) > len(revStr)+1 {
		comment = comment[len(revStr)+1:]
	}

	var f float64 = 0
	if len(comment) > 0 {
		var ferr error
		f, ferr = strconv.ParseFloat(comment, 64)
		if ferr != nil {
			return 0, ferr
		}
	}

	return f, nil
}

func SetTableCommentRev(db *sql.DB, table string, revStr string, rev float64) error {
	q := fmt.Sprintf("ALTER TABLE %s COMMENT = '%s %f'", table, revStr, rev)
	_, err := db.Exec(q)
	return err
}

func ScanIntoMap(db *sql.DB, query string) ([]map[string]interface{}, error) {
	rows, _ := db.Query(query) // Note: Ignoring errors for brevity
	cols, _ := rows.Columns()
	res := make([]map[string]interface{}, 0)

	defer rows.Close()
	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}

		res = append(res, m)
	}

	return res, nil
}

func DoTransact(db *sql.DB, s []string) error {
	// Get new Transaction. See http://golang.org/pkg/database/sql/#DB.Begin
	txn, err := db.Begin()

	if err != nil {
		return err
	}

	defer func() {
		// Rollback the transaction after the function returns.
		// If the transaction was already commited, this will do nothing.
		_ = txn.Rollback()
	}()

	for _, q := range s {
		// Execute the query in the transaction.
		_, err := txn.Exec(q)

		if err != nil {
			return fmt.Errorf("%s. Query=%s", err, q)
		}
	}

	// Commit the transaction.
	return txn.Commit()
}

func TableHasIndex(db *sql.DB, table string, index string) (bool, error) {
	var count int
	q := fmt.Sprintf("SELECT COUNT(1) IndexIsThere FROM INFORMATION_SCHEMA.STATISTICS WHERE table_schema=DATABASE() AND table_name='%s' AND index_name='%s'",
		table,
		index,
	)

	err := db.QueryRow(q).Scan(&count)
	return count > 0, err
}

func TableHasColumn(db *sql.DB, table string, index string) (bool, error) {
	var count int
	q := fmt.Sprintf("SELECT COUNT(1) IndexIsThere FROM INFORMATION_SCHEMA.COLUMNS WHERE table_schema=DATABASE() AND table_name='%s' AND column_name='%s'",
		table,
		index,
	)

	err := db.QueryRow(q).Scan(&count)
	return count > 0, err
}
