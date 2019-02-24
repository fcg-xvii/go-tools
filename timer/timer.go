// Copyright 2017 F.C.G. (Flint's Crew Group). Russia, Moscow
// @author Mark Krass
/*
Таймер для использования в многопоточных приложениях. Обратите внимание!!! функция обратного вызова будет вызвана из другой
рутины, а не из той, где стартовал таймер. Это следует иметь в виду.
*/
/*
Рассмотрим пример, где мы запустим таймер с интервалом в секунду, а затем через 5 секунд отключим его:

  // Создаём объект таймера, передаём интервал и функцию обратного вызова
  var t *timer.Engine
  t = timer.New(time.Second, func() {
    // Выводим сообщение и текущее состояник таймера (будет STATE_STARTED)
    fmt.Println("callBack", t.State())
  })

  // Выведем состояние состояние таймера непосредственно перед запуском (== STATE_STOPPED)
  fmt.Println(t.State())

  // Стартуем его (Start() может быть вызван сколько угодно раз, это не повредит таймеру)
  t.Start()

  // Стартуем рутину, которая остановит таймер (можно было в принципе и в главной рутине, но я по политическим убеждениям
  // сделал в отдельной
  go func() {
    // Ждём 5 секунд
    time.Sleep(time.Second * 5)

    // Останавливаем (можно вызывать сколь угодно раз, если хочется)
    t.Stop()

    // Выведем состояние таймера сразу после остановки (== STATE_PREPARE_TO_STOP)
    fmt.Println(t.State())

    // Подождём, пока он не остановится окончательно
    for t.State() != timer.STATE_STOPPED {

      // Тут будет... Ну вы понели
      fmt.Println(t.State())
    }

    // Ещё раз на всякий случай... будет то же самое...
    fmt.Println(t.State())

    // Теперь, когда он остановился, можно запустить его ещё раз
    t.Start()
  }()

  // Спасибо за внимание :)
}

*/
package timer

import "time"

type TimerState int8 // Тип состояния таймера

const (
	STATE_STOPPED         TimerState = iota // Остановлен, готов к запуску
	STATE_STARTED                           // Запущен
	STATE_PREPARE_TO_STOP                   // Подготовка к остановке. Чтоб запустить, надо дождаться STATE_STOPPED
)

// Строковое отображение состояния
func (s TimerState) String() string {
	switch s {
	case STATE_STOPPED:
		return "STATE_STOPPED"
	case STATE_STARTED:
		return "STATE_STARTED"
	case STATE_PREPARE_TO_STOP:
		return "STATE_PREPARE_TO_STOP"
	default:
		return "STATE_UNDEFINED"
	}
}

// Конструктор таймера. Нужно передать интервал, функцию обратного вызова и флаг, указывающий на необходимость ожидания перед первым тиком
func New(interval time.Duration, callBack func(), runOnce bool) *Engine {
	return &Engine{STATE_STOPPED, interval, nil, callBack, runOnce}
}

// Структура таймера
type Engine struct {
	state    TimerState
	interval time.Duration
	ticker   chan int8
	callBack func()
	runOnce  bool
}

// Запуск таймера. Сработает, если State() == STATE_STOPPED
func (s *Engine) Start() {
	if s.state == STATE_PREPARE_TO_STOP {
		for s.state == STATE_PREPARE_TO_STOP {
			time.Sleep(time.Nanosecond)
		}
	}
	if s.state == STATE_STARTED {
		return
	}
	s.state = STATE_STARTED
	s.ticker = make(chan int8)
	go func() {
		if s.runOnce {
			s.callBack()
		}
	loop:
		for {
			select {
			case _, ok := <-s.ticker:
				if !ok {
					break loop
				}
			case <-time.After(s.interval):
				s.callBack()
			}
		}
		s.ticker = nil
		s.state = STATE_STOPPED
	}()
}

// Остановка (State() должен быть STATE_STARTED)
func (s *Engine) Stop() {
	if s.state != STATE_STARTED {
		return
	}
	s.state = STATE_PREPARE_TO_STOP
	close(s.ticker)
}

// Получаем состояние таймера
func (s *Engine) State() TimerState {
	return s.state
}

// Установка интервала таймера
func (s *Engine) SetInterval(interval time.Duration) {
	s.interval = interval
}
