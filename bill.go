package main

import "fmt"

type Bill struct{
	name string
	items map[string]float64
	tip float64
}

func createBill(name string) Bill {
	data := Bill{
		name : name,
		items: map[string]float64{},
		tip  : 0,
	}

	return data
}

//Bill format

func (b *Bill) format() string {
	fs := "Order name : " + b.name + "\n"
	fs += "Bill breakdown : \n"
	var total float64 = 0

	//list items
	for k, v := range b.items {
		fs += fmt.Sprintf("%-25v ... $%v \n", k + ":", v)
		total += v
	}

	//total 
	fs += fmt.Sprintf("%-25v ... $%v\n", "tip : ", b.tip)


	fs += fmt.Sprintf("---------------------------------------------\n")

	//total 
	fs += fmt.Sprintf("%-25v $%0.2f", "total : ", total + b.tip)

	fs += fmt.Sprintf("\n**********************************************\n")

	return fs
}

func (b *Bill) updateTip(x float64) {
	b.tip = x
}

func (b *Bill) addItem(name string, price float64) {
	b.items[name] = price
}