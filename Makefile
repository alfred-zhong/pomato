VERSION=`cat ./version`

clean:
	rm -rf ./bin

mkdir-bin:
	mkdir -p ./bin
	
build: mkdir-bin
	go build -o ./bin/pomato -ldflags "-X github.com/alfred-zhong/pomato/cmd.version=$(VERSION)" github.com/alfred-zhong/pomato/cmd
	# cp ./pomato.yaml ./bin/

all: clean build

dev-server-run: build
	cd ./bin && ./pomato --pomodoro-time 5 --break-time 3 --long-break-time 4 --time-unit s
