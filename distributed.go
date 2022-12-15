package main

import (
	"fmt"
	"sync"
	"time"
)

var test int = 0
var shuffle_out = map[string][]string{}
var intersection_key [30]string

type Job struct{
	id int
    friends map[string][]string
}

type Result struct {
	job Job
	output map[string][]string
}

var jobs = make(chan Job, 5)
var results = make(chan Result, 5)

func local_Map_Shuffle (input map[string][]string) map[string][]string{
	var Map_out = map[string][]string{}
	var intersection_key [50]string
	var remember [50]string
	var index= 0
	var intersection bool = false
	var intersection_index int = 0
//////////////////////////////////////////////////////////////////////////////
// input check
//////////////////////////////////////////////////////////////////////////////
	for key := range input {
		var orginal_key string
		var temp_key string
		temp_key = key
		for value := range input[key] {
			var temp_value string
			var temp_bool bool
			temp_value = input[key][value]
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
				for orginal_value := range input[key] {
					conc = append(conc, input[key][orginal_value])
				}
				Map_out[orginal_key] = conc
			} else if intersection == true {
				for value := range Map_out[orginal_key] {
					conc = append(conc, Map_out[orginal_key][value])
				}
				for orginal_value := range input[key] {
					conc = append(conc, input[key][orginal_value])
				}
				Map_out[orginal_key] = conc

			}
			remember[index] = orginal_key
			index++
			intersection = false
		}
	}
	for i:=range intersection_key{
		var local_common string
		var send_to_local_reducer string
		var local_intersection_key string
		var local_reducer_output []string
		local_intersection_key = intersection_key[i]
		for j:=range Map_out[local_intersection_key]{
			send_to_local_reducer += Map_out[local_intersection_key][j]
		}
		local_reducer_output,local_common = local_reducer(send_to_local_reducer)
		if local_common != ""{
			fmt.Printf("local common friends between %s is:%s",intersection_key[i],local_common)
		}
        Map_out[local_intersection_key] = local_reducer_output

	}
	time.Sleep(2 * time.Second)
	return Map_out
}


func local_reducer(input string) ([]string,string){
	var local_common string
	var local_reducer_output []string
	var flag bool = false
	var temp = input
	for i:=0;i<len(temp);i++{
		if i != len(temp){
			for j:=i+1;j<len(temp);j++{
				if string(temp[i]) == string(temp[j]){
					flag = true
				}
			}
		}
		if flag == false{
			local_reducer_output = append(local_reducer_output,string(temp[i]))
		}else if flag == true{
			local_common += string(temp[i])
		}
		flag = false
	}
	return local_reducer_output,local_common
}


func order(one string,two string) bool{
	var status bool
	if one > two {
		status = true
	} else if one < two {
		status = false
	}
	return status
}

func worker (wg *sync.WaitGroup)  {
	for job := range jobs{
		output := Result{job,local_Map_Shuffle(job.friends)}
		results <- output
	}
	wg.Done()
}


func CreateWorkerPool(noOfWorkers int){
	var wg sync.WaitGroup
	for i:=0; i < noOfWorkers; i++ {
		wg.Add(1)
		go worker(&wg)
	}
	wg.Wait()
	close(results)
}

func allocate(noOfJobs int){
	for i:= 0;i < noOfJobs; i++{
		var temp_of_init = map[string][]string{}
		var number_of_persons int
		var number_of_friends int
		var person_name string
		var str []string
		fmt.Printf("Please insert number of persons in job %d:\n",i)
		fmt.Scan(&number_of_persons)
		for j:= 0; j < number_of_persons; j++{
			fmt.Printf("Please Enter person %d name in job %d:\n",j+1,i)
			fmt.Scan(&person_name)
			fmt.Printf("Please Enter num of person %d friends in job %d:\n",j+1,i)
			fmt.Scan(&number_of_friends)
			str = initate(person_name,number_of_friends,i)
			temp_of_init[person_name] = str
		}
		job := Job{i , temp_of_init}
		jobs <- job
	}
	close(jobs)
}


func initate(person string,friends_num int , id int) []string {
	var temp_of_friends []string
	var temp string
	for i:=0; i < friends_num ; i++{
		fmt.Printf("please insert %s's %d'st friends name in job %d:\n", person,i+1,id)
		fmt.Scan(&temp)
		temp_of_friends = append(temp_of_friends,temp)
		temp =""
	}
	return temp_of_friends
}



func shuffle_output(done chan bool){
	var remember [30]string
	var index int = 0
	var intersection_index int = 0
	for result := range results{
		for key:=range result.output{
			var intersection_flag bool = false
			if key != ""{
				for i := range remember{
					if key == remember[i]{
						intersection_flag = true
					}
				}
				var temp_str []string
				if intersection_flag == false{
					for value:= range result.output[key]{
						temp_str = append(temp_str,result.output[key][value])
					}
					shuffle_out[key] = temp_str

				}else if intersection_flag == true{
					for value1 := range result.output[key]{
						temp_str = append(temp_str,result.output[key][value1])
					}
					for value2 := range shuffle_out[key]{
						temp_str = append(temp_str,shuffle_out[key][value2])
					}
					shuffle_out[key] = temp_str
					intersection_key[intersection_index]=key
					intersection_index++
				}
				remember[index] = key
				index ++
				intersection_flag = false
			}

		}
	}
	done <- true
}

func reduce_in(){
	var send_to_reduce string
	for key := range intersection_key{
		send_to_reduce = ""
		var temp = intersection_key[key]
		for key2 := range shuffle_out[temp]{
			send_to_reduce += shuffle_out[temp][key2]
		}
		reduce(send_to_reduce,temp)
	}
}


func reduce(reduce_in string,key string){
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
			fmt.Printf("common friends between %s is:%s\n",key,string(temp[i]))

		}
		flag = false
	}
}


func main(){
	var NoOfJobs int
	var NoOfWorkers int
	fmt.Printf("please enter num of jobs:")
	fmt.Scan(&NoOfJobs)
	fmt.Printf("Please Enter num of workers:")
	fmt.Scan(&NoOfWorkers)
	go allocate(NoOfJobs)
	done := make(chan bool)
	go shuffle_output(done)
	CreateWorkerPool(NoOfWorkers)
	<-done
    reduce_in()
}