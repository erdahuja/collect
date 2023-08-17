package v1

import (
	"net/http"

	"collect/app/services/collect-api/handlers/v1/formgrp"
	"collect/app/services/collect-api/handlers/v1/questiongrp"
	"collect/app/services/collect-api/handlers/v1/responsegrp"
	"collect/app/services/collect-api/handlers/v1/usergrp"
	"collect/business/auth"
	"collect/business/core/answer"
	"collect/business/core/answer/repositories/answerdb"
	"collect/business/core/form"
	"collect/business/core/form/repositories/formdb"
	"collect/business/core/question"
	"collect/business/core/question/repositories/questiondb"
	"collect/business/core/response"
	"collect/business/core/response/repositories/responsedb"
	"collect/business/core/user"
	"collect/business/mid"
	"collect/foundation/config"
	"collect/foundation/event"
	"collect/foundation/kafka"
	"collect/foundation/web"

	"collect/business/core/user/repositories/usercache"
	"collect/business/core/user/repositories/userdb"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log          *zap.SugaredLogger
	DB           *sqlx.DB
	Cfg          *config.Configurations
	ServerErrors chan error
}

func Register(app *web.App, cfg Config) {
	const version = "v1"

	writer := kafka.GetKafkaWriter(cfg.ServerErrors, cfg.Cfg.DOMAIN_KAFKA_BROKER, cfg.Cfg.DOMAIN_KAFKA_TOPIC)
	eventCore := event.NewCore(cfg.Log, writer)

	a := auth.New(eventCore, auth.Config{
		Log:    cfg.Log,
		DB:     cfg.DB,
		Secret: []byte(cfg.Cfg.Secret),
	})

	authen := mid.Authenticate(a)
	ruleAdmin := mid.Authorize(a, auth.RuleAdminOnly)
	ruleAny := mid.Authorize(a, auth.RuleAny)
	ruleUser := mid.Authorize(a, auth.RuleUserOnly)
	ruleCollector := mid.Authorize(a, auth.RuleCollectorOnly)

	ugh := usergrp.Handlers{
		User: user.NewCore(eventCore, usercache.NewStore(cfg.Log, userdb.NewStore(cfg.Log, cfg.DB))),
		Auth: a,
	}
	app.Handle(http.MethodGet, version, "/users/token", ugh.Token)                // public api on basic auth, to get auth token
	app.Handle(http.MethodPost, version, "/users", ugh.Create, authen, ruleAdmin) // add users (only by admin)
	app.Handle(http.MethodGet, version, "/users", ugh.Query, authen, ruleAdmin)   // all users

	frmCore := form.NewCore(cfg.Log, eventCore, formdb.NewStore(cfg.Log, cfg.DB))
	queCore := question.NewCore(cfg.Log, eventCore, questiondb.NewStore(cfg.Log, cfg.DB))
	fgh := formgrp.Handlers{
		Form:     frmCore,
		Question: queCore,
		Auth:     a, // TODO: use check user to form permission, since this is resouce specific instead of api, we need to check permission in route handlers
	}
	app.Handle(http.MethodPost, version, "/forms", fgh.Create, authen, ruleAdmin)                                       // create form, only admin
	app.Handle(http.MethodDelete, version, "/forms", fgh.Delete, authen, ruleAdmin)                                     // delete forms
	app.Handle(http.MethodGet, version, "/forms", fgh.Query, authen, ruleAdmin)                                         // view forms
	app.Handle(http.MethodGet, version, "/forms/:form_id/questions", fgh.QueryQuestionsByFormID, authen, ruleCollector) // view questions of a form

	cgh := questiongrp.Handlers{
		Question: queCore,
	}
	app.Handle(http.MethodPost, version, "/questions", cgh.Create, authen, ruleAdmin)               // add a new question
	app.Handle(http.MethodGet, version, "/questions/:quesID", cgh.QueryByID, authen, ruleCollector) // view questions by quesID
	app.Handle(http.MethodDelete, version, "/questions", cgh.DeleteQuestion, authen, ruleAdmin)     // remove question from form


	respCore := response.NewCore(cfg.Log, eventCore, responsedb.NewStore(cfg.Log, cfg.DB))
	rgh := responsegrp.Handlers{
		Response: respCore,
		Answer:   answer.NewCore(cfg.Log, eventCore, answerdb.NewStore(cfg.Log, cfg.DB)),
	}
	app.Handle(http.MethodPost, version, "/response", rgh.Create, authen, ruleCollector)                          // create form, only admin
	app.Handle(http.MethodGet, version, "/response/:form_id", rgh.QueryByFormID, authen, ruleUser, ruleCollector) // view responses of a form
	app.Handle(http.MethodPost, version, "/response/:id/answer", rgh.CreateAnswer, authen, ruleCollector)

	app.Handle(http.MethodGet, version, "/status", ugh.Status, ruleAny)

}
