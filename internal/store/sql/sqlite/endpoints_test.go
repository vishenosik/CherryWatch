package sqlite

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_GetEndpoints_IN(t *testing.T) {

	query, args, err := sqlx.In(
		`
		SELECT 
			e.id,
			e.service_name,
			e.url,
			e.interval_seconds,
			(
				SELECT GROUP_CONCAT(code, ',') 
				FROM endpoint_success_codes 
				WHERE endpoint_id = e.id
			) AS success_codes,
			(
				SELECT GROUP_CONCAT(service_name, ',') 
				FROM endpoint_notification_services 
				WHERE endpoint_id = e.id
			) AS notification_services
		FROM 
			endpoints e
		WHERE 
			e.id IN (?)
	`,
		[]string{"id1", "id2", "id3"},
	)

	assert.NoError(t, err)

	fmt.Println(query, args)

	query = sqlx.Rebind(sqlx.DOLLAR, query) // приводим заполнители в нужный формат
	fmt.Println(query)
}
