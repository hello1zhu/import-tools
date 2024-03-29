package router

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/bangwork/import-tools/serve/controllers"
	"github.com/bangwork/import-tools/serve/middlewares"
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
)

//go:embed dist lang
var FS embed.FS

func Run(port int) {
	gin.SetMode(gin.DebugMode)
	api := gin.Default()
	api.Use(middlewares.Recovery(), middlewares.Logger())
	api.Use(GinI18nLocalize())
	api.Use(middlewares.Cors())

	temple := template.Must(template.New("").ParseFS(FS, "dist/index.html"))
	api.SetHTMLTemplate(temple)

	fe, err := fs.Sub(FS, "dist/assets")
	if err != nil {
		log.Println("embed dist assets err", err)
		return
	}
	api.StaticFS("/assets", http.FS(fe))

	fe, _ = fs.Sub(FS, "dist/public")
	api.StaticFS("/public", http.FS(fe))

	api.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	api.POST("/check_path_exist", controllers.CheckPathExist)
	api.POST("/check_jira_path_exist", controllers.CheckJiraPathExist)
	api.POST("/jira_backup_list", controllers.JiraBackUpList)
	api.POST("/resolve/start", controllers.StartResolve)
	api.GET("/resolve/progress", controllers.ResolveProgress)
	api.POST("/resolve/stop", controllers.StopResolve)
	api.GET("/resolve/result", controllers.ResolveResult)
	api.GET("/project_list", controllers.ProjectList)
	api.POST("/choose_team", controllers.ChooseTeam)
	api.POST("/issue_type_list", controllers.IssueTypeList)
	api.POST("/set_share_disk", controllers.SetShardDisk)

	api.GET("/import/reset", controllers.Reset)
	api.POST("/import/start", controllers.StartImport)
	api.POST("/import/pause", controllers.PauseImport)
	api.POST("/import/continue", controllers.ContinueImport)
	api.POST("/import/stop", controllers.StopImport)
	api.GET("/import/progress", controllers.ImportProgress)
	api.GET("/import/log", controllers.GetAllImportLog)
	api.GET("/import/log/start_line/:start_line", controllers.GetImportLog)
	api.GET("/import/log/download/all", controllers.DownloadLogFile)
	api.GET("/import/log/download/current", controllers.DownloadCurrentLogFile)
	api.GET("/import/scope", controllers.GetScope)

	api.Run(fmt.Sprintf(":%d", port))
}

func GinI18nLocalize() gin.HandlerFunc {
	return i18n.Localize(
		i18n.WithBundle(&i18n.BundleCfg{
			RootPath:         "./lang",
			AcceptLanguage:   []language.Tag{language.Chinese, language.English},
			DefaultLanguage:  language.Chinese,
			FormatBundleFile: "toml",
			UnmarshalFunc:    toml.Unmarshal,
			Loader:           &i18n.EmbedLoader{FS: FS},
		}),
	)
}
