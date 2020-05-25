package main

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"time"
	"sync"	
	"os"
	"bufio"
	"strconv"
)
import s "strings"

//функция выполнения GET-запроса
func MakeRequest(url string) string{
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(body)

}

func main() {
	var wg sync.WaitGroup //переменная для ожидания выполнения всех горутин
	var c chan int = make(chan int, 5) // создаём канал максимум на 5 одновременных значений
	var result int = 0 //переменная для подсчёта результата
	var urls = listenStdin() //массив с url'ами из stdin
	
	fmt.Println("Старт программы") //выводим результат
	//цикл по элементам массив (url'ам)
	for i := 0; i < len(urls); i++ { 
		wg.Add(1) //добавили в группу ожидания 1
		go countGo(urls[i],&wg,c) //выполняем горутину, передавая в неё url,группу ожидания и канал
		result = result + <- c //подсчёт результатов
	}

	fmt.Println("ждём выполнения всех горутин...") //выводим результат
	wg.Wait() //ждём выполнения всех горутин
	fmt.Println("Результат: "+strconv.Itoa(result)) //выводим результат
}

//функция для считывания url'ов из stdin
func listenStdin() []string{
	var urlsArr []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var str = scanner.Text()
		urlsArr = s.Split(str, `\n`)
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	return urlsArr
}

//функция для подсчёта количества вхождений 'go' в теле ответа url'а
func countGo(url string, wg *sync.WaitGroup, c chan int) {
	defer wg.Done() //удаляем из группы ожидания 1
	res := MakeRequest(url) //делаем запрос по url
	time.Sleep(time.Millisecond*300) //ждём 300 мс
	count := s.Count(res,"go") //подсчитываем количество вхождений 'go' в теле ответа
	log.Println("Вхождений 'go' в теле ответа URL: "+strconv.Itoa(count))
	c <- count //передаём количество вхождений в канал
}
