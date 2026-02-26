CREATE TABLE `video` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `title` varchar(255) NOT NULL DEFAULT '' COMMENT '视频标题',
    `cover` varchar(255) NOT NULL DEFAULT '' COMMENT '封面图',
    `author_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '作者ID',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='视频表';

-- 预置一些视频数据
INSERT INTO `video` (`id`, `title`, `cover`, `author_id`) VALUES 
(1001, '周杰伦新歌首发，太好听了！', 'https://images.unsplash.com/photo-1470225620780-dba8ba36b745?w=200', 1),
(1002, '成都大熊猫基地一日游 🐼', 'https://images.unsplash.com/photo-1564349683136-77e08bef1ef1?w=200', 2),
(1003, '家常红烧肉保姆级教程，入口即化', 'https://images.unsplash.com/photo-1544025162-d76694265947?w=200', 3),
(1004, '这就是大自然的鬼斧神工吗？太震撼了', 'https://images.unsplash.com/photo-1464822759023-fed622ff2c3b?w=200', 1);

CREATE TABLE `comment` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '评论ID',
  `video_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '视频ID',
  `user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '用户ID',
  `content` varchar(1000) NOT NULL DEFAULT '' COMMENT '评论内容',
  `parent_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '父评论ID，0表示顶级评论',
  `root_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '根评论ID，用于回复链',
  `like_count` bigint(20) NOT NULL DEFAULT '0' COMMENT '点赞数',
  `reply_count` bigint(20) NOT NULL DEFAULT '0' COMMENT '回复数',
  `ip_location` varchar(50) NOT NULL DEFAULT '' COMMENT 'IP 属地',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除：0=否, 1=是',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_video_id` (`video_id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_root_id` (`root_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='评论表';

-- 用户表
CREATE TABLE `user` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `nickname` varchar(50) NOT NULL DEFAULT '' COMMENT '用户昵称',
    `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '头像地址',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 评论点赞表
CREATE TABLE `comment_like` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `comment_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '评论ID',
    `user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '用户ID',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_comment_user` (`comment_id`, `user_id`),
    INDEX `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='点赞记录表';
