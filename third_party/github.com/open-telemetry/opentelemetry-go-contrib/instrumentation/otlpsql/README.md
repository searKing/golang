# otlpsql

OpenTelemetry SQL database driver wrapper.

Add an otlpsql wrapper to your existing database code to instrument the
interactions with the database.

## installation

go get -u github.com/searKing/opentelemetry-go-contrib/instrumentation/database/otlpsql

## initialize

To use otlpsql with your application, register an otlpsql wrapper of a database
driver as shown below.

Example1:
```go
    import (
        "database/sql"
        _ "github.com/mattn/go-sqlite3"
        "github.com/searKing/opentelemetry-go-contrib/instrumentation/database/otlpsql"
    )
    // Register our sqlite3-otlp wrapper for the provided SQLite3 driver.
	// "sqlite3-otlp" must not be registered, set in func init(){} as recommended.
	sql.Register("sqlite3-otlp", otlpsql.Wrap(&sqlite3.SQLiteDriver{}, otlpsql.WithAllWrapperOptions()))

	// Connect to a SQLite3 database using the otlpsql driver wrapper.
	db, err := sql.Open("sqlite3-otlp", "resource.db")
	if err != nil {
		log.Fatal(err)
	}
```
Example2:
```go
import (
    _ "github.com/mattn/go-sqlite3"
    "github.com/searKing/opentelemetry-go-contrib/instrumentation/database/otlpsql"
)

var (
    driverName string
    err        error
    db         *sql.DB
)

// Register our otlpsql wrapper for the provided SQLite3 driver.
driverName, err = otlpsql.Register("sqlite3", otlpsql.WithAllTraceOptions(), otlpsql.WithInstanceName("resources"))
if err != nil {
    log.Fatalf("unable to register our otlpsql driver: %v\n", err)
}

// Connect to a SQLite3 database using the otlpsql driver wrapper.
db, err = sql.Open(driverName, "resource.db")
```

A more explicit and alternative way to bootstrap the otlpsql wrapper exists as
shown below. This will only work if the actual database driver has its driver
implementation exported.

Example:
```go
import (
    sqlite3 "github.com/mattn/go-sqlite3"
    "github.com/searKing/opentelemetry-go-contrib/instrumentation/database/otlpsql"
)

var (
    driver driver.Driver
    err    error
    db     *sql.DB
)

// Explicitly wrap the SQLite3 driver with otlpsql.
driver = otlpsql.Wrap(&sqlite3.SQLiteDriver{})

// Register our otlpsql wrapper as a database driver.
sql.Register("otlpsql-sqlite3", driver)

// Connect to a SQLite3 database using the otlpsql driver wrapper.
db, err = sql.Open("otlpsql-sqlite3", "resource.db")
```

Projects providing their own abstractions on top of database/sql/driver can also
wrap an existing driver.Conn interface directly with otlpsql.

Example:
```go
import "github.com/searKing/opentelemetry-go-contrib/instrumentation/database/otlpsql"

func GetConn(...) driver.Conn {
    // Create custom driver.Conn.
    conn := initializeConn(...)

    // Wrap with otlpsql.
    return otlpsql.WrapConn(conn, otlpsql.WithAllTraceOptions())    
}
```

Finally database drivers that support the new (Go 1.10+) driver.Connector
interface can be wrapped directly by otlpsql without the need for otlpsql to
register a driver.Driver.

Example:
```go
import(
    "github.com/searKing/opentelemetry-go-contrib/instrumentation/database/otlpsql"
    "github.com/lib/pq"
)

var (
    connector driver.Connector
    err       error
    db        *sql.DB
)

// Get a database driver.Connector for a fixed configuration.
connector, err = pq.NewConnector("postgres://user:passt@host:5432/db")
if err != nil {
    log.Fatalf("unable to create our postgres connector: %v\n", err)
}

// Wrap the driver.Connector with otlpsql.
connector = otlpsql.WrapConnector(connector, otlpsql.WithAllWrapperOptions())

// Use the wrapped driver.Connector.
db = sql.OpenDB(connector)
```

## metrics

Next to tracing, otlpsql also supports Opentelemetrys stats automatically.

From Go 1.11 and up, otlpsql also has the ability to record database connection
pool details. Use the `RecordStats` function and provide a `*sql.DB` to record
details on, as well as the required record interval.

```go
// Connect to a SQLite3 database using the otlpsql driver wrapper.
db, err = sql.Open("otlpsql-sqlite3", "resource.db")

// Record DB stats every 5 seconds until we exit.
defer otlpsql.RecordStats(db, 5 * time.Second)()
```

## Recorded metrics

| Metric                 | Search suffix          | Additional tags            |
|------------------------|------------------------|----------------------------|
| Latency in milliseconds| "go_sql_client_latency_milliseconds"|"method", "error", "status" |

If using RecordStats:

| Metric                                                   | Search suffix                                |
|----------------------------------------------------------|----------------------------------------------|
| Number of open connections                               | "go_sql_connections_open"                 |
| Number of idle connections                               | "go_sql_connections_idle"                 |
| Number of active connections                             | "go_sql_connections_active"               |
| Total number of connections waited for                   | "go_sql_connections_wait_count"           |
| Total time blocked waiting for new connections           | "go_sql_connections_wait_duration_milliseconds"        |
| Total number of closed connections by SetMaxIdleConns    | "go_sql_connections_idle_closed"     |
| Total number of closed connections by SetConnMaxLifetime | "go_sql_connections_lifetime_closed" |

## jmoiron/sqlx

If using the `sqlx` library with named queries you will need to use the
`sqlx.NewDb` function to wrap an existing `*sql.DB` connection. Do not use the
`sqlx.Open` and `sqlx.Connect` methods.
`sqlx` uses the driver name to figure out which database is being used. It uses
this knowledge to convert named queries to the correct bind type (dollar sign,
question mark) if named queries are not supported natively by the
database. Since otlpsql creates a new driver name it will not be recognized by
sqlx and named queries will fail.

Use one of the above methods to first create a `*sql.DB` connection and then
create a `*sqlx.DB` connection by wrapping the `*sql.DB` like this:

```go
    // Register our otlpsql wrapper for the provided Postgres driver.
    driverName, err := otlpsql.Register("postgres", otlpsql.WithAllTraceOptions())
    if err != nil { ... }

    // Connect to a Postgres database using the otlpsql driver wrapper.
    db, err := sql.Open(driverName, "postgres://localhost:5432/my_database")
    if err != nil { ... }

    // Wrap our *sql.DB with sqlx. use the original db driver name!!!
    dbx := sqlx.NewDB(db, "postgres")
```

## context

To really take advantage of otlpsql, all database calls should be made using the
*Context methods. Failing to do so will result in many orphaned otlpsql traces
if the `AllowRoot` TraceOption is set to true. By default AllowRoot is disabled
and will result in otlpsql not tracing the database calls if context or parent
spans are missing.

| Old            | New                   |
|----------------|-----------------------|
| *DB.Begin      | *DB.BeginTx           |
| *DB.Exec       | *DB.ExecContext       |
| *DB.Ping       | *DB.PingContext       |
| *DB.Prepare    | *DB.PrepareContext    |
| *DB.Query      | *DB.QueryContext      |
| *DB.QueryRow   | *DB.QueryRowContext   |
|                |                       |
| *Stmt.Exec     | *Stmt.ExecContext     |
| *Stmt.Query    | *Stmt.QueryContext    |
| *Stmt.QueryRow | *Stmt.QueryRowContext |
|                |                       |
| *Tx.Exec       | *Tx.ExecContext       |
| *Tx.Prepare    | *Tx.PrepareContext    |
| *Tx.Query      | *Tx.QueryContext      |
| *Tx.QueryRow   | *Tx.QueryRowContext   |

Example:
```go
func (s *svc) GetDevice(ctx context.Context, id int) (*Device, error) {
    // Assume we have instrumented our service transports and ctx holds a span.
    var device Device
    if err := s.db.QueryRowContext(
        ctx, "SELECT * FROM device WHERE id = ?", id,
        ).Scan(&device); err != nil {
        return nil, err
    }
    return device
}
```
