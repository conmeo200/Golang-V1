package main

import (
    "math/rand"
)




func main(){
	// PlaceOrder(Order{"t4uQ8@example.com", "0123456789", "123"})

	// time.Sleep(3 * time.Second)

    nums := make([]int, 1000)
    for i    := range nums {
        nums[i] = rand.Intn(1000) + 1
    }

	// ProcessImage(nums)

}

