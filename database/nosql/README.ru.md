<p align="center">
  <span>Русский</span> |
  <a href="README.md">English</a>
</p>

# nosql
> Инструмент позволяет легко вызвать функцию SQL базы данных, принимающую в качестве входного параметра json объект.
> Результирующим значением функции так же должен являться json объект. Всю работу формированию запроса, 
> конвертацию параметров и работу с транзакциями инструмент берёт на себя, позволяя тем самым при минимальном
> количестве кода реализовать полноценное REST API.

### Базы данных
В качестве базы данных может быть использована любая SQL база данных. Тестирование производилось на 
<a href="https://postgrespro.ru">PostgreSQL</a>

### Как использовать
> В качестве примера будет выполнено подключение к базе данных PostgreSQL и вызвана функция подсчёта количества
> элементов в массиве, переданном в качестве аргумента

```golang
  // В качестве базы данных используем PostgreSQL c расширение plv8 (https://plv8.github.io/)
  // В базе данных необходимо определить функцию, которая будет подсчитывает количество элементов в массиве
  // Аргументом функции является json объект { "input": array }, функция вернёт json объект 
  // { "input": array, output: count }
	//
	// CREATE OR REPLACE FUNCTION public.arr_count(data jsonb)
	// RETURNS jsonb AS
	// $BODY$
	//   if(typeof(data.input) != 'object' || !data.input.length) {
	//	   plv8.elog(ERROR, 'Incoming data must be array')
	//   }
	//   data.output = data.input.length
	//   return data
	// $BODY$
	// LANGUAGE plv8 IMMUTABLE STRICT

	// dbConn read from config file before
	// example dbConn string: postgres://postgres:postgres@127.0.0.1/postgres?sslmode=disable&port=5432

  // Строка с параметрами подключения к базе данных
  dbConn := postgres://postgres:postgres@127.0.0.1/postgres?sslmode=disable&port=5432

	// подкоючение к базе данных
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		fmt.Println(err)
		return
	}

	// определение функции открытия транцакции.
	openTX := func() (*sql.Tx, error) {
		return db.Begin()
	}

	// создание объекта для вызова функций базы данных
	api := nosql.New(openTX)

	// инициализация параметров запроса к базе данных
	data := map[string]interface{}{
		"input": []int{1, 2, 3},
	}

	// вызов функции базы данных с параметром
	result, err := api.Call("public.arr_count", data)
	if err == nil {
		fmt.Println(err)
		return
	}

	// запрос выполнен.
	fmt.Println(result)
```

## Лицензия

The MIT License (MIT), подробнее [LICENSE](LICENSE).
