package connect

import(
	"context"
	"os"
	"log"
	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"
	// "rhymald/mag-zeta/base"
	"fmt"
)

func ConnectCacheDB() []*pgx.Conn {
	config, err := pgx.ParseConfig(os.Getenv("CACHEDB_WRITER_URL"))
	if err != nil {
			log.Fatal(err)
	}
	config.RuntimeParams["application_name"] = "$ mag_cached_grid"
	writer, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
			log.Fatal(err)
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
		"stepid UUID PRIMARY KEY default gen_random_uuid()",
		"inserted_at TIMESTAMP default current_timestamp()",
		"id TEXT",
		"t INT",
		"x INT",
		"y INT",
		"INDEX position (t, x, y)",
		"ttl_expire_after = '61 seconds', ttl_job_cron = '*/1 * * * *'",
	)
	if _, err := tx.Exec(ctx, query); err != nil { log.Fatal(err) ; return err }
	return nil
}

func WriteTrace(writer *pgx.Conn, id string, trace *map[int][3]int) error {
	err := crdbpgx.ExecuteTx(context.Background(), writer, pgx.TxOptions{}, func(tx pgx.Tx) error { return writeAllTrace(context.Background(), tx, id, *trace) })
	max, buffer := 0, *trace
	for ts, _ := range buffer { if ts > max { max = ts }}
	latest := buffer[max]
	newTrace := make(map[int][3]int) ; newTrace[max] = latest
	*trace = newTrace
	return err
}
func writeAllTrace(ctx context.Context, tx pgx.Tx, id string, trace map[int][3]int) error {
	query, first := "UPSERT INTO eua (id, t, x, y) VALUES", true
	for ts, rxy := range trace {
		if first {
			query = fmt.Sprintf("%s ('%s', '%d', '%d', '%d')", query, id, ts, rxy[1], rxy[2])
			first = false 
		} else {
			query = fmt.Sprintf("%s, ('%s', '%d', '%d', '%d')", query, id, ts, rxy[1], rxy[2])
		}
	}
	query = fmt.Sprintf("%s;", query)
	log.Println(query)
	if _, err := tx.Exec(ctx, query); err != nil { log.Fatal(err) ; return err }
	return nil
}

func WriteChunk(writer *pgx.Conn, chunk map[string][][3]int) error {
	err := crdbpgx.ExecuteTx(context.Background(), writer, pgx.TxOptions{}, func(tx pgx.Tx) error { return writeChunk(context.Background(), tx, chunk) })
	return err
}
func writeChunk(ctx context.Context, tx pgx.Tx, chunk map[string][][3]int) error {
	query, first := "UPSERT INTO eua (id, t, x, y) VALUES", true
	for id, char := range chunk { for _, txy := range char {
		if first {
			query = fmt.Sprintf("%s ('%s', '%d', '%d', '%d')", query, id, txy[0], txy[1], txy[2])
			first = false 
		} else {
			query = fmt.Sprintf("%s, ('%s', '%d', '%d', '%d')", query, id, txy[0], txy[1], txy[2])
		}
	}}
	query = fmt.Sprintf("%s;", query)
	log.Println(query)
	if _, err := tx.Exec(ctx, query); err != nil { log.Fatal(err) ; return err }
	return nil
}
