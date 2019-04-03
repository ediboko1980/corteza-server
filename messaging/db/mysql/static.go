// Code generated by statik. DO NOT EDIT.

// Package contains static assets.
package mysql

var Asset = "PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x1a\x00	\x0020180704080000.base.up.sqlUT\x05\x00\x01\x80Cm8-- Keeps all known channels\nCREATE TABLE channels (\n  id               BIGINT UNSIGNED NOT NULL,\n  name             TEXT            NOT NULL, -- display name of the channel\n  topic            TEXT            NOT NULL,\n  meta             JSON            NOT NULL,\n\n  type             ENUM ('private', 'public', 'group') NOT NULL DEFAULT 'public',\n\n  rel_organisation BIGINT UNSIGNED NOT NULL REFERENCES organisation(id),\n  rel_creator      BIGINT UNSIGNED NOT NULL,\n\n  created_at       DATETIME        NOT NULL DEFAULT NOW(),\n  updated_at       DATETIME            NULL,\n  archived_at      DATETIME            NULL,\n  deleted_at       DATETIME            NULL, -- channel soft delete\n\n  rel_last_message BIGINT UNSIGNED NOT NULL DEFAULT 0,\n\n  PRIMARY KEY (id)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;\n\n-- handles channel membership\nCREATE TABLE channel_members (\n  rel_channel      BIGINT UNSIGNED NOT NULL REFERENCES channels(id),\n  rel_user         BIGINT UNSIGNED NOT NULL,\n\n  type             ENUM ('owner', 'member', 'invitee') NOT NULL DEFAULT 'member',\n\n  created_at       DATETIME        NOT NULL DEFAULT NOW(),\n  updated_at       DATETIME            NULL,\n\n  PRIMARY KEY (rel_channel, rel_user)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;\n\nCREATE TABLE channel_views (\n  rel_channel      BIGINT UNSIGNED NOT NULL REFERENCES channels(id),\n  rel_user         BIGINT UNSIGNED NOT NULL,\n\n  -- timestamp of last view, should be enough to find out which messaghr\n  viewed_at        DATETIME        NOT NULL DEFAULT NOW(),\n\n  -- new messages count since last view\n  new_since        INT    UNSIGNED NOT NULL DEFAULT 0,\n\n  PRIMARY KEY (rel_user, rel_channel)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;\n\nCREATE TABLE channel_pins (\n  rel_channel      BIGINT UNSIGNED NOT NULL REFERENCES channels(id),\n  rel_message      BIGINT UNSIGNED NOT NULL REFERENCES messages(id),\n  rel_user         BIGINT UNSIGNED NOT NULL,\n\n  created_at       DATETIME        NOT NULL DEFAULT NOW(),\n\n  PRIMARY KEY (rel_channel, rel_message)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;\n\nCREATE TABLE messages (\n  id               BIGINT UNSIGNED NOT NULL,\n  type             TEXT,\n  message          TEXT            NOT NULL,\n  meta             JSON,\n  rel_user         BIGINT UNSIGNED NOT NULL,\n  rel_channel      BIGINT UNSIGNED NOT NULL REFERENCES channels(id),\n  reply_to         BIGINT UNSIGNED     NULL REFERENCES messages(id),\n\n  created_at       DATETIME        NOT NULL DEFAULT NOW(),\n  updated_at       DATETIME            NULL,\n  deleted_at       DATETIME            NULL,\n\n  PRIMARY KEY (id)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;\n\nCREATE TABLE reactions (\n  id               BIGINT UNSIGNED NOT NULL,\n  rel_user         BIGINT UNSIGNED NOT NULL,\n  rel_message      BIGINT UNSIGNED NOT NULL REFERENCES messages(id),\n  rel_channel      BIGINT UNSIGNED NOT NULL REFERENCES channels(id),\n  reaction         TEXT            NOT NULL,\n\n  created_at       DATETIME        NOT NULL DEFAULT NOW(),\n\n  PRIMARY KEY (id)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;\n\nCREATE TABLE attachments (\n  id               BIGINT UNSIGNED NOT NULL,\n  rel_user         BIGINT UNSIGNED NOT NULL,\n\n  url              VARCHAR(512),\n  preview_url      VARCHAR(512),\n\n  size             INT    UNSIGNED,\n  mimetype         VARCHAR(255),\n  name             TEXT,\n\n  meta             JSON,\n\n  created_at       DATETIME        NOT NULL DEFAULT NOW(),\n  updated_at       DATETIME            NULL,\n  deleted_at       DATETIME            NULL,\n\n  PRIMARY KEY (id)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;\n\nCREATE TABLE message_attachment (\n  rel_message      BIGINT UNSIGNED NOT NULL REFERENCES messages(id),\n  rel_attachment   BIGINT UNSIGNED NOT NULL REFERENCES attachment(id),\n\n  PRIMARY KEY (rel_message)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;\n\nCREATE TABLE event_queue (\n  id               BIGINT UNSIGNED NOT NULL,\n  origin           BIGINT UNSIGNED NOT NULL,\n  subscriber       TEXT,\n  payload          JSON,\n\n  PRIMARY KEY (id)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;\n\nCREATE TABLE event_queue_synced (\n  origin           BIGINT UNSIGNED NOT NULL,\n  rel_last         BIGINT UNSIGNED NOT NULL,\n\n  PRIMARY KEY (origin)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;\nPK\x07\x08\xd5\x9c\xef\x89V\x10\x00\x00V\x10\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00$\x00	\x0020181009080000.altering_types.up.sqlUT\x05\x00\x01\x80Cm8update channels set type = 'group' where type = 'direct';\nalter table channels CHANGE type type  enum('private', 'public', 'group');\nalter table channel_members CHANGE type type  enum('owner', 'member', 'invitee');\nPK\x07\x08E1\xf5\xa4\xd7\x00\x00\x00\xd7\x00\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00#\x00	\x0020181013080000.channel_views.up.sqlUT\x05\x00\x01\x80Cm8ALTER TABLE channel_views DROP viewed_at;\nALTER TABLE channel_views ADD rel_last_message_id BIGINT UNSIGNED;\nALTER TABLE channel_views CHANGE new_since new_messages_count INT UNSIGNED;\n\n-- Table structure after these changes:\n-- +---------------------+---------------------+------+-----+---------+-------+\n-- | Field               | Type                | Null | Key | Default | Extra |\n-- +---------------------+---------------------+------+-----+---------+-------+\n-- | rel_channel         | bigint(20) unsigned | NO   | PRI | NULL    |       |\n-- | rel_user            | bigint(20) unsigned | NO   | PRI | NULL    |       |\n-- | rel_last_message_id | bigint(20) unsigned | YES  |     | NULL    |       |\n-- | new_messages_count  | int(10) unsigned    | NO   |     | 0       |       |\n-- +---------------------+---------------------+------+-----+---------+-------+\n\n-- Prefill with data\nINSERT INTO channel_views (rel_channel, rel_user, rel_last_message_id)\n  SELECT cm.rel_channel, cm.rel_user, max(m.ID)\n    FROM channel_members AS cm INNER JOIN messages AS m ON (m.rel_channel = cm.rel_channel)\n  GROUP BY cm.rel_channel, cm.rel_user;\n\nPK\x07\x08`\xcbP\xf9t\x04\x00\x00t\x04\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x1d\x00	\x0020181013080000.replies.up.sqlUT\x05\x00\x01\x80Cm8ALTER TABLE messages CHANGE reply_to reply_to BIGINT UNSIGNED NOT NULL DEFAULT 0;\nALTER TABLE messages ADD replies INT UNSIGNED NOT NULL DEFAULT 0;\nPK\x07\x08m\xedWA\x94\x00\x00\x00\x94\x00\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00(\x00	\x0020181101080000.pins_and_reactions.up.sqlUT\x05\x00\x01\x80Cm8DROP TABLE channel_pins;\nDROP TABLE reactions;\n\nCREATE TABLE message_flags (\n  id               BIGINT UNSIGNED NOT NULL,\n  rel_channel      BIGINT UNSIGNED NOT NULL,\n  rel_message      BIGINT UNSIGNED NOT NULL,\n  rel_user         BIGINT UNSIGNED NOT NULL,\n  flag             TEXT,\n\n  created_at       DATETIME        NOT NULL DEFAULT NOW(),\n\n  PRIMARY KEY (id)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;\nPK\x07\x08eA\x1eo\x90\x01\x00\x00\x90\x01\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x1e\x00	\x0020181107080000.mentions.up.sqlUT\x05\x00\x01\x80Cm8CREATE TABLE mentions (\n  id               BIGINT UNSIGNED NOT NULL,\n  rel_channel      BIGINT UNSIGNED NOT NULL,\n  rel_message      BIGINT UNSIGNED NOT NULL,\n  rel_user         BIGINT UNSIGNED NOT NULL,\n  rel_mentioned_by BIGINT UNSIGNED NOT NULL,\n\n  created_at       DATETIME        NOT NULL DEFAULT NOW(),\n\n  PRIMARY KEY (id)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;\n\nCREATE INDEX lookup_mentions ON mentions (rel_mentioned_by)\nPK\x07\x08\xfb\xe8\x9b\x98\xac\x01\x00\x00\xac\x01\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x1d\x00	\x0020181115080000.unreads.up.sqlUT\x05\x00\x01\x80Cm8ALTER TABLE channel_views RENAME TO unreads;\n\nALTER TABLE unreads ADD     rel_reply_to                        BIGINT UNSIGNED NOT NULL AFTER rel_channel;\nALTER TABLE unreads CHANGE rel_channel         rel_channel      BIGINT UNSIGNED NOT NULL DEFAULT 0;\nALTER TABLE unreads CHANGE rel_user            rel_user         BIGINT UNSIGNED NOT NULL DEFAULT 0;\nALTER TABLE unreads CHANGE rel_last_message_id rel_last_message BIGINT UNSIGNED NOT NULL DEFAULT 0;\nALTER TABLE unreads CHANGE new_messages_count  count            INT    UNSIGNED NOT NULL DEFAULT 0;\n\nPK\x07\x08jf1Q+\x02\x00\x00+\x02\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00*\x00	\x0020181124173028.remove_events_tables.up.sqlUT\x05\x00\x01\x80Cm8DROP TABLE event_queue;\nDROP TABLE event_queue_synced;PK\x07\x08\xdd.y06\x00\x00\x006\x00\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00)\x00	\x0020181205153145.messages-to-utf8mb4.up.sqlUT\x05\x00\x01\x80Cm8alter table messages convert to character set utf8mb4 collate utf8mb4_unicode_ci;PK\x07\x08Ig\xbfOQ\x00\x00\x00Q\x00\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00&\x00	\x0020190122191150.membership-flags.up.sqlUT\x05\x00\x01\x80Cm8ALTER TABLE channel_members ADD flag ENUM ('pinned', 'hidden', 'ignored', '') NOT NULL DEFAULT '' AFTER `type`;\nPK\x07\x084\xfb\xe3\xf4p\x00\x00\x00p\x00\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00#\x00	\x0020190206112022.prefix-tables.up.sqlUT\x05\x00\x01\x80Cm8-- misc tables\n\nALTER TABLE attachments            RENAME TO messaging_attachment;\nALTER TABLE mentions               RENAME TO messaging_mention;\nALTER TABLE unreads                RENAME TO messaging_unread;\n\n-- channel tables\n\nALTER TABLE channels               RENAME TO messaging_channel;\nALTER TABLE channel_members        RENAME TO messaging_channel_member;\n\n-- message tables\n\nALTER TABLE messages               RENAME TO messaging_message;\nALTER TABLE message_attachment     RENAME TO messaging_message_attachment;\nALTER TABLE message_flags          RENAME TO messaging_message_flag;\nPK\x07\x08\x145\xde}Q\x02\x00\x00Q\x02\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0e\x00	\x00migrations.sqlUT\x05\x00\x01\x80Cm8CREATE TABLE IF NOT EXISTS `migrations` (\n `project` varchar(16) NOT NULL COMMENT 'sam, crm, ...',\n `filename` varchar(255) NOT NULL COMMENT 'yyyymmddHHMMSS.sql',\n `statement_index` int(11) NOT NULL COMMENT 'Statement number from SQL file',\n `status` TEXT NOT NULL COMMENT 'ok or full error message',\n PRIMARY KEY (`project`,`filename`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;\n\nPK\x07\x08\x0d\xa5T2x\x01\x00\x00x\x01\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x00\x00!(\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x06\x00	\x00new.shUT\x05\x00\x01\x80Cm8#!/bin/bash\ntouch $(date +%Y%m%d%H%M%S).up.sqlPK\x07\x08s\xd4N*.\x00\x00\x00.\x00\x00\x00PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(\xd5\x9c\xef\x89V\x10\x00\x00V\x10\x00\x00\x1a\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\x00\x00\x00\x0020180704080000.base.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(E1\xf5\xa4\xd7\x00\x00\x00\xd7\x00\x00\x00$\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\xa7\x10\x00\x0020181009080000.altering_types.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(`\xcbP\xf9t\x04\x00\x00t\x04\x00\x00#\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\xd9\x11\x00\x0020181013080000.channel_views.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(m\xedWA\x94\x00\x00\x00\x94\x00\x00\x00\x1d\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\xa7\x16\x00\x0020181013080000.replies.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(eA\x1eo\x90\x01\x00\x00\x90\x01\x00\x00(\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\x8f\x17\x00\x0020181101080000.pins_and_reactions.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(\xfb\xe8\x9b\x98\xac\x01\x00\x00\xac\x01\x00\x00\x1e\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81~\x19\x00\x0020181107080000.mentions.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(jf1Q+\x02\x00\x00+\x02\x00\x00\x1d\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\x7f\x1b\x00\x0020181115080000.unreads.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(\xdd.y06\x00\x00\x006\x00\x00\x00*\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\xfe\x1d\x00\x0020181124173028.remove_events_tables.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(Ig\xbfOQ\x00\x00\x00Q\x00\x00\x00)\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\x95\x1e\x00\x0020181205153145.messages-to-utf8mb4.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(4\xfb\xe3\xf4p\x00\x00\x00p\x00\x00\x00&\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81F\x1f\x00\x0020190122191150.membership-flags.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(\x145\xde}Q\x02\x00\x00Q\x02\x00\x00#\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\x13 \x00\x0020190206112022.prefix-tables.up.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(\x0d\xa5T2x\x01\x00\x00x\x01\x00\x00\x0e\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\xbe\"\x00\x00migrations.sqlUT\x05\x00\x01\x80Cm8PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x00\x00!(s\xd4N*.\x00\x00\x00.\x00\x00\x00\x06\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xed\x81{$\x00\x00new.shUT\x05\x00\x01\x80Cm8PK\x05\x06\x00\x00\x00\x00\x0d\x00\x0d\x00\\\x04\x00\x00\xe6$\x00\x00\x00\x00"
