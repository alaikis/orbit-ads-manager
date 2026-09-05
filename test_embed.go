package main

import (
	"fmt"
	"orbit/apps/api/internal/model"
)

func main() {
	m := model.TenantMember{
		TenantID: 1,
		UserID: 2,
	}
	fmt.Println(m.TenantID)
}