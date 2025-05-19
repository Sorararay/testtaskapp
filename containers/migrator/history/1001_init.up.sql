CREATE TABLE task
(
    id      CHAR(26)     NOT NULL, -- ULID 26bytes
    title   VARCHAR(191) NOT NULL, -- タイトル
    content TEXT         NOT NULL, -- 内容
    status  TEXT CHECK (status IN ('BACKLOG', 'PROGRESS', 'DONE')) NOT NULL, -- ステータス
    created TIMESTAMPTZ  NOT NULL, -- 作成時間
    updated TIMESTAMPTZ  NOT NULL, -- 更新時間
    PRIMARY KEY (id)
);