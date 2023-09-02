CREATE TABLE IF NOT EXISTS `user` (
    id BIGINT PRIMARY KEY,
    email VARCHAR(20) NOT NULL,
    password text NOT NULL,
    salt text NOT NULL,
    -- 二进制状态位，暂无作用
    state INT NOT NULL,
    reg_time BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS `skin` (
    id BIGINT PRIMARY KEY,
    -- 第一个上传的用户
    user_id BIGINT NOT NULL,
    skin_hash VARCHAR(50) NOT NULL,
    `type` VARCHAR(10) NOT NULL,
    variant VARCHAR(10) NOT NULL
);

CREATE TABLE IF NOT EXISTS `user_skin` (
    user_id BIGINT PRIMARY KEY,
    skin_id BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS `user_token` (
    user_id BIGINT PRIMARY KEY,
    token_id INT NOT NULL
);

CREATE TABLE IF NOT EXISTS `user_profile` (
    user_id BIGINT PRIMARY KEY,
    name VARCHAR(20) NOT NULL,
    uuid text NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS name_index ON user_profile (name);