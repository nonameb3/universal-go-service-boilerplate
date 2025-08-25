package main

import (
	"fmt"
	"log"

	"github.com/universal-go-service/boilerplate/cmd/migrations"
	"github.com/universal-go-service/boilerplate/config"
	"github.com/universal-go-service/boilerplate/pkg/providers/database"
)

func main() {
	fmt.Printf("🚀 Universal Go Service Boilerplate\n")
	fmt.Printf("===================================\n\n")

	// Load configuration
	env := config.GetEnvironment()
	cfg := config.GetConfig(env)

	fmt.Printf("📋 Configuration loaded successfully!\n")
	fmt.Printf("🌍 Environment: %s\n", env)
	fmt.Printf("🏠 Server: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("📦 App: %s v%s\n", cfg.App.Name, cfg.App.Version)
	fmt.Printf("🔧 Database: %s:%v@%s:%s\n", cfg.Db.Host, cfg.Db.Port, cfg.Db.User, cfg.Db.Password)
	fmt.Printf("🔧 Debug mode: %t\n\n", cfg.App.Debug)

	fmt.Printf("\n✅ Universal Go Service Boilerplate is ready!\n\n")

	fmt.Printf("🎯 What you get out of the box:\n")
	fmt.Printf("   ✅ Multi-environment configuration (local, dev, staging, prod)\n")
	fmt.Printf("   ✅ Pluggable provider system (logger, metrics, auth, cache, database)\n")
	fmt.Printf("   ✅ Clean architecture structure\n")
	fmt.Printf("   ✅ Production-ready defaults\n")
	fmt.Printf("   ✅ Easy company library integration\n\n")

	// migrate postgres database
	db, err := database.NewPostgres(database.DatabaseConfig{
		Host:     cfg.Db.Host,
		Port:     cfg.Db.Port,
		Username: cfg.Db.User,
		Password: cfg.Db.Password,
		Database: cfg.Db.DBName,
		SSLMode:  cfg.Db.SSLMode,
		Timezone: cfg.Db.TimeZone,
	})
	if err != nil {
		log.Fatalf("Failed to get database: %v", err)
	}
	if cfg.Db.AutoMigrate {
		// Test the database connection first
		err = db.Health()
		if err != nil {
			log.Printf("⚠️ Database health check failed: %v", err)
			log.Printf("⚠️ Migration skipped - no database connection")
		} else {
			migrations.ExecuteMigration(db)
		}
	}

	fmt.Printf("🚀 Next steps to build your service:\n")
	fmt.Printf("   1. Add your domain entities to internal/domain/\n")
	fmt.Printf("   2. Create use cases in internal/usecase/\n")
	fmt.Printf("   3. Implement repositories in internal/repository/\n")
	fmt.Printf("   4. Add HTTP handlers in internal/handler/\n")
	fmt.Printf("   5. Replace default providers with your company's libraries\n\n")

	fmt.Printf("🔧 To integrate your company's libraries:\n")
	fmt.Printf("   • Create provider implementations in pkg/providers/\n")
	fmt.Printf("   • Register them in the factory system\n")
	fmt.Printf("   • Update configuration files\n")
	fmt.Printf("   • No code changes needed elsewhere!\n\n")

	fmt.Printf("💡 Example: TokenX Logger Integration\n")
	fmt.Printf("   See examples/company-logger/ for a complete example\n")
}
