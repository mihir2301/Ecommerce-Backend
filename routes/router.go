package routes

import (
	"golang_project/constants"
	"golang_project/controllers"
	"net/http"
)

var healthCheckRoutes = Router{
	Route{"HealthCheck", http.MethodGet, constants.HealthCheckRoute, controllers.HealthCheck},
}

var userRoutes = Router{
	//email verification routes
	Route{"VerifyEmail", http.MethodPost, constants.VerifyEmailRoute, controllers.VerifyEmail},
	Route{"VerifyOtp", http.MethodPost, constants.VerifyOtpRoute, controllers.VerifyOtp},
	Route{"ResendEmail", http.MethodPost, constants.ResendEmailRoute, controllers.VerifyEmail},

	Route{"Register User", http.MethodPost, constants.UserRegisterRoute, controllers.RegisterUser},
	Route{"Login User", http.MethodPost, constants.UserLoginRoute, controllers.UserLogin},
}

var productroutes = Router{
	Route{"RegisterProduct", http.MethodPost, constants.RegisterProductRoute, controllers.Registerproducts},
	Route{"Update Products", http.MethodPut, constants.UpdateproductRoute, controllers.UpdateProducts},
	Route{"Delete Products", http.MethodDelete, constants.DeleteProductRoute, controllers.DeleteProduct},
}

var globalproducts = Router{
	Route{"List-Product", http.MethodGet, constants.ListProductRoute, controllers.ListProductsController},
	Route{"Search-Product", http.MethodGet, constants.SearchProductRoute, controllers.SearchProducts},
}

var userauthroutes = Router{
	Route{"ADD TO CART", http.MethodPost, constants.AddToCartRoute, controllers.AddToCart},
	Route{"AddAddress", http.MethodPost, constants.AddAddressRoute, controllers.AddAddressOfUser},
	Route{"Get Single User", http.MethodGet, constants.GetSingleUserRoute, controllers.GetSingleUser},
	Route{"Update User", http.MethodPut, constants.UpdateUser, controllers.UpdateUser},
	Route{"Checkout Order", http.MethodPut, constants.CheckOutRoute, controllers.CheckOutOrder},
}
