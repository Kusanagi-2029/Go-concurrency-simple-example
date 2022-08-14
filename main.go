package main

import (
	"fmt"
)

/*
1) В данном случае при выполнении go-рутины планировщик пооставил её в очередь на выполнение
   Поэтому невсегда будет выводиться нужный нам текст.

func main() {
	go fmt.Println("Это - горутина")

	fmt.Println("Это главная горутина - функция main()")
}
*/

/*

2) Блокировка главной go-рутины
func main() {
	go fmt.Println("Это реализация конкуренции в Go через горутины")

	fmt.Println("Это главная горутина - функция main()")

	// Функция Sleep() в языке Go используется для блокирования/остановки главной go-рутины как минимум на указанное время
	time.Sleep(100 * time.Millisecond) // Блокировка/Остановка на 100 миллисекунд главной go-рутины, планировщик переключается на другую, ранее запланированную, go-рутину
}
*/

/*


3) Наглядный пример употребления go-рутины:

3.1) В данном случае программа парсинга URL выполняется около 7 секунд

func main() {
	t := time.Now()
	rand.Seed(t.UnixNano()) // Для генерации рандомных значений

	// Обработка двух сообщений происходит последовательно
	 parseURL("https://www.youtube.com/watch?v=w-eJDx-Lq1g")                   // Моё видео по Docker на Windows
	parseURL("https://github.com/Kusanagi-2029/Dockerfile-simple-DotNetCoreApp") // мой Github-репо по Docker на Windows

	fmt.Printf("Parsing completed. Time Elapesed: %.2f seconds\n", time.Since(t).Seconds())
}


3.2) А в этом случае программа парсинга URL выполняется около 3-4 секунд

func main() {
	t := time.Now()
	rand.Seed(t.UnixNano()) // Для генерации рандомных значений

	// Обработка двух сообщений происходит поОЧЕРЕДНО
	go parseURL("https://www.youtube.com/watch?v=w-eJDx-Lq1g")                   // Моё видео по Docker на Windows
	parseURL("https://github.com/Kusanagi-2029/Dockerfile-simple-DotNetCoreApp") // мой Github-репо по Docker на Windows

	fmt.Printf("Parsing completed. Time Elapesed: %.2f seconds\n", time.Since(t).Seconds())
}

// Функция, которая будет выводить информацию о статусе парсинга (которого на самом деле нет, ведь мы это имитируем)
func parseURL(url string) {
	// цикл на 5 итераций с информацией о статусе "парсинга" (его имитации)
	for i := 0; i < 5; i++ {
		latency := rand.Intn(500) + 500 // Рандомная пауза между итерациями (0.5 - 1 сек)

		time.Sleep(time.Duration(latency) * time.Millisecond)

		fmt.Printf("Parsing <%s> - Step %d - Latency %d\n", url, i+1, latency)
	}
}
*/

// 4) Работа с каналами - примитивами для синхронизации и обмена данными между Go-рутинами
// Горутины могут писать данные в канал и читать с него

/*

func main() {
	// Создадим канал, передающий значения строкового типа
	// Каналы инициализируются, как и слайсы, в main()
	// Неинициализированные каналы будут равны nil
	message := make(chan string)

	// Горутина, записывающая сообщение в канал с блокировкой перед этим на 2 сек.
	go func() {
		time.Sleep(2 * time.Second)
		message <- "hello" // Запись в канал
	}()

	msg := <-message // Чтение из канала
	fmt.Println(msg)

	// либо можно читать из канала сразу в вызове функции:
	//fmt.Println(<-message)
}

*/

/*
4.1) Запись в канал без горутины в функции main() произойдёт ошибка в runtime приложения - DEADLOCK
DEADLOCK происходит, когда при чтении из/записи в канал происходит блокировка, при этом ни одна
из доступных горутин никогда не прочитает данные из канала.

Пример #1: главная горутина - функция main() - ждёт, когда какая-нибудь горутина прочтёт из канала данные,
	       но этого никогда не произойдёт => DEADLOCK !

func main() {
	message := make(chan string)
	message <- "hello"
	fmt.Println(<-message)
}


Пример #2: числа успешно читаются, но после прочтения последненго элемента произойдёт DEADLOCK -
		   цикл пытается прочитать данные из канала, но в него уже никто ничего не пишет


func main() {
	message := make(chan string)
	go func() {
		for i := 1; i <= 10; i++ {
			message <- fmt.Sprintf("%d", i)
			time.Sleep(time.Millisecond * 500)
		}
	}()

	for {
		fmt.Println(<-message)
	}
}

*/

/*

4.2) Поэтому функция, которая выполняет запись в канал, после записи всех элементов должна ЗАКРЫТЬ этот канал.

func main() {
	message := make(chan string)
	go func() {
		for i := 1; i <= 10; i++ {
			message <- fmt.Sprintf("%d", i)
			time.Sleep(time.Millisecond * 500)
		}

		close(message) // ЗАКРЫТИЕ КАНАЛА
	}()

	// Переменная msg = TRUE, если канал инициализирован, FALSE - если его закрыли.
	for msg := range message {
		fmt.Println(msg)
	}
}
*/

// 5) Но канал message сейчас небуферизированный - каждая операция записи в канал БЛОКИРУЕТ горутины.

func main() {
	message := make(chan string, 2) // В качестве второго аргумента прописывается ДЛИНА БУФЕРА
	message <- "hello"

	fmt.Println(<-message) // В этой же (! ГЛАВНОЙ) горутине прочитаем из канала. DEADLOCK'a нет!
}
