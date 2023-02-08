package products

type Product struct {
	ID              int     `json:"id"`
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	AuthUser        string  `json:"auth_user"`
	AuthPass        string  `json:"auth_pass"`
	Price           float64 `json:"price"`
	RenewalDay      int     `json:"renewal_day"`
	UrlNotifSub     string  `json:"url_notif_sub"`
	UrlNotifUnsub   string  `json:"url_notif_unsub"`
	UrlNotifRenewal string  `json:"url_notif_renewal"`
	UrlPostback     string  `json:"url_postback"`
}
