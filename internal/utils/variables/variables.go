package variables

import "time"

var (
	CacheTTL        = 6 * time.Hour
	ItemKey         = "seafood_item:"
	AllItemsKey     = "seafood_items"
	AllSeafoodItems = "seafood_category"
)
