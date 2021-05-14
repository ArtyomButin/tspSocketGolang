package main

import "net"
import "fmt"
import "bufio"
import "strings" // требуется только ниже для обработки примера

func main() {

	fmt.Println("Launching server...")

	// Устанавливаем прослушивание порта
	ln, _ := net.Listen("tcp", ":8081")

	// Открываем порт
	conn, _ := ln.Accept()

	// Запускаем цикл
	for {
		// Будем прослушивать все сообщения разделенные \n
		message, _ := bufio.NewReader(conn).ReadString('\n')
		// Распечатываем полученое сообщение
		fmt.Print("Message Received:", message)
		// Процесс выборки для полученной строки
		newMessage := strings.ToUpper(message)
		// Отправить новую строку обратно клиенту
		_, err := conn.Write([]byte(newMessage + "\n"))
		if err != nil {
			return
		}
	}
}