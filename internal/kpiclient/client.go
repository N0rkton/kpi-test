package kpiclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"newtest/internal/datamodel"
	"strings"

	log "github.com/sirupsen/logrus"
)

const token = "48ab34464a5573519725deb5865cc74c"

// KPIClinet - клиент для взаимодействия с api development.kpi-drive.ru
type KPIClinet struct {
	baseURL string
}

func NewKPIClient(url string) KPIClinet {
	return KPIClinet{baseURL: url}
}

// SendFact - Функция для отправки фактов
func (kpi KPIClinet) SendFact(fact datamodel.Fact) error {
	client := &http.Client{}

	// Формируем тело запроса
	data := url.Values{}
	data.Set("period_start", fact.PeriodStart)
	data.Set("period_end", fact.PeriodEnd)
	data.Set("period_key", fact.PeriodKey)
	data.Set("indicator_to_mo_id", fmt.Sprintf("%d", fact.IndicatorToMoID))
	data.Set("indicator_to_mo_fact_id", fmt.Sprintf("%d", fact.IndicatorToMoFactID))
	data.Set("value", fmt.Sprintf("%d", fact.Value))
	data.Set("fact_time", fact.FactTime)
	data.Set("is_plan", fmt.Sprintf("%d", fact.IsPlan))
	data.Set("auth_user_id", fmt.Sprintf("%d", fact.AuthUserID))
	data.Set("comment", fact.Comment)

	// Создаём HTTP-запрос
	req, err := http.NewRequest(http.MethodPost, kpi.baseURL+"facts/save_fact", strings.NewReader(data.Encode()))
	if err != nil {
		log.Warn("Worker: ошибка создания запроса: ", err)
		return err
	}

	// Устанавливаем заголовки
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "applicatino/x-www-form-urlencoded")

	// Выполняем запрос
	resp, err := client.Do(req)
	if err != nil {
		log.Warn("Worker: ошибка запроса: ", err)
		return err
	}

	// Закрываем тело ответа
	resp.Body.Close()

	// Лог успешной отправки
	if resp.StatusCode != http.StatusOK {
		log.Warn("Worker: ошибка API: ", resp.Status)
		b, _ := io.ReadAll(resp.Body)
		log.Warn(string(b))
		return errors.New("status code is not 200: " + resp.Status)
	}

	log.Info("Worker: факт успешно отправлен: ", fact)
	return nil

}

// CheckFactExists - Функция для проверки существования фактов
func (kpi KPIClinet) CheckFactExists(fact datamodel.Fact) error {
	client := &http.Client{}

	// Формируем тело запроса
	data := url.Values{}
	data.Set("period_start", fact.PeriodStart)
	data.Set("period_end", fact.PeriodEnd)
	data.Set("period_key", fact.PeriodKey)
	data.Set("indicator_to_mo_id", fmt.Sprintf("%d", fact.IndicatorToMoID))

	req, err := http.NewRequest(http.MethodPost, kpi.baseURL+"indicators/get_facts", strings.NewReader(data.Encode()))
	if err != nil {
		log.Warn("Workercheck: ошибка создания запроса: ", err)
		return err
	}

	// Устанавливаем заголовки
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Warn("Workercheck: ошибка запроса: ", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn("Workercheck: ошибка API: ", resp.Status)
		b, _ := io.ReadAll(resp.Body)
		log.Warn(string(b))
		return err
	}

	// Декодируем ответ
	var responseData []datamodel.Fact
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		log.Warn("Workercheck: ошибка декодирования ответа: ", err)
		return err
	}

	// Проверяем, есть ли уже такая запись
	for _, existingFact := range responseData {
		if existingFact.FactTime == fact.FactTime && existingFact.Value == fact.Value {
			log.Info("Workercheck: факт уже существует: ", fact)
			return nil
		}
	}

	log.Warn("fact not found")
	return errors.New("fact not found")

}
