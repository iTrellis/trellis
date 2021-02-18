
/*
-- Query: SELECT * FROM apis.api
-- Date: 2020-09-15 11:59
*/

CREATE TABLE `api` (
  `id` varchar(50) NOT NULL DEFAULT '',
  `name` varchar(200) NOT NULL DEFAULT '',
  `service_domain` varchar(100) NOT NULL DEFAULT '',
  `service_name` varchar(100) NOT NULL DEFAULT '',
  `service_version` varchar(50) NOT NULL DEFAULT '',
  `topic` varchar(100) NOT NULL DEFAULT '',
  `status` varchar(50) NOT NULL DEFAULT 'normal',
  `version` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


INSERT INTO `api` (`id`,`name`,`service_domain`,`service_name`,`service_version`,`topic`,`status`,`version`) VALUES ('1','trellis.test','trellis','serviceA','v1','test1','normal',0);
INSERT INTO `api` (`id`,`name`,`service_domain`,`service_name`,`service_version`,`topic`,`status`,`version`) VALUES ('2','trellis.test_remote','trellis','serviceA','v1','test_remote','normal',0);
INSERT INTO `api` (`id`,`name`,`service_domain`,`service_name`,`service_version`,`topic`,`status`,`version`) VALUES ('3','trellis.test_grpc','trellis','serviceA','v1','test_grpc','normal',0);
INSERT INTO `api` (`id`,`name`,`service_domain`,`service_name`,`service_version`,`topic`,`status`,`version`) VALUES ('4','trellis.testb','trellis','serviceB','v1','test','normal',0);
