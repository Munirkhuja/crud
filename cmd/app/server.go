package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Munirkhuja/crud/pkg/customers"
)

type Server struct {
	mux *http.ServeMux
	customersSvc *customers.Service
}
func NewServer(mux *http.ServeMux,customersSvc *customers.Service)*Server  {
	return &Server{mux: mux,customersSvc: customersSvc}
}
func (s *Server)ServeHTTP(writer http.ResponseWriter,request *http.Request)  {
	s.mux.ServeHTTP(writer,request)
}
func (s *Server)Init()  {
	s.mux.HandleFunc("/customers.getById",s.handleGetCustomersByID)
	s.mux.HandleFunc("/customers.getAll",s.handleGetAll)
	s.mux.HandleFunc("/customers.getAllActive",s.handleGetAllActive)
	s.mux.HandleFunc("/customers.save",s.handleSave)
	s.mux.HandleFunc("/customers.removeById",s.handleRemoveByID)
	s.mux.HandleFunc("/customers.blockById",s.handleBlockByID)
	s.mux.HandleFunc("/customers.unblockById",s.handleUnblockByID)
}
func (s *Server)handleGetCustomersByID(writer http.ResponseWriter,request *http.Request)  {
	idParam:=request.URL.Query().Get("id")
	id,err:=strconv.ParseInt(idParam,10,64)
	if err!=nil {
		log.Print(err)
		http.Error(writer,http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return
	}
	item,err:=s.customersSvc.ByID(request.Context(),id)
	if errors.Is(err,customers.ErrNotFound) {
		http.Error(writer,http.StatusText(http.StatusNotFound),http.StatusNotFound)
		return
	}
	if err!=nil {
		http.Error(writer,http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return		
	}
	
	data,derr:=json.Marshal(item)
	if derr!=nil {
		log.Print(derr)
		http.Error(writer,http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return				
	}
	writer.Header().Set("Content-Type","aplication/json")
	_,err=writer.Write(data)
	if err!=nil {
		log.Print(err)
	}
}
func (s *Server) handleGetAll(writer http.ResponseWriter, request *http.Request) {
	all, err := s.customersSvc.All(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
		return
	}
	data, err := json.Marshal(all)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_, err = writer.Write([]byte(data))
	if err != nil {
		log.Print("Error!: Can't write anything on data.")
	}
}
// func (s *Server)handleGetAll(writer http.ResponseWriter,request *http.Request)  {	
// 	items,err:=s.customersSvc.All(request.Context())
// 	if errors.Is(err,customers.ErrNotFound) {
// 		http.Error(writer,http.StatusText(http.StatusNotFound),http.StatusNotFound)
// 		return
// 	}
// 	if err!=nil {
// 		http.Error(writer,http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
// 		return		
// 	}
	
// 	data,derr:=json.Marshal(items)
// 	if derr!=nil {
// 		log.Print(derr)
// 		http.Error(writer,http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
// 		return				
// 	}
// 	writer.Header().Set("Content-Type","aplication/json")
// 	_,err=writer.Write(data)
// 	if err!=nil {
// 		log.Print(err)
// 	}
// }
func (s *Server)handleGetAllActive(writer http.ResponseWriter,request *http.Request)  {	
	items,err:=s.customersSvc.AllActive(request.Context())
	if errors.Is(err,customers.ErrNotFound) {
		http.Error(writer,http.StatusText(http.StatusNotFound),http.StatusNotFound)
		return
	}
	if err!=nil {
		http.Error(writer,http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return		
	}
	
	data,derr:=json.Marshal(items)
	if derr!=nil {
		log.Print(derr)
		http.Error(writer,http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return				
	}
	writer.Header().Set("Content-Type","aplication/json")
	_,err=writer.Write(data)
	if err!=nil {
		log.Print(err)
	}
}

func (s *Server)handleSave(writer http.ResponseWriter,request *http.Request)  {
	parseerr:=request.ParseMultipartForm(10*1024*1024)
	if parseerr!=nil {
		log.Print(parseerr)
		http.Error(writer,http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return
	}
	postForm:=request.PostForm
	name,phone:="",""
	idParam,exists:=postForm["id"]
	if !exists {
		idParam=[]string{"0"}
	}
	id,pierr:=strconv.ParseInt(idParam[0],10,64)
	if pierr!=nil {
		id=0
		log.Print(pierr)
	}
	if nameParam,exists:=postForm["name"];exists {
		name=nameParam[0]
	}
	if phoneParam,exists:=postForm["phone"];exists {
		phone=phoneParam[0]
	}
	item,err:=s.customersSvc.Save(request.Context(),id,name,phone)
	if errors.Is(err,customers.ErrNotFound) {
		http.Error(writer,http.StatusText(http.StatusNotFound),http.StatusNotFound)
		return
	}
	if err!=nil {
		http.Error(writer,http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return		
	}
	
	data,derr:=json.Marshal(item)
	if derr!=nil {
		log.Print(derr)
		http.Error(writer,http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return				
	}
	writer.Header().Set("Content-Type","aplication/json")
	_,err=writer.Write(data)
	if err!=nil {
		log.Print(err)
	}
}

func (s *Server)handleRemoveByID(writer http.ResponseWriter,request *http.Request)  {
	idParam:=request.URL.Query().Get("id")
	id,err:=strconv.ParseInt(idParam,10,64)
	if err!=nil {
		log.Print(err)
		http.Error(writer,http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return
	}
	err=s.customersSvc.RemoveByID(request.Context(),id)
	if errors.Is(err,customers.ErrNotFound) {
		http.Error(writer,http.StatusText(http.StatusNotFound),http.StatusNotFound)
		return
	}
	if err!=nil {
		http.Error(writer,http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return		
	}
	
	data,derr:=json.Marshal(id)
	if derr!=nil {
		log.Print(derr)
		http.Error(writer,http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return				
	}
	writer.Header().Set("Content-Type","aplication/json")
	_,err=writer.Write(data)
	if err!=nil {
		log.Print(err)
	}
}

func (s *Server)handleBlockByID(writer http.ResponseWriter,request *http.Request)  {
	idParam:=request.URL.Query().Get("id")
	id,err:=strconv.ParseInt(idParam,10,64)
	if err!=nil {
		log.Print(err)
		http.Error(writer,http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return
	}
	item,err:=s.customersSvc.BlockByID(request.Context(),id)
	if errors.Is(err,customers.ErrNotFound) {
		http.Error(writer,http.StatusText(http.StatusNotFound),http.StatusNotFound)
		return
	}
	if err!=nil {
		http.Error(writer,http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return		
	}
	
	data,derr:=json.Marshal(item)
	if derr!=nil {
		log.Print(derr)
		http.Error(writer,http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return				
	}
	writer.Header().Set("Content-Type","aplication/json")
	_,err=writer.Write(data)
	if err!=nil {
		log.Print(err)
	}
}
func (s *Server)handleUnblockByID(writer http.ResponseWriter,request *http.Request)  {
	idParam:=request.URL.Query().Get("id")
	id,err:=strconv.ParseInt(idParam,10,64)
	if err!=nil {
		log.Print(err)
		http.Error(writer,http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return
	}
	item,err:=s.customersSvc.UnblockByID(request.Context(),id)
	if errors.Is(err,customers.ErrNotFound) {
		http.Error(writer,http.StatusText(http.StatusNotFound),http.StatusNotFound)
		return
	}
	if err!=nil {
		http.Error(writer,http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return		
	}
	
	data,derr:=json.Marshal(item)
	if derr!=nil {
		log.Print(derr)
		http.Error(writer,http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		return				
	}
	writer.Header().Set("Content-Type","aplication/json")
	_,err=writer.Write(data)
	if err!=nil {
		log.Print(err)
	}
}