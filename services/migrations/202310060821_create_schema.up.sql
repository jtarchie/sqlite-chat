CREATE TABLE users(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  message TEXT NULL,
  email_address TEXT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ensure created_at and updated_at (TRIGGER) for all
CREATE TRIGGER update_users_timestamp BEFORE
UPDATE
  ON users FOR EACH ROW BEGIN
UPDATE
  users
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE channels(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  creator_user_id INTEGER NOT NULL,
  description TEXT,
  private INTEGER NOT NULL DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(creator_user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- ensure created_at and updated_at (TRIGGER) for all
CREATE TRIGGER update_channels_timestamp BEFORE
UPDATE
  ON channels FOR EACH ROW BEGIN
UPDATE
  channels
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE channel_members(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  channel_id INTEGER NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(channel_id) REFERENCES channels(id) ON DELETE CASCADE,
  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- ensure created_at and updated_at (TRIGGER) for all
CREATE TRIGGER update_channel_members_timestamp BEFORE
UPDATE
  ON channel_members FOR EACH ROW BEGIN
UPDATE
  channel_members
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE messages(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  channel_id INTEGER NOT NULL,
  parent_id INTEGER,
  copy TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(channel_id) REFERENCES channels(id) ON DELETE CASCADE,
  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY(parent_id) REFERENCES messages(id) ON DELETE CASCADE
);

-- ensure created_at and updated_at (TRIGGER) for all
CREATE TRIGGER update_messages_timestamp BEFORE
UPDATE
  ON messages FOR EACH ROW BEGIN
UPDATE
  messages
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

-- initalize the chat
INSERT INTO
  users (name, message, email_address)
VALUES
  (
    "bot",
    "I am the default bot of the chat.",
    "bot@example.com"
  );

INSERT INTO
  channels (
    name,
    creator_user_id,
    description,
    private
  )
VALUES
  ("general", 1, "", 0);

INSERT INTO
  channel_members (user_id, channel_id)
VALUES
  (1, 1);

INSERT INTO
  messages (user_id, channel_id, copy)
VALUES
  (1, 1, "Welcome to the chat!");

CREATE VIEW user_channels AS
SELECT
  u.id AS user_id,
  u.name AS user_name,
  u.email_address,
  c.id AS channel_id,
  c.name AS channel_name,
  c.description,
  c.private
FROM
  users u
  JOIN channel_members cm ON u.id = cm.user_id
  JOIN channels c ON cm.channel_id = c.id;