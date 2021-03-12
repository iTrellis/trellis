
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
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


INSERT INTO `api` (`id`,`name`,`service_domain`,`service_name`,`service_version`,`topic`,`status`) VALUES ('1','trellis.test','trellis','serviceA','v1','test1','normal');
INSERT INTO `api` (`id`,`name`,`service_domain`,`service_name`,`service_version`,`topic`,`status`) VALUES ('2','trellis.test_remote','trellis','serviceA','v1','test_remote','normal');
INSERT INTO `api` (`id`,`name`,`service_domain`,`service_name`,`service_version`,`topic`,`status`) VALUES ('3','trellis.test_grpc','trellis','serviceA','v1','test_grpc','normal');
INSERT INTO `api` (`id`,`name`,`service_domain`,`service_name`,`service_version`,`topic`,`status`) VALUES ('4','trellis.testb','trellis','serviceB','v1','test','normal');
