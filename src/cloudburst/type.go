package main


type Pincode struct {
	Id  uint64 `json:"pincode"`
}

type Restaurant struct {
	Id      uint64 `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type RestaurantList struct {
	Restaurants []Restaurant `json:"restaurantlist"`
}

type Item struct {
	Id          uint64 `json:"id"`
	Name        string `json:"name"`
	Price      float64 `json:"price"`
	Description string `json:"description"`
}

type ItemList struct {
	Menu []Item `json:"menu"`
}

type CartItem struct {
	Id      uint64 `json:"id"` //item id
	Name    string `json:name`
	Quantity uint8 `json:"quantity"`
}

type Cart struct{
	Id             string `json:"id"`
	RestaurantId   uint64 `json:"restaurantId"`
	RestaurantName string `json:"restaurantName"`
	Items      []CartItem `json:"items"`
}

type Order struct {
	UserId       string `json:"userid"`
	RestaurantId uint64 `json:"restaurantId"`
	RestaurantName string `json:"restaurantName"`
	Items    []CartItem `json:"items"`
	Id           string `json:"id"`   //will be same as cart id and cart gets deleted.
	OrderStatus  string `json:"status"`
}

type User struct {
	Id       string `json:"id"`
	Password string `json:"password"`
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
}

type OrderList struct {
	Orders      []string `json:"orderlist"`
}

/*
type Offer struct {
	Offer      string `json:"offer"`
	Validity string `json:"validity"`
}

type OfferList struct {
	Offers      []string `json:"offerslist"`
}
*/