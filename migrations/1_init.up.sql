
CREATE TABLE articles (
  id int UNSIGNED NOT NULL AUTO_INCREMENT,
  title varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  created datetime DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
