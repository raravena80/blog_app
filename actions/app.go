package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/ssl"
	"github.com/gobuffalo/envy"
	"github.com/unrolled/secure"

	"github.com/gobuffalo/buffalo/middleware/csrf"
	"github.com/gobuffalo/buffalo/middleware/i18n"
	"github.com/gobuffalo/packr"
	"github.com/raravena80/blog_app/models"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_blog_app_session",
		})
		// Automatically redirect to SSL
		app.Use(ssl.ForceSSL(secure.Options{
			SSLRedirect:     ENV == "production",
			SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		}))

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}
		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)
		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)
		// Remove to disable this.
		app.Use(middleware.PopTransaction(models.DB))

		// Setup and use translations:
		var err error
		if T, err = i18n.New(packr.NewBox("../locales"), "en-US"); err != nil {
			app.Stop(err)
		}
		app.Use(T.Middleware())
		app.Use(SetCurrentUser)

		app.GET("/", HomeHandler)

		app.ServeFiles("/assets", assetsBox)
		// users routes
		auth := app.Group("/users")
		auth.GET("/register", UsersRegisterGet)
		auth.POST("/register", UsersRegisterPost)
		auth.GET("/login", UsersLoginGet)
		auth.POST("/login", UsersLoginPost)
		auth.GET("/logout", UsersLogout)
		postGroup := app.Group("/posts")
		postGroup.GET("/index", PostsIndex)
		postGroup.GET("/create", AdminRequired(PostsCreateGet))
		postGroup.POST("/create", AdminRequired(PostsCreatePost))
		postGroup.GET("/detail/{pid}", PostsDetail)
		postGroup.GET("/edit/{pid}", AdminRequired(PostsEditGet))
		postGroup.POST("/edit/{pid}", AdminRequired(PostsEditPost))
		postGroup.GET("/delete/{pid}", AdminRequired(PostsDelete))
		commentsGroup := app.Group("/comments")
		commentsGroup.Use(LoginRequired)
		commentsGroup.POST("/create/{pid}", CommentsCreatePost)
		commentsGroup.GET("/edit/{cid}", CommentsEditGet)
		commentsGroup.POST("/edit/{cid}", CommentsEditPost)
		commentsGroup.GET("/delete/{cid}", CommentsDelete)
	}

	return app
}
