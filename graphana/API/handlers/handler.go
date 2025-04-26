package handlers

import (
	"HacatonSovKomBank/AI"
	"HacatonSovKomBank/middleware"
	"HacatonSovKomBank/pkg/jwtutils"
	"HacatonSovKomBank/pkg/postgreSQL"
	"encoding/json"
	"github.com/coalaura/mistral"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

const (
	loginURL               = "/login/"
	registerURL            = "/registration/"
	listResForUsernameURL  = "/receipts/"
	analyticsForReceiptURL = "/receipt/{uuid}"
	analyticsWeekURL       = "/analytics/week/"
	analyticsMonthURL      = "/analytics/month/"
	analyticsYearURL       = "/analytics/year/"
)

type handler struct {
	db       *gorm.DB
	clientAI *mistral.MistralClient
}

type AuthInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// NewHandler создает новый обработчик с подключением к базе данных
func NewHandler(DB *gorm.DB, client *mistral.MistralClient) Handler {
	return &handler{
		db:       DB,
		clientAI: client,
	}
}

// Register регистрирует все маршруты для авторизации, пользователей и чатов.
func (h *handler) Register(router *gin.Engine) {
	// Публичные маршруты (без JWT)
	public := router.Group("/")
	{
		// Регистрация и Login
		public.POST(registerURL, h.registration)
		public.POST(loginURL, h.login)

		// TODO Создание пустого чека с флос значением

		// TODO Изменение контента

	}

	// Защищённые маршруты (с JWT)
	private := router.Group("/")
	private.Use(middleware.AuthMiddleware()) // Проверка JWT
	{
		// Запрос списка чеков на человека
		private.GET(listResForUsernameURL, h.listResForUsername)

		// Запрос метрики по дате week - month - year
		private.GET(analyticsWeekURL, h.analyticsWeek)
		private.GET(analyticsMonthURL, h.analyticsMonth)
		private.GET(analyticsYearURL, h.analyticsYear)

		// Метрики по определенному чеку
		router.GET(analyticsForReceiptURL, h.analyticsForReceipt)

		// TODO Запрос на рекомендации по всем чекам за последний месяц

		// ТЕСТОВЫЙ ДЛЯ JWT
		private.GET("/test/", h.test)

		// Добавление искусственных данных пользователю
		private.GET("/testdata/", h.AddTestReceiptsHandler)

		//TODO Отправка фото чека на бэкенд в формате FormData, я долежн его переправить на ML, где его распарсят в чек
	}

}

// jpg -> JSON передаем ее обратно на валидацию пользователя -> Передаю JSON в рекомендательную систему;

// Регистрация пользователя
// Registration godoc

// @Summary Регистрация нового пользователя
// @Description Регистрирует нового пользователя, хэшируя пароль и сохраняя в БД
// @Tags auth
// @Accept json
// @Produce json
// @Param input body AuthInput true "Данные для регистрации"
// @Success 200 {object} map[string]string "Пользователь успешно зарегистрирован"
// @Failure 400 {object} map[string]string "Неверный ввод или пользователь уже существует"
// @Failure 500 {object} map[string]string "Ошибка сервера при создании пользователя"
// @Router /registration/ [post]
func (h *handler) registration(c *gin.Context) {
	var input AuthInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("[REGISTRATION] Ошибка парсинга JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ввод"})
		return
	}
	log.Printf("[REGISTRATION] Получен запрос на регистрацию: %s", input.Username)

	var user postgreSQL.User
	if err := h.db.Where("username = ?", input.Username).First(&user).Error; err == nil {
		log.Printf("[REGISTRATION] Попытка регистрации существующего пользователя: %s", input.Username)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь уже существует"})
		return
	}

	hashedPassword, err := jwtutils.HashPassword(input.Password)
	if err != nil {
		log.Printf("[REGISTRATION] Ошибка хэширования пароля: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при хэшировании пароля"})
		return
	}

	newUser := postgreSQL.User{
		Username:     input.Username,
		PasswordHash: hashedPassword,
	}

	if err := h.db.Create(&newUser).Error; err != nil {
		log.Printf("[REGISTRATION] Ошибка создания пользователя в БД: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании пользователя"})
		return
	}

	log.Printf("[REGISTRATION] Пользователь успешно зарегистрирован: %s", input.Username)
	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно зарегистрирован"})
}

// Авторизация пользователя
// Login godoc

// @Summary Логин пользователя
// @Description Авторизует пользователя и возвращает JWT токен с username в payload
// @Tags auth
// @Accept json
// @Produce json
// @Param input body AuthInput true "Данные для входа"
// @Success 200 {object} map[string]string "Токен успешно создан"
// @Failure 400 {object} map[string]string "Неверный ввод"
// @Failure 401 {object} map[string]string "Неверный логин или пароль"
// @Failure 500 {object} map[string]string "Ошибка сервера при создании токена"
// @Router /login/ [post]
func (h *handler) login(c *gin.Context) {
	var input AuthInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("[LOGIN] Ошибка парсинга JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ввод"})
		return
	}
	log.Printf("[LOGIN] Получен запрос на логин: %s", input.Username)

	var user postgreSQL.User
	if err := h.db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		log.Printf("[LOGIN] Пользователь не найден: %s", input.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		return
	}

	if !jwtutils.CheckPasswordHash(input.Password, user.PasswordHash) {
		log.Printf("[LOGIN] Неверный пароль для пользователя: %s", input.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		return
	}

	token, err := jwtutils.GenerateJWT(user.Username)
	if err != nil {
		log.Printf("[LOGIN] Ошибка генерации токена для пользователя %s: %v", input.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании токена"})
		return
	}

	log.Printf("[LOGIN] Токен успешно сгенерирован для пользователя: %s", input.Username)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// @Summary Тестовый маршрут
// @Description Этот маршрут используется для проверки работоспособности сервера, доступен только с валидным JWT токеном.
// @Tags test
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]string "Сообщение об успешной проверке"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Неавторизованный доступ"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /test/ [get]
func (h *handler) test(c *gin.Context) {
	usernameInterface, exists := c.Get("username")
	if !exists {
		c.JSON(401, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	username, ok := usernameInterface.(string)
	if !ok {
		c.JSON(500, gin.H{"error": "Ошибка сервера: неверный тип данных username"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Mess": username})
}

// @Summary Получить список всех чеков пользователя и аналитику по ним
// @Description Возвращает список всех чеков для текущего авторизованного пользователя и данные на Dashboard.
// @Tags receipts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "OK"
// @Failure 401 {object} map[string]string "Пользователь не авторизован"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /receipts [get]
func (h *handler) listResForUsername(c *gin.Context) {
	// Достаём username из контекста
	usernameInterface, exists := c.Get("username")
	if !exists {
		c.JSON(401, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	username, ok := usernameInterface.(string)
	if !ok {
		c.JSON(500, gin.H{"error": "Ошибка сервера: неверный тип данных username"})
		return
	}

	// Ищем чеки по username
	var receipts []postgreSQL.Receipt
	if err := h.db.Where("username = ?", username).Find(&receipts).Error; err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при поиске чеков: " + err.Error()})
		return
	}

	// Получаем все UUID чеков
	var uuids []string
	for _, receipt := range receipts {
		uuids = append(uuids, receipt.UUID)
	}

	// Получаем категории из ReceiptCategories на основе найденных UUID
	var receiptCategories []postgreSQL.ReceiptCategories
	if err := h.db.Where("uuid IN ?", uuids).Find(&receiptCategories).Error; err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при получении категорий: " + err.Error()})
		return
	}

	// Подготовим суммирование по категориям
	categorySums := make(map[string]int)

	// Для каждой категории из чеков суммируем значения
	for _, rc := range receiptCategories {
		// Парсим JSON в map[string]int
		var categories map[string]int
		if err := json.Unmarshal(rc.JSONCat, &categories); err != nil {
			log.Printf("Ошибка парсинга JSONCat для чека %s: %v", rc.UUID, err)
			continue
		}

		// Суммируем значения
		for category, value := range categories {
			categorySums[category] += value
		}
	}

	// Структура для отображения аналитики категорий
	type CategoryData struct {
		ID    int    `json:"id"`
		Value int    `json:"value"`
		Label string `json:"label"`
	}

	var analytics []CategoryData
	id := 0

	// Формируем список категорий с суммами
	for category, sum := range categorySums {
		if sum > 0 { // Фильтруем категории с ненулевыми суммами
			analytics = append(analytics, CategoryData{
				ID:    id,
				Value: sum,
				Label: category,
			})
			id++
		}
	}

	// Возвращаем список чеков и аналитики
	c.JSON(200, gin.H{
		"receipts":  receipts,  // Массив всех чеков
		"analytics": analytics, // Аналитика по категориям
	})
}

// @Summary Аналитика по конкретному чеку
// @Description Анализирует JSON чека и возвращает разбивку по категориям товаров и чек в формате JSON
// @Tags analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Результат анализа чека"
// @Failure 500 {object} map[string]string "Ошибка при обработке чека"
// @Router /analytics/receipt [get]
func (h *handler) analyticsForReceipt(c *gin.Context) {
	// Извлекаем UUID из параметра запроса
	receiptUUID := c.Param("uuid")
	if receiptUUID == "" {
		c.JSON(400, gin.H{"error": "UUID не указан"})
		return
	}

	// Ищем чек по UUID в таблице Receipt
	var receipt postgreSQL.Receipt
	if err := h.db.Where("uuid = ?", receiptUUID).First(&receipt).Error; err != nil {
		c.JSON(404, gin.H{"error": "Чек не найден"})
		return
	}

	// Ищем категории для этого UUID в таблице ReceiptCategories
	var receiptCategory postgreSQL.ReceiptCategories
	if err := h.db.Where("uuid = ?", receiptUUID).First(&receiptCategory).Error; err != nil {
		c.JSON(404, gin.H{"error": "Категории для этого чека не найдены"})
		return
	}

	// Формируем ответ
	c.JSON(200, gin.H{
		"receipts":  receipt.JSONCheck,       // JSONCheck из Receipt
		"analytics": receiptCategory.JSONCat, // JSONCat из ReceiptCategories
	})
}

// @Summary Еженедельная аналитика расходов
// @Description Возвращает сумму расходов по категориям за текущую неделю для авторизованного пользователя и все чеки за неделю
// @Tags analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} map[string]string "Результат анализа за неделю"
// @Failure 401 {object} map[string]string "Пользователь не авторизован"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /analytics/week [get]
func (h *handler) analyticsWeek(c *gin.Context) {
	// Достаём username из контекста
	usernameInterface, exists := c.Get("username")
	if !exists {
		c.JSON(401, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	username, ok := usernameInterface.(string)
	if !ok {
		c.JSON(500, gin.H{"error": "Ошибка сервера: неверный тип данных username"})
		return
	}

	// Получаем текущую дату и вычисляем дату начала недели
	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	// Ищем чеки пользователя за последнюю неделю, где данные JSON готовы
	var receipts []postgreSQL.Receipt
	if err := h.db.Where("username = ? AND date >= ? AND date < ? AND is_json_ready = true", username, startOfWeek, endOfWeek).
		Find(&receipts).Error; err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при получении чеков: " + err.Error()})
		return
	}

	// Если нет чеков, возвращаем пустой результат
	if len(receipts) == 0 {
		c.JSON(200, gin.H{"receipts": []any{}, "analytics": []any{}})
		return
	}

	// Массив чеков для ответа
	var receiptsData []datatypes.JSON
	for _, receipt := range receipts {
		receiptsData = append(receiptsData, receipt.JSONCheck)
	}

	// Собираем UUID всех чеков для дальнейшего поиска категорий
	var uuids []string
	for _, receipt := range receipts {
		uuids = append(uuids, receipt.UUID)
	}

	// Ищем категории для всех чеков
	var receiptCategories []postgreSQL.ReceiptCategories
	if err := h.db.Where("uuid IN ?", uuids).Find(&receiptCategories).Error; err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при получении категорий: " + err.Error()})
		return
	}

	// Суммируем значения по категориям
	categorySums := make(map[string]int)
	for _, rc := range receiptCategories {
		// Распаковываем JSON в map
		var categories map[string]int
		if err := json.Unmarshal(rc.JSONCat, &categories); err != nil {
			log.Printf("Ошибка парсинга категорий для чека %s: %v", rc.UUID, err)
			continue
		}

		// Суммируем значения
		for category, value := range categories {
			categorySums[category] += value
		}
	}

	// Формируем результат для категорий
	type CategoryData struct {
		ID    int    `json:"id"`
		Value int    `json:"value"`
		Label string `json:"label"`
	}

	var categoryResult []CategoryData
	id := 0
	for category, sum := range categorySums {
		if sum > 0 {
			categoryResult = append(categoryResult, CategoryData{
				ID:    id,
				Value: sum,
				Label: category,
			})
			id++
		}
	}

	// Возвращаем данные
	c.JSON(200, gin.H{
		"receipts":  receiptsData,
		"analytics": categoryResult,
	})
}

// @Summary Ежемесячная аналитика расходов
// @Description Возвращает сумму расходов по категориям за текущий месяц для авторизованного пользователя и все чеки за месяц
// @Tags analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} map[string]string "Результат анализа за месяц"
// @Failure 401 {object} map[string]string "Пользователь не авторизован"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /analytics/month [get]
func (h *handler) analyticsMonth(c *gin.Context) {
	// Достаём username из контекста
	usernameInterface, exists := c.Get("username")
	if !exists {
		c.JSON(401, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	username, ok := usernameInterface.(string)
	if !ok {
		c.JSON(500, gin.H{"error": "Ошибка сервера: неверный тип данных username"})
		return
	}

	// Получаем текущую дату и вычисляем дату начала месяца
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0) // следующий месяц

	// Ищем чеки пользователя за последний месяц, где данные JSON готовы
	var receipts []postgreSQL.Receipt
	if err := h.db.Where("username = ? AND date >= ? AND date < ? AND is_json_ready = true", username, startOfMonth, endOfMonth).
		Find(&receipts).Error; err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при получении чеков: " + err.Error()})
		return
	}

	// Если нет чеков, возвращаем пустой результат
	if len(receipts) == 0 {
		c.JSON(200, gin.H{"receipts": []any{}, "analytics": []any{}})
		return
	}

	// Массив чеков для ответа
	var receiptsData []datatypes.JSON
	for _, receipt := range receipts {
		receiptsData = append(receiptsData, receipt.JSONCheck)
	}

	// Собираем UUID всех чеков для дальнейшего поиска категорий
	var uuids []string
	for _, receipt := range receipts {
		uuids = append(uuids, receipt.UUID)
	}

	// Ищем категории для всех чеков
	var receiptCategories []postgreSQL.ReceiptCategories
	if err := h.db.Where("uuid IN ?", uuids).Find(&receiptCategories).Error; err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при получении категорий: " + err.Error()})
		return
	}

	// Суммируем значения по категориям
	categorySums := make(map[string]int)
	for _, rc := range receiptCategories {
		// Распаковываем JSON в map
		var categories map[string]int
		if err := json.Unmarshal(rc.JSONCat, &categories); err != nil {
			log.Printf("Ошибка парсинга категорий для чека %s: %v", rc.UUID, err)
			continue
		}

		// Суммируем значения
		for category, value := range categories {
			categorySums[category] += value
		}
	}

	// Формируем результат для категорий
	type CategoryData struct {
		ID    int    `json:"id"`
		Value int    `json:"value"`
		Label string `json:"label"`
	}

	var categoryResult []CategoryData
	id := 0
	for category, sum := range categorySums {
		if sum > 0 {
			categoryResult = append(categoryResult, CategoryData{
				ID:    id,
				Value: sum,
				Label: category,
			})
			id++
		}
	}

	// Возвращаем данные
	c.JSON(200, gin.H{
		"receipts":  receiptsData,
		"analytics": categoryResult,
	})
}

// @Summary Ежегодная аналитика расходов
// @Description Возвращает сумму расходов по категориям за текущий год для авторизованного пользователя и все чеки за год
// @Tags analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} map[string]string "Результат анализа за год"
// @Failure 401 {object} map[string]string "Пользователь не авторизован"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /analytics/year [get]
func (h *handler) analyticsYear(c *gin.Context) {
	// Достаём username из контекста
	usernameInterface, exists := c.Get("username")
	if !exists {
		c.JSON(401, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	username, ok := usernameInterface.(string)
	if !ok {
		c.JSON(500, gin.H{"error": "Ошибка сервера: неверный тип данных username"})
		return
	}

	// Получаем текущую дату и вычисляем дату начала года
	now := time.Now()
	startOfYear := time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, now.Location())
	endOfYear := startOfYear.AddDate(1, 0, 0) // следующий год

	// Ищем чеки пользователя за последний год, где данные JSON готовы
	var receipts []postgreSQL.Receipt
	if err := h.db.Where("username = ? AND date >= ? AND date < ? AND is_json_ready = true", username, startOfYear, endOfYear).
		Find(&receipts).Error; err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при получении чеков: " + err.Error()})
		return
	}

	// Если нет чеков, возвращаем пустой результат
	if len(receipts) == 0 {
		c.JSON(200, gin.H{"receipts": []any{}, "analytics": []any{}})
		return
	}

	// Массив чеков для ответа
	var receiptsData []datatypes.JSON
	for _, receipt := range receipts {
		receiptsData = append(receiptsData, receipt.JSONCheck)
	}

	// Собираем UUID всех чеков для дальнейшего поиска категорий
	var uuids []string
	for _, receipt := range receipts {
		uuids = append(uuids, receipt.UUID)
	}

	// Ищем категории для всех чеков
	var receiptCategories []postgreSQL.ReceiptCategories
	if err := h.db.Where("uuid IN ?", uuids).Find(&receiptCategories).Error; err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при получении категорий: " + err.Error()})
		return
	}

	// Суммируем значения по категориям
	categorySums := make(map[string]int)
	for _, rc := range receiptCategories {
		// Распаковываем JSON в map
		var categories map[string]int
		if err := json.Unmarshal(rc.JSONCat, &categories); err != nil {
			log.Printf("Ошибка парсинга категорий для чека %s: %v", rc.UUID, err)
			continue
		}

		// Суммируем значения
		for category, value := range categories {
			categorySums[category] += value
		}
	}

	// Формируем результат для категорий
	type CategoryData struct {
		ID    int    `json:"id"`
		Value int    `json:"value"`
		Label string `json:"label"`
	}

	var categoryResult []CategoryData
	id := 0
	for category, sum := range categorySums {
		if sum > 0 {
			categoryResult = append(categoryResult, CategoryData{
				ID:    id,
				Value: sum,
				Label: category,
			})
			id++
		}
	}

	// Возвращаем данные
	c.JSON(200, gin.H{
		"receipts":  receiptsData,
		"analytics": categoryResult,
	})
}

// AddTestReceiptsHandler — хендлер для добавления нескольких тестовых чеков в базу данных
// @Summary Добавить несколько тестовых чеков в базу данных
// @Description Этот хендлер создает несколько тестовых чеков и их категории для пользователя, который отправил запрос
// @Tags test
// @Accept json
// @Produce json
// @Success 200 "Чеки успешно добавлены"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /testdata/ [get]
func (h *handler) AddTestReceiptsHandler(c *gin.Context) {
	// Получаем имя пользователя из контекста
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось получить имя пользователя"})
		return
	}

	// Массив тестовых данных для чеков
	testJSONChecks := []string{
		`{
			"gt_parse": {
				"shop": "Монетка",
				"date": "25.04.2025",
				"time": ["22:04:00"],
				"items": [
					{"name": "Фильтр м", "measurement": "л", "count": 3, "price": 87.82, "overall": 263.46},
					{"name": "Электрич", "measurement": "шт", "count": 1, "price": 456.4, "overall": 456.4},
					{"name": "Принтер", "measurement": "шт", "count": 3, "price": 341.27, "overall": 1023.81},
					{"name": "Чай", "measurement": "шт", "count": 3, "price": 253.68, "overall": 761.04},
					{"name": "Хлеб", "measurement": "шт", "count": 2, "price": 440.71, "overall": 881.42}
				],
				"overall": "3386.13"
			}
		}`,
		`{
			"gt_parse": {
				"shop": "Пятёрочка",
				"date": "26.04.2025",
				"time": ["10:30:00"],
				"items": [
					{"name": "Молоко", "measurement": "л", "count": 2, "price": 60.50, "overall": 121.00},
					{"name": "Сыр", "measurement": "кг", "count": 1, "price": 300.00, "overall": 300.00},
					{"name": "Масло", "measurement": "г", "count": 500, "price": 125.00, "overall": 125.00},
					{"name": "Хлеб", "measurement": "шт", "count": 2, "price": 40.00, "overall": 80.00},
					{"name": "Яйца", "measurement": "шт", "count": 12, "price": 4.00, "overall": 48.00}
				],
				"overall": "674.00"
			}
		}`,
		`{
			"gt_parse": {
				"shop": "Лента",
				"date": "27.04.2025",
				"time": ["18:00:00"],
				"items": [
					{"name": "Туфли", "measurement": "шт", "count": 1, "price": 1200.00, "overall": 1200.00},
					{"name": "Куртка", "measurement": "шт", "count": 1, "price": 3500.00, "overall": 3500.00},
					{"name": "Рюкзак", "measurement": "шт", "count": 1, "price": 800.00, "overall": 800.00},
					{"name": "Шарф", "measurement": "шт", "count": 1, "price": 200.00, "overall": 200.00}
				],
				"overall": "5700.00"
			}
		}`,
		`{
			"gt_parse": {
				"shop": "Азбука Вкуса",
				"date": "28.04.2025",
				"time": ["12:15:00"],
				"items": [
					{"name": "Вино", "measurement": "бут", "count": 2, "price": 600.00, "overall": 1200.00},
					{"name": "Сёмга", "measurement": "кг", "count": 1, "price": 1500.00, "overall": 1500.00},
					{"name": "Икра", "measurement": "банка", "count": 1, "price": 950.00, "overall": 950.00},
					{"name": "Креветки", "measurement": "кг", "count": 2, "price": 700.00, "overall": 1400.00}
				],
				"overall": "5050.00"
			}
		}`,
	}

	for _, testJSONCheck := range testJSONChecks {
		// Парсим входной JSON
		var input AI.Input
		err := json.Unmarshal([]byte(testJSONCheck), &input)
		if err != nil {
			log.Printf("Ошибка парсинга JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка парсинга JSON"})
			return
		}

		// Генерируем уникальный UUID для чека
		receiptUUID := uuid.New().String()

		// Создаем новый чек для пользователя
		receipt := postgreSQL.Receipt{
			Username:    username.(string), // Используем username из контекста
			UUID:        receiptUUID,
			Name:        "Тестовый чек",
			Date:        time.Now(),
			JSONCheck:   datatypes.JSON([]byte(testJSONCheck)),
			IsJSONReady: true,
		}

		// Сохраняем чек в базе данных
		if err := h.db.Create(&receipt).Error; err != nil {
			log.Printf("Ошибка сохранения чека в БД: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения чека в БД"})
			return
		}

		// Разбираем категории товаров с помощью AI
		categories, err := AI.FindCategories(h.clientAI, testJSONCheck)
		if err != nil {
			log.Printf("Ошибка категоризации товаров: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка категоризации товаров"})
			return
		}

		// Сохраняем категории чека в базе данных
		categoriesMap := make(postgreSQL.JSONCat)
		for _, category := range categories {
			if category.Label != "" && category.Value > 0 {
				categoriesMap[category.Label] = category.Value
			}
		}

		// Конвертируем map в JSON
		categoriesJSON, err := json.Marshal(categoriesMap)
		if err != nil {
			log.Printf("Ошибка маршалинга категорий: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка обработки категорий",
			})
			return
		}

		receiptCategories := postgreSQL.ReceiptCategories{
			UUID:    receiptUUID,
			JSONCat: datatypes.JSON(categoriesJSON),
		}

		// Валидация перед сохранением
		if len(categoriesMap) == 0 {
			log.Printf("Предупреждение: попытка сохранить чек без категорий: %s", receiptUUID)
		}

		// Сохраняем категории в БД с обработкой ошибок
		if err := h.db.Create(&receiptCategories).Error; err != nil {
			log.Printf("Ошибка сохранения категорий чека %s в БД: %v", receiptUUID, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Ошибка сохранения категорий чека в БД",
				"details": err.Error(),
			})
			return
		}
	}

	// Отправляем успешный ответ
	c.JSON(http.StatusOK, gin.H{"message": "Чеки успешно добавлены"})
}
