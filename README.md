Каждая ошибка имеет тип (NotFound, AccessDenied и т.д.) этот тип переконвертится в HTTP код или GRPC код и уйдет в ответ на запрос ручки.
Код для сохранения ошибки в базе с целью отдать потом пользователю историю мы использовать не будем.
Можно управлять потоком с помощью tag в ошибке.
Можно обогащать ошибку любыми данными при движении вверх по коду. Реализовано двумя ручками - WithPayloadKV - данные для пользователя и WithLogKV - данные для логирования.
Есть конвертеры которые будут создавать наши ошибки из общеизвестных либ - gopg, runtime и т.д. 
Можно создавать новую ошибку с новым кодом и сообщением сохраняя старую в cause.
Есть возможность при создании типа ошибки задавать коллбэк в который будет передан контекст и ошибка, таким образом можно навешивать автологирование ошибок.
Сообщение ошибки всегда уходит пользователю вместе с кодом.
