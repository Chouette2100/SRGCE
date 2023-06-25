--
-- Table structure for table `eventuser`
--

DROP TABLE IF EXISTS `weventuser`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `weventuser` (
  `eventid` varchar(128) NOT NULL,
  `userno` int NOT NULL,
  `istarget` char(1) DEFAULT NULL,
  `iscntrbpoints` char(1) DEFAULT 'N',
  `graph` char(1) DEFAULT NULL,
  `color` char(20) DEFAULT NULL,
  `point` int DEFAULT NULL,
  `vld` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`eventid`,`userno`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
