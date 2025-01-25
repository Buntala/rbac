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
    user_id  uint = 1
)

type IUserController interface {
    UrlHandler(w http.ResponseWriter, r *http.Request)
    UrlWithIdHandler(w http.ResponseWriter, r *http.Request) 
    Delete(w http.ResponseWriter, r *http.Request)
}

type userController struct {
    users *[]models.User
    roles *[]models.Role
}

func NewUserController(user *[]models.User, role *[]models.Role) IUserController{
    return &userController{user, role}
}

func (c *userController) UrlHandler(w http.ResponseWriter, r *http.Request){
    switch r.Method{
    case "POST":
        c.Create(w,r)
        return
    case "GET":
        c.GetAllUser(w,r)
        return
    default:
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Invalid Method"))
    }
}

func (c *userController) UrlWithIdHandler(w http.ResponseWriter, r *http.Request){
    switch r.Method{
    case "GET":
        c.GetUserById(w,r)
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

func (c *userController) GetAllUser(w http.ResponseWriter, r *http.Request){
    resp, err := json.Marshal(c.users)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
    }
    w.Header().Set("Content-Type","application/json")
    w.Write(resp)
}

func (c *userController) GetUserById(w http.ResponseWriter, r *http.Request) {
    query_id := r.URL.Path[len(r.URL.Path)-1]
    id, err := strconv.Atoi(string(query_id))
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Error Parsing Id"))
        return
    }
    for _, r := range *c.users {
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
    w.Header().Set("Content-Type","application/json")
    w.Write([]byte("id not found"))
}

func (c *userController) Create(w http.ResponseWriter, r *http.Request) {
    var input models.UserInput
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    json.Unmarshal(b, &input)

    var roles []*models.Role

    for _, id := range input.RoleIds{
        role, err := c.fetchRolesById(id)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(err.Error()))
            return
        }
        roles = append(roles, role)
    }

    insert_data := models.User{
        Id: user_id,
        Name: input.Name,
        Username: input.Username,
        Password: input.Password,
        Roles: roles,
    }
    *c.users = append(*c.users, insert_data)
    user_id += 1

    resp, err := json.Marshal(insert_data)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Header().Set("Content-Type","application/json")
        w.Write([]byte(err.Error()))
        return
    }
    w.Header().Set("Content-Type","application/json")
    w.Write([]byte(resp))
}

func (c *userController) Update(w http.ResponseWriter, r *http.Request) {
    query_id := r.URL.Path[len(r.URL.Path)-1]
    id, err := strconv.Atoi(string(query_id))

    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Error Parsing Id"))
        return
    }

    var input models.UserInput
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    json.Unmarshal(b, &input)

    var roles []*models.Role

    for _, id := range input.RoleIds{
        role, err := c.fetchRolesById(id)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(err.Error()))
            return
        }
        roles = append(roles, role)
    }

    update_data := models.User{
        Id: uint(id),
        Name: input.Name,
        Username: input.Username,
        Password: input.Password,
        Roles: roles,
    }

    for idx := range *c.users {
        if (*c.users)[idx].Id == uint(id) {
            resp, err := json.Marshal(update_data)
            if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                w.Write([]byte(err.Error()))
                return
            }
            (*c.users)[idx] = update_data
            w.Header().Set("Content-Type","application/json")
            w.Write([]byte(resp))
            return
        }
    }

    return
}

func (c *userController) Delete(w http.ResponseWriter, r *http.Request){
    query_id := r.URL.Path[len(r.URL.Path)-1]
    id, err := strconv.Atoi(string(query_id))
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Error Parsing Id"))
        return
    }

    for idx := range *c.users {
        if (*c.users)[idx].Id == uint(id) {
            var deleted_data models.User
            deleted_data = (*c.users)[idx]
            if idx == len(*c.users) - 1 {
                (*c.users) = (*c.users)[:idx]
            } else {
                *c.users = append((*c.users)[:idx], (*c.users)[idx+1])
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
    w.Write([]byte("id not found"))
}

func (c *userController) fetchRolesById(id uint) (*models.Role, error){
    for i := range *c.roles {
        if (*c.roles)[i].Id == id {
            return &(*c.roles)[i], nil
        }
    }
    return nil, errors.New("id not found")
}
