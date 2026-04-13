package server

import (
	"net/http"

	"github.com/sanaul03/ai-sdlc-backend/internal/handler"
	"github.com/sanaul03/ai-sdlc-backend/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

// New builds and returns an http.Handler with all routes registered.
func New(db *pgxpool.Pool) http.Handler {
	mux := http.NewServeMux()

	companyRepo := repository.NewCompanyRepository(db)
	depotRepo := repository.NewDepotRepository(db)
	vcRepo := repository.NewVehicleCategoryRepository(db)
	vtRepo := repository.NewVehicleTypeRepository(db)
	vehicleRepo := repository.NewVehicleRepository(db)

	companyH := handler.NewCompanyHandler(companyRepo)
	depotH := handler.NewDepotHandler(depotRepo)
	vcH := handler.NewVehicleCategoryHandler(vcRepo)
	vtH := handler.NewVehicleTypeHandler(vtRepo)
	vehicleH := handler.NewVehicleHandler(vehicleRepo)

	// Company routes
	mux.HandleFunc("POST /companies", companyH.Create)
	mux.HandleFunc("GET /companies", companyH.List)
	mux.HandleFunc("GET /companies/{id}", companyH.GetByID)
	mux.HandleFunc("PUT /companies/{id}", companyH.Update)
	mux.HandleFunc("DELETE /companies/{id}", companyH.Delete)

	// Depot routes
	mux.HandleFunc("POST /depots", depotH.Create)
	mux.HandleFunc("GET /depots/{id}", depotH.GetByID)
	mux.HandleFunc("GET /companies/{id}/depots", depotH.ListByCompany)
	mux.HandleFunc("PUT /depots/{id}", depotH.Update)
	mux.HandleFunc("DELETE /depots/{id}", depotH.Delete)

	// Vehicle category routes
	mux.HandleFunc("POST /vehicle-categories", vcH.Create)
	mux.HandleFunc("GET /vehicle-categories", vcH.List)
	mux.HandleFunc("GET /vehicle-categories/{id}", vcH.GetByID)
	mux.HandleFunc("PUT /vehicle-categories/{id}", vcH.Update)
	mux.HandleFunc("DELETE /vehicle-categories/{id}", vcH.Delete)

	// Vehicle type routes
	mux.HandleFunc("POST /vehicle-types", vtH.Create)
	mux.HandleFunc("GET /vehicle-types/{id}", vtH.GetByID)
	mux.HandleFunc("GET /vehicle-categories/{id}/vehicle-types", vtH.ListByCategory)
	mux.HandleFunc("PUT /vehicle-types/{id}", vtH.Update)
	mux.HandleFunc("DELETE /vehicle-types/{id}", vtH.Delete)

	// Vehicle routes
	mux.HandleFunc("POST /vehicles", vehicleH.Create)
	mux.HandleFunc("GET /vehicles/{id}", vehicleH.GetByID)
	mux.HandleFunc("GET /companies/{id}/vehicles", vehicleH.ListByCompany)
	mux.HandleFunc("PUT /vehicles/{id}", vehicleH.Update)
	mux.HandleFunc("DELETE /vehicles/{id}", vehicleH.Delete)

	return mux
}
