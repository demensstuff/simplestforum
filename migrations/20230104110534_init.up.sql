-- users --
CREATE TABLE users
(
    id           BIGSERIAL   PRIMARY KEY,
    nickname     TEXT        UNIQUE NOT NULL,
    password     TEXT        NOT NULL,
    show_info    BOOLEAN     NOT NULL DEFAULT FALSE,
    rank         BIGINT      NOT NULL DEFAULT 1,
    level        TEXT,
    restriction  TEXT,
    count_posts  BIGINT      NOT NULL DEFAULT 0,
    count_topics BIGINT      NOT NULL DEFAULT 0,

    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

-- users_info --
CREATE TABLE users_info
(
    user_id    BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    phone      TEXT,
    email      TEXT,
    first_name TEXT,
    last_name  TEXT
);

INSERT INTO users (id, nickname, password, level) VALUES (1, 'admin', '$2a$10$R44LYcGDBbNkyqfDA1V4e.ykyeDMk28EHSN24oyEqZHSJYEMqs1zO', 'ADMIN');
INSERT INTO users_info (user_id, email) VALUES (1, 'admin@admin.admin');

-- sections --
CREATE TABLE sections
(
    id           BIGSERIAL  PRIMARY KEY,
    name         TEXT       NOT NULL,
    description  TEXT,
    count_topics BIGINT     NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

-- topics --
CREATE TABLE topics
(
    id          BIGSERIAL  PRIMARY KEY,
    section_id  BIGINT     NOT NULL REFERENCES sections (id) ON DELETE CASCADE,
    name        TEXT,
    user_id     BIGINT     NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    count_posts BIGINT     NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

-- posts --
CREATE TABLE posts
(
    id         BIGSERIAL  PRIMARY KEY,
    text       TEXT,
    user_id    BIGINT     NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    topic_id   BIGINT     NOT NULL REFERENCES topics (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- notifications --
CREATE TABLE notifications
(
    id         BIGSERIAL  PRIMARY KEY,
    user_id    BIGINT     NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    text       TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);