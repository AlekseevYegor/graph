package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"graphs/constant"
	"graphs/entity"
	"graphs/entity/postgre"
	xmlentity "graphs/entity/xml"
	"graphs/repository/postges"
	"graphs/repository/receiver"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)

	// init DB
	db, err := connectDB()
	if err != nil {
		fmt.Printf("Error connect to db: %v\n", err)
		return
	}
	defer db.Close()

	graphRepo := postges.NewGraphRepo(db)
	err = downloadXMLGraphToDB(ctx, graphRepo)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Check if graph has cycle via DFS(Depth First Search).
	//Recursive query - add to path viewed edges, and check if we already pass the edge.
	cyclePath, err := graphRepo.GetGraphCycle(ctx)
	if err != nil {
		fmt.Printf("Error GetGraphCycle: %v\n", err)
		return
	}

	if len(cyclePath) > 0 {
		fmt.Printf("Found Cycle in graph: %v\n", cyclePath)
	} else {
		fmt.Println("Cycle in graph not found.")
	}

	// read graph from DB and make graph structure
	graphDB, err := graphRepo.GetGraph(ctx)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	// Start input message listener
	receiver.Receive(ctx, entity.NewGraph(*graphDB))
}

func setupGracefulShutdown(stop func()) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		fmt.Println("Got Interrupt signal")
		stop()
	}()
}

func connectDB() (*sqlx.DB, error) {
	var (
		host     = viper.GetString("DB_HOST")
		port     = viper.GetString("DB_PORT")
		dbName   = viper.GetString("DB_NAME")
		user     = viper.GetString("DB_USER")
		password = viper.GetString("DB_PASSWORD")
		schema   = viper.GetString("DB_SCHEMA")
		ssl      = viper.GetBool("SSL_MODE")
	)

	address := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?search_path=%s&sslmode=", user, password, host, port, dbName, schema)
	if !ssl {
		address += "disable"
	} else {
		address += "require"
	}

	db, err := sql.Open(constant.DBDriverName, address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres DB")
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping to postgres DB")
	}

	fmt.Println("checking DB migrations")

	if err := applySchemaMigrationWithDatabaseInstance(constant.DBDriverName, db); err != nil {
		return nil, fmt.Errorf("failed to migrate postgres DB schema")
	}

	fmt.Println("DB connected")

	return sqlx.NewDb(db, "postgres"), nil
}

// applySchemaMigrationWithDatabaseInstance creates and migrates db schema versions; based on db driver instance
// Doesn't close db connection though this is fully caller's responsibility
func applySchemaMigrationWithDatabaseInstance(databaseName string, db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", databaseName, driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != migrate.ErrNoChange {
		return err
	}

	_, _, err = m.Version()
	return err

}

func downloadXMLGraphToDB(ctx context.Context, graphRepo *postges.GraphRepo) error {
	xmlFilePath := "graph.xml"

	// Read XML
	body, err := os.ReadFile(xmlFilePath)
	if err != nil {
		return fmt.Errorf("error reading XML file: %w", err)
	}

	// Used standard library for XML parsing
	// If needed more performance or/and xsd validation I would rather use library https://github.com/lestrrat-go/libxml2
	var graphXML xmlentity.Graph
	err = xml.Unmarshal(body, &graphXML)
	if err != nil {
		return fmt.Errorf("error unmarshalling XML: %w", err)

	}

	err = graphXML.Validate()
	if err != nil {
		return fmt.Errorf("error validate graph XML: %w", err)
	}

	fmt.Printf("Graph ID: %s\n", graphXML.ID)
	fmt.Printf("Graph Name: %s\n", graphXML.Name)

	graph := postgre.NewGraph(graphXML)

	//save graph to DB in Transactions
	err = graphRepo.UpsertGraph(ctx, graph)
	if err != nil {
		return fmt.Errorf("error upsert graph into DB: %w", err)
	}

	return nil
}
