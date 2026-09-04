package analytics

import (
	"database/sql"
	"fmt"
	"math"
	"sort"
	"strings"
)

// Profile is the summary of a table.
type Profile struct {
	TableName   string
	RowCount    int
	ColumnCount int
	Columns     []ColumnProfile
}

// ColumnProfile contains per-column statistics.
type ColumnProfile struct {
	Name     string
	Type     string
	Missing  float64
	Unique   int
	Min      string
	Max      string
	Mean     string
	Median   string
	Sample   []string
}

// ProfileTable inspects a table and calculates simple statistics.
func ProfileTable(db *sql.DB, tableName string) (*Profile, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is required")
	}
	if strings.TrimSpace(tableName) == "" {
		return nil, fmt.Errorf("table name is required")
	}

	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM %s", quoteIdentifier(tableName))
	var rowCount int
	if err := db.QueryRow(countSQL).Scan(&rowCount); err != nil {
		return nil, fmt.Errorf("count rows: %w", err)
	}

	cols, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", quoteIdentifier(tableName)))
	if err != nil {
		return nil, fmt.Errorf("read schema: %w", err)
	}
	defer cols.Close()

	columnInfos := make([]struct {
		name string
		typ  string
	}, 0)
	for cols.Next() {
		var cid int
		var name, typ string
		var notNull int
		var defaultValue sql.NullString
		var pk int
		if err := cols.Scan(&cid, &name, &typ, &notNull, &defaultValue, &pk); err != nil {
			return nil, fmt.Errorf("scan schema: %w", err)
		}
		columnInfos = append(columnInfos, struct {
			name string
			typ  string
		}{name: name, typ: typ})
	}
	if err := cols.Err(); err != nil {
		return nil, fmt.Errorf("iterate schema: %w", err)
	}

	profile := &Profile{
		TableName:   tableName,
		RowCount:    rowCount,
		ColumnCount: len(columnInfos),
		Columns:     make([]ColumnProfile, 0, len(columnInfos)),
	}

	for _, col := range columnInfos {
		colProf, err := profileColumn(db, tableName, col.name, col.typ)
		if err != nil {
			return nil, err
		}
		profile.Columns = append(profile.Columns, *colProf)
	}
	return profile, nil
}

func profileColumn(db *sql.DB, tableName, columnName, columnType string) (*ColumnProfile, error) {
	selectSQL := fmt.Sprintf("SELECT %s FROM %s", quoteIdentifier(columnName), quoteIdentifier(tableName))
	rows, err := db.Query(selectSQL)
	if err != nil {
		return nil, fmt.Errorf("read column %s: %w", columnName, err)
	}
	defer rows.Close()

	values := make([]float64, 0)
	missing := 0
	uniqueSet := map[string]struct{}{}
	minVal := ""
	maxVal := ""

	for rows.Next() {
		var raw any
		if err := rows.Scan(&raw); err != nil {
			return nil, fmt.Errorf("scan value in %s: %w", columnName, err)
		}
		if raw == nil {
			missing++
			continue
		}
		strVal := fmt.Sprintf("%v", raw)
		uniqueSet[strVal] = struct{}{}
		if minVal == "" || strVal < minVal {
			minVal = strVal
		}
		if maxVal == "" || strVal > maxVal {
			maxVal = strVal
		}
		if numericValue, err := parseNumeric(strVal); err == nil {
			values = append(values, numericValue)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate column %s: %w", columnName, err)
	}

	rowCountQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s IS NOT NULL", quoteIdentifier(tableName), quoteIdentifier(columnName))
	var nonNullCount int
	if err := db.QueryRow(rowCountQuery).Scan(&nonNullCount); err != nil {
		return nil, fmt.Errorf("count non-null values for %s: %w", columnName, err)
	}

	missingRatio := 0.0
	if total := nonNullCount + missing; total > 0 {
		missingRatio = float64(missing) / float64(total) * 100
	}

	col := &ColumnProfile{
		Name:    columnName,
		Type:    columnType,
		Missing: missingRatio,
		Unique:  len(uniqueSet),
		Min:     minVal,
		Max:     maxVal,
	}
	if len(values) > 0 {
		col.Mean = formatFloat(mean(values))
		col.Median = formatFloat(median(values))
		if len(values) == 1 {
			col.Min = formatFloat(values[0])
			col.Max = formatFloat(values[0])
		}
	} else {
		col.Mean = ""
		col.Median = ""
	}
	if col.Min == "" && len(uniqueSet) > 0 {
		for key := range uniqueSet {
			col.Min = key
			break
		}
	}
	if col.Max == "" && len(uniqueSet) > 0 {
		maxKey := ""
		for key := range uniqueSet {
			if maxKey == "" || key > maxKey {
				maxKey = key
			}
		}
		col.Max = maxKey
	}
	return col, nil
}

func quoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func parseNumeric(s string) (float64, error) {
	return strconvToFloat(s)
}

func strconvToFloat(s string) (float64, error) {
	var out float64
	_, err := fmt.Sscanf(s, "%f", &out)
	return out, err
}

func mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func median(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	clone := append([]float64(nil), values...)
	sort.Float64s(clone)
	mid := len(clone) / 2
	if len(clone)%2 == 0 {
		return (clone[mid-1] + clone[mid]) / 2
	}
	return clone[mid]
}

func formatFloat(v float64) string {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return "0"
	}
	if math.Abs(v) >= 1e6 || (math.Abs(v) > 0 && math.Abs(v) < 1e-6) {
		return fmt.Sprintf("%.6f", v)
	}
	return fmt.Sprintf("%.0f", v)
}
