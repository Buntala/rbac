package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Buntala/rbac/models"
)

var (
    role_id  uint = 1
)

type IRoleController interface {
    UrlHandler(w http.ResponseWriter, r *http.Request)
    UrlWithIdHandler(w http.ResponseWriter, r *http.Request) 
    GetAllRole(w http.ResponseWriter, r *http.Request)
}

type roleController struct {
    roles *[]models.Role
    permissions *[]models.Permission
}

func NewRoleController(roles *[]models.Role, permissions *[]models.Permission) IRoleController{
    return &roleController{roles, permissions}
}

func (c *roleController) UrlHandler(w http.ResponseWriter, r *http.Request){
    switch r.Method{
    case "POST":
        c.Create(w,r)
        return
    case "GET":
        c.GetAllRole(w,r)
        return
    default:
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Invalid Method"))
    }
}

func (c *roleController) UrlWithIdHandler(w http.ResponseWriter, r *http.Request){
    switch r.Method{
    case "GET":
        c.GetRoleById(w,r)
        return
    case "PUT":
        c.Update(w, r)
        return
    case "DELETE":
        c.Delete(w, r)
        return
    default:
        w.WriteHeader(http.StatusBadRequest)
        w.Header().Set("Content-Type","application/json")
        w.Write([]byte("Invalid Method"))
    }
}

func (c *roleController) GetAllRole(w http.ResponseWriter, r *http.Request) {
    resp, err := json.Marshal(c.roles)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
    }
    w.Header().Set("Content-Type","application/json")
    w.Write(resp)
    return
}

func (c *roleController) GetRoleById(w http.ResponseWriter, r *http.Request) {
    query_id := r.URL.Path[len(r.URL.Path)-1]
    id, err := strconv.Atoi(string(query_id))
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Error Parsing Id"))
        return
    }
    for _, r := range *c.roles {
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

func (c *roleController) Create(w http.ResponseWriter, r *http.Request) {
    var input models.RoleInput
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    json.Unmarshal(b, &input)
    
    var permissions []*models.Permission
    
    for _, id := range input.PermissionIds{
        permission, err := c.fetchPermissionById(id)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(err.Error()))
            return
        }
        permissions = append(permissions, permission)
    }
    
    insert_data := models.Role{
        Id: role_id,
        Name: input.Name,
        Permissions: permissions,
    }

    *c.roles = append(*c.roles, insert_data)
    role_id += 1

    resp, err := json.Marshal(insert_data)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
    }
    w.Header().Set("Content-Type","application/json")
    w.Write([]byte(resp))
    return
}

func (c *roleController) Update(w http.ResponseWriter, r *http.Request){
    query_id := r.URL.Path[len(r.URL.Path)-1]
    id, err := strconv.Atoi(string(query_id))

    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Error Parsing Id"))
        return
    }

    var input models.RoleInput
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    json.Unmarshal(b, &input)

    var permissions []*models.Permission
    
    for _, id := range input.PermissionIds{
        permission, err := c.fetchPermissionById(id)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(err.Error()))
            return
        }
        permissions = append(permissions, permission)
    }
    
    update_data := models.Role{
        Id: uint(id),
        Name: input.Name,
        Permissions: permissions,
    }

    for idx := range *c.roles {
        if (*c.roles)[idx].Id == uint(id) {
            resp, err := json.Marshal(update_data)
            if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                w.Write([]byte(err.Error()))
                return
            }
            (*c.roles)[idx] = update_data
            w.Header().Set("Content-Type","application/json")
            w.Write([]byte(resp))
            return 
        }
    }

    w.WriteHeader(http.StatusNotFound)
    w.Write([]byte("id not found"))
}

func (c *roleController) Delete(w http.ResponseWriter, r *http.Request) {
    query_id := r.URL.Path[len(r.URL.Path)-1]
    id, err := strconv.Atoi(string(query_id))

    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Error Parsing Id"))
        return
    }

    for idx := range *c.roles {
        if (*c.roles)[idx].Id == uint(id) {
            var deleted_data models.Role
            deleted_data = (*c.roles)[idx]
            if idx == len(*c.roles) - 1 {
                *c.roles = (*c.roles)[:idx]
            } else {
                *c.roles = append((*c.roles)[:idx], (*c.roles)[idx+1])
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

    return
}

func (c *roleController) fetchPermissionById(id uint) (*models.Permission, error){
    for i := range *c.permissions {
        if (*c.permissions)[i].Id == id {
            return &(*c.permissions)[i], nil
        }
    }

    return nil, errors.New("id not found")
}
