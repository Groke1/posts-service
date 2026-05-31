package app

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"posts-service/internal/config"
	graphqlctrl "posts-service/internal/controller/graphql"
	commentinmemory "posts-service/internal/repository/comment/inmemory"
	commentpostgres "posts-service/internal/repository/comment/postgres"
	postinmemory "posts-service/internal/repository/post/inmemory"
	postpostgres "posts-service/internal/repository/post/postgres"
	"posts-service/internal/subscription"
	"posts-service/internal/usecase/comment"
	"posts-service/internal/usecase/post"
	"posts-service/migrations"
	"posts-service/pkg/db"
	"posts-service/pkg/generated/posts/graphql/generated"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func Run(logger *zap.Logger, cfg *config.Config) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	storageType, err := cfg.GetStorageType()
	if err != nil {
		return err
	}

	var (
		postRepoForPostService post.PostRepository
		postRepoForComment     comment.PostRepository
		commentRepo            comment.CommentRepository
		transactor             comment.Transactor
	)

	if storageType == config.Postgres {
		pgxcfg, err := pgxpool.ParseConfig(cfg.ConstructPostgresURL())
		if err != nil {
			logger.Fatal("can not create pgxpool cfg", zap.Error(err))
			return err
		}
		pgxcfg.MaxConns = 8
		pgxcfg.MinConns = 1
		pgxcfg.HealthCheckPeriod = 30 * time.Second
		pgxcfg.MaxConnLifetime = 0
		pgxcfg.MaxConnIdleTime = 5 * time.Minute

		dbPool, err := pgxpool.NewWithConfig(ctx, pgxcfg)
		if err != nil {
			logger.Error("can not create pgxpool", zap.Error(err))
			return err
		}

		defer dbPool.Close()
		migrations.SetupPostgres(dbPool, logger)
		postRepo := postpostgres.NewPostRepository(dbPool)
		postRepoForPostService = postRepo
		postRepoForComment = postRepo
		commentRepo = commentpostgres.NewCommentRepository(dbPool)
		transactor = db.NewTransactor(dbPool)
	} else {
		postRepo := postinmemory.NewPostRepository()
		postRepoForPostService = postRepo
		postRepoForComment = postRepo
		commentRepo = commentinmemory.NewCommentRepository()
		transactor = db.NewNoopTransactor()
	}
	publisher := subscription.NewCommentPubSub()

	postService := post.NewPostService(postRepoForPostService)
	commentService := comment.NewCommentService(commentRepo, postRepoForComment, transactor, publisher)

	resolver := graphqlctrl.NewResolver(postService, commentService, publisher)

	httpServer := runGraphQL(logger, cfg, resolver)

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	return httpServer.Shutdown(shutdownCtx)
}

func runGraphQL(logger *zap.Logger, cfg *config.Config, resolver *graphqlctrl.Resolver) *http.Server {
	srv := handler.New(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: resolver,
			},
		),
	)

	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.Use(extension.Introspection{})

	mux := http.NewServeMux()

	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", srv)

	httpServer := &http.Server{
		Addr:    ":" + cfg.GraphQL.Port,
		Handler: mux,
	}

	go func() {
		logger.Info("graphql server listening", zap.String("port", cfg.GraphQL.Port))

		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("graphql server error", zap.Error(err))
		}
	}()

	return httpServer
}
