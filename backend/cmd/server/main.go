package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/koperasi-gresik/backend/config"
	"github.com/koperasi-gresik/backend/internal/shared/database"
	"github.com/koperasi-gresik/backend/internal/shared/database/seeds"
	"github.com/koperasi-gresik/backend/internal/shared/event"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"github.com/koperasi-gresik/backend/internal/shared/notification"
	"github.com/koperasi-gresik/backend/internal/shared/worker"

	iamHandler "github.com/koperasi-gresik/backend/internal/modules/iam/handler"
	iamRepo "github.com/koperasi-gresik/backend/internal/modules/iam/repository"
	iamService "github.com/koperasi-gresik/backend/internal/modules/iam/service"

	orgHandler "github.com/koperasi-gresik/backend/internal/modules/organization/handler"
	orgRepo "github.com/koperasi-gresik/backend/internal/modules/organization/repository"
	orgService "github.com/koperasi-gresik/backend/internal/modules/organization/service"

	memberHandler "github.com/koperasi-gresik/backend/internal/modules/member/handler"
	memberRepo "github.com/koperasi-gresik/backend/internal/modules/member/repository"
	memberService "github.com/koperasi-gresik/backend/internal/modules/member/service"

	savingHandler "github.com/koperasi-gresik/backend/internal/modules/savings/handler"
	savingRepo "github.com/koperasi-gresik/backend/internal/modules/savings/repository"
	savingService "github.com/koperasi-gresik/backend/internal/modules/savings/service"

	accountingHandler "github.com/koperasi-gresik/backend/internal/modules/accounting/handler"
	accountingRepo "github.com/koperasi-gresik/backend/internal/modules/accounting/repository"
	accountingService "github.com/koperasi-gresik/backend/internal/modules/accounting/service"

	cashHandler "github.com/koperasi-gresik/backend/internal/modules/cash/handler"
	cashRepo "github.com/koperasi-gresik/backend/internal/modules/cash/repository"
	cashService "github.com/koperasi-gresik/backend/internal/modules/cash/service"

	loanHandler "github.com/koperasi-gresik/backend/internal/modules/loan/handler"
	loanRepo "github.com/koperasi-gresik/backend/internal/modules/loan/repository"
	loanService "github.com/koperasi-gresik/backend/internal/modules/loan/service"

	inventoryHandler "github.com/koperasi-gresik/backend/internal/modules/inventory/handler"
	inventoryRepo "github.com/koperasi-gresik/backend/internal/modules/inventory/repository"
	inventoryService "github.com/koperasi-gresik/backend/internal/modules/inventory/service"

	supplierHandler "github.com/koperasi-gresik/backend/internal/modules/supplier/handler"
	supplierRepo "github.com/koperasi-gresik/backend/internal/modules/supplier/repository"
	supplierService "github.com/koperasi-gresik/backend/internal/modules/supplier/service"

	salesHandler "github.com/koperasi-gresik/backend/internal/modules/sales/handler"
	salesRepo "github.com/koperasi-gresik/backend/internal/modules/sales/repository"
	salesService "github.com/koperasi-gresik/backend/internal/modules/sales/service"

	purchasingHandler "github.com/koperasi-gresik/backend/internal/modules/purchasing/handler"
	purchasingRepo "github.com/koperasi-gresik/backend/internal/modules/purchasing/repository"
	purchasingService "github.com/koperasi-gresik/backend/internal/modules/purchasing/service"

	reportHandler "github.com/koperasi-gresik/backend/internal/modules/report/handler"
	reportService "github.com/koperasi-gresik/backend/internal/modules/report/service"

	auditHandler "github.com/koperasi-gresik/backend/internal/modules/audit/handler"
	auditRepo "github.com/koperasi-gresik/backend/internal/modules/audit/repository"

	billingHandler "github.com/koperasi-gresik/backend/internal/modules/billing/handler"
	billingRepo "github.com/koperasi-gresik/backend/internal/modules/billing/repository"

	notificationHandler "github.com/koperasi-gresik/backend/internal/modules/notification/handler"
	notificationRepo "github.com/koperasi-gresik/backend/internal/modules/notification/repository"
	notificationSvc "github.com/koperasi-gresik/backend/internal/modules/notification/service"

	posHandler "github.com/koperasi-gresik/backend/internal/modules/pos/handler"
	posRepo "github.com/koperasi-gresik/backend/internal/modules/pos/repository"
	posService "github.com/koperasi-gresik/backend/internal/modules/pos/service"

	shuHandler "github.com/koperasi-gresik/backend/internal/modules/shu/handler"
	shuRepo "github.com/koperasi-gresik/backend/internal/modules/shu/repository"
	shuService "github.com/koperasi-gresik/backend/internal/modules/shu/service"

	closingHandler "github.com/koperasi-gresik/backend/internal/modules/closing/handler"
	closingRepo "github.com/koperasi-gresik/backend/internal/modules/closing/repository"
	closingService "github.com/koperasi-gresik/backend/internal/modules/closing/service"
)

func main() {
	// Load .env file (ignore error in production)
	_ = godotenv.Load()

	// Load config
	cfg := config.Load()

	// Connect to database
	db := database.NewPostgres(cfg.Database)

	// Run Migrations
	database.RunMigrations(cfg.Database)

	// Connect to Redis
	rdb := database.NewRedis(cfg.Redis)

	// Permission Cache (Redis-backed — zero DB queries in middleware)
	permCache := middleware.NewPermissionCache(rdb, db)

	// Seed global system permissions at startup (idempotent upsert)
	if err := seeds.SeedPermissions(context.Background(), db); err != nil {
		log.Printf("⚠️  Failed to seed permissions: %v", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ErrorHandler: customErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": cfg.App.Name,
		})
	})

	// API v1 Group
	v1 := app.Group("/api/v1")

	// Middleware instances
	authMid := middleware.Auth(cfg.JWT.Secret)
	tenantMid := middleware.Tenant()

	// ─── Module Route Registration ───

	// Repositories
	userRepository := iamRepo.NewUserRepository(db)
	roleRepository := iamRepo.NewRoleRepository(db)
	tokenRepository := iamRepo.NewTokenRepository(db)
	orgRepository := orgRepo.NewOrganizationRepository(db)
	memberRepository := memberRepo.NewMemberRepository(db)
	savingRepository := savingRepo.NewSavingRepository(db)
	accountingRepository := accountingRepo.NewAccountingRepository(db)
	cashRepository := cashRepo.NewCashRepository(db)
	loanRepository := loanRepo.NewLoanRepository(db)
	inventoryRepository := inventoryRepo.NewInventoryRepository(db)
	warehouseRepository := inventoryRepo.NewWarehouseRepository(db)
	supplierRepository := supplierRepo.NewSupplierRepository(db)
	salesRepository := salesRepo.NewSalesRepository(db)
	purchasingRepository := purchasingRepo.NewPurchasingRepository(db)
	auditRepository := auditRepo.NewAuditRepository(db)
	billingRepository := billingRepo.NewBillingRepository(db)
	notificationRepoObj := notificationRepo.NewNotificationRepository(db)
	posRepository := posRepo.NewPOSRepository(db)
	shuRepository := shuRepo.NewSHURepository(db)
	closingRepository := closingRepo.NewClosingRepository(db)

	// Event Bus (Redis Streams for Persistence)
	eventBus := event.NewRedisEventBus(rdb, "publisher")
	notifEventBus := event.NewRedisEventBus(rdb, "notification_group")
	accountingEventBus := event.NewRedisEventBus(rdb, "accounting_group")
	defer eventBus.Close()

	// Services
	authenticationService := iamService.NewAuthService(userRepository, roleRepository, tokenRepository, eventBus, cfg.JWT.Secret, cfg.JWT.ExpirationHours)
	organizationService := orgService.NewOrganizationService(orgRepository, userRepository, roleRepository, db, rdb)
	membershipService := memberService.NewMemberService(memberRepository, authenticationService)
	savingModuleService := savingService.NewSavingService(savingRepository, eventBus)
	reportModuleService := reportService.NewReportService(db, rdb)
	accountingModuleService := accountingService.NewAccountingService(accountingRepository, reportModuleService)
	cashModuleService := cashService.NewCashService(cashRepository)
	loanModuleService := loanService.NewLoanService(loanRepository, eventBus)
	inventoryModuleService := inventoryService.NewInventoryService(inventoryRepository, warehouseRepository)
	supplierModuleService := supplierService.NewSupplierService(supplierRepository)
	salesModuleService := salesService.NewSalesService(salesRepository, inventoryModuleService)
	purchasingModuleService := purchasingService.NewPurchasingService(purchasingRepository, inventoryModuleService)
	shuModuleService := shuService.NewSHUService(shuRepository)
	warehouseModuleService := inventoryService.NewWarehouseService(warehouseRepository, inventoryRepository)
	posModuleService := posService.NewPOSService(posRepository, inventoryModuleService)
	closingModuleService := closingService.NewClosingService(closingRepository, loanRepository, savingRepository, accountingModuleService, orgRepository)
	mobileModuleService := memberService.NewMobileService(memberRepository, savingRepository, loanRepository)

	// Async Background Workers
	taskDistributor := worker.NewRedisTaskDistributor(cfg.Redis)
	workerServer := worker.NewServer(cfg.Redis, cfg.SMTP)
	workerServer.RegisterClosingHandlers(closingModuleService) // Register EOD/EOM task handlers

	// Automated Cron Scheduler (Asynq)
	asynqScheduler := worker.NewScheduler(cfg.Redis)
	// EOD runs every day at 01:00 AM
	_ = asynqScheduler.RegisterTask("0 1 * * *", worker.TypeRunEOD, worker.PayloadRunEOD{Date: time.Now().AddDate(0, 0, -1).Format("2006-01-02")})
	// EOM runs on the 1st of every month at 01:10 AM
	_ = asynqScheduler.RegisterTask("10 1 1 * *", worker.TypeRunEOM, worker.PayloadRunEOM{Month: int(time.Now().Month()), Year: time.Now().Year()})

	go func() {
		if err := workerServer.Start(); err != nil {
			logrus.Fatalf("[Asynq] Worker server failed: %v", err)
		}
	}()
	go func() {
		if err := asynqScheduler.Start(); err != nil {
			logrus.Fatalf("[Asynq] Scheduler failed: %v", err)
		}
	}()
	defer workerServer.Stop()
	defer asynqScheduler.Stop()

	// Start Notification Engine (Event Subscriber)
	pushAdapter := notification.NewFCMAdapter()
	notificationModuleService := notificationSvc.NewNotificationService(notificationRepoObj, tokenRepository, pushAdapter, notifEventBus, taskDistributor)
	notificationModuleService.Start(context.Background())

	// Start Accounting Event Handler
	accountingEventHandler := accountingService.NewAccountingEventHandler(accountingModuleService, accountingEventBus)
	go accountingEventHandler.Start(context.Background())

	// Handlers
	authenticationHandler := iamHandler.NewAuthHandler(authenticationService)
	organizationHandler := orgHandler.NewOrganizationHandler(organizationService)
	membershipHandler := memberHandler.NewMemberHandler(membershipService)
	savingsHandler := savingHandler.NewSavingHandler(savingModuleService)
	accountingModuleHandler := accountingHandler.NewAccountingHandler(accountingModuleService)
	cashModuleHandler := cashHandler.NewCashHandler(cashModuleService)
	loanModuleHandler := loanHandler.NewLoanHandler(loanModuleService)
	inventoryModuleHandler := inventoryHandler.NewInventoryHandler(inventoryModuleService)
	warehouseModuleHandler := inventoryHandler.NewWarehouseHandler(warehouseModuleService)
	supplierModuleHandler := supplierHandler.NewSupplierHandler(supplierModuleService)
	salesModuleHandler := salesHandler.NewSalesHandler(salesModuleService)
	purchasingModuleHandler := purchasingHandler.NewPurchasingHandler(purchasingModuleService)
	reportModuleHandler := reportHandler.NewReportHandler(reportModuleService)
	auditModuleHandler := auditHandler.NewAuditHandler(auditRepository)
	billingModuleHandler := billingHandler.NewBillingHandler(billingRepository)
	notificationModuleHandler := notificationHandler.NewNotificationHandler(notificationRepoObj)
	posModuleHandler := posHandler.NewPOSHandler(posModuleService)
	shuModuleHandler := shuHandler.NewSHUHandler(shuModuleService)
	closingModuleHandler := closingHandler.NewClosingHandler(closingModuleService)
	mobileModuleHandler := memberHandler.NewMobileHandler(mobileModuleService)

	// Register Routes
	iamHandler.RegisterPublicRoutes(v1, authenticationHandler)
	orgHandler.RegisterPublicRoutes(v1, organizationHandler)

	// RBAC Management Routes
	rbacRoleRepo := iamRepo.NewRoleRepository(db)
	rbacRoleService := iamService.NewRoleService(rbacRoleRepo, db, permCache)
	rbacRoleHandler := iamHandler.NewRoleHandler(rbacRoleService)
	iamHandler.RegisterRoleRoutes(v1, rbacRoleHandler, authMid, tenantMid, permCache)

	// Protected Routes (Apply middleware explicitly)
	iamHandler.RegisterRoutes(v1, authenticationHandler, authMid, tenantMid)
	orgHandler.RegisterRoutes(v1, organizationHandler, authMid, tenantMid)
	memberHandler.RegisterRoutes(v1, membershipHandler, authMid, tenantMid)
	memberHandler.RegisterMobileRoutes(v1, mobileModuleHandler, authMid, tenantMid)
	savingHandler.RegisterRoutes(v1, savingsHandler, authMid, tenantMid, middleware.ModuleGuard(db, "savings"))
	accountingHandler.RegisterRoutes(v1, accountingModuleHandler, authMid, tenantMid)
	cashHandler.RegisterRoutes(v1, cashModuleHandler, authMid, tenantMid)
	loanHandler.RegisterRoutes(v1, loanModuleHandler, authMid, tenantMid, middleware.ModuleGuard(db, "loans"))
	inventoryHandler.RegisterRoutes(v1, inventoryModuleHandler, authMid, tenantMid)
	inventoryHandler.RegisterWarehouseRoutes(v1, warehouseModuleHandler, authMid, tenantMid)
	supplierHandler.RegisterRoutes(v1, supplierModuleHandler, authMid, tenantMid)
	salesHandler.RegisterRoutes(v1, salesModuleHandler, authMid, tenantMid)
	purchasingHandler.RegisterRoutes(v1, purchasingModuleHandler, authMid, tenantMid)
	reportHandler.RegisterRoutes(v1, reportModuleHandler, authMid, tenantMid)
	auditHandler.RegisterRoutes(v1, auditModuleHandler, authMid, tenantMid)
	billingHandler.RegisterRoutes(v1, billingModuleHandler, authMid, tenantMid)
	notificationHandler.RegisterRoutes(v1, notificationModuleHandler, authMid, tenantMid)
	posHandler.RegisterRoutes(v1, posModuleHandler, authMid, tenantMid)
	shuHandler.RegisterRoutes(v1, shuModuleHandler, authMid, tenantMid)
	closingHandler.RegisterRoutes(v1, closingModuleHandler, authMid) // Admin manual triggers

	// Period Guard Middleware setup (Optional inclusion in other routes)
	// Example: v1.Use(middleware.PeriodGuard(closingRepository)) 
	// (applied after authMid/tenantMid to ensure we have context)

	// Start server
	go func() {
		addr := fmt.Sprintf(":%s", cfg.App.Port)
		log.Printf("🚀 %s starting on %s", cfg.App.Name, addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Println("✅ Server stopped gracefully")
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": err.Error(),
	})
}
