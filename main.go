package main

import (
	"receipt-processor/domain"
	"receipt-processor/service"

	"github.com/gin-gonic/gin"
)

// set up http layer of application to call service to handle the business logic
func main() {
	service := service.ReceiptService{ReceiptsById: map[string]domain.Receipt{}}
	router := gin.Default()

	router.POST("/receipts/process", func(c *gin.Context) {
		var receipt domain.Receipt
		if error := c.ShouldBindJSON(&receipt); error != nil {
			c.JSON(400, gin.H{
				"error": error.Error(),
			})
			return
		}
		uuid := service.ProcessReceipt(receipt)
		c.JSON(200, gin.H{"id": uuid})
	})

	router.GET("/receipts/:id/points", func(c *gin.Context) {
		points, ok := service.GetPoints(c.Param("id"))
		if ok {
			c.JSON(200, gin.H{"points": points})
		} else {
			c.JSON(404, gin.H{"error": "not found"})
		}
	})

	router.Run(":8080")
}
