package ru

import (
	"context"
	"fmt"
	"strings"
	"time"

	whoisentity "github.com/muratom/domain-monitoring/services/emitter/internal/core/domain/whois"
	"github.com/muratom/domain-monitoring/services/emitter/internal/core/service/whois"
	"github.com/muratom/domain-monitoring/services/emitter/internal/core/service/whois/adapter"
)

// https://tcinet.ru/documents/whois_ru_rf.pdf
type Adapter struct {
	whoisClient adapter.Client
	whoisServer string
}

type domainState string

const (
	// Домен зарегистрирован
	Registered domainState = "REGISTERED"

	// Домен делегирован
	Delegated domainState = "DELEGATED"

	// Домен не делегирован
	NotDelegated domainState = "NOT DELEGATED"

	// Срок регистрации и период преимущественного продления домена истекли,
	// домен ожидает автоматического удаления
	PendingDelete domainState = "pendingDelete"

	// Данные администратора домена проверены регистратором
	Verified domainState = "VERIFIED"

	// Данные администратора не были проверены регистратором
	Unverified domainState = "UNVERIFIED"
)

type Response struct {
	// Domain Доменное имя
	domain string

	// Список DNS-серверов, указанных для домена
	// (если имя сервера содержит имя домена, то указываются также его IP-адреса)
	nserver []string

	// Состояние доменного имени
	state []domainState

	// Идентификатор регистратора, осуществляющего поддержку сведений о доменном имени, в Реестре.
	registrar string

	// Дата и время регистрации домена в формате UTC; не изменяется при продлении срока
	// регистрации или при смене администратора или регистратора домена.
	created time.Time

	// Дата и время окончания срока регистрации домена в формате UTC.
	paidTill time.Time
}

func NewAdapter(whoisClient adapter.Client, whoisProvider whois.ServerProvider) *Adapter {
	server, err := whoisProvider.GetServerByDomain("ru")
	if err != nil {
		return nil
	}
	return &Adapter{
		whoisClient: whoisClient,
		whoisServer: server,
	}
}

func (a *Adapter) MakeRequest(ctx context.Context, req *whois.Request) (*whois.Response, error) {
	resp, err := a.whoisClient.FetchWhois(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch a response from WHOIS server: %w", err)
	}

	return resp, nil
}

func (a *Adapter) PrepareRequest(ctx context.Context, fqdn string) (*whois.Request, error) {
	body := []byte(fmt.Sprintf("%s\r\n", fqdn))

	return &whois.Request{
		WhoisServer: a.whoisServer,
		Body:        body,
	}, nil
}

func (a *Adapter) ParseResponse(ctx context.Context, resp *whois.Response) (*whoisentity.Record, error) {
	ruResponse, err := parseResponse(ctx, resp.Raw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse raw response from RU WHOIS server: %w", err)
	}

	return &whoisentity.Record{
		DomainName:  ruResponse.domain,
		NameServers: ruResponse.nserver,
		Created:     ruResponse.created,
		PaidTill:    ruResponse.paidTill,
	}, nil
}

type whoisResponse []byte

func (r *whoisResponse) findValueFor(field string) string {
	value := ""
	textToSearch := string(*r)
	if valueStartIdx := strings.Index(textToSearch, field); valueStartIdx != -1 {
		valueStartIdx += len(field) + 1 // название поля + двоеточие
		for textToSearch[valueStartIdx] == ' ' {
			valueStartIdx++
		}
		valueEndShift := strings.IndexByte(textToSearch[valueStartIdx:], '\n')
		if valueEndShift > -1 {
			value = textToSearch[valueStartIdx : valueStartIdx+valueEndShift]
		}
	}
	return value
}

// TODO: поддержать парсинг nserver
func parseResponse(ctx context.Context, rawResponse []byte) (*Response, error) {
	response := whoisResponse(rawResponse)

	created, err := time.Parse(time.RFC3339, response.findValueFor("created"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse field 'created': %w", err)
	}
	paidTill, err := time.Parse(time.RFC3339, response.findValueFor("paid-till"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse field 'paid-till': %w", err)
	}

	state := response.findValueFor("state")
	states := strings.Split(state, ", ")
	formattedStates := make([]domainState, len(states))
	for i := 0; i < len(states); i++ {
		formattedStates[i] = domainState(states[i])
	}

	return &Response{
		domain:    response.findValueFor("domain"),
		nserver:   nil,
		state:     formattedStates,
		registrar: response.findValueFor("registrar"),
		created:   created,
		paidTill:  paidTill,
	}, nil
}
