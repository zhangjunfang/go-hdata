package core

import (
	"database/sql"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/smartystreets/go-disruptor"
	"github.com/zhangjunfang/hdata/hdata_core/util"
)

const (
	BufferSize   int64 = (1 << 16) * 2048
	BufferMask   int64 = BufferSize - 1
	Reservations int64 = 1024
)

var ring [BufferSize]*Student

//var ring = [BufferSize]int64{}

type Student struct {
	Id   int64
	Name string
}

func (s *Student) String() string {
	return strconv.Itoa(int(s.Id)) + ":" + s.Name
}

var rows *sql.Rows
var stmt *sql.Stmt
var db *sql.DB
var ws sync.WaitGroup

func init() {
	dataSourceName := "root:@tcp(localhost:3306)/zhangboyu?charset=utf8"
	db, _ = util.GetConnection("mysql", dataSourceName, 10, 10)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//rows, _ = db.Query(" SELECT u.id ,u.`name` FROM `user` as u ")
	rows, _ = db.Query(" SELECT u.id ,u.`name` FROM `dd` as u ")
	stmt, _ = db.Prepare(" INSERT INTO cc(id ,name) VALUES(?,?) ")

}
func Execution() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//	tx, _ := db.Begin()
	//	stmt, _ := tx.Prepare(" INSERT INTO dd(id ,name) VALUES(?,?) ")
	//	for i := 0; i < 1000000; i++ {
	//		stmt.Exec(i, strconv.Itoa(i))
	//	}
	//	tx.Commit()
	//	stmt.Close()

	//	return

	controller := disruptor.Configure(BufferSize).WithConsumerGroup(SampleConsumer{}).
		BuildShared()
		//Build()

	controller.Start()

	started := time.Now()
	//publish(controller.Writer())
	//publishShared(controller.Writer())
	publishSharedBatch(controller.Writer())
	finished := time.Now()
	fmt.Println("耗时：", finished.Sub(started))
	//time.Sleep(1 * time.Hour)

	controller.Stop()

}

func publish(writer *disruptor.Writer) {
	sequence := int64(0)
	var id int64
	var name string
	for rows.Next() {
		sequence = writer.Reserve(1) //每次存储数据量
		//fmt.Println(sequence)
		rows.Scan(&id, &name)
		ring[sequence&BufferMask] = &Student{
			Id:   id,
			Name: name,
		}
		writer.Commit(sequence-1, sequence) //每生产一个 提交一个
	}
	//writer.Commit(0, sequence) //每生产一批 提交批次
}
func publishShared(writer *disruptor.SharedWriter) {
	sequence := int64(0)
	var id int64
	var name string
	for rows.Next() {
		sequence = writer.Reserve(1) ///每次存储数据量
		//fmt.Println(sequence)
		rows.Scan(&id, &name)
		ring[sequence&BufferMask] = &Student{
			Id:   id,
			Name: name,
		}
		writer.Commit(sequence-1, sequence) //每生产一个 提交一个
	}
	//writer.Commit(0, sequence) //每生产一批 提交批次
}
func publishSharedBatch(writer *disruptor.SharedWriter) {
	sequence := int64(0)
	var id int64
	var name string
	for rows.Next() {
		sequence = writer.Reserve(Reservations) ///每次存储数据量
		count := int64(1)
		for lower := sequence - Reservations + 1; lower <= sequence && rows.Next(); lower++ {
			rows.Scan(&id, &name)
			ring[lower&BufferMask] = &Student{
				Id:   id,
				Name: name,
			}
			count++
		}
		writer.Commit(sequence-count+1, sequence) //每生产一批 提交批次
	}
}

type SampleConsumer struct{}

func (this SampleConsumer) Consume(lower, upper int64) {

	started := time.Now()
	for lower <= upper {
		message := ring[lower&BufferMask]
		if message != nil {
			stmt.Exec(message.Id, message.Name)
		}
		lower++
	}
	fmt.Println("耗时：", time.Now().Sub(started))
	stmt.Close()
}

func (this SampleConsumer) Consume3(lower, upper int64) {
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare(" INSERT INTO dd(id ,name) VALUES(?,?) ")
	started := time.Now()
	for lower <= upper {
		message := ring[lower&BufferMask]
		if message != nil {
			stmt.Exec(message.Id, message.Name)
		}
		lower++
	}
	tx.Commit()
	finished := time.Now()
	fmt.Println("耗时：", finished.Sub(started))
	stmt.Close()
}

func (this SampleConsumer) Consume2(lower, upper int64) {
	started := time.Now()
	for lower <= upper {
		message := ring[lower&BufferMask]
		if message != nil {
			stmt.Exec(message.Id, message.Name)
		}
		lower++
	}
	finished := time.Now()
	fmt.Println("耗时：", finished.Sub(started))
}
