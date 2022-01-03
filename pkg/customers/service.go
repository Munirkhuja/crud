package customers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

var ErrNotFound = errors.New("item not found")
var ErrInternal = errors.New("internal error")
type Service struct{
	db *sql.DB
}
func NewService(db *sql.DB)*Service  {
	return &Service{db: db}
}
type Customer struct{
	ID int64 `json:"id"`
	Name string `json:"name"`
	Phone string `json:"phone"`
	Active string `json:"active"`
	Created time.Time `json:"created"`
}
func (s *Service)ByID(ctx context.Context,id int64)(*Customer,error) {
	item:=&Customer{}
	err:=s.db.QueryRowContext(ctx,`
	select id,name,phone,active,created from customers where id=$1
	`,id).Scan(&item.ID,&item.Name,&item.Phone,&item.Active,&item.Created)
	if errors.Is(err,sql.ErrNoRows) {
		return nil,ErrNotFound
	}
	if err!=nil {
		log.Print(err)
		return nil,ErrInternal
	}
	return item,nil
}
func (s *Service)All(ctx context.Context)([]*Customer,error) {
	items:=make([]*Customer,0)
	rows,err:=s.db.QueryContext(ctx,`
	select id,name,phone,active,created from customers`)
	if err!=nil {
		log.Print(err)
		return nil,ErrInternal
	}
	defer rows.Close()
	for rows.Next(){
		item:=&Customer{}
		err=rows.Scan(&item.ID,&item.Name,&item.Phone,&item.Active,&item.Created)
		if err!=nil {
			log.Print(err)
			return nil,err
		}
		items = append(items, item)
	}
	err=rows.Err()
	if err!=nil {
		log.Print(err)
		return nil,err
	}
	return items,nil
}
func (s *Service)AllActive(ctx context.Context)([]*Customer,error) {
	items:=make([]*Customer,0)
	rows,err:=s.db.QueryContext(ctx,`
	select id,name,phone,active,created from customers where active`)
	if err!=nil {
		log.Print(err)
		return nil,ErrInternal
	}
	defer rows.Close()
	for rows.Next(){
		item:=&Customer{}
		err=rows.Scan(&item.ID,&item.Name,&item.Phone,&item.Active,&item.Created)
		if err!=nil {
			log.Print(err)
			return nil,err
		}
		items = append(items, item)
	}
	err=rows.Err()
	if err!=nil {
		log.Print(err)
		return nil,err
	}
	return items,nil
}
func (s *Service)Save(ctx context.Context,id int64,name string,phone string)(*Customer,error) {
	item:=&Customer{}
	if id==0 {			
		err:=s.db.QueryRowContext(ctx,`
		insert into customers(name,phone) values($1,$2) 
		on conflict(phone) do update set name=excluded.name 
		RETURNING id,name,phone,active,created;
		`,name,phone).Scan(&item.ID,&item.Name,&item.Phone,&item.Active,&item.Created)
		if err!=nil {
			log.Print(err)
			return nil,ErrInternal
		}
	}else{
		err:=s.db.QueryRowContext(ctx,`
		update customers set name=$2,phone=$3
		where id=$1
		RETURNING id,name,phone,active,created;
		`,id,name,phone).Scan(&item.ID,&item.Name,&item.Phone,&item.Active,&item.Created)
		if err!=nil {
			log.Print(err)
			return nil,ErrInternal
		}
	}
	return item,nil
}
func (s *Service)RemoveByID(ctx context.Context,id int64)(error) {
	_,err:=s.db.ExecContext(ctx,`
	delete from customers where id=$1
	`,id)
	if err!=nil {
		log.Print(err)
		return ErrInternal
	}
	return nil
}
func (s *Service)BlockByID(ctx context.Context,id int64)(*Customer,error) {
	item:=&Customer{}
	err:=s.db.QueryRowContext(ctx,`
	update customers set active=false where id=$1	
	RETURNING id,name,phone,active,created
	`,id).Scan(&item.ID,&item.Name,&item.Phone,&item.Active,&item.Created)
	if errors.Is(err,sql.ErrNoRows) {
		return nil,ErrNotFound
	}
	if err!=nil {
		log.Print(err)
		return nil,ErrInternal
	}
	return item,nil
}

func (s *Service)UnblockByID(ctx context.Context,id int64)(*Customer,error) {
	item:=&Customer{}
	err:=s.db.QueryRowContext(ctx,`
	update customers set active=true where id=$1	
	RETURNING id,name,phone,active,created
	`,id).Scan(&item.ID,&item.Name,&item.Phone,&item.Active,&item.Created)
	if errors.Is(err,sql.ErrNoRows) {
		return nil,ErrNotFound
	}
	if err!=nil {
		log.Print(err)
		return nil,ErrInternal
	}
	return item,nil
}