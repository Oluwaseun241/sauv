package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type PostgresDB struct {
	conn *pgx.Conn
}

func (pg *PostgresDB) Connect(dbUrl string) error {
	var err error
	pg.conn, err = pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Connected")
	return nil
}

func (pg *PostgresDB) Close() {
	if pg.conn != nil {
		pg.conn.Close(context.Background())
		fmt.Println("Connection closed")
	}
}

func (pg PostgresDB) QueryBackup() error {
	// query the database to get backup
	return nil
}

func (pg *PostgresDB) Backup(dbUrl, destination string) error {
	err := pg.Connect(dbUrl)
	if err != nil {
		return err
	}
	defer pg.Close()

	// Create or open the backup file
	f, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("Unable to create backup file: %v", err)
	}
	defer f.Close()

	connForTables, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		return fmt.Errorf("Unable to connect for table fetching: %v", err)
	}
	defer connForTables.Close(context.Background())

	rows, err := connForTables.Query(context.Background(), "SELECT table_name FROM information_schema.tables WHERE table_schema='public'")
	if err != nil {
		return fmt.Errorf("Error fetching table names: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return fmt.Errorf("Error scanning table name: %v", err)
		}

		err = pg.dumpTableSchema(dbUrl, f, tableName)
		if err != nil {
			return fmt.Errorf("Error dumping schema for table %s: %v", tableName, err)
		}

		err = pg.dumpTableData(dbUrl, f, tableName)
		if err != nil {
			return fmt.Errorf("Error dumping data for table %s: %v", tableName, err)
		}
	}

	if rows.Err() != nil {
		return fmt.Errorf("Error iterating over tables: %v", rows.Err())
	}

	fmt.Println("Backup successful, saved to", destination)
	return nil
}

func (pg *PostgresDB) dumpTableSchema(dbUrl string, f *os.File, tableName string) error {
	// Use a separate connection for dumping schema
	connForSchema, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		return fmt.Errorf("Unable to connect for schema fetching: %v", err)
	}
	defer connForSchema.Close(context.Background())

	schemaQuery := fmt.Sprintf(
		"SELECT 'CREATE TABLE ' || table_name || ' (' || array_to_string(array_agg(column_name || ' ' || data_type), ', ') || ');' FROM information_schema.columns WHERE table_name='%s' GROUP BY table_name;",
		tableName,
	)
	row := connForSchema.QueryRow(context.Background(), schemaQuery)

	var createTableStmt string
	if err := row.Scan(&createTableStmt); err != nil {
		return fmt.Errorf("Error fetching schema for table %s: %v", tableName, err)
	}

	_, err = f.WriteString(createTableStmt + ";\n\n")
	if err != nil {
		return fmt.Errorf("Error writing schema for table %s: %v", tableName, err)
	}
	return nil
}

func (pg *PostgresDB) dumpTableData(dbUrl string, f *os.File, tableName string) error {
	offset := 0
	limit := 100

	connForData, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		return fmt.Errorf("Unable to connect for data fetching: %v", err)
	}
	defer connForData.Close(context.Background())

	for {
		dataQuery := fmt.Sprintf("SELECT * FROM %s LIMIT %d OFFSET %d", tableName, limit, offset)
		dataRows, err := connForData.Query(context.Background(), dataQuery)
		if err != nil {
			return fmt.Errorf("Error querying data from table %s: %v", tableName, err)
		}

		colCount := len(dataRows.FieldDescriptions())

		rowCount := 0
		for dataRows.Next() {
			rowValues, err := dataRows.Values()
			if err != nil {
				return fmt.Errorf("Error scanning row data: %v", err)
			}

			insertStmt := fmt.Sprintf("INSERT INTO %s VALUES (", tableName)
			for i, val := range rowValues {
				if val == nil {
					insertStmt += "NULL"
				} else {
					switch v := val.(type) {
					case string:
						insertStmt += fmt.Sprintf("'%s'", v)
					case int, int64, float64:
						insertStmt += fmt.Sprintf("'%s'", v)
					default:
						insertStmt += fmt.Sprintf("'%s'", v)
					}
				}
				if i < colCount-1 {
					insertStmt += ", "
				}
			}
			insertStmt += ");\n"

			_, err = f.WriteString(insertStmt)
			if err != nil {
				return fmt.Errorf("Error writing data for table %s: %v", tableName, err)
			}
			rowCount++
		}

		dataRows.Close()
		_, err = f.WriteString("\n\n")
		if err != nil {
			return fmt.Errorf("Error writing data for table %s: %v", tableName, err)
		}

		if rowCount < limit {
			break
		}

		offset += limit
	}
	return nil
}
