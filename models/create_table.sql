CREATE TABLE `user`
(
    `id`          bigint(20)                             NOT NULL AUTO_INCREMENT,
    `user_id`     bigint(20)                             NOT NULL,
    `username`    varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
    `password`    varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
    `email`       varchar(64) COLLATE utf8mb4_general_ci,
    `gender`      tinyint(4)                             NOT NULL DEFAULT '0',
    `create_time` timestamp                              NULL     DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp                              NULL     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username` (`username`) USING BTREE,
    UNIQUE KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;


CREATE TABLE server_nodes
(
    `id`          bigint(20)  NOT NULL AUTO_INCREMENT,
    `name`        varchar(64) NOT NULL,
    `host`        varchar(64) NOT NULL,
    `port`        varchar(64) NOT NULL,
    `account`     varchar(64) NOT NULL,
    `password`    varchar(64) NOT NULL,
    `status`      boolean     NOT NULL,
    `remark`      varchar(64),
    `create_time` timestamp   NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp   NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;