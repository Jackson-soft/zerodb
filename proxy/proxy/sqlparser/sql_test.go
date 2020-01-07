package sqlparser

import (
	"bytes"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestConn_Insert(t *testing.T) {
	//c := newTestConn()
	//defer c.Close()
	//c.UseDB("backend_conn_test_database")

	// init db
	db, err := sql.Open("mysql", "zerodb:zerodb@tcp(zerodb-proxy003.pre.2dfire.info:9696)/zerodb?charset=utf8&allowOldPasswords=1")
	if err != nil {
		fmt.Println(err)
		return
	}
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	db.Ping()

	len := 100000

	var wg sync.WaitGroup
	wg.Add(len * 5)

	// 加个缓冲，限制下for循环的速度，不然goroutine一下子创建太多了
	chBuf := make(chan int, 400)
	for i := 1; i < len+1; i++ {
		s1 := fmt.Sprintf("insert into benchyou0 (id, k, c, pad) values(%d, %d, 'testxxxx', 'testxxxx')", i, i)
		s2 := fmt.Sprintf("insert into benchyou1 (id, k, c, pad) values(%d, %d, 'testxxxx', 'testxxxx')", i, i)
		s3 := fmt.Sprintf("insert into benchyou2 (id, k, c, pad) values(%d, %d, 'testxxxx', 'testxxxx')", i, i)
		s4 := fmt.Sprintf("insert into benchyou3 (id, k, c, pad) values(%d, %d, 'testxxxx', 'testxxxx')", i, i)
		s5 := fmt.Sprintf("insert into benchyou4 (id, k, c, pad) values(%d, %d, 'testxxxx', 'testxxxx')", i, i)

		go exe(s1, db, t, &wg, chBuf)
		go exe(s2, db, t, &wg, chBuf)
		go exe(s3, db, t, &wg, chBuf)
		go exe(s4, db, t, &wg, chBuf)
		go exe(s5, db, t, &wg, chBuf)
	}
	println("waiting ###################################")
	wg.Wait()
	println("done ###################################")

}

func exe(sqlString string, db *sql.DB, t *testing.T, wg *sync.WaitGroup, ch chan int) {
	defer func() {
		wg.Done()
		<-ch
	}()

	// 如果缓冲满了，会阻塞
	ch <- 1
	if _, err := db.Exec(sqlString); err != nil {
		t.Fatal(err)
	}

}

func testParse(t *testing.T, sql string) {
	_, err := Parse(sql)
	if err != nil {
		t.Fatal(err)
	}

}

func TestStringToByteArr(t *testing.T) {
	str2 := "select * from order.instance_detail where entity_id = 3123123123123213213"
	fmt.Println(len(str2))
	data2 := []byte(str2)
	fmt.Println(len(data2))
	fmt.Println(data2)
}

func TestMake(t *testing.T) {
	data := make([]byte, 3, 6)
	data[0] = 0
	data[1] = 1
	data[2] = 2
	println(len(data))
	data = append(data, 3)
	println(len(data))
	data = append(data, 4)
	println(len(data))
	data = append(data, 5)
	println(len(data))

	data = append(data, 6)
	println(len(data))
}

func BenchmarkParse(b *testing.B) {
	for n := 0; n < b.N; n++ {
		//_, err := Parse("select 'abcd', 20, 30.0, eid from a where 1=eid and name='3'")
		_, err := Parse("SELECT * FROM benchyou4 WHERE id = 882095016067628781 and a = 111")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBuilder(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var sb strings.Builder
		for i := 0; i < 1000; i++ {
			sb.WriteString(strconv.Itoa(i))
		}

	}
}

func BenchmarkBuffer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var b bytes.Buffer
		for i := 0; i < 1000; i++ {
			b.WriteString(strconv.Itoa(i))
		}
	}
}

func BenchmarkPlus(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var s string
		for i := 0; i < 1000; i++ {
			s += strconv.Itoa(i)
		}
	}
}

func TestSet(t *testing.T) {
	sql := "set names gbk"
	testParse(t, sql)
}

func TestSimpleSelect(t *testing.T) {
	sql := "select last_insert_id() as a"
	testParse(t, sql)
}

func TestMixer(t *testing.T) {
	sql := `admin upnode("node1", "master", "127.0.0.1")`
	testParse(t, sql)

	sql = "show databases"
	testParse(t, sql)

	sql = "show tables from abc"
	testParse(t, sql)

	sql = "show tables from abc like a"
	testParse(t, sql)

	sql = "show tables from abc where a = 1"
	testParse(t, sql)

	sql = "show proxy abc"
	testParse(t, sql)
}
