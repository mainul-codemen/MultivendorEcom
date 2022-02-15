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

# Install Postgres Docker images

1. Pull image

    docker pull postgres:12-alpine

2. Start a postgres instance

    docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

3. Try to connect and exect its console

    docker exec -it postgres12 psql -U root

    root=# select now();
    
    \q

4. View Container logs

    dokcer logs postgres12

5. Go to postgresdb

    docker exec -it postgres12 /bin/sh

        pwd

        ls -l

        createdb --username=root --owner=root simpe_bank 

        psql

        dropdb simple_bank

        exit

6. Crate db from outside of the container

        docker exec -it postgres12 createdb --username=root --owner=root simple_bank

7. Excess Database from command

        docker exec -it postgres12 psql -U root simple_bank


8. Create Makefile

9. Stop and remove container

    docker stop postgres12

    docker ps

    docker rm postgres12

    docker ps -a

10. Run makefile

    make postgres

    docker ps

    make createdb

    