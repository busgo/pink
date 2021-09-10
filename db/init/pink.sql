/*
 Navicat MySQL Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 50731
 Source Host           : localhost:3306
 Source Schema         : pink

 Target Server Type    : MySQL
 Target Server Version : 50731
 File Encoding         : 65001

 Date: 02/09/2021 22:06:29
*/


CREATE DATABASE IF NOT EXISTS pink default charset utf8mb4 COLLATE utf8mb4_general_ci;

use pink;

-- ----------------------------
-- Table structure for execute_snapshot_his
-- ----------------------------
DROP TABLE IF EXISTS `execute_snapshot_his`;
CREATE TABLE `execute_snapshot_his`
(
    `id`            bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
    `job_id`        varchar(32) NOT NULL,
    `snapshot_id`   varchar(32) NOT NULL,
    `job_name`      varchar(32) NOT NULL,
    `group`         varchar(32) NOT NULL,
    `cron`          varchar(255) DEFAULT NULL,
    `target`        varchar(255) DEFAULT NULL,
    `ip`            varchar(32)  DEFAULT NULL,
    `param`         varchar(255) DEFAULT NULL,
    `state`         tinyint(4) DEFAULT NULL,
    `before_time`   varchar(32) NOT NULL,
    `schedule_time` varchar(32) NOT NULL,
    `end_time`      varchar(32)  DEFAULT NULL,
    `times`         bigint(20) NOT NULL,
    `mobile`        varchar(32)  DEFAULT NULL,
    `remark`        varchar(32)  DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;