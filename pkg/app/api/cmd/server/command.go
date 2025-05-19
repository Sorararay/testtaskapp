package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq" // PostgreSQLドライバをインポート
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
	"golang.org/x/sync/errgroup"

	"github.com/gihyodocker/taskapp/pkg/app/api/handler"
	"github.com/gihyodocker/taskapp/pkg/cli"
	"github.com/gihyodocker/taskapp/pkg/config"
	"github.com/gihyodocker/taskapp/pkg/db"
	"github.com/gihyodocker/taskapp/pkg/repository"
	"github.com/gihyodocker/taskapp/pkg/server"
)

type command struct {
	port        int
	gracePeriod time.Duration
	configFile  string
}

func NewCommand() *cobra.Command {
	c := &command{
		port:        8180,
		gracePeriod: 5 * time.Second,
	}
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start up the api server",
		RunE:  cli.WithContext(c.execute),
	}
	cmd.Flags().IntVar(&c.port, "port", c.port, "The port number used to run HTTP api.")
	cmd.Flags().DurationVar(&c.gracePeriod, "grace-period", c.gracePeriod, "How long to wait for graceful shutdown.")
	cmd.Flags().StringVar(&c.configFile, "config-file", c.configFile, "The path to the config file.")

	cmd.MarkFlagRequired("config-file")
	return cmd
}

func (c *command) execute(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	appConfig, err := config.LoadConfigFile(c.configFile)
	if err != nil {
		slog.Error("failed to load api configuration",
			slog.String("config-file", c.configFile),
			err,
		)
		return err
	}
	// Open PostgreSQL connection
	dbConn, err := createPostgreSQL(*appConfig.Database)
	if err != nil {
		slog.Error("failed to open PostgreSQL connection", err)
		return err
	}

	// Initialize repositories
	taskRepo := repository.NewTask(dbConn)

	// Handlers
	taskHandler := handler.NewTask(taskRepo)

	options := []server.Option{
		server.WithGracePeriod(c.gracePeriod),
	}
	httpServer := server.NewHTTPServer(c.port, options...)
	httpServer.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	httpServer.Put("/api/tasks/{id}", taskHandler.Update)
	httpServer.Delete("/api/tasks/{id}", taskHandler.Delete)
	httpServer.Get("/api/tasks/{id}", taskHandler.Get)
	httpServer.Post("/api/tasks", taskHandler.Create)
	httpServer.Get("/api/tasks", taskHandler.List)

	group.Go(func() error {
		return httpServer.Serve(ctx)
	})

	if err := group.Wait(); err != nil {
		slog.Error("failed while running", err)
		return err
	}
	return nil
}

// createPostgreSQL は PostgreSQL データベースへの接続を確立します。
func createPostgreSQL(conf config.Database) (*sql.DB, error) {
	// PostgreSQL用のDSN (Data Source Name) を組み立てます。
	// conf.Host には "ホスト名:ポート番号" (例: "localhost:5432") またはホスト名のみを指定します。
	// 環境やDSNの形式に合わせて調整してください。
	// 一般的なURI形式のDSN:
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		conf.Username,
		conf.Password,
		conf.Host, // 例: "localhost:5432" や "your-postgres-server.com"
		conf.DBName,
		conf.SSLMode,
	)

	// db.OpenDB が期待する db.Datasource オブジェクトを作成します。
	// db.NewDatasource がドライバ名とDSN文字列から Datasource を作成すると仮定します。
	pgDatasource := db.NewPgDatasource("postgres", dsn)

	options := []db.Option{
		db.WithMaxIdleConns(conf.MaxIdleConns),
		db.WithMaxOpenConns(conf.MaxOpenConns),
		db.WithConnMaxLifetime(conf.ConnMaxLifetime),
	}

	// db.OpenDB が第1引数にドライバ名 ("postgres")、第2引数にDSN文字列を取ると仮定しています。
	// 修正: db.OpenDB は第1引数に db.Datasource オブジェクトを取ります。
	return db.OpenDB(pgDatasource, options...)
}
