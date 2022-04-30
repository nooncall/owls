package auth

type Auth interface {
	FilterCluster(username string)([]string, error)
	FilterDB(username, cluster string)([]string, error)
	GetManagers()[]string
}

