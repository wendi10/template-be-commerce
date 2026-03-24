package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	appmiddleware "github.com/template-be-commerce/internal/middleware"
	jwtpkg "github.com/template-be-commerce/pkg/jwt"
)

type Handlers struct {
	Auth     *AuthHandler
	Customer *CustomerHandler
	Product  *ProductHandler
	Banner   *BannerHandler
	Promo    *PromoHandler
	Cart     *CartHandler
	Order    *OrderHandler
	Payment  *PaymentHandler
	Admin    *AdminHandler
}

// NewRouter builds and returns the chi router with all routes registered.
func NewRouter(h Handlers, jwtManager *jwtpkg.Manager) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(appmiddleware.Logger)
	r.Use(appmiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/v1", func(r chi.Router) {
		// ─── Public AUTH ────────────────────────────────────────────────────────
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", h.Auth.Register)
			r.Post("/login", h.Auth.LoginCustomer)
			r.Post("/refresh", h.Auth.RefreshToken)
		})

		// ─── Public PRODUCTS & CATEGORIES ───────────────────────────────────────
		r.Get("/products", h.Product.ListProducts)
		r.Get("/products/{id}", h.Product.GetProduct)
		r.Get("/products/slug/{slug}", h.Product.GetProductBySlug)
		r.Get("/categories", h.Product.ListCategories)

		// ─── Public BANNERS ──────────────────────────────────────────────────────
		r.Get("/banners", h.Banner.ListActiveBanners)

		// ─── Public PROMO validation ─────────────────────────────────────────────
		r.Post("/promos/validate", h.Promo.ValidatePromo)

		// ─── Payment CALLBACK (no auth - signed by provider) ─────────────────────
		r.Post("/payments/callback/{provider}", h.Payment.HandleCallback)

		// ─── Customer (authenticated) ────────────────────────────────────────────
		r.Group(func(r chi.Router) {
			r.Use(appmiddleware.Authenticate(jwtManager))
			r.Use(appmiddleware.RequireCustomer)

			// Profile
			r.Get("/me", h.Customer.GetProfile)
			r.Put("/me", h.Customer.UpdateProfile)

			// Addresses
			r.Get("/me/addresses", h.Customer.ListAddresses)
			r.Post("/me/addresses", h.Customer.CreateAddress)
			r.Put("/me/addresses/{addressID}", h.Customer.UpdateAddress)
			r.Delete("/me/addresses/{addressID}", h.Customer.DeleteAddress)
			r.Patch("/me/addresses/{addressID}/default", h.Customer.SetDefaultAddress)

			// Cart
			r.Get("/cart", h.Cart.GetCart)
			r.Post("/cart", h.Cart.AddToCart)
			r.Put("/cart/{itemID}", h.Cart.UpdateCartItem)
			r.Delete("/cart/{itemID}", h.Cart.RemoveCartItem)

			// Orders
			r.Post("/orders", h.Order.CreateOrder)
			r.Get("/orders", h.Order.ListMyOrders)
			r.Get("/orders/{id}", h.Order.GetOrder)
			r.Post("/orders/{id}/cancel", h.Order.CancelOrder)

			// Payments
			r.Post("/payments", h.Payment.CreatePayment)
			r.Get("/payments/order/{orderID}", h.Payment.GetPaymentByOrder)
		})

		// ─── ADMIN routes ────────────────────────────────────────────────────────
		r.Route("/admin", func(r chi.Router) {
			// Admin auth (no auth middleware needed)
			r.Post("/auth/register", h.Auth.RegisterAdmin)
			r.Post("/auth/login", h.Auth.LoginAdmin)

			// Protected admin routes
			r.Group(func(r chi.Router) {
				r.Use(appmiddleware.Authenticate(jwtManager))
				r.Use(appmiddleware.RequireAdmin)

				// Products
				r.Get("/products", h.Product.AdminListProducts)
				r.Post("/products", h.Product.AdminCreateProduct)
				r.Put("/products/{id}", h.Product.AdminUpdateProduct)
				r.Delete("/products/{id}", h.Product.AdminDeleteProduct)
				r.Post("/products/{id}/images", h.Product.AdminAddProductImage)
				r.Delete("/products/{id}/images/{imageID}", h.Product.AdminDeleteProductImage)

				// Categories
				r.Get("/categories", h.Product.AdminListCategories)
				r.Post("/categories", h.Product.AdminCreateCategory)
				r.Put("/categories/{id}", h.Product.AdminUpdateCategory)
				r.Delete("/categories/{id}", h.Product.AdminDeleteCategory)

				// Banners
				r.Get("/banners", h.Banner.AdminListBanners)
				r.Post("/banners", h.Banner.AdminCreateBanner)
				r.Get("/banners/{id}", h.Banner.AdminGetBanner)
				r.Put("/banners/{id}", h.Banner.AdminUpdateBanner)
				r.Delete("/banners/{id}", h.Banner.AdminDeleteBanner)

				// Promos
				r.Get("/promos", h.Promo.AdminListPromos)
				r.Post("/promos", h.Promo.AdminCreatePromo)
				r.Get("/promos/{id}", h.Promo.AdminGetPromo)
				r.Put("/promos/{id}", h.Promo.AdminUpdatePromo)
				r.Delete("/promos/{id}", h.Promo.AdminDeletePromo)

				// Orders
				r.Get("/orders", h.Order.AdminListOrders)
				r.Get("/orders/{id}", h.Order.AdminGetOrder)
				r.Patch("/orders/{id}/status", h.Order.AdminUpdateOrderStatus)

				// Customers
				r.Get("/customers", h.Admin.ListCustomers)

				// Reports
				r.Get("/reports/sales", h.Order.AdminSalesReport)
			})
		})
	})

	return r
}
