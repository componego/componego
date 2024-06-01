package migration

import (
	"github.com/componego/componego/examples/url-shortener-app/pkg/components/database"
)

const sql = `
CREATE TABLE IF NOT EXISTS redirects (
    entity_id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    url_key VARCHAR(10) NOT NULL,
    url VARCHAR(255) NOT NULL,
    PRIMARY KEY(entity_id),
    UNIQUE KEY message_key(url_key)
) ENGINE = INNODB DEFAULT CHARSET = utf8mb3 COLLATE = utf8mb3_general_ci COMMENT = 'Table with redirects'
`

func Run(dbProvider database.Provider) error {
	// Note that dependency 'dbProvider' is present in the application
	// because we added a component to the application that provides that dependency.
	db, err := dbProvider.Get("main-storage")
	if err == nil {
		_, err = db.Exec(sql)
	}
	return err
}
