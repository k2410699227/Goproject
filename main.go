package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

func main() {

	var wg sync.WaitGroup //使用等待组对所有gproutines进行管理
	var root string       //读取的根目录
	var compare bool      //是否进行对比的标志
	var allPaths []string //储存可读取的文件路径

	goproject_test.main()
	flag.StringVar(&root, "d", ".", "搜索目录")
	flag.BoolVar(&compare, "c", false, "是否对串行和并行进行比较(true or false)")
	//dir := "E:/ASproject" //读取的根目录
	flag.Parse()
	flag.Usage()

	findDir(root, &allPaths) //递归遍历文件夹下有读取权限的所有文件

	start := time.Now()
	for _, file := range allPaths { //每个文件进行计算md5的值

		wg.Add(1)
		go read(file, &wg)
	}

	wg.Wait()
	parallel := time.Since(start).Seconds()

	if !compare {
		fmt.Println("并行运行时间为", parallel)
	} else {
		start = time.Now()
		for _, file := range allPaths { //每个文件进行计算md5的值

			sread(file)
		}

		serial := time.Since(start).Seconds()

		fmt.Println("并行运行时间为", parallel)
		fmt.Println("串行运行时间为", serial)
	}
}

func read(name string, wg *sync.WaitGroup) { //对文件进行读取的goroutine原型

	f, _ := os.ReadFile(name)

	fmt.Printf("%s\t%x\n", name, md5.Sum(f)) //输出文件路径及计算得到的md5值
	defer wg.Done()
}

func sread(name string) { //对文件进行读取的goroutine原型

	f, _ := os.ReadFile(name)

	fmt.Printf("%s\t%x\n", name, md5.Sum(f)) //输出文件路径及计算得到的md5值

}

func findDir(dir string, all *[]string) {
	fileinfo, err := ioutil.ReadDir(dir)
	if err != nil { //若文件夹没有打开权限，停止遍历该文件夹
		return
	}

	// 遍历这个文件夹
	for _, fi := range fileinfo {

		// 判断是不是目录
		if fi.IsDir() {
			findDir(dir+"/"+fi.Name(), all) //递归遍历文件夹
		} else {
			_, e := os.ReadFile(dir + "/" + fi.Name())
			if e == nil { //若有读取文件权限，添加文件路径
				*all = append((*all)[0:], dir+"/"+fi.Name())
			}
		}
	}
}

//未用到
func write(name string, content [16]byte) {
	//os.O_CREATE:创建
	//os.O_WRONLY:只写
	//os.O_APPEND:追加
	//os.O_RDONLY:只读
	//os.O_RDWR:读写
	//os.O_TRUNC:清空

	//0644:文件的权限
	//如果没有test.txt这个文件那么就创建，并且对这个文件只进行写和追加内容。
	file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("文件错误,错误为:%v\n", err)
		return
	}
	defer file.Close()

	file.Write([]byte("Hello World!")) //将str字符串的内容写到文件中，强制转换为byte，因为Write接收的是byte。
	file.WriteString("Hello")          //写字符串
}
