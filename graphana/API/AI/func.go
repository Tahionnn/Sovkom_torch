package AI

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coalaura/mistral"
	"log"
	"strings"
)

type DataForGraphicForUser struct {
	ID    int    `json:"id"`
	Value int    `json:"value"`
	Label string `json:"label"`
}

type Input struct {
	GTParse GTParse `json:"gt_parse"`
}

type GTParse struct {
	Shop    string   `json:"shop"`
	Date    string   `json:"date"`
	Time    []string `json:"time"`
	Items   []Item   `json:"items"`
	Overall string   `json:"overall"`
}

type Item struct {
	Name        string  `json:"name"`
	Measurement string  `json:"measurement"`
	Count       int     `json:"count"`
	Price       float64 `json:"price"`
	Overall     float64 `json:"overall"`
}

type SubSend struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

const prompt = "Ты - система категоризации товаров. \n\nУ тебя есть список товаров в формате [{\"name\": \"Название товара\", \"count\": Количество}]. \nРаспредели все товары по следующим категориям: \n['Одежда и обувь', 'Продукты', 'Автомобильные товары', 'Электроника и бытовая техника', 'Дом и интерьер', 'Красота и уход', 'Здоровье и медицина', 'Спорт и отдых', 'Детские товары', 'Хобби и творчество'].\n\nЕсли товар не подходит ни под одну категорию, отнеси его в категорию 'Другое'.\n\nВ ответе обязательно верни только JSON в формате:\n{\n  \"Категория\": сумма товаров\n}\n\nНикакого дополнительного текста, объяснений или форматирования. Только чистый JSON.\n"

func FindCategories(client *mistral.MistralClient, jsonFile string) ([]DataForGraphicForUser, error) {
	var input Input
	var dataSend []SubSend

	// Распарсим весь входящий JSON
	err := json.Unmarshal([]byte(jsonFile), &input)
	if err != nil {
		log.Fatalf("Ошибка парсинга JSON: %v", err)
	}

	// Теперь достаем только items
	for _, item := range input.GTParse.Items {
		dataSend = append(dataSend, SubSend{
			Name:  item.Name,
			Count: item.Count,
		})
	}

	// Переводим items в JSON для отправки
	jsonData, err := json.Marshal(dataSend)
	if err != nil {
		log.Printf("Ошибка конвертации товаров в JSON: %s", err)
		return nil, err
	}

	// Формируем запрос в модель
	request := mistral.ChatCompletionRequest{
		Model: mistral.ModelMistralSmall,
		Messages: []mistral.Message{
			{
				Role:    mistral.RoleSystem,
				Content: "Система категоризации товаров",
			},
			{
				Role:    mistral.RoleUser,
				Content: prompt + "\nВот мои товары: " + string(jsonData),
			},
		},
	}

	response, err := client.Chat(request)
	if err != nil {
		log.Printf("Ошибка ответа Mistral AI: %s", err)
		return nil, err
	}

	fmt.Println(response.Choices[0].Message.Content)

	result, err := CleanAndParseJSON(response.Choices[0].Message.Content)
	if err != nil {
		log.Fatalf("Ошибка очистки/валидации JSON: %v", err)
	}
	//fmt.Println(result)

	// Заполняем data
	var data []DataForGraphicForUser
	id := 0
	for label, value := range result {
		data = append(data, DataForGraphicForUser{
			ID:    id,
			Value: value,
			Label: label,
		})
		id++
	}

	return data, nil
}

// CleanAndParseJSON очищает строку: убирает первую и последнюю строки и парсит в map[string]int
func CleanAndParseJSON(raw string) (map[string]int, error) {
	// Разбиваем по строкам
	lines := strings.Split(raw, "\n")

	// Если строк меньше 3 — явно недостаточно для нормального JSON
	if len(lines) < 3 {
		return nil, errors.New("слишком мало строк для обработки")
	}

	// Убираем первую и последнюю строки
	lines = lines[1 : len(lines)-1]

	// Склеиваем обратно в одну строку
	cleaned := strings.Join(lines, "\n")

	// Парсим JSON
	result := make(map[string]int)
	err := json.Unmarshal([]byte(cleaned), &result)
	if err != nil {
		return nil, errors.New("ответ невалидный JSON или неправильный формат")
	}

	return result, nil
}
