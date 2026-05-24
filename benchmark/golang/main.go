package main

import (
	"fmt"
	"net/http"
	"runtime"
)

// var logChannel chan string

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Envia o log para o canal (não trava)
	// logChannel <- fmt.Sprintf("Recebida requisição: %s %s", r.Method, r.URL.Path)

	fmt.Fprintln(w, "Hello, World!")
}

// func logWorker() {
// 	for message := range logChannel {
// 		log.Println(message)
// 	}
// }

func main() {
	// Limitar a 1 CPU
	runtime.GOMAXPROCS(8)

	// Criar o canal de logs com buffer para evitar bloqueios
	// logChannel = make(chan string, 100)

	// Iniciar o worker em uma goroutine separada
	// go logWorker()

	http.HandleFunc("/", helloHandler)

	addr := ":8081"
	// logChannel <- fmt.Sprintf("Servidor rodando em http://localhost%s usando 1 CPU")

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		// logChannel <- fmt.Sprintf("Erro ao iniciar o servidor: %v", err)
	}
}
