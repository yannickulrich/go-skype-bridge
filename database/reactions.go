package database

import (
	"database/sql"
	"fmt"
	log "maunium.net/go/maulogger/v2"

	"github.com/kelaresg/matrix-skype/types"
	"maunium.net/go/mautrix/id"
)

type ReactionQuery struct {
	db  *Database
	log log.Logger
}

func (rq *ReactionQuery) New() *Reaction {
	return &Reaction{
		db:  rq.db,
		log: rq.log,
	}
}

func (rq *ReactionQuery) GetForMID(mid string) (messages map[types.SkypeReactionID]*Reaction) {
	rows, err := rq.db.Query("SELECT id, chat_jid, chat_receiver, message_id, mxid, sender, timestamp, content FROM reactions WHERE message_id=$1", mid)
	if err != nil || rows == nil {
		return nil
	}
	defer rows.Close()
	messages = make(map[types.SkypeReactionID]*Reaction)
	for rows.Next() {
		r := rq.New().Scan(rows)
		messages[r.ID] = r
	}
	return
}

type Reaction struct {
	db  *Database
	log log.Logger

	ID        types.SkypeReactionID
	Chat      PortalKey
	MID       types.SkypeMessageID
	MXID      id.EventID
	Sender    types.SkypeID
	Timestamp uint64
	Content   string
}

func (r *Reaction) Scan(row Scannable) *Reaction {
	var content []byte
	err := row.Scan(&r.ID, &r.Chat.JID, &r.Chat.Receiver, &r.MID, &r.MXID, &r.Sender, &r.Timestamp, &content)
	if err != nil {
		if err != sql.ErrNoRows {
			r.log.Errorln("Database scan failed:", err)
		}
		return nil
	}

	r.Content = string(content)

	return r
}

func (r *Reaction) Insert() {
	fmt.Println("calling insert", r)
	_, err := r.db.Exec("INSERT INTO reactions (id, chat_jid, chat_receiver, message_id, mxid, sender, timestamp, content) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		r.ID, r.Chat.JID, r.Chat.Receiver, r.MID, r.MXID, r.Sender, r.Timestamp, r.Content)
	if err != nil {
		r.log.Warnfln("Failed to insert %s@%s: %v", r.Chat, r.MID, err)
	}
}

func (r *Reaction) Delete() {
	fmt.Println("calling delete", r)
	_, err := r.db.Exec("DELETE FROM reactions WHERE chat_jid=$1 AND chat_receiver=$2 AND message_id=$3", r.Chat.JID, r.Chat.Receiver, r.MID)
	if err != nil {
		r.log.Warnfln("Failed to delete %s@%s: %v", r.Chat, r.MID, err)
	}
}
