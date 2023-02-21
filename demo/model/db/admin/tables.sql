CREATE TABLE `staff_info`
(
    `uid`   CHAR(6)    NOT NULL DEFAULT '' COMMENT '用户ID',
    `username` VARCHAR(11) NOT NULL DEFAULT '' COMMENT '登录名',
    `password` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '登录密码',
    `nickname` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '昵称',
    `avatar`   varchar(20) NOT NULL DEFAULT '0' COMMENT '头像',
    `email` varchar(100)    NOT NULL DEFAULT '' COMMENT '电子邮箱',
    `phone` varchar(20)    NOT NULL DEFAULT '' COMMENT '电话号码',
    `roles`  VARCHAR(100) NOT NULL DEFAULT '' COMMENT '权限角色',
    `status`   TINYINT(1) UNSIGNED NOT NULL DEFAULT '0' COMMENT '状态',
    `ct`       bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '注册时间',
    `ut`       bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '修改时间',
    PRIMARY KEY (`uid`),
    UNIQUE KEY `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户基本信息表';


CREATE TABLE `user_token`
(
    `uid` char(6)    NOT NULL DEFAULT '' COMMENT '用户ID',
    `token`  varchar(32) NOT NULL DEFAULT '' COMMENT 'token ID',
    `et`     bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '过期时间',
    `it`     bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '签发时间',
    `ut`     bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户token信息表';


CREATE TABLE `trace_log_2023` (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `module` varchar(50) NOT NULL DEFAULT '',
  `action` varchar(100) NOT NULL DEFAULT '',
  `roles` varchar(200) NOT NULL DEFAULT '',
  `params` text NOT NULL,
  `result` text NOT NULL,
  `uid` varchar(20) NOT NULL DEFAULT '',
  `ts` bigint(20) UNSIGNED NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;
