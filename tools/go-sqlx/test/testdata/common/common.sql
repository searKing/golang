/*
 * Copyright 2020 The searKing Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

CREATE TABLE IF NOT EXISTS common
(
    # 必备字段
    id         BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    created_at DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最近更新时间',
    # 软删除
    is_deleted TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'soft delete, 0 for not deleted, 1 for deleted',
    deleted_at DATETIME                  DEFAULT NULL COMMENT '最近删除时间',
    # 版本控制
    version    int              NOT NULL DEFAULT 0,

    PRIMARY KEY pk_id (id)
) DEFAULT CHARSET = utf8 COMMENT ='common table';