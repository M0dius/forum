-- MySQL Script generated by MySQL Workbench
-- Wed Nov 13 02:00:28 2024
-- Model: New Model    Version: 1.0
-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema Forum
-- -----------------------------------------------------
DROP SCHEMA IF EXISTS `Forum` ;

-- -----------------------------------------------------
-- Schema Forum
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `Forum` DEFAULT CHARACTER SET utf8 ;
SHOW WARNINGS;
USE `Forum` ;

-- -----------------------------------------------------
-- Table `categories`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `categories` ;

SHOW WARNINGS;
CREATE TABLE IF NOT EXISTS `categories` (
  `idcategories` INTEGER PRIMARY KEY AUTO_INCREMENT,
  `name` TEXT NOT NULL,
  `description` TEXT NULL
);

SHOW WARNINGS;

-- -----------------------------------------------------
-- Table `comment`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `comment` ;

SHOW WARNINGS;
CREATE TABLE IF NOT EXISTS `comment` (
  `commentid` INTEGER PRIMARY KEY AUTO_INCREMENT,
  `content` TEXT NULL,
  `comment_at` DATETIME NULL,
  `post_postid` INTEGER NOT NULL,
  `user_userid` INTEGER NOT NULL,
  FOREIGN KEY (`post_postid`) REFERENCES `post`(`postid`),
  FOREIGN KEY (`user_userid`) REFERENCES `user`(`userid`)
);

SHOW WARNINGS;

-- -----------------------------------------------------
-- Table `dislikes`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `dislikes` ;

SHOW WARNINGS;
CREATE TABLE IF NOT EXISTS `dislikes` (
  `dislikeid` INTEGER PRIMARY KEY AUTO_INCREMENT,
  `dislike_at` DATETIME NULL,
  `post_postid` INTEGER NOT NULL,
  `user_userid` INTEGER NOT NULL,
  FOREIGN KEY (`post_postid`) REFERENCES `post`(`postid`),
  FOREIGN KEY (`user_userid`) REFERENCES `user`(`userid`)
);

SHOW WARNINGS;

-- -----------------------------------------------------
-- Table `likes`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `likes` ;

SHOW WARNINGS;
CREATE TABLE IF NOT EXISTS `likes` (
  `likeid` INTEGER PRIMARY KEY AUTO_INCREMENT,
  `like_at` DATETIME NULL,
  `post_postid` INTEGER NOT NULL,
  `user_userid` INTEGER NOT NULL,
  FOREIGN KEY (`post_postid`) REFERENCES `post`(`postid`),
  FOREIGN KEY (`user_userid`) REFERENCES `user`(`userid`)
);

SHOW WARNINGS;

-- -----------------------------------------------------
-- Table `post`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `post` ;

SHOW WARNINGS;
CREATE TABLE IF NOT EXISTS `post` (
  `postid` INTEGER PRIMARY KEY AUTO_INCREMENT,
  `image` TEXT NULL,
  `content` TEXT NULL,
  `post_at` DATETIME NOT NULL,
  `user_userid` INTEGER NOT NULL,
  FOREIGN KEY (`user_userid`) REFERENCES `user`(`userid`)
);

SHOW WARNINGS;

-- -----------------------------------------------------
-- Table `post_has_categories`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `post_has_categories` ;

SHOW WARNINGS;
CREATE TABLE IF NOT EXISTS `post_has_categories` (
  `id` INTEGER PRIMARY KEY AUTO_INCREMENT,
  `post_postid` INTEGER NOT NULL,
  `categories_idcategories` INTEGER NOT NULL,
  FOREIGN KEY (`post_postid`) REFERENCES `post`(`postid`),
  FOREIGN KEY (`categories_idcategories`) REFERENCES `categories`(`idcategories`)
);

SHOW WARNINGS;

-- -----------------------------------------------------
-- Table `session`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `session` ;

SHOW WARNINGS;
CREATE TABLE IF NOT EXISTS `session` (
  `sessionid` INTEGER PRIMARY KEY AUTO_INCREMENT,
  `userid` INTEGER NOT NULL UNIQUE,
  `start` DATETIME NOT NULL,
  FOREIGN KEY (`userid`) REFERENCES `user`(`userid`)
);

SHOW WARNINGS;

-- -----------------------------------------------------
-- Table `user`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `user` ;

SHOW WARNINGS;
CREATE TABLE IF NOT EXISTS `user` (
<<<<<<< HEAD
  `userid` INT NOT NULL,
  `F-name`  NOT NULL,
  `L-name` VARCHAR(45) NOT NULL,
  `Username` VARCHAR(45) NOT NULL,
  `Email` VARCHAR(45) NOT NULL,
  `password` VARCHAR(45) NOT NULL,
  `current_session` INT NOT NULL,
  `role_id` INT NOT NULL,
  PRIMARY KEY (`userid`, `current_session`),
  FOREIGN KEY (`role_id`) REFERENCES `user_roles`(`roleid`))
ENGINE = InnoDB;

SHOW WARNINGS;
CREATE UNIQUE INDEX `Username_UNIQUE` ON `user` (`Username` ASC) VISIBLE;

SHOW WARNINGS;
CREATE UNIQUE INDEX `userid_UNIQUE` ON `user` (`userid` ASC) VISIBLE;
=======
  `userid` INTEGER PRIMARY KEY AUTO_INCREMENT,
  `F_name` TEXT NOT NULL,
  `L_name` TEXT NOT NULL,
  `Username` TEXT NOT NULL,
  `Email` TEXT NOT NULL,
  `password` TEXT NOT NULL,
  `session_sessionid` INTEGER NOT NULL,
  `role_id` INTEGER NOT NULL,
  `Avatar` TEXT,
  FOREIGN KEY (`session_sessionid`) REFERENCES `session`(`sessionid`),
  FOREIGN KEY (`role_id`) REFERENCES `user_roles`(`roleid`)
);
>>>>>>> main

SHOW WARNINGS;

-- -----------------------------------------------------
-- Table `user_roles`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `user_roles` ;

SHOW WARNINGS;
CREATE TABLE IF NOT EXISTS `user_roles` (
  `roleid` INTEGER PRIMARY KEY AUTO_INCREMENT,
  `role_name` TEXT NOT NULL
);

SHOW WARNINGS;
INSERT INTO `user_roles` (`roleid`, `role_name`) VALUES (1, 'admin'), (2, 'moderator'), (3, 'normal_user');

SHOW WARNINGS;

-- -----------------------------------------------------
-- Table `friends`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `friends` ;

SHOW WARNINGS;
CREATE TABLE IF NOT EXISTS `friends` (
  `id` INTEGER PRIMARY KEY AUTO_INCREMENT,
  `user_userid` INTEGER NOT NULL,
  `friend_userid` INTEGER NOT NULL,
  FOREIGN KEY (`user_userid`) REFERENCES `user`(`userid`),
  FOREIGN KEY (`friend_userid`) REFERENCES `user`(`userid`)
);

SHOW WARNINGS;

-- -----------------------------------------------------
-- Table `followers`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `followers` ;

SHOW WARNINGS;
CREATE TABLE IF NOT EXISTS `followers` (
  `id` INTEGER PRIMARY KEY AUTO_INCREMENT,
  `user_userid` INTEGER NOT NULL,
  `follower_userid` INTEGER NOT NULL,
  FOREIGN KEY (`user_userid`) REFERENCES `user`(`userid`),
  FOREIGN KEY (`follower_userid`) REFERENCES `user`(`userid`)
);

SHOW WARNINGS;

-- -----------------------------------------------------
-- Table `notifications`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `notifications` ;

SHOW WARNINGS;
CREATE TABLE IF NOT EXISTS `notifications` (
  `notificationid` INTEGER PRIMARY KEY AUTO_INCREMENT,
  `user_userid` INTEGER NOT NULL,
  `post_id` INTEGER NOT NULL,
  `message` TEXT NOT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`user_userid`) REFERENCES `user`(`userid`),
  FOREIGN KEY (`post_id`) REFERENCES `post`(`postid`)
);

SHOW WARNINGS;

-- -----------------------------------------------------
-- Table `following`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `following` ;

SHOW WARNINGS;
CREATE TABLE IF NOT EXISTS `following` (
  `id` INTEGER PRIMARY KEY AUTO_INCREMENT,
  `user_userid` INTEGER NOT NULL,
  `following_userid` INTEGER NOT NULL,
  FOREIGN KEY (`user_userid`) REFERENCES `user`(`userid`),
  FOREIGN KEY (`following_userid`) REFERENCES `user`(`userid`)
);

SHOW WARNINGS;

-- -----------------------------------------------------
-- Table `reports`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `reports` ;

SHOW WARNINGS;
CREATE TABLE IF NOT EXISTS `reports` (
  `id` INTEGER PRIMARY KEY AUTO_INCREMENT,
  `post_id` INTEGER NOT NULL,
  `reported_by` INTEGER NOT NULL,
  `report_reason` TEXT,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`post_id`) REFERENCES `post`(`postid`),
  FOREIGN KEY (`reported_by`) REFERENCES `user`(`userid`)
);

SHOW WARNINGS;

SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;