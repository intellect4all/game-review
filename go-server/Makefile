ec2 := ec2-user@ec2-18-235-12-50.compute-1.amazonaws.com:/home/ec2-user/servers/game_rev/
identity := ~/.ssh/game_rev.pem

run:
	go run main.go

buildServer:
	env GOARCH=amd64 GOOS=linux
	go build -o build/go-server

pushToEC2:
	scp $(identity) build/go-server $(ec2)
	scp $(identity) -r resources $(ec2)
	scp $(identity) .env $(ec2)
	scp $(identity) -r docs $(ec2)
	scp $(identity) -r private $(ec2)

buildDocker:
	docker build --platform linux/amd64 -t game_rev .
	docker tag game_rev intellect4all/game_rev

pushDocker:
	docker push intellect4all/game_rev

buildAndPushDocker: buildDocker pushDocker

