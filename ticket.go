/*

*/

package zd

import "fmt"

// Ticket struct
type Ticket struct {
	ID                  *float64            `json:"id,omitempty"`
	URL                 *string             `json:"url,omitempty"`
	ExternalID          *string             `json:"external_id,omitempty"`
	Type                *string             `json:"type,omitempty"`
	Subject             *string             `json:"subject,omitempty"`
	Description         *string             `json:"description,omitempty"`
	Priority            *string             `json:"priority,omitempty"`
	Status              *string             `json:"status,omitempty"`
	Recipient           *string             `json:"recipient,omitempty"`
	RequesterID         *float64            `json:"requester_id,omitempty"`
	SubmitterID         *float64            `json:"submitter_id,omitempty"`
	AssigneeID          *float64            `json:"assignee_id,omitempty"`
	OrganizationID      *float64            `json:"organization_id,omitempty"`
	GroupID             *float64            `json:"group_id,omitempty"`
	CollaboratorIDs     []float64           `json:"collaborator_ids,omitempty"`
	ForumTopicID        *float64            `json:"forum_topic_id,omitempty"`
	ProblemID           *float64            `json:"problem_id,omitempty"`
	HasIncidents        *bool               `json:"has_incidents,omitempty"`
	DueAt               *string             `json:"due_at,omitempty"`
	Tags                []string            `json:"tags,omitempty"`
	Via                 *Via                `json:"via,omitempty"`
	CustomFields        []CustomField       `json:"custom_fields,omitempty"`
	SatisfactionRating  *SatisfactionRating `json:"satisfaction_rating,omitempty"`
	SharingAgreementIds []float64           `json:"sharing_agreement_ids,omitempty"`
	FollowupIds         []float64           `json:"followup_ids,omitempty"`
	TicketFormID        *float64            `json:"ticket_form_id,omitempty"`
	BrandID             *float64            `json:"brand_id,omitempty"`
	CreatedAt           *string             `json:"created_at,omitempty"`
	UpdatedAt           *string             `json:"updated_at,omitempty"`
	Comment             *TicketComment      `json:"comment,omitempty"`
	Requester           *TicketRequester    `json:"requester,omitempty"`
}

// Ticket Single Response
type TicketSingleResponse struct {
	Ticket *Ticket `json:"ticket"`
}

type TicketComment struct {
	Body string `json:"body,omitempty"`
}

type TicketRequester struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// Via struct
type Via struct {
	Channel *string      `json:"channel,omitempty"`
	Source  *interface{} `json:"source,omitempty"`
}

// CustomField struct
type CustomField struct {
	ID    *float64     `json:"id,omitempty"`
	Value *interface{} `json:"value,omitempty"`
}

// SatisfactionRating struct
type SatisfactionRating struct {
	ID      *float64 `json:"id,omitempty"`
	Score   *string  `json:"score,omitempty"`
	Comment *string  `json:"comment,omitempty"`
}

// TicketResponse struct
type TicketResponse struct {
	Results  []Ticket `json:"tickets"`
	Next     *string  `json:"next_page,omitempty"`
	Previous *string  `json:"previous_page,omitempty"`
	Count    *int     `json:"count,omitempty"`
}

// TicketUserGroupResponse struct
type TicketUserGroupResponse struct {
	Tickets  []Ticket `json:"tickets,omitempty"`
	Users    []User   `json:"users,omitempty"`
	Groups   []Group  `json:"groups,omitempty"`
	Next     *string  `json:"next_page,omitempty"`
	Previous *string  `json:"previous_page,omitempty"`
	Count    *int     `json:"count,omitempty"`
}

// TicketService struct
type TicketService struct {
	client *Client
}

// List returns a slice of all products
func (s *TicketService) List() ([]Ticket, error) {

	var resource []Ticket

	rp, _, _, err := s.getPage("")

	if err != nil {
		return nil, err
	}

	resource = append(resource, rp...)

	// for next != nil {
	// 	rp, nx, _, err := s.getPage(*next)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	next = nx
	// 	resource = append(resource, *rp...)
	// }

	return resource, err
}

// ListByView function
func (s *TicketService) ListByView(id string) ([]Ticket, error) {

	var resource []Ticket

	url := fmt.Sprintf("views/%s/tickets.json", id)
	rp, next, _, err := s.getPage(url) // Hardcoded "Problem tickets" view
	if err != nil {
		return nil, err
	}

	resource = append(resource, rp...)

	for next != nil {
		rp, nx, _, err := s.getPage(*next)
		if err != nil {
			return nil, err
		}
		next = nx
		resource = append(resource, rp...)
	}

	return resource, err
}

// ListByViewUG function
func (s *TicketService) ListByViewUG(id string) ([]Ticket, []User, []Group, error) {

	var tickets []Ticket
	var users []User
	var groups []Group

	url := fmt.Sprintf("views/%s/tickets.json?include=users,groups", id)
	tkts, usrs, grps, next, _, err := s.getPageUG(url)
	if err != nil {
		return nil, nil, nil, err
	}

	tickets = append(tickets, tkts...)
	users = append(users, usrs...)
	groups = append(groups, grps...)

	for next != nil {
		tkts, usrs, grps, nx, _, err := s.getPageUG(*next)
		if err != nil {
			return nil, nil, nil, err
		}
		next = nx
		tickets = append(tickets, tkts...)
		users = append(users, usrs...)
		groups = append(groups, grps...)
	}

	return tickets, users, groups, err
}

// GetProblemIncidents gets all problem tickets
func (s *TicketService) GetProblemIncidents(id string) ([]Ticket, error) {
	var resource []Ticket

	url := fmt.Sprintf("tickets/%s/incidents.json", id)

	rp, next, _, err := s.getPage(url)
	if err != nil {
		return nil, err
	}

	resource = append(resource, rp...)

	for next != nil {
		rp, nx, _, err := s.getPage(*next)
		if err != nil {
			return nil, err
		}
		next = nx
		resource = append(resource, rp...)
	}

	return resource, err
}

// GetProblemIncidentsCount gets only the first page of tickets which includes the count field
func (s *TicketService) GetProblemIncidentsCount(id string) (int, error) {

	url := fmt.Sprintf("tickets/%s/incidents.json", id)

	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	response := new(TicketResponse)
	_, err = s.client.Do(req, response)
	if err != nil {
		return 0, err
	}

	resource := response.Count
	return *resource, err
}

func (s *TicketService) getPage(url string) ([]Ticket, *string, *Response, error) {

	if url == "" {
		url = "tickets.json"
	}

	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	response := new(TicketResponse)
	resp, err := s.client.Do(req, response)
	if err != nil {
		return nil, nil, resp, err
	}

	next := response.Next
	resource := response.Results
	return resource, next, resp, err
}

func (s *TicketService) getPageUG(url string) ([]Ticket, []User, []Group, *string, *Response, error) {

	if url == "" {
		url = "tickets.json?include=users,groups"
	}

	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	response := new(TicketUserGroupResponse)
	resp, err := s.client.Do(req, response)
	if err != nil {
		return nil, nil, nil, nil, resp, err
	}

	next := response.Next
	tickets := response.Tickets
	users := response.Users
	groups := response.Groups
	return tickets, users, groups, next, resp, err
}

// GetOne method
func (s *TicketService) GetOne(id string) (*Ticket, *Response, error) {
	url := fmt.Sprintf("tickets/%s.json", id)

	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	ticketSingleResponse := new(TicketSingleResponse)
	resp, err := s.client.Do(req, &ticketSingleResponse)
	if err != nil {
		return nil, resp, err
	}

	return ticketSingleResponse.Ticket, resp, err
}

// Create func creates a single Ticket
func (s *TicketService) Create(t *Ticket) (*Ticket, error) {
	var ticket *Ticket
	var err error

	tw := &TicketSingleResponse{Ticket: t}
	url := fmt.Sprintf("tickets.json")

	req, err := s.client.NewRequest("POST", url, tw)
	if err != nil {
		return ticket, err
	}

	result := new(TicketSingleResponse)
	_, err = s.client.Do(req, result)
	if err != nil {
		return ticket, err
	}

	ticket = result.Ticket
	return ticket, err
}
