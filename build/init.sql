SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `pass_wd` varchar(128) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `phone` varchar(11) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `status` tinyint(4) NULL DEFAULT 0,
  `desc` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `refresh_token` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `created_at` bigint(20) NOT NULL,
  `updated_at` bigint(20) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uix_user_name`(`name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

--  reco_mgt 推荐游戏管理表
CREATE TABLE IF NOT EXISTS `reco_mgt`  (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `kind` int(10) NOT NULL,
  `game_id` int(10) NOT NULL,
  `ad_slot_id` int(10) NOT NULL,
  `banner` text CHARACTER SET utf8 COLLATE utf8_general_ci,
  `comment` varchar(256) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `web_url` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `created_at` bigint(20) NOT NULL,
  `updated_at` bigint(20) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `kind`(`kind`) USING BTREE
)

--  ad_slot 推荐位管理表
CREATE TABLE IF NOT EXISTS `ad_slot`  (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `comment` varchar(256) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `created_at` bigint(20) NOT NULL,
  `updated_at` bigint(20) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `name`(`name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;


--  game_class 游戏分类管理表
CREATE TABLE IF NOT EXISTS `game_class`  (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `pid` int(10) NOT NULL,
  `kind` int(10) NOT NULL,
  `seq` int(10)  NOT NULL,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `comment` varchar(256) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `created_at` bigint(20) NOT NULL,
  `updated_at` bigint(20) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `kind`(`kind`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

--  game_pub 游戏管理发布表
CREATE TABLE IF NOT EXISTS `game_pub`  (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `game_id` varchar(60) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `icon` varchar(256) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `web_url` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `down_url` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `banner` varchar(256) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `hot_image` varchar(256) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `image` varchar(256) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `comment` varchar(256) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `status` tinyint(4) NULL DEFAULT 0,
  `game_dealer_id` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `game_flag` varchar(128) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `ratio` int(10) NOT NULL DEFAULT 0,
  `unit` varchar(12) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
  `main_type` int(10) NOT NULL DEFAULT 0,
  `sub_type` int(10)  NOT NULL DEFAULT 0,
  `main_label` int(10)   NOT NULL DEFAULT 0,
  `sub_label` varchar(256) CHARACTER SET utf8 COLLATE utf8_general_ci  NOT NULL DEFAULT 0,
  `seq` int(10)  NOT NULL DEFAULT 0,
  `created_at` bigint(20) NOT NULL,
  `updated_at` bigint(20) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uix_user_name`(`name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;
-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `users` VALUES (1, 'admin', '$2a$10$L5sYpJs2sZNxt8l0f8PSuute3mERbsOqBxVf4Xhocv1X7Ziu1.ZT6', '13288888888', 1, '超级管理员','', 1571360400, 1571360400);
-- ----------------------------
-- Table structure for roles
-- ----------------------------
DROP TABLE IF EXISTS `roles`;
CREATE TABLE `roles`  (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `desc` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `created_at` bigint(20) NOT NULL,
  `updated_at` bigint(20) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 10 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_role
-- ----------------------------
DROP TABLE IF EXISTS `user_role`;
CREATE TABLE `user_role`  (
  `user_id` int(10) NOT NULL,
  `role_id` int(10) NOT NULL,
  PRIMARY KEY (`user_id`, `role_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_role
-- ----------------------------
INSERT INTO `user_role` VALUES (1, 1);
-- ----------------------------

-- ----------------------------
-- Table structure for permissions
-- ----------------------------
DROP TABLE IF EXISTS `permissions`;
CREATE TABLE `permissions`  (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '权限id',
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '权限名(操作名)',
  `method` tinyint(4) NULL DEFAULT 0,
  `path` varchar(256) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '路径',
  `group_id` int(10) NOT NULL,
  `creat ed_at` bigint(20) NOT NULL,
  `updated_at` bigint(20) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of permissions
-- ----------------------------
INSERT INTO `permissions` VALUES (1, '新增用户',1, '/user',1, 1571360400, 1571360400);
INSERT INTO `permissions` VALUES (2, '删除用户',2, '/user',1, 1571360400, 1571360400);
INSERT INTO `permissions` VALUES (3, '修改用户资料',4, '/user',1, 1571360400, 1571360400);
INSERT INTO `permissions` VALUES (4, '查看用户列表',8, '/user',1, 1571360400, 1571360400);
INSERT INTO `permissions` VALUES (5, '新增角色',1, '/role',2, 1571360400, 1571360400);
INSERT INTO `permissions` VALUES (6, '删除角色',2, '/role',2, 1571360400, 1571360400);
INSERT INTO `permissions` VALUES (7, '修改角色资料',4, '/role',2, 1571360400, 1571360400);
INSERT INTO `permissions` VALUES (8, '查看角色列表',8, '/role',2, 1571360400, 1571360400);
-- ----------------------------
-- Table structure for role_perm
-- ----------------------------
DROP TABLE IF EXISTS `role_perm`;
CREATE TABLE `role_perm`  (
  `role_id` int(10) NOT NULL,
  `perm_id` int(10) NOT NULL,
  PRIMARY KEY (`role_id`, `perm_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic

DROP TABLE IF EXISTS `perm_group`;
CREATE TABLE `perm_group`  (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '权限组id',
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '权限组名称',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- Records of perm_group
-- ----------------------------
INSERT INTO `perm_group` VALUES (1, '用户管理');
INSERT INTO `perm_group` VALUES (2, '角色管理');


-- ----------------------------
-- Table structure for banner
-- ----------------------------
DROP TABLE IF EXISTS `banner`;
CREATE TABLE `banner`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `position` tinyint(4) NOT NULL,
  `seq` tinyint(4) NOT NULL,
  `icon` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `url` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `remark` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `status` tinyint(4) NOT NULL DEFAULT 1,
  `created_at` bigint(20) NULL DEFAULT NULL,
  `updated_at` bigint(20) NULL DEFAULT NULL,
  `kind` tinyint(4) NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for banner_publish
-- ----------------------------
DROP TABLE IF EXISTS `banner_publish`;
CREATE TABLE `banner_publish`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `banner_id` bigint(20) UNSIGNED NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for game_label
-- ----------------------------
DROP TABLE IF EXISTS `game_label`;
CREATE TABLE `game_label`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `pid` int(11) UNSIGNED NOT NULL,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `icon` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `comment` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `kind` tinyint(4) NOT NULL,
  `created_at` bigint(20) NULL DEFAULT NULL,
  `updated_at` bigint(20) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uix_label_name`(`name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;


-- ----------------------------
-- Table structure for game_pub
-- ----------------------------
DROP TABLE IF EXISTS `game_pub`;
CREATE TABLE `game_pub`  (
  `id` int(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `game_id` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `ad_slot` tinyint(4) NULL DEFAULT null,
  `main_type` tinyint(4) NULL DEFAULT NULL,
  `sub_type` tinyint(4) NULL DEFAULT NULL,
  `main_label` tinyint(4) NULL DEFAULT NULL,
  `access_type` tinyint(4) NULL DEFAULT NULL,
  `show_type` tinyint(4) NULL DEFAULT NULL,
  `sub_label` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `icon` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `banner` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `status` tinyint(4) NOT NULL,
  `seq` int(11) NOT NULL,
  `comment` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `web_url` varchar(128) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `down_url` varchar(128) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `created_at` bigint(20) NULL DEFAULT NULL,
  `updated_at` bigint(20) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for ad_slot
-- ----------------------------
DROP TABLE IF EXISTS `ad_slot`;
CREATE TABLE `ad_slot`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `comment` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `created_at` bigint(20) NULL DEFAULT NULL,
  `updated_at` bigint(20) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for reco_mgt
-- ----------------------------
DROP TABLE IF EXISTS `reco_mgt`;
CREATE TABLE `reco_mgt`  (
  `id` int(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `game_id` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `kind` tinyint(4) NOT NULL,
  `ad_slot_id` tinyint(4) NOT NULL,
  `banner` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `comment` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `web_url` varchar(128) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `created_at` bigint(20) NULL DEFAULT NULL,
  `updated_at` bigint(20) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

