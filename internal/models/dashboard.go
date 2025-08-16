package models

type ResponseDashboard struct {
	CountCustomer int                      `json:"count_customer"`
	CountLeads    int                      `json:"count_leads"`
	CountDeals    int                      `json:"count_deals"`
	SummaryDeals  []ListCountGroupByStatus `json:"summary_deals"`
}

type ListCountGroupByStatus struct {
	Count int    `json:"count"`
	Label string `json:"label"`
}
