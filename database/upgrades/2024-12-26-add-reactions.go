package upgrades

import (
	"database/sql"
)

func init() {
	upgrades[21] = upgrade{"Add reaction store to database", func(tx *sql.Tx, ctx context) error {
		_, err := tx.Exec(`CREATE TABLE reactions (
			id            CHAR(13),
			chat_jid      VARCHAR(255),
			chat_receiver VARCHAR(255),
			message_id    CHAR(13),
			mxid          VARCHAR(255) NOT NULL UNIQUE,
			sender        VARCHAR(255) NOT NULL,
			timestamp     BIGINT       NOT NULL,
			content       bytea        NOT NULL,
			PRIMARY KEY (chat_jid, chat_receiver, message_id, id),
			FOREIGN KEY (chat_jid, chat_receiver) REFERENCES portal(jid, receiver) ON DELETE CASCADE
		)`)
		if err != nil {
			return err
		}

		return nil
	}}
}

