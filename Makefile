.PHONY: build frontend clean

frontend:
	pnpm --dir frontend install --frozen-lockfile
	pnpm --dir frontend build
	rm -rf internal/frontend/dist/*
	cp -R frontend/out/. internal/frontend/dist/

build: frontend
	mkdir -p bin
	go build -o bin/nilchan ./cmd

clean:
	rm -rf bin frontend/out internal/frontend/dist/*
