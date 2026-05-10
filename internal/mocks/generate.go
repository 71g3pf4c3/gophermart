package mocks

//go:generate mockgen -source=../controller/restapi/user.go -destination=restapi_user_mock.go -package=mocks User
//go:generate mockgen -source=../controller/restapi/orders.go -destination=restapi_order_mock.go -package=mocks Order
//go:generate mockgen -source=../controller/restapi/balance.go -destination=restapi_balance_mock.go -package=mocks Balance
//go:generate mockgen -source=../repo/repo.go -destination=repo_user_mock.go -package=mocks UserRepo
