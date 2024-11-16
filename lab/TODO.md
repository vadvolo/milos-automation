# TODO

## netbox

- [x] netbox: разобраться с корректным импортом данных

## annet

- [x] annet: поправить конфиг, чтобы annet работал
- [x] annet: перенести токен в конфиг (есть в `annet/config.yaml`, но всё ещё раскидан по compose'ам)
- [ ] annet: переписать все комментарии на русском языке на английский язык
- [ ] annet: https://st.yandex-team.ru/NOCDEV-16217 При deploy не показывает наличие оставшегося diff
- [ ] annet: https://st.yandex-team.ru/NOCDEV-16218 deploy не показывает progress

## scripts

- [x] netsshsetup: спрятать DEBUG
- [x] netsshsetup: выпилить `vendor` из репо
- [ ] netsshsetup: распараллелить
- [x] netsshsetup: сделать сборку под целевую архитектуру

## emulators

- [ ] dynamips: добавить healthcheck на готовность ВМ

## Makefile

- [x] Makefile: проверять что в папке ./vm_images есть образ
- [ ] Makefile: убрать SUDO

## lab00

- [x] lab00: смержить генераторы в README
- [x] lab00: поменять адреса устройств с 0, 1, 2 на 1, 2, 3

## lab01

- [ ] ...

## lab10

- [x] Поменять адресацию на линкнетах
- [x] Перенести redistribute в mesh

## lab12

- [x] Добавить поддержку Arista в генераторы
- [x] Добавить поддержку FRR в генераторы

## general

- [x] Репо: в README заменить python3 -m ... на annet
- [ ] Репо: сделать в main чистовую версию кода без лишних обростков
- [ ] Репо: проверить все креды, заменить на annet/annet
  - [x] lab00
  - [x] lab01
  - [ ] lab02
  - [ ] lab03
- [x] Репо: переписать корневой README
- [ ] Репо: в README лаб Cisco и Arista добавить "подождать ~45 сек после старта"
- [x] Репо: попробовать разархивировать образ Cisco, возможно будет быстрее собираться
- [ ] Репо: вычистить мусор из .gitignore (после переезда в main)
- [ ] Репо: в README добавить `make clean`, который приведёт её к изначальному состоянию
- [ ] Репо: в README добавить скриншоты
