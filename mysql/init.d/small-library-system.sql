SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- データベース: `small_library_system`
--

-- --------------------------------------------------------

--
-- テーブルの構造 `rent_history`
--
CREATE TABLE `rent_history` (
  `id` int NOT NULL,
  `collection_list_id` int NOT NULL,
  `rent_by` varchar(200) COLLATE utf8mb4_general_ci NOT NULL,
  `start_date` date NOT NULL,
  `estimate_end_date` date NOT NULL,
  `real_end_date` date DEFAULT NULL,
  `last_modified_date` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
--
-- テーブルのインデックス `rent_history`
--
ALTER TABLE `rent_history`
  ADD PRIMARY KEY (`id`);
--
-- テーブルのAUTO_INCREMENT `rent_history`
--
ALTER TABLE `rent_history`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;


-- --------------------------------------------------------

--
-- テーブルの構造 `collection_list`
--
CREATE TABLE `collection_list` (
  `id` int NOT NULL, -- 連番
  `qr_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `created_by` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `created_at` varchar(200) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `note` text COLLATE utf8mb4_general_ci,

  `rent_by` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `rent_start_date` date DEFAULT NULL,
  `estimate_rent_end_date` date DEFAULT NULL,

  `is_rent` tinyint(1) NOT NULL DEFAULT '0',
  `is_registered` tinyint(1) NOT NULL DEFAULT '0',
  `is_delete` tinyint(1) NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
--
-- テーブルのインデックス `collection_list`
--
ALTER TABLE `collection_list`
  ADD PRIMARY KEY (`id`);
--
-- テーブルのAUTO_INCREMENT `collection_list`
--
ALTER TABLE `collection_list`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- テーブルの構造 `tag_list`
--
CREATE TABLE `tag_list` (
  `id` int NOT NULL,
  `tag_name` varchar(200) COLLATE utf8mb4_general_ci NOT NULL,
  `genre` tinyint NOT NULL DEFAULT '0',
  `isDelete` tinyint(1) NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
--
-- テーブルのインデックス `tag_list`
--
ALTER TABLE `tag_list`
  ADD PRIMARY KEY (`id`);
--
-- テーブルのAUTO_INCREMENT `tag_list`
--
ALTER TABLE `tag_list`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

-- --------------------------------------------------------

--
-- テーブルの構造 `tag_tagging`
--

CREATE TABLE `tag_tagging` (
  `id` int NOT NULL,
  `collection_list_id` int NOT NULL,
  `tag_id` int NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
--
-- テーブルのインデックス `tag_tagging`
--
ALTER TABLE `tag_tagging`
  ADD PRIMARY KEY (`id`);
--
-- テーブルのAUTO_INCREMENT `tag_tagging`
--
ALTER TABLE `tag_tagging`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;


COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
