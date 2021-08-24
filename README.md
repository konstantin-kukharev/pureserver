# Application server
Сервер приложений (англ. application server) — это программная платформа (фреймворк), предназначенная для эффективного исполнения процедур (программ, скриптов), на которых построены приложения. Сервер приложений действует как набор компонентов, доступных разработчику программного обеспечения через API (интерфейс прикладного программирования), определённый самой платформой.

Для веб-приложений основная задача компонентов сервера — обеспечивать создание динамических страниц. Однако современные серверы приложений включают в себя и поддержку кластеризации, повышенную отказоустойчивость, балансировку нагрузки, позволяя таким образом разработчикам сфокусироваться только на реализации бизнес-логики.

## Описание

Реализация tcp/http сервера приложений, для приложений в единой ИС.

Предполагается что на более высоком уровне расположен web сервер или группа серверов, в зону ответственности которых входит: 
- обработка и первичная фильтрация внешнего трафика маршрутизация до целевого application server
- безопасность (ssl, tls, ...)
- маршрутизацию трафика между приложениями (app server)

Основные цели:
- упростить реализацию корпоративных приложений, позволив разработчику сосредоточиться на реализации бизнес логики
- отвечать за прием и маршрутизацию сообщений внутри приложения
- предоставить набор интерфейсов для возможности реализации корпоративных приложений без привязки к конкретной библиотеке или реализации http/app server
- предоставить единый интерфейс для взаимодействия с системой оркестрации контейнеров: 
  - набор handler`ов с информацией о текущем статусе приложения
  - возможность отправки push сообщений в случае инцидентов (большая нагрузка, ...)
  - набор интерфейсов для трассировки, без привязки к конкретному инструменту или библиотеке
  - набор инструментов для логирования
- реализация graceful shutdown на уровне пода в системе оркестрации

## Инструменты и библиотеки

За основу были взяты:
- [pat (formerly pat.go) - A Sinatra style pattern muxer for Go's net/http library](https://pkg.go.dev/github.com/bmizerany/pat)
- [evio (event loop networking framework)](https://github.com/tidwall/evio)

## Hello World
```golang
package main

import (
  "encoding/json"
  "fmt"
  ps "github.com/konstantin-kukharev/pureserver"
  "strconv"
)

const DefaultPort = 8080

type TestResponse struct {
  A int `json:"A"`
  B int `json:"B"`
  C int `json:"C"`
}

func HelloWorldHandler(w ps.ResponseWriter, req ps.HttpRequestInterface) {
  var body TestResponse
  _ = json.Unmarshal([]byte(req.GetBody()), &body)
  name, _ := req.GetParam("increment")
  app, _ := strconv.Atoi(name)
  body.A += app
  body.B *= app
  body.C = body.C / app
  result, _ := json.Marshal(body)
  w.SetBody(result)
}

func main() {
  mux := ps.NewMux()
  mux.Post("/hello/:increment", ps.HandlerFunc(HelloWorldHandler))

  server := ps.NewHttp(mux)
  server.SetPort(DefaultPort)
  fmt.Println(server.Serve())
}
```
## Test
Для тестов использовалась библиотека [vegeta](https://github.com/tsenart/vegeta)
![alt text](doc/img.png)
cmd/ps
```
Requests      [total, rate, throughput]  3000, 50.02, 50.02
Duration      [total, attack, wait]      59.979318458s, 59.978972458s, 346µs
Latencies     [mean, 50, 95, 99, max]    375.157µs, 356.943µs, 491.728µs, 667.772µs, 7.478125ms
Bytes In      [total, mean]              102000, 34.00
Bytes Out     [total, mean]              102000, 34.00
Success       [ratio]                    100.00%
Status Codes  [code:count]               200:3000  
Error Set:

Bucket           #     %       Histogram
[0s,     200µs]  16    0.53%   
[200µs,  300µs]  308   10.27%  #######
[300µs,  400µs]  1979  65.97%  #################################################
[400µs,  500µs]  558   18.60%  #############
[500µs,  10ms]   139   4.63%   ###
[10ms,   100ms]  0     0.00%   
[100ms,  +Inf]   0     0.00%   
```
cmd/http
```
Requests      [total, rate, throughput]  3000, 50.02, 0.00
Duration      [total, attack, wait]      59.979894333s, 59.978111792s, 1.782541ms
Latencies     [mean, 50, 95, 99, max]    2.259091ms, 1.763455ms, 2.400878ms, 14.803734ms, 127.604458ms
Bytes In      [total, mean]              0, 0.00
Bytes Out     [total, mean]              0, 0.00
Success       [ratio]                    0.00%
Status Codes  [code:count]               0:3000  
Error Set:

Bucket           #     %       Histogram
[0s,     200µs]  0     0.00%   
[200µs,  300µs]  0     0.00%   
[300µs,  400µs]  0     0.00%   
[400µs,  500µs]  0     0.00%   
[500µs,  10ms]   2967  98.90%  ##########################################################################
[10ms,   100ms]  30    1.00%   
[100ms,  +Inf]   3     0.10%   
```
## FYI
Использование данных библиотек на момент выбора было продиктовано простотой реализации и результатами сравнительных метрик производительности

В дальнейшем предполагается возможность замены или добавления других инструментов, в зависимости от результатов эксплуатации, поэтому при интеграции модуля следует использовать интерфейсы и избегать сильной связанности 

В данный момент наибольшие сомнения вызывает наивная реализация парсера заголовков и тела запроса, несмотря на полное покрытие библиотеки тестами.

для корректного использования private репозитория github необходимо в $HOMEDIR/.gitconf добавить блок
```
[url "git@github.com:"]
    insteadOf = https://github.com/
```
так же необходимо добавить GOPRIVATE в go env
```
GOPRIVATE="github.com/konstantin-kukharev/pureserver"
```
