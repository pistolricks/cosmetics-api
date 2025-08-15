package main

import (
	"context"
	"database/sql"
	"expvar"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/joho/godotenv"
	"github.com/pistolricks/cosmetics-api/internal/chromium"
	"github.com/pistolricks/cosmetics-api/internal/data"
	"github.com/pistolricks/cosmetics-api/internal/mailer"
	"github.com/pistolricks/cosmetics-api/internal/riman"
	"github.com/pistolricks/cosmetics-api/internal/shopify"
	v2 "github.com/pistolricks/cosmetics-api/internal/v2"
	"github.com/pistolricks/cosmetics-api/internal/vcs"
)

type Envars struct {
	StoreName      string
	RimanStoreName string
	LoginUrl       string
	Username       string
	Password       string
	ShopifyToken   string
	ShopifyKey     string
	ShopifySecret  string
	Token          string
}

var (
	version = vcs.Version()
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	cors struct {
		trustedOrigins []string
	}
}

type CreditCard struct {
}

type application struct {
	config   config
	logger   *slog.Logger
	envars   *Envars
	page     *rod.Page
	browser  *rod.Browser
	cookies  []*proto.NetworkCookie
	models   data.Models
	riman    riman.Riman
	shopify  shopify.ShopClient
	chromium chromium.ChromeConnector
	v2       v2.Api
	client   *riman.Client
	session  *riman.Session
	mailer   mailer.Mailer
	wg       sync.WaitGroup
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	storeName := os.Getenv("STORE_NAME")
	if storeName == "" {
		fmt.Println("missing store name")
		return
	}

	rimanStoreName := os.Getenv("RIMAN_STORE_NAME")
	if rimanStoreName == "" {
		fmt.Println("missing riman store name")
		return
	}

	loginUrl := os.Getenv("LOGIN_URL")
	if loginUrl == "" {
		fmt.Println("missing login url")
		return
	}

	username := os.Getenv("USERNAME")
	if username == "" {
		fmt.Println("missing username")
		return
	}

	password := os.Getenv("PASSWORD")
	if password == "" {
		fmt.Println("missing password")
		return
	}

	shopifyToken := os.Getenv("SHOPIFY_TOKEN")
	if shopifyToken == "" {
		fmt.Println("missing shop token")
		return
	}

	shopifyKey := os.Getenv("SHOPIFY_KEY")
	if shopifyKey == "" {
		fmt.Println("missing shop key")
		return
	}

	shopifySecret := os.Getenv("SHOPIFY_SECRET")
	if shopifySecret == "" {
		fmt.Println("missing shop secret")
		return
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		fmt.Println("missing token")
		return
	}

	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "PostgreSQL DSN")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "a7420fc0883489", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "e75ffd0a3aa5ec", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Greenlight <no-reply@greenlight.alexedwards.net>", "SMTP sender")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	ShopifyV2 := v2.ShopifyV2()
	chromeConfig := &chromium.ChromeConfig{Browser: chromium.ChromeBrowser()}
	shopClient := shopify.NewShopClient(shopify.ShopConfig),
		chromium.NewChromeConnector(chromeConfig),

		logger.Info("database connection pool established")

	expvar.NewString("version").Set(version)

	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))

	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))

	vars := &Envars{StoreName: storeName, RimanStoreName: rimanStoreName, LoginUrl: loginUrl, Username: username, Password: password, ShopifyToken: shopifyToken, ShopifyKey: shopifyKey, ShopifySecret: shopifySecret, Token: token}

	fmt.Println(vars)

	app := &application{
		config: cfg,
		logger: logger,
		envars: vars,
		models: data.NewModels(db),
		riman:  riman.NewRiman(db),

		chromium: chromium.NewChromeConnector(chromeConfig),
		mailer:   mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
