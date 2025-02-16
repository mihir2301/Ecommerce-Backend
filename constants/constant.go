package constants

const (
	Sender   = "mihirg495423@gmail.com"
	Database = "ecommerse"
	DbUri    = "localhost:27017"

	HealthCheckRoute = "/health"

	//email verification constants

	VerifyEmailRoute = "/verify-email"
	VerifyOtpRoute   = "/verify-otp"
	ResendEmailRoute = "/resend-email"

	UserRegisterRoute = "/user-register"
	UserLoginRoute    = "/login"

	//product routes

	RegisterProductRoute = "/product-register"
	ListProductRoute     = "/list-products"
	SearchProductRoute   = "/search"
	UpdateproductRoute   = "/update-product"
	DeleteProductRoute   = "/delte-product"
	AddToCartRoute       = "/cart"
	AddAddressRoute      = "/address"
	GetSingleUserRoute   = "/user/:id"
	UpdateUser           = "/update-user"
	CheckOutRoute        = "/user/:id"
)

const (
	NormalUser = "user"
	AdminUser  = "admin"
)

const (
	Otpvalidation = 60
)

// collection Name
const (
	Verification      = "Verification"
	UsersCollection   = "User"
	ProductCollection = "products"
)
