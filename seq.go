package main

import "fmt"

var friends = map[string][]string{}
var Map_out = map[string][]string{}
var intersection_key [50]string
//var shuffle_out = map[string]string{}

func initate(person string,friends_num int) {
	var temp_of_friends []string
	var temp string
	for i:=0; i < friends_num ; i++{
		fmt.Printf("please insert %s's %d'st friends name:\n", person,i+1)
		fmt.Scan(&temp)
		temp_of_friends = append(temp_of_friends,temp)
		temp =""
	}
	friends[person] = temp_of_friends
}

func Map_shuffle_func() {
	var remember [50]string
	var index= 0
	var intersection bool = false
	var intersection_index int = 0
	for key := range friends {
		var orginal_key string
		var temp_key string
		temp_key = key
		for value := range friends[key] {
			var temp_value string
			var temp_bool bool
			temp_value = friends[key][value]
			temp_bool = order(temp_key, temp_value)
			if temp_bool == false {
				orginal_key = temp_key + temp_value
			} else {
				orginal_key = temp_value + temp_key
			}
			var conc []string
			for flag := range remember {
				if orginal_key == remember[flag] {
					intersection = true
					intersection_key[intersection_index]=orginal_key
					intersection_index += 1
				}
			}
			if intersection == false {
				for orginal_value := range friends[key] {
					conc = append(conc, friends[key][orginal_value])
				}
				Map_out[orginal_key] = conc
			} else if intersection == true {
				for value := range Map_out[orginal_key] {
					conc = append(conc, Map_out[orginal_key][value])
				}
				for orginal_value := range friends[key] {
					conc = append(conc, friends[key][orginal_value])
				}
				Map_out[orginal_key] = conc

			}
			remember[index] = orginal_key
			index++
			intersection = false
		}
	}
	for key:=range Map_out{
		print("key is:")
		print(key)
		print("     value is:")
		for value:= range Map_out[key]{
			print(Map_out[key][value])
		}
		print("\n")
	}
	print("intersect is:")
	for i:= range intersection_key{
		print(intersection_key[i]+"   ")
	}
}


/*func shuffle(){
	var remember [20]string
	var index = 0
	var cont bool = true
	var temp_string string
	for key := range Map_out{
		var temp_key string
		temp_key = key
		//fmt.Printf(temp_key)
		for flag := range remember{
			if temp_key == remember[flag]{
				cont = false
			}
		}
		if cont == true{
			for key2 := range Map_out{
				if key2 == temp_key{
					for value := range Map_out[key2]{
						temp_string += Map_out[key2][value]
					}
				}
			}
			remember[index] = temp_key
			index += 1
			shuffle_out[temp_key] = temp_string
		}
		cont =true
		temp_string =""
	}
}
*/

func order(one string,two string) bool{
	var status bool
	if one > two {
		status = true
	} else if one < two {
		status = false
	}
	return status
}

func reduce_in(){
	var send_to_reduce string
	for key := range intersection_key{
		send_to_reduce = ""
		var temp = intersection_key[key]
		for key2 := range Map_out[temp]{
			send_to_reduce += Map_out[temp][key2]
		}
		reduce(send_to_reduce)
	}
}


func reduce(reduce_in string){
	var flag bool = false
	var temp = reduce_in
	for i:=0;i<len(temp);i++{
		if i != len(temp){
			for j:=i+1;j<len(temp);j++{
				if string(temp[i]) == string(temp[j]){
					flag = true
				}
			}
		}
		if flag == true{
			print("common is:")
			print(string(temp[i])+"\n")
		}
		flag = false
	}
}



func main(){
	var number_of_persons int
	var number_of_friends int
	var person_name string
	fmt.Println("Please insert number of persons")
	fmt.Scan(&number_of_persons)
	for i:= 0; i < number_of_persons; i++{
		fmt.Printf("Please Enter person %d name:\n",i+1)
		fmt.Scan(&person_name)
		fmt.Printf("Please Enter num of person %d friends:\n",i+1)
		fmt.Scan(&number_of_friends)
		initate(person_name,number_of_friends)
	}
	Map_shuffle_func()
	reduce_in()
}

