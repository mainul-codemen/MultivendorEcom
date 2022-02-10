# MultivendorEcom

1. Making http Connection:
```go
package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func main() {
	fmt.Println("Hello")
	http.HandleFunc("/hello", hello)
	http.ListenAndServe(":8090", nil)
}
```

## Important link share between us

https://opencart.opencartworks.com/themes/so_emarket/index.php?route=common/home

http://preview.themeforest.net/item/nest-multipurpose-ecommerce-html-template/full_screen_preview/33948410?_ga=2.199365225.69734089.1642430589-1728743351.1633148608

https://docs.google.com/document/d/1XV2MTPnYtKC9ln54rDg64B_Xmuy2P_kdAqVeWoNNRQI/edit

Date and month calculation :

https://www.sqlines.com/postgresql/how-to/datediff

https://www.postgresqltutorial.com/postgresql-joins/