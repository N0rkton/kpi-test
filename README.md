Тестовое задание для backend Golang

Сервер запускается на порту 8080
Доступен хендлер "/upload-facts", метод Post, ожидаемое тело запроса:

```[

  {
  
    "period_start": "2024-12-01",
    
    "period_end": "2024-12-31",
    
    "period_key": "month",
    
    "indicator_to_mo_id": 227373,
    
    "indicator_to_mo_fact_id": 0,
    
    "value": 3,
    
    "fact_time": "2024-12-31",
    
    "is_plan": 0,
    
    "auth_user_id": 40,
    
    "comment": "buffer Shvarts git: N0rkton"
    
  }
  
]```



Для реализации буфера был выбран паттерн worker pool. Он позволяет распаралелить обработку фактов что повышает производительнось. 
Для продуктового окружения можно заменить его на кафку. 
