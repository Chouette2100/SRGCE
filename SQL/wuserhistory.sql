DROP TABLE IF EXISTS `wuserhistory`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `wuserhistory` (
  `userno` int NOT NULL,
  `user_name` char(200) DEFAULT NULL,
  `genre` char(100) DEFAULT NULL,
  `rank` char(16) DEFAULT NULL,
  `nrank` varchar(60) DEFAULT '-',
  `prank` varchar(60) DEFAULT '-',
  `level` int DEFAULT NULL,
  `followers` int DEFAULT NULL,
  `fans` int DEFAULT '-1',
  `fans_lst` int DEFAULT '-1',
  `ts` datetime NOT NULL,
  PRIMARY KEY (`userno`,`ts`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
