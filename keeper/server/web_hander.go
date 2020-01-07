package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) creatHander() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	r.POST("/login", s.login)
	r.GET("/check_health", s.checkHealth)
	r.GET("/proxy_list", s.proxyList)

	api := r.Group("/api")
	//api.Use(JWTAuth())
	{
		api.POST("/config/init", s.shardConfInit)
		api.GET("/config/full", s.apiConfig)
		api.GET("/config/basic", s.apiBasic)
		api.GET("/config/stop_service", s.apiStopService)
		api.GET("/config/switch", s.apiSwitch)
		api.GET("/config/hostgroups", s.apiHostGroups)
		api.GET("/config/host_group_clusters", s.apiHostGroupClusters)
		api.GET("/config/schema_configs", s.apiSchemaConfigs)
		api.GET("/config/tables", s.apiTableConfigs)
		api.GET("/config/export", s.exportConfig)

		api.GET("/proxy_cluster/status", s.clusterStatus)
		api.POST("/proxy_cluster/switch", s.switchDB)
		api.POST("/proxy/unregister", s.unregister)

		api.POST("/snapshot_config", s.snapshotConfig)
		api.GET("/snapshot_list", s.snapshotList)
		api.GET("/cluster_list", s.clusterList)

		api.POST("/update_basic", s.updateBasic)
		api.POST("/update_switch", s.updateSwitch)

		api.POST("/add_hostgroupcluster", s.AddHostGroupCluster)
		api.POST("/delete_hostgroupcluster", s.DelHostGroupCluster)
		api.POST("/update_hostgroupcluster", s.UpdateHostGroupCluster)
		api.POST("/add_hostgroups", s.addHostGroups)
		api.POST("/delete_hostgroups", s.delHostGroups)
		api.POST("/update_stopservice", s.updateStopService)
		api.POST("/add_schema", s.addSchema)
		api.POST("/delete_schema", s.deleteSchema)
		api.POST("/add_table", s.addTable)
		api.POST("/delete_table", s.deleteTable)
		api.POST("/update_shardrw", s.updateShardRW)
		api.POST("/push_config", s.pushConfig)
		api.POST("/rollback_config", s.rollbackConfig)
	}

	return r
}
