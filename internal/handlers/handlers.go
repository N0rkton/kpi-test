package handlers

import (
	"encoding/json"
	"net/http"
	"newtest/internal/datamodel"
	"newtest/internal/kpiclient"
)

// количество одновременно работающих воркеров
const workers = 10

// handlerWrapper - структура, содержащая канал для задач (jobs) и клиент для работы с API
type handlerWrapper struct {
	jobs   chan datamodel.Fact
	client kpiclient.KPIClinet
}

// Init - функция инициализации handlerWrapper.
func Init(client kpiclient.KPIClinet) handlerWrapper {
	jobs := make(chan datamodel.Fact, 1000)

	hw := handlerWrapper{jobs: jobs, client: client}

	//Создаем воркеры для работы с api
	for w := 0; w < workers; w++ {
		go hw.worker()
	}

	return hw
}

// NewFacts - принимает факты в формате json в теле запроса
func (hw handlerWrapper) NewFacts(res http.ResponseWriter, req *http.Request) {
	var body []datamodel.Fact

	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := 0; i < len(body); i++ {
		hw.jobs <- body[i]
	}
}

// worker - каналы jobs может хранить в буфере до 1000 фактов, в продуктовой среде его можно заменить кафкой.
// При получении нового запроса отправлять его информацию в кафку далее брать запись из кафки и работать с ней
func (hw handlerWrapper) worker() {
	for fact := range hw.jobs {
		if err := hw.client.SendFact(fact); err != nil {
			continue
		}

		hw.client.CheckFactExists(fact)

	}
}
