package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")

	// for i := 0; i < 10; i++ {
	// 	fmt.Println(i)
	// }
	var request HttpRequest = HttpRequest{ContentType: "json"}
	fmt.Println(request.ContentType)
	var server listener = listener{port: ":8080", protocol: "tcp"}
	server.start()
}
