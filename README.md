# mosbrew_scales

## Компиляция:

```
env GOOS=windows GOARCH=386 go build -v main.go
env GOOS=windows GOARCH=386 go build -ldflags -H=windowsgui -v main.go
```

Отладка через виртуальный com-порт:
```
alexey-r-laptop:~ alexey$ socat -d -d pty,raw,echo=0 pty,raw,echo=0
2019/03/10 23:09:48 socat[30669] N PTY is /dev/ttys000
2019/03/10 23:09:48 socat[30669] N PTY is /dev/ttys003
2019/03/10 23:09:48 socat[30669] N starting data transfer loop with FDs [5,5] and [7,7]

echo "TEST" > /dev/ttys003
echo $(( ( RANDOM % 100 )  + 1 )) > /dev/ttys003

cat < /dev/ttys000 - можно посмотреть вход
```

go run main.go

Примеры данных:
```
35 received len(6): 5   : 00110101 00000000 00000000 00000000 00001101 00001010  : ok!
40 received len(6): @  : 01000000 00000000 00000000 00100000 00001101 00001010  : not ok: strconv.Atoi: parsing "": invalid syntax
45 received len(6): E  : 01000101 00000000 00000000 00100000 00001101 00001010  : not ok: strconv.Atoi: parsing "": invalid syntax
50 received len(6): P   : 01010000 00000000 00000000 00000000 00001101 00001010  : not ok: strconv.Atoi: parsing "": invalid syntax
70 received len(6): u   : 01110101 00000000 00000000 00000000 00001101 00001010  : not ok: strconv.Atoi: parsing "": invalid syntax
80 received len(6): �  : 10000000 00000000 00000000 00100000 00001101 00001010  : not ok: strconv.Atoi: parsing "": invalid syntax
85 received len(6): �   : 10000101 00000000 00000000 00000000 00001101 00001010  : not ok: strconv.Atoi: parsing "": invalid syntax
95 received len(6): �   : 10010101 00000000 00000000 00000000 00001101 00001010  : not ok: strconv.Atoi: parsing "": invalid syntax
100 received len(6):    : 00000000 00000001 00000000 00000000 00001101 00001010  : not ok: strconv.Atoi: parsing "": invalid syntax
```
