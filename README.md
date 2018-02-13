# lift-emulator

Go app for emulate managing 1 lift from dispatcher.

It's very light application which has only business-logic for PoC.

There is no tests, cool logger, super env variables parsing, docker-schmoker, etc.. :)

## Task definition (from customer):

> 
> Предлагаем вам решить тестовое задание — написать программу «симулятор лифта». 
> 
> Программа запускается из командной строки, в качестве параметров задается:
> 
> - кол-во этажей в подъезде — N (от 5 до 20);
> - высота одного этажа;
> - скорость лифта при движении в метрах в секунду (ускорением пренебрегаем, считаем, что когда лифт едет — он сразу едет с определенной > > скоростью);
> - время между открытием и закрытием дверей.
> 
> После запуска программа должна постоянно ожидать ввода от пользователя и выводить действия лифта в реальном времени. События, которые нужно > выводить:
> 
> - лифт проезжает некоторый этаж;
> - лифт открыл двери;
> - лифт закрыл двери.
> 
> Возможный ввод пользователя:
> 
> - вызов лифта на этаж из подъезда;
> - нажать на кнопку этажа внутри лифта.
> 
> Считаем, что пользователь не может помешать лифту закрыть двери.
> 
> Все данные, которых не хватает в задаче, можно выбрать на свое усмотрение.
> 
> Решение можно прислать в виде ссылки на любой публичный git-репозиторий: GitHub, Bitbucket, GitLab и т.п.
> 
> Желаем успехов!

## How use it:

In one terminal run

```bash
go run main.go -d=2s -s=3 -c=10 -h=6
```

In other terminal run

```bash
curl -d "floor=10" -H "Content-Type: application/x-www-form-urlencoded" -X POST http://localhost:9090/external
curl -d "floor=7" -H "Content-Type: application/x-www-form-urlencoded" -X POST http://localhost:9090/internal
```
