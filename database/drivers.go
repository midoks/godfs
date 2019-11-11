package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type DB struct {
	db *sql.DB
}

type BinFileGroup struct {
	Id      int64
	Md5     string
	Created string
}

type BinFile struct {
	Id      int64
	Md5     string
	NodeNum int
	Node    string
	Path    string
	Size    int64
	Created string
}

func (this *DB) Init(path string) {
	this.db, _ = sql.Open("sqlite3", path)

	sql_table_bin_file_group := `
    CREATE TABLE IF NOT EXISTS bin_file_group (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        md5 CHAR(32) NULL UNIQUE,
        created DATE NULL
    );
`
	this.db.Exec(sql_table_bin_file_group)
	this.db.Exec("create index bin_file_group_md5 on bin_file_group(md5);")

	sql_table_bin_file := `
    CREATE TABLE IF NOT EXISTS bin_file (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        gid INTEGER NULL,
        md5 CHAR(32) NULL UNIQUE,
        path TEXT NULL,
        node_num INT NULL,
        node TEXT NULL,
        size INTEGER NULL,
        updated DATE NULL,
        created DATE NULL
    );
`
	this.db.Exec(sql_table_bin_file)
	this.db.Exec("create index bin_file_md5 on bin_file(md5);")
}

func (this *DB) FindFileGroupGetId(md5 string) int64 {
	data, isFind := this.FindFileGroup(md5)
	if isFind {
		return data.Id
	}

	insertId, err := this.AddFileGroup(md5)
	if err == nil {
		return insertId
	}
	return 0
}

func (this *DB) FindFileGroup(md5 string) (*BinFileGroup, bool) {
	var (
		err  error
		rows *sql.Rows
	)
	di := &BinFileGroup{}

	rows, err = this.db.Query(fmt.Sprintf("select id,md5,created from bin_file_group where md5='%s' limit 1", md5))
	if err != nil {
		return di, false
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&di.Id, &di.Md5, &di.Created)
		if err != nil {
			return di, false
		}
		return di, true
	}

	return di, false
}

func (this *DB) AddFileGroup(md5 string) (int64, error) {
	var (
		err  error
		stmt *sql.Stmt
		res  sql.Result
		id   int64
	)

	stmt, err = this.db.Prepare("INSERT INTO bin_file_group(md5,created) values(?,?)")
	if err != nil {
		return 0, err
	}
	created := time.Now().Format("2006-01-02T15:04:05Z")
	res, err = stmt.Exec(md5, created)
	if err != nil {
		return 0, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (this *DB) DeleteFileGroupById(id int64) error {
	var (
		err  error
		stmt *sql.Stmt
	)
	stmt, err = this.db.Prepare("DELETE FROM bin_file_group where id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func (this *DB) DeleteFileGroupByMd5(md5 string) error {
	var (
		err  error
		stmt *sql.Stmt
	)
	stmt, err = this.db.Prepare("DELETE FROM bin_file_group where md5=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(md5)
	if err != nil {
		return err
	}
	return nil
}

func (this *DB) FindFileByMd5(md5 string) (*BinFile, bool) {
	var (
		err  error
		rows *sql.Rows
	)
	di := &BinFile{}

	rows, err = this.db.Query(fmt.Sprintf("select id,md5,path,node_num,node,size,created from bin_file where md5='%s' limit 1", md5))
	if err != nil {
		return di, false
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&di.Id, &di.Md5, &di.Path, &di.NodeNum, &di.Node, &di.Size, &di.Created)
		if err != nil {
			return di, false
		}
		return di, true
	}

	return di, false
}

func (this *DB) AddFileRow(md5 string, gid int64, path string, node_num int, node string, size int64) error {
	var (
		err  error
		stmt *sql.Stmt
	)
	stmt, err = this.db.Prepare("INSERT INTO bin_file(md5,gid,path,node_num,node,size,created) values(?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	created := time.Now().Format("2006-01-02T15:04:05Z")
	_, err = stmt.Exec(md5, gid, path, node_num, node, size, created)
	if err != nil {
		return err
	}
	return nil
}

func (this *DB) UpdateFileNode(bf *BinFile) error {
	var (
		err  error
		stmt *sql.Stmt
	)
	stmt, err = this.db.Prepare("UPDATE bin_file set node_num=?, node=? where id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(&bf.NodeNum, &bf.Node, &bf.Id)
	if err != nil {
		return err
	}
	return nil
}

func (this *DB) DeleteRowById(id int64) error {
	var (
		err  error
		stmt *sql.Stmt
	)
	stmt, err = this.db.Prepare("DELETE FROM bin_file where id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func (this *DB) DeleteRowByMd5(md5 string) error {
	var (
		err  error
		stmt *sql.Stmt
	)
	stmt, err = this.db.Prepare("DELETE FROM bin_file where md5=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(md5)
	if err != nil {
		return err
	}
	return nil
}

func Open(path string) *DB {
	ds := &DB{}
	ds.Init(path)
	return ds
}
