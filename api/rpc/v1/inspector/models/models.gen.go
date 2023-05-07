// Package models provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package models

import (
	"time"
)

// Defines values for FieldType.
const (
	Dns   FieldType = "dns"
	Fqdn  FieldType = "fqdn"
	Whois FieldType = "whois"
)

// Defines values for OperationType.
const (
	Create OperationType = "create"
	Delete OperationType = "delete"
	Update OperationType = "update"
)

// AnyValue Can be any value, including `null`.
type AnyValue = interface{}

// Changelog defines model for Changelog.
type Changelog struct {
	FieldType FieldType `json:"field_type"`

	// From Can be any value, including `null`.
	From          AnyValue      `json:"from"`
	OperationType OperationType `json:"operation_type"`
	Path          []string      `json:"path"`

	// To Can be any value, including `null`.
	To AnyValue `json:"to"`
}

// Changelogs defines model for Changelogs.
type Changelogs = []Changelog

// Domain defines model for Domain.
type Domain struct {
	Dns   ResourceRecords `json:"dns"`
	Fqdn  string          `json:"fqdn"`
	Whois WhoisRecords    `json:"whois"`
}

// Error defines model for Error.
type Error struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// FieldType defines model for FieldType.
type FieldType string

// MX defines model for MX.
type MX struct {
	Host string `json:"host"`

	// Pref Host preference
	Pref int64 `json:"pref"`
}

// NS defines model for NS.
type NS struct {
	Host string `json:"host"`
}

// OperationType defines model for OperationType.
type OperationType string

// ResourceRecords defines model for ResourceRecords.
type ResourceRecords struct {
	A     []string  `json:"A"`
	AAAA  *[]string `json:"AAAA,omitempty"`
	CNAME *string   `json:"CNAME,omitempty"`
	MX    *[]MX     `json:"MX,omitempty"`
	NS    *[]NS     `json:"NS,omitempty"`
	SRV   *[]SRV    `json:"SRV,omitempty"`
	TXT   *[]TXT    `json:"TXT,omitempty"`
}

// SRV defines model for SRV.
type SRV struct {
	// Port порт TCP или UDP, на котором работает сервис
	Port int64 `json:"port"`

	// Priority приоритет целевого хоста, более низкое значение означает более предпочтительный
	Priority int64 `json:"priority"`

	// Target канонические имя машины, предоставляющей сервис
	Target string `json:"target"`

	// Weight относительный вес для записей c одинаковым приоритетом
	Weight int64 `json:"weight"`
}

// TXT defines model for TXT.
type TXT = string

// WhoisRecords defines model for WhoisRecords.
type WhoisRecords struct {
	// Created Дата и время регистрации домена в формате UTC
	Created *time.Time `json:"created,omitempty"`

	// DomainName Доменное имя
	DomainName string `json:"domainName"`

	// NameServers Список DNS-серверов, указанных для домена
	NameServers *[]string `json:"nameServers,omitempty"`

	// PaidTill Дата и время окончания срока регистрации домена в формате UTC
	PaidTill time.Time `json:"paidTill"`

	// Registrar Регситратор домена
	Registrar *string `json:"registrar,omitempty"`
}

// AddDomainParams defines parameters for AddDomain.
type AddDomainParams struct {
	Fqdn string `form:"fqdn" json:"fqdn"`
}

// GetChangelogsParams defines parameters for GetChangelogs.
type GetChangelogsParams struct {
	Fqdn string `form:"fqdn" json:"fqdn"`
}

// DeleteDomainParams defines parameters for DeleteDomain.
type DeleteDomainParams struct {
	Fqdn string `form:"fqdn" json:"fqdn"`
}

// GetDomainParams defines parameters for GetDomain.
type GetDomainParams struct {
	Fqdn string `form:"fqdn" json:"fqdn"`
}

// UpdateDomainParams defines parameters for UpdateDomain.
type UpdateDomainParams struct {
	Fqdn string `form:"fqdn" json:"fqdn"`
}
