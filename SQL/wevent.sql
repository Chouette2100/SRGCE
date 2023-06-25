--
-- Table structure for table `wevent`
--

DROP TABLE IF EXISTS `wevent`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `wevent` (
  `eventid` char(100) NOT NULL,
  `ieventid` int DEFAULT '0',
  `event_name` char(200) DEFAULT NULL,
  `period` char(100) DEFAULT NULL,
  `starttime` datetime DEFAULT NULL,
  `endtime` datetime DEFAULT NULL,
  `noentry` int DEFAULT '0',
  `intervalmin` int DEFAULT '5',
  `modmin` int DEFAULT '4',
  `modsec` int DEFAULT '10',
  `fromorder` int DEFAULT '1',
  `toorder` int DEFAULT '10',
  `resethh` int DEFAULT '4',
  `resetmm` int DEFAULT '0',
  `nobasis` int DEFAULT '999999',
  `maxdsp` int DEFAULT '10',
  `cmap` int DEFAULT '1',
  `target` int DEFAULT '0',
  `rstatus` varchar(20) DEFAULT '',
  `maxpoint` int DEFAULT '0',
  `achk` int DEFAULT '0',
  `aclr` int DEFAULT '0',
  PRIMARY KEY (`eventid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
