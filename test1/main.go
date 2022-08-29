package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/yaml.v2"
)

type Conf struct { // кстати, тут была буква C русская, поэтому у тебя возникала проблема, когда программа не могла найти
	// структуру "Conf" и тебе приходилось копировать название и вставлять
	Update int `yaml:"update"`
	Every  int `yaml:"every"`
}

// для себя в будущем и для других людей, которые будут смотреть этот код лучше всего писать описание функций, например:
// getConf - возвращает структуру конфиг файла
func GetConf(fileName string) (*Conf, error) { // желательно, чтобы функции, в названии которых имеется "get" что-то отдавали
	var c = new(Conf)

	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		// fmt.Println("yamlFile.Get err   %v ", err)

		// вместо строчки выше лучше использовать ту, что ниже, т.к. иначе при возникновении ошибки код у тебя не остановится,
		// а продолжит выполняться, что приведёт к критической ошибке

		return nil, err // при возникновении ошибки возвращаем пустую структуру и саму ошибку
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func main() {
	// обработка паники
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r) // если возникнет критическая ошибка - программа выведет её содержимое и завершится на следующей строчке
			os.Exit(1)                       // код 1 для выхода из программы через функцию os.Exit() означает, что программу завершаем с ошибкой
			// для того, чтобы в любом моменте кода завершить программу без ошибки можно вызвать os.Exit(0)
		}
	}()

	c, err := GetConf("conf.yaml")
	if err != nil {
		fmt.Printf("ошибка чтения конфига: %v\n", err) // при вызове fmt.Printf сочетание "%v" выведет любое значение после запятой в истинном виде (%d для числа, %s для строки)
		return                                         // return в функции main() практически то же самое, что и os.Exit(0) , т.е. выход из программы без ошибки
	}

	go c.timer()
	go c.update()

	// следующий код нужен только для того, чтобы программа завершалась только тогда, когда
	// ты в консоли нажмёшь "ctrl + c". Обычно я его просто копирую во все свои программы и вставляю в функцию main

	sigs := make(chan os.Signal, 1)
	signal.Notify(
		sigs,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	<-sigs
}

// timer - функция для подсчёта и вывода времени с начала запуска программы
func (c *Conf) timer() {
	/*
		Кстати, длинные комментарии можно писать и таким способом.
		Для подсчёта времени лучше использовать функцию time.After(), т.к. она работает побыстрее и выглядит в коде более наглядно (легче понять, что происходит в методе timer())
	*/
	var startTime = time.Now()

	// в данном случае структура select будет ожидать, когда пройдёт c.Every секунд, после этого
	// начнёт выполняться код в теле case
	for {
		select {
		case <-time.After(time.Second * time.Duration(c.Every)):
			timeCount := time.Since(startTime).Seconds()
			fmt.Printf("время с запуска %v секунд\n", int(timeCount))
		}
	}
}

// update - функция для перечитывания конфиг файла каждые c.Update секунд
func (c *Conf) update() {
	for {
		select {
		case <-time.After(time.Second * time.Duration(c.Every)):
			newConf, err := GetConf("conf.yaml")
			if err != nil {
				fmt.Printf("Ошибка перечитывания конфига: %v\n", err)
				os.Exit(1) // завершаем программу при возникновении ошибки чтения конфига
			}

			// нельзя просто написать "c = newConf", т.к. функция update() вызвана у уже существующей структуры "c" и замена
			// этой переменной не повлияет на все уже запущенные функции этой структуры, поэтому просто изменим переменные нашей структуры
			c.Every = newConf.Every
			c.Update = newConf.Update
		}
	}
}
