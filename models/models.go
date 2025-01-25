package models

type User struct {
    Id         uint
    Name       string
    Username   string
    Password   string
    Roles      []*Role
}

type Role struct {
    Id     uint
    Name   string
    Permissions []*Permission
}

type Permission struct {
    Id     uint
    Name   string
    Url    string
}

type UserInput struct {
    Name       string
    Username   string
    Password   string
    RoleIds    []uint `json:"role_id"`
}

type RoleInput struct {
    Name           string
    PermissionIds  []uint `json:"permission_id"`
}
