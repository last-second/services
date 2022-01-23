get_user:
	go build -o get_user ./lambda/get_user/main.go

lambda: get_user
	echo "built lambdas"
