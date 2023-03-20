ec2 := ec2-user@ec2-18-235-12-50.compute-1.amazonaws.com:/home/ec2-user/servers/game_rev/

run:
	go run main.go

buildServer:
	env GOARCH=amd64 GOOS=linux
	go build -o build/go-server

pushToEC2:
	scp build/go-server $(ec2)
	scp -r resources $(ec2)
	scp .env $(ec2)
	scp -r docs $(ec2)
	scp -r private $(ec2)

