build:
	docker build -t resume .

run:
	docker run -p 54321:49999 resume 