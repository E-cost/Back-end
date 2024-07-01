package main

import (
	"Ecost/internal/api/middleware/dns"
	"Ecost/internal/config"
	contactsHandler "Ecost/internal/contact/handler"
	"Ecost/internal/contact/mail"
	contactsRepository "Ecost/internal/contact/repository"
	service1 "Ecost/internal/contact/service"
	"Ecost/internal/contact/storage"
	"Ecost/internal/database/client/postgresql"
	"Ecost/internal/facades"
	service4 "Ecost/internal/facades/service"
	itemsHandler "Ecost/internal/item/handler"
	itemsRepository "Ecost/internal/item/repository"
	service2 "Ecost/internal/item/service"
	storage2 "Ecost/internal/item/storage"
	ordersHandler "Ecost/internal/order/handler"
	ordersRepository "Ecost/internal/order/reposiroty"
	service3 "Ecost/internal/order/service"
	storage3 "Ecost/internal/order/storage"
	service6 "Ecost/internal/redis/service"
	"Ecost/internal/utils/cors"
	yandexHandler "Ecost/internal/yandex/handler"
	service5 "Ecost/internal/yandex/service"
	"Ecost/pkg/logging"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	config2 "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
)

func NewRedisClient(cfg *config.Config, logger *logging.Logger) *redis.Client {
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}

	client := redis.NewClient(options)
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		logger.Fatalf("Could not connect to Redis: %v", err)
	}

	return client
}

func registerRepositories(
	logger *logging.Logger,
	postgresClient *pgxpool.Pool,
) (
	storage.ContactsRepository,
	storage2.ItemsRepository,
	storage3.OrdersRepository,
) {
	contacts := contactsRepository.NewRepository(logger, postgresClient)
	items := itemsRepository.NewRepository(logger, postgresClient)
	orders := ordersRepository.NewRepository(logger, postgresClient)

	return contacts, items, orders
}

func initializeServices(
	logger *logging.Logger,
	yandexClient *s3.PresignClient,
	redisClient *redis.Client,
	cfg *config.Config,
	contacts storage.ContactsRepository,
	items storage2.ItemsRepository,
	orders storage3.OrdersRepository,
) (
	service1.ContactService,
	service2.ItemService,
	service3.OrderService,
	service5.YandexService,
	service6.RedisIpRequestService,
	service6.RedisItemService,
) {
	redisSvc := service6.NewRedisIpRequestService(logger, redisClient)
	redisItemsSvc := service6.NewRedisItemsService(logger, redisClient)
	contactSvc := service1.NewContactService(logger, contacts, nil)
	itemSvc := service2.NewItemService(logger, items, redisItemsSvc, nil)
	orderSvc := service3.NewOrderService(logger, orders, nil, nil, nil)
	yandexSvc := service5.NewYandexService(logger, yandexClient, cfg)

	return contactSvc, itemSvc, orderSvc, yandexSvc, redisSvc, redisItemsSvc
}

func initializeMiddleware(logger *logging.Logger, redisSvc service6.RedisIpRequestService) *dns.Middleware {
	redisMiddleware := dns.NewIpRequestMiddleware(logger, redisSvc)

	return redisMiddleware
}

func initializeFacades(
	contactSvc service1.ContactService,
	itemSvc service2.ItemService,
	orderSvc service3.OrderService,
	yandexSvc service5.YandexService,
) (
	facades.ContactOrderFacade,
	facades.OrderItemFacade,
	facades.YandexItemFacade,
	facades.YandexOrderFacade,
) {
	orderContactFacade := service4.NewOrderContactFacade(orderSvc, contactSvc)
	orderItemFacade := service4.NewOrderItemFacade(orderSvc, itemSvc)
	yandexItemFacade := service4.NewYandexItemFacade(yandexSvc, itemSvc)
	yandexOrderFacade := service4.NewYandexOrderFacade(yandexSvc, orderSvc)

	orderSvc.SetContactFacade(orderContactFacade)
	orderSvc.SetItemFacade(orderItemFacade)
	orderSvc.SetYandexFacade(yandexOrderFacade)
	contactSvc.SetYandexFacade(yandexOrderFacade)
	itemSvc.SetYandexFacade(yandexItemFacade)

	return orderContactFacade, orderItemFacade, yandexItemFacade, yandexOrderFacade
}

func initializeHandlers(
	logger *logging.Logger,
	ipMiddleware *dns.Middleware,
	router *httprouter.Router,
	yandexClient *s3.PresignClient,
	cfg *config.Config,
	contacts storage.ContactsRepository,
	items storage2.ItemsRepository,
	orders storage3.OrdersRepository,
	redisItemsSvc service6.RedisItemService,
	contactFacade facades.ContactOrderFacade,
	itemFacade facades.OrderItemFacade,
	yandexFacade facades.YandexItemFacade,
	yandexOrderFacade facades.YandexOrderFacade,
) {
	contactHandler := contactsHandler.NewContactHandler(logger, ipMiddleware, contacts, yandexOrderFacade)
	itemHandler := itemsHandler.NewItemHandler(items, logger, redisItemsSvc, yandexFacade)
	orderHandler := ordersHandler.NewOrderHandler(logger, ipMiddleware, contactFacade, itemFacade, yandexOrderFacade, orders)
	cloudHandler := yandexHandler.NewYandexHandler(logger, yandexClient, cfg)
	contactHandler.SaveContact(router)
	itemHandler.SaveItem(router)
	orderHandler.SaveOrder(router)
	cloudHandler.SaveYandexStorage(router)
}

func start(logger *logging.Logger, router *httprouter.Router, cfg *config.Config) {
	logger.Info("Starting the application...")

	var listener net.Listener
	var listenErr error

	if cfg.Server.Type == "sock" {
		logger.Info("Detecting the path of an application...")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}

		logger.Info("Initializing a socket...")
		socketPath := path.Join(appDir, "app.sock")

		logger.Info("Listening unix socket...")
		listener, listenErr = net.Listen("unix", socketPath)

		logger.Infof("Server is listening unix socket: %s", socketPath)
	} else {
		logger.Info("Listening TCP...")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Server.BindIP, cfg.Server.Port))

		logger.Infof("Server is listening on host: %s:%s", cfg.Server.BindIP, cfg.Server.Port)
	}

	if listenErr != nil {
		panic(listenErr)
	}

	mail.NewGmailSender(cfg.Email.Sender, cfg.Email.Address, cfg.Email.Password)

	corsSettings := cors.CorsSettings()
	handler := corsSettings.Handler(router)

	server := &http.Server{
		Handler:      handler,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}

func main() {
	logger := logging.GetLogger()
	cfg := config.GetConfig()

	sdkConfig, err := config2.LoadDefaultConfig(context.TODO())
	if err != nil {
		logger.Fatalf("%v", err)
	}

	redisClient := NewRedisClient(cfg, logger)

	yandexClient := s3.NewFromConfig(sdkConfig)
	presignedClient := s3.NewPresignClient(yandexClient)

	postgresClient, err := postgresql.NewClient(context.TODO(), cfg.Postgresql)
	if err != nil {
		logger.Fatalf("%v", err)
	}

	logger.Info("Initializing router...")
	router := httprouter.New()

	logger.Info("Initializing repositories...")
	contacts, items, orders := registerRepositories(logger, postgresClient)

	logger.Info("Initializing services...")
	contactSvc, itemSvc, orderSvc, yandexSvc, redisSvc, redisItemsSvc := initializeServices(logger, presignedClient, redisClient, cfg, contacts, items, orders)
	orderContactFacade, orderItemFacade, yandexItemFacade, yandexOrderFacade := initializeFacades(contactSvc, itemSvc, orderSvc, yandexSvc)

	logger.Info("Initializing middleware...")
	ipMiddleware := initializeMiddleware(logger, redisSvc)

	logger.Info("Initializing handlers...")
	initializeHandlers(logger, ipMiddleware, router, presignedClient, cfg, contacts, items, orders, redisItemsSvc, orderContactFacade, orderItemFacade, yandexItemFacade, yandexOrderFacade)

	start(logger, router, cfg)

	select {}
}
