package role

//go:generate go run ../../../../cmd/injecttrace -file admin.go -receiver Role -service Role

const Admin = "admin"
