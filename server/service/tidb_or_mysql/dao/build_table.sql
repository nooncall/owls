DROP TABLE IF EXISTS `owl_task`;
CREATE TABLE IF NOT EXISTS `owl_task`
(
    `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `name`           VARCHAR(150)     NOT NULL DEFAULT '' COMMENT '名称',
    `status`         VARCHAR(40)     NOT NULL DEFAULT '' COMMENT '状态',
    `creator`        VARCHAR(60)     NOT NULL DEFAULT '' COMMENT '创建者',
    `reviewer`       VARCHAR(256)     NOT NULL DEFAULT '' COMMENT '审查者',
    `executor`       VARCHAR(60)     NOT NULL DEFAULT '' COMMENT '执行者',
    `exec_info`      VARCHAR(300)     NOT NULL DEFAULT '' COMMENT '执行信息',
    `reject_content` varchar(500)    NOT NULL DEFAULT '' COMMENT '驳回信息',

    `ct`             BIGINT          NOT NULL DEFAULT 0 COMMENT '创建时间',
    `ut`             BIGINT          NOT NULL DEFAULT 0 COMMENT '更改时间',
    `et`             BIGINT          NOT NULL DEFAULT 0 COMMENT '执行时间',
    `ft`             BIGINT          NOT NULL DEFAULT 0 COMMENT '执行结束时间',

    PRIMARY KEY (`id`),
    KEY idx_csc (`creator`, `status`, `ct`),
    KEY idx_rsc (`reviewer`, `status`, `ct`),
    KEY idx_sc (`status`, `ct`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin
  AUTO_INCREMENT = 7077 COMMENT '任务表';

DROP TABLE IF EXISTS `owl_subtask`;
CREATE TABLE IF NOT EXISTS `owl_subtask`
(
    `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `task_id`      BIGINT UNSIGNED NOT NULL COMMENT 'task id',
    `task_type`    VARCHAR(20)     NOT NULL DEFAULT '' COMMENT '任务类型',
    `db_name`      VARCHAR(120)    NOT NULL DEFAULT '' COMMENT '数据库名称',
    `cluster_name` VARCHAR(120)    NOT NULL DEFAULT '' COMMENT '数据库集群',

    PRIMARY KEY (`id`),
    KEY t_id (`task_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin
  AUTO_INCREMENT = 7077 COMMENT '子任务表';

DROP TABLE IF EXISTS `owl_exec_item`;
CREATE TABLE IF NOT EXISTS `owl_exec_item`
(
    `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `task_id`       BIGINT UNSIGNED NOT NULL COMMENT 'task id',
    `subtask_id`    BIGINT UNSIGNED NOT NULL COMMENT 'subtask id',
    `sql_content`   TEXT COMMENT 'sql',
    `remark`        VARCHAR(400)    NOT NULL DEFAULT '' COMMENT '备注',
    `affect_rows`   INT             NOT NULL DEFAULT 0 COMMENT '影响行数',
    `rule_comments` VARCHAR(300)     NOT NULL DEFAULT '' COMMENT '规则审核结果',
    `status`        VARCHAR(40)     NOT NULL DEFAULT '' COMMENT '子项状态',
    `exec_info`     VARCHAR(200)    NOT NULL DEFAULT '' COMMENT '执行信息',
    `backup_status` VARCHAR(40)     NOT NULL DEFAULT '' COMMENT '子项备份状态',
    `backup_info`   VARCHAR(200)    NOT NULL DEFAULT '' COMMENT '备份异常信息',
    `backup_id`     BIGINT UNSIGNED NOT NULL COMMENT 'backup id',

    `ut`            BIGINT          NOT NULL DEFAULT 0 COMMENT '更改时间',
    `et`            BIGINT          NOT NULL DEFAULT 0 COMMENT '执行时间',

    PRIMARY KEY (`id`),
    KEY idx_t_s_id (`task_id`, `subtask_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin
  AUTO_INCREMENT = 7077 COMMENT '子任务里的一条sql，执行单位';

DROP TABLE IF EXISTS `owl_rule_status`;
CREATE TABLE IF NOT EXISTS `owl_rule_status`
(
    `id`      BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `name`    VARCHAR(120)     NOT NULL DEFAULT '' COMMENT '规则名',
    `close`  TINYINT         NOT NULL DEFAULT 0 COMMENT '开关',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin
  AUTO_INCREMENT = 7077 COMMENT '规则状态任务表';

DROP TABLE IF EXISTS `owl_cluster`;
CREATE TABLE IF NOT EXISTS `owl_cluster`
(
    `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `name`        VARCHAR(160)    NOT NULL DEFAULT '' COMMENT '名称',
    `description` VARCHAR(200)    NOT NULL DEFAULT '' COMMENT '描述',
    `addr`        VARCHAR(120)     NOT NULL DEFAULT '' COMMENT '访问地址',
    `user`        VARCHAR(120)    NOT NULL DEFAULT '' COMMENT '用户',
    `pwd`         VARCHAR(200)    NOT NULL DEFAULT '' COMMENT '加密后秘钥',

    `ct`          BIGINT          NOT NULL DEFAULT 0 COMMENT '创建时间',
    `ut`          BIGINT          NOT NULL DEFAULT 0 COMMENT '更改时间',
    `operator`    VARCHAR(60)     NOT NULL DEFAULT '' COMMENT '操作人',

    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin
  AUTO_INCREMENT = 7077 COMMENT '集群信息表';

DROP TABLE IF EXISTS `owl_backup`;
CREATE TABLE IF NOT EXISTS `owl_backup`
(
    `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `data`           MEDIUMTEXT      NOT NULL COMMENT '备份数据',
    `rollback_user`  VARCHAR(60)     NOT NULL DEFAULT '' COMMENT '恢复执行人',
    `is_rolled_back` TINYINT         NOT NULL DEFAULT 0 COMMENT '是否已恢复',

    `ct`             BIGINT          NOT NULL DEFAULT 0 COMMENT '创建时间',
    `rollback_time`  BIGINT          NOT NULL DEFAULT 0 COMMENT '回滚时间',

    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin
  AUTO_INCREMENT = 7077 COMMENT '备份信息表';

DROP TABLE IF EXISTS `owl_admin`;
CREATE TABLE IF NOT EXISTS `owl_admin`
(
    `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `username`       VARCHAR(40)     NOT NULL DEFAULT '' COMMENT 'ldap账号',
    `description`    VARCHAR(200)    NOT NULL DEFAULT '' COMMENT '描述',

    `ct`             BIGINT          NOT NULL DEFAULT 0 COMMENT '创建时间',
    `creator`        VARCHAR(40)     NOT NULL DEFAULT '' COMMENT '操作人',

PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin
  AUTO_INCREMENT = 7077 COMMENT '管理员表';