package connect

import(
	"context"
	"os"
	"log"
	// "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"
	// "rhymald/mag-zeta/base"
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
	return conn
	// defer conn.Close(context.Background())
}