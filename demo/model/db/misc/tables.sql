CREATE TABLE `user_info`
(
    `userid`   CHAR(10)    NOT NULL DEFAULT '' COMMENT '用户ID',
    `username` VARCHAR(11) NOT NULL DEFAULT '' COMMENT '登录名',
    `password` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '登录密码',
    `nid`      VARCHAR(10) NOT NULL DEFAULT '' COMMENT '靓号',
    `nickname` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '昵称',
    `sex`      TINYINT(1) UNSIGNED NOT NULL DEFAULT '0' COMMENT '性别',
    `avatar`   varchar(20) NOT NULL DEFAULT '0' COMMENT '头像',
    `birthday` char(10)    NOT NULL DEFAULT '' COMMENT '出生日期',
    `inviter`  VARCHAR(20) NOT NULL DEFAULT '' COMMENT '邀请人',
    `status`   TINYINT(1) UNSIGNED NOT NULL DEFAULT '0' COMMENT '状态',
    `ct`       bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '注册时间',
    `ut`       bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '修改时间',
    PRIMARY KEY (`userid`),
    UNIQUE KEY `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户基本信息表';


CREATE TABLE `user_token`
(
    `userid` char(10)    NOT NULL DEFAULT '' COMMENT '用户ID',
    `token`  varchar(32) NOT NULL DEFAULT '' COMMENT 'token ID',
    `et`     bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '过期时间',
    `it`     bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '签发时间',
    PRIMARY KEY (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户token信息表';