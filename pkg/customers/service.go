package customers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var ErrNotFound = errors.New("item not found")
var ErrInternal = errors.New("internal error")
type Service struct{
	pool *pgxpool.Pool
}
func NewService(pool *pgxpool.Pool)*Service  {
	return &Service{pool: pool}
}
type Customer struct{
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Phone   string    `json:"phone"`
	Active  bool      `json:"active"`
	Created time.Time `json:"created"`
}
func (s *Service)ByID(ctx context.Context,id int64)(*Customer,error) {
	item:=&Customer{}
	err:=s.pool.QueryRow(ctx,`
	select id,name,phone,active,created from customers where id=$1
	`,id).Scan(&item.ID,&item.Name,&item.Phone,&item.Active,&item.Created)
	if errors.Is(err,pgx.ErrNoRows) {
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
	rows,err:=s.pool.Query(ctx,`
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
	rows,err:=s.pool.Query(ctx,`
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
		err:=s.pool.QueryRow(ctx,`
		insert into customers(name,phone) values($1,$2) 
		on conflict(phone) do update set name=excluded.name 
		returning id,name,phone,active,created
		`,name,phone).Scan(&item.ID,&item.Name,&item.Phone,&item.Active,&item.Created)
		log.Print("in save id 0")
		if err!=nil {
			log.Print(err)
			return nil,ErrInternal
		}
		return item,nil
	}
	err:=s.pool.QueryRow(ctx,`
	update customers set name=$2,phone=$3
	where id=$1
	RETURNING id,name,phone,active,created
	`,id,name,phone).Scan(&item.ID,&item.Name,&item.Phone,&item.Active,&item.Created)
	if err!=nil {
		log.Print(err)
		return nil,ErrInternal
	}
	return item,nil
}
func (s *Service)RemoveByID(ctx context.Context,id int64)(error) {
	_,err:=s.pool.Exec(ctx,`
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
	err:=s.pool.QueryRow(ctx,`
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
	err:=s.pool.QueryRow(ctx,`
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