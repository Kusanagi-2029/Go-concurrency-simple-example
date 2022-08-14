##### Содержание  
* [Конкуренция в Go](#concurrency) 
* [Горутины](#goroutines)
* [Каналы](#channels)
* [DEADLOCK](#deadlock)
* [Буферизированный канал](#buffered_channel)

<a name="concurrency"><h2>Конкуренция в Go</h2></a>
Большие программы часто состоят из множества более мелких подпрограмм. Например, веб-сервер обрабатывает запросы, сделанные веб-браузерами, и в ответ обслуживает веб-страницы HTML. Каждый запрос обрабатывается как небольшая программа.

Было бы идеально, если бы подобные программы могли одновременно запускать свои более мелкие компоненты (в случае веб-сервера для обработки нескольких запросов). Одновременное выполнение более чем одной задачи называется параллелизмом. Go имеет богатую поддержку параллелизма с использованием горутин и каналов.

<a name="goroutines"><h2>Горутины</h2></a>
Горутина — это функция, которая может работать одновременно с другими функциями. Чтобы создать горутину, мы используем ключевое слово go, за которым следует вызов функции:

![Concurrency 1](https://user-images.githubusercontent.com/71845085/184555262-eedb577f-8da0-47a3-9bdd-256b6bf93b80.jpg)
В данном случае при выполнении go-рутины планировщик пооставил её в очередь на выполнение. Поэтому невсегда будет выводиться нужный нам текст.

### Блокировка главной go-рутины
Функция main() - главная горутина. После её завершения завершается всё остальное. 

![Concurrency 2](https://user-images.githubusercontent.com/71845085/184555339-22f32b76-c4a3-4882-9803-0d9c9d080612.jpg)
Функция Sleep() в языке Go используется для блокирования/остановки главной go-рутины как минимум на указанное время.
В этом же случае при блокировании главной горутины, выведется всё нам необходимое.

### Наглядный пример употребления горутины
#### Программа без горутины
* В данном случае программа парсинга URL выполняется около 7 секунд.
* Обработка двух сообщений происходит последовательно - сначала по первому URL, затем - по второму.

![Concurrency 3](https://user-images.githubusercontent.com/71845085/184555462-6017075a-535d-4393-86a4-0fdfdcd42745.jpg)

#### Программа с горутиной
* А в этом случае программа парсинга URL выполняется около 3-4 секунд.
* Обработка двух сообщений происходит поОЧЕРЕДНО.

![Concurrency 4](https://user-images.githubusercontent.com/71845085/184555508-7298d80f-6bde-41ac-b0ca-594ca1cf9f8e.jpg)

<a name="channels"><h2>Каналы</h2></a>
Каналы - примитивы для синхронизации и обмена данными между Go-рутинами
* Горутины могут писать данные в канал и читать из него.
* Каналы инициализируются, как и слайсы, в main().
* Неинициализированные каналы будут равны nil.

Создадим канал, передающий значения строкового типа:
```go
func main() {
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
```
Всё работает, запись "hello" выведется через 2 секунды.

<a name="deadlock"><h2>DEADLOCK</h2></a>
При записи в канал без горутины в функции main() произойдёт ошибка в runtime приложения - DEADLOCK.
DEADLOCK происходит, когда при чтении из/записи в канал происходит блокировка, при этом ни одна из доступных горутин никогда не прочитает данные из канала.
```go
Пример #1: главная горутина - функция main() - ждёт, когда какая-нибудь горутина прочтёт из канала данные,
	       но этого никогда не произойдёт => DEADLOCK !

func main() {
	message := make(chan string)
	message <- "hello"
	fmt.Println(<-message)
}
```

```go
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
```
![Concurrency 5](https://user-images.githubusercontent.com/71845085/184555819-5b144dd3-6e39-4103-bf1a-a3f60eafd354.jpg)

**Решение DEADLOCK'а по примеру 2:**
Функция, которая выполняет запись в канал, после записи всех элементов должна ЗАКРЫТЬ этот канал.
![Concurrency 6](https://user-images.githubusercontent.com/71845085/184555884-de62787f-3a5d-4da1-9fd4-a84d58dc7600.jpg)

* [Буферизированный канал](#buffered_channel)

<a name="buffered_channel"><h2>Буферизированные каналы</h2></a>
Данный канал message сейчас ***небуферизированный*** - каждая операция записи в канал БЛОКИРУЕТ горутины.

![Concurrency 7](https://user-images.githubusercontent.com/71845085/184555940-3b8fc0f7-6243-4b5e-9254-1d1c6618b5bd.jpg)

Как видно из скриншота, DEADLOCK'а не происходит.
