package pkg

import (
	"context"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

func getMaxConnPool(ctx context.Context, url string) (int, error) {
	dbPool, err := pgxpool.New(ctx, url)
	if err != nil {
		return 0, err
	}
	defer dbPool.Close()

	var getMaxConn string
	query := `SHOW max_connections`
	err = dbPool.QueryRow(ctx, query).Scan(&getMaxConn)
	if err != nil {
		return 0, err
	}

	maxConn, err := strconv.ParseFloat(getMaxConn, 64)
	if err != nil {
		return 0, err
	}
	maxConnPool := int(math.Floor(maxConn * 0.9))
	if maxConnPool < 1 {
		maxConnPool = 1
	}

	return maxConnPool, nil
}

func CreateConnPool(ctx context.Context, url string) *pgxpool.Pool {

	connConf, err := pgxpool.ParseConfig(url)
	if err != nil {
		panic(err)
	}

	maxPool, err := getMaxConnPool(ctx, url)
	if err != nil {
		panic(err)
	}

	connConf.MaxConns = int32(maxPool)

	pool, err := pgxpool.NewWithConfig(ctx, connConf)

	if err != nil {
		log.Fatalln("Cannot connect database: ", err)
		os.Exit(1)
	}

	if err := pool.Ping(context.Background()); err == nil {
		log.Println("Database connected")
	} else {
		log.Fatalln("Cannot connect database: ", err)
		os.Exit(1)
	}

	return pool
}
