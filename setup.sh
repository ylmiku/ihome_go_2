#启动redisserver
redis-server ./conf/redis.conf

#启动fastdfs - tracker
fdfs_trackerd ./conf/tracker.conf restart
#启动fastdfs - storage
fdfs_storaged /home/itcast/workspace/go/src/ihome_go_2/conf/storage.conf restart
