package connect

import(
	"context"
	"os"
	"log"
	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"
	"rhymald/mag-zeta/base"
	"rhymald/mag-zeta/play"
	"fmt"
)

const AcceptableReadLatency = 200

func ConnectCacheDB() []*pgx.Conn {
	con_string :=  os.Getenv("CACHEDB_WRITER_URL")
	if len(con_string) <= 5 { con_string = "postgresql://root@localhost:26257/grid?sslmode=disable" }
	config, err := pgx.ParseConfig(con_string)
	if err != nil { log.Fatal(err) }
	config.RuntimeParams["application_name"] = "$ mag_cached_grid"
	writer, err := pgx.ConnectConfig(context.Background(), config)
	for err != nil {
		log.Println("Waiting for DB:", err)
		writer, err = pgx.ConnectConfig(context.Background(), config)
		base.Wait(1618)
	}
	// defer conn.Close(context.Background())
	err = crdbpgx.ExecuteTx(context.Background(), writer, pgx.TxOptions{}, func(tx pgx.Tx) error { return initTable(context.Background(), tx) })
	return []*pgx.Conn{ writer } 
}

func initTable(ctx context.Context, tx pgx.Tx) error {
	// Dropping existing table if it exists
	log.Println(" => Drop existing table:")
	if _, err := tx.Exec(ctx, "DROP TABLE IF EXISTS eua;"); err != nil {
			return err
	}
	// Create the grid table
	log.Println(" => Creating grid table: eua...")
	query := fmt.Sprintf("CREATE TABLE %s (%s, %s, %s, %s, %s, %s, %s) WITH (%s);",
		"eua",
		"stepid STRING PRIMARY KEY AS (CONCAT(id, '@', CAST(t AS STRING))) STORED",// default gen_random_uuid()",
		"inserted_at TIMESTAMP default current_timestamp()",
		"id TEXT",
		"t INT",
		"x FLOAT",
		"y FLOAT",
		"INDEX position (t, x, y)",
		"ttl_expire_after = '61 seconds', ttl_job_cron = '*/1 * * * *'",
	)
	if _, err := tx.Exec(ctx, query); err != nil { log.Fatal(err) ; return err }
	return nil
}

// func WriteTrace(writer *pgx.Conn, id string, trace *map[int][3]int) error {
// 	err := crdbpgx.ExecuteTx(context.Background(), writer, pgx.TxOptions{}, func(tx pgx.Tx) error { return writeAllTrace(context.Background(), tx, id, *trace) })
// 	max, buffer := 0, *trace
// 	for ts, _ := range buffer { if ts > max { max = ts }}
// 	latest := buffer[max]
// 	newTrace := make(map[int][3]int) ; newTrace[max] = latest
// 	*trace = newTrace
// 	return err
// }
// func writeAllTrace(ctx context.Context, tx pgx.Tx, id string, trace map[int][3]int) error {
// 	query, first := "UPSERT INTO eua (id, t, x, y) VALUES", true
// 	for ts, rxy := range trace {
// 		if first {
// 			query = fmt.Sprintf("%s ('%s', '%d', '%d', '%d')", query, id, ts, rxy[1], rxy[2])
// 			first = false 
// 		} else {
// 			query = fmt.Sprintf("%s, ('%s', '%d', '%d', '%d')", query, id, ts, rxy[1], rxy[2])
// 		}
// 	}
// 	query = fmt.Sprintf("%s;", query)
// 	// log.Println(query)
// 	if _, err := tx.Exec(ctx, query); err != nil { log.Fatal(err) ; return err }
// 	return nil
// }

func WriteChunk(writer *pgx.Conn, chunk map[string][][3]int) error {
	err := crdbpgx.ExecuteTx(context.Background(), writer, pgx.TxOptions{}, func(tx pgx.Tx) error { return writeChunk(context.Background(), tx, chunk) })
	return err
}
func writeChunk(ctx context.Context, tx pgx.Tx, chunk map[string][][3]int) error {
	query, first := "UPSERT INTO eua (id, t, x, y, inserted_at) VALUES", true
	for id, char := range chunk { for _, txy := range char {
		if first {
			query = fmt.Sprintf("%s ('%s', '%d', '%d', '%d', current_timestamp())", query, id, txy[0], txy[1], txy[2])
			first = false 
		} else {
			query = fmt.Sprintf("%s, ('%s', '%d', '%d', '%d', current_timestamp())", query, id, txy[0], txy[1], txy[2])
		}
	}}
	query = fmt.Sprintf("%s;", query)
	// log.Println(query)
	if _, err := tx.Exec(ctx, query); err != nil { log.Fatal(err) ; return err }
	return nil
}

func ReadRound(writer *pgx.Conn, x, y, r, t int) map[int][]string {
	// Read the balance.
	// var list string
	// if err := tx.QueryRow(ctx,
	// 		"SELECT id FROM uae WHERE x < $1 AND x > $2 AND y < $3 AND y > $4 AND SQRT(SQR($5-x)+SQR($6-y)) < $7", x+r, x-r, y+r, y-r, x, y, r).Scan(&list); err != nil {
	// 		return err
	// }
	list, err := writer.Query(context.Background(), "SELECT id, t, date_part('epoch', inserted_at) FROM eua WHERE x < $1 AND x > $2 AND y < $3 AND y > $4 AND SQRT(POW(($5-x),2)+POW(($6-y),2)) < $7 AND t < $8+3 AND t > $8-3", x+r, x-r, y+r, y-r, x, y, r, t)
	defer list.Close()
	if err != nil {log.Fatal(err)}
	now := play.TAxis()
	log.Printf("Characters within %d from [%d,%d] at %d:\n", r, x, y, now) // from play
	start:=base.StartEpoch/1000000+base.Epoch()
	buffer :=  make(map[int][]string)
	for list.Next() {
		var id string
		var t int
		var ts float64
		if err := list.Scan(&id, &t, &ts); err != nil { log.Fatal(err) }
		diff := int(1000*ts)-start
		if t == now {
			log.Printf(" => %3d: %s | %+d\n", t, id, diff)
		} else {
			log.Printf("    %3d: %s | %+d\n", t, id, diff)
		}
		// found := map[string][2]int{ id: [2]int{ t, diff }}
		if diff < AcceptableReadLatency && diff > -AcceptableReadLatency { 
			renew := buffer
			renewList := renew[t]
			renewList = append(renewList, id)
			renew[t] = renewList
			buffer = renew
			// buffer[t] = append(buffer[t], id) 
		}
	}

	// // Perform the transfer.
	// log.Printf("Transferring funds from account with ID %s to account with ID %s...", from, to)
	// if _, err := tx.Exec(ctx,
	// 		"UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, from); err != nil {
	// 		return err
	// }
	// if _, err := tx.Exec(ctx,
	// 		"UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, to); err != nil {
	// 		return err
	// }
	return buffer
}

// func printBalances(conn *pgx.Conn) error {
// 	rows, err := conn.Query(context.Background(), "SELECT id, balance FROM accounts")
// 	if err != nil {
// 			log.Fatal(err)
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 			var id uuid.UUID
// 			var balance int
// 			if err := rows.Scan(&id, &balance); err != nil {
// 					log.Fatal(err)
// 			}
// 			log.Printf("%s: %d\n", id, balance)
// 	}
// 	return nil
// }
