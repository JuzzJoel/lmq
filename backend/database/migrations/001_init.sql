CREATE TABLE IF NOT EXISTS links (
    id          BIGSERIAL PRIMARY KEY,
    token       VARCHAR(20) UNIQUE NOT NULL,
    long_url    TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at  TIMESTAMPTZ,
    click_count BIGINT NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_links_token ON links(token);

CREATE TABLE IF NOT EXISTS click_events (
    id          BIGSERIAL PRIMARY KEY,
    link_id     BIGINT NOT NULL REFERENCES links(id) ON DELETE CASCADE,
    clicked_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ip_address  INET,
    city        VARCHAR(100) DEFAULT 'Unknown',
    region      VARCHAR(100) DEFAULT 'Unknown',
    country_code VARCHAR(2) DEFAULT 'XX',
    user_agent  TEXT,
    browser     VARCHAR(100) NOT NULL DEFAULT '',
    os          VARCHAR(100) NOT NULL DEFAULT '',
    is_mobile   BOOLEAN NOT NULL DEFAULT FALSE,
    referer     TEXT
);

CREATE INDEX IF NOT EXISTS idx_click_events_link_id ON click_events(link_id);
CREATE INDEX IF NOT EXISTS idx_click_events_clicked_at ON click_events(clicked_at);
CREATE INDEX IF NOT EXISTS idx_click_events_city ON click_events(city);
