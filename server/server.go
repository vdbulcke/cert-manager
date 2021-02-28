package server

import (
	"net/http"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/vdbulcke/cert-manager/data"
	"github.com/vdbulcke/cert-manager/handlers/api"
	"github.com/vdbulcke/cert-manager/handlers/certificate"
	"github.com/vdbulcke/cert-manager/handlers/tag"
)

func makeServer(l hclog.Logger, v *data.Validation, certBackend *data.CertBackend) http.Server {
	// create handlers
	apiHandler := api.NewAPI()
	certHandler := certificate.NewAPICertificateHandler(l, v, certBackend)
	tagHandler := tag.NewAPITagHandler(l, v, certBackend)

	// API Base Path
	apiBasePath := "/api/beta2"

	// create a new serve mux and register the handlers
	mainRouter := mux.NewRouter()

	// handlers for API
	apiRouter := mainRouter.PathPrefix(apiBasePath).Subrouter()
	apiRouter.Use(apiHandler.CommonAPIMiddleware)

	//
	// Create Routes
	//

	// Certificates API
	// GET
	apiRouter.Handle(
		"/certificate/GetCertificateByID/{id}",
		api.Handler{Handler: certHandler.GetCertificateByID}).
		Methods(http.MethodGet)

	apiRouter.Handle(
		"/certificate/GetCertificateByFingerprint/{id}",
		api.Handler{Handler: certHandler.GetCertificateByFingerprint}).
		Methods(http.MethodGet)

	apiRouter.Handle(
		"/certificate/ListCerts",
		api.Handler{Handler: certHandler.ListCerts}).
		Methods(http.MethodGet)

	// POST
	certAPIPost := apiRouter.Methods(http.MethodPost, http.MethodOptions).Subrouter()
	certAPIPost.Handle(
		"/certificate/CreateCertificate",
		api.Handler{Handler: certHandler.CreateCertificate})
	certAPIPost.Use(certHandler.MiddlewareValidateCertificateInput)

	// PUT
	certAPIPut := apiRouter.Methods(http.MethodPut).Subrouter()

	certAPIPut.Handle(
		"/certificate/UpdateCertificateTag/{id}",
		api.Handler{Handler: certHandler.UpdateCertificateTag})
	certAPIPut.Use(certHandler.MiddlewareValidateCertificateTagInput)

	// DELETE
	certAPIDelete := apiRouter.Methods(http.MethodDelete).Subrouter()
	certAPIDelete.Handle(
		"/certificate/DeleteCertificateTagsByID/{id}",
		api.Handler{Handler: certHandler.DeleteCertificateTagsByID})
	certAPIDelete.Use(certHandler.MiddlewareValidateCertificateTagInput)

	apiRouter.Handle(
		"/certificate/DeleteCertificateByID/{id}",
		api.Handler{Handler: certHandler.DeleteCertificateByID}).
		Methods(http.MethodDelete)

	// Tag API
	// GET

	apiRouter.Handle(
		"/tag/GetTagByID/{id}",
		api.Handler{Handler: tagHandler.GetTagByID}).
		Methods(http.MethodGet)

	apiRouter.Handle(
		"/tag/GetTagByName/{name}",
		api.Handler{Handler: tagHandler.GetTagByName}).
		Methods(http.MethodGet)

	apiRouter.Handle(
		"/tag/ListTags",
		api.Handler{Handler: tagHandler.ListTags}).
		Methods(http.MethodGet)

	// POST
	apiPost := apiRouter.Methods(http.MethodPost).Subrouter()
	apiPost.Handle(
		"/tag/CreateTag",
		api.Handler{Handler: tagHandler.CreateTag})
	apiPost.Use(tagHandler.MiddlewareValidateTagInput)

	// PUT
	apiPut := apiRouter.Methods(http.MethodPost).Subrouter()
	apiPut.Handle(
		"/tag/UpdateTagDescription/{id}",
		api.Handler{Handler: tagHandler.UpdateTagDescription})
	apiPut.Use(tagHandler.MiddlewareValidateTagDescriptionInput)

	// DELETE
	apiRouter.Handle(
		"/tag/DeleteTagByID/{id}",
		api.Handler{Handler: tagHandler.DeleteTagByID}).
		Methods(http.MethodDelete)

	//
	// Swagger
	//

	// handler for documentation
	// redocOpts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	// redocHandler := middleware.Redoc(redocOpts, nil)
	// handler for swagger ui
	swaggerOpts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	swaggerHandler := middleware.SwaggerUI(swaggerOpts, nil)

	// // Doc Handers
	// getDoc := mainRouter.Methods(http.MethodGet).Subrouter()
	// getDoc.Handle("/docs", redocHandler)

	// Swagger UI Handers
	getSwagger := mainRouter.Methods(http.MethodGet).Subrouter()
	getSwagger.Handle("/docs", swaggerHandler)

	//
	// TODO: handler swagger.yaml
	getSwaggerYaml := mainRouter.Methods(http.MethodGet).Subrouter()
	getSwaggerYaml.Handle("/swagger.yaml", http.FileServer(http.Dir("./static/")))

	//
	// CORS
	//
	// TODO: get cors from config
	corsAllowedHeader := handlers.AllowedHeaders([]string{"X-Requested-With"})
	corsAllowedOrigin := handlers.AllowedOrigins([]string{"*"})
	corsAllowedMethod := handlers.AllowedMethods([]string{
		http.MethodPost,
		http.MethodGet,
		http.MethodPut,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodHead,
	})
	corsHandler := handlers.CORS(corsAllowedMethod, corsAllowedOrigin, corsAllowedHeader)

	//
	// Http Server
	//
	// create a new server
	s := http.Server{
		Addr:         "0.0.0.0:9393",                                   // configure the bind address
		Handler:      corsHandler(mainRouter),                          // set the default handler
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{}), // set the logger for the server
		ReadTimeout:  5 * time.Second,                                  // max time to read request from the client
		WriteTimeout: 10 * time.Second,                                 // max time to write response to the client
		IdleTimeout:  120 * time.Second,                                // max time for connections using TCP Keep-Alive
	}

	return s

}
