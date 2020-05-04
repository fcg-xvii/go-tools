<p align="center">
  <span>Русский</span> |
  <a href="README.md#go-tools">English</a>  
</p>

# yask
> Инструмент для работы с сервисом синтеза и распознавания речи <b>Yandex Speech Kit</b> (подробнее о сервисе <a href="https://cloud.yandex.ru/docs/speechkit/" target="_blank">https://cloud.yandex.ru/docs/speechkit/</a>) для языка программирования <b>golang</b>. Инструмент позволяет синтезировать речь из тескта, а так же распознавать текст из звукового потока.

Перед началом работы необходимо зарегистрироваться на <a href="https://cloud.yandex.ru/" target="_blank">https://cloud.yandex.ru/</a> для получения API-ключа и идентификатора директирии (подробнее <a href="https://cloud.yandex.ru/docs" target="_blank">https://cloud.yandex.ru/docs</a>).

### Форматы аудиопотока
<ul>
    <li><b>OGG</b> <a href="https://ru.wikipedia.org/wiki/Ogg" target="_blank">https://ru.wikipedia.org/wiki/Ogg</a></li>
    <li><b>PCM</b> <a href="https://ru.wikipedia.org/wiki/Импульсно-кодовая_модуляция" target="_blank">https://ru.wikipedia.org/wiki/Импульсно-кодовая_модуляция</a> (при паспознавании текста в параметром формата lpcm, может быть использован поток формата wav</li>
</ul>

### Синтез речи
```golang
import (
	"log"
	"os"

	"github.com/fcg-xvii/go-tools/speech/yask"
)

func main() {
	yaFolderID := "b1g..."                              // идентификатор директории в yandex
	yaAPIKey := "AQVNy..."            // ключ api yandex
	text := "Привет, это тест синтеза речи с помощью сервиса Яндекса" // текст для синтеза

	// инициализация конфигурации для синтеза
	config := yask.TTSDefaultConfigText(yaFolderID, yaAPIKey, text)

	//
	r, err := yask.TextToSpeech(config)
	if err != nil {
		log.Println(err)
		return
	}

	f, err := os.OpenFile("tts.wav", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0655)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	if err := yask.EncodePCMToWav(r, f, config.Rate, 16, 1); err != nil {
		log.Println(err)
		return
	}
}
```




