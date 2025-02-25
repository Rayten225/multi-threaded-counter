package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
)

var (
	filePath      string
	wg            sync.WaitGroup
	mutex         sync.Mutex
	maxGoroutines = 10
	workers       int
)

// Структура для сортировки
type keyValue struct {
	Key   string
	Value int
}

func main() {
	fmt.Print("Введите путь к файлу: ")
	_, err := fmt.Scanln(&filePath)
	HandleErr(err)

	file, err := os.Open(filePath)
	defer file.Close()
	HandleErr(err)

	WordCount(file)
}

// Функция подсчета и сортировки слов
func WordCount(file *os.File) map[string]int {
	if file == nil {
		return nil
	}

	var sortedStruct []keyValue
	sem := make(chan struct{}, 10)
	words := make(map[string]int)
	scanner := bufio.NewScanner(file)
	var re = regexp.MustCompile(`[[:punct:]]`)

	for i := 0; scanner.Scan(); i++ {
		wg.Add(1)
		sem <- struct{}{}
		//Получает строку без спец символов
		str := strings.Split(re.ReplaceAllString(scanner.Text(), ""), " ")

		//Запускаем горутину на каждую строчку файла
		if maxGoroutines < workers {
			wg.Wait()
		} else {
			go func(str []string) {
				workers++
				defer wg.Done()
				defer func() { <-sem }()
				for _, word := range str {
					mutex.Lock()
					words[strings.ToLower(word)]++
					mutex.Unlock()
					workers--
				}
			}(str)
		}

	}
	wg.Wait()

	//Сортируем и выводим слова по их количеству
	for word, count := range words {
		sortedStruct = append(sortedStruct, keyValue{word, count})
	}
	sort.Slice(sortedStruct, func(i, j int) bool {
		return sortedStruct[i].Value > sortedStruct[j].Value
	})

	for _, kv := range sortedStruct {
		fmt.Printf("%s: %d\n", kv.Key, kv.Value)
	}
	return words
}

func HandleErr(err error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Ошибка открытия файла: %v", err)
	}
	defer file.Close()
}
