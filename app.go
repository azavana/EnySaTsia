package main

import (
	R "voting/routes"

	"github.com/kataras/iris"
)

// App -- the application root
func App() *iris.Application {
	app := iris.New()

	app.Get("/", func(ctx iris.Context) {
		ctx.JSON("every thing ok")
	})

	s := app.Party("/session")
	{
		s.Post("/new", R.NewSession)
		s.Put("/update", R.UpdateSession)
		s.Post("/start", R.StartSession)
		s.Post("/close", R.CloseSession)
		s.Get("/", R.GetAllSession)
		s.Get("/{id: string}", R.GetSession)
	}

	q := app.Party("/question")
	{
		q.Post("/new", R.NewQuestion)
		q.Put("/update", R.UpdateQuestion)
		q.Post("/startVote", R.VoteStart)
		q.Post("/closeVote", R.VoteClose)
		q.Post("/vote", R.Vote)
		q.Get("/{question: string}", R.GetQuestion)
		q.Get("/session/{session: string}", R.GetQuestionOfSession)
	}
	return app
}

func main() {
	app := App()
	app.Run(iris.Addr(":8080"))
}
