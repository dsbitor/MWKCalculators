// Package csvtable exports a SQLite table to CSV and replaces a
// table's contents from CSV, the mechanism the user database's
// machine-specific tables (change gears, dividing head setups, and
// similar) are populated through, since those tables are never
// seeded with a fabricated default.
//
// Import requires the CSV header to name columns that already exist
// on the target table, created ahead of time by a migration; this
// package never creates or alters a table's schema, only its rows.
package csvtable

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

// Export writes every row of table to w as CSV, with a header row of
// column names taken from the table itself.
func Export(ctx context.Context, db *sql.DB, table string, w io.Writer) error {
	columns, err := tableColumns(ctx, db, table)
	if err != nil {
		return err
	}

	rows, err := db.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM %s", quoteColumnList(columns), quoteIdentifier(table)))
	if err != nil {
		return fmt.Errorf("query table %s for export: %w", table, err)
	}
	defer rows.Close()

	writer := csv.NewWriter(w)
	if err := writer.Write(columns); err != nil {
		return fmt.Errorf("write CSV header for table %s: %w", table, err)
	}

	if err := writeRows(rows, columns, writer); err != nil {
		return err
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate rows of table %s: %w", table, err)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("flush CSV output for table %s: %w", table, err)
	}
	return nil
}

func writeRows(rows *sql.Rows, columns []string, writer *csv.Writer) error {
	values := make([]any, len(columns))
	pointers := make([]any, len(columns))
	for i := range values {
		pointers[i] = &values[i]
	}

	record := make([]string, len(columns))
	for rows.Next() {
		if err := rows.Scan(pointers...); err != nil {
			return fmt.Errorf("scan row for export: %w", err)
		}
		for i, v := range values {
			record[i] = formatValue(v)
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("write CSV row: %w", err)
		}
	}
	return nil
}

func formatValue(v any) string {
	if v == nil {
		return ""
	}
	if b, ok := v.([]byte); ok {
		return string(b)
	}
	return fmt.Sprint(v)
}

// Import replaces every row of table with the rows read from r,
// which must be CSV with a header row naming columns that already
// exist on table. The replacement happens inside a single
// transaction: a malformed row, or a column name in the header that
// the table does not have, leaves table's existing contents
// untouched rather than partially replaced.
func Import(ctx context.Context, db *sql.DB, table string, r io.Reader) error {
	reader := csv.NewReader(r)
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("read CSV header for table %s: %w", table, err)
	}

	existingColumns, err := tableColumns(ctx, db, table)
	if err != nil {
		return err
	}
	if err := validateHeader(table, header, existingColumns); err != nil {
		return err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("read CSV rows for table %s: %w", table, err)
	}

	return importRecords(ctx, db, table, header, records)
}

func validateHeader(table string, header, existingColumns []string) error {
	known := make(map[string]bool, len(existingColumns))
	for _, column := range existingColumns {
		known[column] = true
	}
	for _, column := range header {
		if !known[column] {
			return fmt.Errorf("CSV column %q is not a column of table %s", column, table)
		}
	}
	return nil
}

func importRecords(ctx context.Context, db *sql.DB, table string, header []string, records [][]string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction to import table %s: %w", table, err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s", quoteIdentifier(table))); err != nil {
		return fmt.Errorf("clear table %s before import: %w", table, err)
	}

	insertSQL := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		quoteIdentifier(table),
		quoteColumnList(header),
		placeholders(len(header)),
	)

	for rowIndex, record := range records {
		if len(record) != len(header) {
			return fmt.Errorf("row %d of table %s has %d values, want %d", rowIndex+1, table, len(record), len(header))
		}
		args := make([]any, len(record))
		for i, value := range record {
			args[i] = value
		}
		if _, err := tx.ExecContext(ctx, insertSQL, args...); err != nil {
			return fmt.Errorf("insert row %d into table %s: %w", rowIndex+1, table, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit import of table %s: %w", table, err)
	}
	return nil
}

func tableColumns(ctx context.Context, db *sql.DB, table string) ([]string, error) {
	rows, err := db.QueryContext(ctx, fmt.Sprintf("PRAGMA table_info(%s)", quoteIdentifier(table)))
	if err != nil {
		return nil, fmt.Errorf("read schema of table %s: %w", table, err)
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var (
			cid        int
			name       string
			columnType string
			notNull    int
			defaultVal any
			pk         int
		)
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultVal, &pk); err != nil {
			return nil, fmt.Errorf("scan schema of table %s: %w", table, err)
		}
		columns = append(columns, name)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate schema of table %s: %w", table, err)
	}
	if len(columns) == 0 {
		return nil, fmt.Errorf("table %s does not exist or has no columns", table)
	}
	return columns, nil
}

func quoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func quoteColumnList(columns []string) string {
	quoted := make([]string, len(columns))
	for i, c := range columns {
		quoted[i] = quoteIdentifier(c)
	}
	return strings.Join(quoted, ", ")
}

func placeholders(n int) string {
	parts := make([]string, n)
	for i := range parts {
		parts[i] = "?"
	}
	return strings.Join(parts, ", ")
}
