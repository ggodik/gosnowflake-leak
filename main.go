package main

import (
	"context"
	"database/sql/driver"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/snowflakedb/gosnowflake"
)

const (
	queryStmt = "SELECT seq4() A, uniform(1, 10, RANDOM(12)) B, uniform(1, 100, RANDOM(121)) C, randstr(15, random()), randstr(25, random()), randstr(35, random())FROM TABLE(GENERATOR(ROWCOUNT=>%d))"
)

var (
	rowsCount = flag.Int("rows", 100, "number of rows to generate")
)

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer func() {
		if err := mem(); err != nil {
			log.Fatal(err)
		}
	}()

	go printMemUsage(ctx, 5)

	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Fatal("DSN is empty")
	}
	db, err := getDB(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	batches, err := query(db, *rowsCount)
	if err != nil {
		log.Fatal(err)
	}

	for _, b := range batches {
		recs, err := b.Fetch(ctx)
		if err != nil {
			log.Fatal(err)
		}

		for _, rec := range *recs {
			rec.Release()
		}
	}

	log.Println("DONE")
}

func getDB(ctx context.Context, dsn string) (driver.Conn, error) {
	ctx = gosnowflake.WithStreamDownloader(gosnowflake.WithArrowBatches(ctx))
	scfg, err := gosnowflake.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}

	d := gosnowflake.SnowflakeDriver{}
	return d.OpenWithConfig(ctx, *scfg)
}

// note - snowflake is a pain in the ass here
// only way to get Arrow out of it is to use a prepared statement that returns a Result interface that can be cast
func query(db driver.Conn, rows int) ([]*gosnowflake.ArrowBatch, error) {
	q := fmt.Sprintf(queryStmt, rows)
	log.Printf("query=%s", q)
	st, err := db.Prepare(q)
	if err != nil {
		return nil, err
	}
	defer st.Close()

	res, err := st.Query([]driver.Value{})
	if err != nil {
		return nil, err
	}
	defer res.Close()

	sres, ok := res.(gosnowflake.SnowflakeRows)
	if !ok {
		return nil, fmt.Errorf("not snowflake")
	}

	return sres.GetArrowBatches()
}

func printMemUsage(ctx context.Context, interval int) {
	ticker := time.NewTicker(time.Second * time.Duration(interval))

	go func() {
		for {
			select {
			case <-ticker.C:
				var m runtime.MemStats
				runtime.ReadMemStats(&m)
				// For info on each, see: https://golang.org/pkg/runtime/#MemStats
				log.Printf("memory usage alloc=%s totalAlloc=%s system=%s numGC=%d",
					byteCountSI(int64(m.Alloc)), byteCountSI(int64(m.TotalAlloc)), byteCountSI(int64(m.Sys)), m.NumGC)
			case <-ctx.Done():
				return
			}
		}
	}()
}

const (
	unit = 1000
)

// ByteCountSI format signed int64 into human readable largest units and shorten
func byteCountSI(b int64) string {
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func mem() error {
	name := fmt.Sprintf("profile-%s.mem", time.Now().Format("20060102150405"))
	log.Printf("run the following command to see memory profiling\ngo tool pprof -http :8000 %s", name)
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return pprof.WriteHeapProfile(f)
}
