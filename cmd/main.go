package main

import (
	"log"
	"time"

	goprogressbar "github.com/RodrigoPetter/go-progress-bar"
)

func main() {
	// Create a simple progress bar instance
	bar := goprogressbar.NewProgressBar("Barra de testes 1")

	//configure the log writter
	log.SetOutput(bar.GetLogWriter())

	bar2 := bar.NewSubBar("Subtask da barra 1")
	bar3 := bar.NewSubBar("este é um nome longo para testar o padding do do nome da var ")
	bar4 := bar.NewSubBar("Teste 1234")
	bar5 := bar4.NewSubBar("Mais um teste para verificar o que ta rolando", goprogressbar.WithStages(5))

	goprogressbar.StartRenderLoop(bar)

	for i := 0; i <= 100; i++ {

		log.Println(i, "Hello from standard logger!")

		bar.Increment(2)
		bar2.Increment(2)
		bar3.Increment(2)

		if i == 3 {
			bar.UpdateProgressBar(goprogressbar.WithStages(30))
			log.Println("BAR 5 CONCLUÍDA")
			bar5.Finish()
			bar4.NewSubBar("Mais um teste para verificar o que ta rolando", goprogressbar.WithStages(5))
		} else if i < 3 {
			log.Println("Incrementando bar5")
			bar5.Increment(1)
		}
		time.Sleep(1 * time.Second)
	}
	select {}

}
