package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"github.com/koperasi-gresik/backend/config"
	"github.com/koperasi-gresik/backend/internal/shared/database"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"

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

	posHandler "github.com/koperasi-gresik/backend/internal/modules/pos/handler"
	posRepo "github.com/koperasi-gresik/backend/internal/modules/pos/repository"

	shuHandler "github.com/koperasi-gresik/backend/internal/modules/shu/handler"
	shuRepo "github.com/koperasi-gresik/backend/internal/modules/shu/repository"
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
	_ = rdb // Will be used by modules

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
	orgRepository := orgRepo.NewOrganizationRepository(db)
	memberRepository := memberRepo.NewMemberRepository(db)
	savingRepository := savingRepo.NewSavingRepository(db)
	accountingRepository := accountingRepo.NewAccountingRepository(db)
	cashRepository := cashRepo.NewCashRepository(db)
	loanRepository := loanRepo.NewLoanRepository(db)
	inventoryRepository := inventoryRepo.NewInventoryRepository(db)
	supplierRepository := supplierRepo.NewSupplierRepository(db)
	salesRepository := salesRepo.NewSalesRepository(db)
	purchasingRepository := purchasingRepo.NewPurchasingRepository(db)
	auditRepository := auditRepo.NewAuditRepository(db)
	billingRepository := billingRepo.NewBillingRepository(db)
	notificationRepository := notificationRepo.NewNotificationRepository(db)
	posRepository := posRepo.NewPOSRepository(db)
	shuRepository := shuRepo.NewSHURepository(db)

	// Services
	authenticationService := iamService.NewAuthService(userRepository, roleRepository, cfg.JWT.Secret, cfg.JWT.ExpirationHours)
	organizationService := orgService.NewOrganizationService(orgRepository, userRepository, roleRepository, db)
	membershipService := memberService.NewMemberService(memberRepository, authenticationService)
	savingModuleService := savingService.NewSavingService(savingRepository)
	accountingModuleService := accountingService.NewAccountingService(accountingRepository)
	cashModuleService := cashService.NewCashService(cashRepository)
	loanModuleService := loanService.NewLoanService(loanRepository)
	inventoryModuleService := inventoryService.NewInventoryService(inventoryRepository)
	supplierModuleService := supplierService.NewSupplierService(supplierRepository)
	salesModuleService := salesService.NewSalesService(salesRepository, inventoryModuleService)
	purchasingModuleService := purchasingService.NewPurchasingService(purchasingRepository, inventoryModuleService)
	reportModuleService := reportService.NewReportService(db)

	// Handlers
	authenticationHandler := iamHandler.NewAuthHandler(authenticationService)
	organizationHandler := orgHandler.NewOrganizationHandler(organizationService)
	membershipHandler := memberHandler.NewMemberHandler(membershipService)
	savingsHandler := savingHandler.NewSavingHandler(savingModuleService)
	accountingModuleHandler := accountingHandler.NewAccountingHandler(accountingModuleService)
	cashModuleHandler := cashHandler.NewCashHandler(cashModuleService)
	loanModuleHandler := loanHandler.NewLoanHandler(loanModuleService)
	inventoryModuleHandler := inventoryHandler.NewInventoryHandler(inventoryModuleService)
	supplierModuleHandler := supplierHandler.NewSupplierHandler(supplierModuleService)
	salesModuleHandler := salesHandler.NewSalesHandler(salesModuleService)
	purchasingModuleHandler := purchasingHandler.NewPurchasingHandler(purchasingModuleService)
	reportModuleHandler := reportHandler.NewReportHandler(reportModuleService)
	auditModuleHandler := auditHandler.NewAuditHandler(auditRepository)
	billingModuleHandler := billingHandler.NewBillingHandler(billingRepository)
	notificationModuleHandler := notificationHandler.NewNotificationHandler(notificationRepository)
	posModuleHandler := posHandler.NewPOSHandler(posRepository)
	shuModuleHandler := shuHandler.NewSHUHandler(shuRepository)

	// Register Routes
	iamHandler.RegisterPublicRoutes(v1, authenticationHandler)
	orgHandler.RegisterPublicRoutes(v1, organizationHandler)

	// Protected Routes (Apply middleware explicitly)
	iamHandler.RegisterRoutes(v1, authenticationHandler, authMid, tenantMid)
	orgHandler.RegisterRoutes(v1, organizationHandler, authMid, tenantMid)
	memberHandler.RegisterRoutes(v1, membershipHandler, authMid, tenantMid)
	savingHandler.RegisterRoutes(v1, savingsHandler, authMid, tenantMid, middleware.ModuleGuard(db, "savings"))
	accountingHandler.RegisterRoutes(v1, accountingModuleHandler, authMid, tenantMid)
	cashHandler.RegisterRoutes(v1, cashModuleHandler, authMid, tenantMid)
	loanHandler.RegisterRoutes(v1, loanModuleHandler, authMid, tenantMid, middleware.ModuleGuard(db, "loans"))
	inventoryHandler.RegisterRoutes(v1, inventoryModuleHandler, authMid, tenantMid)
	supplierHandler.RegisterRoutes(v1, supplierModuleHandler, authMid, tenantMid)
	salesHandler.RegisterRoutes(v1, salesModuleHandler, authMid, tenantMid)
	purchasingHandler.RegisterRoutes(v1, purchasingModuleHandler, authMid, tenantMid)
	reportHandler.RegisterRoutes(v1, reportModuleHandler, authMid, tenantMid)
	auditHandler.RegisterRoutes(v1, auditModuleHandler, authMid, tenantMid)
	billingHandler.RegisterRoutes(v1, billingModuleHandler, authMid, tenantMid)
	notificationHandler.RegisterRoutes(v1, notificationModuleHandler, authMid, tenantMid)
	posHandler.RegisterRoutes(v1, posModuleHandler, authMid, tenantMid)
	shuHandler.RegisterRoutes(v1, shuModuleHandler, authMid, tenantMid)

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
