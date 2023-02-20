-- create "users" table
CREATE TABLE `users` (`id` bigint NOT NULL AUTO_INCREMENT, `name` varchar(255) NOT NULL, `email` varchar(255) NOT NULL, `created_at` timestamp NOT NULL, PRIMARY KEY (`id`), UNIQUE INDEX `email` (`email`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- create "posts" table
CREATE TABLE `posts` (`id` bigint NOT NULL AUTO_INCREMENT, `title` varchar(255) NOT NULL, `body` longtext NOT NULL, `created_at` timestamp NOT NULL, `post_author` bigint NULL, PRIMARY KEY (`id`), INDEX `posts_users_author` (`post_author`), CONSTRAINT `posts_users_author` FOREIGN KEY (`post_author`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL) CHARSET utf8mb4 COLLATE utf8mb4_bin;
