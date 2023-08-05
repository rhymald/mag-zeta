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

func ConnectCacheDB() *pgx.Conn {
	config, err := pgx.ParseConfig(os.Getenv("CACHEDB_URL"))
	if err != nil {
			log.Fatal(err)
	}
	config.RuntimeParams["application_name"] = "$ docs_simplecrud_gopgx"
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
			log.Fatal(err)
	}
	// defer conn.Close(context.Background())
	err = crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error { return initTable(context.Background(), tx) })
	return conn
}

func initTable(ctx context.Context, tx pgx.Tx) error {
	// Dropping existing table if it exists

	log.Println(" => Drop existing table:")
	if _, err := tx.Exec(ctx, "DROP TABLE IF EXISTS eua;"); err != nil {
			return err
	}

	// Create the accounts table
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
		"ttl_expire_after = '5 minutes', ttl_job_cron = '*/2 * * * *'",
	)
	if _, err := tx.Exec(ctx, query); err != nil { log.Println(" => ERROR[Failed to create table]:", err) ; return err }
	return nil
}
