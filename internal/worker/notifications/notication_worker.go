package notifications

import (
	"fmt"
	"sync"
	"time"
)

type NotificationWorker struct {
	
}

func NewNotificationWorker() *NotificationWorker {
	return &NotificationWorker{}
}

type Order struct {
	UserEmail string
	Phone string
	ID  string
}

func PlaceOrder(order Order) {
    // Lưu đơn hàng vào DB (bắt buộc xong trước)
    saveOrder(order)

    // Gửi thông báo — không cần chờ, dùng goroutine
    go sendEmail(order.UserEmail, "Đặt hàng thành công!")
    go sendSMS(order.Phone, "Đơn hàng #"+order.ID)
    go updateAnalytics(order)

    // Trả response về cho user ngay, không chờ email/SMS
    fmt.Printf("Đặt hàng thành công! \n")
}

func sendEmail(email string, message string) {
	time.Sleep(2 * time.Second)
	fmt.Printf("Send Emai %v thành công!, message: %v \n", email, message)
}

func sendSMS(phone string, message string) {
	time.Sleep(2 * time.Second)
	fmt.Printf("Send Phone %v thành công!, message: %v \n", phone, message)
}

func updateAnalytics(order Order) {
	time.Sleep(1 * time.Second)
	fmt.Printf("Update Analytics: %v\n", order)
}

func saveOrder	(order Order) {
	time.Sleep(3 * time.Second)
	fmt.Println("Save Order:", order.ID)
}


func ProcessImage(array_int []int) {

    jobs := make(chan int, len(array_int))
    
      // 3 worker chạy song song
    var wg sync.WaitGroup
    for i := 0; i < 3; i++ {
        wg.Add(1)
        workerID := i 
        go func() {
            defer wg.Done()
            for size := range jobs {
				fmt.Printf("Worker %v handler messages %v\n", workerID, size)
            }
        }()
    }

    for _, s := range array_int {
        jobs <- s
    }

    close(jobs)
    wg.Wait()
    fmt.Println("Xử lý ảnh xong!")
}

func resizeAndUpload(size string) {
	time.Sleep(1 * time.Second)
	fmt.Printf("Resize and upload %v \n", size)
}