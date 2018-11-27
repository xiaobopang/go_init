CREATE TABLE `user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) DEFAULT '' COMMENT '账号',
  `age` int(11) NOT NULL DEFAULT '0' COMMENT '年龄',
  `email` varchar(64) NOT NULL DEFAULT '' COMMENT '邮箱',
  `gender` tinyint(4) NOT NULL DEFAULT '0' COMMENT '性别：1-男，2-女，0-未知',
  `is_deleted` tinyint(4) NOT NULL DEFAULT '0' COMMENT '删除：0-不删除，1-删除',
  `password` varchar(50) NOT NULL DEFAULT '' COMMENT '密码',
  `created_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;