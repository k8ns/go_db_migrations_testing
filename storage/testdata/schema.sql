-- Code generated by MigrationCombiner. DO NOT EDIT.
-- Source: ../migrations
--
-- Generated by this command:
--
--  go run cmd/genmin/main.go -dir=../migrations -out=testdata/schema.sql
--
-- Generation timestamp: Fri, 25 Oct 2024 00:58:58 BST
--
-- 1_init.up.sql

CREATE TABLE articles (
  id int UNSIGNED NOT NULL AUTO_INCREMENT,
  title varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  created datetime DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


-- 2_authors.up.sql

CREATE TABLE authors (
  id int UNSIGNED NOT NULL AUTO_INCREMENT,
  name varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  created datetime DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

