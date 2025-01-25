package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Buntala/rbac/models"
)

var (
    permission_id  uint = 1
)

type IPermissionController interface {
    UrlHandler(w http.ResponseWriter, r *http.Request)
    UrlWithIdHandler(w http.ResponseWriter, r *http.Request) 
    GetAllPermission(w http.ResponseWriter, r *http.Request)
    GetPermissionById(w http.ResponseWriter, r *http.Request)
    Create(w http.ResponseWriter, r *http.Request)
}

type permissionController struct {
    permissions *[]models.Permission
}

func NewPermissionController(permissions *[]models.Permission) IPermissionController{
    return &permissionController{permissions}
}

func (c *permissionController) UrlHandler(w http.ResponseWriter, r *http.Request){
    switch r.Method{
    case "POST":
        c.Create(w,r)
        return
    case "GET":
        c.GetAllPermission(w,r)
        return
    default:
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Invalid Method"))
    }
}

func (c *permissionController) UrlWithIdHandler(w http.ResponseWriter, r *http.Request){
    switch r.Method{
    case "GET":
        c.GetPermissionById(w,r)
        return
    case "PUT":
        c.Update(w, r)
        return
    case "DELETE":
        c.Delete(w, r)
        return
    default:
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Invalid Method"))
    }
}

func (c *permissionController) GetAllPermission(w http.ResponseWriter, r *http.Request){
    resp, err := json.Marshal(c.permissions)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
    }
    w.Header().Set("Content-Type","application/json")
    w.Write(resp)
}

func (c *permissionController) GetPermissionById(w http.ResponseWriter, r *http.Request) {
    query_id := r.URL.Path[len(r.URL.Path)-1]
    id, err := strconv.Atoi(string(query_id))

    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Error Parsing Id"))
        return
    }

    for _, r := range *c.permissions {
        if r.Id == uint(id) {
            resp, err := json.Marshal(r)
            if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                w.Write([]byte(err.Error()))
                return
            }
            w.Header().Set("Content-Type","application/json")
            w.Write(resp)
            return
        }
    }
    w.WriteHeader(http.StatusNotFound)
    w.Write([]byte("id not found"))
}

func (c *permissionController) Create(w http.ResponseWriter, r *http.Request) {
    var input models.Permission
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    json.Unmarshal(b, &input)
    input.Id = permission_id

    *c.permissions = append(*c.permissions, input)
    permission_id += 1

    resp, err := json.Marshal(input)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
    }
    w.Header().Set("Content-Type","application/json")
    w.Write([]byte(resp))
    return
}

func (c *permissionController) Update(w http.ResponseWriter, r *http.Request) {
    query_id := r.URL.Path[len(r.URL.Path)-1]
    id, err := strconv.Atoi(string(query_id))

    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Error Parsing Id"))
        return
    }

    var input models.Permission
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    json.Unmarshal(b, &input)
    input.Id = uint(id)

    for idx := range *c.permissions {
        if (*c.permissions)[idx].Id == uint(id) {
            resp, err := json.Marshal(input)
            if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                w.Write([]byte(err.Error()))
                return
            }
            (*c.permissions)[idx] = input
            w.Header().Set("Content-Type","application/json")
            w.Write([]byte(resp))
            return
        }
    }

    w.WriteHeader(http.StatusNotFound)
    w.Header().Set("Content-Type","application/json")
    w.Write([]byte("id not found"))
}

func (c *permissionController) Delete(w http.ResponseWriter, r *http.Request){
    query_id := r.URL.Path[len(r.URL.Path)-1]
    id, err := strconv.Atoi(string(query_id))

    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Error Parsing Id"))
        return
    }

    for idx := range *c.permissions {
        if (*c.permissions)[idx].Id == uint(id) {
            var deleted_data models.Permission
            deleted_data = (*c.permissions)[idx]
            if idx == len(*c.permissions) - 1 {
                *c.permissions = (*c.permissions)[:idx]
            } else {
                *c.permissions = append((*c.permissions)[:idx], (*c.permissions)[idx+1])
            }

            resp, err := json.Marshal(deleted_data)
            if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                w.Write([]byte(err.Error()))
                return
            }
            w.Header().Set("Content-Type","application/json")
            w.Write([]byte(resp))
            return
        }
    }

    w.WriteHeader(http.StatusNotFound)
    w.Header().Set("Content-Type","application/json")
    w.Write([]byte("id not found"))
}
