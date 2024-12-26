build:
	go build -o slg cmd/slg/main.go
install:
	go build -ldflags "-X main.assetsDir=/usr/local/share/slg/assets/en/ugly/" -o slg cmd/slg/main.go
	chmod 755 slg
	cp -f slg /usr/local/bin
	mkdir /usr/local/share/slg
	cp -r assets /usr/local/share/slg
uninstall:
	rm -f slg
	rm -f /usr/local/bin/slg
	rm -rf /usr/local/share/slg
